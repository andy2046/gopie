package batch

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	var (
		callcount int32
		wg        sync.WaitGroup
		loop      = 2
	)
	wg.Add(loop)

	b := New(1*time.Second, 10, func(batch []interface{}) {
		atomic.AddInt32(&callcount, 1)
		if len(batch) != loop {
			t.Errorf("batch size got %v want %v", len(batch), loop)
		}
	})
	defer b.Close()

	for i := 0; i < loop; i++ {
		go func(i int) {
			b.Batch(i)
			wg.Done()
		}(i)
	}

	wg.Wait()
	time.Sleep(2 * time.Second)

	cc := atomic.LoadInt32(&callcount)
	if cc != 1 {
		t.Errorf("callcount got %v want %v", cc, 1)
	}
}

func TestCount(t *testing.T) {
	var (
		callcount int32
		wg        sync.WaitGroup
		loop      = 10
	)
	wg.Add(loop)

	b := New(10*time.Second, 10, func(batch []interface{}) {
		atomic.AddInt32(&callcount, 1)
		if len(batch) != loop {
			t.Errorf("batch size got %v want %v", len(batch), loop)
		}
	})
	defer b.Close()

	for i := 0; i < loop; i++ {
		go func(i int) {
			b.Batch(i)
			wg.Done()
		}(i)
	}

	wg.Wait()
	time.Sleep(1 * time.Second)

	cc := atomic.LoadInt32(&callcount)
	if cc != 1 {
		t.Errorf("callcount got %v want %v", cc, 1)
	}
}

func TestLoad(t *testing.T) {
	var (
		wg   sync.WaitGroup
		loop = 10000
	)
	wg.Add(loop)

	b := New(1*time.Second, 100, func(batch []interface{}) {
	})
	defer b.Close()

	for i := 0; i < loop; i++ {
		go func(i int) {
			b.Batch(i)
			wg.Done()
		}(i)
	}

	wg.Wait()
	time.Sleep(1 * time.Second)
}
