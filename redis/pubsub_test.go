// The MIT License (MIT)

// Copyright (c) 2016 Roy Xu

package redis_test

import (
	"github.com/qqbuby/goredis/redis"
	"reflect"
	"testing"
)

func TestPubSub(t *testing.T) {
	sc, err := redis.NewPubSub(url)
	if err != nil {
		t.Fatalf("Could not connect to Redis at %s: Connection refused.", url)
	}
	defer sc.Close()

	pc, err := redis.NewPubSub(url)
	if err != nil {
		t.Fatalf("Could not connect to Redis at %s: Connection refused.", url)
	}
	defer pc.Close()

	sc.Subscribe("c1")
	r := sc.Receive()
	if !reflect.DeepEqual(r, redis.PubSubMessage{Kind: "SUBSCRIBE", Channel: "c1", Num: 1}) {
		t.Fatal("TestPubSub failure.")
	}
	pc.Publish("c1", "Hi")
	m := sc.Receive()
	if !reflect.DeepEqual(m, redis.Message{Channel: "c1", Text: "Hi"}) {
		t.Fatal("TestPubSub failure.")
	}
}
