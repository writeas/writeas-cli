package api

import (
	"fmt"
	"net/http"

	"code.as/core/socks"
)

var (
	TorPort = 9150
)

// TODO: never used?
func torClient() *http.Client {
	dialSocksProxy := socks.DialSocksProxy(socks.SOCKS5, fmt.Sprintf("127.0.0.1:%d", TorPort))
	transport := &http.Transport{Dial: dialSocksProxy}
	return &http.Client{Transport: transport}
}
