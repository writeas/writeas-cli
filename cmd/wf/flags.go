package main

import (
	"gopkg.in/urfave/cli.v1"
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:  "host, H",
		Usage: "Operate against a custom hostname",
	},
	cli.StringFlag{
		Name:  "user, u",
		Usage: "Use authenticated user, other than default",
	},
}
