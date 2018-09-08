package main

import (
	"gopkg.in/urfave/cli.v1"
)

func userAgent(c *cli.Context) string {
	ua := c.String("user-agent")
	if ua == "" {
		return defaultUserAgent
	}
	return ua + " (" + defaultUserAgent + ")"
}

func isTor(c *cli.Context) bool {
	return c.Bool("tor") || c.Bool("t")
}
