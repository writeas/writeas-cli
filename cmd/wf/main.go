package main

import (
	"os"

	"github.com/writeas/writeas-cli/commands"
	"github.com/writeas/writeas-cli/config"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	appInfo := map[string]string{
		"configDir": configDir,
		"version":   "1.0",
	}
	config.DirMustExist(config.UserDataDir(appInfo["configDir"]))
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the version",
	}

	// Run the app
	app := cli.NewApp()
	app.Name = "wf"
	app.Version = appInfo["version"]
	app.Usage = "Publish to any WriteFreely instance from the command-line."
	// TODO: who is the author? the contributors? link to GH?
	app.Authors = []cli.Author{
		{
			Name:  "Write.as",
			Email: "hello@write.as",
		},
	}
	app.ExtraInfo = func() map[string]string {
		return appInfo
	}
	app.Action = requireAuth(commands.CmdPost, "publish")
	app.Flags = append(config.PostFlags, flags...)
	app.Commands = []cli.Command{
		{
			Name:   "post",
			Usage:  "Alias for default action: create post from stdin",
			Action: requireAuth(commands.CmdPost, "publish"),
			Flags:  config.PostFlags,
			Description: `Create a new post on WriteFreely from stdin.

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
   
   If posting fails for any reason, 'wf' will show you the temporary file
   location and how to pipe it to 'wf' to retry.`,
			Action: requireAuth(commands.CmdNew, "publish"),
			Flags:  config.PostFlags,
		},
		{
			Name:   "publish",
			Usage:  "Publish a file",
			Action: requireAuth(commands.CmdPublish, "publish"),
			Flags:  config.PostFlags,
		},
		{
			Name:   "delete",
			Usage:  "Delete a post",
			Action: requireAuth(commands.CmdDelete, "delete"),
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
			Action: requireAuth(commands.CmdUpdate, "update"),
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
			Name:        "posts",
			Usage:       "List all of your posts",
			Description: "This will list only local posts.",
			Action:      requireAuth(commands.CmdListPosts, "posts"),
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
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "Show verbose post listing, including Edit Tokens",
				},
			},
		}, {
			Name:   "blogs",
			Usage:  "List blogs",
			Action: requireAuth(commands.CmdCollections, "blogs"),
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
					Name:  "url",
					Usage: "Show list with URLs",
				},
			},
		}, {
			Name:        "claim",
			Usage:       "Claim local unsynced posts",
			Action:      requireAuth(commands.CmdClaim, "claim"),
			Description: "This will claim any unsynced posts local to this machine. To see which posts these are run: wf posts.",
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
		}, {
			Name:   "auth",
			Usage:  "Authenticate with a WriteFreely instance",
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
			Usage:  "Log out of a WriteFreely instance",
			Action: requireAuth(cmdLogOut, "logout"),
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
   wf {{.Name}}{{if .Flags}} [command options]{{end}} [arguments...]{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if .Flags}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{ end }}
`
	app.Run(os.Args)
}
