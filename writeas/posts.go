package main

import (
	"fmt"
	"github.com/writeas/writeas-cli/utils"
	"os"
	"path/filepath"
	"strings"
)

const (
	POSTS_FILE = "posts.psv"
	SEPARATOR  = `|`
)

type Post struct {
	ID        string
	EditToken string
}

func userDataDir() string {
	return filepath.Join(parentDataDir(), DATA_DIR_NAME)
}

func dataDirExists() bool {
	return fileutils.Exists(userDataDir())
}

func createDataDir() {
	err := os.Mkdir(userDataDir(), 0700)
	if err != nil && DEBUG {
		panic(err)
	}
}

func addPost(id, token string) {
	f, err := os.OpenFile(filepath.Join(userDataDir(), POSTS_FILE), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		if DEBUG {
			panic(err)
		} else {
			return
		}
	}
	defer f.Close()

	l := fmt.Sprintf("%s%s%s\n", id, SEPARATOR, token)

	if _, err = f.WriteString(l); err != nil && DEBUG {
		panic(err)
	}
}

func tokenFromID(id string) string {
	post := fileutils.FindLine(filepath.Join(userDataDir(), POSTS_FILE), id)
	if post == "" {
		return ""
	}

	parts := strings.Split(post, SEPARATOR)
	if len(parts) < 2 {
		return ""
	}

	return parts[1]
}

func removePost(id string) {
	fileutils.RemoveLine(filepath.Join(userDataDir(), POSTS_FILE), id)
}

func getPosts() *[]Post {
	lines := fileutils.ReadData(filepath.Join(userDataDir(), POSTS_FILE))

	posts := []Post{}
	parts := make([]string, 2)

	for _, l := range *lines {
		parts = strings.Split(l, SEPARATOR)
		if len(parts) < 2 {
			continue
		}
		posts = append(posts, Post{ID: parts[0], EditToken: parts[1]})
	}

	return &posts

}
