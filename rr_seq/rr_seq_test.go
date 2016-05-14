package rrs

import (
	"fmt"
	"testing"
	"time"

	"github.com/facebookgo/ensure"
	"github.com/teh-cmc/seq"
)

// -----------------------------------------------------------------------------

// NOTE: run these tests with `go test -race -cpu 1,8,32`

func TestRRSeq_New_BufSize(t *testing.T) {
	var rrseq *RRSeq
	var err error

	rrseq, err = NewRRSeq(testingRRSeqDefaultName, -42, testingRRServerAddrs...)
	if err != nil {
		t.Fatal(err)
	}
	ensure.DeepEqual(t, cap(rrseq.ids), 0)

	rrseq, err = NewRRSeq(testingRRSeqDefaultName, 0, testingRRServerAddrs...)
	if err != nil {
		t.Fatal(err)
	}
	ensure.DeepEqual(t, cap(rrseq.ids), 0)

	rrseq, err = NewRRSeq(testingRRSeqDefaultName, 1, testingRRServerAddrs...)
	if err != nil {
		t.Fatal(err)
	}
	ensure.DeepEqual(t, cap(rrseq.ids), 1)

	rrseq, err = NewRRSeq(testingRRSeqDefaultName, 1e6, testingRRServerAddrs...)
	if err != nil {
		t.Fatal(err)
	}
	ensure.DeepEqual(t, cap(rrseq.ids), int(1e6))
}

func TestRRSeq_FirstID(t *testing.T) {
	rrseq, err := NewRRSeq("TestRRSeq_FirstID", 1e2, testingRRServerAddrs...)
	if err != nil {
		t.Fatal(err)
	}
	ensure.DeepEqual(t, <-rrseq.GetStream(), seq.ID(1))
}

// -----------------------------------------------------------------------------

func testRRSeq_SingleClient(bufSize int, t *testing.T) {
	name := fmt.Sprintf("testRRSeq_SingleClient(%d)", bufSize)
	s, err := NewRRSeq(name, bufSize, testingRRServerAddrs...)
	if err != nil {
		t.Fatal(err)
	}
	lastID := seq.ID(0)

	go func() {
		<-time.After(time.Millisecond * 250)
		_ = s.Close()
	}()

	for id := range s.GetStream() {
		ensure.DeepEqual(t, id, lastID+1)
		lastID = id
	}
}

func TestRRSeq_BufSize0_SingleClient(t *testing.T) {
	testRRSeq_SingleClient(0, t)
}

func TestRRSeq_BufSize1_SingleClient(t *testing.T) {
	testRRSeq_SingleClient(1, t)
}

func TestRRSeq_BufSize1024_SingleClient(t *testing.T) {
	testRRSeq_SingleClient(1024, t)
}
