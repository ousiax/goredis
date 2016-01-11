package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	hostPtr := flag.String("host", "127.0.0.1", "Listener host, (default: 127.0.0.1)")
	portPtr := flag.String("p", "2000", "Listener port (default: 2000)")
	flag.Parse()

	address := *hostPtr + ":" + *portPtr

	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	ch := make(chan int)
	go func() {
		for {
			// Wati for a connection.
			conn, err := l.Accept()
			if err != nil {
				log.Fatal(err)
			}
			go func(c net.Conn) {
				defer func() {
					c.Close()
					recover()
				}()
				rd := bufio.NewReader(c)
				for {
					if s, err := rd.ReadString('\n'); err == nil {
						fmt.Print(s)
						v := strings.ToLower(strings.Trim(s, "\r\n "))
						if v == "shutdown" {
							ch <- 1
							break
						} else if v == "close" {
							break
						}
					}
				}
			}(conn)
		}
	}()

	<-ch
}
