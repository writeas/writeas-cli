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

	"github.com/writeas/writeas-cli/config"
	"github.com/writeas/writeas-cli/fileutils"
	"github.com/writeas/writeas-cli/log"
	writeas "go.code.as/writeas.v2"
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
	Updated time.Time
}

func AddPost(c *cli.Context, id, token string) error {
	f, err := os.OpenFile(filepath.Join(config.UserDataDir(c.App.ExtraInfo()["configDir"]), postsFile), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
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

func TokenFromID(c *cli.Context, id string) string {
	post := fileutils.FindLine(filepath.Join(config.UserDataDir(c.App.ExtraInfo()["configDir"]), postsFile), id)
	if post == "" {
		return ""
	}

	parts := strings.Split(post, separator)
	if len(parts) < 2 {
		return ""
	}

	return parts[1]
}

func removePost(path, id string) {
	fileutils.RemoveLine(filepath.Join(config.UserDataDir(path), postsFile), id)
}

func GetPosts(c *cli.Context) *[]Post {
	lines := fileutils.ReadData(filepath.Join(config.UserDataDir(c.App.ExtraInfo()["configDir"]), postsFile))

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

func GetUserPosts(c *cli.Context) ([]RemotePost, error) {
	waposts, err := DoFetchPosts(c)
	if err != nil {
		return nil, err
	}

	if len(waposts) == 0 {
		return nil, nil
	}

	posts := []RemotePost{}
	for _, p := range waposts {
		post := RemotePost{
			Post: Post{
				ID:        p.ID,
				EditToken: p.Token,
			},
			Title:   p.Title,
			Excerpt: getExcerpt(p.Content),
			Slug:    p.Slug,
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

func HandlePost(fullPost []byte, c *cli.Context) (*writeas.Post, error) {
	tor := config.IsTor(c)
	if c.Int("tor-port") != 0 {
		TorPort = c.Int("tor-port")
	}
	if tor {
		log.Info(c, "Posting to hidden service...")
	} else {
		log.Info(c, "Posting...")
	}

	return DoPost(c, fullPost, c.String("font"), false, tor, c.Bool("code"))
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
