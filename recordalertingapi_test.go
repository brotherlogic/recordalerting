package main

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"

	pbgd "github.com/brotherlogic/godiscogs"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	pbro "github.com/brotherlogic/recordsorganiser/proto"
)

type testRo struct {
	fail bool
}

func (t *testRo) getLocation(ctx context.Context, name string) (*pbro.Location, error) {
	if t.fail {
		return nil, fmt.Errorf("Built to fail")
	}
	if name == "Listening Box" || name == "Listening Pile" {
		return &pbro.Location{ReleasesLocation: []*pbro.ReleasePlacement{&pbro.ReleasePlacement{InstanceId: 1234}}}, nil
	}
	return nil, fmt.Errorf("Unknown Location")
}

type testGh struct {
	count int
	fail  bool
}

func (gh *testGh) alert(ctx context.Context, r *pbrc.Record, text string) error {
	if gh.fail {
		return fmt.Errorf("Built to fail")
	}
	gh.count++
	return nil
}

type testRc struct {
	fail    bool
	order   bool
	missing bool
	failAll bool
	iidMap  map[int32]*pbrc.Record
}

func (rc *testRc) addRecord(iid int32, r *pbrc.Record) {
	r.GetRelease().InstanceId = iid
	rc.iidMap[iid] = r
}

func (rc *testRc) clean(ctx context.Context, instanceID int32) error {
	return nil
}

func (rc *testRc) getRecord(ctx context.Context, instanceID int32) (*pbrc.Record, error) {
	if val, ok := rc.iidMap[instanceID]; ok {
		return val, nil
	}
	if rc.fail || rc.failAll {
		return nil, fmt.Errorf("Built to fail")
	}
	if instanceID == 1233 {
		return &pbrc.Record{Release: &pbgd.Release{Title: "Madeup"}, Metadata: &pbrc.ReleaseMetadata{DateAdded: 1234, Category: pbrc.ReleaseMetadata_PURCHASED}}, nil
	}

	if instanceID == 1234 {
		return &pbrc.Record{Release: &pbgd.Release{Title: "Madeup"}, Metadata: &pbrc.ReleaseMetadata{DateAdded: 1234, Category: pbrc.ReleaseMetadata_LISTED_TO_SELL}}, nil
	}
	return nil, fmt.Errorf("Unknown record")
}

func (rc *testRc) getRecordsInPurgatory(ctx context.Context) ([]*pbrc.Record, error) {
	if rc.fail {
		return []*pbrc.Record{}, fmt.Errorf("Built to fail")
	}
	return []*pbrc.Record{&pbrc.Record{Release: &pbgd.Release{Title: "MadeUp"}, Metadata: &pbrc.ReleaseMetadata{}}}, nil
}

func (rc *testRc) getRecordsInFolder(ctx context.Context, folder int32) ([]int32, error) {
	if rc.fail {
		return []int32{}, fmt.Errorf("Built to fail")
	}
	return []int32{1234, 1233}, nil
}

func (rc *testRc) getRecords(ctx context.Context) ([]*pbrc.Record, error) {
	if rc.failAll {
		return []*pbrc.Record{}, fmt.Errorf("Built to fail")
	}
	return []*pbrc.Record{&pbrc.Record{Release: &pbgd.Release{Title: "MadeUp", FolderId: 812802}, Metadata: &pbrc.ReleaseMetadata{}}}, nil
}

func (rc *testRc) getSaleRecords(ctx context.Context) ([]*pbrc.Record, error) {
	if rc.fail {
		return []*pbrc.Record{}, fmt.Errorf("Built to fail")
	}
	return []*pbrc.Record{&pbrc.Record{Release: &pbgd.Release{Title: "MadeUp"}, Metadata: &pbrc.ReleaseMetadata{Category: pbrc.ReleaseMetadata_LISTED_TO_SELL}}}, nil
}

