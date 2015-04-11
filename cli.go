package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const (
	apiUrl       = "http://i.write.as"
	hiddenApiUrl = "http://writeas7pm7rcdqg.onion"
	readApiUrl   = "http://i.write.as"
	VERSION      = "0.1.0"
)

func main() {
	app := cli.NewApp()
	app.Name = "writeas"
	app.Version = VERSION
	app.Usage = "Simple text pasting and publishing"
	app.Authors = []cli.Author{
		{
			Name:  "Write.as",
			Email: "hello@write.as",
		},
	}
	app.Action = cmdPost
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "tor, t",
			Usage: "Perform action on Tor hidden service",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "post",
			Usage:  "Alias for default action: create post from stdin",
			Action: cmdPost,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "tor, t",
					Usage: "Post to Tor hidden service",
				},
			},
		},
		{
			Name:   "delete",
			Usage:  "Delete a post",
			Action: cmdDelete,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "tor, t",
					Usage: "Delete from Tor hidden service",
				},
			},
		},
		{
			Name:   "get",
			Usage:  "Read a raw post",
			Action: cmdGet,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "tor, t",
					Usage: "Get from Tor hidden service",
				},
			},
		},
	}

	app.Run(os.Args)
}

func readStdIn() []byte {
	numBytes, numChunks := int64(0), int64(0)
	r := bufio.NewReader(os.Stdin)
	fullPost := []byte{}
	buf := make([]byte, 0, 1024)
	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		numChunks++
		numBytes += int64(len(buf))

		fullPost = append(fullPost, buf...)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}

	return fullPost
}

func check(err error) {
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}

func getPass() []byte {
	// TODO: don't show passphrase in the terminal
	var p string
	_, err := fmt.Scanln(&p)
	check(err)
	return []byte(p)
}

func cmdPost(c *cli.Context) {
	fullPost := readStdIn()

	tor := c.Bool("tor") || c.Bool("t")
	if tor {
		fmt.Println("Posting to hidden service...")
	} else {
		fmt.Println("Posting...")
	}

	DoPost(fullPost, false, tor)
}

func cmdDelete(c *cli.Context) {
	friendlyId := c.Args().Get(0)
	token := c.Args().Get(1)
	if friendlyId == "" || token == "" {
		fmt.Println("usage: writeas delete <postId> <token>")
		os.Exit(1)
	}

	tor := c.Bool("tor") || c.Bool("t")

	DoDelete(friendlyId, token, tor)
}

func cmdGet(c *cli.Context) {
	friendlyId := c.Args().Get(0)
	if friendlyId == "" {
		fmt.Println("usage: writeas get <postId>")
		os.Exit(1)
	}

	tor := c.Bool("tor") || c.Bool("t")

	DoFetch(friendlyId, tor)
}

func client(read, tor bool, path, query string) (string, *http.Client) {
	var u *url.URL
	var client *http.Client
	if tor {
		u, _ = url.ParseRequestURI(hiddenApiUrl)
		if read {
			u.Path = "/" + path
		} else {
			u.Path = "/api"
		}
		client = torClient()
	} else {
		if read {
			u, _ = url.ParseRequestURI(readApiUrl)
		} else {
			u, _ = url.ParseRequestURI(apiUrl)
		}
		u.Path = "/" + path
		client = &http.Client{}
	}
	if query != "" {
		u.RawQuery = query
	}
	urlStr := fmt.Sprintf("%v", u)

	return urlStr, client
}

func DoFetch(friendlyId string, tor bool) {
	path := friendlyId
	if len(friendlyId) == 12 {
		// Original (pre-alpha) plain text URLs
		path = friendlyId
	} else if len(friendlyId) == 13 {
		// Alpha phase HTML-based URLs
		path = friendlyId + ".txt"
	} else {
		// Fallback path. Plan is to always support .txt file for raw files
		path = friendlyId + ".txt"
	}

	urlStr, client := client(true, tor, path, "")

	r, _ := http.NewRequest("GET", urlStr, nil)
	r.Header.Add("User-Agent", "writeas-cli v"+VERSION)

	resp, err := client.Do(r)
	check(err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)
		check(err)
		fmt.Printf("%s\n", string(content))
	} else {
		fmt.Printf("Problem getting post: %s\n", resp.Status)
	}
}

func DoPost(post []byte, encrypt, tor bool) {
	data := url.Values{}
	data.Set("w", string(post))
	if encrypt {
		data.Add("e", "")
	}

	urlStr, client := client(false, tor, "", "")

	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	check(err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)
		check(err)
		fmt.Printf("%s\n", string(content))
	} else {
		fmt.Printf("Unable to post: %s\n", resp.Status)
	}
}

func DoDelete(friendlyId, token string, tor bool) {
	urlStr, client := client(false, tor, "", fmt.Sprintf("id=%s&t=%s", friendlyId, token))

	r, _ := http.NewRequest("DELETE", urlStr, nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	check(err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		if tor {
			fmt.Println("Post deleted from hidden service.")
		} else {
			fmt.Println("Post deleted.")
		}
	} else {
		fmt.Printf("Problem deleting: %s\n", resp.Status)
	}
}
