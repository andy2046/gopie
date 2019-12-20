

# sequence
`import "github.com/andy2046/gopie/pkg/sequence"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package sequence implements Iceflake sequence generator interface.
```bash
Iceflake is the interface for snowflake similar sequence generator.

Iceflake algorithm:

+-------+--------------------+----------+
| sign  | delta milliseconds | sequence |
+-------+--------------------+----------+
| 1 bit | 63-n bits          | n bits   |

sequence (n bits)
The last custom n bits, represents sequence within the one millisecond.

delta milliseconds (63-n bits)
The next 63-n bits, represents delta milliseconds since a custom epoch.
```




## <a name="pkg-index">Index</a>
* [type Iceflake](#Iceflake)
* [type MemFlake](#MemFlake)
  * [func NewMemFlake(startTime time.Time, bitLenSequence uint8, machineID uint64) *MemFlake](#NewMemFlake)
  * [func (m *MemFlake) BitLenSequence() uint8](#MemFlake.BitLenSequence)
  * [func (m *MemFlake) MachineID() uint64](#MemFlake.MachineID)
  * [func (m *MemFlake) Next() (uint64, error)](#MemFlake.Next)
  * [func (m *MemFlake) NextN(n int) (uint64, error)](#MemFlake.NextN)
  * [func (m *MemFlake) StartTime() time.Time](#MemFlake.StartTime)
* [type MemSeq](#MemSeq)
  * [func NewMemSeq(machineID uint64) *MemSeq](#NewMemSeq)
  * [func (m *MemSeq) MachineID() uint64](#MemSeq.MachineID)
  * [func (m *MemSeq) Next() (uint64, error)](#MemSeq.Next)
  * [func (m *MemSeq) NextN(n int) (uint64, error)](#MemSeq.NextN)
* [type Sequencer](#Sequencer)


#### <a name="pkg-files">Package files</a>
[iceflake.go](/src/github.com/andy2046/gopie/pkg/sequence/iceflake.go) [memflake.go](/src/github.com/andy2046/gopie/pkg/sequence/memflake.go) [memsequence.go](/src/github.com/andy2046/gopie/pkg/sequence/memsequence.go) [sequencer.go](/src/github.com/andy2046/gopie/pkg/sequence/sequencer.go) 






## <a name="Iceflake">type</a> [Iceflake](/src/target/iceflake.go?s=109:410#L8)
``` go
type Iceflake interface {
    Sequencer
    // StartTime defines the time since which
    // the Iceflake time is defined as the elapsed time.
    StartTime() time.Time
    // BitLenSequence defines the bit length of sequence number,
    // and the bit length of time is 63 - BitLenSequence().
    BitLenSequence() uint8
}
```
Iceflake is the interface for snowflake similar sequence generator.










## <a name="MemFlake">type</a> [MemFlake](/src/target/memflake.go?s=112:264#L10)
``` go
type MemFlake struct {
    sync.Mutex
    // contains filtered or unexported fields
}
```
MemFlake is an implementation of in-memory Iceflake.







### <a name="NewMemFlake">func</a> [NewMemFlake](/src/target/memflake.go?s=482:569#L26)
``` go
func NewMemFlake(startTime time.Time, bitLenSequence uint8, machineID uint64) *MemFlake
```
NewMemFlake creates a MemFlake.





### <a name="MemFlake.BitLenSequence">func</a> (\*MemFlake) [BitLenSequence](/src/target/memflake.go?s=1897:1938#L93)
``` go
func (m *MemFlake) BitLenSequence() uint8
```
BitLenSequence ...




### <a name="MemFlake.MachineID">func</a> (\*MemFlake) [MachineID](/src/target/memflake.go?s=1729:1766#L83)
``` go
func (m *MemFlake) MachineID() uint64
```
MachineID ...




### <a name="MemFlake.Next">func</a> (\*MemFlake) [Next](/src/target/memflake.go?s=1130:1171#L54)
``` go
func (m *MemFlake) Next() (uint64, error)
```
Next ...




### <a name="MemFlake.NextN">func</a> (\*MemFlake) [NextN](/src/target/memflake.go?s=1209:1256#L59)
``` go
func (m *MemFlake) NextN(n int) (uint64, error)
```
NextN ...




### <a name="MemFlake.StartTime">func</a> (\*MemFlake) [StartTime](/src/target/memflake.go?s=1809:1849#L88)
``` go
func (m *MemFlake) StartTime() time.Time
```
StartTime ...




## <a name="MemSeq">type</a> [MemSeq](/src/target/memsequence.go?s=103:169#L9)
``` go
type MemSeq struct {
    sync.Mutex
    // contains filtered or unexported fields
}
```
MemSeq is an implementation of in-memory Sequencer.







### <a name="NewMemSeq">func</a> [NewMemSeq](/src/target/memsequence.go?s=321:361#L21)
``` go
func NewMemSeq(machineID uint64) *MemSeq
```
NewMemSeq creates a MemSeq.





### <a name="MemSeq.MachineID">func</a> (\*MemSeq) [MachineID](/src/target/memsequence.go?s=685:720#L45)
``` go
func (m *MemSeq) MachineID() uint64
```
MachineID ...




### <a name="MemSeq.Next">func</a> (\*MemSeq) [Next](/src/target/memsequence.go?s=421:460#L28)
``` go
func (m *MemSeq) Next() (uint64, error)
```
Next ...




### <a name="MemSeq.NextN">func</a> (\*MemSeq) [NextN](/src/target/memsequence.go?s=498:543#L33)
``` go
func (m *MemSeq) NextN(n int) (uint64, error)
```
NextN ...




## <a name="Sequencer">type</a> [Sequencer](/src/target/sequencer.go?s=651:946#L21)
``` go
type Sequencer interface {
    // Next returns the next sequence.
    Next() (uint64, error)
    // NextN reserves the next `n` sequences and returns the first one,
    // `n` should not be less than 1.
    NextN(n int) (uint64, error)
    // MachineID returns the unique ID of the instance.
    MachineID() uint64
}
```
Sequencer is the interface for sequence generator.














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
