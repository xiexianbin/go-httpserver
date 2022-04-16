package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
)

func response(conn net.Conn) {
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
    resp := req.Do()
    fmt.Println(resp.String())
    err = resp.Do(conn)
    if err != nil {
        fmt.Println(err.Error())
    }
    return
}

func main() {
    addr := "0.0.0.0:8888"
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
        go response(conn)
        //response(conn)
    }
}
