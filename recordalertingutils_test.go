package main

import (
	"testing"
	"time"

	pbgd "github.com/brotherlogic/godiscogs"
	pbrc "github.com/brotherlogic/recordcollection/proto"
)

var movetests = []struct {
	in *pbrc.Record
}{
	{&pbrc.Record{Release: &pbgd.Release{FolderId: 812802, Rating: 4, SleeveCondition: "blah", RecordCondition: "blah"},
		Metadata: &pbrc.ReleaseMetadata{Category: pbrc.ReleaseMetadata_STAGED_TO_SELL, GoalFolder: 242017, RecordWidth: 81.2, SalePriceUpdate: time.Now().Unix()}}},
}

func TestAllMoves(t *testing.T) {
	for _, test := range movetests {
		s := InitTest()

		s.assessRecord(test.in)
		if s.alertCount == 0 {
			t.Errorf("%v did not cause an issue", test.in)
		}
	}
}

func TestNoSale(t *testing.T) {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
	s.rc = &testRc{}

	s.alertForMissingSaleID(&pbrc.Record{Metadata: &pbrc.ReleaseMetadata{Category: pbrc.ReleaseMetadata_LISTED_TO_SELL}})
}

func TestPurgatory(t *testing.T) {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
	s.rc = &testRc{}

	s.alertForPurgatory(&pbrc.Record{Release: &pbgd.Release{FolderId: 1362206}})
}

func TestPurgatoryFail(t *testing.T) {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
	s.rc = &testRc{fail: true}

	s.alertForPurgatory(&pbrc.Record{})

}

func InitTest() *Server {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
	s.rc = &testRc{}
	s.ro = &testRo{}

	return s
}

func TestInvalid(t *testing.T) {
	s := InitTest()
	s.validateRecord(&pbrc.Record{})
	s.validateRecord(&pbrc.Record{Metadata: &pbrc.ReleaseMetadata{Category: pbrc.ReleaseMetadata_PURCHASED}})
	s.assessRecord(&pbrc.Record{})
}

func TestInvalidListeningPile(t *testing.T) {
	s := InitTest()
	s.validateRecord(&pbrc.Record{})
	s.validateRecord(&pbrc.Record{Release: &pbgd.Release{FolderId: 812802}, Metadata: &pbrc.ReleaseMetadata{}})
	s.assessRecord(&pbrc.Record{})
}
