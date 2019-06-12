package commands

import (
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/writeas/writeas-cli/api"
	"github.com/writeas/writeas-cli/config"
	cli "gopkg.in/urfave/cli.v1"
)

const (
	dashBar = "--------------------------------------------------------------------------------"
)

func printLocalPosts(c *cli.Context, tw *tabwriter.Writer, posts *[]api.Post) {
	ids := c.Bool("id")
	urls := c.Bool("url")
	numPosts := len(*posts)
	if ids || !urls && numPosts != 0 {
		fmt.Fprintf(tw, "Local\t%s\t%s\t\n", "ID", "Token")
	} else if numPosts != 0 {
		fmt.Fprintf(tw, "Local\t%s\t%s\t\n", "URL", "Token")
	} else {
		fmt.Fprintf(tw, "No local posts found\n")
	}
	for i := range *posts {
		p := (*posts)[numPosts-1-i]
		if ids || urls {
			fmt.Fprintf(tw, "unsynced\t%s\t%s\n", p.ID, p.EditToken)
		} else {
			fmt.Fprintf(tw, "unsynced\t%s\t%s\n", getPostURL(c, p.ID), p.EditToken)
		}
	}
}

func printRemotePosts(c *cli.Context, tw *tabwriter.Writer) error {
	ids := c.Bool("id")
	urls := c.Bool("url")
	details := c.Bool("d")
	u, _ := config.LoadUser(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
	if u != nil {
		remotePosts, err := api.GetUserPosts(c)
		if err != nil {
			return err
		}

		if len(remotePosts) > 0 {
			identifier := "URL"
			if ids || !urls {
				identifier = "ID"
			}
			if !details {
				fmt.Fprintf(tw, "\nAccount\t%s\t%s\n", identifier, "Title")
			}
		}
		for _, p := range remotePosts {
			identifier := getPostURL(c, p.ID)
			if ids || !urls {
				identifier = p.ID
			}
			if details {
				slug := p.ID
				if p.Slug != "" {
					slug = p.Slug
					if p.Collection != "" {
						slug = p.Collection + "/" + slug
					}
				}
				url := getPostURL(c, slug)
				if p.Slug == "" {
					p.Slug = "no-slug"
				}
				if p.Collection == "" {
					p.Collection = "no-blog"
				}
				fmt.Fprintf(tw, "\n%s\t%s\t%s\n", "Title: ", p.Title, " ")
				fmt.Fprintf(tw, "%s\t%s\t%s\n", "Last Updated: ", prettyDate(p.Updated), "")
				fmt.Fprintf(tw, "%s\t%s\t%s\n", "Blog: ", p.Collection, " ")
				fmt.Fprintf(tw, "%s\t%s\t%s\n", "Slug/ID: ", p.Slug+" / "+p.ID, " ")
				fmt.Fprintf(tw, "%s\t%s\t%s\n\n", "URL: ", url, " ")
			}
			synced := "unsynced"
			if p.Synced {
				synced = "synced"
			}
			if details {
				fmt.Fprintln(tw, p.Excerpt)
				fmt.Fprintln(tw, dashBar)
			} else {
				fmt.Fprintf(tw, "%s\t%s\t%s\n", synced, identifier, p.Title)
			}
		}
	}
	return nil
}

func prettyDate(date time.Time) string {
	return date.Local().Format(time.RFC822)
}
