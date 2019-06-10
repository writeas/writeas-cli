package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/tabwriter"

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
		return cli.NewExitError(fmt.Sprintf("Couldn't delete remote copy: %v", err), 1)
	}

	// Delete local file, if necessary
	cfg, err := config.LoadConfig(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
	if cfg.Posts.Directory != "" {
		// TODO: handle deleting blog posts
		err = fileutils.DeleteFile(filepath.Join(cfg.Posts.Directory, friendlyID+api.PostFileExt))
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Couldn't delete local copy: %v", err), 1)
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

func CmdListPosts(c *cli.Context) error {
	urls := c.Bool("url")
	ids := c.Bool("id")

	var p api.Post
	posts := api.GetPosts(c)
	tw := tabwriter.NewWriter(os.Stdout, 10, 0, 2, ' ', tabwriter.TabIndent)
	numPosts := len(*posts)
	if ids || !urls && numPosts != 0 {
		fmt.Fprintf(tw, "Local\t%s\t%s\t\n", "ID", "Token")
	} else if numPosts != 0 {
		fmt.Fprintf(tw, "Local\t%s\t%s\t\n", "URL", "Token")
	} else {
		fmt.Fprintf(tw, "No local posts found\n")
	}
	for i := range *posts {
		p = (*posts)[numPosts-1-i]
		if ids || !urls {
			fmt.Fprintf(tw, "unsynced\t%s\t%s\t\n", p.ID, p.EditToken)
		} else {
			fmt.Fprintf(tw, "unsynced\t%s\t%s\t\n", getPostURL(c, p.ID), p.EditToken)
		}
	}
	u, _ := config.LoadUser(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
	if u != nil {
		remotePosts, err := api.GetUserPosts(c)
		if err != nil {
			fmt.Println(err)
		}

		if len(remotePosts) > 0 {
			identifier := "URL"
			if ids || !urls {
				identifier = "ID"
			}
			fmt.Fprintf(tw, "\nAccount\t%s\t%s\t\n", identifier, "Title")
		}
		for _, p := range remotePosts {
			identifier := getPostURL(c, p.ID)
			if ids || !urls {
				identifier = p.ID
			}
			synced := "unsynced"
			if p.Synced {
				synced = "synced"
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\t\n", synced, identifier, p.Title)
		}
	}
	return tw.Flush()
}

func getPostURL(c *cli.Context, slug string) string {
	base := config.WriteasBaseURL
	if config.IsDev() {
		base = config.DevBaseURL
	}
	ext := ""
	// Output URL in requested format
	if c.Bool("md") {
		ext = ".md"
	}
	return fmt.Sprintf("%s/%s%s", base, slug, ext)
}

func CmdCollections(c *cli.Context) error {
	u, err := config.LoadUser(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("couldn't load config: %v", err), 1)
	}
	if u == nil {
		return cli.NewExitError("You must be authenticated to view collections.\nLog in first with: writeas auth <username>", 1)
	}
	colls, err := api.DoFetchCollections(c)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't get collections for user %s: %v", u.User.Username, err), 1)
	}
	urls := c.Bool("url")
	tw := tabwriter.NewWriter(os.Stdout, 8, 0, 2, ' ', tabwriter.TabIndent)
	detail := "Title"
	if urls {
		detail = "URL"
	}
	fmt.Fprintf(tw, "%s\t%s\t\n", "Alias", detail)
	for _, c := range colls {
		dData := c.Title
		if urls {
			dData = c.URL
		}
		fmt.Fprintf(tw, "%s\t%s\t\n", c.Alias, dData)
	}
	tw.Flush()
	return nil
}

func CmdClaim(c *cli.Context) error {
	u, err := config.LoadUser(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("couldn't load config: %v", err), 1)
	}
	if u == nil {
		return cli.NewExitError("You must be authenticated to claim local posts.\nLog in first with: writeas auth <username>", 1)
	}

	localPosts := api.GetPosts(c)
	if len(*localPosts) == 0 {
		return nil
	}

	results, err := api.ClaimPosts(c, localPosts)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to claim posts: %v", err), 1)
	}

	for _, r := range *results {
		fmt.Printf("Adding %s to user %s..", r.Post.ID, u.User.Username)
		if r.ErrorMessage != "" {
			fmt.Printf(" Failed\n")
			if config.Debug() {
				log.Errorln("Failed claiming post %s: %v", r.ID, r.ErrorMessage)
			}
		} else {
			fmt.Printf(" OK\n")
			// only delete local if successful
			api.RemovePost(c.App.ExtraInfo()["configDir"], r.Post.ID)
		}
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
