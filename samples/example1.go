// The MIT License (MIT)
//
// Copyright (c) 2016 Roy Xu

package main

import (
	"flag"
	"fmt"
	"github.com/qqbuby/goredis/redis"
	"strconv"
)

func main() {
	hostnamePtr := flag.String("h", "127.0.0.1", "Server hostname (default: 127.0.0.1).")
	portPtr := flag.Int("p", 6379, "Server port (default: 6379).")
	flag.Parse()

	network := "tcp"
	address := *hostnamePtr + ":" + strconv.Itoa(*portPtr)

	client, err := redis.NewClient(network, address)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	var key string = "key"
	client.Set(key, "hello world")
	v, _ := client.Get(key)
	fmt.Println(v)
}
