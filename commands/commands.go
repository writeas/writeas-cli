package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/howeyc/gopass"
	"github.com/writeas/writeas-cli/api"
	"github.com/writeas/writeas-cli/config"
	"github.com/writeas/writeas-cli/fileutils"
	"github.com/writeas/writeas-cli/log"
	cli "gopkg.in/urfave/cli.v1"
)

func CmdPost(c *cli.Context) error {
	_, err := api.HandlePost(api.ReadStdIn(), c)
	return err
}

func CmdNew(c *cli.Context) error {
	fname, p := api.ComposeNewPost()
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

		log.InfolnQuit("Empty post. Bye!")
	}

	_, err := api.HandlePost(*p, c)
	if err != nil {
		log.Errorln("Error posting: %s\n%s", err, config.MessageRetryCompose(fname))
		return cli.NewExitError("", 1)
	}

	// Clean up temporary post
	if fname != "" {
		os.Remove(fname)
	}

	return nil
}

func CmdPublish(c *cli.Context) error {
	filename := c.Args().Get(0)
	if filename == "" {
		return cli.NewExitError("usage: writeas publish <filename>", 1)
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	p, err := api.HandlePost(content, c)
	if err != nil {
		return err
	}

	// Save post to posts folder
	cfg, err := config.LoadConfig(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
	if cfg.Posts.Directory != "" {
		err = api.WritePost(cfg.Posts.Directory, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func CmdDelete(c *cli.Context) error {
	friendlyID := c.Args().Get(0)
	token := c.Args().Get(1)
	if friendlyID == "" {
		return cli.NewExitError("usage: writeas delete <postId> [<token>]", 1)
	}

	u, _ := config.LoadUser(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
	if token == "" {
		// Search for the token locally
		token = api.TokenFromID(c, friendlyID)
		if token == "" && u == nil {
			log.Errorln("Couldn't find an edit token locally. Did you create this post here?")
			log.ErrorlnQuit("If you have an edit token, use: writeas delete %s <token>", friendlyID)
		}
	}

	tor := config.IsTor(c)
	if c.Int("tor-port") != 0 {
		api.TorPort = c.Int("tor-port")
	}
	if tor {
		log.Info(c, "Deleting via hidden service...")
	} else {
		log.Info(c, "Deleting...")
	}

	err := api.DoDelete(c, friendlyID, token, tor)
	if err != nil {
		return err
	}

	// Delete local file, if necessary
	cfg, err := config.LoadConfig(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
	if cfg.Posts.Directory != "" {
		// TODO: handle deleting blog posts
		err = fileutils.DeleteFile(filepath.Join(cfg.Posts.Directory, friendlyID+api.PostFileExt))
		if err != nil {
			return err
		}
	}

	return nil
}

func CmdUpdate(c *cli.Context) error {
	friendlyID := c.Args().Get(0)
	token := c.Args().Get(1)
	if friendlyID == "" {
		return cli.NewExitError("usage: writeas update <postId> [<token>]", 1)
	}

	u, _ := config.LoadUser(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
	if token == "" {
		// Search for the token locally
		token = api.TokenFromID(c, friendlyID)
		if token == "" && u == nil {
			log.Errorln("Couldn't find an edit token locally. Did you create this post here?")
			log.ErrorlnQuit("If you have an edit token, use: writeas update %s <token>", friendlyID)
		}
	}

	// Read post body
	fullPost := api.ReadStdIn()

	tor := config.IsTor(c)
	if c.Int("tor-port") != 0 {
		api.TorPort = c.Int("tor-port")
	}
	if tor {
		log.Info(c, "Updating via hidden service...")
	} else {
		log.Info(c, "Updating...")
	}

	return api.DoUpdate(c, fullPost, friendlyID, token, c.String("font"), tor, c.Bool("code"))
}

func CmdGet(c *cli.Context) error {
	friendlyID := c.Args().Get(0)
	if friendlyID == "" {
		return cli.NewExitError("usage: writeas get <postId>", 1)
	}

	tor := config.IsTor(c)
	if c.Int("tor-port") != 0 {
		api.TorPort = c.Int("tor-port")
	}
	if tor {
		log.Info(c, "Getting via hidden service...")
	} else {
		log.Info(c, "Getting...")
	}

	return api.DoFetch(friendlyID, config.UserAgent(c), tor)
}

func CmdAdd(c *cli.Context) error {
	friendlyID := c.Args().Get(0)
	token := c.Args().Get(1)
	if friendlyID == "" || token == "" {
		return cli.NewExitError("usage: writeas add <postId> <token>", 1)
	}

	err := api.AddPost(c, friendlyID, token)
	return err
}

func CmdList(c *cli.Context) error {
	urls := c.Bool("url")
	ids := c.Bool("id")

	var p api.Post
	posts := api.GetPosts(c)
	for i := range *posts {
		p = (*posts)[len(*posts)-1-i]
		if ids || !urls {
			fmt.Printf("%s ", p.ID)
		}
		if urls {
			base := config.WriteasBaseURL
			if config.IsDev() {
				base = config.DevBaseURL
			}
			ext := ""
			// Output URL in requested format
			if c.Bool("md") {
				ext = ".md"
			}
			fmt.Printf("%s/%s%s ", base, p.ID, ext)
		}
		fmt.Print("\n")
	}
	return nil
}

func CmdAuth(c *cli.Context) error {
	// Check configuration
	u, err := config.LoadUser(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("couldn't load config: %v", err), 1)
	}
	if u != nil && u.AccessToken != "" {
		return cli.NewExitError("You're already authenticated as "+u.User.Username+". Log out with: writeas logout", 1)
	}

	// Validate arguments and get password
	username := c.Args().Get(0)
	if username == "" {
		return cli.NewExitError("usage: writeas auth <username>", 1)
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
	err = api.DoLogIn(c, username, string(pass))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("error logging in: %v", err), 1)
	}

	return nil
}

func CmdLogOut(c *cli.Context) error {
	return api.DoLogOut(c)
}
