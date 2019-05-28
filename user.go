package writeascli

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/writeas/writeas-cli/fileutils"
	"go.code.as/writeas.v2"
)

const UserFile = "user.json"

func LoadUser(dataDir string) (*writeas.AuthUser, error) {
	fname := filepath.Join(dataDir, UserFile)
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

func SaveUser(dataDir string, u *writeas.AuthUser) error {
	// Marshal struct into pretty-printed JSON
	userJSON, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return err
	}

	// Save file
	err = ioutil.WriteFile(filepath.Join(dataDir, UserFile), userJSON, 0600)
	if err != nil {
		return err
	}
	return nil
}
