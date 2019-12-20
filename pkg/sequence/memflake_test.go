package sequence_test

import (
	. "github.com/andy2046/gopie/pkg/sequence"
	"testing"
	"time"
)

func TestIceflake_Next(t *testing.T) {
	s := NewMemFlake(time.Time{}, 0, 0)
	if s == nil {
		t.Fatal("error: nil Iceflake")
	}
	test := tester{t, s}

	n0 := test.Next()
	for i := 0; i < 3; i++ {
		test.Next()
	}
	n := test.Next()
	if n <= n0 {
		t.Errorf("error: expected %d greater than %d", n, n0)
	}

	n0 = test.NextN(1000)
	if n0 <= n {
		t.Errorf("error: expected %d greater than %d", n0, n)
	}

	n = test.Next()
	if n <= n0 {
		t.Errorf("error: expected %d greater than %d", n, n0)
	}
}

func TestIceflake_Parallel(t *testing.T) {
	s := NewMemFlake(time.Time{}, 0, 0)
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
