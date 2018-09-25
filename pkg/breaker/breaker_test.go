package breaker

import (
	"fmt"
	"testing"
	"time"
)

var (
	defaultCB *CircuitBreaker
	customCB  *CircuitBreaker
)

func init() {
	defaultCB = New()
	customCB = newCustom()
}

func sleep(cb *CircuitBreaker, period time.Duration) {
	if !cb.expiry.IsZero() {
		cb.expiry = cb.expiry.Add(-period)
	}
}

func succeed(cb *CircuitBreaker) error {
	_, err := cb.Execute(func() (interface{}, error) { return nil, nil })
	return err
}

func fail(cb *CircuitBreaker) error {
	msg := "fail"
	_, err := cb.Execute(func() (interface{}, error) { return nil, fmt.Errorf(msg) })
	if err.Error() == msg {
		return nil
	}
	return err
}

func succeedLater(cb *CircuitBreaker, delay time.Duration) <-chan error {
	ch := make(chan error)
	go func() {
		_, err := cb.Execute(func() (interface{}, error) {
			time.Sleep(delay)
			return nil, nil
		})
		ch <- err
	}()
	return ch
}

func newCustom() *CircuitBreaker {
	opt := func(st *Settings) error {
		st.Name = "breaker"
		st.MaxRequests = 3
		st.Interval = 30 * time.Second
		st.Timeout = 90 * time.Second
		st.ShouldTrip = func(counts Counts) bool {
			return counts.Requests == 3 && counts.TotalFailures == 2
		}
		return nil
	}

	return New(opt)
}

func assert(t *testing.T) func(bool) {
	return func(equal bool) {
		if !equal {
			t.Error("fail to assert")
		}
	}
}

func TestDefaultCircuitBreaker(t *testing.T) {
	eq := assert(t)
	eq("CircuitBreaker" == defaultCB.Name())

	for i := 0; i < 5; i++ {
		eq(nil == fail(defaultCB))
	}
	eq(StateClosed == defaultCB.State())
	eq("{Requests:5 TotalSuccesses:0 TotalFailures:5 ConsecutiveSuccesses:0 ConsecutiveFailures:5}" ==
		fmt.Sprintf("%+v", defaultCB.counts))
	eq(nil == succeed(defaultCB))
	eq(StateClosed == defaultCB.State())
	eq("{Requests:6 TotalSuccesses:1 TotalFailures:5 ConsecutiveSuccesses:1 ConsecutiveFailures:0}" ==
		fmt.Sprintf("%+v", defaultCB.counts))

	eq(nil == fail(defaultCB))
	eq(StateClosed == defaultCB.State())
	eq("{Requests:7 TotalSuccesses:1 TotalFailures:6 ConsecutiveSuccesses:0 ConsecutiveFailures:1}" ==
		fmt.Sprintf("%+v", defaultCB.counts))

	// StateClosed to StateOpen
	for i := 0; i < 5; i++ {
		eq(nil == fail(defaultCB)) // 6 consecutive failures
	}
	eq(StateOpen == defaultCB.State())
	eq("{Requests:0 TotalSuccesses:0 TotalFailures:0 ConsecutiveSuccesses:0 ConsecutiveFailures:0}" ==
		fmt.Sprintf("%+v", defaultCB.counts))
	eq(false == defaultCB.expiry.IsZero())

	eq(nil != succeed(defaultCB))
	eq(nil != fail(defaultCB))
	eq("{Requests:0 TotalSuccesses:0 TotalFailures:0 ConsecutiveSuccesses:0 ConsecutiveFailures:0}" ==
		fmt.Sprintf("%+v", defaultCB.counts))

	sleep(defaultCB, 59*time.Second)
	eq(StateOpen == defaultCB.State())

	// StateOpen to StateHalfOpen
	sleep(defaultCB, 1*time.Second) // over Timeout 60 seconds
	eq(StateHalfOpen == defaultCB.State())
	eq(true == defaultCB.expiry.IsZero())

	// StateHalfOpen to StateOpen
	eq(nil == fail(defaultCB))
	eq(StateOpen == defaultCB.State())
	eq("{Requests:0 TotalSuccesses:0 TotalFailures:0 ConsecutiveSuccesses:0 ConsecutiveFailures:0}" ==
		fmt.Sprintf("%+v", defaultCB.counts))
	eq(false == defaultCB.expiry.IsZero())

	// StateOpen to StateHalfOpen
	sleep(defaultCB, 60*time.Second)
	eq(StateHalfOpen == defaultCB.State())
	eq(true == defaultCB.expiry.IsZero())

	// StateHalfOpen to StateClosed
	eq(nil == succeed(defaultCB))
	eq(StateClosed == defaultCB.State())
	eq("{Requests:0 TotalSuccesses:0 TotalFailures:0 ConsecutiveSuccesses:0 ConsecutiveFailures:0}" ==
		fmt.Sprintf("%+v", defaultCB.counts))
	eq(true == defaultCB.expiry.IsZero())
}

