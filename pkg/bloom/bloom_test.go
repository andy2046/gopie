package bloom_test

import (
	"encoding/binary"
	"testing"

	"github.com/andy2046/gopie/pkg/bloom"
)

func TestBasic(t *testing.T) {
	f := bloom.New(1000, 4)
	e1 := []byte("Boss")
	e2 := []byte("Joke")
	e3 := []byte("Emotion")
	f.Add(e1)
	e3b := f.Exist(e3)
	e1a := f.Exist(e1)
	e2a := f.Exist(e2)
	f.Add(e3)
	e3a := f.Exist(e3)
	if !e1a {
		t.Errorf("%q should Exist.", e1)
	}
	if e2a {
		t.Errorf("%q should not Exist.", e2)
	}
	if e3b {
		t.Errorf("%q should not Exist the first time we check.", e3)
	}
	if !e3a {
		t.Errorf("%q should Exist the second time we check.", e3)
	}
}

func TestUint(t *testing.T) {
	f := bloom.New(1000, 4)
	n1 := make([]byte, 4)
	n2 := make([]byte, 4)
	n3 := make([]byte, 4)
	n4 := make([]byte, 4)
	binary.BigEndian.PutUint32(n1, 100)
	binary.BigEndian.PutUint32(n2, 101)
	binary.BigEndian.PutUint32(n3, 102)
	binary.BigEndian.PutUint32(n4, 103)
	f.Add(n1)
	n3b := f.Exist(n3)
	n1a := f.Exist(n1)
	n2a := f.Exist(n2)
	f.Add(n3)
	n3a := f.Exist(n3)
	n4a := f.Exist(n4)
	if !n1a {
		t.Errorf("%q should Exist.", n1)
	}
	if n2a {
		t.Errorf("%q should not Exist.", n2)
	}
	if n3b {
		t.Errorf("%q should not Exist the first time we check.", n3)
	}
	if !n3a {
		t.Errorf("%q should Exist the second time we check.", n3)
	}
	if n4a {
		t.Errorf("%q should not Exist.", n4)
	}
}

func TestString(t *testing.T) {
	f := bloom.NewGuess(1000, 0.001)
	s1 := "Filter"
	s2 := "is"
	s3 := "in"
	s4 := "bloom"
	f.AddString(s1)
	s3b := f.ExistString(s3)
	s1a := f.ExistString(s1)
	s2a := f.ExistString(s2)
	f.AddString(s3)
	s3a := f.ExistString(s3)
	s4a := f.ExistString(s4)
	if !s1a {
		t.Errorf("%q should Exist.", s1)
	}
	if s2a {
		t.Errorf("%q should not Exist.", s2)
	}
	if s3b {
		t.Errorf("%q should not Exist the first time we check.", s3)
	}
	if !s3a {
		t.Errorf("%q should Exist the second time we check.", s3)
	}
	if s4a {
		t.Errorf("%q should not Exist.", s4)
	}
}

func TestGuessFalsePositive(t *testing.T) {
	n, p := uint64(100000), float64(0.001)
	m, k := bloom.Guess(n, p)
	f := bloom.NewGuess(n, p)
	fp := f.GuessFalsePositive(n)
	t.Logf("m=%v k=%v n=%v p=%v fp=%v", m, k, n, p, fp)
	if fp > p {
		t.Errorf("False Positive too high")
	}
}

func TestM(t *testing.T) {
	f := bloom.New(1000, 4)
	if f.M() != 1024 {
		t.Errorf("M() %v is not correct", f.M())
	}
}

func TestK(t *testing.T) {
	f := bloom.New(1000, 4)
	if f.K() != 4 {
		t.Errorf("K() %v is not correct", f.K())
	}
}

func BenchmarkAddExist(b *testing.B) {
	f := bloom.NewGuess(uint64(b.N), 0.0001)
	key := make([]byte, 8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		binary.BigEndian.PutUint64(key, uint64(i))
		f.Add(key)
		f.Exist(key)
	}
}
