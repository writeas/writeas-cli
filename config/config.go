package config

import (
	"os"
	"path/filepath"

	ini "gopkg.in/ini.v1"
)

const (
	UserConfigFile = "config.ini"
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

func LoadConfig(dataDir string) (*UserConfig, error) {
	// TODO: load config to var shared across app
	cfg, err := ini.LooseLoad(filepath.Join(dataDir, UserConfigFile))
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

func SaveConfig(dataDir string, uc *UserConfig) error {
	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, uc)
	if err != nil {
		return err
	}

	return cfg.SaveTo(filepath.Join(dataDir, UserConfigFile))
}

var editors = []string{"WRITEAS_EDITOR", "EDITOR"}

func GetConfiguredEditor() string {
	for _, v := range editors {
		if e := os.Getenv(v); e != "" {
			return e
		}
	}
	return ""
}
