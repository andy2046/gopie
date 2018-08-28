package barrier

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	num := 10
	b := New(num)
	n, nWaiting := b.N(), b.NWaiting()
	isBroken := b.IsBroken()

	if n != num {
		t.Error("number of barrier members, expected", num, ", got", n)
	}
	if nWaiting != 0 {
		t.Error("number of barrier members waiting, expected", 0, ", got", nWaiting)
	}
	if isBroken != false {
		t.Error("barrier isBroken, expected", false, ", got", isBroken)
	}

	defer func() {
		if recover() == nil {
			t.Error("panic expected")
		}
	}()
	_ = New(0)
	_ = New(-1)
}

func TestAwaitOnce(t *testing.T) {
	num := 100
	b := New(num)
	ctx := context.Background()

	wg := sync.WaitGroup{}
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			err := b.Await(ctx)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	n, nWaiting := b.N(), b.NWaiting()
	isBroken := b.IsBroken()

	if n != num {
		t.Error("number of barrier members, expected", num, ", got", n)
	}
	if nWaiting != 0 {
		t.Error("number of barrier members waiting, expected", 0, ", got", nWaiting)
	}
	if isBroken != false {
		t.Error("barrier isBroken, expected", false, ", got", isBroken)
	}
}

func TestAwaitMany(t *testing.T) {
	num := 100
	m := 1000
	b := New(num)
	ctx := context.Background()
	wg := sync.WaitGroup{}

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < m; j++ {
				err := b.Await(ctx)
				if err != nil {
					panic(err)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	n, nWaiting := b.N(), b.NWaiting()
	isBroken := b.IsBroken()

	if n != num {
		t.Error("number of barrier members, expected", num, ", got", n)
	}
	if nWaiting != 0 {
		t.Error("number of barrier members waiting, expected", 0, ", got", nWaiting)
	}
	if isBroken != false {
		t.Error("barrier isBroken, expected", false, ", got", isBroken)
	}
}

func TestAwaitTooMany(t *testing.T) {
	num := 100
	m := 1000
	b := New(1)
	ctx := context.Background()
	wg := sync.WaitGroup{}

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < m; j++ {
				err := b.Await(ctx)
				if err != nil {
					panic(err)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	n, nWaiting := b.N(), b.NWaiting()
	isBroken := b.IsBroken()

	if n != 1 {
		t.Error("number of barrier members, expected", num, ", got", n)
	}
	if nWaiting != 0 {
		t.Error("number of barrier members waiting, expected", 0, ", got", nWaiting)
	}
	if isBroken != false {
		t.Error("barrier isBroken, expected", false, ", got", isBroken)
	}
}

func TestReset(t *testing.T) {
	num := 100
	b := New(num + 1) // members are more than goroutines so all goroutines will wait
	ctx := context.Background()

	go func() {
		time.Sleep(30 * time.Millisecond)
		b.Reset()
	}()

	wg := sync.WaitGroup{}
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			err := b.Await(ctx)
			if err != ErrBroken {
				panic(err)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	n, nWaiting := b.N(), b.NWaiting()
	isBroken := b.IsBroken()

	if n != num+1 {
		t.Error("number of barrier members, expected", num+1, ", got", n)
	}
	if nWaiting != 0 {
		t.Error("number of barrier members waiting, expected", 0, ", got", nWaiting)
	}
	if isBroken != false {
		t.Error("barrier isBroken, expected", false, ", got", isBroken)
	}
}

func TestAwaitOnceCtxDone(t *testing.T) {
	num := 100
	b := New(num + 1) // members are more than goroutines so all goroutines will wait
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	var deadlineCount, brokenBarrierCount int32
	wg := sync.WaitGroup{}

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			err := b.Await(ctx)
			if err == context.DeadlineExceeded {
				atomic.AddInt32(&deadlineCount, 1)
			} else if err == ErrBroken {
				atomic.AddInt32(&brokenBarrierCount, 1)
			} else {
				panic("must be either context.DeadlineExceeded or ErrBroken error")
			}
			wg.Done()
		}()
	}

	wg.Wait()

	n, nWaiting := b.N(), b.NWaiting()
	isBroken := b.IsBroken()

	if n != num+1 {
		t.Error("number of barrier members, expected", num+1, ", got", n)
	}
	if nWaiting != 100 {
		t.Error("number of barrier members waiting, expected", 100, ", got", nWaiting)
	}
	if isBroken != true {
		t.Error("barrier isBroken, expected", true, ", got", isBroken)
	}

	if deadlineCount == 0 {
		t.Error("number of context.DeadlineExceeded errors must be more than 0, got", deadlineCount)
	}
	if deadlineCount+brokenBarrierCount != int32(num) {
		t.Error("number of context.DeadlineExceeded and ErrBroken errors, expected", num, ", got", deadlineCount+brokenBarrierCount)
	}
}

func TestAwaitManyCtxDone(t *testing.T) {
	num := 100
	b := New(num)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	wg := sync.WaitGroup{}

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			for {
				err := b.Await(ctx)
				if err != nil {
					if err != context.DeadlineExceeded && err != ErrBroken {
						panic("must be either context.DeadlineExceeded or ErrBroken error")
					}
					break
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	n := b.N()
	isBroken := b.IsBroken()

	if n != num {
		t.Error("number of barrier members, expected", num, ", got", n)
	}
	if isBroken != true {
		t.Error("barrier isBroken, expected", true, ", got", isBroken)
	}
}
