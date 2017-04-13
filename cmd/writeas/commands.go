package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func cmdPost(c *cli.Context) error {
	err := handlePost(readStdIn(), c)
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

		fmt.Println("Empty post. Bye!")
		os.Exit(0)
	}

	err := handlePost(*p, c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error posting: %s\n", err)
		fmt.Fprintf(os.Stderr, messageRetryCompose(fname))
		return cli.NewExitError("", 1)
	}

	// Clean up temporary post
	if fname != "" {
		os.Remove(fname)
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
			fmt.Println("Couldn't find an edit token locally. Did you create this post here?")
			fmt.Printf("If you have an edit token, use: writeas delete %s <token>\n", friendlyID)
			os.Exit(1)
		}
	}

	tor := c.Bool("tor") || c.Bool("t")
	if c.Int("tor-port") != 0 {
		torPort = c.Int("tor-port")
	}
	if tor {
		fmt.Println("Deleting via hidden service...")
	} else {
		fmt.Println("Deleting...")
	}

	return DoDelete(friendlyID, token, tor)
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
			fmt.Println("Couldn't find an edit token locally. Did you create this post here?")
			fmt.Printf("If you have an edit token, use: writeas update %s <token>\n", friendlyID)
			os.Exit(1)
		}
	}

	// Read post body
	fullPost := readStdIn()

	tor := c.Bool("tor") || c.Bool("t")
	if c.Int("tor-port") != 0 {
		torPort = c.Int("tor-port")
	}
	if tor {
		fmt.Println("Updating via hidden service...")
	} else {
		fmt.Println("Updating...")
	}

	return DoUpdate(fullPost, friendlyID, token, c.String("font"), tor, c.Bool("code"))
}

func cmdGet(c *cli.Context) error {
	friendlyID := c.Args().Get(0)
	if friendlyID == "" {
		return cli.NewExitError("usage: writeas get <postId>", 1)
	}

	tor := c.Bool("tor") || c.Bool("t")
	if c.Int("tor-port") != 0 {
		torPort = c.Int("tor-port")
	}
	if tor {
		fmt.Println("Getting via hidden service...")
	} else {
		fmt.Println("Getting...")
	}

	return DoFetch(friendlyID, tor)
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
