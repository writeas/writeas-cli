package impart

import (
	"net/http"
)

// HTTPError holds an HTTP status code and an error message.
type HTTPError struct {
	Status  int
	Message string
}

// Error displays the HTTPError's error message and satisfies the error
// interface.
func (h HTTPError) Error() string {
	if h.Message == "" {
		return http.StatusText(h.Status)
	}
	return h.Message
}
