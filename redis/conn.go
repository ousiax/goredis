package redis

import (
	"bufio"
	"net"
)

type conn interface {
	Receive() (reply interface{}, err error)
	SyncSend(commandName string, args ...interface{}) error
	AsyncSend(commandName string, args ...interface{}) error
	Flush() error
}

type Conn struct {
	conn   *net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
}

func Dial(rawUrl string) (c Conn, err error) {
	return Conn{}, nil
}
