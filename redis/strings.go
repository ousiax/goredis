// The MIT License (MIT)
//
// Copyright (c) 2016 Roy Xu

package redis

type respStringer interface {
	// APPEND key value Append a value to a key
	Append(key, value string) (int, error)

	// BITCOUNT key [start end] Count set bits in a string
	BitCount(key string, p ...int) (int, error)

	// BITOP operation destkey key [key ...] Perform bitwise operations between strings
	BitOp(operation, destkey, key string, keys ...string) (int, error)

	// BITPOS key bit [start] [end] Find first bit set or clear in a string
	BitPOs(key string, bit int, p ...int) (int, error)

	// DECR key Decrement the integer value of a key by one
	Decr(key string) (int, error)

	// DECRBY key decrement Decrement the integer value of a key by the given number
	DecrBy(key string, decrement int) (int, error)

	// GET key Get the value of a key
	Get(key string) (string, error)

	// GETBIT key offset Returns the bit value at offset in the string value stored at key
	GetBit(key string, offset int) (int, error)

	// GETRANGE key start end Get a substring of the string stored at a key
	GetRange(key string, start, end int) (string, error)

	// GETSET key value Set the string value of a key and return its old value
	GetSet(key, value string) (string, error)

	// INCR key Increment the integer value of a key by one
	Incr(key string) (int, error)

	// INCRBY key increment Increment the integer value of a key by the given amount
	IncrBy(key string, decrement int) (int, error)

	// INCRBYFLOAT key increment Increment the float value of a key by the given amount
	IncrByFloat(key string, decrement float64) (float64, error)

	// MGET key [key ...] Get the values of all the given keys
	MGet(key string, keys ...string) ([]interface{}, error)

	// MSET key value [key value ...] Set multiple keys to multiple values
	MSet(key, value string, p ...string) (string, error)

	// MSETNX key value [key value ...] Set multiple keys to multiple values, only if none of the keys exist
	MSetNx(key, value string, p ...string) (string, error)

	// PSETEX key milliseconds value Set the value and expiration in milliseconds of a key
	PSetEx(key string, milliseconds int, value string) (string, error)

	// SET key value [EX seconds] [PX milliseconds] [NX|XX] Set the string value of a key
	Set(key, value string, p ...interface{}) (string, error)

	// SETBIT key offset value Sets or clears the bit at offset in the string value stored at key
	SetBit(key string, offset, value int) (int, error)

	// SETEX key seconds value Set the value and expiration of a key
	SetEx(key string, seconds int, value string) (string, error)

	// SETNX key value Set the value of a key, only if the key does not exist
	SetNx(key, value string) (int, error)

	// SETRANGE key offset value Overwrite part of a string at key starting at the specified offset
	SetRange(key string, offset, value int) (int, error)

	// STRLEN key Get the length of the value stored in a key
	StrLen(key string) (int, error)
}
