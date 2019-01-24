package bloom_test

import (
	"encoding/binary"
	"testing"

	"github.com/andy2046/gopie/pkg/bloom"
)

func TestScaleBasic(t *testing.T) {
	f := bloom.NewS(0.01)
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

func TestScaleUint(t *testing.T) {
	f := bloom.NewS(0.01)
	n1 := make([]byte, 4)
	n2 := make([]byte, 4)
	n3 := make([]byte, 4)
	n4 := make([]byte, 4)
	binary.BigEndian.PutUint32(n1, 10000)
	binary.BigEndian.PutUint32(n2, 10001)
	binary.BigEndian.PutUint32(n3, 10002)
	binary.BigEndian.PutUint32(n4, 10003)
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

	for i := uint32(1); i < 1000; i++ {
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, i)
		f.Add(buf)
	}
	count := 0
	for i := uint32(1000); i < 4000; i++ {
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, i)
		if f.Exist(buf) {
			count++
		}
	}
	t.Logf("FP rate is %v", f.FalsePositive())
	if f.FalsePositive() > 0.01 {
		t.Errorf("False Positive rate should not be > 0.01")
	}

	if count > 1 {
		t.Errorf("Actual FP %d greater than expected FP", count)
	}
}

func TestScaleString(t *testing.T) {
	f := bloom.NewS(0.01)
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

func TestScaleM(t *testing.T) {
	f := bloom.NewS(0.01)
	if f.M() < 512 {
		t.Errorf("M() %v is not correct", f.M())
	}
}

func TestScaleK(t *testing.T) {
	f := bloom.NewS(0.01)
	if f.K() < 3 {
		t.Errorf("K() %v is not correct", f.K())
	}
}

func BenchmarkScaleAddExist(b *testing.B) {
	f := bloom.NewSGuess(uint64(b.N), 0.0001, 0.9)
	key := make([]byte, 8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		binary.BigEndian.PutUint64(key, uint64(i))
		f.Add(key)
		f.Exist(key)
	}
}
