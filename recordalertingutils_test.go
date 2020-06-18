package main

import (
	"testing"

	pbrc "github.com/brotherlogic/recordcollection/proto"
)

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

	s.alertForPurgatory(&pbrc.Record{})

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

func TestOldListeningBox(t *testing.T) {
	s := InitTest()

	s.alertForOldListeningBoxRecord(&pbrc.Record{})

}

func TestOldListeningPile(t *testing.T) {
	s := InitTest()

	s.alertForOldListeningPileRecord(&pbrc.Record{})

}

func TestOldListeningPileFailRO(t *testing.T) {
	s := InitTest()
	s.ro = &testRo{fail: true}

	s.alertForOldListeningPileRecord(&pbrc.Record{})

}

func TestOldListeningPileFailRC(t *testing.T) {
	s := InitTest()
	s.rc = &testRc{fail: true}

	s.alertForOldListeningPileRecord(&pbrc.Record{})

}

func TestInvalid(t *testing.T) {
	s := InitTest()
	s.validateRecord(&pbrc.Record{})
	s.validateRecord(&pbrc.Record{Metadata: &pbrc.ReleaseMetadata{Category: pbrc.ReleaseMetadata_PURCHASED}})
	s.assessRecord(&pbrc.Record{})
}
