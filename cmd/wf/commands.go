package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/writeas/writeas-cli/api"
	"github.com/writeas/writeas-cli/commands"
	"github.com/writeas/writeas-cli/config"
	"github.com/writeas/writeas-cli/executable"
	"github.com/writeas/writeas-cli/log"
	cli "gopkg.in/urfave/cli.v1"
)

func requireAuth(f cli.ActionFunc, action string) cli.ActionFunc {
	return func(c *cli.Context) error {
		// check for logged in users when host is provided without user
		if c.GlobalIsSet("host") && !c.GlobalIsSet("user") {
			// multiple users should display a list
			if num, users, err := usersLoggedIn(c); num > 1 && err == nil {
				return cli.NewExitError(fmt.Sprintf("Multiple logged in users, please use '-u' or '-user' to specify one of:\n%s", strings.Join(users, ", ")), 1)
			} else if num == 1 && err == nil {
				// single user found for host should be set as user flag so LoadUser can
				// succeed, and notify the client
				if err := c.GlobalSet("user", users[0]); err != nil {
					return cli.NewExitError(fmt.Sprintf("Failed to set user flag for only logged in user at host %s: %v", users[0], err), 1)
				}
				log.Info(c, "Host specified without user flag, using logged in user: %s\n", users[0])
			} else if err != nil {
				return cli.NewExitError(fmt.Sprintf("Failed to check for logged in users: %v", err), 1)
			}
		}
		u, err := config.LoadUser(c)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Couldn't load user: %v", err), 1)
		}
		if u == nil {
			return cli.NewExitError("You must be authenticated to "+action+".\nLog in first with: "+executable.Name()+" auth <username>", 1)
		}

		return f(c)
	}
}

// usersLoggedIn checks for logged in users for the set host flag
// it returns the number of users and a slice of usernames
func usersLoggedIn(c *cli.Context) (int, []string, error) {
	path, err := config.UserHostDir(c)
	if err != nil {
		return 0, nil, err
	}
	dir, err := os.Open(path)
	if err != nil {
		return 0, nil, err
	}
	contents, err := dir.Readdir(0)
	if err != nil {
		return 0, nil, err
	}
	var names []string
	for _, file := range contents {
		if file.IsDir() {
			// stat user.json
			if _, err := os.Stat(filepath.Join(path, file.Name(), "user.json")); err == nil {
				names = append(names, file.Name())
			}
		}
	}
	return len(names), names, nil
}

func cmdAuth(c *cli.Context) error {
	err := commands.CmdAuth(c)
	if err != nil {
		return err
	}

	// Get the username from the command, just like commands.CmdAuth does
	username := c.Args().Get(0)

	// Update config if this is user's first auth
	cfg, err := config.LoadConfig(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
	if err != nil {
		log.Errorln("Not saving config. Unable to load config: %s", err)
		return err
	}
	if cfg.Default.Host == "" && cfg.Default.User == "" {
		// This is user's first auth, so save defaults
		cfg.Default.Host = api.HostURL(c)
		cfg.Default.User = username
		err = config.SaveConfig(config.UserDataDir(c.App.ExtraInfo()["configDir"]), cfg)
		if err != nil {
			log.Errorln("Not saving config. Unable to save config: %s", err)
			return err
		}
		fmt.Printf("Set %s on %s as default account.\n", username, c.GlobalString("host"))
	}

	return nil
}

func cmdLogOut(c *cli.Context) error {
	err := commands.CmdLogOut(c)
	if err != nil {
		return err
	}

	// Remove this from config if it's the default account
	cfg, err := config.LoadConfig(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
	if err != nil {
		log.Errorln("Not updating config. Unable to load: %s", err)
		return err
	}
	username, err := config.CurrentUser(c)
	if err != nil {
		log.Errorln("Not updating config. Unable to load current user: %s", err)
		return err
	}
	reqHost := api.HostURL(c)
	if reqHost == "" {
		// No --host given, so we're using the default host
		reqHost = cfg.Default.Host
	}
	if cfg.Default.Host == reqHost && cfg.Default.User == username {
		// We're logging out of default username + host, so remove from config file
		cfg.Default.Host = ""
		cfg.Default.User = ""
		err = config.SaveConfig(config.UserDataDir(c.App.ExtraInfo()["configDir"]), cfg)
		if err != nil {
			log.Errorln("Not updating config. Unable to save config: %s", err)
			return err
		}
	}

	return nil
}
