// The MIT License (MIT)

// Copyright (c) 2016 Roy Xu

package redis

import (
	"errors"
	"fmt"
	"strings"
)

type PubSub struct {
	cn       Conn
	channels []string
	numSub   int
}

func NewPubSub(url string) (PubSub, error) {
	c, err := Dial(url)
	cli := PubSub{cn: c, channels: []string{}}
	return cli, err
}

type PubSubMessage struct {
	Kind    string
	Channel string
	Num     int
}

func (sm PubSubMessage) String() string {
	return fmt.Sprintf("%s\n%s\n%d\n", sm.Kind, sm.Channel, sm.Num)
}

type Message struct {
	Channel string
	Text    string
}

func (m Message) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n", "MESSAGE", m.Channel, m.Text)
}

type PMessage struct {
	Pattern string
	Channel string
	Text    string
}

func (p PMessage) String() string {
	return fmt.Sprintf("%s\n%v\n%s\n%s\n", "PMESSAGE", p.Pattern, p.Channel, p.Text)
}

// Channels returns the list of the subscribed channels
func (p *PubSub) Channels() []string {
	return p.channels
}

// NumSub returns the number of subscribers for the all channels.
func (p *PubSub) NumSub() int {
	return p.numSub
}

// PSUBSCRIBE pattern [pattern ...]
// Listen for messages published to channels matching the given patterns
func (p *PubSub) PSubscribe(pattern string, patterns ...interface{}) error {
	p.cn.Pipe("SUBSCRIBE", MakeSlice(patterns, pattern)...)
	return p.cn.Flush()
}

// PUBLISH channel message
// Post a message to a channel
// Integer reply: the number of clients that received the message.
func (p *PubSub) Publish(channel, message string) (int, error) {
	r, e := p.cn.Send("PUBLISH", channel, message)
	v, _ := Int(r)
	return v, e
}

// PUNSUBSCRIBE [pattern [pattern ...]]
// Stop listening for messages posted to channels matching the given patterns
func (p *PubSub) PUnsubscribe(pattern string, patterns ...interface{}) error {
	p.cn.Pipe("SUBSCRIBE", MakeSlice(patterns, pattern)...)
	return p.cn.Flush()
}

// SUBSCRIBE channel [channel ...]
// Listen for messages published to the given channels
func (p *PubSub) Subscribe(channel string, channels ...interface{}) error {
	p.cn.Pipe("SUBSCRIBE", MakeSlice(channels, channel)...)
	return p.cn.Flush()
}

// UNSUBSCRIBE [channel [channel ...]]
// Stop listening for messages posted to the given channels
func (p *PubSub) Unsubscribe(channel string, channels ...interface{}) error {
	p.cn.Pipe("SUBSCRIBE", MakeSlice(channels, channel)...)
	return p.cn.Flush()
}

func (p *PubSub) Receive() interface{} {
	r, err := p.cn.Receive()
	if err != nil {
		return err
	}
	if v, err := r.([]interface{}); err {
		s, _ := String(v[0])
		switch k := strings.ToUpper(s); k {
		case "MESSAGE":
			c, _ := String(v[1])
			t, _ := String(v[2])
			return Message{Channel: c, Text: t}
		case "PMESSAGE":
			p, _ := String(v[1])
			c, _ := String(v[2])
			t, _ := String(v[3])
			return PMessage{Pattern: p, Channel: c, Text: t}
		case "SUBSCRIBE", "UNSUBSCRIBE", "PSUBSCRIBE", "PUNSUBSCRIBE":
			c, _ := String(v[1])
			n, _ := Int(v[2])
			return PubSubMessage{Kind: k, Channel: c, Num: n}
		default:
			return errors.New(fmt.Sprintf("pubsub.Receive: Protocol error. (%v)", v))
		}
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
