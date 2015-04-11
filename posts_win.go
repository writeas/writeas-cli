// +build windows

package main

import (
	"github.com/luisiturrios/gowin"
)

const DATA_DIR_NAME = "Write.as"

func parentDataDir() string {
	folders := gowin.ShellFolders{gowin.USER}
	return folders.AppData()
}