func (rc *testRc) getLibraryRecords(ctx context.Context) ([]*pbrc.Record, error) {
	if rc.fail {
		return []*pbrc.Record{}, fmt.Errorf("Built to fail")
	}

	if !rc.order {
		return []*pbrc.Record{
			&pbrc.Record{Release: &pbgd.Release{Title: "Jazz Moderne"}},
			&pbrc.Record{Release: &pbgd.Release{Title: "Action Charme Espace"}},
			&pbrc.Record{Release: &pbgd.Release{Title: "Paysages, Evasion, Melancolie"}},
			&pbrc.Record{Release: &pbgd.Release{Title: "Sports Et Action"}},
		}, nil
	}

	if rc.missing {
		return []*pbrc.Record{
			&pbrc.Record{Release: &pbgd.Release{Title: "Jazz Moderne"}},
			&pbrc.Record{Release: &pbgd.Release{Title: "Paysages, Evasion, Melancolie"}},
			&pbrc.Record{Release: &pbgd.Release{Title: "Sports Et Action"}},
		}, nil
	}

	return []*pbrc.Record{
		&pbrc.Record{Release: &pbgd.Release{Title: "Action Charme Espace"}},
		&pbrc.Record{Release: &pbgd.Release{Title: "Paysages, Evasion, Melancolie"}},
		&pbrc.Record{Release: &pbgd.Release{Title: "Jazz Moderne"}},
		&pbrc.Record{Release: &pbgd.Release{Title: "Sports Et Action"}},
	}, nil
}

func TestFilledRecordIntoCollectionButNoWeight(t *testing.T) {
	s, rc := InitTest()

	// Record is moving from listening pile into collection
	rc.addRecord(1234, &pbrc.Record{
		Release: &pbgd.Release{FolderId: 812802, Rating: 5},
		Metadata: &pbrc.ReleaseMetadata{
			MoveFolder: 242017,
			Category:   pbrc.ReleaseMetadata_IN_COLLECTION,
			Keep:       pbrc.ReleaseMetadata_DIGITAL_KEEPER,
			FiledUnder: pbrc.ReleaseMetadata_FILE_12_INCH,
		},
	})

	_, err := s.ClientUpdate(context.Background(), &pbrc.ClientUpdateRequest{InstanceId: 1234})
	if err == nil {
		t.Errorf("We expected this to be an error condition")
	}

	rc.addRecord(1234, &pbrc.Record{
		Release: &pbgd.Release{FolderId: 812802, Rating: 5},
		Metadata: &pbrc.ReleaseMetadata{
			MoveFolder:    242017,
			Category:      pbrc.ReleaseMetadata_IN_COLLECTION,
			WeightInGrams: 300,
			Keep:          pbrc.ReleaseMetadata_DIGITAL_KEEPER,
			FiledUnder:    pbrc.ReleaseMetadata_FILE_12_INCH,
		},
	})

	_, err = s.ClientUpdate(context.Background(), &pbrc.ClientUpdateRequest{InstanceId: 1234})
	if err != nil {
		t.Errorf("An unexpected error on this move: %v", err)
	}
}

func TestFilledRecordIntoCollectionButNoDigitalKeep(t *testing.T) {
	s, rc := InitTest()

	// Record is moving from listening pile into collection
	rc.addRecord(1234, &pbrc.Record{
		Release: &pbgd.Release{FolderId: 812802, Rating: 5},
		Metadata: &pbrc.ReleaseMetadata{
			MoveFolder:    242017,
			Category:      pbrc.ReleaseMetadata_IN_COLLECTION,
			FiledUnder:    pbrc.ReleaseMetadata_FILE_12_INCH,
			WeightInGrams: 300,
		},
	})

	_, err := s.ClientUpdate(context.Background(), &pbrc.ClientUpdateRequest{InstanceId: 1234})
	if err == nil {
		t.Errorf("We expected this to be an error condition")
	}

	rc.addRecord(1234, &pbrc.Record{
		Release: &pbgd.Release{FolderId: 812802, Rating: 5},
		Metadata: &pbrc.ReleaseMetadata{
			MoveFolder:    242017,
			Category:      pbrc.ReleaseMetadata_IN_COLLECTION,
			WeightInGrams: 300,
			Keep:          pbrc.ReleaseMetadata_DIGITAL_KEEPER,
			FiledUnder:    pbrc.ReleaseMetadata_FILE_12_INCH,
		},
	})

	_, err = s.ClientUpdate(context.Background(), &pbrc.ClientUpdateRequest{InstanceId: 1234})
	if err != nil {
		t.Errorf("An unexpected error on this move: %v", err)
	}
}
