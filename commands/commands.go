package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/howeyc/gopass"
	"github.com/writeas/writeas-cli/api"
	"github.com/writeas/writeas-cli/config"
	"github.com/writeas/writeas-cli/executable"
	"github.com/writeas/writeas-cli/log"
	cli "gopkg.in/urfave/cli.v1"
)

func CmdPost(c *cli.Context) error {
	if config.IsTor(c) {
		log.Info(c, "Publishing via hidden service...")
	} else {
		log.Info(c, "Publishing...")
	}

	_, err := api.DoPost(c, api.ReadStdIn(), c.String("font"), false, c.Bool("code"))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
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

	if config.IsTor(c) {
		log.Info(c, "Publishing via hidden service...")
	} else {
		log.Info(c, "Publishing...")
	}

	_, err := api.DoPost(c, *p, c.String("font"), false, c.Bool("code"))
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
		return cli.NewExitError("usage: "+executable.Name()+" publish <filename>", 1)
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if config.IsTor(c) {
		log.Info(c, "Publishing via hidden service...")
	} else {
		log.Info(c, "Publishing...")
	}
	_, err = api.DoPost(c, content, c.String("font"), false, c.Bool("code"))
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	// TODO: write local file if directory is set
	return nil
}

func CmdDelete(c *cli.Context) error {
	friendlyID := c.Args().Get(0)
	token := c.Args().Get(1)
	if friendlyID == "" {
		return cli.NewExitError("usage: "+executable.Name()+" delete <postId> [<token>]", 1)
	}

	u, _ := config.LoadUser(c)
	if token == "" {
		// Search for the token locally
		token = api.TokenFromID(c, friendlyID)
		if token == "" && u == nil {
			log.Errorln("Couldn't find an edit token locally. Did you create this post here?")
			log.ErrorlnQuit("If you have an edit token, use: "+executable.Name()+" delete %s <token>", friendlyID)
		}
	}

	if config.IsTor(c) {
		log.Info(c, "Deleting via hidden service...")
	} else {
		log.Info(c, "Deleting...")
	}

	err := api.DoDelete(c, friendlyID, token)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Couldn't delete post: %v", err), 1)
	}

	// TODO: Delete local file, if necessary
	return nil
}

func CmdUpdate(c *cli.Context) error {
	friendlyID := c.Args().Get(0)
	token := c.Args().Get(1)
	if friendlyID == "" {
		return cli.NewExitError("usage: "+executable.Name()+" update <postId> [<token>]", 1)
	}

	u, _ := config.LoadUser(c)
	if token == "" {
		// Search for the token locally
		token = api.TokenFromID(c, friendlyID)
		if token == "" && u == nil {
			log.Errorln("Couldn't find an edit token locally. Did you create this post here?")
			log.ErrorlnQuit("If you have an edit token, use: "+executable.Name()+" update %s <token>", friendlyID)
		}
	}

	// Read post body
	fullPost := api.ReadStdIn()

	if config.IsTor(c) {
		log.Info(c, "Updating via hidden service...")
	} else {
		log.Info(c, "Updating...")
	}
	err := api.DoUpdate(c, fullPost, friendlyID, token, c.String("font"), c.Bool("code"))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("%v", err), 1)
	}
	return nil
}

func CmdGet(c *cli.Context) error {
	friendlyID := c.Args().Get(0)
	if friendlyID == "" {
		return cli.NewExitError("usage: "+executable.Name()+" get <postId>", 1)
	}

	if config.IsTor(c) {
		log.Info(c, "Getting via hidden service...")
	} else {
		log.Info(c, "Getting...")
	}

	err := api.DoFetch(c, friendlyID)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("%v", err), 1)
	}
	return nil
}

func CmdAdd(c *cli.Context) error {
	friendlyID := c.Args().Get(0)
	token := c.Args().Get(1)
	if friendlyID == "" || token == "" {
		return cli.NewExitError("usage: "+executable.Name()+" add <postId> <token>", 1)
	}

	err := api.AddPost(c, friendlyID, token)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("%v", err), 1)
	}
	return nil
}

func CmdListPosts(c *cli.Context) error {
	urls := c.Bool("url")
	ids := c.Bool("id")
	details := c.Bool("v")

	posts := api.GetPosts(c)

	u, _ := config.LoadUser(c)
	if u != nil {
		if config.IsTor(c) {
			log.Info(c, "Getting posts via hidden service...")
		} else {
			log.Info(c, "Getting posts...")
		}
		remotePosts, err := api.GetUserPosts(c, true)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("error getting posts: %v", err), 1)
		}

		if len(remotePosts) > 0 {
			fmt.Println("Anonymous Posts")
			if details {
				identifier := "URL"
				if ids || !urls {
					identifier = "ID"
				}
				fmt.Println(identifier)
			}
		}
		for _, p := range remotePosts {
			identifier := getPostURL(c, p.ID)
			if ids || !urls {
				identifier = p.ID
			}

			fmt.Println(identifier)
		}

		if len(*posts) > 0 {
			fmt.Printf("\nUnclaimed Posts\n")
		}
	}

	if details {
		var p api.Post
		tw := tabwriter.NewWriter(os.Stdout, 10, 0, 2, ' ', tabwriter.TabIndent)
		numPosts := len(*posts)
		if ids || !urls && numPosts != 0 {
			fmt.Fprintf(tw, "%s\t%s\t\n", "ID", "Token")
		} else if numPosts != 0 {
			fmt.Fprintf(tw, "%s\t%s\t\n", "URL", "Token")
		} else {
			fmt.Fprintf(tw, "No local posts found\n")
		}
		for i := range *posts {
			p = (*posts)[numPosts-1-i]
			if ids || !urls {
				fmt.Fprintf(tw, "%s\t%s\t\n", p.ID, p.EditToken)
			} else {
				fmt.Fprintf(tw, "%s\t%s\t\n", getPostURL(c, p.ID), p.EditToken)
			}
		}
		return tw.Flush()
	}

	for _, p := range *posts {
		if ids || !urls {
			fmt.Printf("%s\n", p.ID)
		} else {
			fmt.Printf("%s\n", getPostURL(c, p.ID))
		}
	}
	return nil
}

