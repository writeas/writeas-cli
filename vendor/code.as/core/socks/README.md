SOCKS
=====

[![GoDoc](https://godoc.org/code.as/core/socks?status.svg)](https://godoc.org/code.as/core/socks)

SOCKS is a SOCKS4, SOCKS4A and SOCKS5 proxy package for Go, forked from [h12w/socks](https://github.com/h12w/socks) and patched so it's `go get`able.

## Quick Start
### Get the package

    go get -u "code.as/core/socks"

### Import the package

    import "code.as/core/socks"

### Create a SOCKS proxy dialing function

    dialSocksProxy := socks.DialSocksProxy(socks.SOCKS5, "127.0.0.1:1080")
    tr := &http.Transport{Dial: dialSocksProxy}
    httpClient := &http.Client{Transport: tr}

## Example

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"code.as/core/socks"
)

func main() {
	dialSocksProxy := socks.DialSocksProxy(socks.SOCKS5, "127.0.0.1:1080")
	tr := &http.Transport{Dial: dialSocksProxy}
	httpClient := &http.Client{Transport: tr}
	resp, err := httpClient.Get("http://www.google.com")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.StatusCode)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))
}
```

## Alternatives
http://godoc.org/golang.org/x/net/proxy
