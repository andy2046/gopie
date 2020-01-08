package spsc_test

import (
	. "github.com/andy2046/gopie/pkg/spsc"
	"runtime"
	"sync"
	"testing"
	"time"
)

func BenchmarkSPSC_NonBlocking(b *testing.B) {
	var sp = New(8192)
	var wg sync.WaitGroup
	wg.Add(2)
	type T struct {
		i int
	}
	var v = T{5}
	b.ResetTimer()
	go func(n int) {
		runtime.LockOSThread()
		for i := 0; i < n; i++ {
			for !sp.Offer(&v) {
				runtime.Gosched()
			}
		}
		wg.Done()
	}(b.N)

	go func(n int) {
		runtime.LockOSThread()
		var v *T
		for i := 0; i < n; i++ {
			for !sp.Poll(&v) {
				runtime.Gosched()
			}
			_ = *v
		}
		wg.Done()
	}(b.N)

	wg.Wait()
}

func BenchmarkSPSC_NonBlocking2(b *testing.B) {
	var sp = New(8192)
	var wg sync.WaitGroup
	wg.Add(2)
	type T struct {
		i int
	}
	p := runtime.GOMAXPROCS(1)
	var v = T{5}
	b.ResetTimer()
	go func(n int) {
		for i := 0; i < n; i++ {
			for !sp.Offer(&v) {
				runtime.Gosched()
			}
		}
		wg.Done()
	}(b.N)

	go func(n int) {
		var v *T
		for i := 0; i < n; i++ {
			for !sp.Poll(&v) {
				runtime.Gosched()
			}
			_ = *v
		}
		wg.Done()
	}(b.N)

	wg.Wait()
	runtime.GOMAXPROCS(p)
}

func BenchmarkSPSC_Blocking(b *testing.B) {
	var sp = New(8192)
	var wg sync.WaitGroup
	wg.Add(2)
	type T struct {
		i int
	}
	var v = T{5}
	b.ResetTimer()
	go func(n int) {
		runtime.LockOSThread()
		for i := 0; i < n; i++ {
			sp.Put(&v)
		}
		sp.Close()
		wg.Done()
	}(b.N)

	go func(n int) {
		runtime.LockOSThread()
		var v *T
		for i := 0; sp.Get(&v); i++ {
			_ = *v
		}
		wg.Done()
	}(b.N)

	wg.Wait()
}

func BenchmarkSPSC_Blocking2(b *testing.B) {
	var sp = New(8192)
	var wg sync.WaitGroup
	wg.Add(2)
	type T struct {
		i int
	}
	p := runtime.GOMAXPROCS(1)
	var v = T{5}
	b.ResetTimer()
	go func(n int) {
		for i := 0; i < b.N; i++ {
			sp.Put(&v)
		}
		sp.Close()
		wg.Done()
	}(b.N)

	go func(n int) {
		var v *T
		for i := 0; sp.Get(&v); i++ {
			_ = *v
		}
		wg.Done()
	}(b.N)

	wg.Wait()
	runtime.GOMAXPROCS(p)
}

func TestSPSC_NonBlocking(t *testing.T) {
	const N = 1000
	var sp = New(4)
	var wg sync.WaitGroup
	wg.Add(2)
	t1 := time.Now()
	go func(n int) {
		defer wg.Done()
		for i := 0; i < n; i++ {
			for !sp.Offer(int64(i)) {
				runtime.Gosched()
			}
			time.Sleep(1 * time.Microsecond)
		}
	}(N)
	go func(n int) {
		defer wg.Done()
		var v *int
		for i := 0; i < n; i++ {
			for !sp.Poll(&v) {
				runtime.Gosched()
			}
			if i != *v {
				t.Fatalf("Expected %d, but got %d", i, *v)
				panic(i)
			}
		}
	}(N)
	wg.Wait()
	t.Log(time.Since(t1))
}

func TestSPSC_Blocking(t *testing.T) {
	const N = 1000
	var sp = New(4)
	var wg sync.WaitGroup
	wg.Add(2)
	t1 := time.Now()
	go func(n int) {
		defer wg.Done()
		for i := 0; i < n; i++ {
			sp.Put(int64(i))
			time.Sleep(1 * time.Microsecond)
		}
		sp.Close()
	}(N)
	go func(n int) {
		defer wg.Done()
		var v *int
		for i := 0; sp.Get(&v); i++ {
			if i != *v {
				t.Fatalf("Expected %d, but got %d", i, *v)
				panic(i)
			}
		}
	}(N)
	wg.Wait()
	t.Log(time.Since(t1))
}
