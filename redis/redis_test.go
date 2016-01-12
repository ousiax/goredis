// The MIT License (MIT)
//
// Copyright (c) 2016 Roy Xu

package redis_test

import (
	"flag"
	"github.com/qqbuby/goredis/redis"
	"os"
	"testing"
)

var client redis.Client

const (
	network = "tcp"
	// address = "127.0.0.1:6379"
	address = "192.168.1.130:6379"
)

func setup() error {
	flag.Parse()
	var err error
	client, err = redis.NewClient(network, address)
	return err
}

func tearDown() error {
	keys, _ := client.Keys("TEST:*")
	client.Del(keys[0], keys[1:]...)
	client.Quit()
	return client.Close()
}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		os.Exit(-1)
	}
	code := m.Run()
	tearDown()
	os.Exit(code)
}

// [BEGIN] RESP CONNECTION

func TestAuth(t *testing.T) {
	const (
		password = "password"
	)
	rsp, err := client.Auth(password)
	if err != nil {
		t.Logf("Auth: %v", err)
	} else if rsp != "OK" {
		t.Logf("Auth: %v", rsp)
	}
}

func TestEcho(t *testing.T) {
	const message = "Hello world!"
	rsp, err := client.Echo(message)
	if err != nil {
		t.Errorf("Echo does not work properly. R:%s", err.Error())
	}
	if rsp != message {
		t.Errorf("Echo does not work properly. E:%s,R:%s", message, rsp)
	}
}

func TestPing(t *testing.T) {
	const message = "PONG"
	rsp, err := client.Ping()
	if err != nil {
		t.Errorf("Ping does not work properly. R:%s", err.Error())
	}
	if rsp != message {
		t.Errorf("Ping does not work properly. E:%s,R:%s", message, rsp)
	}
}

func TestSelect(t *testing.T) {
	rsp, err := client.Select(0)
	if err != nil {
		t.Errorf("Select does not work properly. R:%s", err.Error())
	} else if rsp != "OK" {
		t.Errorf("Select does not work properly. E:%s,R:%s", "OK", rsp)
	}
}

func TestQuit(t *testing.T) {
	client, e := redis.NewClient(network, address)
	if e != nil {
		t.Errorf("Quit: %s", e.Error())
	}
	defer client.Close()

	rsp, err := client.Quit()
	if err != nil {
		t.Errorf("Quit does not work properly. R:%s", err.Error())
	} else if rsp != "OK" {
		t.Errorf("Quit does not work properly. E:%s,R:%s", "OK", rsp)
	}
}

// [END] RESP CONNECTION

// [BEGIN] RESP KEYS

func TestKeys(t *testing.T) {
	const (
		key     = "TEST:KEYS"
		pattern = key
	)
	client.Set(key, pattern)
	rsp, _ := client.Keys(pattern)
	s := rsp[0].(string)
	if s != key {
		t.Errorf("Keys did not work properly. E:%s, R:%s", key, s)
	}

}

func TestDel(t *testing.T) {
	const (
		key   = "TEST:DEL"
		value = "value"
	)
	client.Set(key, value)
	s, _ := client.Del(key)
	if s != 1 {
		t.Error("Del did not work properly.")
	}
}

// [END] RESP KEYS

// [BEGIN] RESP STRINGS

func TestAppend(t *testing.T) {
	const (
		key   = "TEST:APPEND"
		value = "foobuzz"
	)
	client.Set(key, value)
	s, _ := client.Append(key, value)
	if s != len(value)*2 {
		t.Errorf("Append dit not work properly. [%d]", s)
	}
}

func TestSetBit(t *testing.T) {
	const (
		key    = "TEST:SETBIT"
		offset = 10
		bit    = 1
	)
	client.SetBit(key, offset, bit)
	b, _ := client.SetBit(key, offset, bit)
	if b != bit {
		t.Error("SetBit did not work properly.")
	}
}

func TestBitCount(t *testing.T) {
	const (
		key   = "TEST:BITCOUNT"
		count = 10
		bit   = 1
	)
	for i := 0; i < count; i++ {
		client.SetBit(key, i, 1)
	}
	s, _ := client.BitCount(key)
	if s != count {
		t.Errorf("BitCount dit not work properly. R:%d", s)
	}
}

func TestSet(t *testing.T) {
	const (
		key   = "TEST:SET"
		value = "foobuzz"
	)
	s, e := client.Set(key, value)
	if e != nil {
		t.Errorf("Set dit not work properly: %s", e.Error())
	}
	if s != "OK" {
		t.Errorf("Set dit not work properly. [%s]", s)
	}
}

func TestGet(t *testing.T) {
	const (
		key   = "TEST:GET"
		value = "foobuzz"
	)
	args := make([]interface{}, 2)
	args[0] = key
	args[1] = value
	stat, _ := client.Set(key, value)
	if stat != "OK" {
		t.Fatalf("Get dit not work properly, result: %s, expected: %s", stat, "OK")
	}
	v, _ := client.Get(key)
	if v != value {
		t.Fatalf("Get dit not work properly, result: [%s], expected: %s.", v, value)
	}
}

func TestDecr(t *testing.T) {
	const (
		key   = "TEST:DECR"
		value = 100
	)
	client.Set(key, value)
	v, _ := client.Decr(key)
	if v != value-1 {
		t.Fatalf("Decr dit not work properly.")
	}
}

func TestDecrBy(t *testing.T) {
	const (
		key       = "TEST:DECRBY"
		value     = 100
		decrement = 10
	)
	client.Set(key, value)
	v, _ := client.DecrBy(key, decrement)
	if v != value-decrement {
		t.Fatalf("DecrBy dit not work properly.")
	}
}

func TestIncr(t *testing.T) {
	const (
		key   = "TEST:INCR"
		value = 100
	)
	client.Set(key, value)
	v, e := client.Incr(key)
	if e != nil || v != value+1 {
		t.Fatalf("Incr dit not work properly.")
	}
}

func TestIncrBy(t *testing.T) {
	const (
		key       = "TEST:INCRBY"
		value     = 100
		increment = 10
	)
	client.Set(key, value)
	v, e := client.IncrBy(key, increment)
	if e != nil || v != value+increment {
		t.Fatalf("IncrBy dit not work properly.")
	}
}

func TestIncrByFloat(t *testing.T) {
	const (
		key       = "TEST:INCRBYFLOAT"
		value     = 100
		increment = 1.25
	)
	client.Set(key, value)
	v, _ := client.IncrByFloat(key, increment)
	if v != value+increment {
		t.Fatalf("IncrByFloat dit not work properly.")
	}
}

// [END] RESP STRINGS
