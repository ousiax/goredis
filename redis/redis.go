// The MIT License (MIT)

// Copyright (c) 2016 Roy Xu

package redis

import (
	"errors"
)

type Client struct {
	cn Conn
}

func NewClient(urlstring string) (Client, error) {
	c, err := Dial(urlstring)
	cli := Client{cn: c}
	return cli, err
}

func (cli *Client) Close() (err error) {
	return cli.cn.Close()
}

func (cli *Client) Send(cmd string, args ...interface{}) (reply interface{}, err error) {
	rsp, err := cli.cn.Send(cmd, args...)
	return rsp, err
}

// CONNECTION:BEGIN

// AUTH password
// Authenticate to the server
// Simple string reply
func (cli *Client) Auth(password string) (string, error) {
	rsp, err := cli.Send("AUTH", password)
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// ECHO message
// Echo the given string
// Simple string reply
func (cli *Client) Echo(message string) (string, error) {
	rsp, err := cli.Send("ECHO", message)
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// PING
// Ping the server
// Simple string reply
func (cli *Client) Ping() (string, error) {
	rsp, err := cli.Send("PING")
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// QUIT
// Close the connection
// Simple string reply: always OK.
func (cli *Client) Quit() (string, error) {
	rsp, err := cli.Send("QUIT")
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// SELECT index
// Change the selected database for the current connection
// Simple string reply
func (cli *Client) Select(index int) (string, error) {
	rsp, err := cli.Send("SELECT", index)
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// CONNECTION:END

// KEYS:BEGIN

// DEL key [key ...]
// Delete a key
// Integer reply: The number of keys that were removed.
func (cli *Client) Del(key interface{}, keys ...interface{}) (int, error) {
	p := constructParameters(keys, key)
	rsp, err := cli.Send("DEL", p...)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// DUMP key
// Return a serialized version of the value stored at the specified key.
// Bulk string reply: the serialized value.
func (cli *Client) Dump(key interface{}) (interface{}, error) {
	rsp, err := cli.Send("DUMP", key)
	if err != nil {
		return nil, err
	}
	v, e := parseStringEx(rsp)
	return v, e
}

// EXISTS key [key ...]
// Determine if a key exists
// Integer reply: The number of keys existing among the ones specified as arguments.
func (cli *Client) Exists(key interface{}, keys ...interface{}) (int, error) {
	p := constructParameters(keys, key)
	rsp, err := cli.Send("EXISTS", p...)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// EXPIRE key seconds
// Set a key's time to live in seconds
// Integer reply, specifically:
//     1 if the timeout was set.
//     0 if key does not exist or the timeout could not be set.
func (cli *Client) Expire(key interface{}, seconds int) (int, error) {
	rsp, err := cli.Send("EXPIRE", key, seconds)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// EXPIREAT key timestamp
// Set the expiration for a key as a UNIX timestamp
// Integer reply, specifically:
//     1 if the timeout was set.
//     0 if key does not exist or the timeout could not be set.
func (cli *Client) ExpireAt(key interface{}, timestamp int) (int, error) {
	rsp, err := cli.Send("EXPIREAT", key, timestamp)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// KEYS pattern
// Find all keys matching the given pattern
// Array reply: list of keys matching pattern.
func (cli *Client) Keys(pattern interface{}) ([]interface{}, error) {
	rsp, err := cli.Send("KEYS", pattern)
	if err != nil {
		return nil, err
	}
	v, e := parseStrings(rsp)
	return v, e
}

// MIGRATE host port key destination-db timeout [COPY] [REPLACE]
// Atomically transfer a key from a Redis instance to another one.
// Simple string reply: The command returns OK on success.

// MOVE key db
// Move a key to another database
// Integer reply, specifically:
//     1 if key was moved.
//     0 if key was not moved.
func (cli *Client) Move(key interface{}, db int) (int, error) {
	rsp, err := cli.Send("MOVE", key, db)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// OBJECT subcommand [arguments [arguments ...]]
// Inspect the internals of Redis objects
// The OBJECT command supports multiple sub commands:
//     OBJECT REFCOUNT <key> returns the number of references of the value associated with the specified key. This command is mainly useful for debugging.
//     OBJECT ENCODING <key> returns the kind of internal representation used in order to store the value associated with a key.
//     OBJECT IDLETIME <key> returns the number of seconds since the object stored at the specified key is idle (not requested by read or write operations).
//     While the value is returned in seconds the actual resolution of this timer is 10 seconds, but may vary in future implementations.
// Objects can be encoded in different ways:
//     Strings can be encoded as raw (normal string encoding) or int (strings representing integers in a 64 bit signed interval are encoded in this way in order to save space).
//     Lists can be encoded as ziplist or linkedlist. The ziplist is the special representation that is used to save space for small lists.
//     Sets can be encoded as intset or hashtable. The intset is a special encoding used for small sets composed solely of integers.
//     Hashes can be encoded as zipmap or hashtable. The zipmap is a special encoding used for small hashes.
//     Sorted Sets can be encoded as ziplist or skiplist format. As for the List type small sorted sets can be specially encoded using ziplist, while the skiplist encoding is the one that works with sorted sets of any size.
// Different return values are used for different subcommands.
//     Subcommands refcount and idletime return integers.
//     Subcommand encoding returns a bulk reply.
// If the object you try to inspect is missing, a null bulk reply is returned.
// Object(subcommand string, arguments ...interface{}) (interface{}, error)

// PERSIST key
// Remove the expiration from a key
// Integer reply, specifically:
//     1 if the timeout was removed.
//     0 if key does not exist or does not have an associated timeout.
func (cli *Client) Persist(key interface{}) (int, error) {
	rsp, err := cli.Send("PERSIST", key)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// PEXPIRE key milliseconds
// Set a key's time to live in milliseconds
// Integer reply, specifically:
//     1 if the timeout was set.
//     0 if key does not exist or the timeout could not be set.
func (cli *Client) PExpire(key string, milliseconds int) (int, error) {
	rsp, err := cli.Send("PEXPIRE", key, milliseconds)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// PEXPIREAT key milliseconds-timestamp
// Set the expiration for a key as a UNIX timestamp specified in milliseconds
// Integer reply, specifically:
//     1 if the timeout was set.
//     0 if key does not exist or the timeout could not be set (see: EXPIRE).
func (cli *Client) PExpireAt(key string, millisecondTimestamp int) (int, error) {
	rsp, err := cli.Send("PEXPIREAT", key, millisecondTimestamp)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// PTTL key
// Get the time to live for a key in milliseconds
// Integer reply: TTL in milliseconds, or a negative value in order to signal an error.
func (cli *Client) Pttl(key string) (int, error) {
	rsp, err := cli.Send("PTTL", key)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// RANDOMKEY
// Return a random key from the keyspace
// Bulk string reply: the random key, or nil when the database is empty.
func (cli *Client) RandomKey() (interface{}, error) {
	rsp, err := cli.Send("RANDOMKEY")
	if err != nil {
		return nil, err
	}
	v, e := parseStringEx(rsp)
	return v, e
}

// RENAME key newkey
// Rename a key
// Simple string reply
func (cli *Client) Rename(key, newkey string) (string, error) {
	rsp, err := cli.Send("RENAME", key, newkey)
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// RENAMENX key newkey
// Rename a key, only if the new key does not exist
// Integer reply, specifically:
//     1 if key was renamed to newkey.
//     0 if newkey already exists.
func (cli *Client) RenameNx(key, newkey string) (int, error) {
	rsp, err := cli.Send("RENAMENX", key, newkey)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// RESTORE key ttl serialized-value [REPLACE]
// Create a key using the provided serialized value, previously obtained using DUMP.
// Simple string reply: The command returns OK on success.
func (cli *Client) Restore(key interface{}, ttl int, serializedValue interface{}, replace bool) (string, error) {
	var rsp interface{}
	var err error
	if replace {
		rsp, err = cli.Send("RESTORE", key, ttl, serializedValue, "REPLACE")
	} else {
		rsp, err = cli.Send("RESTORE", key, ttl, serializedValue)
	}
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// SCAN cursor [MATCH pattern] [COUNT count]
// Incrementally iterate the keys space

// SORT key [BY pattern] [LIMIT offset count] [GET pattern [GET pattern ...]] [ASC|DESC] [ALPHA] [STORE destination]
// Sort the elements in a list, set or sorted set
// Array reply: list of sorted elements.

// TTL key
// Get the time to live for a key
// Integer reply: TTL in seconds, or a negative value in order to signal an error.
func (cli *Client) Ttl(key interface{}) (int, error) {
	rsp, err := cli.Send("TTL", key)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// TYPE key
// Determine the type stored at key
// Simple string reply: type of key, or none when key does not exist.
func (cli *Client) Type(key interface{}) (string, error) {
	rsp, err := cli.Send("TYPE", key)
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// WAIT numslaves timeout
// Wait for the synchronous replication of all the write commands sent in the context of the current connection
// Integer reply: The command returns the number of slaves reached by all the writes performed in the context of the current connection.
func (cli *Client) Wait(numslaves, timeout int) (int, error) {
	rsp, err := cli.Send("WAIT", numslaves, timeout)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// KEYS:END

// PUBSUB:BEGIN

// PSUBSCRIBE pattern [pattern ...]
// Listen for messages published to channels matching the given patterns
func (cli *Client) PSubscribe(pattern interface{}, patterns ...interface{}) error {
	return nil
}

// PUBSUB subcommand [argument [argument ...]]
// Inspect the state of the Pub/Sub subsystem

// PUBSUB CHANNELS [pattern]
// Lists the currently active channels.
// Array reply: a list of active channels, optionally matching the specified pattern.

// PUBSUB NUMSUB [channel-1 ... channel-N]
// Returns the number of subscribers (not counting clients subscribed to patterns) for the specified channels.
// Array reply:
//     a list of channels and number of subscribers for every channel.
//     The format is channel, count, channel, count, ..., so the list is flat.
//     The order in which the channels are listed is the same as the order of the channels specified in the command call.
// Note that it is valid to call this command without channels. In this case it will just return an empty list.

// PUBSUB NUMPAT
// Returns the number of subscriptions to patterns (that are performed using the PSUBSCRIBE command).
// Note that this is not just the count of clients subscribed to patterns but the total number of patterns all the clients are subscribed to.
// Integer reply: the number of patterns all the clients are subscribed to.

// PUBLISH channel message
// Post a message to a channel
// Integer reply: the number of clients that received the message.
func (cli *Client) Publish(channel, message interface{}) (int, error) {
	rsp, err := cli.Send("PUBLISH", channel, message)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// PUNSUBSCRIBE [pattern [pattern ...]]
// Stop listening for messages posted to channels matching the given patterns
func (cli *Client) PUnsubscribe(pattern interface{}, patterns ...interface{}) error {
	return nil
}

// SUBSCRIBE channel [channel ...]
// Listen for messages published to the given channels
func (cli *Client) Subscribe(channel interface{}, channels ...interface{}) error {
	return nil
}

// UNSUBSCRIBE [channel [channel ...]]
// Stop listening for messages posted to the given channels
func (cli *Client) Unsubscribe(channel interface{}, channels ...interface{}) error {
	return nil
}

// PUBSUB:BEGIN

// STRINGS:BEGIN

// APPEND key value
// Append a value to a key
// Integer reply: the length of the string after the append operation.
func (cli *Client) Append(key, value interface{}) (int, error) {
	rsp, err := cli.Send("APPEND", key, value)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// BITCOUNT key [start end]
// Count set bits in a string
// Integer reply: The number of bits set to 1.
func (cli *Client) BitCount(key interface{}, p ...int) (int, error) {
	args := make([]interface{}, 1+len(p))
	args[0] = key
	rsp, err := cli.Send("BITCOUNT", args...)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// BITOP operation destkey key [key ...]
// Perform bitwise operations between strings
// The BITOP command supports four bitwise operations: AND, OR, XOR and NOT.
// Integer reply: The size of the string stored in the destination key, that is equal to the size of the longest input string.
func (cli *Client) BitOp(operation, destkey, key interface{}, keys ...interface{}) (int, error) {
	args := constructParameters(keys, operation, destkey, key)
	rsp, err := cli.Send("BITOP", args...)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// BITPOS key bit [start] [end]
// Find first bit set or clear in a string
// Integer reply: The command returns the position of the first bit set to 1 or 0 according to the request.
func (cli *Client) BitPOs(key interface{}, bit int, p ...int) (int, error) {
	args := make([]interface{}, 2+len(p))
	args[0] = key
	args[1] = bit
	for i := 0; i < len(p); i++ {
		args[i+1] = p[i]
	}
	rsp, err := cli.Send("BITPOS", args...)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// DECR key
// Decrement the integer value of a key by one
// Integer reply: the value of key after the decrement
func (cli *Client) Decr(key interface{}) (int, error) {
	return cli.DecrBy(key, 1)
}

// DECRBY key decrement
// Decrement the integer value of a key by the given number
// Integer reply: the value of key after the decrement
func (cli *Client) DecrBy(key interface{}, decrement int) (int, error) {
	rsp, err := cli.Send("DECRBY", key, decrement)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// GET key
// Get the value of a key
// Bulk string reply: the value of key, or nil when key does not exist.
func (cli *Client) Get(key interface{}) (value interface{}, err error) {
	rsp, err := cli.Send("GET", key)
	if err != nil {
		return nil, err
	}
	v, e := parseStringEx(rsp)
	return v, e
}

// GETBIT key offset
// Returns the bit value at offset in the string value stored at key
// Integer reply: the bit value stored at offset.
func (cli *Client) GetBit(key interface{}, offset int) (int, error) {
	rsp, err := cli.Send("GETBIT", key, offset)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// GETRANGE key start end
// Get a substring of the string stored at a key
// Bulk string reply: the substring of the string stored
func (cli *Client) GetRange(key interface{}, start, end int) (string, error) {
	rsp, err := cli.Send("GETRANGE", key, start, end)
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// GETSET key value
// Set the string value of a key and return its old value
// Bulk string reply: the old value stored at key, or nil when key did not exist.
func (cli *Client) GetSet(key, value interface{}) (interface{}, error) {
	rsp, err := cli.Send("GETSET", key, value)
	if err != nil {
		return nil, err
	}
	v, e := parseStringEx(rsp)
	return v, e
}

// INCR key
// Increment the integer value of a key by one
// Integer reply: the value of key after the increment
func (cli *Client) Incr(key interface{}) (int, error) {
	return cli.IncrBy(key, 1)
}

// INCRBY key increment
// Increment the integer value of a key by the given amount
// Integer reply: the value of key after the increment
func (cli *Client) IncrBy(key interface{}, decrement int) (int, error) {
	rsp, err := cli.Send("INCRBY", key, decrement)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// INCRBYFLOAT key increment
// Increment the float value of a key by the given amount
// Bulk string reply: the value of key after the increment.
func (cli *Client) IncrByFloat(key interface{}, decrement float64) (float64, error) {
	rsp, err := cli.Send("INCRBYFLOAT", key, decrement)
	if err != nil {
		return -1.0, err
	}
	v, e := parseFloat(rsp)
	return v, e
}

// MGET key [key ...]
// Get the values of all the given keys
// Array reply: list of values at the specified keys.
func (cli *Client) MGet(key interface{}, keys ...interface{}) ([]interface{}, error) {
	args := constructParameters(keys, key)
	rsp, err := cli.Send("MGET", args...)
	if err != nil {
		return nil, err
	}
	v, e := rsp.([]interface{})
	if e {
		return v, nil
	} else {
		return nil, errors.New("redis.MGet: not a valid reply.")
	}
}

// MSET key value [key value ...]
// Set multiple keys to multiple values
// Simple string reply: always OK since MSET can't fail.
func (cli *Client) MSet(key, value interface{}, p ...interface{}) (string, error) {
	args := constructParameters(p, key, value)
	rsp, err := cli.Send("MSET", args...)
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// MSETNX key value [key value ...]
// Set multiple keys to multiple values, only if none of the keys exist
// Integer reply, specifically:
//    1 if the all the keys were set.
//    0 if no key was set (at least one key already existed).
func (cli *Client) MSetNx(key, value interface{}, p ...interface{}) (int, error) {
	args := constructParameters(p, key, value)
	rsp, err := cli.Send("MSETNX", args...)
	if err != nil {
		return 0, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// PSETEX key milliseconds value
// Set the value and expiration in milliseconds of a key
// Simple string reply: OK if SET was executed correctly.
func (cli *Client) PSetEx(key interface{}, milliseconds int, value interface{}) (string, error) {
	rsp, err := cli.Send("PSETEX", key, milliseconds, value)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// SET key value [EX seconds] [PX milliseconds] [NX|XX]
// Set the string value of a key
// Simple string reply: OK if SET was executed correctly.
// Null reply: a Null Bulk Reply is returned if the SET operation was not performed because the user specified the NX or XX option but the condition was not met.
func (cli *Client) Set(key, value interface{}, p ...interface{}) (string, error) {
	args := make([]interface{}, 2+len(p))
	args[0] = key
	args[1] = value
	for i := 0; i < len(p); i++ {
		args[i+2] = p[i]
	}
	rsp, e := cli.Send("SET", args...)
	if e != nil {
		return "", e
	}
	v, e := parseString(rsp)
	return v, e
}

// SETBIT key offset value
// Sets or clears the bit at offset in the string value stored at key
// Integer reply: the original bit value stored at offset.
func (cli *Client) SetBit(key interface{}, offset, value int) (int, error) {
	rsp, err := cli.Send("SETBIT", key, offset, value)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// SETEX key seconds value
// Set the value and expiration of a key
// Simple string reply: OK if SET was executed correctly.
func (cli *Client) SetEx(key interface{}, seconds int, value interface{}) (string, error) {
	rsp, err := cli.Send("SETEX", key, seconds, value)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// SETNX key value
// Set the value of a key, only if the key does not exist
// Integer reply, specifically:
//    1 if the key was set
//    0 if the key was not set
func (cli *Client) SetNx(key, value interface{}) (int, error) {
	rsp, err := cli.Send("SETNX", key, value)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// SETRANGE key offset value
// Overwrite part of a string at key starting at the specified offset
// Integer reply: the length of the string after it was modified by the command.
func (cli *Client) SetRange(key interface{}, offset int, value interface{}) (int, error) {
	rsp, err := cli.Send("SETRANGE", key, offset, value)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// STRLEN key
// Get the length of the value stored in a key
// Integer reply: the length of the string at key, or 0 when key does not exist.
func (cli *Client) StrLen(key interface{}) (int, error) {
	rsp, err := cli.Send("StrLen", key)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// STRINGS:END
