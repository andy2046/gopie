package ratelimit

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type (
	// SlidingWindowLimiter implements a limiter with sliding window counter.
	SlidingWindowLimiter struct {
		limit  Limit
		expire int
		store  Store
	}

	// Store represents a store for limiter state.
	Store interface {
		// Incr add `increment` to field `timestamp` in `key`
		Incr(key string, timestamp int64, increment int) error
		// SetIncr set `key` and add `increment` to field `timestamp` in `key`
		SetIncr(key string, timestamp int64, increment int) error
		// Expire set `key` to expire in `expire` seconds
		Expire(key string, expire int) error
		// Get returns value of field `timestamp` in `key`
		Get(key string, timestamp int64) int
		// Exists check if `key` exists
		Exists(key string) bool
	}

	redisStore struct {
		clientOpts *redis.Options
		client     *redis.Client
		mu         sync.RWMutex
	}
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

// NotFound will be returned if it fails to get value.
const NotFound = -1

// NewSlidingWindowLimiter returns a new Limiter at rate `r` tokens per second,
// and the key expires in `expire` seconds.
func NewSlidingWindowLimiter(r Limit, expire int, store Store) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		limit:  r,
		store:  store,
		expire: expire,
	}
}

// Allow is the shortcut for AllowN(time.Now(), key, 1).
func (s *SlidingWindowLimiter) Allow(key string) bool {
	return s.AllowN(time.Now(), key, 1)
}

// AllowN checks whether `n` requests for `key` may happen at time `now`.
func (s *SlidingWindowLimiter) AllowN(now time.Time, key string, n int) bool {
	sec := timeToSeconds(now)
	var err error

	if existed := s.store.Exists(key); existed {
		err = s.store.Incr(key, sec, n)
	} else {
		err = s.store.SetIncr(key, sec, n)
		if err == nil {
			s.store.Expire(key, s.expire)
		}
	}

	if err != nil {
		logger.Println(err)
		return false
	}

	if count := s.store.Get(key, sec); count == NotFound || float64(count) > float64(s.limit) {
		return false
	}
	return true
}

// NewRedisStore returns a new Redis Store.
func NewRedisStore(clientOptions *redis.Options) (Store, error) {
	r := &redisStore{
		clientOpts: clientOptions,
	}
	err := r.newConnection()
	return r, err
}

func (r *redisStore) newConnection() error {
	r.mu.RLock()
	if r.client != nil {
		r.mu.RUnlock()
		return nil
	}
	r.mu.RUnlock()

	client := redis.NewClient(r.clientOpts)
	_, err := client.Ping().Result()
	if err != nil {
		return err
	}

	r.mu.Lock()
	r.client = client
	r.mu.Unlock()
	return nil
}

func (r *redisStore) Incr(key string, field int64, increment int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	hIncrBy := r.client.HIncrBy(key, strconv.FormatInt(field, 10), int64(increment))
	if err := hIncrBy.Err(); err != nil {
		val := hIncrBy.Val()
		return fmt.Errorf("ratelimit: Incr val=%v error=%v", val, err)
	}
	return nil
}
func (r *redisStore) SetIncr(key string, field int64, increment int) error {
	return r.Incr(key, field, increment)
}

func (r *redisStore) Expire(key string, timeout int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	expire := r.client.Expire(key, time.Duration(timeout)*time.Second)
	if err := expire.Err(); err != nil {
		val, ttl := expire.Val(), r.client.TTL(key)
		return fmt.Errorf("ratelimit: Expire val=%v ttl=%v error=%v", val, ttl.Val(), err)
	}
	return nil
}

func (r *redisStore) Get(key string, field int64) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	hGet := r.client.HGet(key, strconv.FormatInt(field, 10))
	val := hGet.Val()
	if err := hGet.Err(); err != nil || err == redis.Nil || val == "" {
		logger.Printf("ratelimit: Get val=%v error=%v\n", val, err)
		return NotFound
	}
	n, _ := strconv.Atoi(val)
	return n
}

func (r *redisStore) Exists(key string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	n, err := r.client.Exists(key).Result()
	if err != nil || n == 0 {
		logger.Printf("ratelimit: Exists record count=%v error=%v\n", n, err)
		return false
	}
	return true
}

func timeToSeconds(t time.Time) int64 {
	newT := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
	return newT.Unix()
}
