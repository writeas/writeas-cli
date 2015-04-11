// +build windows

package main

import (
	"os"
)

const DATA_DIR_NAME = "Write.as"

func parentDataDir() string {
	return os.Getenv("APPDATA")
}
