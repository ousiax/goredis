// The MIT License (MIT)
//
// Copyright (c) 2016 Roy Xu

package redis_test

import (
	"github.com/qqbuby/goredis/redis"
	"testing"
)

const (
	network = "tcp"
	// address = "127.0.0.1:6379"
	address = "192.168.1.130:6379"
)

const (
	testKey   = "foobar"
	testValue = "foobuzz"
)

func TestSet(t *testing.T) {
	client, err := redis.NewClient(network, address)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	s, e := client.Set(testKey, testValue)
	if e != nil {
		t.Errorf("Set dit not work properly: %s", e.Error())
	}
	if s != "OK" {
		t.Errorf("Set dit not work properly. [%s]", s)
	}
}

func TestGet(t *testing.T) {
	client, err := redis.NewClient(network, address)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	args := make([]interface{}, 2)
	args[0] = testKey
	args[1] = testValue
	stat, _ := client.Send("SET", args...)
	if stat != "OK" {
		t.Fatalf("Set dit not work properly, result: %s, expected: %s.", stat, "OK")
	}
	v, _ := client.Get(testKey)
	if v != testValue {
		t.Fatalf("Set dit not work properly, result: [%s], expected: %s.", v, testValue)
	}
}

func TestDecr(t *testing.T) {
	client, err := redis.NewClient(network, address)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	client.Set("k", "100")
	v, e := client.Decr("k")
	if e != nil || v != 99 {
		t.Fatalf("Decr dit not work properly.")
	}
}

func TestDecrBy(t *testing.T) {
	client, err := redis.NewClient(network, address)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	client.Set("k", "100")
	v, e := client.DecrBy("k", 10)
	if e != nil || v != 90 {
		t.Fatalf("DecrBy dit not work properly.")
	}
}

func TestIncr(t *testing.T) {
	client, err := redis.NewClient(network, address)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	client.Set("k", "100")
	v, e := client.Incr("k")
	if e != nil || v != 101 {
		t.Fatalf("Incr dit not work properly.")
	}
}

func TestIncrBy(t *testing.T) {
	client, err := redis.NewClient(network, address)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	client.Set("k", "100")
	v, e := client.IncrBy("k", 10)
	if e != nil || v != 110 {
		t.Fatalf("IncrBy dit not work properly.")
	}
}

func TestIncrByFloat(t *testing.T) {
	client, err := redis.NewClient(network, address)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	client.Set("k", "100")
	v, e := client.IncrByFloat("k", float64(2.5))
	if e != nil || v != 102.5 {
		t.Fatalf("IncrByFloat dit not work properly.")
	}
}
