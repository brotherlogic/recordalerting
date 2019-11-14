package main

import (
	"testing"

	"golang.org/x/net/context"
)

func TestNoSale(t *testing.T) {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
	gh := &testGh{}
	s.rc = &testRc{}
	s.gh = gh

	s.alertForMissingSaleID(context.Background())

	if gh.count == 0 {
		t.Errorf("No errors sent!: %v", gh.count)
	}
}

func TestPurgatory(t *testing.T) {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
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
	s.SkipIssue = true
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
	s.SkipIssue = true
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
	s.SkipIssue = true
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
	s.SkipIssue = true
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
	s.SkipIssue = true
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
	s.SkipIssue = true
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

func TestOldListeningPileFailRO(t *testing.T) {
	s := InitTest()
	s.ro = &testRo{fail: true}

	err := s.alertForOldListeningPileRecord(context.Background())

	if err == nil {
		t.Errorf("Did not error")
	}
}

func TestOldListeningPileFailRC(t *testing.T) {
	s := InitTest()
	s.rc = &testRc{fail: true}

	err := s.alertForOldListeningPileRecord(context.Background())

	if err == nil {
		t.Errorf("Did not error")
	}
}

func TestInvalid(t *testing.T) {
	s := InitTest()

	s.validateRecords(context.Background())

	if s.invalidRecords == 0 {
		t.Errorf("No invalidation")
	}
}

func TestInvalidFail(t *testing.T) {
	s := InitTest()
	s.rc = &testRc{failAll: true}

	err := s.validateRecords(context.Background())

	if err == nil {
		t.Errorf("Did not error")
	}
}

func TestInvalidFailQuery(t *testing.T) {
	s := InitTest()
	s.rc = &testRc{fail: true}

	err := s.validateRecords(context.Background())

	if err == nil {
		t.Errorf("Did not error")
	}
}
