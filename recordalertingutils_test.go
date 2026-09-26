package main

import (
	"context"
	"testing"
	"time"

	pbds "github.com/brotherlogic/dstore/client"
	pbgd "github.com/brotherlogic/godiscogs/proto"
	pb "github.com/brotherlogic/recordalerting/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
)

var movetests = []struct {
	in *pbrc.Record
}{
	{&pbrc.Record{Release: &pbgd.Release{FolderId: 812802, Rating: 4, SleeveCondition: "blah", RecordCondition: "blah"},
		Metadata: &pbrc.ReleaseMetadata{Category: pbrc.ReleaseMetadata_STAGED_TO_SELL, GoalFolder: 242017, RecordWidth: 81.2, SalePriceUpdate: time.Now().Unix()}}},
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
	trc := &testRc{iidMap: make(map[int64]*pbrc.Record)}
	s.rc = trc
	s.ro = &testRo{}

	s.dstoreClient = &pbds.DStoreClient{Test: true}

	// Init an empty config
	s.saveConfig(context.Background(), &pb.Config{})

	return s, trc
}
func TestNeedsNotes(t *testing.T) {
	s, _ := InitTest()
 
	// Record in allowed folder with no notes and 0 rating -> should pass (return nil)
	r1 := &pbrc.Record{
		Release: &pbgd.Release{FolderId: 7651472, Rating: 0},
		Metadata: &pbrc.ReleaseMetadata{Notes: ""},
	}
	err1 := s.needsNotes(context.Background(), &pb.Config{}, r1)
	if err1 != nil {
		t.Errorf("Should have passed for missing notes in folder 7651472 with 0 rating: %v", err1)
	}
 
	// Record in disallowed folder with no notes -> should pass (return nil)
	r2 := &pbrc.Record{
		Release: &pbgd.Release{FolderId: 1234},
		Metadata: &pbrc.ReleaseMetadata{Notes: ""},
	}
	err2 := s.needsNotes(context.Background(), &pb.Config{}, r2)
	if err2 != nil {
		t.Errorf("Should have passed for missing notes in folder 1234: %v", err2)
	}
 
	// Record in allowed folder with notes -> should pass (return nil)
	r3 := &pbrc.Record{
		Release: &pbgd.Release{FolderId: 7651472},
		Metadata: &pbrc.ReleaseMetadata{Notes: "Some notes"},
	}
	err3 := s.needsNotes(context.Background(), &pb.Config{}, r3)
	if err3 != nil {
		t.Errorf("Should have passed with notes in folder 7651472: %v", err3)
	}

	// Record in allowed folder with no notes and > 0 rating -> should fail (return error)
	r4 := &pbrc.Record{
		Release: &pbgd.Release{FolderId: 7651472, Rating: 1},
		Metadata: &pbrc.ReleaseMetadata{Notes: ""},
	}
	err4 := s.needsNotes(context.Background(), &pb.Config{}, r4)
	if err4 == nil {
		t.Errorf("Should have failed for missing notes in folder 7651472 with > 0 rating")
	}
}

func TestMissingPackageScoreEnum(t *testing.T) {
	if int32(pb.Problem_MISSING_PACKAGE_SCORE) != 14 {
		t.Errorf("Expected Problem_MISSING_PACKAGE_SCORE to be 14, got %d", pb.Problem_MISSING_PACKAGE_SCORE)
	}
	if name, ok := pb.Problem_ProblemType_name[14]; !ok || name != "MISSING_PACKAGE_SCORE" {
		t.Errorf("Expected ProblemType name for 14 to be MISSING_PACKAGE_SCORE, got %q", name)
	}

	metadata := &pbrc.ReleaseMetadata{
		PackageScore: 78,
	}
	if metadata.GetPackageScore() != 78 {
		t.Errorf("Expected GetPackageScore to return 78, got %d", metadata.GetPackageScore())
	}
}

