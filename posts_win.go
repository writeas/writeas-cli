// +build windows

package writeascli

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	dataDirName = "Write.as"
	noEditorErr = "Error getting default editor. You shouldn't see this, so let us know you did: hello@write.as"
)

func parentDataDir() string {
	return os.Getenv("APPDATA")
}

func editPostCmd(fname string) *exec.Cmd {
	// NOTE this won't work if fname contains spaces.
	return exec.Command("cmd", "/C copy con "+fname)
}

func messageRetryCompose(fname string) string {
	return fmt.Sprintf("To retry this post, run:\n  type %s | writeas.exe", fname)
}
