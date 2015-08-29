package main

import (
	"fmt"
	"github.com/writeas/writeas-cli/utils"
	"io/ioutil"
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
	if err != nil {
		if DEBUG {
			panic(err)
		} else {
			fmt.Printf("Error creating data directory: %s\n", err)
			return
		}
	}
}

func addPost(id, token string) {
	f, err := os.OpenFile(filepath.Join(userDataDir(), POSTS_FILE), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		if DEBUG {
			panic(err)
		} else {
			fmt.Printf("Error creating local posts list: %s\n", err)
			return
		}
	}
	defer f.Close()

	l := fmt.Sprintf("%s%s%s\n", id, SEPARATOR, token)

	if _, err = f.WriteString(l); err != nil {
		if DEBUG {
			panic(err)
		} else {
			fmt.Printf("Error writing to local posts list: %s\n", err)
			return
		}
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

func composeNewPost() (string, *[]byte) {
	f, err := fileutils.TempFile(os.TempDir(), "WApost", "txt")
	if err != nil {
		if DEBUG {
			panic(err)
		} else {
			fmt.Printf("Error creating temp file: %s\n", err)
			return "", nil
		}
	}
	f.Close()

	cmd := editPostCmd(f.Name())
	if cmd == nil {
		os.Remove(f.Name())

		fmt.Println(NO_EDITOR_ERR)
		return "", nil
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		os.Remove(f.Name())

		if DEBUG {
			panic(err)
		} else {
			fmt.Printf("Error starting editor: %s\n", err)
			return "", nil
		}
	}

	// If something fails past this point, the temporary post file won't be
	// removed automatically. Calling function should handle this.
	if err := cmd.Wait(); err != nil {
		if DEBUG {
			panic(err)
		} else {
			fmt.Printf("Editor finished with error: %s\n", err)
			return "", nil
		}
	}

	post, err := ioutil.ReadFile(f.Name())
	if err != nil {
		if DEBUG {
			panic(err)
		} else {
			fmt.Printf("Error reading post: %s\n", err)
			return "", nil
		}
	}
	return f.Name(), &post
}
