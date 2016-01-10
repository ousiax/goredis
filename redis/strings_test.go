package redis

import (
	"testing"
)

const (
	network = "tcp"
	address = "192.168.1.130:6379"
)

func TestGet(t *testing.T) {
	client, err := NewClient(network, address)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	client.Send("*3\r\n$3\r\nset\r\n$8\r\ntest_key\r\n$10\r\ntest_value\r\n")
	stat, _ := client.Receive()
	if stat != "OK" {
		t.Fatalf("redis.TestGet: result: %s, expected: %s.", stat, "OK")
	}
	v, _ := client.Get("testKey")
	if v != "testValue" {
		t.Fatalf("redis.TestGet: result: %s, expected: %s.", v, "test_value")
	}
}
