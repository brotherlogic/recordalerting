package main

import "testing"

func TestPurgatory(t *testing.T) {
	s := Init()
	gh := &testGh{}
	s.rc = &testRc{}
	s.gh = gh

	s.alertForPurgatory()

	if gh.count == 0 {
		t.Errorf("No errors sent!")
	}
}

func TestPurgatoryFail(t *testing.T) {
	s := Init()
	s.SkipLog = true
	gh := &testGh{}
	s.rc = &testRc{fail: true}
	s.gh = gh

	s.alertForPurgatory()

	if gh.count != 0 {
		t.Errorf("Errors have been sent!")
	}
}
