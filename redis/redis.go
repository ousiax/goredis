package redis

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Client interface {
	Send(s string) (n int, err error)
	RawReply() (resp []byte, err error)
	Reply() (resp interface{}, err error)
	Set(key, value []byte) (err error)
	Get(key string) (value []byte)
}

type client struct {
	Conn   net.Conn
	Writer *bufio.Writer
	Reader *bufio.Reader
}

func (c *client) Send(s string) (n int, err error) {
	n, err = c.Writer.WriteString(s)
	c.Writer.Flush()
	return
}

func (c *client) reply() (symbl byte, response []byte, err error) {
	line, err := c.Reader.ReadSlice('\n')
	if err != nil {
		return byte('\x00'), nil, err
	}
	return line[0], line[1 : len(line)-3], nil
}

func (c *client) RawReply() ([]byte, error) {
	resp, err := c.Reader.ReadSlice('\n')
	if err != nil {
		return resp, err
	}
	symbol := resp[0] // +,-,:,$,*
	switch symbol {
	case '+', '-', ':':
		break
	case '$':
		line, _ := c.Reader.ReadSlice('\n')
		resp = append(resp, line...)
	case '*':
		d := strings.Trim(string(resp), "*\r\n")
		length, _ := strconv.Atoi(d)
		for i := 0; i < length; i++ {
			line, _ := c.RawReply()
			resp = append(resp, line...)
		}
	default:
		panic(fmt.Sprintf("Protocol Error: %s, %s", string(symbol), string(resp)))
	}
	return resp, err
}

func (c *client) Reply() (interface{}, error) {
	symbl, resp, err := c.reply()
	if err != nil {
		return nil, err
	}
	switch symbl {
	case '+', '-':
		return string(resp), nil
	case ':':
		return strconv.Atoi(string(resp))
	case '$':
		length, _ := strconv.Atoi(string(resp))
		if length == -1 {
			return nil, nil
		} else {
			return c.Reply()
		}
	case '*':
		length, _ := strconv.Atoi(string(resp))
		if length == -1 {
			return nil, nil
		}
		reslt := make([]interface{}, 0)
		for i := 0; i < length; i++ {
			rep, _ := c.Reply()
			reslt = append(reslt, rep)
		}
		return reslt, nil
	default:
		panic(fmt.Sprintf("Protocol Error: %s, %s", string(symbl), string(resp)))
	}
}

func (c *client) Set(key, value []byte) (err error) {
	return nil
}

func (c *client) Get(key string) (value []byte) {
	//line, _ := c.br.ReadSlice('\n')
	return nil
}

func NewClient(network, address string) (cmd Client, err error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	w := bufio.NewWriter(conn)
	rd := bufio.NewReader(conn)
	c := client{Conn: conn, Writer: w, Reader: rd}
	return Client(&c), nil
}
