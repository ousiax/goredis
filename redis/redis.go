package redis

import (
	"bufio"
	"net"
	//"strconv"
)

type Commander interface {
	Send(s string) (n int, err error)
	Reply() (reply interface{}, err error)
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

func (c *client) Reply() (interface{}, error) {
	line, err := c.Reader.ReadSlice('\n')

	switch line[0] {
	case '+', '-', ':', '$', '*':
		// n, _ = string.Atoi(string(line[1 : len(line)-3]))
		// for i := 0; i < n; i++ {
		// }
	}
	return line, err

}

func (c *client) Set(key, value []byte) (err error) {
	return nil
}

func (c *client) Get(key string) (value []byte) {
	//line, _ := c.br.ReadSlice('\n')
	return nil
}

func NewClient(network, address string) (cmd Commander, err error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	w := bufio.NewWriter(conn)
	rd := bufio.NewReader(conn)
	c := client{Conn: conn, Writer: w, Reader: rd}
	return Commander(&c), nil
}
