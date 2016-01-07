// A client sends to the Redis server a RESP Array consiting of just Bulk Strings.
// A Redis server replies to clients sending any valid RESP data type as reply
// For Simple Strings the first byte of the reply is "+"
// For Errors the first byte of the reply is "-"
// For Integers the first byte of the reply is ":"
// For Bulk Strings the first byte of the reply is "$"
// For Arrays the first byte of the reply is "*"
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("net: could not connect to redis server.")
		return
	}
	defer conn.Close()
	r := bufio.NewReader(os.Stdin)
	fmt.Println("press 'quit' to exit.")
	for {
		fmt.Print(">")
		raw, _ := r.ReadString('\n')
		raw = strings.Trim(raw, "\n ")
		if strings.ToLower(raw) == "quit" {
			break
		}
		a := strings.Split(raw, " ")
		cmd := make([]byte, len(raw))
		cmd = append(cmd, []byte(fmt.Sprintf("*%d\r\n", len(a)))...)
		for _, s := range a {
			cmd = append(cmd, []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))...)
		}
		fmt.Println(string(cmd))
	}
}
