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
	cli.StringFlag{
		Name:  "host, H",
		Usage: "Operate against a custom hostname",
	},
}
