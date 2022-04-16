package main

import (
    "bufio"
    "fmt"
    "io"
    "net"
)

type Response struct {
    Protocol   string
    StatusCode int
    Status     string
    Headers    map[string]string
    Body       interface{}
}

func (resp Response) String() string {
    return fmt.Sprintf("Protocol: %s, StatusCode: %d, Status: %s, Headers: %v, Body: %v", resp.Protocol, resp.StatusCode, resp.Status, resp.Headers, resp.Body)
}

func (resp Response) Do(conn net.Conn) error {
    //fmt.Fprint(conn, fmt.Sprintf("%s %d %s\n", resp.Protocol, resp.StatusCode, resp.Status))
    //for k, v := range resp.Headers {
    //    fmt.Fprint(conn, fmt.Sprintf("%s: %s\n", k, v))
    //}
    //fmt.Fprint(conn, "\r\n")
    //fmt.Fprint(conn, resp.Body, "\r\n")

    var err error
    var respWriter *bufio.Writer = bufio.NewWriter(conn)

    //response info
    _, err = respWriter.WriteString(fmt.Sprintf("%s %d %s\r\n", resp.Protocol, resp.StatusCode, resp.Status))
    if err != nil {
        fmt.Println(err.Error())
    }

    //response header
    for k, v := range resp.Headers {
        _, err := respWriter.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
        if err != nil {
            fmt.Println(err.Error())
        }
    }

    _, err = respWriter.WriteString("\r\n")
    if err != nil {
        fmt.Println(err.Error())
    }

    //response body
    switch resp.Body.(type) {
    case string:
        body := resp.Body.(string)
        _, err = respWriter.WriteString(fmt.Sprintf("%s\r\n", body))
        if err != nil {
            fmt.Println(err.Error())
        }
    case io.Reader:
        body := resp.Body.(io.Reader)
        _, err := respWriter.ReadFrom(body)
        if err != nil {
            fmt.Println(err.Error())
        }
    }

    err = respWriter.Flush()
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    return nil
}
