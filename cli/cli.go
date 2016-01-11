package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/qqbuby/redis-go/redis"
	"os"
	"strconv"
	"strings"
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
		// command := make([]byte, 0) // if the cap > 0, slice always insert a \x00, why?
		// command = append(command, fmt.Sprintf("*%s\r\n", strconv.Itoa(len(s)))...)
		// for _, p := range s {
		// 	command = append(command, fmt.Sprintf("$%d\r\n%s\r\n", len(p), p)...)
		// }
		// //fmt.Printf("%q\n", command) // for debug to output raw command bytes
		// resp, e := client.Send(string(command))
		args := make([]interface{}, len(s[1:]))
		for i := 0; i < len(args); i++ {
			args[i] = s[1+i]
		}
		resp, e := client.Send(s[0], args...)
		if e == nil {
			fmt.Printf("%v\n", resp)
		} else {
			fmt.Printf("%v\n", err.Error())
		}
	}
}
