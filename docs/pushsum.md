

# pushsum
`import "github.com/andy2046/gopie/pkg/pushsum"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package pushsum implements Push-Sum Protocol.




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [type Config](#Config)
* [type Gossiper](#Gossiper)
* [type Message](#Message)
* [type PushSum](#PushSum)
  * [func New(self net.Addr, peers []net.Addr, cfg Config) *PushSum](#New)
  * [func (ps *PushSum) Close()](#PushSum.Close)
  * [func (ps *PushSum) Estimate(key uint64) (float64, error)](#PushSum.Estimate)
  * [func (ps *PushSum) IsActive() bool](#PushSum.IsActive)
  * [func (ps *PushSum) IsClosed() bool](#PushSum.IsClosed)
  * [func (ps *PushSum) NPeers() int](#PushSum.NPeers)
  * [func (ps *PushSum) OnMessage(msg Message)](#PushSum.OnMessage)
  * [func (ps *PushSum) Pause()](#PushSum.Pause)
  * [func (ps *PushSum) Resume()](#PushSum.Resume)
  * [func (ps *PushSum) SetPeers(peers []net.Addr)](#PushSum.SetPeers)
  * [func (ps *PushSum) TryPause() bool](#PushSum.TryPause)
* [type Scheduler](#Scheduler)
* [type ValueReader](#ValueReader)


#### <a name="pkg-files">Package files</a>
[pushsum.go](/src/github.com/andy2046/gopie/pkg/pushsum/pushsum.go) [store.go](/src/github.com/andy2046/gopie/pkg/pushsum/store.go) 



## <a name="pkg-variables">Variables</a>
``` go
var ErrNotFound = errors.New("Not Found")
```
ErrNotFound for not found error.




## <a name="Config">type</a> [Config](/src/target/pushsum.go?s=801:1058#L44)
``` go
type Config struct {
    Scheduler            Scheduler
    ValueReader          ValueReader
    Gossiper             Gossiper
    IntervalInMS         int
    UpdateSteps          int
    StoreLen             int
    ConvergenceCount     int
    ConvergenceThreshold float64
}
```
Config for PushSum.










## <a name="Gossiper">type</a> [Gossiper](/src/target/pushsum.go?s=332:392#L21)
``` go
type Gossiper interface {
    Gossip(addr net.Addr, msg Message)
}
```
Gossiper communicate message to other nodes.










## <a name="Message">type</a> [Message](/src/target/pushsum.go?s=679:774#L37)
``` go
type Message struct {
    Key    uint64 // Key should not be zero
    Value  float64
    Weight float64
}
```
Message for value / weight information of specific key.










## <a name="PushSum">type</a> [PushSum](/src/target/pushsum.go?s=1116:1729#L56)
``` go
type PushSum struct {
    // contains filtered or unexported fields
}
```
PushSum for gossip based computation of aggregate.







### <a name="New">func</a> [New](/src/target/pushsum.go?s=2244:2306#L110)
``` go
func New(self net.Addr, peers []net.Addr, cfg Config) *PushSum
```
New returns a new PushSum instance.





### <a name="PushSum.Close">func</a> (\*PushSum) [Close](/src/target/pushsum.go?s=4101:4127#L203)
``` go
func (ps *PushSum) Close()
```
Close stop PushSum.




### <a name="PushSum.Estimate">func</a> (\*PushSum) [Estimate](/src/target/pushsum.go?s=4894:4950#L246)
``` go
func (ps *PushSum) Estimate(key uint64) (float64, error)
```
Estimate returns the estimated average value of all nodes.




### <a name="PushSum.IsActive">func</a> (\*PushSum) [IsActive](/src/target/pushsum.go?s=3556:3590#L180)
``` go
func (ps *PushSum) IsActive() bool
```
IsActive returns true if PushSum is active, false otherwise.




### <a name="PushSum.IsClosed">func</a> (\*PushSum) [IsClosed](/src/target/pushsum.go?s=4302:4336#L212)
``` go
func (ps *PushSum) IsClosed() bool
```
IsClosed returns true if closed, false otherwise.




### <a name="PushSum.NPeers">func</a> (\*PushSum) [NPeers](/src/target/pushsum.go?s=4425:4456#L217)
``` go
func (ps *PushSum) NPeers() int
```
NPeers returns the number of peers.




### <a name="PushSum.OnMessage">func</a> (\*PushSum) [OnMessage](/src/target/pushsum.go?s=5297:5338#L267)
``` go
func (ps *PushSum) OnMessage(msg Message)
```
OnMessage process message from other nodes.




### <a name="PushSum.Pause">func</a> (\*PushSum) [Pause](/src/target/pushsum.go?s=3968:3994#L196)
``` go
func (ps *PushSum) Pause()
```
Pause wait in a loop to pause the PushSum.




### <a name="PushSum.Resume">func</a> (\*PushSum) [Resume](/src/target/pushsum.go?s=3672:3699#L185)
``` go
func (ps *PushSum) Resume()
```
Resume activate the PushSum.




### <a name="PushSum.SetPeers">func</a> (\*PushSum) [SetPeers](/src/target/pushsum.go?s=4551:4596#L225)
``` go
func (ps *PushSum) SetPeers(peers []net.Addr)
```
SetPeers reset peers.




### <a name="PushSum.TryPause">func</a> (\*PushSum) [TryPause](/src/target/pushsum.go?s=3827:3861#L191)
``` go
func (ps *PushSum) TryPause() bool
```
TryPause try to pause the PushSum,
it returns true if succeed, false otherwise.




## <a name="Scheduler">type</a> [Scheduler](/src/target/pushsum.go?s=447:520#L26)
``` go
type Scheduler interface {
    Schedule(interval uint64, cb func())
    Close()
}
```
Scheduler schedule tasks to run in an interval.










## <a name="ValueReader">type</a> [ValueReader](/src/target/pushsum.go?s=563:616#L32)
``` go
type ValueReader interface {
    Read(key uint64) float64
}
```
ValueReader returns the true value.














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
