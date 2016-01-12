// The MIT License (MIT)
//
// Copyright (c) 2016 Roy Xu

package redis

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strconv"
)

type Client interface {
	respStringer
	Send(commandName string, args ...interface{}) (reply interface{}, err error)
	//Receive() (resp interface{}, err error)
	//RawReceive() (resp []byte, err error)
	Close() (err error)
}

type client struct {
	conn net.Conn
	bw   *bufio.Writer
	br   *bufio.Reader
}

func (c *client) Close() (err error) {
	return c.conn.Close()
}

func (c *client) Send(commandName string, args ...interface{}) (reply interface{}, err error) {
	return c.executeCommand(commandName, args...)
}

// receive returns the raw result without the symbols(+-:$*) and CR&LF.
// One-dimensional byte slice, []byte:
//    For Simple Strings the first byte of the reply is "+"
//    For Errors the first byte of the reply is "-"
//    For Integers the first byte of the reply is ":"
//    For Bulk Strings the first byte of the reply is "$"
// Two-dimensional byte slice, [][]byte:
//    For Arrays the first byte of the reply is "*"
// nil :
//    For Null Bulk String
//    For Null Array
func (c *client) receive() (interface{}, error) {
	line, err := c.br.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	symbol, resp := line[0], line[1:len(line)-2] // trim first bit and CRLF
	switch symbol {
	case '+', '-', ':':
		return resp, nil
	case '$':
		length, _ := strconv.Atoi(string(resp))
		if length == -1 {
			return nil, nil
		} else {
			s, _ := c.br.ReadSlice('\n')
			return s[0 : len(s)-2], nil
		}
	case '*':
		length, _ := strconv.Atoi(string(resp))
		if length == -1 {
			return nil, nil
		}
		reslt := make([]interface{}, 0)
		for i := 0; i < length; i++ {
			rep, _ := c.receive()
			reslt = append(reslt, rep)
		}
		return reslt, nil
	default:
		panic(fmt.Sprintf("Protocol Error: %s, %s", string(symbol), string(resp)))
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

	switch p.(type) {
	case string:
		err = c.writeString(p.(string))
	case []byte:
		err = c.writeBytes(p.([]byte))
	case int:
		err = c.writeInt(int64(p.(int)))
	case float32:
		err = c.writeFloat(float64(p.(float32)))
	case float64:
		err = c.writeFloat(p.(float64))
	case bool:
		if p.(bool) {
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

// [BEGIN] RESP Strings

// APPEND key value
// Append a value to a key
// Integer reply: the length of the string after the append operation.
func (c *client) Append(key, value interface{}) (int, error) {
	resp, err := c.executeCommand("APPEND", key, value)
	if err != nil {
		return -1, err
	}
	v := resp.(int)
	return v, nil
}

// BITCOUNT key [start end]
// Count set bits in a string
// Integer reply: The number of bits set to 1.
func (c *client) BitCount(key interface{}, p ...int) (int, error) {
	args := make([]interface{}, 1+len(p))
	args[0] = key
	resp, err := c.executeCommand("BITCOUNT", args...)
	if err != nil {
		return -1, err
	}
	v := resp.(int)
	return v, nil
}

// BITOP operation destkey key [key ...]
// Perform bitwise operations between strings
// The BITOP command supports four bitwise operations: AND, OR, XOR and NOT.
// Integer reply: The size of the string stored in the destination key, that is equal to the size of the longest input string.
func (c *client) BitOp(operation, destkey, key interface{}, keys ...interface{}) (int, error) {
	args := make([]interface{}, 3+len(keys))
	args[0] = operation
	args[1] = destkey
	args[2] = key
	for i := 0; i < len(keys); i++ {
		args[i+1] = keys[i]
	}
	resp, err := c.executeCommand("BITOP", args...)
	if err != nil {
		return -1, err
	}
	return resp.(int), err
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
	resp, err := c.executeCommand("BITPOS", args...)
	if err != nil {
		return -1, err
	}
	return resp.(int), err
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
	resp, err := c.executeCommand("DECRBY", key, decrement)
	if err != nil {
		return -1, err
	}
	v := resp.(int)
	return v, nil
}

// GET key
// Get the value of a key
// Bulk string reply: the value of key, or nil when key does not exist.
func (c *client) Get(key interface{}) (value interface{}, err error) {
	resp, err := c.executeCommand("GET", key)
	if err != nil {
		return "", err
	}
	v := resp.(string)
	return v, nil
}

// GETBIT key offset
// Returns the bit value at offset in the string value stored at key
// Integer reply: the bit value stored at offset.
func (c *client) GetBit(key interface{}, offset int) (int, error) {
	resp, err := c.executeCommand("GETBIT", key, offset)
	if err != nil {
		return -1, err
	}
	v := resp.(int)
	return v, nil
}

// GETRANGE key start end
// Get a substring of the string stored at a key
// Bulk string reply: the substring of the string stored
func (c *client) GetRange(key interface{}, start, end int) (interface{}, error) {
	resp, err := c.executeCommand("GETRANGE", key, start, end)
	if err != nil {
		return "", err
	}
	v := resp.(string)
	return v, nil
}

// GETSET key value
// Set the string value of a key and return its old value
// Bulk string reply: the old value stored at key, or nil when key did not exist.
func (c *client) GetSet(key, value interface{}) (interface{}, error) {
	resp, err := c.executeCommand("GETSET", key, value)
	if err != nil {
		return "", err
	}
	v := resp.(string)
	return v, nil
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
	resp, err := c.executeCommand("INCRBY", key, decrement)
	if err != nil {
		return -1, err
	}
	v := resp.(int)
	return v, nil
}

// INCRBYFLOAT key increment
// Increment the float value of a key by the given amount
// Bulk string reply: the value of key after the increment.
func (c *client) IncrByFloat(key interface{}, decrement float64) (float64, error) {
	resp, err := c.executeCommand("INCRBYFLOAT", key, decrement)
	if err != nil {
		return -1.0, err
	}
	v, e := strconv.ParseFloat(resp.(string), 10)
	if e != nil {
		return -1.0, e
	}
	return v, nil
}

// MGET key [key ...]
// Get the values of all the given keys
// Array reply: list of values at the specified keys.
func (c *client) MGet(key interface{}, keys ...interface{}) ([]interface{}, error) {
	// args := make([]interface{}, 1+len(keys))
	// args[0] = key
	// for i := 0; i < len(keys); i++ {
	// 	args[i+1] = keys[i]
	// }
	// return c.executeCommand("MGET", args...)
	return nil, nil
}

// MSET key value [key value ...]
// Set multiple keys to multiple values
// Simple string reply: always OK since MSET can't fail.
func (c *client) MSet(key, value interface{}, p ...interface{}) (string, error) {
	args := make([]interface{}, 2+len(p))
	args[0] = key
	args[1] = value
	for i := 0; i < len(p); i++ {
		args[i+2] = p[i]
	}
	resp, err := c.executeCommand("MSET", args...)
	if err != nil {
		return "", err
	}
	return resp.(string), nil
}

// MSETNX key value [key value ...]
// Set multiple keys to multiple values, only if none of the keys exist
// Integer reply, specifically:
//    1 if the all the keys were set.
//    0 if no key was set (at least one key already existed).
func (c *client) MSetNx(key, value interface{}, p ...interface{}) (int, error) {
	args := make([]interface{}, 2+len(p))
	args[0] = key
	args[1] = value
	for i := 0; i < len(p); i++ {
		args[i+2] = p[i]
	}
	resp, err := c.executeCommand("MSETNX", args...)
	if err != nil {
		return 0, err
	}
	return resp.(int), nil
}

// PSETEX key milliseconds value
// Set the value and expiration in milliseconds of a key
// Simple string reply
func (c *client) PSetEx(key interface{}, milliseconds int, value interface{}) (string, error) {
	resp, err := c.executeCommand("PSETEX", key, milliseconds, value)
	if err != nil {
		return "", err
	}
	v := resp.(string)
	return v, nil
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
	s, e := c.executeCommand("SET", args...)
	if e != nil {
		return "", e
	}
	switch s.(type) {
	case nil:
		return "", nil
	case string:
		v := s.(string)
		return v, nil
	default:
		panic("Set do not work properly with an unexpected value.")
	}
}

// SETBIT key offset value
// Sets or clears the bit at offset in the string value stored at key
// Integer reply: the original bit value stored at offset.
func (c *client) SetBit(key interface{}, offset, value int) (int, error) {
	resp, err := c.executeCommand("SETBIT", key, offset, value)
	if err != nil {
		return -1, err
	}
	v := resp.(int)
	return v, nil
}

// SETEX key seconds value
// Set the value and expiration of a key
// Simple string reply
func (c *client) SetEx(key interface{}, seconds int, value interface{}) (string, error) {
	resp, err := c.executeCommand("SETEX", key, seconds, value)
	if err != nil {
		return "", err
	}
	v := resp.(string)
	return v, nil
}

// SETNX key value
// Set the value of a key, only if the key does not exist
// Integer reply, specifically:
//    1 if the key was set
//    0 if the key was not set
func (c *client) SetNx(key, value interface{}) (int, error) {
	resp, err := c.executeCommand("SETNX", key, value)
	if err != nil {
		return 0, err
	}
	v := resp.(int)
	return v, nil
}

// SETRANGE key offset value
// Overwrite part of a string at key starting at the specified offset
// Integer reply: the length of the string after it was modified by the command.
func (c *client) SetRange(key interface{}, offset, value int) (int, error) {
	resp, err := c.executeCommand("SETRANGE", key, offset, value)
	if err != nil {
		return 0, err
	}
	v := resp.(int)
	return v, nil
}

// STRLEN key
// Get the length of the value stored in a key
// Integer reply: the length of the string at key, or 0 when key does not exist.
func (c *client) StrLen(key interface{}) (int, error) {
	resp, err := c.executeCommand("StrLen", key)
	if err != nil {
		return 0, err
	}
	v := resp.(int)
	return v, nil
}

// [END] RESP Strings
