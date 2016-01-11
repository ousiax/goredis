package redis

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strconv"
)

type Client interface {
	respStrings
	Send(s string) (reply interface{}, err error)
	//Receive() (resp interface{}, err error)
	//RawReceive() (resp []byte, err error)
	Close() (err error)
}

type client struct {
	conn net.Conn
	bw   *bufio.Writer
	br   *bufio.Reader
}

func (c *client) Send(s string) (reply interface{}, err error) {
	c.bw.WriteString(s)
	c.bw.Flush()
	return c.receive()
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

// func (c *client) RawReceive() ([]byte, error) {
// 	resp, err := c.br.ReadSlice('\n')
// 	if err != nil {
// 		return resp, err
// 	}
// 	symbol := resp[0] // +,-,:,$,*
// 	switch symbol {
// 	case '+', '-', ':':
// 		break
// 	case '$':
// 		line, _ := c.br.ReadSlice('\n')
// 		resp = append(resp, line...)
// 	case '*':
// 		d := strings.Trim(string(resp), "*\r\n")
// 		length, _ := strconv.Atoi(d)
// 		for i := 0; i < length; i++ {
// 			line, _ := c.RawReceive()
// 			resp = append(resp, line...)
// 		}
// 	default:
// 		panic(fmt.Sprintf("Protocol Error: %s, %s", string(symbol), string(resp)))
// 	}
// 	return resp, err
// }

func (c *client) Close() (err error) {
	return c.conn.Close()
}

func (c *client) writeBytes(p []byte) error {
	err := c.writeHeader('$', int64(len(p)))
	if err != nil {
		return err
	}
	_, err = c.bw.Write(p)
	if err == nil {
		_, err = c.bw.Write([]byte{'\r', '\n'})
	}
	return err
}

//  2\r\n
func (c *client) writeInt(n int64) error {
	var p []byte
	p = strconv.AppendInt(p, n, 10)
	return c.writeBytes(p)
}

func (c *client) writeFloat(f float64) error {
	var p []byte
	p = strconv.AppendFloat(p, f, 'g', -1, 64)
	return c.writeBytes(p)
}

//  *2\r\n
func (c *client) writeHeader(symbol byte, n int64) error {
	err := c.bw.WriteByte(symbol)
	if err == nil {
		err = c.writeInt(n)
	}
	return err
}

//  LLEN\r\n
func (c *client) writeString(s string) error {
	return c.writeBytes([]byte(s))
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
		switch arg.(type) {
		case string:
			err = c.writeString(arg.(string))
		case []byte:
			err = c.writeBytes(arg.([]byte))
		case int:
			err = c.writeInt(int64(arg.(int)))
		case float32:
			err = c.writeFloat(float64(arg.(float32)))
		case float64:
			err = c.writeFloat(arg.(float64))
		case bool:
			if arg.(bool) {
				err = c.writeString("1")
			} else {
				err = c.writeString("0")
			}
		case nil:
			err = c.writeString("")
		default:
			var buf bytes.Buffer
			fmt.Fprint(&buf, arg)
			err = c.writeBytes(buf.Bytes())
		}
	}
	if err == nil {
		err = c.bw.Flush()
	}
	if err == nil {
		return c.receive()
	}
	return nil, err
}

func (c *client) Get(key string) (value interface{}, err error) {
	return c.executeCommand("GET", key)
}

func (c *client) Set(key string, args ...interface{}) (string, error) {
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
