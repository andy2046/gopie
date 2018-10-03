

# ratelimit
`import "github.com/andy2046/gopie/pkg/ratelimit"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package ratelimit implements a rate limiter.




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [type Limit](#Limit)
  * [func Every(interval time.Duration) Limit](#Every)
* [type Limiter](#Limiter)
  * [func New(r Limit, b int) *Limiter](#New)
  * [func (l *Limiter) Allow() bool](#Limiter.Allow)
  * [func (l *Limiter) AllowN(now time.Time, n int) bool](#Limiter.AllowN)
  * [func (l *Limiter) Burst() int](#Limiter.Burst)
  * [func (l *Limiter) Limit() Limit](#Limiter.Limit)
  * [func (l *Limiter) Wait() (time.Duration, error)](#Limiter.Wait)
  * [func (l *Limiter) WaitN(now time.Time, n int) (time.Duration, error)](#Limiter.WaitN)
* [type SlidingWindowLimiter](#SlidingWindowLimiter)
  * [func NewSlidingWindowLimiter(r Limit, expire int, store Store) *SlidingWindowLimiter](#NewSlidingWindowLimiter)
  * [func (s *SlidingWindowLimiter) Allow(key string) bool](#SlidingWindowLimiter.Allow)
  * [func (s *SlidingWindowLimiter) AllowN(now time.Time, key string, n int) bool](#SlidingWindowLimiter.AllowN)
* [type Store](#Store)
  * [func NewRedisStore(clientOptions *redis.Options) (Store, error)](#NewRedisStore)


#### <a name="pkg-files">Package files</a>
[ratelimit.go](/src/github.com/andy2046/gopie/pkg/ratelimit/ratelimit.go) [sliding_window.go](/src/github.com/andy2046/gopie/pkg/ratelimit/sliding_window.go) 


## <a name="pkg-constants">Constants</a>
``` go
const NotFound = -1
```
NotFound will be returned if it fails to get value.





## <a name="Limit">type</a> [Limit](/src/target/ratelimit.go?s=162:180#L11)
``` go
type Limit float64
```
Limit defines the maximum number of requests per second.







### <a name="Every">func</a> [Every](/src/target/ratelimit.go?s=245:285#L14)
``` go
func Every(interval time.Duration) Limit
```
Every converts a time interval between requests to a Limit.





## <a name="Limiter">type</a> [Limiter](/src/target/ratelimit.go?s=512:742#L22)
``` go
type Limiter struct {
    // contains filtered or unexported fields
}
```
Limiter implements a token bucket limiter at rate `r` tokens per second with burst size of `b` tokens.







### <a name="New">func</a> [New](/src/target/ratelimit.go?s=1013:1046#L44)
``` go
func New(r Limit, b int) *Limiter
```
New returns a new Limiter at rate `r` tokens per second with burst of `b` tokens.





### <a name="Limiter.Allow">func</a> (\*Limiter) [Allow](/src/target/ratelimit.go?s=1149:1179#L52)
``` go
func (l *Limiter) Allow() bool
```
Allow is the shortcut for AllowN(time.Now(), 1).




### <a name="Limiter.AllowN">func</a> (\*Limiter) [AllowN](/src/target/ratelimit.go?s=1281:1332#L57)
``` go
func (l *Limiter) AllowN(now time.Time, n int) bool
```
AllowN checks whether `n` requests may happen at time `now`.




### <a name="Limiter.Burst">func</a> (\*Limiter) [Burst](/src/target/ratelimit.go?s=877:906#L39)
``` go
func (l *Limiter) Burst() int
```
Burst returns the Limiter's burst size.




### <a name="Limiter.Limit">func</a> (\*Limiter) [Limit](/src/target/ratelimit.go?s=781:812#L34)
``` go
func (l *Limiter) Limit() Limit
```
Limit returns the Limiter's rate.




### <a name="Limiter.Wait">func</a> (\*Limiter) [Wait](/src/target/ratelimit.go?s=1421:1468#L62)
``` go
func (l *Limiter) Wait() (time.Duration, error)
```
Wait is the shortcut for WaitN(time.Now(), 1).




### <a name="Limiter.WaitN">func</a> (\*Limiter) [WaitN](/src/target/ratelimit.go?s=1597:1665#L67)
``` go
func (l *Limiter) WaitN(now time.Time, n int) (time.Duration, error)
```
WaitN calculates the time duration to wait before `n` requests may happen at time `now`.




## <a name="SlidingWindowLimiter">type</a> [SlidingWindowLimiter](/src/target/sliding_window.go?s=191:266#L16)
``` go
type SlidingWindowLimiter struct {
    // contains filtered or unexported fields
}
```
SlidingWindowLimiter implements a limiter with sliding window counter.







### <a name="NewSlidingWindowLimiter">func</a> [NewSlidingWindowLimiter](/src/target/sliding_window.go?s=1209:1293#L50)
``` go
func NewSlidingWindowLimiter(r Limit, expire int, store Store) *SlidingWindowLimiter
```
NewSlidingWindowLimiter returns a new Limiter at rate `r` tokens per second,
and the key expires in `expire` seconds.





### <a name="SlidingWindowLimiter.Allow">func</a> (\*SlidingWindowLimiter) [Allow](/src/target/sliding_window.go?s=1438:1491#L59)
``` go
func (s *SlidingWindowLimiter) Allow(key string) bool
```
Allow is the shortcut for AllowN(time.Now(), key, 1).




### <a name="SlidingWindowLimiter.AllowN">func</a> (\*SlidingWindowLimiter) [AllowN](/src/target/sliding_window.go?s=1608:1684#L64)
``` go
func (s *SlidingWindowLimiter) AllowN(now time.Time, key string, n int) bool
```
AllowN checks whether `n` requests for `key` may happen at time `now`.




## <a name="Store">type</a> [Store](/src/target/sliding_window.go?s=317:826#L23)
``` go
type Store interface {
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
```
Store represents a store for limiter state.







### <a name="NewRedisStore">func</a> [NewRedisStore](/src/target/sliding_window.go?s=2145:2208#L89)
``` go
func NewRedisStore(clientOptions *redis.Options) (Store, error)
```
NewRedisStore returns a new Redis Store.









- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
