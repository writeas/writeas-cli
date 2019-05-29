package config

import (
	"os"
	"path/filepath"

	"github.com/writeas/writeas-cli/fileutils"
)

func UserDataDir(dataDirName string) string {
	return filepath.Join(parentDataDir(), dataDirName)
}

func DataDirExists(dataDirName string) bool {
	return fileutils.Exists(UserDataDir(dataDirName))
}

func CreateDataDir(dataDirName string) error {
	return os.Mkdir(UserDataDir(dataDirName), 0700)
}
