package ratelimit

import (
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func TestSlidingWindowLimiter(t *testing.T) {
	t.Skip("only works with a local running Redis server")
	opts := &redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}
	key := "user_1"
	expire := 10
	var once sync.Once
	redStore, _ := NewRedisStore(opts)
	sLimiter := NewSlidingWindowLimiter(1, expire, redStore)

	t1 := t0.Add(200 * time.Millisecond)
	t2 := t0.Add(time.Duration(expire) * time.Second)
	cases := []struct {
		t  time.Time
		n  int
		ok bool
	}{
		{t0, 1, true},
		{t0, 1, false},
		{t0, 1, false},
		{t1, 1, false},
		{t1, 1, false},
		{t2, 1, true},
		{t2, 1, false},
	}

	for _, c := range cases {
		if c.t == t2 {
			once.Do(func() {
				time.Sleep(time.Duration(expire) * time.Second)
			})
		}
		ok := sLimiter.AllowN(c.t, key, c.n)
		if ok != c.ok {
			t.Errorf("AllowN(%v, %v, %v) = %v want %v",
				c.t, key, c.n, ok, c.ok)
		}
	}
}
