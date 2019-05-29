package config

import (
	"github.com/cloudfoundry/jibber_jabber"
	"github.com/writeas/writeas-cli/log"
	cli "gopkg.in/urfave/cli.v1"
)

// Application constants.
const (
	Version          = "1.99-dev"
	defaultUserAgent = "writeas-cli v" + Version
	// Defaults for posts on Write.as.
	DefaultFont    = PostFontMono
	WriteasBaseURL = "https://write.as"
	DevBaseURL     = "https://development.write.as"
	TorBaseURL     = "http://writeas7pm7rcdqg.onion"
)

func UserAgent(c *cli.Context) string {
	ua := c.String("user-agent")
	if ua == "" {
		return defaultUserAgent
	}
	return ua + " (" + defaultUserAgent + ")"
}

func IsTor(c *cli.Context) bool {
	return c.Bool("tor") || c.Bool("t")
}

func Language(c *cli.Context, auto bool) string {
	if l := c.String("lang"); l != "" {
		return l
	}
	if !auto {
		return ""
	}
	// Automatically detect language
	l, err := jibber_jabber.DetectLanguage()
	if err != nil {
		log.Info(c, "Language detection failed: %s", err)
		return ""
	}
	return l
}

func Collection(c *cli.Context) string {
	if coll := c.String("c"); coll != "" {
		return coll
	}
	if coll := c.String("b"); coll != "" {
		return coll
	}
	return ""
}
