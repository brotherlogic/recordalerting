package main

import (
	"fmt"
	"time"

	"github.com/brotherlogic/recordsorganiser/locator"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gd "github.com/brotherlogic/godiscogs/proto"
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
		detail := fmt.Sprintf("This one [%v]: https://www.discogs.com/madeup/release/%v\n", r.GetRelease().GetInstanceId(), r.GetRelease().GetId())
		if class == pb.Problem_MISSING_FILED {
			location, err := locator.ReadableLocation(ctx, s.FDialServer, r.GetRelease().GetInstanceId(), true)
			if err == nil {
				detail = fmt.Sprintf("This one [%v]: https://www.discogs.com/madeup/release/%v\nLocation: \n%v", r.GetRelease().GetInstanceId(), r.GetRelease().GetId(), location)
			} else {
				s.CtxLog(ctx, fmt.Sprintf("Error reading location: %v -> %v", err, location))
			}
		}

		issue, err := s.ImmediateIssue(ctx, fmt.Sprintf("%v [%v] %v", r.GetRelease().GetTitle(), r.GetRelease().GetInstanceId(), errorMessage), detail,
			(class == pb.Problem_MISSING_FILED || class == pb.Problem_MISSING_WEIGHT), true)
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
		s.CtxLog(ctx, fmt.Sprintf("Deleting the record %v for %v -> %v", number, r.GetRelease().GetInstanceId(), err))

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

	if needs &&
		(class == pb.Problem_MISSING_WEIGHT ||
			class == pb.Problem_MISSING_WIDTH ||
			class == pb.Problem_MISSING_FILED ||
			class == pb.Problem_NEEDS_DIGITAL ||
			class == pb.Problem_NEEDS_KEEPER ||
			class == pb.Problem_NEEDS_SOLD_DETAILS ||
			class == pb.Problem_BAD_BANDCAMP ||
			class == pb.Problem_MISSING_SLEEVE) {
		return status.Errorf(codes.FailedPrecondition, "Record %v fails validation - please fix (%v)", r.GetRelease().GetInstanceId(), class)
	}
	return nil
}

func (s *Server) needsWeight(ctx context.Context, config *pb.Config, r *pbrc.Record) error {
	s.CtxLog(ctx, fmt.Sprintf("HERE %v, %v and %v", r.GetRelease().GetFolderId(), r.GetMetadata().GetMoveFolder(), r.GetMetadata().GetWeightInGrams()))
	return s.adjustState(ctx, config, r,
		(r.GetMetadata().GetCategory() == pbrc.ReleaseMetadata_SOLD && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_DIGITAL && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_UNKNOWN && r.GetMetadata().GetWeightInGrams() == 0) ||
			(r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_DIGITAL && r.GetRelease().GetFolderId() == 812802 && r.GetMetadata().GetMoveFolder() > 0 && r.GetMetadata().GetWeightInGrams() <= 0),
		pb.Problem_MISSING_WEIGHT, "needs weight")
}
func (s *Server) needsWidth(ctx context.Context, config *pb.Config, r *pbrc.Record) error {
	return s.adjustState(ctx, config, r,
		r.GetMetadata().GetMoveFolder() == 812802 && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_DIGITAL && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_UNKNOWN && r.GetMetadata().GetRecordWidth() <= 0,
		pb.Problem_MISSING_WIDTH, "needs width")
}

func (s *Server) badBandcamp(ctx context.Context, config *pb.Config, r *pbrc.Record) error {
	file := false
	for _, format := range r.GetRelease().GetFormats() {
		if format.GetName() == "File" {
			file = true
		}
	}

	// We can't assess if there are no formats yet.
	if len(r.GetRelease().GetFormats()) == 0 {
		file = true
	}

	return s.adjustState(ctx, config, r,
		r.GetMetadata().GetGoalFolder() == 1782105 && !file,
		pb.Problem_BAD_BANDCAMP, "bad bandcamp")
}

