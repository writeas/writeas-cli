package main

import (
	"gopkg.in/urfave/cli.v1"
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:  "host, H",
		Usage: "Use the given WriteFreely instance hostname",
	},
	cli.StringFlag{
		Name:  "user, u",
		Usage: "Use the given account username",
	},
}
