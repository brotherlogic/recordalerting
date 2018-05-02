package main

import (
	"fmt"

	"golang.org/x/net/context"
)

func (s *Server) alertForPurgatory(ctx context.Context) {
	records, err := s.rc.getRecordsInPurgatory()
	if err == nil {
		if len(records) > 0 {
			for _, r := range records {
				if !r.GetMetadata().GetDirty() {
					s.gh.alert(records[0], fmt.Sprintf("%v is in Purgatory!", records[0].GetRelease().Title))
					break
				}
			}
		}
	} else {
		s.Log(fmt.Sprintf("Error getting records: %v", err))
	}
}

func (s *Server) alertForMisorderedMPI(ctx context.Context) {
	records, err := s.rc.getLibraryRecords()

	if err != nil {
		s.Log(fmt.Sprintf("Error getting library records: %v", err))
		return
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
		s.gh.alert(nil, fmt.Sprintf("MPI is not ordered correctly!"))
	}
	return
}
