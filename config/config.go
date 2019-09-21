package config

import (
	"os"
	"path/filepath"

	ini "gopkg.in/ini.v1"
)

const (
	// ConfigFile is the full filename for application configuration files
	ConfigFile = "config.ini"
)

type (
	// APIConfig is not currently used
	APIConfig struct {
	}

	// PostsConfig stores the directory for the user post cache
	PostsConfig struct {
		Directory string `ini:"directory"`
	}

	// DefaultConfig stores the default host and user to authenticate with
	DefaultConfig struct {
		Host string `ini:"host"`
		User string `ini:"user"`
	}

	// Config represents the entire base configuration
	Config struct {
		API     APIConfig     `ini:"api"`
		Default DefaultConfig `ini:"default"`
		Posts   PostsConfig   `ini:"posts"`
	}
)

func LoadConfig(dataDir string) (*Config, error) {
	// TODO: load config to var shared across app
	cfg, err := ini.LooseLoad(filepath.Join(dataDir, ConfigFile))
	if err != nil {
		return nil, err
	}

	// Parse INI file
	uc := &Config{}
	err = cfg.MapTo(uc)
	if err != nil {
		return nil, err
	}
	return uc, nil
}

func SaveConfig(dataDir string, uc *Config) error {
	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, uc)
	if err != nil {
		return err
	}

	return cfg.SaveTo(filepath.Join(dataDir, ConfigFile))
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
