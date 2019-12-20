package sequence_test

import (
	. "github.com/andy2046/gopie/pkg/sequence"
	"testing"
)

type tester struct {
	T   *testing.T
	seq Sequencer
}

func (t tester) Next() uint64 {
	id, err := t.seq.Next()
	if err != nil {
		t.T.Fatal("error:", err)
	}
	return id
}

func (t tester) NextN(n int) uint64 {
	id, err := t.seq.NextN(n)
	if err != nil {
		t.T.Fatal("error:", err)
	}
	return id
}

func TestSequence_Next(t *testing.T) {
	s := NewMemSeq(0)
	test := tester{t, s}

	loop := uint64(3)
	for i := uint64(0); i < loop; i++ {
		test.Next()
	}
	n := test.Next()
	if n != loop {
		t.Errorf("error: expected %d got %d", loop, n)
	}

	n = test.NextN(1000)
	if n != loop+1 {
		t.Errorf("error: expected %d got %d", loop+1, n)
	}

	n = test.Next()
	if n != loop+1001 {
		t.Errorf("error: expected %d got %d", loop+1001, n)
	}
}

func TestSequence_Parallel(t *testing.T) {
	s := NewMemSeq(0)
	test := tester{t, s}

	consumer := make(chan uint64)
	const numID = 10000
	generate := func() {
		for i := 0; i < numID; i++ {
			consumer <- test.Next()
		}
	}
	const numGenerator = 10
	for i := 0; i < numGenerator; i++ {
		go generate()
	}
	set := make(map[uint64]bool)
	for i := 0; i < numID*numGenerator; i++ {
		id := <-consumer
		if _, exist := set[id]; exist {
			t.Fatal("error: duplicated")
		}
		set[id] = true
	}
}