func getPostURL(c *cli.Context, slug string) string {
	var base string
	if c.App.Name == "writeas" {
		if config.IsDev() {
			base = config.DevBaseURL
		} else {
			base = config.WriteasBaseURL
		}
	} else {
		if host := api.HostURL(c); host != "" {
			base = host
		} else {
			// TODO handle error, or load config globally, see T601
			// https://phabricator.write.as/T601
			cfg, _ := config.LoadConfig(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
			if cfg.Default.Host != "" && cfg.Default.User != "" {
				if parts := strings.Split(cfg.Default.Host, "://"); len(parts) > 1 {
					base = cfg.Default.Host
				} else {
					base = "https://" + cfg.Default.Host
				}
			}
		}
	}
	ext := ""
	// Output URL in requested format
	if c.Bool("md") {
		ext = ".md"
	}
	return fmt.Sprintf("%s/%s%s", base, slug, ext)
}

func CmdCollections(c *cli.Context) error {
	u, err := config.LoadUser(c)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("couldn't load config: %v", err), 1)
	}
	if u == nil {
		return cli.NewExitError("You must be authenticated to view collections.\nLog in first with: "+executable.Name()+" auth <username>", 1)
	}
	if config.IsTor(c) {
		log.Info(c, "Getting blogs via hidden service...")
	} else {
		log.Info(c, "Getting blogs...")
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
	u, err := config.LoadUser(c)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("couldn't load config: %v", err), 1)
	}
	if u == nil {
		return cli.NewExitError("You must be authenticated to claim local posts.\nLog in first with: "+executable.Name()+" auth <username>", 1)
	}

	localPosts := api.GetPosts(c)
	if len(*localPosts) == 0 {
		return nil
	}

	if config.IsTor(c) {
		log.Info(c, "Claiming %d post(s) for %s via hidden service...", len(*localPosts), u.User.Username)
	} else {
		log.Info(c, "Claiming %d post(s) for %s...", len(*localPosts), u.User.Username)
	}

	results, err := api.ClaimPosts(c, localPosts)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to claim posts: %v", err), 1)
	}

	var okCount, errCount int
	for _, r := range *results {
		id := r.ID
		if id == "" {
			// No top-level ID, so the claim was successful
			id = r.Post.ID
		}
		status := fmt.Sprintf("Post %s...", id)
		if r.ErrorMessage != "" {
			log.Errorln("%serror: %v", status, r.ErrorMessage)
			errCount++
		} else {
			log.Info(c, "%sOK", status)
			okCount++
			// only delete local if successful
			api.RemovePost(c, id)
		}
	}
	log.Info(c, "%d claimed, %d failed", okCount, errCount)
	return nil
}

func CmdAuth(c *cli.Context) error {
	username := c.Args().Get(0)
	if username == "" && c.GlobalIsSet("user") {
		username = c.GlobalString("user")
	}
	// Check configuration
	u, err := config.LoadUser(c)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("couldn't load config: %v", err), 1)
	}
	if u != nil && u.AccessToken != "" && username == u.User.Username {
		return cli.NewExitError("You're already authenticated as "+u.User.Username, 1)
	}

	// Validate arguments and get password
	if username == "" {
		cfg, err := config.LoadConfig(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Failed to load config: %v", err), 1)
		}
		if cfg.Default.Host != "" && cfg.Default.User != "" {
			username = cfg.Default.User
			fmt.Printf("No user provided, using default user %s for host %s...\n", cfg.Default.User, cfg.Default.Host)
		} else {
			return cli.NewExitError("usage: "+executable.Name()+" auth <username>", 1)
		}
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

	if config.IsTor(c) {
		log.Info(c, "Logging in via hidden service...")
	} else {
		log.Info(c, "Logging in...")
	}
	err = api.DoLogIn(c, username, string(pass))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("error logging in: %v", err), 1)
	}

	return nil
}

func CmdLogOut(c *cli.Context) error {
	if config.IsTor(c) {
		log.Info(c, "Logging out via hidden service...")
	} else {
		log.Info(c, "Logging out...")
	}
	err := api.DoLogOut(c)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("error logging out: %v", err), 1)
	}
	return nil
}
