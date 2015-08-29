// +build windows

package main

import (
	"os"
	"os/exec"
)

const (
	DATA_DIR_NAME = "Write.as"
	NO_EDITOR_ERR = "Error getting default editor. You shouldn't see this, so let us know you did: hello@write.as"
)

func parentDataDir() string {
	return os.Getenv("APPDATA")
}

func editPostCmd(fname string) *exec.Cmd {
	return exec.Command(fname)
}
