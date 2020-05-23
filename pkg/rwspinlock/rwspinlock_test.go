package rwspinlock

import (
	"sync"
	"testing"
)

func TestRLocker(t *testing.T) {
	threads, loops := 8, 1000
	var wg sync.WaitGroup
	wg.Add(threads)
	l := New().RLocker()

	for i := 0; i < threads; i++ {
		go func() {
			for i := 0; i < loops; i++ {
				l.Lock()
				l.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestLock(t *testing.T) {
	threads, loops, count := 8, 1000, 0
	var wg sync.WaitGroup
	wg.Add(threads)
	l := New()

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

	if count != threads*loops {
		t.Errorf("expected %d got %d", threads*loops, count)
	}
}

func TestRLock(t *testing.T) {
	threads, loops := 8, 1000
	var wg sync.WaitGroup
	wg.Add(threads)
	l := New()

	for i := 0; i < threads; i++ {
		go func() {
			for i := 0; i < loops; i++ {
				l.RLock()
				l.RUnlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestRWLock(t *testing.T) {
	threadsR, threadsW, loopsR, loopsW := 8, 2, 1000, 100
	var wg sync.WaitGroup
	wg.Add(threadsR + threadsW)
	l := New()

	for i := 0; i < threadsR; i++ {
		go func() {
			for i := 0; i < loopsR; i++ {
				l.RLock()
				l.RUnlock()
			}
			wg.Done()
		}()
	}

	for i := 0; i < threadsW; i++ {
		go func() {
			for i := 0; i < loopsW; i++ {
				l.Lock()
				l.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkRLock(b *testing.B) {
	l := New()
	for n := 0; n < b.N; n++ {
		l.RLock()
		l.RUnlock()
	}
}

func BenchmarkLock(b *testing.B) {
	l := New()
	for n := 0; n < b.N; n++ {
		l.Lock()
		l.Unlock()
	}
}

func BenchmarkRWLock(b *testing.B) {
	l := New()
	for n := 0; n < b.N; n++ {
		if n&3 == 0 {
			go func() {
				l.Lock()
				l.Unlock()
			}()
		}
		l.RLock()
		l.RUnlock()
	}
}

func BenchmarkRMutex(b *testing.B) {
	var l sync.RWMutex
	for n := 0; n < b.N; n++ {
		l.RLock()
		l.RUnlock()
	}
}

func BenchmarkMutex(b *testing.B) {
	var l sync.RWMutex
	for n := 0; n < b.N; n++ {
		l.Lock()
		l.Unlock()
	}
}

func BenchmarkRWMutex(b *testing.B) {
	var l sync.RWMutex
	for n := 0; n < b.N; n++ {
		if n&3 == 0 {
			go func() {
				l.Lock()
				l.Unlock()
			}()
		}
		l.RLock()
		l.RUnlock()
	}
}