func (s *Server) needsSaleBudget(ctx context.Context, config *pb.Config, r *pbrc.Record) error {
	return s.adjustState(ctx, config, r,
		r.GetMetadata().GetMoveFolder() == 1613206 && r.GetMetadata().GetSaleBudget() == "",
		pb.Problem_NEEDS_SALE_BUDGET, "needs sale budget")
}

func (s *Server) needsSleeve(ctx context.Context, config *pb.Config, r *pbrc.Record) error {
	return s.adjustState(ctx, config, r,
		r.GetMetadata().GetMoveFolder() == 812802 && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_DIGITAL && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FileSize(pbrc.ReleaseMetadata_VINYL_STORAGE_NO_INNER) && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_UNKNOWN && r.GetMetadata().GetSleeve() == pbrc.ReleaseMetadata_SLEEVE_UNKNOWN,
		pb.Problem_MISSING_SLEEVE, "needs sleeve")
}

func (s *Server) needsSold(ctx context.Context, config *pb.Config, r *pbrc.Record) error {
	return s.adjustState(ctx, config, r,
		r.GetMetadata().GetCategory() == pbrc.ReleaseMetadata_SOLD_ARCHIVE && (r.GetMetadata().GetSoldDate() == 0 || r.GetMetadata().GetSoldPrice() == 0),
		pb.Problem_NEEDS_SOLD_DETAILS, "needs sold details")
}

func (s *Server) needsKeeperJudgement(ctx context.Context, config *pb.Config, r *pbrc.Record) error {
	return s.adjustState(ctx, config, r,
		r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_DIGITAL && r.GetRelease().GetFolderId() == 812802 && r.GetMetadata().GetCategory() == pbrc.ReleaseMetadata_IN_COLLECTION && r.GetMetadata().GetKeep() == pbrc.ReleaseMetadata_KEEP_UNKNOWN, pb.Problem_NEEDS_KEEPER, "needs keeper judgement")
}

func (s *Server) needsDigitalAssess(ctx context.Context, config *pb.Config, r *pbrc.Record) error {
	return s.adjustState(ctx, config, r,
		r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_CD && r.GetRelease().GetFolderId() == 812802 && r.GetMetadata().GetCategory() == pbrc.ReleaseMetadata_IN_COLLECTION && r.GetRelease().GetRating() <= 4 && r.GetMetadata().GetDigitalAvailability() == pbrc.ReleaseMetadata_AVAILABILITY_UNKNOWN, pb.Problem_NEEDS_DIGITAL, "needs digital avail")
}

func (s *Server) needsFiled(ctx context.Context, config *pb.Config, r *pbrc.Record) error {
	s.CtxLog(ctx, fmt.Sprintf("1: %v", r.GetMetadata().GetFiledUnder() == pbrc.ReleaseMetadata_FILE_UNKNOWN))
	s.CtxLog(ctx, fmt.Sprintf("2: %v", r.GetMetadata().GetNewBoxState() == pbrc.ReleaseMetadata_OUT_OF_BOX))
	s.CtxLog(ctx, fmt.Sprintf("3: %v", r.GetMetadata().GetBoxState() == pbrc.ReleaseMetadata_OUT_OF_BOX))
	s.CtxLog(ctx, fmt.Sprintf("4: %v", r.GetMetadata().GetFiledUnder() == pbrc.ReleaseMetadata_FILE_UNKNOWN))
	return s.adjustState(ctx, config, r,
		r.GetMetadata().GetFiledUnder() == pbrc.ReleaseMetadata_FILE_UNKNOWN && (r.GetMetadata().GetNewBoxState() == pbrc.ReleaseMetadata_OUT_OF_BOX || r.GetMetadata().GetBoxState() == pbrc.ReleaseMetadata_OUT_OF_BOX || r.GetMetadata().GetBoxState() == pbrc.ReleaseMetadata_BOX_UNKNOWN),
		pb.Problem_MISSING_FILED, "needs filling")
}
func (s *Server) needsCondition(ctx context.Context, config *pb.Config, r *pbrc.Record) error {
	return s.adjustState(ctx, config, r,
		r.GetMetadata().GetCategory() != pbrc.ReleaseMetadata_UNKNOWN && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_UNKNOWN && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_DIGITAL && r.GetMetadata().GetFiledUnder() != pbrc.ReleaseMetadata_FILE_UNKNOWN && r.GetRelease().GetRecordCondition() == "",
		pb.Problem_MISSING_CONDITION, "needs condition")
}

