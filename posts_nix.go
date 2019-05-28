// +build !windows

package writeascli

import (
	"fmt"
	"os/exec"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	dataDirName = ".writeas"
	noEditorErr = "Couldn't find default editor. Try setting $EDITOR environment variable in ~/.profile"
)

func parentDataDir() string {
	dir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	return dir
}

func editPostCmd(fname string) *exec.Cmd {
	editor := GetConfiguredEditor()
	if editor == "" {
		// Fall back to default editor
		path, err := exec.LookPath("vim")
		if err != nil {
			path, err = exec.LookPath("nano")
			if err != nil {
				return nil
			}
		}
		editor = path
	}
	return exec.Command(editor, fname)
}

func messageRetryCompose(fname string) string {
	return fmt.Sprintf("To retry this post, run:\n  cat %s | writeas", fname)
}
