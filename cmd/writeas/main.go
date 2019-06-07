package main

import (
	"os"

	"github.com/writeas/writeas-cli/api"
	"github.com/writeas/writeas-cli/commands"
	"github.com/writeas/writeas-cli/config"
	"github.com/writeas/writeas-cli/log"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	initialize(appInfo["configDir"])
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the version",
	}

	// Run the app
	app := cli.NewApp()
	app.Name = "writeas"
	app.Version = config.Version
	app.Usage = "Publish text quickly"
	app.Authors = []cli.Author{
		{
			Name:  "Write.as",
			Email: "hello@write.as",
		},
	}
	app.ExtraInfo = func() map[string]string {
		return appInfo
	}
	app.Action = commands.CmdPost
	app.Flags = config.PostFlags
	app.Commands = []cli.Command{
		{
			Name:   "post",
			Usage:  "Alias for default action: create post from stdin",
			Action: commands.CmdPost,
			Flags:  config.PostFlags,
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
   value set to the WRITEAS_EDITOR or EDITOR environment variable, or vim, or
   finally nano.

   Use the --code flag to indicate that the post should use syntax 
   highlighting. Or use the --font [value] argument to set the post's 
   appearance, where [value] is mono, monospace (default), wrap (monospace 
   font with word wrapping), serif, or sans.
   
   If posting fails for any reason, 'writeas' will show you the temporary file
   location and how to pipe it to 'writeas' to retry.`,
			Action: commands.CmdNew,
			Flags:  config.PostFlags,
		},
		{
			Name:   "publish",
			Usage:  "Publish a file to Write.as",
			Action: commands.CmdPublish,
			Flags:  config.PostFlags,
		},
		{
			Name:   "delete",
			Usage:  "Delete a post",
			Action: commands.CmdDelete,
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
			Action: commands.CmdUpdate,
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
			Action: commands.CmdGet,
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
			Action: commands.CmdAdd,
		},
		{
			Name:        "posts",
			Usage:       "List all of your posts",
			Description: "This will list only local posts when not currently authenticated. To list remote posts as well, first run: writeas auth <username>.",
			Action:      commands.CmdListPosts,
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
		}, {
			Name:   "blogs",
			Usage:  "List blogs",
			Action: commands.CmdCollections,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "url",
					Usage: "Show list with URLs",
				},
			},
		},
		{
			Name:   "fetch",
			Usage:  "Fetch authenticated user's Write.as posts",
			Action: api.CmdPull,
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
			Action: commands.CmdAuth,
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
			Action: commands.CmdLogOut,
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

func initialize(dataDirName string) {
	// Ensure we have a data directory to use
	if !config.DataDirExists(dataDirName) {
		err := config.CreateDataDir(dataDirName)
		if err != nil {
			if config.Debug() {
				panic(err)
			} else {
				log.Errorln("Error creating data directory: %s", err)
				return
			}
		}
	}
}
