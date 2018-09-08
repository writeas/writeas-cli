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
