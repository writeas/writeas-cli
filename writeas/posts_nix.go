// +build !windows

package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"os/exec"
)

const (
	DATA_DIR_NAME = ".writeas"
	NO_EDITOR_ERR = "Couldn't find default editor. Try setting $EDITOR environment variable in ~/.profile"
)

func parentDataDir() string {
	dir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	return dir
}

func editPostCmd(fname string) *exec.Cmd {
	editor := os.Getenv("EDITOR")
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
