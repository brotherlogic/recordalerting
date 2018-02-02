package main

import "fmt"
import pbrc "github.com/brotherlogic/recordcollection/proto"

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