func (s *Server) assessRecord(ctx context.Context, config *pb.Config, r *pbrc.Record) error {
	// We don't alert on boxed records
	s.CtxLog(ctx, fmt.Sprintf("Assessing the %v", r))
	if r.GetMetadata().GetBoxState() != pbrc.ReleaseMetadata_OUT_OF_BOX && r.GetMetadata().GetBoxState() != pbrc.ReleaseMetadata_BOX_UNKNOWN {
		s.CtxLog(ctx, fmt.Sprintf("Quitting early because of boxstate: %v", r.GetMetadata().GetBoxState()))
		return nil
	}

	// We don't alert on parents records
	if r.GetMetadata().GetCategory() == pbrc.ReleaseMetadata_PARENTS || r.GetMetadata().GetCategory() == pbrc.ReleaseMetadata_SOLD_ARCHIVE {
		return nil
	}

	err1 := s.needsFiled(ctx, config, r)
	err2 := s.needsWeight(ctx, config, r)
	err3 := s.needsWidth(ctx, config, r)
	err4 := s.needsCondition(ctx, config, r)
	err5 := s.needsSleeve(ctx, config, r)
	err6 := s.needsDigitalAssess(ctx, config, r)
	err7 := s.needsKeeperJudgement(ctx, config, r)
	err8 := s.needsSaleBudget(ctx, config, r)
	err9 := s.needsSold(ctx, config, r)
	s.badBandcamp(ctx, config, r)

	s.CtxLog(ctx, fmt.Sprintf("Run assess: %v, %v, %v, %v, %v, %v, %v", err1, err2, err3, err4, err5, err6, err7))

	// Only fail
	if r.GetMetadata().GetCategory() != pbrc.ReleaseMetadata_UNKNOWN {
		if (r.GetMetadata().GetFiledUnder() == pbrc.ReleaseMetadata_FILE_12_INCH || r.GetMetadata().GetFiledUnder() == pbrc.ReleaseMetadata_FILE_7_INCH) && r.GetMetadata().GetSaleState() != gd.SaleState_SOLD && r.GetMetadata().GetCategory() != pbrc.ReleaseMetadata_SOLD_ARCHIVE {
			if time.Since(time.Unix(r.GetMetadata().GetLastCleanDate(), 0)) > time.Hour*24*365*3 {
				if r.GetRelease().GetFolderId() != 3386035 && r.GetMetadata().GetMoveFolder() != 3386035 {
					err := s.rc.clean(ctx, r.GetRelease().GetInstanceId())
					if err != nil {
						return err
					}
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
	if err5 != nil {
		return err5
	}
	if err6 != nil {
		return err6
	}
	if err7 != nil {
		return err7
	}

	if err8 != nil {
		return err8
	}

	if err9 != nil {
		return err9
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
}

func (s *Server) alertForPurgatory(r *pbrc.Record) {
	if !r.GetMetadata().GetDirty() && r.GetRelease().GetFolderId() == 1362206 {
		s.alertCount++
		s.RaiseIssue(fmt.Sprintf("%v is a problematic record - purg", r.GetRelease().GetInstanceId()), fmt.Sprintf("[%v]. %v is in Purgatory!", r.GetRelease().GetId(), r.GetRelease().GetTitle()))
	}
}
