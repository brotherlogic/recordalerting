package main

import "fmt"

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
	}
}
