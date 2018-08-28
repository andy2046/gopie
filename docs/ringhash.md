

# ringhash
`import "github.com/andy2046/gopie/pkg/ringhash"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package ringhash provides a ring hash implementation.




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [type Config](#Config)
* [type Hash](#Hash)
* [type Node](#Node)
* [type Option](#Option)
* [type Ring](#Ring)
  * [func New(options ...Option) *Ring](#New)
  * [func (r *Ring) Add(node string) bool](#Ring.Add)
  * [func (r *Ring) AddNode(keys ...string)](#Ring.AddNode)
  * [func (r *Ring) Done(node string) bool](#Ring.Done)
  * [func (r *Ring) GetLeastNode(key string) (string, error)](#Ring.GetLeastNode)
  * [func (r *Ring) GetNode(key string) (string, error)](#Ring.GetNode)
  * [func (r *Ring) IsEmpty() bool](#Ring.IsEmpty)
  * [func (r *Ring) Loads() map[string]int64](#Ring.Loads)
  * [func (r *Ring) MaxLoad() int64](#Ring.MaxLoad)
  * [func (r *Ring) Nodes() (nodes []string)](#Ring.Nodes)
  * [func (r *Ring) RemoveNode(node string) bool](#Ring.RemoveNode)
  * [func (r *Ring) UpdateLoad(node string, load int64)](#Ring.UpdateLoad)


#### <a name="pkg-files">Package files</a>
[ringhash.go](/src/github.com/andy2046/gopie/pkg/ringhash/ringhash.go) 



## <a name="pkg-variables">Variables</a>
``` go
var (

    // ErrNoNode when there is no node added into the hash ring.
    ErrNoNode = errors.New("no node added")
    // ErrNodeNotFound when no node found in LoadMap.
    ErrNodeNotFound = errors.New("node not found in LoadMap")
    // DefaultConfig is the default config for hash ring.
    DefaultConfig = Config{
        HashFn:          hash,
        Replicas:        10,
        BalancingFactor: 1.25,
    }
)
```



## <a name="Config">type</a> [Config](/src/target/ringhash.go?s=619:708#L38)
``` go
type Config struct {
    HashFn          Hash
    Replicas        int
    BalancingFactor float64
}
```
Config is the config for hash ring.










## <a name="Hash">type</a> [Hash](/src/target/ringhash.go?s=191:219#L16)
``` go
type Hash func(key string) uint64
```
Hash is the hash function.










## <a name="Node">type</a> [Node](/src/target/ringhash.go?s=256:299#L19)
``` go
type Node struct {
    Name string
    Load int64
}
```
Node is the node in the ring.










## <a name="Option">type</a> [Option](/src/target/ringhash.go?s=748:776#L45)
``` go
type Option = func(*Config) error
```
Option applies config to Config.










## <a name="Ring">type</a> [Ring](/src/target/ringhash.go?s=348:576#L25)
``` go
type Ring struct {
    // contains filtered or unexported fields
}
```
Ring is the data store for keys hash map.







### <a name="New">func</a> [New](/src/target/ringhash.go?s=1480:1513#L73)
``` go
func New(options ...Option) *Ring
```
New returns a new Ring.





### <a name="Ring.Add">func</a> (\*Ring) [Add](/src/target/ringhash.go?s=3684:3720#L179)
``` go
func (r *Ring) Add(node string) bool
```
Add increases load of the given node by 1,
should only be used with GetLeast.




### <a name="Ring.AddNode">func</a> (\*Ring) [AddNode](/src/target/ringhash.go?s=1983:2021#L93)
``` go
func (r *Ring) AddNode(keys ...string)
```
AddNode adds Node with key as name to the hash ring.




### <a name="Ring.Done">func</a> (\*Ring) [Done](/src/target/ringhash.go?s=3958:3995#L193)
``` go
func (r *Ring) Done(node string) bool
```
Done decreases load of the given node by 1,
should only be used with GetLeast.




### <a name="Ring.GetLeastNode">func</a> (\*Ring) [GetLeastNode](/src/target/ringhash.go?s=2862:2917#L135)
``` go
func (r *Ring) GetLeastNode(key string) (string, error)
```
GetLeastNode uses consistent hashing with bounded loads to get the least loaded node.




### <a name="Ring.GetNode">func</a> (\*Ring) [GetNode](/src/target/ringhash.go?s=2554:2604#L121)
``` go
func (r *Ring) GetNode(key string) (string, error)
```
GetNode returns the closest node in the hash ring to the provided key.




### <a name="Ring.IsEmpty">func</a> (\*Ring) [IsEmpty](/src/target/ringhash.go?s=1861:1890#L88)
``` go
func (r *Ring) IsEmpty() bool
```
IsEmpty returns true if there is no node in the ring.




### <a name="Ring.Loads">func</a> (\*Ring) [Loads](/src/target/ringhash.go?s=4772:4811#L235)
``` go
func (r *Ring) Loads() map[string]int64
```
Loads returns the loads of all the nodes in the hash ring.




### <a name="Ring.MaxLoad">func</a> (\*Ring) [MaxLoad](/src/target/ringhash.go?s=5081:5111#L248)
``` go
func (r *Ring) MaxLoad() int64
```
MaxLoad returns the maximum load for a single node in the hash ring,
which is (totalLoad/numberOfNodes)*balancingFactor.




### <a name="Ring.Nodes">func</a> (\*Ring) [Nodes](/src/target/ringhash.go?s=4559:4598#L224)
``` go
func (r *Ring) Nodes() (nodes []string)
```
Nodes returns the list of nodes in the hash ring.




### <a name="Ring.RemoveNode">func</a> (\*Ring) [RemoveNode](/src/target/ringhash.go?s=4195:4238#L206)
``` go
func (r *Ring) RemoveNode(node string) bool
```
RemoveNode deletes node from the hash ring.




### <a name="Ring.UpdateLoad">func</a> (\*Ring) [UpdateLoad](/src/target/ringhash.go?s=3365:3415#L165)
``` go
func (r *Ring) UpdateLoad(node string, load int64)
```
UpdateLoad sets load of the given node to the given load.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
