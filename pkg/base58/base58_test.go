package base58_test

import (
	. "github.com/andy2046/gopie/pkg/base58"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	testcases := map[uint64]string{
		0:              "1",
		32:             "Z",
		57:             "z",
		math.MaxUint8:  "5Q",
		math.MaxUint16: "LUv",
		math.MaxUint32: "7YXq9G",
		math.MaxUint64: "jpXCZedGfVQ",
	}
	for k, v := range testcases {
		r := Encode(k)
		if v != r {
			t.Errorf("expected %s got %s", v, r)
		}
	}
}

func TestDecode(t *testing.T) {
	_, err := Decode("")
	if err == nil {
		t.Fail()
	}
	_, err = Decode("0")
	if err == nil {
		t.Fail()
	}

	testcases := map[uint64]string{
		0:              "1",
		32:             "Z",
		57:             "z",
		math.MaxUint8:  "5Q",
		math.MaxUint16: "LUv",
		math.MaxUint32: "7YXq9G",
		math.MaxUint64: "jpXCZedGfVQ",
	}
	for k, v := range testcases {
		r, err := Decode(v)
		if err != nil {
			t.Fatal(err)
		}
		if k != r {
			t.Errorf("expected %d got %d", k, r)
		}
	}
}

func TestReverse(t *testing.T) {
	testcases := map[string]string{
		"":    "",
		"1":   "1",
		"ABC": "CBA",
		"xyz": "zyx",
	}
	for k, v := range testcases {
		r := []byte(k)
		Reverse(r)
		if v != string(r) {
			t.Errorf("expected %s got %s", v, string(r))
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(uint64(s.Int63()))
	}
}

func BenchmarkDecode(b *testing.B) {
	arr := []string{"1", "Z", "z", "5Q", "LUv", "7YXq9G", "jpXCZedGfVQ"}
	ln := len(arr)
	s := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Decode(arr[s.Intn(ln)])
		if err != nil {
			b.Fatal(err)
		}
	}
}
