// The MIT License (MIT)

// Copyright (c) 2016 Roy Xu

// A simple tcp transparent proxy.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	network    = "tcp"
	remoteHost = "192.168.128.134"
	remoteport = 6379
	proxyHost  = "127.0.0.1"
	proxyPort  = 6379
)

func main() {
	rHostPtr := flag.String("rhost", remoteHost, fmt.Sprintf("Remote Server host. (Default: %s", remoteHost))
	rPortPtr := flag.Int("rport", remoteport, fmt.Sprintf("Remote Server port (default: %d).", remoteport))
	hPtr := flag.String("host", proxyHost, fmt.Sprintf("Proxy Server host. (Default: %s", proxyHost))
	pPtr := flag.Int("port", proxyPort, fmt.Sprintf("Proxy Server port (default: %d).", proxyPort))
	flag.Parse()
	proxy, err := net.Listen(network, *hPtr+":"+strconv.Itoa(*pPtr))
	if err != nil {
		log.Fatal(err)
	}
	defer proxy.Close()

	heart, err := net.Dial(network, *rHostPtr+":"+strconv.Itoa(*rPortPtr))
	if err != nil {
		log.Fatalf("CONNECTION(REMOTE) REFUSED:%s", err)
	}
	defer heart.Close()

	go func() {
		var ping = []byte("PING")
		for {
			time.Sleep(time.Millisecond * 1000)
			heart.SetWriteDeadline(time.Now().Add(time.Millisecond * 150))
			_, err := heart.Write(ping)
			if err != nil {
				log.Fatalf("DISCONNECTED(REMOTE(HEART)):%s", err)
			}
			log.Printf(". H . E . A . R . T . B . E . A . T .")
		}
	}()

	for {
		lcn, err := proxy.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("CONNECTED(PROXY): %s", lcn.RemoteAddr())
		go func() {
			rcn, err := net.Dial(network, *rHostPtr+":"+strconv.Itoa(*rPortPtr))
			if err != nil {
				log.Printf("CONNECTION(REMOTE)  REFUSED:%s", err)
			}
			log.Printf("CONNECTED(REMOTE): %s", rcn.RemoteAddr())
			defer func() {
				rcn.Close()
				log.Printf("DISCONNECTED(REMOTE): %s", rcn.RemoteAddr())
				lcn.Close()
				log.Printf("DISCONNECTED(PROXY): %s", lcn.RemoteAddr())
			}()

			go netcopy(rcn, lcn) //io.Copy doesn't support timeout ?
			netcopy(lcn, rcn)
		}()
	}
}

func netcopy(dst net.Conn, src net.Conn) (err error) {
	var buf = make([]byte, 512)
	n := 0
	for {
		src.SetReadDeadline(time.Now().Add(time.Second * 150))
		n, err = src.Read(buf)
		if err != nil {
			break
		}
		dst.SetWriteDeadline(time.Now().Add(time.Millisecond * 150))
		_, err = dst.Write(buf[0:n])
		if err != nil {
			break
		}
	}
	return err
}
