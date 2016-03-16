package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/atotto/clipboard"
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

// API constants for communicating with Write.as.
const (
	apiURL       = "https://write.as"
	hiddenAPIURL = "http://writeas7pm7rcdqg.onion"
	readAPIURL   = "https://write.as"
)

// Application constants.
const (
	version = "0.4"
)

// Defaults for posts on Write.as.
const (
	defaultFont = PostFontMono
)

// Available flags for creating posts
var postFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "tor, t",
		Usage: "Perform action on Tor hidden service",
	},
	cli.IntFlag{
		Name:  "tor-port",
		Usage: "Use a different port to connect to Tor",
		Value: 9150,
	},
	cli.BoolFlag{
		Name:  "code",
		Usage: "Specifies this post is code",
	},
	cli.StringFlag{
		Name:  "font",
		Usage: "Sets post font to given value",
		Value: defaultFont,
	},
}

func main() {
	initialize()

	// Run the app
	app := cli.NewApp()
	app.Name = "writeas"
	app.Version = version
	app.Usage = "Simple text pasting and publishing"
	app.Authors = []cli.Author{
		{
			Name:  "Write.as",
			Email: "hello@write.as",
		},
	}
	app.Action = cmdPost
	app.Flags = postFlags
	app.Commands = []cli.Command{
		{
			Name:   "post",
			Usage:  "Alias for default action: create post from stdin",
			Action: cmdPost,
			Flags:  postFlags,
			Description: `Create a new post on Write.as from stdin.

   Use the --code flag to indicate that the post should use syntax 
   highlighting. Or use the --font [value] argument to set the post's 
   appearance, where [value] is mono, monospace (default), wrap (monospace 
   font with word wrapping), serif, or sans.`,
		},
		{
			Name:  "new",
			Usage: "Compose a new post from the command-line and publish",
			Description: `An alternative to piping data to the program.

   On Windows, this will use 'copy con' to start reading what you input from the
   prompt. Press F6 or Ctrl-Z then Enter to end input.
   On *nix, this will use the best available text editor, starting with the 
   value set to the EDITOR environment variable, or vim, or finally nano.

   Use the --code flag to indicate that the post should use syntax 
   highlighting. Or use the --font [value] argument to set the post's 
   appearance, where [value] is mono, monospace (default), wrap (monospace 
   font with word wrapping), serif, or sans.
   
   If posting fails for any reason, 'writeas' will show you the temporary file
   location and how to pipe it to 'writeas' to retry.`,
			Action: cmdNew,
			Flags:  postFlags,
		},
		{
			Name:   "delete",
			Usage:  "Delete a post",
			Action: cmdDelete,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "tor, t",
					Usage: "Delete via Tor hidden service",
				},
				cli.IntFlag{
					Name:  "tor-port",
					Usage: "Use a different port to connect to Tor",
					Value: 9150,
				},
			},
		},
		{
			Name:   "update",
			Usage:  "Update (overwrite) a post",
			Action: cmdUpdate,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "tor, t",
					Usage: "Update via Tor hidden service",
				},
				cli.IntFlag{
					Name:  "tor-port",
					Usage: "Use a different port to connect to Tor",
					Value: 9150,
				},
				cli.BoolFlag{
					Name:  "code",
					Usage: "Specifies this post is code",
				},
				cli.StringFlag{
					Name:  "font",
					Usage: "Sets post font to given value",
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
		{
			Name:  "add",
			Usage: "Add an existing post locally",
			Description: `A way to add an existing post to your local store for easy editing later.
			
   This requires a post ID (from https://write.as/[ID]) and an Edit Token
   (exported from another Write.as client, such as the Android app).
`,
			Action: cmdAdd,
		},
		{
			Name:   "list",
			Usage:  "List local posts",
			Action: cmdList,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "id",
					Usage: "Show list with post IDs (default)",
				},
				cli.BoolFlag{
					Name:  "url",
					Usage: "Show list with URLs",
				},
			},
		},
	}

	cli.CommandHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   writeas {{.Name}}{{if .Flags}} [command options]{{end}} [arguments...]{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if .Flags}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{ end }}
`

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
		fmt.Fprintf(os.Stderr, "writeas: %s\n", err)
		os.Exit(1)
	}
}

func handlePost(fullPost []byte, c *cli.Context) error {
	tor := c.Bool("tor") || c.Bool("t")
	if c.Int("tor-port") != 0 {
		torPort = c.Int("tor-port")
	}
	if tor {
		fmt.Println("Posting to hidden service...")
	} else {
		fmt.Println("Posting...")
	}

	return DoPost(fullPost, c.String("font"), false, tor, c.Bool("code"))
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
func DoFetch(friendlyID string, tor bool) {
	path := friendlyID

	urlStr, client := client(true, tor, path, "")

	r, _ := http.NewRequest("GET", urlStr, nil)
	r.Header.Add("User-Agent", "writeas-cli v"+version)

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

// DoPost creates a Write.as post, returning an error if it was
// unsuccessful.
func DoPost(post []byte, font string, encrypt, tor, code bool) error {
	data := url.Values{}
	data.Set("w", string(post))
	if encrypt {
		data.Add("e", "")
	}
	data.Add("font", getFont(code, font))

	urlStr, client := client(false, tor, "", "")

	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	r.Header.Add("User-Agent", "writeas-cli v"+version)
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

		// Copy URL to clipboard
		err = clipboard.WriteAll(string(url))
		if err != nil {
			fmt.Fprintf(os.Stderr, "writeas: Didn't copy to clipboard: %s\n", err)
		} else {
			fmt.Println("Copied to clipboard.")
		}

		// Output URL
		fmt.Printf("%s\n", url)
	} else {
		return fmt.Errorf("Unable to post: %s", resp.Status)
	}

	return nil
}

// DoUpdate updates the given post on Write.as.
func DoUpdate(post []byte, friendlyID, token, font string, tor, code bool) {
	urlStr, client := client(false, tor, friendlyID, fmt.Sprintf("t=%s", token))

	data := url.Values{}
	data.Set("w", string(post))

	if code || font != "" {
		// Only update font if explicitly changed
		data.Add("font", getFont(code, font))
	}

	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	r.Header.Add("User-Agent", "writeas-cli v"+version)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	check(err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		if tor {
			fmt.Println("Post updated via hidden service.")
		} else {
			fmt.Println("Post updated.")
		}
	} else {
		if debug {
			fmt.Printf("Problem updating: %s\n", resp.Status)
		} else {
			fmt.Printf("Post doesn't exist, or bad edit token given.\n")
		}
	}
}

// DoDelete deletes the given post on Write.as.
func DoDelete(friendlyID, token string, tor bool) {
	urlStr, client := client(false, tor, friendlyID, fmt.Sprintf("t=%s", token))

	r, _ := http.NewRequest("DELETE", urlStr, nil)
	r.Header.Add("User-Agent", "writeas-cli v"+version)
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
		removePost(friendlyID)
	} else {
		if debug {
			fmt.Printf("Problem deleting: %s\n", resp.Status)
		} else {
			fmt.Printf("Post doesn't exist, or bad edit token given.\n")
		}
	}
}
