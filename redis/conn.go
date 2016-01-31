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
	c.execute(cmd, args...)
	err = c.Flush()
	if err != nil {
		return nil, err
	}
	return c.Receive()
}

func (c *conn) Pipe(cmd string, args ...interface{}) error {
	return c.execute(cmd, args...)
}

func (c *conn) Flush() error {
	return c.bw.Flush()
}

// Receive returns the raw RESP reply without the symbols(+-:$*) and CR&LF.
// []byte
//    For Simple Strings the first byte of the reply is "+"
//    For Integers the first byte of the reply is ":"
//    For Bulk Strings the first byte of the reply is "$"
// error
//    For Errors the first byte of the reply is "-"
// []interface{} (=[][]byte)
//    For Arrays the first byte of the reply is "*"
// nil
//    For Null Bulk String
//    For Null Array
func (c *conn) Receive() (reply interface{}, err error) {
	p, m, e := c.response()
	if e != nil {
		return nil, e
	}
	switch p {
	case '+', ':':
		return m, nil
	case '-':
		return errors.New(string(m)), nil
	case '$':
		l, _ := strconv.Atoi(string(m)) // todo panic protocol error.
		if l == -1 {
			return nil, nil
		} else {
			m, e := c.br.ReadBytes('\n')
			if e != nil {
				return nil, e
			}
			return m[0 : len(m)-2], nil
		}
	case '*':
		l, _ := strconv.Atoi(string(m)) // todo panic protocol error.
		if l == -1 {
			return nil, nil
		} else {
			a := make([]interface{}, l)
			for i := 0; i < l; i++ {
				m, e := c.Receive()
				a[i] = m
				if e != nil {
					return nil, e
				}
			}
			return a, nil
		}
	default:
		panic(fmt.Sprintf("Protocol Error: %s. [+-:*]", string(p)))
	}
}

// response parses a RESP reply to [prefix],[payloadmessage],[\r\n].
func (c *conn) response() (prefix byte, message []byte, err error) {
	c.cn.SetReadDeadline(time.Now().Add(c.timeout))
	b, e := c.br.ReadSlice('\n')
	if e != nil {
		return '\000', []byte{}, e
	}
	p, m := b[0], b[1:len(b)-2]
	return p, m, nil
}

func (c *conn) execute(cmd string, args ...interface{}) (err error) {
	l := 1 + len(args)
	c.wLen('*', l)
	c.wBulkString(cmd)
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			err = c.wBulkString(v)
		case int:
			err = c.wBulkInt(v)
		case int32:
			err = c.wBulkInt32(v)
		case int64:
			err = c.wBulkInt64(v)
		case float32:
			err = c.wBulkFloat32(v)
		case float64:
			err = c.wBulkFloat64(v)
		case bool:
			err = c.wBulkBool(v)
		case nil:
			err = c.wBulkNil()
		case uint:
			err = c.wBulkUint(v)
		case uint8:
			err = c.wBulkUint8(v)
		case uint16:
			err = c.wBulkUint16(v)
		case uint32:
			err = c.wBulkUint32(v)
		case uint64:
			err = c.wBulkUint64(v)
		default:
			err = c.wBulkOther(v)
		}
	}
	return
}

func (c *conn) wLen(p byte, i int) (err error) {
	c.bw.WriteByte(p)
	c.wInt(i)
	return c.wCRLF()
}

func (c *conn) wBulkString(s string) (err error) {
	l := len(s)
	c.wLen('$', l)
	c.bw.WriteString(s)
	err = c.wCRLF()
	return
}

func (c *conn) wBulkByte(b byte) (err error) {
	c.wLen('$', 1)
	c.bw.WriteByte(b)
	err = c.wCRLF()
	return
}

func (c *conn) wBulkBytes(b []byte) (err error) {
	l := len(b)
	c.wLen('$', l)
	c.bw.Write(b)
	err = c.wCRLF()
	return
}

func (c *conn) wBulkInt(i int) (err error) {
	return c.wBulkInt64(int64(i))
}

func (c *conn) wBulkInt32(i int32) (err error) {
	return c.wBulkInt64(int64(i))
}

func (c *conn) wBulkInt64(i int64) (err error) {
	b := strconv.AppendInt([]byte{}, i, 10)
	l := len(b)
	c.wLen('$', l)
	c.bw.Write(b)
	err = c.wCRLF()
	return
}

func (c *conn) wBulkUint(i uint) (err error) {
	return c.wBulkUint64(uint64(i))
}

func (c *conn) wBulkUint8(i uint8) (err error) {
	return c.wBulkUint64(uint64(i))
}

func (c *conn) wBulkUint16(i uint16) (err error) {
	return c.wBulkUint64(uint64(i))
}

func (c *conn) wBulkUint32(i uint32) (err error) {
	return c.wBulkUint64(uint64(i))
}

func (c *conn) wBulkUint64(i uint64) (err error) {
	b := strconv.AppendUint([]byte{}, i, 10)
	l := len(b)
	c.wLen('$', l)
	c.bw.Write(b)
	err = c.wCRLF()
	return
}

func (c *conn) wBulkFloat32(f float32) (err error) {
	return c.wBulkFloat64(float64(f))
}

func (c *conn) wBulkFloat64(f float64) (err error) {
	b := strconv.AppendFloat([]byte{}, f, 'g', -1, 64)
	l := len(b)
	c.wLen('$', l)
	c.bw.Write(b)
	err = c.wCRLF()
	return
}

func (c *conn) wBulkBool(b bool) (err error) {
	c.wLen('$', 1)
	if b {
		err = c.bw.WriteByte('1')
	} else {
		err = c.bw.WriteByte('0')
	}
	return
}

func (c *conn) wBulkNil() (err error) {
	c.wLen('$', 0)
	_, err = c.bw.WriteString("")
	return
}

func (c *conn) wBulkOther(o interface{}) (err error) {
	var buf bytes.Buffer
	fmt.Fprint(&buf, o)
	b := buf.Bytes()
	l := len(b)
	c.wLen('$', l)
	c.wBulkBytes(b)
	return
}

func (c *conn) wCRLF() (err error) {
	const crlf = "\r\n"
	_, err = c.bw.WriteString(crlf)
	return
}

func (c *conn) wInt(i int) (err error) {
	return c.wInt64(int64(i))
}

func (c *conn) wInt64(i int64) (err error) {
	m := strconv.AppendInt([]byte{}, i, 10)
	_, err = c.bw.Write(m)
	return
}
