package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
)

func (s *Server) validateRecords(ctx context.Context) error {
	records, err := s.rc.getRecords(ctx)
	if err != nil {
		return err
	}
	missingCondition := []int32{}
	for _, r := range records {
		if r.GetMetadata().GoalFolder == 0 {
			s.RaiseIssue(ctx, "Invalid Record", fmt.Sprintf("%v has no goal folder", r.GetRelease().InstanceId), false)
			s.invalidRecords++
		}

		if r.GetRelease().FolderId == 812802 && (len(r.GetRelease().RecordCondition) == 0 || len(r.GetRelease().SleeveCondition) == 0) {
			missingCondition = append(missingCondition, r.GetRelease().InstanceId)
		}
	}

	if len(missingCondition) > 0 {
		s.RaiseIssue(ctx, "Missing Conditions", fmt.Sprintf("%v are missing condition", missingCondition), false)
	}

	return nil
}

func (s *Server) alertForMissingSaleID(ctx context.Context) error {
	records, err := s.rc.getSaleRecords(ctx)
	s.Log(fmt.Sprintf("Found %v records for sale", len(records)))
	if err == nil {
		for _, r := range records {
			if r.GetMetadata().SaleId == 0 {
				s.gh.alert(ctx, r, fmt.Sprintf("%v (%v) is missing the sale id", r.GetRelease().Id, r.GetRelease().InstanceId))
			}
		}
	}
	return nil
}

func (s *Server) alertForPurgatory(ctx context.Context) error {
	records, err := s.rc.getRecordsInPurgatory(ctx)
	if err == nil {
		if len(records) > 0 {
			for _, r := range records {
				if !r.GetMetadata().GetDirty() {
					s.gh.alert(ctx, records[0], fmt.Sprintf("[%v]. %v is in Purgatory!", records[0].GetRelease().Id, records[0].GetRelease().Title))
					break
				}
			}
		}
	} else {
		s.Log(fmt.Sprintf("Error getting records: %v", err))
	}

	return nil
}

func (s *Server) alertForMisorderedMPI(ctx context.Context) error {
	records, err := s.rc.getLibraryRecords(ctx)

	if err != nil {
		s.Log(fmt.Sprintf("Error getting library records: %v", err))
		return err
	}

	order := []string{"Jazz Moderne", "Action Charme Espace", "Paysages, Evasion, Melancolie", "Sports Et Action"}
	lastSeen := -1

	fail := false
	for _, r := range records {
		for i, o := range order {
			if r.GetRelease().Title == o {
				if i < lastSeen {
					fail = true
				}
				lastSeen = i
			}
		}
	}

	if fail {
		s.gh.alert(ctx, nil, fmt.Sprintf("MPI is not ordered correctly!: %v", records))
	}
	return nil
}

func (s *Server) alertForOldListeningBoxRecord(ctx context.Context) error {
	records, err := s.ro.getLocation(ctx, "Listening Box")
	if err == nil {
		s.Log(fmt.Sprintf("FOUND %v records", len(records.ReleasesLocation)))
		for _, r := range records.ReleasesLocation {
			rec, err := s.rc.getRecord(ctx, r.InstanceId)
			if err == nil {
				s.Log(fmt.Sprintf("Record %v has %v", rec.GetRelease().Title, time.Now().Sub(time.Unix(rec.GetMetadata().DateAdded, 0))))
			}
			if err == nil && time.Now().Sub(time.Unix(rec.GetMetadata().DateAdded, 0)) > time.Hour*24*30*4 {
				s.gh.alert(ctx, nil, fmt.Sprintf("Record %v has been in the listening box for %v", rec.GetRelease().Title, time.Now().Sub(time.Unix(rec.GetMetadata().DateAdded, 0))))
			}
			s.Log(fmt.Sprintf("Error in get record? %v", err))
		}
	}

	return err
}

func (s *Server) alertForOldListeningPileRecord(ctx context.Context) error {
	records, err := s.ro.getLocation(ctx, "Listening Pile")
	if err != nil {
		return err
	}

	for _, r := range records.ReleasesLocation {
		rec, err := s.rc.getRecord(ctx, r.InstanceId)
		if err != nil {
			return err
		}
		if time.Now().Sub(time.Unix(rec.GetMetadata().LastMoveTime, 0)) > time.Hour*24*30 {
			s.gh.alert(ctx, nil, fmt.Sprintf("Record %v has been in the listening pile for %v", rec.GetRelease().Title, time.Now().Sub(time.Unix(rec.GetMetadata().LastMoveTime, 0))))
		}
	}

	return nil
}
