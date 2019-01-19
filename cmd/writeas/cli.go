package main

import (
	"bufio"
	"go.code.as/writeas.v2"
	"gopkg.in/urfave/cli.v1"
	"io"
	"log"
	"os"
)

// API constants for communicating with Write.as.
const (
	writeasBaseURL = "https://write.as"
	devBaseURL     = "https://development.write.as"
	torBaseURL     = "http://writeas7pm7rcdqg.onion"
)

// Application constants.
const (
	version = "1.99-dev"
)

// Defaults for posts on Write.as.
const (
	defaultFont = PostFontMono
)

// Available flags for creating posts
var postFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "c, b",
		Usage: "Optional blog to post to",
		Value: "",
	},
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
	cli.BoolFlag{
		Name:  "md",
		Usage: "Returns post URL with Markdown enabled",
	},
	cli.BoolFlag{
		Name:  "verbose, v",
		Usage: "Make the operation more talkative",
	},
	cli.StringFlag{
		Name:  "font",
		Usage: "Sets post font to given value",
	},
	cli.StringFlag{
		Name:  "lang",
		Usage: "Sets post language to given ISO 639-1 language code",
		Value: "",
	},
	cli.StringFlag{
		Name:  "user-agent",
		Usage: "Sets the User-Agent for API requests",
		Value: "",
	},
}

func main() {
	initialize()

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the version",
	}

	// Run the app
	app := cli.NewApp()
	app.Name = "writeas"
	app.Version = version
	app.Usage = "Publish text quickly"
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
   value set to the WRITAS_EDITOR or EDITOR environment variable, or vim, or
   finally nano.

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
			Name:   "publish",
			Usage:  "Publish a file to Write.as",
			Action: cmdPublish,
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
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "Make the operation more talkative",
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
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "Make the operation more talkative",
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
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "Make the operation more talkative",
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
					Name:  "md",
					Usage: "Use with --url to return URLs with Markdown enabled",
				},
				cli.BoolFlag{
					Name:  "url",
					Usage: "Show list with URLs",
				},
			},
		},
		{
			Name:   "fetch",
			Usage:  "Fetch authenticated user's Write.as posts",
			Action: cmdPull,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "tor, t",
					Usage: "Authenticate via Tor hidden service",
				},
				cli.IntFlag{
					Name:  "tor-port",
					Usage: "Use a different port to connect to Tor",
					Value: 9150,
				},
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "Make the operation more talkative",
				},
			},
		},
		{
			Name:   "auth",
			Usage:  "Authenticate with Write.as",
			Action: cmdAuth,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "tor, t",
					Usage: "Authenticate via Tor hidden service",
				},
				cli.IntFlag{
					Name:  "tor-port",
					Usage: "Use a different port to connect to Tor",
					Value: 9150,
				},
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "Make the operation more talkative",
				},
			},
		},
		{
			Name:   "logout",
			Usage:  "Log out of Write.as",
			Action: cmdLogOut,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "tor, t",
					Usage: "Authenticate via Tor hidden service",
				},
				cli.IntFlag{
					Name:  "tor-port",
					Usage: "Use a different port to connect to Tor",
					Value: 9150,
				},
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "Make the operation more talkative",
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

func handlePost(fullPost []byte, c *cli.Context) (*writeas.Post, error) {
	tor := isTor(c)
	if c.Int("tor-port") != 0 {
		torPort = c.Int("tor-port")
	}
	if tor {
		Info(c, "Posting to hidden service...")
	} else {
		Info(c, "Posting...")
	}

	return DoPost(c, fullPost, c.String("font"), false, tor, c.Bool("code"))
}
