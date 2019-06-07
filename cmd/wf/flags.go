package main

import (
	"gopkg.in/urfave/cli.v1"
)

var globalFlags = []cli.Flag{
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
		Name:  "verbose, v",
		Usage: "Make the operation more talkative",
	},
	cli.StringFlag{
		Name:  "user-agent",
		Usage: "Sets the User-Agent for API requests",
		Value: "",
	},
	cli.StringFlag{
		Name:  "host, H",
		Usage: "Operate against a custom hostname",
	},
}
