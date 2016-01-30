// The MIT License (MIT)

// Copyright (c) 2016 Roy Xu

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	network = "tcp"
	host    = "127.0.0.1"
	port    = 2433
	timeout = 30
)

func main() {
	hPtr := flag.String("host", host, fmt.Sprintf("Echo Server host. (Default: %s", host))
	pPtr := flag.Int("port", port, fmt.Sprintf("Server port (default: %d).", port))
	flag.Parse()
	l, err := net.Listen("tcp", *hPtr+":"+strconv.Itoa(*pPtr))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		cn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Connected: %s", cn.RemoteAddr())
		go func() {
			defer func() {
				cn.Close()
				log.Printf("Disconnected: %s", cn.RemoteAddr())
			}()

			r := bufio.NewReader(cn)
			w := bufio.NewWriter(cn)
			for {
				cn.SetReadDeadline(time.Now().Add(time.Second * timeout))
				s, e := r.ReadString('\n')
				if e != nil {
					log.Printf("Disconnected:%s, Error:%s", cn.RemoteAddr(), e)
					break
				}
				if cmd := strings.ToLower(strings.TrimSpace(s)); cmd == "quit" {
					break
				} else if cmd == "shutdown" {
					l.Close()
				}
				cn.SetWriteDeadline(time.Now().Add(time.Second * timeout))
				w.WriteString(s)
				if w.Flush() != nil {
					break
				}
			}
		}()
	}
}
