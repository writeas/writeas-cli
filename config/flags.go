package config

import (
	"gopkg.in/urfave/cli.v1"
)

// Available flags for creating posts
var PostFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "c, b",
		Usage: "Optional blog to post to",
		Value: "",
	},
	cli.BoolFlag{
		Name:  "insecure",
		Usage: "Send request insecurely.",
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
		Value: DefaultFont,
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
