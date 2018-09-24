package main

import (
	"encoding/json"
	"github.com/writeas/writeas-cli/fileutils"
	"go.code.as/writeas.v2"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"path/filepath"
)

const (
	userConfigFile = "config.ini"
	userFile       = "user.json"
)

type (
	APIConfig struct {
	}

	PostsConfig struct {
		Directory string `ini:"directory"`
	}

	UserConfig struct {
		API   APIConfig   `ini:"api"`
		Posts PostsConfig `ini:"posts"`
	}
)

func loadConfig() (*UserConfig, error) {
	// TODO: load config to var shared across app
	cfg, err := ini.LooseLoad(filepath.Join(userDataDir(), userConfigFile))
	if err != nil {
		return nil, err
	}

	// Parse INI file
	uc := &UserConfig{}
	err = cfg.MapTo(uc)
	if err != nil {
		return nil, err
	}
	return uc, nil
}

func saveConfig(uc *UserConfig) error {
	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, uc)
	if err != nil {
		return err
	}

	return cfg.SaveTo(filepath.Join(userDataDir(), userConfigFile))
}

func loadUser() (*writeas.AuthUser, error) {
	fname := filepath.Join(userDataDir(), userFile)
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

func saveUser(u *writeas.AuthUser) error {
	// Marshal struct into pretty-printed JSON
	userJSON, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return err
	}

	// Save file
	err = ioutil.WriteFile(filepath.Join(userDataDir(), userFile), userJSON, 0600)
	if err != nil {
		return err
	}
	return nil
}
