package main

import (
	"github.com/cloudfoundry/jibber_jabber"
	cli "gopkg.in/urfave/cli.v1"
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

func language(c *cli.Context, auto bool) string {
	if l := c.String("lang"); l != "" {
		return l
	}
	if !auto {
		return ""
	}
	// Automatically detect language
	l, err := jibber_jabber.DetectLanguage()
	if err != nil {
		Info(c, "Language detection failed: %s", err)
		return ""
	}
	return l
}

func collection(c *cli.Context) string {
	if coll := c.String("c"); coll != "" {
		return coll
	}
	if coll := c.String("b"); coll != "" {
		return coll
	}
	return ""
}
