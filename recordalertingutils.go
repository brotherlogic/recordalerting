package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gd "github.com/brotherlogic/godiscogs"
	pb "github.com/brotherlogic/recordalerting/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
)

func (s *Server) adjustState(ctx context.Context, config *pb.Config, r *pbrc.Record, needs bool, class pb.Problem_ProblemType, errorMessage string) error {
	// Does this record need a weight
	alreadySeen := false
	var number int32
	for _, problem := range config.GetProblems() {
		if problem.GetType() == class && problem.GetInstanceId() == r.GetRelease().GetInstanceId() {
			alreadySeen = true
			number = problem.GetIssueNumber()
		}
	}
	if needs && !alreadySeen {
		issue, err := s.ImmediateIssue(ctx, fmt.Sprintf("%v [%v] %v", r.GetRelease().GetTitle(), r.GetRelease().GetInstanceId(), errorMessage), fmt.Sprintf("This one [%v]: https://www.discogs.com/madeup/release/%v", r.GetRelease().GetInstanceId(), r.GetRelease().GetId()))
		if err != nil {
			return err
		}

		config.Problems = append(config.Problems, &pb.Problem{
			Type:        class,
			IssueNumber: issue.GetNumber(),
			InstanceId:  r.GetRelease().GetInstanceId()})
		err = s.saveConfig(ctx, config)
		if err != nil {
			return err
		}
	} else if !needs && alreadySeen {
		err := s.DeleteIssue(ctx, number)
		s.Log(fmt.Sprintf("Deleting the record %v for %v -> %v", number, r.GetRelease().GetInstanceId(), err))

		if err != nil {
			return err
		}

		var problems []*pb.Problem
		for _, p := range config.GetProblems() {
			if p.GetInstanceId() != r.GetRelease().GetInstanceId() || p.GetType() != class {
				problems = append(problems, p)
			}
		}
		config.Problems = problems
		err = s.saveConfig(ctx, config)
		if err != nil {
			return err
		}
	}

	if needs && (class == pb.Problem_MISSING_WEIGHT || class == pb.Problem_MISSING_WIDTH || class == pb.Problem_MISSING_FILED) {
		return status.Errorf(codes.FailedPrecondition, "Record %v fails validation - please fix", r.GetRelease().GetInstanceId())
	}
	return nil
}

func (s *Server) needsWeight(ctx context.Context, config *pb.Config, r *pbrc.Record) error {
	return s.adjustState(ctx, config, r,
		r.GetMetadata().GetMoveFolder() == 488127 && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_DIGITAL && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_UNKNOWN && r.GetMetadata().GetWeightInGrams() == 0,
		pb.Problem_MISSING_WEIGHT, "needs weight")
}
func (s *Server) needsWidth(ctx context.Context, config *pb.Config, r *pbrc.Record) error {
	return s.adjustState(ctx, config, r,
		r.GetMetadata().GetMoveFolder() == 812802 && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_DIGITAL && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_UNKNOWN && r.GetMetadata().GetRecordWidth() == 0,
		pb.Problem_MISSING_WIDTH, "needs width")
}
func (s *Server) needsFiled(ctx context.Context, config *pb.Config, r *pbrc.Record) error {
	return s.adjustState(ctx, config, r,
		r.GetMetadata().GetFiledUnder() == pbrc.ReleaseMetadata_FILE_UNKNOWN,
		pb.Problem_MISSING_FILED, "needs filling")
}
func (s *Server) needsCondition(ctx context.Context, config *pb.Config, r *pbrc.Record) error {
	return s.adjustState(ctx, config, r,
		r.GetMetadata().GetDateArrived() > 0 && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_UNKNOWN && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_DIGITAL && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_UNKNOWN && r.GetRelease().GetRecordCondition() == "",
		pb.Problem_MISSING_CONDITION, "needs condition")
}

