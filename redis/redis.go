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
	case '+', '-':
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

// parseStringEx returns a nil or string, otherwise a nil when a error occured.
func parseStringEx(p interface{}) (interface{}, error) {
	switch v := p.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	case nil:
		return nil, nil
	default:
		return nil, strconv.ErrRange
	}
}

// parseString return a empty string if p is nil or a error occured, otherwise a string.
// usually, the p is a string or a nil (i.e. a zero value).
func parseString(p interface{}) string {
	s, _ := p.(string)
	return s
}

// [BEGIN] RESP Strings

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
	args := make([]interface{}, 3+len(keys))
	args[0] = operation
	args[1] = destkey
	args[2] = key
	for i := 0; i < len(keys); i++ {
		args[i+1] = keys[i]
	}
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
		return "", err
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
	v := parseString(rsp)
	return v, nil
}

// GETSET key value
// Set the string value of a key and return its old value
// Bulk string reply: the old value stored at key, or nil when key did not exist.
func (c *client) GetSet(key, value interface{}) (interface{}, error) {
	rsp, err := c.executeCommand("GETSET", key, value)
	if err != nil {
		return "", err
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
	rsp, err := c.executeCommand("MSET", args...)
	if err != nil {
		return "", err
	}
	v := parseString(rsp)
	return v, nil
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
	rsp, err := c.executeCommand("MSETNX", args...)
	if err != nil {
		return 0, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// PSETEX key milliseconds value
// Set the value and expiration in milliseconds of a key
// Simple string reply
func (c *client) PSetEx(key interface{}, milliseconds int, value interface{}) (string, error) {
	rsp, err := c.executeCommand("PSETEX", key, milliseconds, value)
	if err != nil {
		return "", err
	}
	v := parseString(rsp)
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
	rsp, e := c.executeCommand("SET", args...)
	if e != nil {
		return "", e
	}
	v := parseString(rsp)
	return v, nil
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
// Simple string reply
func (c *client) SetEx(key interface{}, seconds int, value interface{}) (string, error) {
	rsp, err := c.executeCommand("SETEX", key, seconds, value)
	if err != nil {
		return "", err
	}
	v := parseString(rsp)
	return v, nil
}

// SETNX key value
// Set the value of a key, only if the key does not exist
// Integer reply, specifically:
//    1 if the key was set
//    0 if the key was not set
func (c *client) SetNx(key, value interface{}) (int, error) {
	rsp, err := c.executeCommand("SETNX", key, value)
	if err != nil {
		return 0, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// SETRANGE key offset value
// Overwrite part of a string at key starting at the specified offset
// Integer reply: the length of the string after it was modified by the command.
func (c *client) SetRange(key interface{}, offset, value int) (int, error) {
	rsp, err := c.executeCommand("SETRANGE", key, offset, value)
	if err != nil {
		return 0, err
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
		return 0, err
	}
	v, e := parseInt(rsp)
	return v, e
}

// [END] RESP Strings
