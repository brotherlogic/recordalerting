package main

import (
	"testing"

	"golang.org/x/net/context"
)

func TestPurgatory(t *testing.T) {
	s := Init()
	s.SkipLog = true
	gh := &testGh{}
	s.rc = &testRc{}
	s.gh = gh

	s.alertForPurgatory(context.Background())

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

	s.alertForPurgatory(context.Background())

	if gh.count != 0 {
		t.Errorf("Errors have been sent!")
	}
}

func TestMPI(t *testing.T) {
	s := Init()
	s.SkipLog = true
	gh := &testGh{}
	s.rc = &testRc{}
	s.gh = gh

	s.alertForMisorderedMPI(context.Background())

	if gh.count != 0 {
		t.Errorf("Errors have been sent!")
	}
}

func TestMPIFail(t *testing.T) {
	s := Init()
	s.SkipLog = true
	gh := &testGh{}
	s.rc = &testRc{fail: true}
	s.gh = gh

	s.alertForMisorderedMPI(context.Background())
	if gh.count != 0 {
		t.Errorf("Errors have been sent!")
	}
}

func TestMPIOrder(t *testing.T) {
	s := Init()
	s.SkipLog = true
	gh := &testGh{}
	s.rc = &testRc{order: true}
	s.gh = gh

	s.alertForMisorderedMPI(context.Background())

	if gh.count == 0 {
		t.Errorf("No errors sent!")
	}
}

func TestMPIMissing(t *testing.T) {
	s := Init()
	s.SkipLog = true
	gh := &testGh{}
	s.rc = &testRc{missing: true}
	s.gh = gh

	s.alertForMisorderedMPI(context.Background())

	if gh.count != 0 {
		t.Errorf("Errors have been sent!")
	}
}
