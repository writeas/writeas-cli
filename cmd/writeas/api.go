package main

import (
	"bytes"
	"fmt"
	"github.com/atotto/clipboard"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	defaultUserAgent = "writeas-cli v" + version
)

func userAgent(c *cli.Context) string {
	ua := c.String("user-agent")
	if ua == "" {
		return defaultUserAgent
	}
	return ua + " (" + defaultUserAgent + ")"
}

func client(read, tor bool, path, query string) (string, *http.Client) {
	var u *url.URL
	var client *http.Client
	if tor {
		u, _ = url.ParseRequestURI(hiddenAPIURL)
		u.Path = "/api/" + path
		client = torClient()
	} else {
		u, _ = url.ParseRequestURI(apiURL)
		u.Path = "/api/" + path
		client = &http.Client{}
	}
	if query != "" {
		u.RawQuery = query
	}
	urlStr := fmt.Sprintf("%v", u)

	return urlStr, client
}

// DoFetch retrieves the Write.as post with the given friendlyID,
// optionally via the Tor hidden service.
func DoFetch(friendlyID, ua string, tor bool) error {
	path := friendlyID

	urlStr, client := client(true, tor, path, "")

	r, _ := http.NewRequest("GET", urlStr, nil)
	r.Header.Add("User-Agent", ua)

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", string(content))
	} else if resp.StatusCode == http.StatusNotFound {
		return ErrPostNotFound
	} else if resp.StatusCode == http.StatusGone {
	} else {
		return fmt.Errorf("Unable to get post: %s", resp.Status)
	}
	return nil
}

// DoPost creates a Write.as post, returning an error if it was
// unsuccessful.
func DoPost(c *cli.Context, post []byte, font string, encrypt, tor, code bool) error {
	data := url.Values{}
	data.Set("w", string(post))
	if encrypt {
		data.Add("e", "")
	}
	data.Add("font", getFont(code, font))

	urlStr, client := client(false, tor, "", "")

	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	r.Header.Add("User-Agent", userAgent(c))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		nlPos := strings.Index(string(content), "\n")
		url := content[:nlPos]
		idPos := strings.LastIndex(string(url), "/") + 1
		id := string(url[idPos:])
		token := string(content[nlPos+1 : len(content)-1])

		addPost(id, token)

		// Output URL in requested format
		if c.Bool("md") {
			url = append(url, []byte(".md")...)
		}
		// Copy URL to clipboard
		err = clipboard.WriteAll(string(url))
		if err != nil {
			Errorln("writeas: Didn't copy to clipboard: %s", err)
		} else {
			Info(c, "Copied to clipboard.")
		}

		// Output URL
		fmt.Printf("%s\n", url)
	} else {
		return fmt.Errorf("Unable to post: %s", resp.Status)
	}

	return nil
}

// DoUpdate updates the given post on Write.as.
func DoUpdate(c *cli.Context, post []byte, friendlyID, token, font string, tor, code bool) error {
	urlStr, client := client(false, tor, friendlyID, fmt.Sprintf("t=%s", token))

	data := url.Values{}
	data.Set("w", string(post))

	if code || font != "" {
		// Only update font if explicitly changed
		data.Add("font", getFont(code, font))
	}

	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	r.Header.Add("User-Agent", userAgent(c))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		if tor {
			Info(c, "Post updated via hidden service.")
		} else {
			Info(c, "Post updated.")
		}
	} else {
		if debug {
			ErrorlnQuit("Problem updating: %s", resp.Status)
		} else {
			return fmt.Errorf("Post doesn't exist, or bad edit token given.")
		}
	}
	return nil
}

// DoDelete deletes the given post on Write.as.
func DoDelete(c *cli.Context, friendlyID, token string, tor bool) error {
	urlStr, client := client(false, tor, friendlyID, fmt.Sprintf("t=%s", token))

	r, _ := http.NewRequest("DELETE", urlStr, nil)
	r.Header.Add("User-Agent", userAgent(c))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		if tor {
			Info(c, "Post deleted from hidden service.")
		} else {
			Info(c, "Post deleted.")
		}
		removePost(friendlyID)
	} else {
		if debug {
			ErrorlnQuit("Problem deleting: %s", resp.Status)
		} else {
			return fmt.Errorf("Post doesn't exist, or bad edit token given.")
		}
	}

	return nil
}
