package main

import (
	"encoding/json"
	"github.com/writeas/writeas-cli/fileutils"
	"go.code.as/writeas.v2"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"fmt"
	"os"
	"strings"
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

// Prints all values of the given struct
// For subfields: the field name is separated with dots (ex: Posts.Directory=)
func printConfig(x interface{}, prefix string, showEmptyValues bool) {
	values := reflect.ValueOf(x)

	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}
	typ := values.Type()

	for i := 0; i < typ.NumField(); i++ {
		val  := values.Field(i)
		name := typ.Field(i).Name

		if prefix != "" {
			name = prefix + "." + name
		}
		if(val.Kind() == reflect.Struct) {
			printConfig(val.Interface(), name, showEmptyValues)
		} else {
			if showEmptyValues || val.Interface() != reflect.Zero(val.Type()).Interface() {
				fmt.Printf("%s=%v\n", name, val)
			}
		}
	}
}

// Get the value of a given field
// For subfields: the name should be separated with dots (ex: Posts.Directory)
func getConfigField(x interface{}, name string) (*reflect.Value, error) {
	path   := strings.Split(name, ".")
	values := reflect.ValueOf(x)

	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}
	for _, part := range path {
		values = values.FieldByName(part)

		if !values.IsValid() {
			err := fmt.Errorf("error: key does not contain a section: %v", name)
			return nil, err
		}
	}
	if values.Kind() == reflect.Struct {
		err := fmt.Errorf("error: key does not contain a section: %v", name)
		return nil, err
	}
	return &values, nil
}

// Opens an editor to modify the config file
func composeConfig() error {
	filename := filepath.Join(userDataDir(), userConfigFile)

	// Open the editor
	cmd := editPostCmd(filename)
	if cmd == nil {
		fmt.Println(noEditorErr)
		return nil
	}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Start(); err != nil {
		if debug {
			panic(err)
		} else {
			Errorln("Error starting editor: %s", err)
			return nil
		}
	}

	// Wait until the editor is closed
	if err := cmd.Wait(); err != nil {
		if debug {
			panic(err)
		} else {
			Errorln("Editor finished with error: %s", err)
			return nil
		}
	}

	// Check if the config file is valid
	_, err := loadConfig()
	if err != nil {
		if debug {
			panic(err)
		} else {
			Errorln("Error loading config: %s", err)
			return nil
		}
	}
	return nil
}
