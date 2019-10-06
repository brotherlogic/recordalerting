package main

import (
	"fmt"

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
}

func (rc *testRc) getRecord(ctx context.Context, instanceID int32) (*pbrc.Record, error) {
	if rc.fail || rc.failAll {
		return nil, fmt.Errorf("Built to fail")
	}
	if instanceID == 1234 {
		return &pbrc.Record{Release: &pbgd.Release{Title: "Madeup"}, Metadata: &pbrc.ReleaseMetadata{DateAdded: 1234}}, nil
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
	return []int32{1234}, nil
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
