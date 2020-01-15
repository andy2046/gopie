package semaphore

import (
	"context"
	"testing"
	"time"
)

func work() {
	time.Sleep(100 * time.Millisecond)
}

func TestSemaphore(t *testing.T) {
	s := New(1)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	if err := s.Acquire(ctx); err != nil {
		t.Fatal(err)
	}
	work()
	if err := s.Release(ctx); err != nil {
		t.Fatal(err)
	}
	cancel()
}

func BenchmarkSemaphore(b *testing.B) {
	s := New(1)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := s.Acquire(ctx); err != nil {
			b.Fatal(err)
		}
		if err := s.Release(ctx); err != nil {
			b.Fatal(err)
		}
	}
	cancel()
}
