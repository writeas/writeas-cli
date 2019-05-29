// +build windows

package config

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	NoEditorErr = "Error getting default editor. You shouldn't see this, so let us know you did: hello@write.as"
)

func parentDataDir() string {
	return os.Getenv("APPDATA")
}

func EditPostCmd(fname string) *exec.Cmd {
	// NOTE this won't work if fname contains spaces.
	return exec.Command("cmd", "/C copy con "+fname)
}

func MessageRetryCompose(fname string) string {
	return fmt.Sprintf("To retry this post, run:\n  type %s | writeas.exe", fname)
}
