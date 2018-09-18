package main

import "fmt"
import pbrc "github.com/brotherlogic/recordcollection/proto"
import pbgd "github.com/brotherlogic/godiscogs"

type testGh struct {
	count int
	fail  bool
}

func (gh *testGh) alert(r *pbrc.Record, text string) error {
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
}

func (rc *testRc) getRecordsInPurgatory() ([]*pbrc.Record, error) {
	if rc.fail {
		return []*pbrc.Record{}, fmt.Errorf("Built to fail")
	}
	return []*pbrc.Record{&pbrc.Record{Release: &pbgd.Release{Title: "MadeUp"}, Metadata: &pbrc.ReleaseMetadata{}}}, nil
}

func (rc *testRc) getSaleRecords() ([]*pbrc.Record, error) {
	if rc.fail {
		return []*pbrc.Record{}, fmt.Errorf("Built to fail")
	}
	return []*pbrc.Record{&pbrc.Record{Release: &pbgd.Release{Title: "MadeUp"}, Metadata: &pbrc.ReleaseMetadata{}}}, nil
}

func (rc *testRc) getLibraryRecords() ([]*pbrc.Record, error) {
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
