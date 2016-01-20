// The MIT License (MIT)

// Copyright (c) 2016 Roy Xu

package redis

type connectionPooling struct {
}

func (p *connectionPooling) Pop() (*connection, error) {
	return nil, nil
}

func (p *connectionPooling) Relese(*connection) error {
	return nil
}

type connection struct {
	host                 string
	port                 int
	db                   int
	password             string
	socketTimeout        int
	socketConnectTimeout int
}

func (c *connection) Connect() error {
	return nil
}

func (c *connection) Disconnect() error {
	return nil
}

type sslconnection struct {
}

type unixDomainSocketConnection struct {
}

func NewPooling(urlstring string) (*connectionPooling, error) {
	return nil, nil
}
