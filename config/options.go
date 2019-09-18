package config

import (
	"strings"

	"github.com/cloudfoundry/jibber_jabber"
	"github.com/writeas/writeas-cli/log"
	cli "gopkg.in/urfave/cli.v1"
)

// Application constants.
const (
	writeasUserAgent = "writeas-cli v"
	wfUserAgent      = "wf-cli v"
	// Defaults for posts on Write.as.
	DefaultFont    = PostFontMono
	WriteasBaseURL = "https://write.as"
	DevBaseURL     = "https://development.write.as"
	TorBaseURL     = "http://writeas7pm7rcdqg.onion"
	torPort        = 9150
)

func UserAgent(c *cli.Context) string {
	client := wfUserAgent
	if c.App.Name == "writeas" {
		client = writeasUserAgent
	}

	ua := c.String("user-agent")
	if ua == "" {
		return client + c.App.ExtraInfo()["version"]
	}
	return ua + " (" + client + c.App.ExtraInfo()["version"] + ")"
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

func TorURL(c *cli.Context) string {
	flagHost := c.String("host")
	if flagHost != "" && strings.HasSuffix(flagHost, "onion") {
		return flagHost
	}
	cfg, _ := LoadConfig(c.App.ExtraInfo()["configDir"])
	if cfg != nil && cfg.Default.Host != "" && strings.HasSuffix(cfg.Default.Host, "onion") {
		return cfg.Default.Host
	}
	return TorBaseURL
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

// HostDirectory returns the sub directory string for the host. Order of
// precedence is a host flag if any, then the configured default, if any
func HostDirectory(c *cli.Context) (string, error) {
	cfg, err := LoadConfig(UserDataDir(c.App.ExtraInfo()["configDir"]))
	if err != nil {
		return "", err
	}
	// flag takes precedence over defaults
	if hostFlag := c.GlobalString("host"); hostFlag != "" {
		if parts := strings.Split(hostFlag, "://"); len(parts) > 1 {
			return parts[1], nil
		}
		return hostFlag, nil
	}

	if cfg.Default.Host != "" && cfg.Default.User != "" {
		if parts := strings.Split(cfg.Default.Host, "://"); len(parts) > 1 {
			return parts[1], nil
		}
		return cfg.Default.Host, nil
	}

	return "", nil
}
