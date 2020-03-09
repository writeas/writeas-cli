impart
======

![MIT license](https://img.shields.io/github/license/writeas/impart.svg) [![#writeas on freenode](https://img.shields.io/badge/freenode-%23writeas-blue.svg)](http://webchat.freenode.net/?channels=writeas) [![Public Slack discussion](http://slack.write.as/badge.svg)](http://slack.write.as/)

**impart** is a library for the final layer between the API and the consumer. It's used in the latest [Write.as](https://write.as) and [HTMLhouse](https://html.house) APIs.

We're still in the early stages of development, so there may be breaking changes.

## Example use

```go
package main

import (
	"fmt"
	"github.com/writeas/impart"
	"net/http"
)

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	http.HandleFunc("/", handle(index))
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "Hello world!")

	return nil
}

func handle(f handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleError(w, r, func() error {
			// Do authentication...

			// Handle the request
			err := f(w, r)

			// Log the request and result...

			return err
		}())
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	if err, ok := err.(impart.HTTPError); ok {
		impart.WriteError(w, err)
		return
	}

	impart.WriteError(w, impart.HTTPError{http.StatusInternalServerError, "Internal server error :("})
}
```
