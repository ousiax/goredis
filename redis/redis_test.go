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

const (
	password = "password"
)

func TestAuth(t *testing.T) {
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

// [BEGIN] RESP STRINGS

const (
	testKey   = "foobar"
	testValue = "foobuzz"
)

func TestSet(t *testing.T) {
	if client == nil {
		t.Fatal("Hello world!")
	}
	s, e := client.Set(testKey, testValue)
	if e != nil {
		t.Errorf("Set dit not work properly: %s", e.Error())
	}
	if s != "OK" {
		t.Errorf("Set dit not work properly. [%s]", s)
	}
}

func TestGet(t *testing.T) {
	args := make([]interface{}, 2)
	args[0] = testKey
	args[1] = testValue
	stat, _ := client.Set(testKey, testValue)
	if stat != "OK" {
		t.Fatalf("Get dit not work properly, result: %s, expected: %s", stat, "OK")
	}
	v, _ := client.Get(testKey)
	if v != testValue {
		t.Fatalf("Get dit not work properly, result: [%s], expected: %s.", v, testValue)
	}
}

func TestDecr(t *testing.T) {
	client.Set("k", "100")
	v, e := client.Decr("k")
	if e != nil || v != 99 {
		t.Fatalf("Decr dit not work properly.")
	}
}

func TestDecrBy(t *testing.T) {
	client.Set("k", "100")
	v, e := client.DecrBy("k", 10)
	if e != nil || v != 90 {
		t.Fatalf("DecrBy dit not work properly.")
	}
}

func TestIncr(t *testing.T) {
	client.Set("k", "100")
	v, e := client.Incr("k")
	if e != nil || v != 101 {
		t.Fatalf("Incr dit not work properly.")
	}
}

func TestIncrBy(t *testing.T) {
	client.Set("k", "100")
	v, e := client.IncrBy("k", 10)
	if e != nil || v != 110 {
		t.Fatalf("IncrBy dit not work properly.")
	}
}

func TestIncrByFloat(t *testing.T) {
	client.Set("k", "100")
	v, e := client.IncrByFloat("k", float64(2.5))
	if e != nil {
		t.Fatalf("IncrByFloat dit not work properly. [%v]", e)
	}

	if v != 102.5 {
		t.Fatalf("IncrByFloat dit not work properly.")
	}
}

// [END] RESP STRINGS
