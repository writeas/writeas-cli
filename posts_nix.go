// +build !windows

package main

import (
	"github.com/mitchellh/go-homedir"
)

const DATA_DIR_NAME = ".writeas"

func parentDataDir() string {
	dir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	return dir
}
