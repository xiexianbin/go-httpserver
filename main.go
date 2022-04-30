package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

var (
	host string
	port int
	dir  string
	help bool
)

func init() {
	flag.StringVar(&host, "h", "0.0.0.0", "listen host")
	flag.IntVar(&port, "p", 8000, "listen port")
	flag.StringVar(&dir, "d", ".", "static files dir")
	flag.BoolVar(&help, "help", false, "show help message")

	flag.Parse()

	flag.Usage = func() {
		fmt.Println("Run: go-httpserver -h 0.0.0.0 -p 8000 -d ./example-site\n" +
			"Usage:")
		flag.PrintDefaults()
	}
}

func handler(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("client from", conn.RemoteAddr(), "closed")
	}(conn)
	fmt.Println("client connected from", conn.RemoteAddr())

	// read request from client net.Conn
	reader := bufio.NewReader(conn)
	req, err := NewRequest(reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(req.String())

	// response
	resp := req.Parse()
	fmt.Println(resp.String())
	err = resp.Do(conn)
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}

func main() {
	if help == true {
		flag.Usage()
		return
	}

	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("close server err", err.Error())
		}
		fmt.Println("server closed.")
	}(listener)
	fmt.Println("listen on:", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		//handler(conn)
		go handler(conn)
	}
}
