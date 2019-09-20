package config

import (
	"os"
	"path/filepath"

	"github.com/writeas/writeas-cli/fileutils"
	"github.com/writeas/writeas-cli/log"
)

// UserDataDir returns a platform specific directory under the user's home
// directory
func UserDataDir(dataDirName string) string {
	return filepath.Join(parentDataDir(), dataDirName)
}

func dataDirExists(dataDirName string) bool {
	return fileutils.Exists(dataDirName)
}

func createDataDir(dataDirName string) error {
	return os.Mkdir(dataDirName, 0700)
}

// DirMustExist checks for a directory, creates it if not found and either
// panics or logs and error depending on the status of Debug
func DirMustExist(dataDirName string) {
	// Ensure we have a data directory to use
	if !dataDirExists(dataDirName) {
		err := createDataDir(dataDirName)
		if err != nil {
			if Debug() {
				panic(err)
			} else {
				log.Errorln("Error creating data directory: %s", err)
				return
			}
		}
	}
}
