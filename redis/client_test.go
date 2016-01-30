// The MIT License (MIT)

// Copyright (c) 2016 Roy Xu

package redis_test

import (
	"github.com/qqbuby/goredis/redis"
	"os"
	"testing"
)

var client redis.Client

func setup() error {
	var err error
	client, err = redis.NewClient(url)
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
	client, e := redis.NewClient(url)
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

func TestDump(t *testing.T) {
	const (
		key   = "TEST:DUMP"
		value = key
	)
	client.Set(key, value)
	serial, _ := client.Dump(key)
	client.Del(key)
	client.Restore(key, 0, serial, false)
	s, _ := client.Get(key)
	if s != value {
		t.Errorf("Dump did not work properly. E:%s, R:%s", value, s)
	}
}

func TestExists(t *testing.T) {
	const (
		key0  = "TEST:EXISTS0"
		key1  = "TEST:EXISTS1"
		key2  = "EST:EXISTS2"
		value = "TEST:EXISTS"
	)
	client.Set(key0, value)
	client.Set(key1, value)
	client.Del(key2)
	r, _ := client.Exists(key0, key1, key2)
	if r != 2 {
		t.Errorf("Exists did not work properly. E:%s, R:%s", 2, r)
	}
}

func TestExpire(t *testing.T) {
	const (
		key     = "TEST:EXPIRE"
		value   = key
		seconds = 10
	)
	client.Del(key)
	r, _ := client.Expire(key, seconds)
	if r != 0 {
		t.Error("Expire did not work properly.")
	}
	client.Set(key, value)
	r, _ = client.Expire(key, seconds)
	if r != 1 {
		t.Error("Expire did not work properly.")
	}
}

func TestExpireAt(t *testing.T) {
	const (
		key       = "TEST:EXPIREAT"
		value     = key
		timestamp = 10000
	)
	client.Del(key)
	r, _ := client.ExpireAt(key, timestamp)
	if r != 0 {
		t.Error("ExpireAt did not work properly.")
	}
	client.Set(key, value)
	r, _ = client.ExpireAt(key, timestamp)
	if r != 1 {
		t.Error("ExpireAt did not work properly.")
	}
}

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

func TestMigrate(t *testing.T) {}

func TestMove(t *testing.T) {
	const (
		key   = "TEST:MOVE"
		value = key
		db    = 1
	)
	client.Select(0)
	client.Set(key, value)
	client.Move(key, db)
	v, _ := client.Get(key)
	if v != nil {
		t.Errorf("Move did not work properly. E:%s, R:%s", nil, v)
	}
	client.Select(db)
	v, _ = client.Get(key)
	if v != value {
		t.Errorf("Move did not work properly. E:%s, R:%s", value, v)
	}
}

func TestObject(t *testing.T) {}

func TestPersist(t *testing.T) {
	const (
		key   = "TEST:PERSIST"
		value = key
	)
	client.Del(key)
	r, _ := client.Persist(key)
	if r != 0 {
		t.Error("Persist did not work properly.")
	}
	client.Set(key, value)
	r, _ = client.Persist(key)
	if r != 0 {
		t.Error("Persist did not work properly.")
	}
	client.Expire(key, 10000)
	r, _ = client.Persist(key)
	if r != 1 {
		v, _ := client.Get(key)
		t.Errorf("Persist did not work properly. V:%s", v)
	}
}

func TestPexpire(t *testing.T) {
	const (
		key     = "TEST:PEXPIRE"
		value   = key
		seconds = 10
	)
	client.Del(key)
	r, _ := client.PExpire(key, seconds)
	if r != 0 {
		t.Error("PExpire did not work properly.")
	}
	client.Set(key, value)
	r, _ = client.PExpire(key, seconds)
	if r != 1 {
		t.Error("PExpire did not work properly.")
	}
}

func TestPexpireAt(t *testing.T) {
	const (
		key       = "TEST:PEXPIREAT"
		value     = key
		timestamp = 10000
	)
	client.Del(key)
	r, _ := client.PExpireAt(key, timestamp)
	if r != 0 {
		t.Error("PExpireAt did not work properly.")
	}
	client.Set(key, value)
	r, _ = client.PExpireAt(key, timestamp)
	if r != 1 {
		t.Error("PExpireAt did not work properly.")
	}
}

func TestPTtl(t *testing.T) {
	const (
		key   = "TEST:PTTL"
		value = key
	)
	client.Del(key)
	v, _ := client.Pttl(key)
	if v != -2 {
		t.Errorf("Keys did not work properly. E:%s, R:%s", -2, v)
	}
	client.Set(key, value)
	v, _ = client.Pttl(key)
	if v != -1 {
		t.Errorf("Keys did not work properly. E:%s, R:%s", -1, v)
	}
	client.Expire(key, 10000)
	v, _ = client.Pttl(key)
	if v <= 0 {
		t.Errorf("Keys did not work properly. R:%s", v)
	}
}

func TestRandomKey(t *testing.T) {
	const (
		key   = "TEST:RANDOMKEY"
		value = key
	)
	client.Set(key, value)
	k, _ := client.RandomKey()
	if k == nil {
		t.Errorf("RandomKey did not work properly. R:%v", k)
	}
}

func TestRename(t *testing.T) {
	const (
		key    = "TEST:RENAME"
		value  = key
		newkey = "TEST:RENAME:NEWKEY"
	)
	client.Set(key, value)
	client.Rename(key, newkey)
	v, _ := client.Get(newkey)
	if v != value {
		t.Errorf("Rename did not work properly. E:%s, R:%s", value, v)
	}
	client.Del(key)
}

func TestRenameNx(t *testing.T) {
	const (
		key    = "TEST:RENAMENX"
		value  = key
		newkey = "TEST:RENAMENX:NEWKEY"
	)
	client.Set(key, value)
	client.Set(newkey, value)
	s, _ := client.RenameNx(key, newkey)
	if s == 1 {
		t.Errorf("RenameNx did not work properly. R:%s", s)
	}
	client.Del(newkey)
	s, _ = client.RenameNx(key, newkey)
	if s == 0 {
		t.Errorf("Rename did not work properly. R:%s", s)
	}
}

func TestRestore(t *testing.T) {
	const (
		key   = "TEST:RESTORE"
		value = key
	)
	client.Set(key, value)
	serial, _ := client.Dump(key)
	e, _ := client.Restore(key, 0, serial, false)
	if e == "OK" {
		t.Error("Restore did not work properly.")
	}
	client.Del(key)
	client.Restore(key, 0, serial, false)
	s, _ := client.Get(key)
	if s != value {
		t.Errorf("Restore did not work properly. E:%s, R:%s", value, s)
	}
}

func TestSort(t *testing.T) {}

func TestTtl(t *testing.T) {
	const (
		key   = "TEST:TTL"
		value = key
	)
	client.Del(key)
	v, _ := client.Ttl(key)
	if v != -2 {
		t.Errorf("Ttl did not work properly. E:%s, R:%s", -2, v)
	}
	client.Set(key, value)
	v, _ = client.Ttl(key)
	if v != -1 {
		t.Errorf("Ttl did not work properly. E:%s, R:%s", -1, v)
	}
	client.Expire(key, 10000)
	v, _ = client.Ttl(key)
	if v <= 0 {
		t.Errorf("Ttl did not work properly. R:%s", v)
	}
}

func TestType(t *testing.T) {
	const (
		key   = "TEST:TYPE"
		value = key
	)
	client.Set(key, value)
	v, _ := client.Type(key)
	if v != "string" {
		t.Errorf("Type did not work properly. E:%s, R:%s", "string", v)
	}
	client.Del(key)
	v, _ = client.Type(key)
	if v != "none" {
		t.Errorf("Type did not work properly. E:%s, R:%s", "none", v)
	}
}

func TestWait(t *testing.T) {}

func TestScan(t *testing.T) {}

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

func TestBitOp(t *testing.T) {
	const (
		operation = "AND"
		destkey   = "TEST:BITOP:DESTKEY"
		key1      = "TEST:BITOP:KEY1"
		key2      = "TEST:BITOP:KEY2"
		value1    = "foobar"
		value2    = "abcdef"
		value3    = "`bc`ab"
	)
	client.Set(key1, value1)
	client.Set(key2, value2)
	client.BitOp(operation, destkey, key1, key2)
	s, _ := client.Get(destkey)
	if s != value3 {
		t.Errorf("BitOp did not work properly. E:%s, R:%s", value3, s)
	}
}

func TestBitPOs(t *testing.T) {
	const (
		key = "TEST:BITPOS"
		bit = 1
		pos = 7
	)
	client.SetBit(key, pos, bit)
	p, _ := client.BitPOs(key, bit)
	if p != pos {
		t.Errorf("BitPOs did not work properly. E:%s, R:%s", pos, p)
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

func TestGetBit(t *testing.T) {
	const (
		key = "TEST:GETBIT"
		bit = 1
		pos = 7
	)
	client.SetBit(key, pos, bit)
	b, _ := client.GetBit(key, pos)
	if b != bit {
		t.Errorf("GetBit did not work properly. E:%s, R:%s", bit, b)
	}
}

func TestGetRange(t *testing.T) {
	const (
		key   = "TEST:GETRANGE"
		value = key
		start = 5
		end   = len(key)
	)
	client.Set(key, value)
	s, _ := client.GetRange(key, start, end)
	if s != value[start:end] {
		t.Errorf("GetRange did not work properly. E:%s, R:%s", value[start:end], s)
	}
}

func TestGetSet(t *testing.T) {
	const (
		key      = "TEST:GETSET"
		value    = key
		newValue = "TEST:GETSET:OLD"
	)
	v, _ := client.GetSet(key, value)
	if v != nil {
		t.Errorf("GetSet did not work properly. R:%s", v)
	}
	v, _ = client.GetSet(key, newValue)
	if v != value {
		t.Errorf("GetSet did not work properly. E:%s, R:%s", value, v)
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

func TestMGet(t *testing.T) {
	const (
		key1   = "TEST:MGET:KEY1"
		value1 = key1
		key2   = "TEST:MGET:KEY2"
		value2 = key2
	)
	client.MSet(key1, value1, key2, value2)
	v, _ := client.MGet(key1, key2)
	if !(string(v[0].([]byte)) == value1 && string(v[1].([]byte)) == value2) {
		t.Error("MGet did not work properly.")
	}
}

func TestMSet(t *testing.T) {
	const (
		key1   = "TEST:MSET:KEY1"
		value1 = key1
		key2   = "TEST:MSET:KEY2"
		value2 = key2
	)
	v, _ := client.MSet(key1, value1, key2, value2)
	if v != "OK" {
		t.Errorf("MSet did not work properly. E:%s, R:%s", "OK", v)
	}
}

func TestMSetNx(t *testing.T) {
	const (
		key1   = "TEST:MSETNX:KEY1"
		value1 = key1
		key2   = "TEST:MSETNX:KEY2"
		value2 = key2
	)
	v, _ := client.MSetNx(key1, value1, key2, value2)
	if v != 1 {
		t.Errorf("MSetNx did not work properly. E:%d, R:%d", 1, v)
	}
	v, _ = client.MSetNx(key1, value1, key2, value2)
	if v != 0 {
		t.Errorf("MSetNx did not work properly. E:%d, R:%d", 0, v)
	}
}

func TestPSetEx(t *testing.T) {
	const (
		key     = "TEST:PSETEX"
		value   = key
		seconds = 1000
	)
	client.Set(key, value)
	_, e := client.PSetEx(key, -1, value)
	if e == nil {
		t.Error("PSetEx did not work properly.")
	}

	s, _ := client.PSetEx(key, seconds, value)
	if s != "OK" {
		t.Errorf("PSetEx did not work properly. E:%s, R:%s", "OK", s)
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

func TestSetEx(t *testing.T) {
	const (
		key     = "TEST:SETEX"
		value   = key
		seconds = 1000
	)
	client.Set(key, value)
	_, e := client.SetEx(key, -1, value)
	if e == nil {
		t.Error("SetEx did not work properly.")
	}
	s, _ := client.SetEx(key, seconds, value)
	if s != "OK" {
		t.Errorf("SetEx did not work properly. E:%s, R:%s", "OK", s)
	}
}

func TestSetNx(t *testing.T) {
	const (
		key   = "TEST:SETNX"
		value = key
	)
	i, _ := client.SetNx(key, value)
	if i == 0 {
		t.Errorf("SetNx did not work properly. E:%d, R:%d", 1, i)
	}

	i, _ = client.SetNx(key, value)
	if i == 1 {
		t.Errorf("SetNx did not work properly. E:%d, R:%d", 0, i)
	}
}

func TestSetRange(t *testing.T) {
	const (
		key    = "TEST:SETRANGE"
		value  = key
		offset = 10
	)
	i, _ := client.SetRange(key, offset, value)
	if i != len(value)+offset {
		t.Errorf("SetRange did not work properly. E:%d, R:%d", len(value)+offset, i)
	}
}

func TestStrLen(t *testing.T) {
	const (
		key   = "TEST:STRLEN"
		value = key
	)
	client.Set(key, value)
	i, _ := client.StrLen(key)
	if i != len(value) {
		t.Errorf("SetRange did not work properly. E:%d, R:%d", len(value), i)
	}
}

// [END] RESP STRINGS

// [BEGIN] MISCELLANEOUS

func TestReceive(t *testing.T) {
	const (
		key   = "TEST:RECEIVE"
		value = key
	)
	_, e := client.Restore(key, -1, value, true)
	if e == nil {
		t.Errorf("Miscellaneous did not work properly.")
	}
}

// [END] MISCELLANEOUS
