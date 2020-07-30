package main

import (
	"fmt"
	"time"

	pbrc "github.com/brotherlogic/recordcollection/proto"
)

func (s *Server) assessRecord(r *pbrc.Record) error {
	s.validateRecord(r)
	s.alertForMissingSaleID(r)
	s.alertForPurgatory(r)
	s.alertForOldListeningBoxRecord(r)
	s.alertForOldListeningPileRecord(r)

	return nil
}

func (s *Server) validateRecord(r *pbrc.Record) {
	if r.GetMetadata().GetGoalFolder() == 0 {
		s.RaiseIssue(fmt.Sprintf("%v Invalid Record", r.GetRelease().GetInstanceId()), fmt.Sprintf("%v has no goal folder", r.GetRelease().GetInstanceId()))
		s.invalidRecords++
	}

	if (len(r.GetRelease().GetRecordCondition()) == 0 || len(r.GetRelease().GetSleeveCondition()) == 0) && r.GetMetadata().GetGoalFolder() != 1782105 {
		s.RaiseIssue(fmt.Sprintf("%v Missing Condition", r.GetRelease().GetInstanceId()), fmt.Sprintf("%v is missing conditions", r.GetRelease().GetInstanceId()))
	}

	if r.GetMetadata().GetCategory() == pbrc.ReleaseMetadata_PURCHASED && time.Now().Sub(time.Unix(r.GetMetadata().GetLastUpdateTime(), 0)) > time.Hour*24 {
		s.RaiseIssue(fmt.Sprintf("%v Stale Purchase", r.GetRelease().GetInstanceId()), fmt.Sprintf("%v has staled", r.GetRelease().GetInstanceId()))
	}
}

func (s *Server) alertForMissingSaleID(r *pbrc.Record) {
	if r.GetMetadata().GetSaleId() <= 0 && r.GetMetadata().GetCategory() == pbrc.ReleaseMetadata_LISTED_TO_SELL {
		s.RaiseIssue(fmt.Sprintf("%v is a problematic record sale id", r.GetRelease().GetInstanceId()), fmt.Sprintf("%v (%v) is missing the sale id", r.GetRelease().GetId(), r.GetRelease().GetInstanceId()))
	}

	//Also validate the current sale price
	if time.Now().Sub(time.Unix(r.GetMetadata().GetSalePriceUpdate(), 0)) > time.Hour*24*7 {
		s.RaiseIssue(fmt.Sprintf("%v is a problematic record price", r.GetRelease().GetInstanceId()), fmt.Sprintf("%v (%v) has not had a sale price update since %v", r.GetRelease().GetId(), r.GetRelease().GetInstanceId(), time.Unix(r.GetMetadata().GetSalePriceUpdate(), 0)))
	}
}

func (s *Server) alertForPurgatory(r *pbrc.Record) {
	if !r.GetMetadata().GetDirty() && r.GetRelease().GetFolderId() == 1362206 {
		s.RaiseIssue(fmt.Sprintf("%v is a problematic record - purg", r.GetRelease().GetInstanceId()), fmt.Sprintf("[%v]. %v is in Purgatory!", r.GetRelease().GetId(), r.GetRelease().GetTitle()))
	}
}

func (s *Server) alertForOldListeningBoxRecord(rec *pbrc.Record) {
	if time.Now().Sub(time.Unix(rec.GetMetadata().GetDateAdded(), 0)) > time.Hour*24*30*4 {
		s.RaiseIssue(fmt.Sprintf("%v old listening box", rec.GetRelease().GetInstanceId()), fmt.Sprintf("Record %v [%v] has been in the listening box for %v", rec.GetRelease().GetTitle(), rec.GetRelease().GetInstanceId(), time.Now().Sub(time.Unix(rec.GetMetadata().GetDateAdded(), 0))))
	}
}

func (s *Server) alertForOldListeningPileRecord(rec *pbrc.Record) {
	if time.Now().Sub(time.Unix(rec.GetMetadata().GetLastMoveTime(), 0)) > time.Hour*24*30 {
		s.RaiseIssue(fmt.Sprintf("%v old listening pile", rec.GetRelease().GetInstanceId()), fmt.Sprintf("Record %v has been in the listening pile for %v", rec.GetRelease().GetTitle(), time.Now().Sub(time.Unix(rec.GetMetadata().GetLastMoveTime(), 0))))
	}
}
