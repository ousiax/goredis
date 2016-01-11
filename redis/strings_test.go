package redis_test

import (
	"github.com/qqbuby/redis-go/redis"
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
		t.Fatalf("redis.TestGet: result: %s, expected: %s.", stat, "OK")
	}
	v, _ := client.Get(testKey)
	if v != testValue {
		t.Fatalf("redis.TestGet: result: [%s], expected: %s.", v, testValue)
	}
}
