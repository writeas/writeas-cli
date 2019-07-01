package api

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	writeas "github.com/writeas/go-writeas/v2"
	"github.com/writeas/writeas-cli/config"
	"github.com/writeas/writeas-cli/fileutils"
	"github.com/writeas/writeas-cli/log"
	cli "gopkg.in/urfave/cli.v1"
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

// RemotePost holds addition information about published
// posts
type RemotePost struct {
	Post
	Title,
	Excerpt,
	Slug,
	Collection,
	EditToken string
	Synced  bool
	Updated time.Time
}

func AddPost(c *cli.Context, id, token string) error {
	hostDir, err := config.HostDirectory(c)
	if err != nil {
		return fmt.Errorf("Error checking for host directory: %v", err)
	}
	f, err := os.OpenFile(filepath.Join(config.UserDataDir(c.App.ExtraInfo()["configDir"]), hostDir, postsFile), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
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

// ClaimPost adds a local post to the authenticated user's account and deletes
// the local reference
func ClaimPosts(c *cli.Context, localPosts *[]Post) (*[]writeas.ClaimPostResult, error) {
	cl, err := newClient(c, true)
	if err != nil {
		return nil, err
	}
	postsToClaim := make([]writeas.OwnedPostParams, len(*localPosts))
	for i, post := range *localPosts {
		postsToClaim[i] = writeas.OwnedPostParams{
			ID:    post.ID,
			Token: post.EditToken,
		}
	}

	return cl.ClaimPosts(&postsToClaim)
}

func TokenFromID(c *cli.Context, id string) string {
	hostDir, _ := config.HostDirectory(c)
	post := fileutils.FindLine(filepath.Join(config.UserDataDir(c.App.ExtraInfo()["configDir"]), hostDir, postsFile), id)
	if post == "" {
		return ""
	}

	parts := strings.Split(post, separator)
	if len(parts) < 2 {
		return ""
	}

	return parts[1]
}

func RemovePost(c *cli.Context, id string) {
	hostDir, _ := config.HostDirectory(c)
	fullPath := filepath.Join(config.UserDataDir(c.App.ExtraInfo()["configDir"]), hostDir, postsFile)
	fileutils.RemoveLine(fullPath, id)
}

func GetPosts(c *cli.Context) *[]Post {
	hostDir, _ := config.HostDirectory(c)
	lines := fileutils.ReadData(filepath.Join(config.UserDataDir(c.App.ExtraInfo()["configDir"]), hostDir, postsFile))

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

func GetUserPosts(c *cli.Context, draftsOnly bool) ([]RemotePost, error) {
	waposts, err := DoFetchPosts(c)
	if err != nil {
		return nil, err
	}

	if len(waposts) == 0 {
		return nil, nil
	}

	posts := []RemotePost{}
	for _, p := range waposts {
		if draftsOnly && p.Collection != nil {
			continue
		}
		post := RemotePost{
			Post: Post{
				ID:        p.ID,
				EditToken: p.Token,
			},
			Title:   p.Title,
			Excerpt: getExcerpt(p.Content),
			Slug:    p.Slug,
			Synced:  p.Slug != "",
			Updated: p.Updated,
		}
		if p.Collection != nil {
			post.Collection = p.Collection.Alias
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// getExcerpt takes in a content string and returns
// a concatenated version. limited to no more than
// two lines of 80 chars each. delimited by '...'
func getExcerpt(input string) string {
	length := len(input)

	if length <= 80 {
		return input
	} else if length < 160 {
		ln1, idx := trimToLength(input, 80)
		if idx == -1 {
			idx = 80
		}
		ln2, _ := trimToLength(input[idx:], 80)
		return ln1 + "\n" + ln2
	} else {
		excerpt := input[:158]
		ln1, idx := trimToLength(excerpt, 80)
		if idx == -1 {
			idx = 80
		}
		ln2, _ := trimToLength(excerpt[idx:], 80)
		return ln1 + "\n" + ln2 + "..."
	}
}

func trimToLength(in string, l int) (string, int) {
	c := []rune(in)
	spaceIdx := -1
	length := len(c)
	if length <= l {
		return in, spaceIdx
	}

	for i := l; i > 0; i-- {
		if c[i] == ' ' {
			spaceIdx = i
			break
		}
	}
	if spaceIdx > -1 {
		c = c[:spaceIdx]
	}
	return string(c), spaceIdx
}

func ComposeNewPost() (string, *[]byte) {
	f, err := fileutils.TempFile(os.TempDir(), "WApost", "txt")
	if err != nil {
		if config.Debug() {
			panic(err)
		} else {
			log.Errorln("Error creating temp file: %s", err)
			return "", nil
		}
	}
	f.Close()

	cmd := config.EditPostCmd(f.Name())
	if cmd == nil {
		os.Remove(f.Name())

		fmt.Println(config.NoEditorErr)
		return "", nil
	}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Start(); err != nil {
		os.Remove(f.Name())

		if config.Debug() {
			panic(err)
		} else {
			log.Errorln("Error starting editor: %s", err)
			return "", nil
		}
	}

	// If something fails past this point, the temporary post file won't be
	// removed automatically. Calling function should handle this.
	if err := cmd.Wait(); err != nil {
		if config.Debug() {
			panic(err)
		} else {
			log.Errorln("Editor finished with error: %s", err)
			return "", nil
		}
	}

	post, err := ioutil.ReadFile(f.Name())
	if err != nil {
		if config.Debug() {
			panic(err)
		} else {
			log.Errorln("Error reading post: %s", err)
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
	postFilename += PostFileExt

	txtFile := p.Content
	if p.Title != "" {
		txtFile = "# " + p.Title + "\n\n" + txtFile
	}
	return ioutil.WriteFile(filepath.Join(postsDir, collDir, postFilename), []byte(txtFile), 0644)
}

func ReadStdIn() []byte {
	numBytes, numChunks := int64(0), int64(0)
	r := bufio.NewReader(os.Stdin)
	fullPost := []byte{}
	buf := make([]byte, 0, 1024)
	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.ErrorlnQuit("Error reading from stdin: %v", err)
		}
		numChunks++
		numBytes += int64(len(buf))

		fullPost = append(fullPost, buf...)
		if err != nil && err != io.EOF {
			log.ErrorlnQuit("Error appending to end of post: %v", err)
		}
	}

	return fullPost
}
