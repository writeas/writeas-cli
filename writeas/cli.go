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
	"strings"
)

const (
	apiUrl       = "http://i.write.as"
	hiddenApiUrl = "http://writeas7pm7rcdqg.onion"
	readApiUrl   = "http://i.write.as"
	VERSION      = "0.1.0"
)

func main() {
	initialize()

	// Run the app
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
		cli.IntFlag{
			Name:  "tor-port",
			Usage: "Use a different port to connect to Tor",
			Value: 9150,
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
				cli.IntFlag{
					Name:  "tor-port",
					Usage: "Use a different port to connect to Tor",
					Value: 9150,
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
				cli.IntFlag{
					Name:  "tor-port",
					Usage: "Use a different port to connect to Tor",
					Value: 9150,
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
				cli.IntFlag{
					Name:  "tor-port",
					Usage: "Use a different port to connect to Tor",
					Value: 9150,
				},
			},
		},
	}

	app.Run(os.Args)
}

func initialize() {
	// Ensure we have a data directory to use
	if !dataDirExists() {
		createDataDir()
	}
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
	if c.Int("tor-port") != 0 {
		torPort = c.Int("tor-port")
	}
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
	if friendlyId == "" {
		fmt.Println("usage: writeas delete <postId> [<token>]")
		os.Exit(1)
	}

	if token == "" {
		// Search for the token locally
		token = tokenFromID(friendlyId)
		if token == "" {
			fmt.Println("Couldn't find an edit token locally. Did you create this post here?")
			fmt.Printf("If you have an edit token, use: writeas delete %s <token>\n", friendlyId)
			os.Exit(1)
		}
	}

	tor := c.Bool("tor") || c.Bool("t")
	if c.Int("tor-port") != 0 {
		torPort = c.Int("tor-port")
	}

	DoDelete(friendlyId, token, tor)
}

func cmdGet(c *cli.Context) {
	friendlyId := c.Args().Get(0)
	if friendlyId == "" {
		fmt.Println("usage: writeas get <postId>")
		os.Exit(1)
	}

	tor := c.Bool("tor") || c.Bool("t")
	if c.Int("tor-port") != 0 {
		torPort = c.Int("tor-port")
	}

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
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("Post not found.\n")
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

		nlPos := strings.Index(string(content), "\n")
		url := content[:nlPos]
		idPos := strings.LastIndex(string(url), "/") + 1
		id := string(url[idPos:])
		token := string(content[nlPos+1 : len(content)-1])

		addPost(id, token)

		fmt.Printf("%s\n", url)
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
		removePost(friendlyId)
	} else {
		if DEBUG {
			fmt.Printf("Problem deleting: %s\n", resp.Status)
		} else {
			fmt.Printf("Post doesn't exist, or bad edit token given.")
		}
	}
}
