// The MIT License (MIT)
//
// Copyright (c) 2016 Roy Xu

package redis

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
)

type Client interface {
	respCluster
	respConnection
	respKeys
	respLists
	respServer
	respStrings
	respTransaction
	respPubSub
	respHyperLogLog
	Send(commandName string, args ...interface{}) (reply interface{}, err error)
	Close() (err error)
}

type client struct {
	conn net.Conn
	bw   *bufio.Writer
	br   *bufio.Reader
}

func NewClient(network, address string) (Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)
	cli := client{conn: conn, bw: w, br: r}
	return Client(&cli), nil
}

func (c *client) Close() (err error) {
	return c.conn.Close()
}

func (c *client) Send(commandName string, args ...interface{}) (reply interface{}, err error) {
	rsp, err := c.executeCommand(commandName, args...)
	return rsp, err
}

// receive returns the raw result without the symbols(+-:$*) and CR&LF.
// string
//    For Simple Strings the first byte of the reply is "+"
// error
//    For Errors the first byte of the reply is "-"
// int
//    For Integers the first byte of the reply is ":"
// []byte
//    For Bulk Strings the first byte of the reply is "$"
// []interface{}
//    For Arrays the first byte of the reply is "*"
// nil
//    For Null Bulk String
//    For Null Array
func (c *client) receive() (interface{}, error) {
	line, err := c.br.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	symbol, rsp := line[0], line[1:len(line)-2] // trim first bit and CRLF
	switch symbol {
	case '+':
		return string(rsp), nil
	case ':':
		return strconv.Atoi(string(rsp))
	case '$':
		length, _ := strconv.Atoi(string(rsp))
		if length == -1 {
			return nil, nil
		} else {
			s, _ := c.br.ReadSlice('\n')
			return s[0 : len(s)-2], nil
		}
	case '*':
		length, _ := strconv.Atoi(string(rsp))
		if length == -1 {
			return nil, nil
		}
		reslt := make([]interface{}, 0, length)
		for i := 0; i < length; i++ {
			rep, _ := c.receive()
			reslt = append(reslt, rep)
		}
		return reslt, nil
	case '-':
		return nil, errors.New(string(rsp))
	default:
		panic(fmt.Sprintf("Protocol Error: %s, %s", string(symbol), string(rsp)))
	}
}

func (c *client) executeCommand(commandName string, args ...interface{}) (reply interface{}, err error) {
	err = c.writeHeader('*', int64(1+len(args)))
	if err != nil {
		return nil, err
	}
	err = c.writeString(commandName)
	for _, arg := range args {
		if err != nil {
			break
		}
		err = c.writeInterface(arg)
	}
	err = c.bw.Flush()
	if err == nil {
		return c.receive()
	}
	return nil, err
}

func (c *client) writeCRLF() error {
	_, err := c.bw.Write([]byte{'\r', '\n'})
	return err
}

// writeBytes write an RESP head header(*$) to buffer and appends with CR&LF. e.g. $3\r\n, *3\r\n
func (c *client) writeHeader(symbol byte, n int64) error {
	c.bw.WriteByte(symbol)
	// err = c.writeInt(n) // bugged, recursive calling.
	p := strconv.AppendInt([]byte{}, n, 10)
	c.bw.Write(p)
	err := c.writeCRLF()
	return err
}

// writeBytes writes a []byte to buffer with RESP header($) and appends with CR&LF. e.g. $3\r\nSET\r\n
func (c *client) writeBytes(p []byte) error {
	c.writeHeader('$', int64(len(p)))
	c.bw.Write(p)
	err := c.writeCRLF()
	return err
}

// writeInt writes an integer to buffer with RESP header($) and and appends with CR&LF. e.g. $3\r\n125\r\n
func (c *client) writeInt(n int64) error {
	p := strconv.AppendInt([]byte{}, n, 10)
	return c.writeBytes(p)
}

