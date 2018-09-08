package main

import (
	"gopkg.in/ini.v1"
	"path/filepath"
)

const (
	userConfigFile = "config.ini"
)

type (
	APIConfig struct {
		Token string `ini:"token"`
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
