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
		Font string `ini:"font"`
		Lang string `ini:"lang"`
		IsRTL bool `ini:"rtl"`
		Collection string `ini:"collection"`
	}

	UserConfig struct {
		API   APIConfig   `ini:"api"`
		Posts PostsConfig `ini:"posts"`
	}

	ConfigSingleton struct {
		uc *UserConfig
		err error
	}
)
var _instance *ConfigSingleton = nil

// Only load config file once
func loadConfig() (*UserConfig, error) {
	if _instance == nil {
		uc, err  := reloadConfig()
		_instance = &ConfigSingleton{uc, err}
	}
	return _instance.uc, _instance.err
}

func reloadConfig() (*UserConfig, error) {
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
