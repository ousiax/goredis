package main

import (
	"bufio"
	"fmt"
	"github.com/qqbuby/redis-go/redis"
	"os"
	"strconv"
	"strings"
)

func main() {
	const network = "tcp"
	const address = "127.0.0.1:6379"

	client, err := redis.NewClient(network, address)
	if err != nil {
		fmt.Println(err)
		return
	}

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
		command := make([]byte, 0) // if the cap > 0, slice always insert a \x00, why?
		command = append(command, fmt.Sprintf("*%s\r\n", strconv.Itoa(len(s)))...)
		for _, p := range s {
			command = append(command, fmt.Sprintf("$%d\r\n%s\r\n", len(p), p)...)
		}
		// fmt.Printf("%s\n") // for debug to output raw command bytes
		client.Send(string(command))
		resp, e := client.RawReply()
		if e == nil {
			fmt.Print(string(resp))
		}
	}
}
