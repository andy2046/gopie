package pushsum

import (
	"math"
	"net"
	"sync"
	"testing"
	"time"
)

type (
	netAddr struct {
		addr string
	}

	gossiper struct{}

	reader struct {
		receiver chan float64
		sender   chan<- float64
	}
)

var (
	g = gossiper{}

	node0 = netAddr{"0"}
	node1 = netAddr{"1"}
	node2 = netAddr{"2"}

	r0 *reader
	r1 *reader
	r2 *reader

	ps0 *PushSum
	ps1 *PushSum
	ps2 *PushSum
)

func TestPushSum(t *testing.T) {
	threshold := float64(0.2)
	key := uint64(2020)
	wg := &sync.WaitGroup{}
	wg.Add(3)

	r0 = newReader()
	r1 = newReader()
	r2 = newReader()

	ps0 = New(node0, []net.Addr{node1, node2}, nil, r0, g, 100, 1, 64)
	ps1 = New(node1, []net.Addr{node0, node2}, nil, r1, g, 100, 1, 64)
	ps2 = New(node2, []net.Addr{node1, node0}, nil, r2, g, 100, 1, 64)

	defer func() {
		ps0.Close()
		ps1.Close()
		ps2.Close()

		r0.close()
		r1.close()
		r2.close()
	}()

	msg := Message{Key: key, Value: 0, Weight: 0}
	ps0.OnMessage(msg)
	ps1.OnMessage(msg)
	ps2.OnMessage(msg)

	// init with 15/25/35 then add 30 to each
	// average should be 55
	update(r0, key, 15, wg)
	update(r1, key, 25, wg)
	update(r2, key, 35, wg)

	wg.Wait()
	time.Sleep(1000 * time.Millisecond)

	v0, err := ps0.Estimate(key)
	if err != nil {
		t.Error(err)
	}

	v1, err := ps1.Estimate(key)
	if err != nil {
		t.Error(err)
	}

	v2, err := ps2.Estimate(key)
	if err != nil {
		t.Error(err)
	}

	if math.Abs(v0-v1) > threshold ||
		math.Abs(v2-v1) > threshold ||
		math.Abs(v0-v2) > threshold {
		t.Fail()
	}

	t.Logf("v0: %f v1: %f v2: %f", v0, v1, v2)
}

func update(r *reader, k uint64, v float64, wg *sync.WaitGroup) {
	go func() {
		i := float64(0)
		for i < 3 {
			r.update(k, v+15*i)
			time.Sleep(100 * time.Millisecond)
			i++
		}
		wg.Done()
	}()
}

func (na netAddr) Network() string {
	return "pushsum"
}

func (na netAddr) String() string {
	return na.addr
}

func (g gossiper) Gossip(addr net.Addr, msg Message) {
	switch addr.String() {
	case "0":
		ps0.OnMessage(msg)
	case "1":
		ps1.OnMessage(msg)
	case "2":
		ps2.OnMessage(msg)
	default:
	}
}

func newReader() *reader {
	rcv := make(chan float64)
	snd := latest(rcv)
	snd <- 0
	return &reader{
		receiver: rcv,
		sender:   snd,
	}
}

func (r *reader) Read(key uint64) float64 {
	return <-r.receiver
}

func (r *reader) update(key uint64, v float64) {
	r.sender <- v
}

func (r *reader) close() {
	close(r.sender)
}

func latest(receiver chan<- float64) chan<- float64 {
	sender := make(chan float64)

	go func() {
		var (
			latest float64
			ok     bool
			temp   chan<- float64
		)
		for {
			select {
			case latest, ok = <-sender:
				if !ok {
					return
				}
				if temp == nil {
					temp = receiver
				}
				continue
			case temp <- latest:
				break
			}
		}
	}()

	return sender
}
