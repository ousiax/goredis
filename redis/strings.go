// The MIT License (MIT)
//
// Copyright (c) 2016 Roy Xu

package redis

type respStringer interface {
	// APPEND key value
	// Append a value to a key
	// Integer reply: the length of the string after the append operation.
	Append(key, value interface{}) (int, error)

	// BITCOUNT key [start end]
	// Count set bits in a string
	// Integer reply: The number of bits set to 1.
	BitCount(key interface{}, p ...int) (int, error)

	// BITOP operation destkey key [key ...]
	// Perform bitwise operations between strings
	// The BITOP command supports four bitwise operations: AND, OR, XOR and NOT.
	// Integer reply: The size of the string stored in the destination key, that is equal to the size of the longest input string.
	BitOp(operation, destkey, key interface{}, keys ...interface{}) (int, error)

	// BITPOS key bit [start] [end]
	// Find first bit set or clear in a string
	// Integer reply: The command returns the position of the first bit set to 1 or 0 according to the request.
	BitPOs(key interface{}, bit int, p ...int) (int, error)

	// DECR key
	// Decrement the integer value of a key by one
	// Integer reply: the value of key after the decrement
	Decr(key interface{}) (int, error)

	// DECRBY key decrement
	// Decrement the integer value of a key by the given number
	// Integer reply: the value of key after the decrement
	DecrBy(key interface{}, decrement int) (int, error)

	// GET key
	// Get the value of a key
	// Bulk string reply: the value of key, or nil when key does not exist.
	Get(key interface{}) (interface{}, error)

	// GETBIT key offset
	// Returns the bit value at offset in the string value stored at key
	// Integer reply: the bit value stored at offset.
	GetBit(key interface{}, offset int) (int, error)

	// GETRANGE key start end
	// Get a substring of the string stored at a key
	// Bulk string reply: the substring of the string stored
	GetRange(key interface{}, start, end int) (interface{}, error)

	// GETSET key value
	// Set the string value of a key and return its old value
	// Bulk string reply: the old value stored at key, or nil when key did not exist.
	GetSet(key, value interface{}) (interface{}, error)

	// INCR key
	// Increment the integer value of a key by one
	// Integer reply: the value of key after the increment
	Incr(key interface{}) (int, error)

	// INCRBY key increment
	// Increment the integer value of a key by the given amount
	// Integer reply: the value of key after the increment
	IncrBy(key interface{}, decrement int) (int, error)

	// INCRBYFLOAT key increment
	// Increment the float value of a key by the given amount
	// Bulk string reply: the value of key after the increment.
	IncrByFloat(key interface{}, decrement float64) (float64, error)

	// MGET key [key ...]
	// Get the values of all the given keys
	// Array reply: list of values at the specified keys.
	MGet(key interface{}, keys ...interface{}) ([]interface{}, error)

	// MSET key value [key value ...]
	// Set multiple keys to multiple values
	// Simple string reply: always OK since MSET can't fail.
	MSet(key, value interface{}, p ...interface{}) (string, error)

	// MSETNX key value [key value ...]
	// Set multiple keys to multiple values, only if none of the keys exist
	// Integer reply, specifically:
	//    1 if the all the keys were set.
	//    0 if no key was set (at least one key already existed).
	MSetNx(key, value interface{}, p ...interface{}) (int, error)

	// PSETEX key milliseconds value
	// Set the value and expiration in milliseconds of a key
	// Simple string reply
	PSetEx(key interface{}, milliseconds int, value interface{}) (string, error)

	// SET key value [EX seconds] [PX milliseconds] [NX|XX]
	// Set the string value of a key
	// Simple string reply: OK if SET was executed correctly.
	// Null reply: a Null Bulk Reply is returned if the SET operation was not performed because the user specified the NX or XX option but the condition was not met.
	Set(key, value interface{}, p ...interface{}) (string, error)

	// SETBIT key offset value
	// Sets or clears the bit at offset in the string value stored at key
	// Integer reply: the original bit value stored at offset.
	SetBit(key interface{}, offset, value int) (int, error)

	// SETEX key seconds value
	// Set the value and expiration of a key
	// Simple string reply
	SetEx(key interface{}, seconds int, value interface{}) (string, error)

	// SETNX key value
	// Set the value of a key, only if the key does not exist
	// Integer reply, specifically:
	//    1 if the key was set
	//    0 if the key was not set
	SetNx(key, value interface{}) (int, error)

	// SETRANGE key offset value
	// Overwrite part of a string at key starting at the specified offset
	// Integer reply: the length of the string after it was modified by the command.
	SetRange(key interface{}, offset, value int) (int, error)

	// STRLEN key
	// Get the length of the value stored in a key
	// Integer reply: the length of the string at key, or 0 when key does not exist.
	StrLen(key interface{}) (int, error)
}
