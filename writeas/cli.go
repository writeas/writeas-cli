package main

import (
	"bufio"
	"bytes"
	"errors"
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

const (
	apiUrl       = "https://write.as"
	hiddenApiUrl = "http://writeas7pm7rcdqg.onion"
	readApiUrl   = "https://write.as"
	VERSION      = "0.3"
)

// Defaults for posts on Write.as.
const (
	defaultFont = PostFontMono
)

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
	app.Version = VERSION
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
		fmt.Printf("%s\n", err)
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

func cmdPost(c *cli.Context) {
	err := handlePost(readStdIn(), c)
	check(err)
}

func cmdNew(c *cli.Context) {
	fname, p := composeNewPost()
	if p == nil {
		// Assume composeNewPost already told us what the error was. Abort now.
		os.Exit(1)
	}

	// Ensure we have something to post
	if len(*p) == 0 {
		// Clean up temporary post
		if fname != "" {
			os.Remove(fname)
		}

		fmt.Println("Empty post. Bye!")
		os.Exit(0)
	}

	err := handlePost(*p, c)
	if err != nil {
		fmt.Printf("Error posting: %s\n", err)
		fmt.Println(messageRetryCompose(fname))
		os.Exit(1)
	}

	// Clean up temporary post
	if fname != "" {
		os.Remove(fname)
	}
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
	if tor {
		fmt.Println("Deleting via hidden service...")
	} else {
		fmt.Println("Deleting...")
	}

	DoDelete(friendlyId, token, tor)
}

func cmdUpdate(c *cli.Context) {
	friendlyId := c.Args().Get(0)
	token := c.Args().Get(1)
	if friendlyId == "" {
		fmt.Println("usage: writeas update <postId> [<token>]")
		os.Exit(1)
	}

	if token == "" {
		// Search for the token locally
		token = tokenFromID(friendlyId)
		if token == "" {
			fmt.Println("Couldn't find an edit token locally. Did you create this post here?")
			fmt.Printf("If you have an edit token, use: writeas update %s <token>\n", friendlyId)
			os.Exit(1)
		}
	}

	// Read post body
	fullPost := readStdIn()

	tor := c.Bool("tor") || c.Bool("t")
	if c.Int("tor-port") != 0 {
		torPort = c.Int("tor-port")
	}
	if tor {
		fmt.Println("Updating via hidden service...")
	} else {
		fmt.Println("Updating...")
	}

	DoUpdate(fullPost, friendlyId, token, c.String("font"), tor, c.Bool("code"))
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
	if tor {
		fmt.Println("Getting via hidden service...")
	} else {
		fmt.Println("Getting...")
	}

	DoFetch(friendlyId, tor)
}

func cmdAdd(c *cli.Context) {
	friendlyId := c.Args().Get(0)
	token := c.Args().Get(1)
	if friendlyId == "" || token == "" {
		fmt.Println("usage: writeas add <postId> <token>")
		os.Exit(1)
	}

	addPost(friendlyId, token)
}

func cmdList(c *cli.Context) {
	urls := c.Bool("url")
	ids := c.Bool("id")

	var p Post
	posts := getPosts()
	for i := range *posts {
		p = (*posts)[len(*posts)-1-i]
		if ids || !urls {
			fmt.Printf("%s ", p.ID)
		}
		if urls {
			fmt.Printf("https://write.as/%s ", p.ID)
		}
		fmt.Print("\n")
	}
}

func client(read, tor bool, path, query string) (string, *http.Client) {
	var u *url.URL
	var client *http.Client
	if tor {
		u, _ = url.ParseRequestURI(hiddenApiUrl)
		u.Path = "/api/" + path
		client = torClient()
	} else {
		u, _ = url.ParseRequestURI(apiUrl)
		u.Path = "/api/" + path
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

func DoPost(post []byte, font string, encrypt, tor, code bool) error {
	data := url.Values{}
	data.Set("w", string(post))
	if encrypt {
		data.Add("e", "")
	}
	data.Add("font", getFont(code, font))

	urlStr, client := client(false, tor, "", "")

	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	r.Header.Add("User-Agent", "writeas-cli v"+VERSION)
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
			fmt.Printf("Didn't copy to clipboard: %s\n", err)
		} else {
			fmt.Println("Copied to clipboard.")
		}

		// Output URL
		fmt.Printf("%s\n", url)
	} else {
		return errors.New(fmt.Sprintf("Unable to post: %s", resp.Status))
	}

	return nil
}

func DoUpdate(post []byte, friendlyId, token, font string, tor, code bool) {
	urlStr, client := client(false, tor, friendlyId, fmt.Sprintf("t=%s", token))

	data := url.Values{}
	data.Set("w", string(post))

	if code || font != "" {
		// Only update font if explicitly changed
		data.Add("font", getFont(code, font))
	}

	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	r.Header.Add("User-Agent", "writeas-cli v"+VERSION)
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

func DoDelete(friendlyId, token string, tor bool) {
	urlStr, client := client(false, tor, friendlyId, fmt.Sprintf("t=%s", token))

	r, _ := http.NewRequest("DELETE", urlStr, nil)
	r.Header.Add("User-Agent", "writeas-cli v"+VERSION)
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
		if debug {
			fmt.Printf("Problem deleting: %s\n", resp.Status)
		} else {
			fmt.Printf("Post doesn't exist, or bad edit token given.\n")
		}
	}
}
