package multilane_test

import (
	. "github.com/andy2046/gopie/pkg/multilane"
	"go.uber.org/goleak"
	"runtime"
	"sync"
	"testing"
	"time"
)

var config = Config{
	LaneWidth: 8,
	QueueSize: 1024,
}

func BenchmarkMultiLane_Blocking(b *testing.B) {
	var (
		m                            = New(config)
		wgPut, wgGet                 sync.WaitGroup
		concurrentGet, concurrentPut = 8, 8
	)
	wgPut.Add(concurrentPut)
	wgGet.Add(concurrentGet)
	type T struct {
		i int
	}
	var v = T{5}
	b.ResetTimer()
	for c := 0; c < concurrentPut; c++ {
		go func(n int) {
			// runtime.LockOSThread()
			for i := 0; i < n; i++ {
				m.Put(&v)
			}
			wgPut.Done()
		}(b.N/concurrentPut + 1)
	}
	for c := 0; c < concurrentGet; c++ {
		go func(n int) {
			// runtime.LockOSThread()
			var v *T
			for i := 0; m.Get(&v); i++ {
				_ = *v
			}
			wgGet.Done()
		}(b.N/concurrentGet + 1)
	}

	wgPut.Wait()
	m.Close()
	wgGet.Wait()
}

func BenchmarkMultiLane_BlockingLane(b *testing.B) {
	var (
		m                            = New(config)
		wgPut, wgGet                 sync.WaitGroup
		concurrentGet, concurrentPut = 8, 8
	)
	wgPut.Add(concurrentPut)
	wgGet.Add(concurrentGet)
	type T struct {
		i int
	}
	var v = T{5}
	b.ResetTimer()
	for c := 0; c < concurrentPut; c++ {
		go func(n int, lane uint32) {
			// runtime.LockOSThread()
			for i := 0; i < n; i++ {
				m.PutLane(lane, &v)
			}
			wgPut.Done()
		}(b.N/concurrentPut+1, uint32(c))
	}
	for c := 0; c < concurrentGet; c++ {
		go func(n int, lane uint32) {
			// runtime.LockOSThread()
			var v *T
			for i := 0; m.GetLane(lane, &v); i++ {
				_ = *v
			}
			wgGet.Done()
		}(b.N/concurrentGet+1, uint32(c))
	}

	wgPut.Wait()
	m.Close()
	wgGet.Wait()
}

func BenchmarkMultiLane_Blocking2(b *testing.B) {
	var (
		m                            = New(config)
		wgPut, wgGet                 sync.WaitGroup
		concurrentGet, concurrentPut = 8, 8
	)
	wgPut.Add(concurrentPut)
	wgGet.Add(concurrentGet)
	type T struct {
		i int
	}
	p := runtime.GOMAXPROCS(concurrentGet + concurrentPut)
	var v = T{5}
	b.ResetTimer()
	for c := 0; c < concurrentPut; c++ {
		go func(n int) {
			for i := 0; i < n; i++ {
				m.Put(&v)
			}
			wgPut.Done()
		}(b.N/concurrentPut + 1)
	}
	for c := 0; c < concurrentGet; c++ {
		go func(n int) {
			var v *T
			for i := 0; m.Get(&v); i++ {
				_ = *v
			}
			wgGet.Done()
		}(b.N/concurrentGet + 1)
	}

	wgPut.Wait()
	m.Close()
	wgGet.Wait()
	runtime.GOMAXPROCS(p)
}

func BenchmarkMultiLane_BlockingLane2(b *testing.B) {
	var (
		m                            = New(config)
		wgPut, wgGet                 sync.WaitGroup
		concurrentGet, concurrentPut = 8, 8
	)
	wgPut.Add(concurrentPut)
	wgGet.Add(concurrentGet)
	type T struct {
		i int
	}
	p := runtime.GOMAXPROCS(concurrentGet + concurrentPut)
	var v = T{5}
	b.ResetTimer()
	for c := 0; c < concurrentPut; c++ {
		go func(n int, lane uint32) {
			for i := 0; i < n; i++ {
				m.PutLane(lane, &v)
			}
			wgPut.Done()
		}(b.N/concurrentPut+1, uint32(c))
	}
	for c := 0; c < concurrentGet; c++ {
		go func(n int, lane uint32) {
			var v *T
			for i := 0; m.GetLane(lane, &v); i++ {
				_ = *v
			}
			wgGet.Done()
		}(b.N/concurrentGet+1, uint32(c))
	}

	wgPut.Wait()
	m.Close()
	wgGet.Wait()
	runtime.GOMAXPROCS(p)
}

func TestMultiLane_Blocking(t *testing.T) {
	defer goleak.VerifyNone(t)
	const N = 1000
	var m = New(Config{
		LaneWidth: 2,
		QueueSize: 8,
	})
	var wg sync.WaitGroup
	wg.Add(2)
	t1 := time.Now()
	go func(n int) {
		defer wg.Done()
		for i := 0; i < n; i++ {
			m.Put(int64(i))
			time.Sleep(1 * time.Microsecond)
		}
		m.Close()
	}(N)
	go func(n int) {
		defer wg.Done()
		var v *int
		for i := 0; m.Get(&v); i++ {
			if i != *v {
				t.Fatalf("Expected %d, but got %d", i, *v)
				panic(i)
			}
		}
	}(N)
	wg.Wait()
	t.Log(time.Since(t1))
}
