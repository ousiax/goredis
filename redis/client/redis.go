package client

import (
	"github.com/qqbuby/goredis/redis"
)

type Client struct {
	conn *redis.Conn
}

func New(rawurl string) (cli Client, err error) {
	return Client{}, nil
}

func (c *Client) Send(commandName string, args ...interface{}) (reply interface{}, err error) {
	return nil, nil
}
