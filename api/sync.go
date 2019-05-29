package api

import (
	//"github.com/writeas/writeas-cli/sync"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/writeas/writeas-cli/config"
	"github.com/writeas/writeas-cli/fileutils"
	"github.com/writeas/writeas-cli/log"
	cli "gopkg.in/urfave/cli.v1"
)

const (
	PostFileExt  = ".txt"
	userFilename = "writeas_user"
)

func CmdPull(c *cli.Context) error {
	cfg, err := config.LoadConfig(config.UserDataDir(c.App.ExtraInfo()["configDir"]))
	if err != nil {
		return err
	}
	// Create posts directory if needed
	if cfg.Posts.Directory == "" {
		syncSetUp(c.App.ExtraInfo()["configDir"], cfg)
	}

	// Fetch posts
	cl, err := NewClient(c, true)
	if err != nil {
		return err
	}

	posts, err := cl.GetUserPosts()
	if err != nil {
		return err
	}

	for _, p := range *posts {
		postFilename := p.ID
		collDir := ""
		if p.Collection != nil {
			postFilename = p.Slug
			// Create directory for collection
			collDir = p.Collection.Alias
			if !fileutils.Exists(filepath.Join(cfg.Posts.Directory, collDir)) {
				log.Info(c, "Creating folder "+collDir)
				err = os.Mkdir(filepath.Join(cfg.Posts.Directory, collDir), 0755)
				if err != nil {
					log.Errorln("Error creating blog directory %s: %s. Skipping post %s.", collDir, err, postFilename)
					continue
				}
			}
		}
		postFilename += PostFileExt

		// Write file
		txtFile := p.Content
		if p.Title != "" {
			txtFile = "# " + p.Title + "\n\n" + txtFile
		}
		err = ioutil.WriteFile(filepath.Join(cfg.Posts.Directory, collDir, postFilename), []byte(txtFile), 0644)
		if err != nil {
			log.Errorln("Error creating file %s: %s", postFilename, err)
		}
		log.Info(c, "Saved post "+postFilename)

		// Update mtime and atime on files
		modTime := p.Updated.Local()
		err = os.Chtimes(filepath.Join(cfg.Posts.Directory, collDir, postFilename), modTime, modTime)
		if err != nil {
			log.Errorln("Error setting time on %s: %s", postFilename, err)
		}
	}

	return nil
}

func syncSetUp(path string, cfg *config.UserConfig) error {
	// Get user information and fail early (before we make the user do
	// anything), if we're going to
	u, err := config.LoadUser(config.UserDataDir(path))
	if err != nil {
		return err
	}

	// Prompt for posts directory
	defaultDir, err := os.Getwd()
	if err != nil {
		return err
	}
	var dir string
	fmt.Printf("Posts directory? [%s]: ", defaultDir)
	fmt.Scanln(&dir)
	if dir == "" {
		dir = defaultDir
	}

	// FIXME: This only works on non-Windows OSes (fix: https://www.reddit.com/r/golang/comments/5t3ezd/hidden_files_directories/)
	userFilepath := filepath.Join(dir, "."+userFilename)

	// Create directory if needed
	if !fileutils.Exists(dir) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			if config.Debug() {
				log.Errorln("Error creating data directory: %s", err)
			}
			return err
		}
		// Create username file in directory
		err = ioutil.WriteFile(userFilepath, []byte(u.User.Username), 0644)
		fmt.Println("Created posts directory.")
	}

	// Save preference
	cfg.Posts.Directory = dir
	err = config.SaveConfig(config.UserDataDir(path), cfg)
	if err != nil {
		if config.Debug() {
			log.Errorln("Unable to save config: %s", err)
		}
		return err
	}
	fmt.Println("Saved config.")

	return nil
}
