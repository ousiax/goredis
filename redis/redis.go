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

func (c *client) receive() (interface{}, error) {
	line, err := c.br.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	symbol, resp := line[0], line[1:len(line)-2] // trim first bit and CRLF
	switch symbol {
	case '+', '-':
		return string(resp), nil
	case ':':
		return strconv.Atoi(string(resp))
	case '$':
		length, _ := strconv.Atoi(string(resp))
		if length == -1 {
			return nil, nil
		} else {
			s, _ := c.br.ReadSlice('\n')
			return string(s[0 : len(s)-2]), nil
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

// APPEND key value Append a value to a key
func (c *client) Append(key, value string) (int, error) {
	resp, err := c.executeCommand("APPEND", key, value)
	if err != nil {
		return -1, err
	}
	v := resp.(int)
	return v, nil
}

// BITCOUNT key [start end] Count set bits in a string
func (c *client) BitCount(key string, p ...int) (int, error) {
	args := make([]interface{}, 1+len(p))
	args[0] = key
	resp, err := c.executeCommand("BITCOUNT", args...)
	if err != nil {
		return -1, err
	}
	v := resp.(int)
	return v, nil
}

// BITOP operation destkey key [key ...] Perform bitwise operations between strings
func (c *client) BitOp(operation, destkey, key string, keys ...string) (int, error) {
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

// BITPOS key bit [start] [end] Find first bit set or clear in a string
func (c *client) BitPOs(key string, bit int, p ...int) (int, error) {
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

// DECR key Decrement the integer value of a key by one
func (c *client) Decr(key string) (int, error) {
	return c.DecrBy(key, 1)
}

// DECRBY key decrement Decrement the integer value of a key by the given number
func (c *client) DecrBy(key string, decrement int) (int, error) {
	resp, err := c.executeCommand("DECRBY", key, decrement)
	if err != nil {
		return -1, err
	}
	v := resp.(int)
	return v, nil
}

// GET key Get the value of a key
func (c *client) Get(key string) (value string, err error) {
	resp, err := c.executeCommand("GET", key)
	if err != nil {
		return "", err
	}
	v := resp.(string)
	return v, nil
}

// GETBIT key offset Returns the bit value at offset in the string value stored at key
func (c *client) GetBit(key string, offset int) (int, error) {
	resp, err := c.executeCommand("GETBIT", key, offset)
	if err != nil {
		return -1, err
	}
	v := resp.(int)
	return v, nil
}

// GETRANGE key start end Get a substring of the string stored at a key
func (c *client) GetRange(key string, start, end int) (string, error) {
	resp, err := c.executeCommand("GETRANGE", key, start, end)
	if err != nil {
		return "", err
	}
	v := resp.(string)
	return v, nil
}

// GETSET key value Set the string value of a key and return its old value
func (c *client) GetSet(key, value string) (string, error) {
	resp, err := c.executeCommand("GETSET", key, value)
	if err != nil {
		return "", err
	}
	v := resp.(string)
	return v, nil
}

// INCR key Increment the integer value of a key by one
func (c *client) Incr(key string) (int, error) {
	return c.IncrBy(key, 1)
}

// INCRBY key increment Increment the integer value of a key by the given amount
func (c *client) IncrBy(key string, decrement int) (int, error) {
	resp, err := c.executeCommand("INCRBY", key, decrement)
	if err != nil {
		return -1, err
	}
	v := resp.(int)
	return v, nil
}

// INCRBYFLOAT key increment Increment the float value of a key by the given amount
func (c *client) IncrByFloat(key string, decrement float64) (float64, error) {
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

// MGET key [key ...] Get the values of all the given keys
func (c *client) MGet(key string, keys ...string) ([]interface{}, error) {
	// args := make([]interface{}, 1+len(keys))
	// args[0] = key
	// for i := 0; i < len(keys); i++ {
	// 	args[i+1] = keys[i]
	// }
	// return c.executeCommand("MGET", args...)
	return nil, nil
}

// MSET key value [key value ...] Set multiple keys to multiple values
func (c *client) MSet(key, value string, p ...string) (string, error) {
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

// MSETNX key value [key value ...] Set multiple keys to multiple values, only if none of the keys exist
func (c *client) MSetNx(key, value string, p ...string) (string, error) {
	args := make([]interface{}, 2+len(p))
	args[0] = key
	args[1] = value
	for i := 0; i < len(p); i++ {
		args[i+2] = p[i]
	}
	resp, err := c.executeCommand("MSETNX", args...)
	if err != nil {
		return "", err
	}
	return resp.(string), nil
}

// PSETEX key milliseconds value Set the value and expiration in milliseconds of a key
func (c *client) PSetEx(key string, milliseconds int, value string) (string, error) {
	resp, err := c.executeCommand("PSETEX", key, milliseconds, value)
	if err != nil {
		return "", err
	}
	v := resp.(string)
	return v, nil
}

// SET key value [EX seconds] [PX milliseconds] [NX|XX] Set the string value of a key
func (c *client) Set(key, value string, p ...interface{}) (string, error) {
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

// SETBIT key offset value Sets or clears the bit at offset in the string value stored at key
func (c *client) SetBit(key string, offset, value int) (int, error) {
	resp, err := c.executeCommand("SETBIT", key, offset, value)
	if err != nil {
		return -1, err
	}
	v := resp.(int)
	return v, nil
}

// SETEX key seconds value Set the value and expiration of a key
func (c *client) SetEx(key string, seconds int, value string) (string, error) {
	resp, err := c.executeCommand("SETEX", key, seconds, value)
	if err != nil {
		return "", err
	}
	v := resp.(string)
	return v, nil
}

// SETNX key value Set the value of a key, only if the key does not exist
func (c *client) SetNx(key, value string) (int, error) {
	resp, err := c.executeCommand("SETNX", key, value)
	if err != nil {
		return 0, err
	}
	v := resp.(int)
	return v, nil
}

// SETRANGE key offset value Overwrite part of a string at key starting at the specified offset
func (c *client) SetRange(key string, offset, value int) (int, error) {
	resp, err := c.executeCommand("SETRANGE", key, offset, value)
	if err != nil {
		return 0, err
	}
	v := resp.(int)
	return v, nil
}

// STRLEN key Get the length of the value stored in a key
func (c *client) StrLen(key string) (int, error) {
	resp, err := c.executeCommand("StrLen", key)
	if err != nil {
		return 0, err
	}
	v := resp.(int)
	return v, nil
}
