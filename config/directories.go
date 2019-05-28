package config

import (
	"os"
	"path/filepath"

	"github.com/writeas/writeas-cli/fileutils"
)

func UserDataDir() string {
	return filepath.Join(parentDataDir(), dataDirName)
}

func DataDirExists() bool {
	return fileutils.Exists(UserDataDir())
}

func CreateDataDir() error {
	return os.Mkdir(UserDataDir(), 0700)
}
