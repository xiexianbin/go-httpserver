package main

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

func GetPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return filepath.Dir(path)
}

func ParseUrl(url string) string {
	if strings.LastIndex(url, "/") == len(url)-1 {
		return fmt.Sprintf("%s%s", url, INDEX)
	}

	return url
}

func GetFile(subPath string) (string, error) {
	fPath, _ := filepath.Abs(filepath.Join(dir, subPath))
	info, err := os.Stat(fPath)
	if err != nil {
		if os.IsNotExist(err) {
			return PAGE404, nil
		}
		fmt.Println(err.Error())
		return "", err
	}
	// dir or file
	if info.IsDir() {
		// dir
		fPath = filepath.Join(fPath, INDEX)
		_, err = os.Stat(fPath)
		if os.IsNotExist(err) {
			return PAGE404, nil
		}
	}

	return fPath, nil
}

func GetFileStream(path string) (io.Reader, error) {
	// is file exist
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return os.Open(path)
}

func TypeByExtension(ext string) string {
	return mime.TypeByExtension(ext)
}
