// The MIT License (MIT)
//
// Copyright (c) 2016 Roy Xu

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/qqbuby/goredis/redis"
	"os"
	"strconv"
	"strings"
)

func print(p interface{}) {
	switch v := p.(type) {
	case string, int:
		fmt.Println(v)
	case []byte:
		fmt.Println(string(v))
	case []interface{}:
		for _, k := range v {
			print(k)
		}
	default:
		fmt.Println(v)
	}
}

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
	defer func() {
		client.Close()
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s>", address)

		raw, _ := reader.ReadString('\n')
		raw = strings.Trim(raw, "\n ")
		if len(raw) == 0 {
			continue
		}
		if strings.ToLower(raw) == "quit" || strings.ToLower(raw) == "exit" {
			break
		}
		s := strings.Split(raw, " ")
		args := make([]interface{}, len(s[1:]))
		for i := 0; i < len(args); i++ {
			args[i] = s[1+i]
		}
		rsp, e := client.Send(s[0], args...)
		if e == nil {
			print(rsp)
		} else {
			fmt.Printf("%v\n", err.Error())
		}
	}
}
