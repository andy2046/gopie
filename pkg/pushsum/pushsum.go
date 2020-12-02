// Package pushsum implements Push-Sum Protocol.
package pushsum

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"errors"
	"math"
	"math/rand"
	"net"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/andy2046/tik"
)

// ErrNotFound for not found error.
var ErrNotFound = errors.New("Not Found")

type (
	// Gossiper communicate message to other nodes.
	Gossiper interface {
		Gossip(addr net.Addr, msg Message)
	}

	// ValueReader returns the true value.
	ValueReader interface {
		Read(key uint64) float64
	}

	// Message for value / weight information of specific key.
	Message struct {
		Key    uint64 // Key should not be zero
		Value  float64
		Weight float64
	}

	// Config for PushSum.
	Config struct {
		Ticker               *tik.Ticker
		ValueReader          ValueReader
		Gossiper             Gossiper
		IntervalInMS         int
		UpdateSteps          int
		StoreLen             int
		ConvergenceCount     int
		ConvergenceThreshold float64
	}

	// PushSum for gossip based computation of aggregate.
	PushSum struct {
		self        net.Addr
		peers       []net.Addr
		tk          *tik.Ticker
		valueReader ValueReader
		gossiper    Gossiper
		active      uintptr
		closed      uintptr

		store *store
		len   uint64 // power of 2
		pos   uint64

		msgCh  chan Message
		closer chan struct{}

		// interval to send messages to other random node
		intervalInMS uint64

		// the number of steps between updating node value
		// if set to 0 the node value will never be updated
		// it is power of 2
		updateSteps uint64

		convergenceThreshold float64
		convergenceCount     int

		mut  sync.RWMutex
		once sync.Once
	}

	valueWeight struct {
		// current step, step increased 1 for each interval
		step         uint64
		key          uint64
		value        float64
		weight       float64
		trueValue    float64
		inTransition bool
		previous     float64
		count        int
	}
)

