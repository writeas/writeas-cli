package config

import (
	"net/url"

	"github.com/cloudfoundry/jibber_jabber"
	"github.com/writeas/writeas-cli/log"
	cli "gopkg.in/urfave/cli.v1"
)

// Application constants.
const (
	defaultUserAgent = "writeas-cli v"
	// Defaults for posts on Write.as.
	DefaultFont    = PostFontMono
	WriteasBaseURL = "https://write.as"
	DevBaseURL     = "https://development.write.as"
	TorBaseURL     = "http://writeas7pm7rcdqg.onion"
	torPort        = 9150
)

func UserAgent(c *cli.Context) string {
	ua := c.String("user-agent")
	if ua == "" {
		return defaultUserAgent + c.App.ExtraInfo()["version"]
	}
	return ua + " (" + defaultUserAgent + c.App.ExtraInfo()["version"] + ")"
}

func IsTor(c *cli.Context) bool {
	return c.Bool("tor") || c.Bool("t")
}

func TorPort(c *cli.Context) int {
	if c.IsSet("tor-port") && c.Int("tor-port") != 0 {
		return c.Int("tor-port")
	}
	return torPort
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
	u, _ := LoadUser(c)
	if u != nil {
		return u.User.Username
	}
	return ""
}

// HostDirectory returns the sub directory string for the host. Order of
// precedence is a host flag if any, then the configured default, if any
func HostDirectory(c *cli.Context) (string, error) {
	cfg, err := LoadConfig(UserDataDir(c.App.ExtraInfo()["configDir"]))
	if err != nil {
		return "", err
	}
	// flag takes precedence over defaults
	if hostFlag := c.GlobalString("host"); hostFlag != "" {
		u, err := url.Parse(hostFlag)
		if err != nil {
			return "", err
		}
		return u.Hostname(), nil
	}

	u, err := url.Parse(cfg.Default.Host)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}
