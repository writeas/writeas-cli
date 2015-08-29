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
	// NOTE this won't work if fname contains spaces.
	return exec.Command("cmd", "/C start /WAIT "+fname)
}

func messageRetryCompose(fname string) string {
	return fmt.Sprintf("To retry this post, run:\n  type %s | writeas.exe", fname)
}
