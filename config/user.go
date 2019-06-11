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
	dir, err := userHostDir(c)
	if err != nil {
		return nil, err
	}
	username, err := currentUser(c)
	if err != nil {
		return nil, err
	}
	fname := filepath.Join(dir, username+".json")
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
	dir, err := userHostDir(c)
	if err != nil {
		return err
	}

	username, err := currentUser(c)
	if err != nil {
		return err
	}

	return fileutils.DeleteFile(filepath.Join(dir, username+".json"))
}

func SaveUser(c *cli.Context, u *writeas.AuthUser) error {
	// Marshal struct into pretty-printed JSON
	userJSON, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return err
	}

	dir, err := userHostDir(c)
	if err != nil {
		return err
	}
	DirMustExist(dir)
	// Save file
	err = ioutil.WriteFile(filepath.Join(dir, u.User.Username+".json"), userJSON, 0600)
	if err != nil {
		return err
	}
	return nil
}

// userHostDir returns the path to the user data directory with the host based
// subpath if the host flag is set
func userHostDir(c *cli.Context) (string, error) {
	dataDir := UserDataDir(c.App.ExtraInfo()["configDir"])
	hostDir, err := HostDirectory(c)
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, hostDir), nil
}

func currentUser(c *cli.Context) (string, error) {
	cfg, err := LoadConfig(UserDataDir(c.App.ExtraInfo()["configDir"]))
	if err != nil {
		return "", err
	}

	if c.GlobalString("user") != "" {
		return c.GlobalString("user"), nil
	}

	return cfg.Default.User, nil
}
