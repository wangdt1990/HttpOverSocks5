package main

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

func main() {
	fmt.Println("Initiate HttpOverSocks5")

	nl, err := net.Listen("tcp", "localhost:1081")
	if err != nil {
	}

	for {
		nc, err := nl.Accept()
		if err != nil {
		}

		go handle(nc)

		time.Sleep(time.Millisecond * 100)
	}
}

func handle(ncFromClient net.Conn) {
	fmt.Print("\n\n\n")
	fmt.Println("****************************************")
	fmt.Println("HTTP Request From " + ncFromClient.RemoteAddr().String())
	fmt.Println("****************************************")

	var sAddress string
	b := make([]byte, 1024)

	i, err := ncFromClient.Read(b)
	if err != nil {
	}
	fmt.Println("HTTP Request")
	fmt.Println(string(b))

	s := strings.Split(string(b[:]), string('\n'))[0]
	sMethod := strings.Split(s, " ")[0]
	sHost := strings.Split(s, " ")[1]

	//Parse destination host
	uuHost, err := url.Parse(sHost)
	if err != nil {
	}
	fmt.Println("uuHost:", uuHost)
	fmt.Println("uuHost.Scheme:" + uuHost.Scheme)
	fmt.Println("uuHost.Host:" + uuHost.Host)
	fmt.Println("uuHost.Opaque:", uuHost.Opaque)

	if uuHost.Opaque == "443" {
		sAddress = uuHost.Scheme + ":443"
	} else {
		if strings.Index(uuHost.Host, ":") == -1 {
			sAddress = uuHost.Host + ":80"
		} else {
			sAddress = uuHost.Host
		}
	}

	fmt.Println("sAddress:", sAddress)

	//Assemble Socks5 Dialer
	uuSocks5, err := url.Parse("socks5://127.0.0.1:1080")
	if err != nil {
	}

	pdSocks5, err := proxy.FromURL(uuSocks5, proxy.Direct)
	if err != nil {
	}

	ncToServer, err := pdSocks5.Dial("tcp", sAddress)

	if sMethod == "CONNECT" {
		ncFromClient.Write([]byte("HTTP/1.1 200 Connection established\n\n"))
	} else {
		ncToServer.Write(b[:i])
	}

	go io.Copy(ncToServer, ncFromClient)
	io.Copy(ncFromClient, ncToServer)

	ncFromClient.Close()
}
