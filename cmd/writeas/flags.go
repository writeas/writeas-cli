package main

import (
	"gopkg.in/urfave/cli.v1"
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:   "user, u",
		Hidden: true,
		Value:  "user",
	},
}
