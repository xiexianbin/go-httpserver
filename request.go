package main

import (
    "bufio"
    "fmt"
    "io"
    "path/filepath"
    "strings"
)

type Request struct {
    Method   string
    Url      string
    Protocol string
    Headers  map[string]string
    Body     interface{}
}

func NewRequest(reqReader *bufio.Reader) (*Request, error) {
    var req = Request{
        Headers: map[string]string{},
    }
    var builder strings.Builder

    // read request info
    infoLine, err := reqReader.ReadString('\n')
    if err != nil {
        fmt.Println(err.Error())
    }
    infoLine = strings.TrimSpace(infoLine)
    //requestInfo := strings.Split(infoLine, " ")
    requestInfo := strings.Fields(infoLine)
    if len(requestInfo) == 3 {
        req.Method = requestInfo[0]
        req.Url = requestInfo[1]
        req.Protocol = requestInfo[2]
    } else {
        return nil, fmt.Errorf("read wrong request info: %s", infoLine)
    }

    // read request header
    for {
        headerLine, err := reqReader.ReadString('\n')
        if err != nil {
            fmt.Println(err.Error())
        }
        headerLine = strings.TrimSpace(headerLine)
        if headerLine == "" {
            break
        } else if strings.Contains(headerLine, ":") {
            h := strings.Split(headerLine, ":")
            req.Headers[h[0]] = h[1]
        }
    }

    // read request body
    bs := make([]byte, 1024)
    switch req.Method {
    case POST, PUT:
        for {
            n, err := reqReader.Read(bs)
            if err != nil {
                if err == io.EOF {
                    break
                }
                fmt.Println(err.Error())
            }
            builder.Write(bs[:n])
        }
    }
    req.Body = builder.String()
    return &req, nil
}

func (req Request) String() string {
    return fmt.Sprintf("Method: %s, Url: %s, Protocal: %s, Headers: %v, Body: %v", req.Method, req.Url, req.Protocol, req.Headers, req.Body)
}

func (req Request) Do() *Response {
    switch req.Method {
    case GET:
        // read current dir file by url, and return
        filePath := ParseUrl(req.Url)

        file, err := GetFile(filePath)
        if err != nil {
            return &Response{
                Protocol:   "http/1.1",
                StatusCode: 404,
                Status:     "Not Found",
                Headers: map[string]string{
                    "content-type": "text/html",
                    "server":       "go-httpserver",
                },
            }
        }

        body, err := GetFileStream(file)
        if err != nil {
            return &Response{
                Protocol:   "http/1.1",
                StatusCode: 500,
                Status:     "Server Error",
                Headers: map[string]string{
                    "content-type": "text/html",
                    "server":       "go-httpserver",
                },
            }
        }
        return &Response{
            Protocol:   "http/1.1",
            StatusCode: 200,
            Status:     "OK",
            Headers: map[string]string{
                "content-type": TypeByExtension(filepath.Ext(file)),
                "server":       "go-httpserver",
            },
            Body: body,
        }
    }
    return nil
}