func init() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("fail to seed math/rand pkg with crypto/rand random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

// New returns a new PushSum instance.
func New(self net.Addr, peers []net.Addr, cfg Config) *PushSum {
	if self == nil {
		panic("nil self")
	}

	if len(peers) == 0 {
		panic("empty peers")
	}

	if cfg.ValueReader == nil {
		panic("nil ValueReader")
	}

	if cfg.Gossiper == nil {
		panic("nil Gossiper")
	}

	ct := cfg.ConvergenceThreshold
	if ct == 0 {
		ct = math.Pow10(-10)
	}

	cc := cfg.ConvergenceCount
	if cc < 3 {
		cc = 3
	}

	steps := nextPowerOf2(uint32(cfg.UpdateSteps))

	slen := cfg.StoreLen
	if slen < 8 {
		slen = 8
	}
	size := nextPowerOf2(uint32(slen))

	tk := cfg.Ticker
	if tk == nil {
		tk = tik.New()
	}

	ps := &PushSum{
		convergenceThreshold: ct,
		convergenceCount:     cc,
		intervalInMS:         uint64(cfg.IntervalInMS),
		msgCh:                make(chan Message),
		closer:               make(chan struct{}),
		store:                newStore(size),
		len:                  uint64(size),
		updateSteps:          uint64(steps),
		self:                 self,
		peers:                peers,
		tk:                   tk,
		valueReader:          cfg.ValueReader,
		gossiper:             cfg.Gossiper,
	}

	go func() {
		for {
			select {
			case <-ps.closer:
				return
			case m := <-ps.msgCh:
				ps.onPushSumMessage(m)
			}
		}
	}()

	return ps
}

// IsActive returns true if PushSum is active, false otherwise.
func (ps *PushSum) IsActive() bool {
	return atomic.LoadUintptr(&ps.active) == 0
}

// Resume activate the PushSum.
func (ps *PushSum) Resume() {
	atomic.StoreUintptr(&ps.active, 0)
}

// TryPause try to pause the PushSum,
// it returns true if succeed, false otherwise.
func (ps *PushSum) TryPause() bool {
	return atomic.CompareAndSwapUintptr(&ps.active, 0, 1)
}

// Pause wait in a loop to pause the PushSum.
func (ps *PushSum) Pause() {
	for !atomic.CompareAndSwapUintptr(&ps.active, 0, 1) {
		runtime.Gosched()
	}
}

// Close stop PushSum.
func (ps *PushSum) Close() {
	ps.once.Do(func() {
		ps.tk.Close()
		close(ps.closer)
		atomic.CompareAndSwapUintptr(&ps.closed, 0, 1)
	})
}

// IsClosed returns true if closed, false otherwise.
func (ps *PushSum) IsClosed() bool {
	return atomic.LoadUintptr(&ps.closed) == 1
}

// NPeers returns the number of peers.
func (ps *PushSum) NPeers() int {
	ps.mut.RLock()
	l := len(ps.peers)
	ps.mut.RUnlock()
	return l
}

// SetPeers reset peers.
func (ps *PushSum) SetPeers(peers []net.Addr) {
	l := len(peers)
	if l == 0 {
		return
	}

	self := ps.self.String()
	np := make([]net.Addr, 0, l)

	for _, p := range peers {
		if p.String() != self {
			np = append(np, p)
		}
	}

	ps.mut.Lock()
	ps.peers = np
	ps.mut.Unlock()
}

// Estimate returns the estimated average value of all nodes.
func (ps *PushSum) Estimate(key uint64) (float64, error) {
	_, vw, ok := ps.store.get(key)
	if !ok {
		return 0, ErrNotFound
	}

	if vw.inTransition {
		return vw.trueValue, nil
	}
	return vw.value / vw.weight, nil
}

func (ps *PushSum) randomNode() net.Addr {
	ps.mut.RLock()
	defer ps.mut.RUnlock()

	l := len(ps.peers)
	return ps.peers[rand.Intn(l)]
}

// OnMessage process message from other nodes.
func (ps *PushSum) OnMessage(msg Message) {
	ps.msgCh <- msg
}

// onNextStep is called to send message periodically.
func (ps *PushSum) onNextStep(k uint64) bool {
	pos, vw, ok := ps.store.get(k)
	if !ok {
		// stop callback timer for this key
		return false
	}

	vw.step++
	v, w := vw.value, vw.weight

	// update value
	if ps.updateSteps > 0 && (vw.step&(ps.updateSteps-1)) == 0 {
		newV := ps.valueReader.Read(k)
		v += (newV - vw.trueValue)
		vw.trueValue = newV
		vw.inTransition = (newV != vw.trueValue)
	}

	// send to self
	vw.value = v / 2
	vw.weight = w / 2

	// if ratio v/w did not change more than threshold in
	// 3 consecutive rounds then it's converged
	estimate := vw.value / vw.weight
	if math.Abs(vw.previous-estimate) < ps.convergenceThreshold {
		vw.count++
	}
	vw.previous = estimate

	if ok := ps.store.compareAndSet(pos, vw); !ok {
		// k deleted from index?
		return false
	}

	// send to random node
	ps.gossiper.Gossip(ps.randomNode(), Message{
		Key:    k,
		Value:  v / 2,
		Weight: w / 2,
	})

	if vw.count == ps.convergenceCount {
		// converged
		return false
	}

	return true
}

// onPushSumMessage is called to receive message.
func (ps *PushSum) onPushSumMessage(msg Message) {
	if !ps.IsActive() {
		// if stopped participating,
		// forward message to another random node
		ps.gossiper.Gossip(ps.randomNode(), msg)
		return
	}

	k := msg.Key

	// update value if key exist in ring
	if ok := ps.store.update(k, msg.Value, msg.Weight, false); ok {
		return
	}

	// key not exist
	pos := atomic.LoadUint64(&ps.pos) & (ps.len - 1)
	v := ps.valueReader.Read(k)
	vw := valueWeight{
		key:       k,
		value:     v,
		weight:    1.0,
		trueValue: v,
	}

	vw.value += msg.Value
	vw.weight += msg.Weight

	ps.store.insert(pos, vw)
	atomic.AddUint64(&ps.pos, 1) // TODO: truncate pos

	_ = ps.tk.Schedule(ps.intervalInMS, ps.callback(k))
}

func (ps *PushSum) callback(k uint64) tik.Callback {
	return func() {
		if ok := ps.onNextStep(k); ok {
			_ = ps.tk.Schedule(ps.intervalInMS, ps.callback(k))
		}
	}
}

func nextPowerOf2(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	return v + 1
}
