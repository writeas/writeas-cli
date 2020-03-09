package impart

import (
	"mime"
	"net/http"
)

// ReqJSON returns whether or not the given Request is sending JSON, based on
// the Content-Type header being application/json.
func ReqJSON(r *http.Request) bool {
	ct, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	return ct == "application/json"
}
