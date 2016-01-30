// The MIT License (MIT)

// Copyright (c) 2016 Roy Xu

package redis

type PubSub struct {
	cn Conn
}

func NewPubSub(url string) (PubSub, error) {
	c, err := Dial(url)
	cli := PubSub{cn: c}
	return cli, err
}

type Kind int

const (
	SUBSCRIBE Kind = iota
	UNSUBSCRIBE
	PSUBSCRIBE
	PUNSUBSCRIBE
)

var kinds = []string{
	"SUBSCRIBE",
	"UNSUBSCRIBE",
	"PSUBSCRIBE",
	"PUNSUBSCRIBE",
}

func (k Kind) String() string {
	return kinds[k]
}

type Subscription struct {
}

type Message struct {
	Channel interface{}
	Text    interface{}
}

type PMessage struct {
	Pattern interface{}
	Channel interface{}
	Text    interface{}
}

// PSUBSCRIBE pattern [pattern ...]
// Listen for messages published to channels matching the given patterns
func (p *PubSub) PSubscribe(pattern interface{}, patterns ...interface{}) error {
	p.cn.Pipe("PSUBSCRIBE", pattern)
	p.cn.Pipe("PSUBSCRIBE", patterns...)
	return p.cn.Flush()
}

// PUBSUB subcommand [argument [argument ...]]
// Inspect the state of the Pub/Sub subsystem

// PUBSUB CHANNELS [pattern]
// Lists the currently active channels.
// Array reply: a list of active channels, optionally matching the specified pattern.

// PUBSUB NUMSUB [channel-1 ... channel-N]
// Returns the number of subscribers (not counting clients subscribed to patterns) for the specified channels.
// Array reply:
//     a list of channels and number of subscribers for every channel.
//     The format is channel, count, channel, count, ..., so the list is flat.
//     The order in which the channels are listed is the same as the order of the channels specified in the command call.
// Note that it is valid to call this command without channels. In this case it will just return an empty list.

// PUBSUB NUMPAT
// Returns the number of subscriptions to patterns (that are performed using the PSUBSCRIBE command).
// Note that this is not just the count of clients subscribed to patterns but the total number of patterns all the clients are subscribed to.
// Integer reply: the number of patterns all the clients are subscribed to.

// PUBLISH channel message
// Post a message to a channel
// Integer reply: the number of clients that received the message.
func (p *PubSub) Publish(channel, message interface{}) (int, error) {
	return 0, nil
}

// PUNSUBSCRIBE [pattern [pattern ...]]
// Stop listening for messages posted to channels matching the given patterns
func (p *PubSub) PUnsubscribe(pattern interface{}, patterns ...interface{}) error {
	p.cn.Pipe("PUNSUBSCRIBE", pattern)
	p.cn.Pipe("PUNSUBSCRIBE", patterns...)
	return p.cn.Flush()
}

// SUBSCRIBE channel [channel ...]
// Listen for messages published to the given channels
func (p *PubSub) Subscribe(channel interface{}, channels ...interface{}) error {
	p.cn.Pipe("SUBSCRIBE", channel)
	p.cn.Pipe("SUBSCRIBE", channels...)
	return p.cn.Flush()
}

// UNSUBSCRIBE [channel [channel ...]]
// Stop listening for messages posted to the given channels
func (p *PubSub) Unsubscribe(channel interface{}, channels ...interface{}) error {
	p.cn.Pipe("UNSUBSCRIBE", channel)
	p.cn.Pipe("UNSUBSCRIBE", channels...)
	return p.cn.Flush()
}

func (p *PubSub) Receive() interface{} {
	r, e := p.cn.Receive()
	if e != nil {
		return e
	}
	return r
}

func (p *PubSub) Ping() (string, error) {
	v, err := p.cn.Send("PING")
	s, _ := v.(string)
	return s, err
}

func (p *PubSub) Close() error {
	return p.cn.Close()
}
