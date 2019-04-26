package main

import (
	"testing"

	"golang.org/x/net/context"
)

func TestNoSale(t *testing.T) {
	s := Init()
	s.SkipLog = true
	gh := &testGh{}
	s.rc = &testRc{}
	s.gh = gh

	s.alertForMissingSaleID(context.Background())

	if gh.count == 0 {
		t.Errorf("No errors sent!")
	}
}

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

func InitTest() *Server {
	s := Init()
	s.SkipLog = true
	s.gh = &testGh{}
	s.rc = &testRc{}
	s.ro = &testRo{}

	return s
}

func TestOldListeningBox(t *testing.T) {
	s := InitTest()
	gh := &testGh{}
	s.gh = gh

	s.alertForOldListeningBoxRecord(context.Background())

	if gh.count == 0 {
		t.Errorf("No errors have been sent!")
	}
}

func TestOldListeningPile(t *testing.T) {
	s := InitTest()
	gh := &testGh{}
	s.gh = gh

	s.alertForOldListeningPileRecord(context.Background())

	if gh.count == 0 {
		t.Errorf("No errors have been sent!")
	}
}
