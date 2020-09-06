package main

import (
	"testing"

	pbgd "github.com/brotherlogic/godiscogs"
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
