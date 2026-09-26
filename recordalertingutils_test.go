package main

import (
	"context"
	"testing"
	"time"

	pbds "github.com/brotherlogic/dstore/client"
	pbgd "github.com/brotherlogic/godiscogs/proto"
	pb "github.com/brotherlogic/recordalerting/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func TestNeedsPackageScore_MissingScoreFails(t *testing.T) {
	s, _ := InitTest()
	config := &pb.Config{}

	r := &pbrc.Record{
		Release: &pbgd.Release{
			InstanceId: 1001,
			FolderId:   812802,
			Rating:     4,
			Title:      "Test Album",
		},
		Metadata: &pbrc.ReleaseMetadata{
			PackageScore: -1,
		},
	}

	err := s.needsPackageScore(context.Background(), config, r)
	if err == nil {
		t.Fatalf("Expected error for missing package score on qualified record, got nil")
	}
	if status.Code(err) != codes.FailedPrecondition {
		t.Errorf("Expected FailedPrecondition status code, got %v", status.Code(err))
	}

	found := false
	for _, p := range config.GetProblems() {
		if p.GetType() == pb.Problem_MISSING_PACKAGE_SCORE && p.GetInstanceId() == 1001 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected problem MISSING_PACKAGE_SCORE registered in config.Problems for instance 1001, got %+v", config.GetProblems())
	}
}

func TestNeedsPackageScore_ZeroScorePasses(t *testing.T) {
	s, _ := InitTest()
	config := &pb.Config{}

	r := &pbrc.Record{
		Release: &pbgd.Release{
			InstanceId: 1002,
			FolderId:   812802,
			Rating:     4,
			Title:      "Zero Score Album",
		},
		Metadata: &pbrc.ReleaseMetadata{
			PackageScore: 0,
		},
	}

	err := s.needsPackageScore(context.Background(), config, r)
	if err != nil {
		t.Errorf("Expected nil error for zero package score, got %v", err)
	}
}

func TestNeedsPackageScore_PositiveScorePasses(t *testing.T) {
	s, _ := InitTest()
	config := &pb.Config{}

	r := &pbrc.Record{
		Release: &pbgd.Release{
			InstanceId: 1003,
			FolderId:   812802,
			Rating:     5,
			Title:      "Positive Score Album",
		},
		Metadata: &pbrc.ReleaseMetadata{
			PackageScore: 4,
		},
	}

	err := s.needsPackageScore(context.Background(), config, r)
	if err != nil {
		t.Errorf("Expected nil error for positive package score, got %v", err)
	}
}

func TestNeedsPackageScore_LowRatingExempt(t *testing.T) {
	s, _ := InitTest()
	config := &pb.Config{}

	for _, rating := range []int32{0, 1} {
		r := &pbrc.Record{
			Release: &pbgd.Release{
				InstanceId: int64(2000 + rating),
				FolderId:   812802,
				Rating:     rating,
				Title:      "Low Rating Album",
			},
			Metadata: &pbrc.ReleaseMetadata{
				PackageScore: -1,
			},
		}

		err := s.needsPackageScore(context.Background(), config, r)
		if err != nil {
			t.Errorf("Expected nil error for low rating (%d), got %v", rating, err)
		}
	}
}

