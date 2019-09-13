package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/hashicorp/go-multierror"
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
		} else if !c.GlobalIsSet("host") && !c.GlobalIsSet("user") {
			// check for global configured pair host/user
			cfg, err := config.LoadConfig(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("Failed to load config from file: %v", err), 1)
				// set flags if found
			}
			// set flags if both were found in config
			if cfg.Default.Host != "" && cfg.Default.User != "" {
				err = c.GlobalSet("host", cfg.Default.Host)
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("Failed to set host from global config: %v", err), 1)
				}
				err = c.GlobalSet("user", cfg.Default.User)
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("Failed to set user from global config: %v", err), 1)
				}
			} else {
				num, err := totalUsersLoggedIn(c)
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("Failed to check for logged in users: %v", err), 1)
				} else if num > 0 {
					return cli.NewExitError("You are authenticated, but have no default user/host set. Supply -user and -host flags.", 1)
				}
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

// totalUsersLoggedIn checks for logged in users for any host
// it returns the number of users and an error if any
func totalUsersLoggedIn(c *cli.Context) (int, error) {
	path := config.UserDataDir(c.App.ExtraInfo()["configDir"])
	dir, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	contents, err := dir.Readdir(0)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, file := range contents {
		if file.IsDir() {
			subDir, err := os.Open(filepath.Join(path, file.Name()))
			if err != nil {
				return 0, err
			}
			subContents, err := subDir.Readdir(0)
			if err != nil {
				return 0, err
			}
			for _, subFile := range subContents {
				if subFile.IsDir() {
					if _, err := os.Stat(filepath.Join(path, file.Name(), subFile.Name(), "user.json")); err == nil {
						count++
					}
				}
			}
		}
	}
	return count, nil
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

func cmdAccounts(c *cli.Context) error {
	// get user config dir
	userDir := config.UserDataDir(c.App.ExtraInfo()["configDir"])
	// load defaults
	cfg, err := config.LoadConfig(userDir)
	if err != nil {
		return cli.NewExitError("Could not load default user configuration", 1)
	}
	defaultUser := cfg.Default.User
	defaultHost := cfg.Default.Host
	if parts := strings.Split(defaultHost, "://"); len(parts) > 1 {
		defaultHost = parts[1]
	}
	// get each host dir
	files, err := ioutil.ReadDir(userDir)
	if err != nil {
		return cli.NewExitError("Could not read user configuration directory", 1)
	}
	// accounts will be a slice of slices of string. the first string in
	// a subslice should always be the hostname
	accounts := [][]string{}
	for _, file := range files {
		if file.IsDir() {
			dirName := file.Name()
			// get each user in host dir
			users, err := usersFromDir(filepath.Join(userDir, dirName))
			if err != nil {
				log.Info(c, "Failed to get users from %s: %v", dirName, err)
				continue
			}
			if len(users) != 0 {
				// append the slice of users as a new slice in accounts w/ the host prepended
				accounts = append(accounts, append([]string{dirName}, users...))
			}
		}
	}

	// print out all logged in accounts
	tw := tabwriter.NewWriter(os.Stdout, 10, 2, 2, ' ', tabwriter.TabIndent)
	if len(accounts) == 0 {
		fmt.Fprintf(tw, "%s\t", "No authenticated accounts found.")
	}
	for _, userList := range accounts {
		host := userList[0]
		for _, username := range userList[1:] {
			if host == defaultHost && username == defaultUser {
				fmt.Fprintf(tw, "[%s]\t%s (default)\n", host, username)
				continue
			}
			fmt.Fprintf(tw, "[%s]\t%s\n", host, username)
		}
	}
	return tw.Flush()
}

func usersFromDir(path string) ([]string, error) {
	users := make([]string, 0, 4)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var errs error
	for _, file := range files {
		if file.IsDir() {
			_, err := os.Stat(filepath.Join(path, file.Name(), "user.json"))
			if err != nil {
				err = multierror.Append(errs, err)
				continue
			}
			users = append(users, file.Name())
		}
	}
	return users, errs
}
