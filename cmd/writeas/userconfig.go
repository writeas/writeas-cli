package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/writeas/writeas-cli/fileutils"
	writeas "go.code.as/writeas.v2"
	ini "gopkg.in/ini.v1"
)

const (
	userConfigFile = "config.ini"
	userFile       = "user.json"
)

type (
	APIConfig struct {
	}

	UserConfig struct {
		API APIConfig `ini:"api"`
	}
)

func loadConfig() (*UserConfig, error) {
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