func TestNeedsPackageScore_DigitalExempt(t *testing.T) {
	s, _ := InitTest()
	config := &pb.Config{}

	// Case 1: FiledUnder == FILE_DIGITAL
	r1 := &pbrc.Record{
		Release: &pbgd.Release{
			InstanceId: 3001,
			FolderId:   812802,
			Rating:     4,
			Title:      "Digital Release (FiledUnder)",
		},
		Metadata: &pbrc.ReleaseMetadata{
			PackageScore: -1,
			FiledUnder:   pbrc.ReleaseMetadata_FILE_DIGITAL,
		},
	}
	if err := s.needsPackageScore(context.Background(), config, r1); err != nil {
		t.Errorf("Expected digital release (FiledUnder) to be exempt, got %v", err)
	}

	// Case 2: Format name == "File"
	r2 := &pbrc.Record{
		Release: &pbgd.Release{
			InstanceId: 3002,
			FolderId:   812802,
			Rating:     4,
			Title:      "Digital Release (Format File)",
			Formats: []*pbgd.Format{
				{Name: "File"},
			},
		},
		Metadata: &pbrc.ReleaseMetadata{
			PackageScore: -1,
		},
	}
	if err := s.needsPackageScore(context.Background(), config, r2); err != nil {
		t.Errorf("Expected digital release (format File) to be exempt, got %v", err)
	}
}

func TestNeedsPackageScore_NonListeningFolderExempt(t *testing.T) {
	s, _ := InitTest()
	config := &pb.Config{}

	r := &pbrc.Record{
		Release: &pbgd.Release{
			InstanceId: 4001,
			FolderId:   1234,
			Rating:     5,
			Title:      "General Folder Record",
		},
		Metadata: &pbrc.ReleaseMetadata{
			PackageScore: -1,
		},
	}

	err := s.needsPackageScore(context.Background(), config, r)
	if err != nil {
		t.Errorf("Expected non-listening folder to be exempt, got %v", err)
	}
}

func TestNeedsPackageScore_AllListeningFolders(t *testing.T) {
	s, _ := InitTest()
	config := &pb.Config{}

	listeningFolders := []int32{812802, 7665013, 7664293, 7651472}
	for _, folder := range listeningFolders {
		if !isListeningPile(folder) {
			t.Errorf("Expected folder %d to be recognized as listening pile", folder)
		}

		r := &pbrc.Record{
			Release: &pbgd.Release{
				InstanceId: int64(5000 + folder),
				FolderId:   folder,
				Rating:     3,
				Title:      "Listening Pile Record",
			},
			Metadata: &pbrc.ReleaseMetadata{
				PackageScore: -1,
			},
		}

		err := s.needsPackageScore(context.Background(), config, r)
		if err == nil {
			t.Errorf("Expected error for folder %d with missing package score, got nil", folder)
		}
	}

	nonListeningFolders := []int32{0, 1234, 3380098, 242017}
	for _, folder := range nonListeningFolders {
		if isListeningPile(folder) {
			t.Errorf("Expected folder %d to NOT be recognized as listening pile", folder)
		}
	}
}

func TestNeedsPackageScore_ResolutionCleansUpProblem(t *testing.T) {
	s, _ := InitTest()
	config := &pb.Config{}

	r := &pbrc.Record{
		Release: &pbgd.Release{
			InstanceId: 6001,
			FolderId:   812802,
			Rating:     4,
			Title:      "Resolving Record",
		},
		Metadata: &pbrc.ReleaseMetadata{
			PackageScore: -1,
		},
	}

	// 1. Initial assessment fails and adds problem
	err := s.needsPackageScore(context.Background(), config, r)
	if err == nil {
		t.Fatalf("Expected error for initial missing package score, got nil")
	}

	found := false
	for _, p := range config.GetProblems() {
		if p.GetType() == pb.Problem_MISSING_PACKAGE_SCORE && p.GetInstanceId() == 6001 {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Expected problem to be registered in config.Problems")
	}

	// 2. Score is resolved to valid value (e.g., 3)
	r.Metadata.PackageScore = 3
	err = s.needsPackageScore(context.Background(), config, r)
	if err != nil {
		t.Errorf("Expected nil error after resolving package score, got %v", err)
	}

	// 3. Verify problem was cleaned up from config.Problems
	for _, p := range config.GetProblems() {
		if p.GetType() == pb.Problem_MISSING_PACKAGE_SCORE && p.GetInstanceId() == 6001 {
			t.Errorf("Expected problem to be removed from config.Problems after resolution, but it still exists")
		}
	}
}

