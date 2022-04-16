package main

import (
    "fmt"
    "testing"
)

func TestGetPath(t *testing.T) {
    fmt.Println(GetPath())
}

func TestParseUrl(t *testing.T) {
    testUrl := map[string]string{
        "/":     "/index.html",
        "/abc/": "/abc/index.html",
    }
    for k, v := range testUrl {
        if p := ParseUrl(k); v != p {
            t.Errorf("%v != %v", v, p)
        }
    }
}

func TestTypeByExtension(t *testing.T) {
    fmt.Println(TypeByExtension(".png"))
    fmt.Println(TypeByExtension(".html"))
    fmt.Println(TypeByExtension(".ico"))
}
