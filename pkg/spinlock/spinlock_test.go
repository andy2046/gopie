package spinlock_test

import (
	. "github.com/andy2046/gopie/pkg/spinlock"
	"sync"
	"testing"
	"time"
)

func TestLocker(t *testing.T) {
	threads, loops, count := 8, 1000000, 0
	var wg sync.WaitGroup
	wg.Add(threads)

	l := New()
	start := time.Now()
	for i := 0; i < threads; i++ {
		go func() {
			for i := 0; i < loops; i++ {
				l.Lock()
				count++
				l.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	duration := time.Since(start)
	t.Logf("duration: %4.2f Seconds\n", duration.Seconds())
	if count != threads*loops {
		t.Errorf("expected %d got %d", threads*loops, count)
	}
}

func TestNoCopy(t *testing.T) {
	// go vet fails
	var l1 Locker
	l2 := l1
	var l3 = l1
	l2 = l1
	_, _ = l2, l3
	t.Log("go vet fails here")
}
