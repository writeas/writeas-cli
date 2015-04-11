package main

import (
	"fmt"
	"github.com/hailiang/gosocks"
	"net/http"
)

var (
	torPort = 9150
)

func torClient() *http.Client {
	dialSocksProxy := socks.DialSocksProxy(socks.SOCKS5, fmt.Sprintf("127.0.0.1:%d", torPort))
	transport := &http.Transport{Dial: dialSocksProxy}
	return &http.Client{Transport: transport}
}