func (s *Server) assessRecord(ctx context.Context, config *pb.Config, r *pbrc.Record) error {
	// We don't alert on boxed records
	if r.GetMetadata().GetBoxState() != pbrc.ReleaseMetadata_OUT_OF_BOX && r.GetMetadata().GetBoxState() != pbrc.ReleaseMetadata_BOX_UNKNOWN {
		return nil
	}

	err1 := s.needsFiled(ctx, config, r)
	err2 := s.needsWeight(ctx, config, r)
	err3 := s.needsWidth(ctx, config, r)
	err4 := s.needsCondition(ctx, config, r)

	// Only fail
	if r.GetMetadata().GetCategory() != pbrc.ReleaseMetadata_UNKNOWN {
		if (r.GetMetadata().GetFiledUnder() == pbrc.ReleaseMetadata_FILE_12_INCH || r.GetMetadata().GetFiledUnder() == pbrc.ReleaseMetadata_FILE_7_INCH) && r.GetMetadata().GetSaleState() != gd.SaleState_SOLD {
			if time.Since(time.Unix(r.GetMetadata().GetLastCleanDate(), 0)) > time.Hour*24*365*3 {
				err := s.rc.clean(ctx, r.GetRelease().GetInstanceId())
				if err != nil {
					return err
				}
			}
		}
	}

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	if err3 != nil {
		return err3
	}
	if err4 != nil {
		return err4
	}

	s.validateRecord(r)
	s.alertForMissingSaleID(r)
	s.alertForPurgatory(r)

	return nil
}

func (s *Server) validateRecord(r *pbrc.Record) {
	if r.GetMetadata().GetGoalFolder() == 0 {
		s.alertCount++
		s.RaiseIssue(fmt.Sprintf("%v Invalid Record", r.GetRelease().GetInstanceId()), fmt.Sprintf("%v has no goal folder", r.GetRelease().GetInstanceId()))
		s.invalidRecords++
	}

	if r.GetMetadata().GetCategory() == pbrc.ReleaseMetadata_STAGED_TO_SELL && r.GetRelease().GetRating() == 4 {
		s.alertCount++
		s.RaiseIssue(fmt.Sprintf("%v Needs your Attention", r.GetRelease().GetInstanceId()), fmt.Sprintf("%v is stuck staged to sell", r.GetRelease().GetTitle()))
	}

	if r.GetMetadata().GetGoalFolder() == 1613206 && r.GetMetadata().GetCategory() != pbrc.ReleaseMetadata_SOLD_ARCHIVE {
		s.RaiseIssue(fmt.Sprintf("%v is marked as sold, but is not SOLD_ARCHIVE", r.GetRelease().GetInstanceId()), fmt.Sprintf("It's %v ->  https://www.discogs.com/madeup/release/%v", r.GetMetadata().GetCategory(), r.GetRelease().GetId()))
	}

	if r.GetMetadata().GetCategory() == pbrc.ReleaseMetadata_PURCHASED && time.Now().Sub(time.Unix(r.GetMetadata().GetLastUpdateTime(), 0)) > time.Hour*24 {
		s.alertCount++
		s.RaiseIssue(fmt.Sprintf("%v Stale Purchase", r.GetRelease().GetInstanceId()), fmt.Sprintf("%v has staled", r.GetRelease().GetInstanceId()))
	}
}

func (s *Server) alertForMissingSaleID(r *pbrc.Record) {
	if r.GetMetadata().GetSaleId() <= 0 && r.GetMetadata().GetCategory() == pbrc.ReleaseMetadata_LISTED_TO_SELL {
		s.alertCount++
		s.RaiseIssue(fmt.Sprintf("%v is a problematic record sale id", r.GetRelease().GetInstanceId()), fmt.Sprintf("%v (%v) is missing the sale id", r.GetRelease().GetId(), r.GetRelease().GetInstanceId()))
	}

	//Also validate the current sale price
	if r.GetMetadata().GetSalePriceUpdate() != 0 && time.Since(time.Unix(r.GetMetadata().GetSalePriceUpdate(), 0)) > time.Hour*24*30 && r.GetRelease().GetFolderId() == 812802 {
		s.alertCount++
		s.RaiseIssue(fmt.Sprintf("%v is a problematic record price", r.GetRelease().GetInstanceId()), fmt.Sprintf("%v (%v) has not had a sale price update since %v", r.GetRelease().GetId(), r.GetRelease().GetInstanceId(), time.Unix(r.GetMetadata().GetSalePriceUpdate(), 0)))
	}
}

func (s *Server) alertForPurgatory(r *pbrc.Record) {
	if !r.GetMetadata().GetDirty() && r.GetRelease().GetFolderId() == 1362206 {
		s.alertCount++
		s.RaiseIssue(fmt.Sprintf("%v is a problematic record - purg", r.GetRelease().GetInstanceId()), fmt.Sprintf("[%v]. %v is in Purgatory!", r.GetRelease().GetId(), r.GetRelease().GetTitle()))
	}
}
