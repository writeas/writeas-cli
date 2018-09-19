package main

import (
	"fmt"
	"github.com/writeas/go-writeas"
	"github.com/writeas/writeas-cli/fileutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	postsFile = "posts.psv"
	separator = `|`
)

// Post holds the basic authentication information for a Write.as post.
type Post struct {
	ID        string
	EditToken string
}

func userDataDir() string {
	return filepath.Join(parentDataDir(), dataDirName)
}

func dataDirExists() bool {
	return fileutils.Exists(userDataDir())
}

func createDataDir() {
	err := os.Mkdir(userDataDir(), 0700)
	if err != nil {
		if debug {
			panic(err)
		} else {
			Errorln("Error creating data directory: %s", err)
			return
		}
	}
}

func addPost(id, token string) error {
	f, err := os.OpenFile(filepath.Join(userDataDir(), postsFile), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("Error creating local posts list: %s", err)
	}
	defer f.Close()

	l := fmt.Sprintf("%s%s%s\n", id, separator, token)

	if _, err = f.WriteString(l); err != nil {
		return fmt.Errorf("Error writing to local posts list: %s", err)
	}

	return nil
}

func tokenFromID(id string) string {
	post := fileutils.FindLine(filepath.Join(userDataDir(), postsFile), id)
	if post == "" {
		return ""
	}

	parts := strings.Split(post, separator)
	if len(parts) < 2 {
		return ""
	}

	return parts[1]
}

func removePost(id string) {
	fileutils.RemoveLine(filepath.Join(userDataDir(), postsFile), id)
}

func getPosts() *[]Post {
	lines := fileutils.ReadData(filepath.Join(userDataDir(), postsFile))

	posts := []Post{}

	if lines != nil && len(*lines) > 0 {
		parts := make([]string, 2)
		for _, l := range *lines {
			parts = strings.Split(l, separator)
			if len(parts) < 2 {
				continue
			}
			posts = append(posts, Post{ID: parts[0], EditToken: parts[1]})
		}
	}

	return &posts
}

func composeNewPost() (string, *[]byte) {
	f, err := fileutils.TempFile(os.TempDir(), "WApost", "txt")
	if err != nil {
		if debug {
			panic(err)
		} else {
			Errorln("Error creating temp file: %s", err)
			return "", nil
		}
	}
	f.Close()

	cmd := editPostCmd(f.Name())
	if cmd == nil {
		os.Remove(f.Name())

		fmt.Println(noEditorErr)
		return "", nil
	}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Start(); err != nil {
		os.Remove(f.Name())

		if debug {
			panic(err)
		} else {
			Errorln("Error starting editor: %s", err)
			return "", nil
		}
	}

	// If something fails past this point, the temporary post file won't be
	// removed automatically. Calling function should handle this.
	if err := cmd.Wait(); err != nil {
		if debug {
			panic(err)
		} else {
			Errorln("Editor finished with error: %s", err)
			return "", nil
		}
	}

	post, err := ioutil.ReadFile(f.Name())
	if err != nil {
		if debug {
			panic(err)
		} else {
			Errorln("Error reading post: %s", err)
			return "", nil
		}
	}
	return f.Name(), &post
}

func WritePost(postsDir string, p *writeas.Post) error {
	postFilename := p.ID
	collDir := ""
	if p.Collection != nil {
		postFilename = p.Slug
		collDir = p.Collection.Alias
	}
	postFilename += postFileExt

	txtFile := p.Content
	if p.Title != "" {
		txtFile = "# " + p.Title + "\n\n" + txtFile
	}
	return ioutil.WriteFile(filepath.Join(postsDir, collDir, postFilename), []byte(txtFile), 0644)
}
