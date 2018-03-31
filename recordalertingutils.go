package main

import (
	"fmt"
)

func (s *Server) alertForPurgatory() {
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

func (s *Server) alertForMisorderedMPI() {
	records, err := s.rc.getLibraryRecords()

	if err != nil {
		s.Log(fmt.Sprintf("Error getting library records: %v", err))
		return
	}

	order := []string{"Jazz Moderne", "Action Charme Espace", "Paysages, Evasion, Melancolie", "Sports Et Action"}
	pointer := 0

	for _, r := range records {
		if r.GetRelease().Title == order[pointer] {
			pointer++

			if pointer >= len(order) {
				return
			}
		}
	}

	s.gh.alert(nil, fmt.Sprintf("MPI is not ordered correctly!"))
	return
}
