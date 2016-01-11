package redis

import (
	"fmt"
	"testing"
)

const (
	network = "tcp"
	address = "192.168.1.130:6379"
)

const (
	testKey   = "foobar"
	testValue = "foobuzz"
)

func TestSet(t *testing.T) {
	client, err := NewClient(network, address)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	s, e := client.Set(testKey, testValue)
	if e != nil {
		t.Errorf("Set dit not work properly: %s", e.Error())
	}
	if s != "OK" {
		t.Error("Set dit not work properly")
	}
}

func TestGet(t *testing.T) {
	client, err := NewClient(network, address)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	stat, _ := client.Send(fmt.Sprintf("*3\r\n$3\r\nset\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(testKey), testKey, len(testValue), testValue))
	if stat != "OK" {
		t.Fatalf("redis.TestGet: result: %s, expected: %s.", stat, "OK")
	}
	v, _ := client.Get("testKey")
	if v != "testValue" {
		t.Fatalf("redis.TestGet: result: %s, expected: %s.", v, testValue)
	}
}
