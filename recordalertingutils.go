package main

import (
	"fmt"
	"time"

	pbrc "github.com/brotherlogic/recordcollection/proto"
)

func (s *Server) assessRecord(r *pbrc.Record) error {
	// We don't alert on boxed records
	if r.GetMetadata().GetBoxState() != pbrc.ReleaseMetadata_OUT_OF_BOX && r.GetMetadata().GetBoxState() != pbrc.ReleaseMetadata_BOX_UNKNOWN {
		return nil
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

	if (len(r.GetRelease().GetRecordCondition()) == 0 || len(r.GetRelease().GetSleeveCondition()) == 0) &&
		r.GetMetadata().GetGoalFolder() != 1782105 &&
		r.GetMetadata().GetGoalFolder() != 1433217 &&
		r.GetMetadata().GetGoalFolder() != 2274270 &&
		r.GetMetadata().GetGoalFolder() != 268147 {
		if r.GetRelease().GetFolderId() == 812802 {
			s.alertCount++
			s.RaiseIssue(fmt.Sprintf("%v Missing Condition", r.GetRelease().GetInstanceId()), fmt.Sprintf("%v is missing conditions. Fix it here https://www.discogs.com/madeup/release/%v", r.GetRelease().GetInstanceId(), r.GetRelease().GetId()))
		}
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
	if time.Now().Sub(time.Unix(r.GetMetadata().GetSalePriceUpdate(), 0)) > time.Hour*24*7 {
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
