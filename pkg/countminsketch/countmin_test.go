package countminsketch

import (
	"strconv"
	"testing"
)

func TestCount(t *testing.T) {
	cms, _ := NewGuess(0.001, 0.99)

	for i := 0; i < 100; i++ {
		cms.Add([]byte(strconv.Itoa(i)))
	}

	if count := cms.Count(); count != 100 {
		t.Errorf("expected 100, got %d", count)
	}
}

func TestEstimate(t *testing.T) {
	cms, _ := NewGuess(0.001, 0.99)
	cms.Add([]byte(`a`))
	cms.Add([]byte(`b`), 1)
	cms.Add([]byte(`c`), 1)
	cms.Add([]byte(`b`), 1)

	if count := cms.Estimate([]byte(`a`)); count != 1 {
		t.Errorf("expected 1, got %d", count)
	}

	if count := cms.Estimate([]byte(`b`)); count != 2 {
		t.Errorf("expected 2, got %d", count)
	}

	if count := cms.Estimate([]byte(`c`)); count != 1 {
		t.Errorf("expected 1, got %d", count)
	}

	if count := cms.Estimate([]byte(`x`)); count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}

func TestMerge(t *testing.T) {
	cms, _ := NewGuess(0.001, 0.99)
	cms.Add([]byte(`a`))
	cms.Add([]byte(`b`), 1)
	cms.Add([]byte(`c`), 1)
	cms.Add([]byte(`b`), 1)
	cms.Add([]byte(`d`), 1)

	other, _ := NewGuess(0.001, 0.99)
	other.Add([]byte(`b`), 1)
	other.Add([]byte(`c`), 1)
	other.Add([]byte(`b`), 1)

	if err := cms.Merge(other); err != nil {
		t.Error(err)
	}

	if count := cms.Estimate([]byte(`a`)); count != 1 {
		t.Errorf("expected 1, got %d", count)
	}

	if count := cms.Estimate([]byte(`b`)); count != 4 {
		t.Errorf("expected 4, got %d", count)
	}

	if count := cms.Estimate([]byte(`c`)); count != 2 {
		t.Errorf("expected 2, got %d", count)
	}

	if count := cms.Estimate([]byte(`d`)); count != 1 {
		t.Errorf("expected 1, got %d", count)
	}

	if count := cms.Estimate([]byte(`x`)); count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}

func TestReset(t *testing.T) {
	cms, _ := NewGuess(0.001, 0.99)
	cms.Add([]byte(`a`))
	cms.Add([]byte(`b`), 1)
	cms.Add([]byte(`c`), 1)
	cms.Add([]byte(`b`), 1)
	cms.Add([]byte(`d`), 1)

	cms.Reset()

	for i := uint(0); i < cms.depth; i++ {
		for j := uint(0); j < cms.width; j++ {
			if x := cms.matrix[i][j]; x != 0 {
				t.Errorf("expected matrix to be empty, got %d", x)
			}
		}
	}
}

func BenchmarkAdd(b *testing.B) {
	b.StopTimer()
	cms, _ := NewGuess(0.001, 0.99)
	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = []byte(strconv.Itoa(i))
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		cms.Add(data[n])
	}
}

func BenchmarkCount(b *testing.B) {
	b.StopTimer()
	cms, _ := NewGuess(0.001, 0.99)
	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = []byte(strconv.Itoa(i))
		cms.Add([]byte(strconv.Itoa(i)))
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		cms.Estimate(data[n])
	}
}
