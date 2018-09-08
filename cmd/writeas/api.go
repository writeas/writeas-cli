package main

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/writeas/go-writeas"
	"gopkg.in/urfave/cli.v1"
)

const (
	defaultUserAgent = "writeas-cli v" + version
)

func client(userAgent string, tor bool) *writeas.Client {
	var client *writeas.Client
	if tor {
		client = writeas.NewTorClient(torPort)
	} else {
		client = writeas.NewClient()
	}
	client.UserAgent = userAgent

	return client
}

// DoFetch retrieves the Write.as post with the given friendlyID,
// optionally via the Tor hidden service.
func DoFetch(friendlyID, ua string, tor bool) error {
	cl := client(ua, tor)

	p, err := cl.GetPost(friendlyID)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", string(p.Content))
	return nil
}

// DoPost creates a Write.as post, returning an error if it was
// unsuccessful.
func DoPost(c *cli.Context, post []byte, font string, encrypt, tor, code bool) error {
	cl := client(userAgent(c), tor)

	p, err := cl.CreatePost(&writeas.PostParams{
		// TODO: extract title
		Content: string(post),
		Font:    getFont(code, font),
	})
	if err != nil {
		return fmt.Errorf("Unable to post: %v", err)
	}

	url := writeasBaseURL
	if tor {
		url = torBaseURL
	}
	url += "/" + p.ID

	// Store post locally
	addPost(p.ID, p.Token)

	// Copy URL to clipboard
	err = clipboard.WriteAll(string(url))
	if err != nil {
		Errorln("writeas: Didn't copy to clipboard: %s", err)
	} else {
		Info(c, "Copied to clipboard.")
	}

	// Output URL
	fmt.Printf("%s\n", url)

	return nil
}

// DoUpdate updates the given post on Write.as.
func DoUpdate(c *cli.Context, post []byte, friendlyID, token, font string, tor, code bool) error {
	cl := client(userAgent(c), tor)

	params := writeas.PostParams{
		ID:      friendlyID,
		Token:   token,
		Content: string(post),
		// TODO: extract title
	}
	if code || font != "" {
		params.Font = getFont(code, font)
	}

	_, err := cl.UpdatePost(&params)
	if err != nil {
		if debug {
			ErrorlnQuit("Problem updating: %v", err)
		}
		return fmt.Errorf("Post doesn't exist, or bad edit token given.")
	}

	if tor {
		Info(c, "Post updated via hidden service.")
	} else {
		Info(c, "Post updated.")
	}
	return nil
}

// DoDelete deletes the given post on Write.as.
func DoDelete(c *cli.Context, friendlyID, token string, tor bool) error {
	cl := client(userAgent(c), tor)

	err := cl.DeletePost(&writeas.PostParams{
		ID:    friendlyID,
		Token: token,
	})
	if err != nil {
		if debug {
			ErrorlnQuit("Problem deleting: %v", err)
		}
		return fmt.Errorf("Post doesn't exist, or bad edit token given.")
	}

	if tor {
		Info(c, "Post deleted from hidden service.")
	} else {
		Info(c, "Post deleted.")
	}
	removePost(friendlyID)

	return nil
}

func DoLogIn(c *cli.Context, uc *UserConfig, username, password string) error {
	cl := client(userAgent(c), isTor(c))

	u, err := cl.LogIn(username, password)
	if err != nil {
		if debug {
			ErrorlnQuit("Problem logging in: %v", err)
		}
		return err
	}

	uc.API.Token = u.AccessToken
	Info(c, "Success.")
	return saveConfig(uc)
}
