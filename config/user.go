package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	writeas "github.com/writeas/go-writeas/v2"
	"github.com/writeas/writeas-cli/fileutils"
	"gopkg.in/urfave/cli.v1"
)

func LoadUser(c *cli.Context) (*writeas.AuthUser, error) {
	dir, err := UserHostDir(c)
	if err != nil {
		return nil, err
	}
	DirMustExist(dir)
	username, err := CurrentUser(c)
	if err != nil {
		return nil, err
	}
	if username == "user" {
		username = ""
	}
	fname := filepath.Join(dir, username, "user.json")
	userJSON, err := ioutil.ReadFile(fname)
	if err != nil {
		if !fileutils.Exists(fname) {
			// Don't return a file-not-found error
			return nil, nil
		}
		return nil, err
	}

	// Parse JSON file
	u := &writeas.AuthUser{}
	err = json.Unmarshal(userJSON, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func DeleteUser(c *cli.Context) error {
	dir, err := UserHostDir(c)
	if err != nil {
		return err
	}

	username, err := CurrentUser(c)
	if err != nil {
		return err
	}

	if username == "user" {
		username = ""
	}

	// Delete user data
	err = fileutils.DeleteFile(filepath.Join(dir, username, "user.json"))
	if err != nil {
		return err
	}

	// Do additional cleanup in wf-cli
	if c.App.Name == "wf" {
		// Delete user dir if it's empty
		userEmpty, err := fileutils.IsEmpty(filepath.Join(dir, username))
		if err != nil {
			return err
		}
		if !userEmpty {
			return nil
		}
		err = fileutils.DeleteFile(filepath.Join(dir, username))
		if err != nil {
			return err
		}

		// Delete host dir if it's empty
		hostEmpty, err := fileutils.IsEmpty(dir)
		if err != nil {
			return err
		}
		if !hostEmpty {
			return nil
		}
		err = fileutils.DeleteFile(dir)
		if err != nil {
			return err
		}
	}

	return nil
}

func SaveUser(c *cli.Context, u *writeas.AuthUser) error {
	// Marshal struct into pretty-printed JSON
	userJSON, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return err
	}

	dir, err := UserHostDir(c)
	if err != nil {
		return err
	}
	// Save file
	username, err := CurrentUser(c)
	if err != nil {
		return err
	}
	if username != "user" {
		dir = filepath.Join(dir, u.User.Username)
	}
	DirMustExist(dir)
	err = ioutil.WriteFile(filepath.Join(dir, "user.json"), userJSON, 0600)
	if err != nil {
		return err
	}
	return nil
}

// UserHostDir returns the path to the user data directory with the host based
// subpath if the host flag is set
func UserHostDir(c *cli.Context) (string, error) {
	dataDir := UserDataDir(c.App.ExtraInfo()["configDir"])
	hostDir, err := HostDirectory(c)
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, hostDir), nil
}

// CurrentUser returns the username of the user taking action in the current
// cli.Context.
func CurrentUser(c *cli.Context) (string, error) {
	if c.App.Name == "writeas" {
		return "user", nil
	}
	// Use user flag value
	if c.GlobalString("user") != "" {
		return c.GlobalString("user"), nil
	}

	// Load host-level config, if host flag is set
	hostDir, err := UserHostDir(c)
	if err != nil {
		return "", err
	}
	cfg, err := LoadConfig(hostDir)
	if err != nil {
		return "", err
	}
	if cfg.Default.User == "" {
		// Load app-level config
		globalCFG, err := LoadConfig(UserDataDir(c.App.ExtraInfo()["configDir"]))
		if err != nil {
			return "", err
		}
		// only user global defaults when both are set and hosts match
		if globalCFG.Default.User != "" &&
			globalCFG.Default.Host != "" &&
			c.GlobalIsSet("host") &&
			globalCFG.Default.Host == c.GlobalString("host") {
			cfg = globalCFG
		}
	}

	return cfg.Default.User, nil
}
