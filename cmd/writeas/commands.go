package main

import (
	"fmt"
	"github.com/howeyc/gopass"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
)

func cmdPost(c *cli.Context) error {
	_, err := handlePost(readStdIn(), c)
	return err
}

func cmdNew(c *cli.Context) error {
	fname, p := composeNewPost()
	if p == nil {
		// Assume composeNewPost already told us what the error was. Abort now.
		os.Exit(1)
	}

	// Ensure we have something to post
	if len(*p) == 0 {
		// Clean up temporary post
		if fname != "" {
			os.Remove(fname)
		}

		InfolnQuit("Empty post. Bye!")
	}

	_, err := handlePost(*p, c)
	if err != nil {
		Errorln("Error posting: %s", err)
		Errorln(messageRetryCompose(fname))
		return cli.NewExitError("", 1)
	}

	// Clean up temporary post
	if fname != "" {
		os.Remove(fname)
	}

	return nil
}

func cmdPublish(c *cli.Context) error {
	filename := c.Args().Get(0)
	if filename == "" {
		return cli.NewExitError("usage: writeas publish <filename>", 1)
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	p, err := handlePost(content, c)
	if err != nil {
		return err
	}

	// Save post to posts folder
	cfg, err := loadConfig()
	if cfg.Posts.Directory != "" {
		err = WritePost(cfg.Posts.Directory, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func cmdDelete(c *cli.Context) error {
	friendlyID := c.Args().Get(0)
	token := c.Args().Get(1)
	if friendlyID == "" {
		return cli.NewExitError("usage: writeas delete <postId> [<token>]", 1)
	}

	if token == "" {
		// Search for the token locally
		token = tokenFromID(friendlyID)
		if token == "" {
			Errorln("Couldn't find an edit token locally. Did you create this post here?")
			ErrorlnQuit("If you have an edit token, use: writeas delete %s <token>", friendlyID)
		}
	}

	tor := isTor(c)
	if c.Int("tor-port") != 0 {
		torPort = c.Int("tor-port")
	}
	if tor {
		Info(c, "Deleting via hidden service...")
	} else {
		Info(c, "Deleting...")
	}

	return DoDelete(c, friendlyID, token, tor)
}

func cmdUpdate(c *cli.Context) error {
	friendlyID := c.Args().Get(0)
	token := c.Args().Get(1)
	if friendlyID == "" {
		return cli.NewExitError("usage: writeas update <postId> [<token>]", 1)
	}

	if token == "" {
		// Search for the token locally
		token = tokenFromID(friendlyID)
		if token == "" {
			Errorln("Couldn't find an edit token locally. Did you create this post here?")
			ErrorlnQuit("If you have an edit token, use: writeas update %s <token>", friendlyID)
		}
	}

	// Read post body
	fullPost := readStdIn()

	tor := isTor(c)
	if c.Int("tor-port") != 0 {
		torPort = c.Int("tor-port")
	}
	if tor {
		Info(c, "Updating via hidden service...")
	} else {
		Info(c, "Updating...")
	}

	return DoUpdate(c, fullPost, friendlyID, token, c.String("font"), tor, c.Bool("code"))
}

func cmdGet(c *cli.Context) error {
	friendlyID := c.Args().Get(0)
	if friendlyID == "" {
		return cli.NewExitError("usage: writeas get <postId>", 1)
	}

	tor := isTor(c)
	if c.Int("tor-port") != 0 {
		torPort = c.Int("tor-port")
	}
	if tor {
		Info(c, "Getting via hidden service...")
	} else {
		Info(c, "Getting...")
	}

	return DoFetch(friendlyID, userAgent(c), tor)
}

func cmdAdd(c *cli.Context) error {
	friendlyID := c.Args().Get(0)
	token := c.Args().Get(1)
	if friendlyID == "" || token == "" {
		return cli.NewExitError("usage: writeas add <postId> <token>", 1)
	}

	err := addPost(friendlyID, token)
	return err
}

func cmdList(c *cli.Context) error {
	urls := c.Bool("url")
	ids := c.Bool("id")

	var p Post
	posts := getPosts()
	for i := range *posts {
		p = (*posts)[len(*posts)-1-i]
		if ids || !urls {
			fmt.Printf("%s ", p.ID)
		}
		if urls {
			fmt.Printf("https://write.as/%s ", p.ID)
		}
		fmt.Print("\n")
	}
	return nil
}

func cmdAuth(c *cli.Context) error {
	// Check configuration
	u, err := loadUser()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("couldn't load config: %v", err), 1)
	}
	if u != nil && u.AccessToken != "" {
		return cli.NewExitError("You're already authenticated.", 1)
	}

	// Validate arguments and get password
	username := c.String("u")
	if username == "" {
		return cli.NewExitError("usage: writeas auth -u <username>", 1)
	}

	fmt.Print("Password: ")
	pass, err := gopass.GetPasswdMasked()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("error reading password: %v", err), 1)
	}

	// Validate password
	if len(pass) == 0 {
		return cli.NewExitError("Please enter your password.", 1)
	}
	return DoLogIn(c, username, string(pass))
}

func cmdLogOut(c *cli.Context) error {
	return DoLogOut(c)
}
