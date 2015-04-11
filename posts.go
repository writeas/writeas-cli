package main

import (
	"fmt"
	"github.com/writeas/writeas-cli/utils"
	"os"
)

const (
	POSTS_FILE = "posts.psv"
	SEPARATOR  = `|`
)

func userDataDir() string {
	return fmt.Sprintf("%s/%s", parentDataDir(), DATA_DIR_NAME)
}

func dataDirExists() bool {
	return fileutils.Exists(userDataDir())
}

func createDataDir() {
	os.Mkdir(userDataDir(), 0700)
}

func addPost(id, token string) {
	f, err := os.OpenFile(userDataDir()+"/"+POSTS_FILE, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	l := fmt.Sprintf("%s%s%s\n", id, SEPARATOR, token)

	if _, err = f.WriteString(l); err != nil {
		panic(err)
	}
}
