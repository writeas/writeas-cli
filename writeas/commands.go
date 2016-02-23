package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func cmdPost(c *cli.Context) {
	err := handlePost(readStdIn(), c)
	check(err)
}

func cmdNew(c *cli.Context) {
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
		fmt.Printf("Error posting: %s\n", err)
		fmt.Println(messageRetryCompose(fname))
		os.Exit(1)
	}

	// Clean up temporary post
	if fname != "" {
		os.Remove(fname)
	}
}

func cmdDelete(c *cli.Context) {
	friendlyId := c.Args().Get(0)
	token := c.Args().Get(1)
	if friendlyId == "" {
		fmt.Println("usage: writeas delete <postId> [<token>]")
		os.Exit(1)
	}

	if token == "" {
		// Search for the token locally
		token = tokenFromID(friendlyId)
		if token == "" {
			fmt.Println("Couldn't find an edit token locally. Did you create this post here?")
			fmt.Printf("If you have an edit token, use: writeas delete %s <token>\n", friendlyId)
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

	DoDelete(friendlyId, token, tor)
}

func cmdUpdate(c *cli.Context) {
	friendlyId := c.Args().Get(0)
	token := c.Args().Get(1)
	if friendlyId == "" {
		fmt.Println("usage: writeas update <postId> [<token>]")
		os.Exit(1)
	}

	if token == "" {
		// Search for the token locally
		token = tokenFromID(friendlyId)
		if token == "" {
			fmt.Println("Couldn't find an edit token locally. Did you create this post here?")
			fmt.Printf("If you have an edit token, use: writeas update %s <token>\n", friendlyId)
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

	DoUpdate(fullPost, friendlyId, token, c.String("font"), tor, c.Bool("code"))
}

func cmdGet(c *cli.Context) {
	friendlyId := c.Args().Get(0)
	if friendlyId == "" {
		fmt.Println("usage: writeas get <postId>")
		os.Exit(1)
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

	DoFetch(friendlyId, tor)
}

func cmdAdd(c *cli.Context) {
	friendlyId := c.Args().Get(0)
	token := c.Args().Get(1)
	if friendlyId == "" || token == "" {
		fmt.Println("usage: writeas add <postId> <token>")
		os.Exit(1)
	}

	addPost(friendlyId, token)
}

func cmdList(c *cli.Context) {
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
}
