package main

import (
	"context"
	"testing"
	"time"

	pbds "github.com/brotherlogic/dstore/client"
	pbgd "github.com/brotherlogic/godiscogs"
	pb "github.com/brotherlogic/recordalerting/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
)

var movetests = []struct {
	in *pbrc.Record
}{
	{&pbrc.Record{Release: &pbgd.Release{FolderId: 812802, Rating: 4, SleeveCondition: "blah", RecordCondition: "blah"},
		Metadata: &pbrc.ReleaseMetadata{Category: pbrc.ReleaseMetadata_STAGED_TO_SELL, GoalFolder: 242017, RecordWidth: 81.2, SalePriceUpdate: time.Now().Unix()}}},
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

func InitTest() (*Server, *testRc) {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
	trc := &testRc{iidMap: make(map[int32]*pbrc.Record)}
	s.rc = trc
	s.ro = &testRo{}

	s.dstoreClient = &pbds.DStoreClient{Test: true}

	// Init an empty config
	s.saveConfig(context.Background(), &pb.Config{})

	return s, trc
}
