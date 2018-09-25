

# breaker
`import "github.com/andy2046/gopie/pkg/breaker"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package breaker implements Circuit Breaker pattern.




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [type CircuitBreaker](#CircuitBreaker)
  * [func New(options ...Option) *CircuitBreaker](#New)
  * [func (cb *CircuitBreaker) Execute(request func() (interface{}, error)) (interface{}, error)](#CircuitBreaker.Execute)
  * [func (cb *CircuitBreaker) Name() string](#CircuitBreaker.Name)
  * [func (cb *CircuitBreaker) State() State](#CircuitBreaker.State)
* [type Counts](#Counts)
* [type Option](#Option)
* [type Settings](#Settings)
* [type State](#State)
  * [func (s State) String() string](#State.String)


#### <a name="pkg-files">Package files</a>
[breaker.go](/src/github.com/andy2046/gopie/pkg/breaker/breaker.go) 



## <a name="pkg-variables">Variables</a>
``` go
var (
    // ErrTooManyRequests is returned when the state is half open
    // and the requests count is more the maxRequests.
    ErrTooManyRequests = errors.New("too many requests, requests count is more the maxRequests in half open state")
    // ErrOpenState is returned when the state is open.
    ErrOpenState = errors.New("circuit breaker is open")
)
```
``` go
var DefaultSettings = Settings{
    Name:        "CircuitBreaker",
    MaxRequests: 1,
    Interval:    0,
    Timeout:     60 * time.Second,
    ShouldTrip: func(counts Counts) bool {
        return counts.ConsecutiveFailures > 5
    },
    OnStateChange: nil,
}
```
DefaultSettings is the default CircuitBreaker Settings.




## <a name="CircuitBreaker">type</a> [CircuitBreaker](/src/target/breaker.go?s=3218:3548#L109)
``` go
type CircuitBreaker struct {
    // contains filtered or unexported fields
}
```
CircuitBreaker prevent an application repeatedly trying to execute an operation that is likely to fail.







### <a name="New">func</a> [New](/src/target/breaker.go?s=4166:4209#L150)
``` go
func New(options ...Option) *CircuitBreaker
```
New returns a new CircuitBreaker with options applied.





### <a name="CircuitBreaker.Execute">func</a> (\*CircuitBreaker) [Execute](/src/target/breaker.go?s=5245:5336#L190)
``` go
func (cb *CircuitBreaker) Execute(request func() (interface{}, error)) (interface{}, error)
```
Execute runs the given request if the CircuitBreaker accepts it.
Execute returns an error instantly if the CircuitBreaker rejects the request.
Otherwise, Execute returns the result of the request.
If a panic occurs in the request, the CircuitBreaker handles it as an error
and causes the same panic again.




### <a name="CircuitBreaker.Name">func</a> (\*CircuitBreaker) [Name](/src/target/breaker.go?s=4650:4689#L171)
``` go
func (cb *CircuitBreaker) Name() string
```
Name returns the name of the CircuitBreaker.




### <a name="CircuitBreaker.State">func</a> (\*CircuitBreaker) [State](/src/target/breaker.go?s=4769:4808#L176)
``` go
func (cb *CircuitBreaker) State() State
```
State returns the current state of the CircuitBreaker.




## <a name="Counts">type</a> [Counts](/src/target/breaker.go?s=1244:1411#L50)
``` go
type Counts struct {
    Requests             uint64
    TotalSuccesses       uint64
    TotalFailures        uint64
    ConsecutiveSuccesses uint64
    ConsecutiveFailures  uint64
}
```
Counts holds the numbers of requests and their successes/failures.
CircuitBreaker clears the internal Counts either
on the change of the state or at the closed-state intervals.
Counts ignores the results of the requests sent before clearing.










## <a name="Option">type</a> [Option](/src/target/breaker.go?s=3607:3637#L125)
``` go
type Option = func(*Settings) error
```
Option applies settings to CircuitBreaker Settings.










## <a name="Settings">type</a> [Settings](/src/target/breaker.go?s=1879:3107#L84)
``` go
type Settings struct {
    // Name is the name of the CircuitBreaker.
    Name string
    // MaxRequests is the maximum number of requests allowed to pass through
    // when the CircuitBreaker is half-open.
    // If MaxRequests is 0, the CircuitBreaker allows only 1 request.
    MaxRequests uint64
    // Interval is the cyclic period of the closed state
    // for the CircuitBreaker to clear the internal Counts.
    // If Interval is 0, the CircuitBreaker doesn't clear internal Counts during the closed state.
    Interval time.Duration
    // Timeout is the period of the open state,
    // after which the state of the CircuitBreaker becomes half-open.
    // If Timeout is 0, the timeout for the CircuitBreaker is 60 seconds.
    Timeout time.Duration
    // ShouldTrip is called with a copy of Counts whenever a request fails in the closed state.
    // If ShouldTrip returns true, the CircuitBreaker will be placed into the open state.
    // If ShouldTrip is nil, default ShouldTrip is used.
    // Default ShouldTrip returns true when the number of consecutive failures is more than 5.
    ShouldTrip func(counts Counts) bool
    // OnStateChange is called whenever the state of the CircuitBreaker changes.
    OnStateChange func(name string, from, to State)
}
```
Settings represents settings for CircuitBreaker.










## <a name="State">type</a> [State](/src/target/breaker.go?s=185:199#L13)
``` go
type State int
```
State is the type representing a state of CircuitBreaker.


``` go
const (
    // StateClosed represents Closed State.
    StateClosed State = iota
    // StateHalfOpen represents HalfOpen State.
    StateHalfOpen
    // StateOpen represents Open State.
    StateOpen
)
```









### <a name="State.String">func</a> (State) [String](/src/target/breaker.go?s=772:802#L33)
``` go
func (s State) String() string
```
String implements stringer interface.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
