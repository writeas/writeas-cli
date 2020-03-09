package impart

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type (
	// Envelope contains metadata and optional data for a response object.
	// Responses will always contain a status code and either:
	// - response Data on a 2xx response, or
	// - an ErrorMessage on non-2xx responses
	//
	// ErrorType is not currently used.
	Envelope struct {
		Code         int         `json:"code"`
		ErrorType    string      `json:"error_type,omitempty"`
		ErrorMessage string      `json:"error_msg,omitempty"`
		Data         interface{} `json:"data,omitempty"`
	}
)

func writeBody(w http.ResponseWriter, body []byte, status int, contentType string) error {
	w.Header().Set("Content-Type", contentType+"; charset=UTF-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(status)
	_, err := w.Write(body)
	return err
}

func RenderActivityJSON(w http.ResponseWriter, value interface{}, status int) error {
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return writeBody(w, body, status, "application/activity+json")
}

func renderJSON(w http.ResponseWriter, value interface{}, status int) error {
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return writeBody(w, body, status, "application/json")
}

func renderString(w http.ResponseWriter, status int, msg string) error {
	return writeBody(w, []byte(msg), status, "text/plain")
}

// WriteSuccess writes the successful data and metadata to the ResponseWriter as
// JSON.
func WriteSuccess(w http.ResponseWriter, data interface{}, status int) error {
	env := &Envelope{
		Code: status,
		Data: data,
	}
	return renderJSON(w, env, status)
}

// WriteError writes the error to the ResponseWriter as JSON.
func WriteError(w http.ResponseWriter, e HTTPError) error {
	env := &Envelope{
		Code:         e.Status,
		ErrorMessage: e.Message,
	}
	return renderJSON(w, env, e.Status)
}

// WriteRedirect sends a redirect
func WriteRedirect(w http.ResponseWriter, e HTTPError) int {
	w.Header().Set("Location", e.Message)
	w.WriteHeader(e.Status)
	return e.Status
}