// writeFloat writes a float to buffer with a RESP header($) and and appends with CR&LF. e.g. $3\r\n1.5\r\n
func (c *client) writeFloat(f float64) error {
	p := strconv.AppendFloat([]byte{}, f, 'g', -1, 64)
	return c.writeBytes(p)
}

// writeString writes a string to buffer with RESP header($) and and appends with CR&LF. e.g. $4\r\nLLEN\r\n
func (c *client) writeString(s string) error {
	err := c.writeBytes([]byte(s))
	return err
}

// writeInterface parses parameter 'p' and writes a it to buffer with RESP header($) and and appends with CR&LF. e.g. $4\r\nLLEN\r\n
func (c *client) writeInterface(p interface{}) (err error) {
	switch v := p.(type) {
	case string:
		err = c.writeString(v)
	case []byte:
		err = c.writeBytes(v)
	case int:
		err = c.writeInt(int64(v))
	case float32:
		err = c.writeFloat(float64(v))
	case float64:
		err = c.writeFloat(v)
	case bool:
		if v {
			err = c.writeString("1")
		} else {
			err = c.writeString("0")
		}
	case nil:
		err = c.writeString("")
	default:
		var buf bytes.Buffer
		fmt.Fprint(&buf, p)
		err = c.writeBytes(buf.Bytes())
	}
	return err
}

func parseInt(p interface{}) (int, error) {
	v, e := p.(int)
	if e {
		return v, nil
	}
	return v, strconv.ErrRange
}

// parseFloat parses a RESP Bulk String to a float64 number.
func parseFloat(p interface{}) (float64, error) {
	switch v := p.(type) {
	case []byte:
		return strconv.ParseFloat(string(v), 64)
	default:
		return 0.0, strconv.ErrRange
	}
}

// parseStringEx parses a RESP Bulk String or a Simple String to a string or a nil, otherwise a nil when a error occured.
func parseStringEx(p interface{}) (interface{}, error) {
	switch v := p.(type) {
	case []byte:
		return string(v), nil
	case nil:
		return nil, nil
	case string:
		return v, nil
	default:
		return nil, errors.New(fmt.Sprintf("redis.parseStringEx: interface conversion, interface is %T, not string.", p))
	}
}

// parseString returns a string if p is a string type, otherwise a empty string.
// usually, the p is a string or a nil (i.e. a zero value).
func parseString(p interface{}) (string, error) {
	rsp, err := parseStringEx(p)
	s, _ := rsp.(string)
	return s, err
}

// parseStrings parses a RESP reply (array apply) to a string arrray (may contain null value).
func parseStrings(p interface{}) ([]interface{}, error) {
	rsp, e := p.([]interface{})
	if !e {
		return nil, errors.New(fmt.Sprintf("redis.parseStrings: interface conversion, interface is %T, not []interface{}.", p))
	}
	for i, v := range rsp {
		p, _ := parseStringEx(v) // Error check should be not required for performance.
		rsp[i] = p
	}
	return rsp, nil
}

func constructParameters(opt []interface{}, p ...interface{}) []interface{} {
	l := len(opt) + len(p)
	a := make([]interface{}, 0, l)
	for _, v := range p {
		a = append(a, v)
	}
	for _, v := range opt {
		a = append(a, v)
	}
	return a
}

// CONNECTION:BEGIN

