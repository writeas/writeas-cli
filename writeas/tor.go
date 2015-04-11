package main

import (
	"github.com/hailiang/gosocks"
	"net/http"
)

func torClient() *http.Client {
	dialSocksProxy := socks.DialSocksProxy(socks.SOCKS5, "127.0.0.1:9150")
	transport := &http.Transport{Dial: dialSocksProxy}
	return &http.Client{Transport: transport}
}