func TestCustomCircuitBreaker(t *testing.T) {
	eq := assert(t)
	eq("breaker" == customCB.Name())

	for i := 0; i < 5; i++ {
		eq(nil == succeed(customCB))
		eq(nil == fail(customCB))
	}
	eq(StateClosed == customCB.State())
	eq("{Requests:10 TotalSuccesses:5 TotalFailures:5 ConsecutiveSuccesses:0 ConsecutiveFailures:1}" ==
		fmt.Sprintf("%+v", customCB.counts))

	sleep(customCB, 29*time.Second)
	eq(nil == succeed(customCB))
	eq(StateClosed == customCB.State())
	eq("{Requests:11 TotalSuccesses:6 TotalFailures:5 ConsecutiveSuccesses:1 ConsecutiveFailures:0}" ==
		fmt.Sprintf("%+v", customCB.counts))

	sleep(customCB, 1*time.Second) // over Interval
	eq(nil == fail(customCB))
	eq(StateClosed == customCB.State())
	eq("{Requests:1 TotalSuccesses:0 TotalFailures:1 ConsecutiveSuccesses:0 ConsecutiveFailures:1}" ==
		fmt.Sprintf("%+v", customCB.counts))

	// StateClosed to StateOpen
	eq(nil == succeed(customCB))
	eq(nil == fail(customCB)) // ShouldTrip triggered
	eq(StateOpen == customCB.State())
	eq("{Requests:0 TotalSuccesses:0 TotalFailures:0 ConsecutiveSuccesses:0 ConsecutiveFailures:0}" ==
		fmt.Sprintf("%+v", customCB.counts))
	eq(false == customCB.expiry.IsZero())

	// StateOpen to StateHalfOpen
	sleep(customCB, 90*time.Second)
	eq(StateHalfOpen == customCB.State())
	eq(true == customCB.expiry.IsZero())

	eq(nil == succeed(customCB))
	eq(nil == succeed(customCB))
	eq(StateHalfOpen == customCB.State())
	eq("{Requests:2 TotalSuccesses:2 TotalFailures:0 ConsecutiveSuccesses:2 ConsecutiveFailures:0}" ==
		fmt.Sprintf("%+v", customCB.counts))

	// StateHalfOpen to StateClosed
	ch := succeedLater(customCB, 100*time.Millisecond) // 3 consecutive successes
	time.Sleep(50 * time.Millisecond)
	customCB.mutex.Lock()
	eq("{Requests:3 TotalSuccesses:2 TotalFailures:0 ConsecutiveSuccesses:2 ConsecutiveFailures:0}" ==
		fmt.Sprintf("%+v", customCB.counts))
	customCB.mutex.Unlock()
	eq(nil != succeed(customCB)) // over MaxRequests
	eq(nil == <-ch)
	eq(StateClosed == customCB.State())
	eq("{Requests:0 TotalSuccesses:0 TotalFailures:0 ConsecutiveSuccesses:0 ConsecutiveFailures:0}" ==
		fmt.Sprintf("%+v", customCB.counts))
	eq(false == customCB.expiry.IsZero())
}