// AUTH password
// Authenticate to the server
// Simple string reply
func (c *client) Auth(password string) (string, error) {
	rsp, err := c.executeCommand("AUTH", password)
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// ECHO message
// Echo the given string
// Simple string reply
func (c *client) Echo(message string) (string, error) {
	rsp, err := c.executeCommand("ECHO", message)
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// PING
// Ping the server
// Simple string reply
func (c *client) Ping() (string, error) {
	rsp, err := c.executeCommand("PING")
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// QUIT
// Close the connection
// Simple string reply: always OK.
func (c *client) Quit() (string, error) {
	rsp, err := c.executeCommand("QUIT")
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// SELECT index
// Change the selected database for the current connection
// Simple string reply
func (c *client) Select(index int) (string, error) {
	rsp, err := c.executeCommand("SELECT", index)
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
func (c *client) Del(key interface{}, keys ...interface{}) (int, error) {
	p := constructParameters(keys, key)
	rsp, err := c.executeCommand("DEL", p...)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// DUMP key
// Return a serialized version of the value stored at the specified key.
// Bulk string reply: the serialized value.
func (c *client) Dump(key interface{}) (interface{}, error) {
	rsp, err := c.executeCommand("DUMP", key)
	if err != nil {
		return nil, err
	}
	v, e := parseStringEx(rsp)
	return v, e
}

// EXISTS key [key ...]
// Determine if a key exists
// Integer reply: The number of keys existing among the ones specified as arguments.
func (c *client) Exists(key interface{}, keys ...interface{}) (int, error) {
	p := constructParameters(keys, key)
	rsp, err := c.executeCommand("EXISTS", p...)
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
func (c *client) Expire(key interface{}, seconds int) (int, error) {
	rsp, err := c.executeCommand("EXPIRE", key, seconds)
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
func (c *client) ExpireAt(key interface{}, timestamp int) (int, error) {
	rsp, err := c.executeCommand("EXPIREAT", key, timestamp)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// KEYS pattern
// Find all keys matching the given pattern
// Array reply: list of keys matching pattern.
func (c *client) Keys(pattern interface{}) ([]interface{}, error) {
	rsp, err := c.executeCommand("KEYS", pattern)
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
func (c *client) Move(key interface{}, db int) (int, error) {
	rsp, err := c.executeCommand("MOVE", key, db)
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
func (c *client) Persist(key interface{}) (int, error) {
	rsp, err := c.executeCommand("PERSIST", key)
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
func (c *client) PExpire(key string, milliseconds int) (int, error) {
	rsp, err := c.executeCommand("PEXPIRE", key, milliseconds)
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
func (c *client) PExpireAt(key string, millisecondTimestamp int) (int, error) {
	rsp, err := c.executeCommand("PEXPIREAT", key, millisecondTimestamp)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// PTTL key
// Get the time to live for a key in milliseconds
// Integer reply: TTL in milliseconds, or a negative value in order to signal an error.
func (c *client) Pttl(key string) (int, error) {
	rsp, err := c.executeCommand("PTTL", key)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// RANDOMKEY
// Return a random key from the keyspace
// Bulk string reply: the random key, or nil when the database is empty.
func (c *client) RandomKey() (interface{}, error) {
	rsp, err := c.executeCommand("RANDOMKEY")
	if err != nil {
		return nil, err
	}
	v, e := parseStringEx(rsp)
	return v, e
}

// RENAME key newkey
// Rename a key
// Simple string reply
func (c *client) Rename(key, newkey string) (string, error) {
	rsp, err := c.executeCommand("RENAME", key, newkey)
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
func (c *client) RenameNx(key, newkey string) (int, error) {
	rsp, err := c.executeCommand("RENAMENX", key, newkey)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// RESTORE key ttl serialized-value [REPLACE]
// Create a key using the provided serialized value, previously obtained using DUMP.
// Simple string reply: The command returns OK on success.
func (c *client) Restore(key interface{}, ttl int, serializedValue interface{}, replace bool) (string, error) {
	var rsp interface{}
	var err error
	if replace {
		rsp, err = c.executeCommand("RESTORE", key, ttl, serializedValue, "REPLACE")
	} else {
		rsp, err = c.executeCommand("RESTORE", key, ttl, serializedValue)
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
func (c *client) Ttl(key interface{}) (int, error) {
	rsp, err := c.executeCommand("TTL", key)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// TYPE key
// Determine the type stored at key
// Simple string reply: type of key, or none when key does not exist.
func (c *client) Type(key interface{}) (string, error) {
	rsp, err := c.executeCommand("TYPE", key)
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// WAIT numslaves timeout
// Wait for the synchronous replication of all the write commands sent in the context of the current connection
// Integer reply: The command returns the number of slaves reached by all the writes performed in the context of the current connection.
func (c *client) Wait(numslaves, timeout int) (int, error) {
	rsp, err := c.executeCommand("WAIT", numslaves, timeout)
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
func (c *client) PSubscribe(pattern interface{}, patterns ...interface{}) error {
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
func (c *client) Publish(channel, message interface{}) (int, error) {
	rsp, err := c.executeCommand("PUBLISH", channel, message)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// PUNSUBSCRIBE [pattern [pattern ...]]
// Stop listening for messages posted to channels matching the given patterns
func (c *client) PUnsubscribe(pattern interface{}, patterns ...interface{}) error {
	return nil
}

// SUBSCRIBE channel [channel ...]
// Listen for messages published to the given channels
func (c *client) Subscribe(channel interface{}, channels ...interface{}) error {
	return nil
}

// UNSUBSCRIBE [channel [channel ...]]
// Stop listening for messages posted to the given channels
func (c *client) Unsubscribe(channel interface{}, channels ...interface{}) error {
	return nil
}

// PUBSUB:BEGIN

// STRINGS:BEGIN

// APPEND key value
// Append a value to a key
// Integer reply: the length of the string after the append operation.
func (c *client) Append(key, value interface{}) (int, error) {
	rsp, err := c.executeCommand("APPEND", key, value)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// BITCOUNT key [start end]
// Count set bits in a string
// Integer reply: The number of bits set to 1.
func (c *client) BitCount(key interface{}, p ...int) (int, error) {
	args := make([]interface{}, 1+len(p))
	args[0] = key
	rsp, err := c.executeCommand("BITCOUNT", args...)
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
func (c *client) BitOp(operation, destkey, key interface{}, keys ...interface{}) (int, error) {
	args := constructParameters(keys, operation, destkey, key)
	rsp, err := c.executeCommand("BITOP", args...)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// BITPOS key bit [start] [end]
// Find first bit set or clear in a string
// Integer reply: The command returns the position of the first bit set to 1 or 0 according to the request.
func (c *client) BitPOs(key interface{}, bit int, p ...int) (int, error) {
	args := make([]interface{}, 2+len(p))
	args[0] = key
	args[1] = bit
	for i := 0; i < len(p); i++ {
		args[i+1] = p[i]
	}
	rsp, err := c.executeCommand("BITPOS", args...)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// DECR key
// Decrement the integer value of a key by one
// Integer reply: the value of key after the decrement
func (c *client) Decr(key interface{}) (int, error) {
	return c.DecrBy(key, 1)
}

// DECRBY key decrement
// Decrement the integer value of a key by the given number
// Integer reply: the value of key after the decrement
func (c *client) DecrBy(key interface{}, decrement int) (int, error) {
	rsp, err := c.executeCommand("DECRBY", key, decrement)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// GET key
// Get the value of a key
// Bulk string reply: the value of key, or nil when key does not exist.
func (c *client) Get(key interface{}) (value interface{}, err error) {
	rsp, err := c.executeCommand("GET", key)
	if err != nil {
		return nil, err
	}
	v, e := parseStringEx(rsp)
	return v, e
}

// GETBIT key offset
// Returns the bit value at offset in the string value stored at key
// Integer reply: the bit value stored at offset.
func (c *client) GetBit(key interface{}, offset int) (int, error) {
	rsp, err := c.executeCommand("GETBIT", key, offset)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// GETRANGE key start end
// Get a substring of the string stored at a key
// Bulk string reply: the substring of the string stored
func (c *client) GetRange(key interface{}, start, end int) (string, error) {
	rsp, err := c.executeCommand("GETRANGE", key, start, end)
	if err != nil {
		return "", err
	}
	v, e := parseString(rsp)
	return v, e
}

// GETSET key value
// Set the string value of a key and return its old value
// Bulk string reply: the old value stored at key, or nil when key did not exist.
func (c *client) GetSet(key, value interface{}) (interface{}, error) {
	rsp, err := c.executeCommand("GETSET", key, value)
	if err != nil {
		return nil, err
	}
	v, e := parseStringEx(rsp)
	return v, e
}

// INCR key
// Increment the integer value of a key by one
// Integer reply: the value of key after the increment
func (c *client) Incr(key interface{}) (int, error) {
	return c.IncrBy(key, 1)
}

// INCRBY key increment
// Increment the integer value of a key by the given amount
// Integer reply: the value of key after the increment
func (c *client) IncrBy(key interface{}, decrement int) (int, error) {
	rsp, err := c.executeCommand("INCRBY", key, decrement)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// INCRBYFLOAT key increment
// Increment the float value of a key by the given amount
// Bulk string reply: the value of key after the increment.
func (c *client) IncrByFloat(key interface{}, decrement float64) (float64, error) {
	rsp, err := c.executeCommand("INCRBYFLOAT", key, decrement)
	if err != nil {
		return -1.0, err
	}
	v, e := parseFloat(rsp)
	return v, e
}

// MGET key [key ...]
// Get the values of all the given keys
// Array reply: list of values at the specified keys.
func (c *client) MGet(key interface{}, keys ...interface{}) ([]interface{}, error) {
	args := constructParameters(keys, key)
	rsp, err := c.executeCommand("MGET", args...)
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
func (c *client) MSet(key, value interface{}, p ...interface{}) (string, error) {
	args := constructParameters(p, key, value)
	rsp, err := c.executeCommand("MSET", args...)
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
func (c *client) MSetNx(key, value interface{}, p ...interface{}) (int, error) {
	args := constructParameters(p, key, value)
	rsp, err := c.executeCommand("MSETNX", args...)
	if err != nil {
		return 0, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// PSETEX key milliseconds value
// Set the value and expiration in milliseconds of a key
// Simple string reply: OK if SET was executed correctly.
func (c *client) PSetEx(key interface{}, milliseconds int, value interface{}) (string, error) {
	rsp, err := c.executeCommand("PSETEX", key, milliseconds, value)
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
func (c *client) Set(key, value interface{}, p ...interface{}) (string, error) {
	args := make([]interface{}, 2+len(p))
	args[0] = key
	args[1] = value
	for i := 0; i < len(p); i++ {
		args[i+2] = p[i]
	}
	rsp, e := c.executeCommand("SET", args...)
	if e != nil {
		return "", e
	}
	v, e := parseString(rsp)
	return v, e
}

// SETBIT key offset value
// Sets or clears the bit at offset in the string value stored at key
// Integer reply: the original bit value stored at offset.
func (c *client) SetBit(key interface{}, offset, value int) (int, error) {
	rsp, err := c.executeCommand("SETBIT", key, offset, value)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// SETEX key seconds value
// Set the value and expiration of a key
// Simple string reply: OK if SET was executed correctly.
func (c *client) SetEx(key interface{}, seconds int, value interface{}) (string, error) {
	rsp, err := c.executeCommand("SETEX", key, seconds, value)
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
func (c *client) SetNx(key, value interface{}) (int, error) {
	rsp, err := c.executeCommand("SETNX", key, value)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// SETRANGE key offset value
// Overwrite part of a string at key starting at the specified offset
// Integer reply: the length of the string after it was modified by the command.
func (c *client) SetRange(key interface{}, offset int, value interface{}) (int, error) {
	rsp, err := c.executeCommand("SETRANGE", key, offset, value)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// STRLEN key
// Get the length of the value stored in a key
// Integer reply: the length of the string at key, or 0 when key does not exist.
func (c *client) StrLen(key interface{}) (int, error) {
	rsp, err := c.executeCommand("StrLen", key)
	if err != nil {
		return -1, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// STRINGS:END
