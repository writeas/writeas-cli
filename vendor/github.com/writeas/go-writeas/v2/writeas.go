// Package writeas provides the binding for the Write.as API
package writeas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"code.as/core/socks"
	"github.com/writeas/impart"
)

const (
	apiURL    = "https://write.as/api"
	devAPIURL = "https://development.write.as/api"
	torAPIURL = "http://writeas7pm7rcdqg.onion/api"

	// Current go-writeas version
	Version = "2-dev"
)

// Client is used to interact with the Write.as API. It can be used to make
// authenticated or unauthenticated calls.
type Client struct {
	baseURL string

	// Access token for the user making requests.
	token string
	// Client making requests to the API
	client *http.Client

	// UserAgent overrides the default User-Agent header
	UserAgent string
}

// defaultHTTPTimeout is the default http.Client timeout.
const defaultHTTPTimeout = 10 * time.Second

// NewClient creates a new API client. By default, all requests are made
// unauthenticated. To optionally make authenticated requests, call `SetToken`.
//
//     c := writeas.NewClient()
//     c.SetToken("00000000-0000-0000-0000-000000000000")
func NewClient() *Client {
	return NewClientWith(Config{URL: apiURL})
}

// NewTorClient creates a new API client for communicating with the Write.as
// Tor hidden service, using the given port to connect to the local SOCKS
// proxy.
func NewTorClient(port int) *Client {
	return NewClientWith(Config{URL: torAPIURL, TorPort: port})
}

// NewDevClient creates a new API client for development and testing. It'll
// communicate with our development servers, and SHOULD NOT be used in
// production.
func NewDevClient() *Client {
	return NewClientWith(Config{URL: devAPIURL})
}

// Config configures a Write.as client.
type Config struct {
	// URL of the Write.as API service. Defaults to https://write.as/api.
	URL string

	// If specified, the API client will communicate with the Write.as Tor
	// hidden service using the provided port to connect to the local SOCKS
	// proxy.
	TorPort int

	// If specified, requests will be authenticated using this user token.
	// This may be provided after making a few anonymous requests with
	// SetToken.
	Token string
}

// NewClientWith builds a new API client with the provided configuration.
func NewClientWith(c Config) *Client {
	if c.URL == "" {
		c.URL = apiURL
	}

	httpClient := &http.Client{Timeout: defaultHTTPTimeout}
	if c.TorPort > 0 {
		dialSocksProxy := socks.DialSocksProxy(socks.SOCKS5, fmt.Sprintf("127.0.0.1:%d", c.TorPort))
		httpClient.Transport = &http.Transport{Dial: dialSocksProxy}
	}

	return &Client{
		client:  httpClient,
		baseURL: c.URL,
		token:   c.Token,
	}
}

// SetToken sets the user token for all future Client requests. Setting this to
// an empty string will change back to unauthenticated requests.
func (c *Client) SetToken(token string) {
	c.token = token
}

// Token returns the user token currently set to the Client.
func (c *Client) Token() string {
	return c.token
}

func (c *Client) get(path string, r interface{}) (*impart.Envelope, error) {
	method := "GET"
	if method != "GET" && method != "HEAD" {
		return nil, fmt.Errorf("Method %s not currently supported by library (only HEAD and GET).\n", method)
	}

	return c.request(method, path, nil, r)
}

func (c *Client) post(path string, data, r interface{}) (*impart.Envelope, error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(data)
	return c.request("POST", path, b, r)
}

func (c *Client) put(path string, data, r interface{}) (*impart.Envelope, error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(data)
	return c.request("PUT", path, b, r)
}

func (c *Client) delete(path string, data map[string]string) (*impart.Envelope, error) {
	r, err := c.buildRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	q := r.URL.Query()
	for k, v := range data {
		q.Add(k, v)
	}
	r.URL.RawQuery = q.Encode()

	return c.doRequest(r, nil)
}

func (c *Client) request(method, path string, data io.Reader, result interface{}) (*impart.Envelope, error) {
	r, err := c.buildRequest(method, path, data)
	if err != nil {
		return nil, err
	}

	return c.doRequest(r, result)
}

func (c *Client) buildRequest(method, path string, data io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	r, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, fmt.Errorf("Create request: %v", err)
	}
	c.prepareRequest(r)

	return r, nil
}

func (c *Client) doRequest(r *http.Request, result interface{}) (*impart.Envelope, error) {
	resp, err := c.client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("Request: %v", err)
	}
	defer resp.Body.Close()

	env := &impart.Envelope{
		Code: resp.StatusCode,
	}
	if result != nil {
		env.Data = result

		err = json.NewDecoder(resp.Body).Decode(&env)
		if err != nil {
			return nil, err
		}
	}

	return env, nil
}

func (c *Client) prepareRequest(r *http.Request) {
	ua := c.UserAgent
	if ua == "" {
		ua = "go-writeas v" + Version
	}
	r.Header.Set("User-Agent", ua)
	r.Header.Add("Content-Type", "application/json")
	if c.token != "" {
		r.Header.Add("Authorization", "Token "+c.token)
	}
}
