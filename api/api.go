package api

import (
	"fmt"
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/writeas/web-core/posts"
	"github.com/writeas/writeas-cli/config"
	"github.com/writeas/writeas-cli/fileutils"
	"github.com/writeas/writeas-cli/log"
	writeas "go.code.as/writeas.v2"
	cli "gopkg.in/urfave/cli.v1"
)

func client(userAgent string, tor bool) *writeas.Client {
	var client *writeas.Client
	if tor {
		client = writeas.NewTorClient(TorPort)
	} else {
		if config.IsDev() {
			client = writeas.NewDevClient()
		} else {
			client = writeas.NewClient()
		}
	}
	client.UserAgent = userAgent

	return client
}

func NewClient(c *cli.Context, authRequired bool) (*writeas.Client, error) {
	var client *writeas.Client
	if config.IsTor(c) {
		client = writeas.NewTorClient(TorPort)
	} else {
		if config.IsDev() {
			client = writeas.NewDevClient()
		} else {
			client = writeas.NewClient()
		}
	}
	client.UserAgent = config.UserAgent(c)
	// TODO: load user into var shared across the app
	u, _ := config.LoadUser(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
	if u != nil {
		client.SetToken(u.AccessToken)
	} else if authRequired {
		return nil, fmt.Errorf("Not currently logged in. Authenticate with: writeas auth <username>")
	}

	return client, nil
}

// DoFetch retrieves the Write.as post with the given friendlyID,
// optionally via the Tor hidden service.
func DoFetch(friendlyID, ua string, tor bool) error {
	cl := client(ua, tor)

	p, err := cl.GetPost(friendlyID)
	if err != nil {
		return err
	}

	if p.Title != "" {
		fmt.Printf("# %s\n\n", string(p.Title))
	}
	fmt.Printf("%s\n", string(p.Content))
	return nil
}

// DoFetchPosts retrieves all remote posts for the
// authenticated user
func DoFetchPosts(c *cli.Context) ([]writeas.Post, error) {
	cl, err := NewClient(c, true)
	if err != nil {
		return nil, err
	}

	posts, err := cl.GetUserPosts()
	if err != nil {
		return nil, err
	}

	return *posts, nil
}

// DoPost creates a Write.as post, returning an error if it was
// unsuccessful.
func DoPost(c *cli.Context, post []byte, font string, encrypt, tor, code bool) (*writeas.Post, error) {
	cl, _ := NewClient(c, false)

	pp := &writeas.PostParams{
		Font:       config.GetFont(code, font),
		Collection: config.Collection(c),
	}
	pp.Title, pp.Content = posts.ExtractTitle(string(post))
	if lang := config.Language(c, true); lang != "" {
		pp.Language = &lang
	}
	p, err := cl.CreatePost(pp)
	if err != nil {
		return nil, fmt.Errorf("Unable to post: %v", err)
	}

	var url string
	if p.Collection != nil {
		url = p.Collection.URL + p.Slug
	} else {
		if tor {
			url = config.TorBaseURL
		} else if config.IsDev() {
			url = config.DevBaseURL
		} else {
			url = config.WriteasBaseURL
		}
		url += "/" + p.ID
		// Output URL in requested format
		if c.Bool("md") {
			url += ".md"
		}
	}

	if cl.Token() == "" {
		// Store post locally, since we're not authenticated
		AddPost(c, p.ID, p.Token)
	}

	// Copy URL to clipboard
	err = clipboard.WriteAll(string(url))
	if err != nil {
		log.Errorln("writeas: Didn't copy to clipboard: %s", err)
	} else {
		log.Info(c, "Copied to clipboard.")
	}

	// Output URL
	fmt.Printf("%s\n", url)

	return p, nil
}

// DoFetchCollections retrieves a list of the currently logged in users
// collections.
func DoFetchCollections(c *cli.Context) ([]RemoteColl, error) {
	cl, err := NewClient(c, true)
	if err != nil {
		if config.Debug() {
			log.ErrorlnQuit("could not create new client: %v", err)
		}
		return nil, fmt.Errorf("Couldn't create new client")
	}

	colls, err := cl.GetUserCollections()
	if err != nil {
		if config.Debug() {
			log.ErrorlnQuit("failed fetching user collections: %v", err)
		}
		return nil, fmt.Errorf("Couldn't get user collections")
	}

	out := make([]RemoteColl, len(*colls))

	for i, c := range *colls {
		coll := RemoteColl{
			Alias: c.Alias,
			Title: c.Title,
			URL:   c.URL,
		}
		out[i] = coll
	}

	return out, nil
}

// DoUpdate updates the given post on Write.as.
func DoUpdate(c *cli.Context, post []byte, friendlyID, token, font string, tor, code bool) error {
	cl, _ := NewClient(c, false)

	params := writeas.PostParams{}
	params.Title, params.Content = posts.ExtractTitle(string(post))
	if lang := config.Language(c, false); lang != "" {
		params.Language = &lang
	}
	if code || font != "" {
		params.Font = config.GetFont(code, font)
	}

	_, err := cl.UpdatePost(friendlyID, token, &params)
	if err != nil {
		if config.Debug() {
			log.ErrorlnQuit("Problem updating: %v", err)
		}
		return fmt.Errorf("Post doesn't exist, or bad edit token given.")
	}

	if tor {
		log.Info(c, "Post updated via hidden service.")
	} else {
		log.Info(c, "Post updated.")
	}
	return nil
}

// DoDelete deletes the given post on Write.as, and removes any local references
func DoDelete(c *cli.Context, friendlyID, token string, tor bool) error {
	cl, _ := NewClient(c, false)

	err := cl.DeletePost(friendlyID, token)
	if err != nil {
		if config.Debug() {
			log.ErrorlnQuit("Problem deleting: %v", err)
		}
		return fmt.Errorf("Post doesn't exist, or bad edit token given.")
	}

	if tor {
		log.Info(c, "Post deleted from hidden service.")
	} else {
		log.Info(c, "Post deleted.")
	}
	removePost(c.App.ExtraInfo()["configDir"], friendlyID)

	return nil
}

func DoLogIn(c *cli.Context, username, password string) error {
	cl := client(config.UserAgent(c), config.IsTor(c))

	u, err := cl.LogIn(username, password)
	if err != nil {
		if config.Debug() {
			log.ErrorlnQuit("Problem logging in: %v", err)
		}
		return err
	}

	err = config.SaveUser(config.UserDataDir(c.App.ExtraInfo()["configDir"]), u)
	if err != nil {
		return err
	}
	log.Info(c, "Logged in as %s.\n", u.User.Username)
	return nil
}

func DoLogOut(c *cli.Context) error {
	cl, err := NewClient(c, true)
	if err != nil {
		return err
	}

	err = cl.LogOut()
	if err != nil {
		if config.Debug() {
			log.ErrorlnQuit("Problem logging out: %v", err)
		}
		return err
	}

	// Delete local user data
	err = fileutils.DeleteFile(filepath.Join(config.UserDataDir(c.App.ExtraInfo()["configDir"]), config.UserFile))
	if err != nil {
		return err
	}

	return nil
}
