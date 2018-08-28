

# pubsub
`import "github.com/andy2046/gopie/pkg/pubsub"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package pubsub provides a pubsub implementation.




## <a name="pkg-index">Index</a>
* [func UUID() string](#UUID)
* [type Message](#Message)
* [type PubSub](#PubSub)
  * [func New(project string) *PubSub](#New)
  * [func (p *PubSub) Name() string](#PubSub.Name)
  * [func (p *PubSub) NewTopic(name string, size int, numGoroutines int) (*Topic, error)](#PubSub.NewTopic)
  * [func (p *PubSub) Topic(name string) *Topic](#PubSub.Topic)
  * [func (p *PubSub) Topics() []string](#PubSub.Topics)
* [type PublishError](#PublishError)
  * [func (pe PublishError) Error() string](#PublishError.Error)
* [type Subscription](#Subscription)
  * [func (s *Subscription) Delete()](#Subscription.Delete)
  * [func (s *Subscription) Receive(f func(*Message))](#Subscription.Receive)
* [type Topic](#Topic)
  * [func (t *Topic) Delete()](#Topic.Delete)
  * [func (t *Topic) Name() string](#Topic.Name)
  * [func (t *Topic) NewSubscription(numGoroutines int) (*Subscription, error)](#Topic.NewSubscription)
  * [func (t *Topic) Publish(ctx context.Context, msg *Message) error](#Topic.Publish)
  * [func (t *Topic) Stop()](#Topic.Stop)
  * [func (t *Topic) Subscription(name string) *Subscription](#Topic.Subscription)
  * [func (t *Topic) Subscriptions() []string](#Topic.Subscriptions)


#### <a name="pkg-files">Package files</a>
[pubsub.go](/src/github.com/andy2046/gopie/pkg/pubsub/pubsub.go) 





## <a name="UUID">func</a> [UUID](/src/target/pubsub.go?s=1639:1657#L90)
``` go
func UUID() string
```
UUID generates uuid.




## <a name="Message">type</a> [Message](/src/target/pubsub.go?s=1314:1361#L73)
``` go
type Message struct {
    ID   string
    Data []byte
}
```
Message represents a Pub/Sub message.










## <a name="PubSub">type</a> [PubSub](/src/target/pubsub.go?s=201:293#L15)
``` go
type PubSub struct {
    // contains filtered or unexported fields
}
```
PubSub is a Pub/Sub instance for a single project.







### <a name="New">func</a> [New](/src/target/pubsub.go?s=1808:1840#L97)
``` go
func New(project string) *PubSub
```
New creates a new PubSub.





### <a name="PubSub.Name">func</a> (\*PubSub) [Name](/src/target/pubsub.go?s=1972:2002#L105)
``` go
func (p *PubSub) Name() string
```
Name returns the full name for the PubSub.




### <a name="PubSub.NewTopic">func</a> (\*PubSub) [NewTopic](/src/target/pubsub.go?s=2253:2336#L112)
``` go
func (p *PubSub) NewTopic(name string, size int, numGoroutines int) (*Topic, error)
```
NewTopic creates a new Topic with the given name,
size is the channel buffer size for topic message chan,
numGoroutines is the number of goroutines it will spawn to push msg concurrently.




### <a name="PubSub.Topic">func</a> (\*PubSub) [Topic](/src/target/pubsub.go?s=2807:2849#L134)
``` go
func (p *PubSub) Topic(name string) *Topic
```
Topic returns the topic by name.




### <a name="PubSub.Topics">func</a> (\*PubSub) [Topics](/src/target/pubsub.go?s=2996:3030#L144)
``` go
func (p *PubSub) Topics() []string
```
Topics list all the topics in the PubSub.




## <a name="PublishError">type</a> [PublishError](/src/target/pubsub.go?s=1440:1491#L79)
``` go
type PublishError struct {
    Msg *Message
    Err error
}
```
PublishError is the error generated when it fails to publish a message.










### <a name="PublishError.Error">func</a> (PublishError) [Error](/src/target/pubsub.go?s=1495:1532#L85)
``` go
func (pe PublishError) Error() string
```



## <a name="Subscription">type</a> [Subscription](/src/target/pubsub.go?s=890:1269#L51)
``` go
type Subscription struct {
    // contains filtered or unexported fields
}
```
Subscription represents a PubSub subscription.










### <a name="Subscription.Delete">func</a> (\*Subscription) [Delete](/src/target/pubsub.go?s=6366:6397#L316)
``` go
func (s *Subscription) Delete()
```
Delete unsubscribes itself from topic.




### <a name="Subscription.Receive">func</a> (\*Subscription) [Receive](/src/target/pubsub.go?s=5985:6033#L292)
``` go
func (s *Subscription) Receive(f func(*Message))
```
Receive receives message for this subscription.




## <a name="Topic">type</a> [Topic](/src/target/pubsub.go?s=333:836#L22)
``` go
type Topic struct {

    // Errors is the error output channel back to the user. You MUST read from this
    // channel or the Publish will deadlock when the channel is full.
    Errors chan PublishError
    // contains filtered or unexported fields
}
```
Topic represents a PubSub topic.










### <a name="Topic.Delete">func</a> (\*Topic) [Delete](/src/target/pubsub.go?s=3628:3652#L177)
``` go
func (t *Topic) Delete()
```
Delete removes itself from PubSuband stop it.




### <a name="Topic.Name">func</a> (\*Topic) [Name](/src/target/pubsub.go?s=3845:3874#L189)
``` go
func (t *Topic) Name() string
```
Name returns the full name for the topic.




### <a name="Topic.NewSubscription">func</a> (\*Topic) [NewSubscription](/src/target/pubsub.go?s=5230:5303#L260)
``` go
func (t *Topic) NewSubscription(numGoroutines int) (*Subscription, error)
```
NewSubscription creates a new Subscription to this topic,
numGoroutines is the number of goroutines it will spawn to pull msg concurrently.




### <a name="Topic.Publish">func</a> (\*Topic) [Publish](/src/target/pubsub.go?s=3228:3292#L155)
``` go
func (t *Topic) Publish(ctx context.Context, msg *Message) error
```
Publish publishes msg to the topic asynchronously.




### <a name="Topic.Stop">func</a> (\*Topic) [Stop](/src/target/pubsub.go?s=3978:4000#L194)
``` go
func (t *Topic) Stop()
```
Stop stops the topic.




### <a name="Topic.Subscription">func</a> (\*Topic) [Subscription](/src/target/pubsub.go?s=5770:5825#L282)
``` go
func (t *Topic) Subscription(name string) *Subscription
```
Subscription returns the subscription by name..




### <a name="Topic.Subscriptions">func</a> (\*Topic) [Subscriptions](/src/target/pubsub.go?s=4882:4922#L248)
``` go
func (t *Topic) Subscriptions() []string
```
Subscriptions list all the subscriptions to this topic.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
