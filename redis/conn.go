// The MIT License (MIT)

// Copyright (c) 2016 Roy Xu

package redis

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"time"
)

type Conn interface {
	Send(cmd string, args ...interface{}) (reply interface{}, err error)
	Pipe(cmd string, args ...interface{}) error
	Flush() error
	Receive() (reply interface{}, err error)
	Close() error
}

type conn struct {
	cn      net.Conn
	bw      *bufio.Writer
	br      *bufio.Reader
	timeout time.Duration
}

func Dial(urlstring string) (Conn, error) {
	u, e := url.Parse(urlstring)
	if e != nil {
		log.Fatal(e)
	}
	network := u.Scheme
	address := u.Host
	cn, err := net.DialTimeout(network, address, time.Second*10)
	if err != nil {
		return nil, err
	}
	w := bufio.NewWriter(cn)
	r := bufio.NewReader(cn)
	cli := &conn{cn: cn, bw: w, br: r, timeout: time.Second * 10}
	return cli, nil
}

func (c *conn) Close() error {
	return c.cn.Close()
}

func (c *conn) Send(cmd string, args ...interface{}) (reply interface{}, err error) {
	c.executeCommand(cmd, args...)
	err = c.Flush()
	if err != nil {
		return nil, err
	}
	return c.Receive()
}

func (c *conn) Pipe(cmd string, args ...interface{}) error {
	return c.executeCommand(cmd, args...)
}

func (c *conn) Flush() error {
	return c.bw.Flush()
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
func (c *conn) Receive() (reply interface{}, err error) {
	c.cn.SetWriteDeadline(time.Now().Add(c.timeout))
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
			rep, _ := c.Receive()
			reslt = append(reslt, rep)
		}
		return reslt, nil
	case '-':
		return nil, errors.New(string(rsp))
	default:
		panic(fmt.Sprintf("Protocol Error: %s, %s", string(symbol), string(rsp)))
	}
}

func (c *conn) executeCommand(cmd string, args ...interface{}) (err error) {
	err = c.writeHeader('*', int64(1+len(args)))
	if err != nil {
		return err
	}
	err = c.writeString(cmd)
	for _, arg := range args {
		if err != nil {
			break
		}
		err = c.writeInterface(arg)
	}
	return
}

func (c *conn) write(p []byte) (err error) {
	c.cn.SetWriteDeadline(time.Now().Add(c.timeout))
	_, err = c.bw.Write(p)
	return
}

func (c *conn) writeCRLF() error {
	_, err := c.bw.Write([]byte{'\r', '\n'})
	return err
}

// writeBytes write an RESP head header(*$) to buffer and appends with CR&LF. e.g. $3\r\n, *3\r\n
func (c *conn) writeHeader(symbol byte, n int64) error {
	c.bw.WriteByte(symbol)
	// err = c.writeInt(n) // bugged, recursive calling.
	p := strconv.AppendInt([]byte{}, n, 10)
	c.write(p)
	err := c.writeCRLF()
	return err
}

// writeBytes writes a []byte to buffer with RESP header($) and appends with CR&LF. e.g. $3\r\nSET\r\n
func (c *conn) writeBytes(p []byte) error {
	c.writeHeader('$', int64(len(p)))
	c.write(p)
	err := c.writeCRLF()
	return err
}

// writeInt writes an integer to buffer with RESP header($) and and appends with CR&LF. e.g. $3\r\n125\r\n
func (c *conn) writeInt(n int64) error {
	p := strconv.AppendInt([]byte{}, n, 10)
	return c.writeBytes(p)
}

// writeFloat writes a float to buffer with a RESP header($) and and appends with CR&LF. e.g. $3\r\n1.5\r\n
func (c *conn) writeFloat(f float64) error {
	p := strconv.AppendFloat([]byte{}, f, 'g', -1, 64)
	return c.writeBytes(p)
}

// writeString writes a string to buffer with RESP header($) and and appends with CR&LF. e.g. $4\r\nLLEN\r\n
func (c *conn) writeString(s string) error {
	err := c.writeBytes([]byte(s))
	return err
}

// writeInterface parses parameter 'p' and writes a it to buffer with RESP header($) and and appends with CR&LF. e.g. $4\r\nLLEN\r\n
func (c *conn) writeInterface(p interface{}) (err error) {
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
