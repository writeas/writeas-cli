package main

import (
	//"github.com/writeas/writeas-cli/sync"
	"fmt"
	"github.com/writeas/writeas-cli/fileutils"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	postFileExt = ".txt"
)

func cmdPull(c *cli.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	// Create posts directory if needed
	if cfg.Posts.Directory == "" {
		syncSetUp(cfg)
	}

	// Fetch posts
	cl, err := newClient(c, true)
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
				Info(c, "Creating folder "+collDir)
				err = os.Mkdir(filepath.Join(cfg.Posts.Directory, collDir), 0755)
				if err != nil {
					Errorln("Error creating blog directory %s: %s. Skipping post %s.", collDir, err, postFilename)
					continue
				}
			}
		}
		postFilename += postFileExt

		// Write file
		txtFile := p.Content
		if p.Title != "" {
			txtFile = "# " + p.Title + "\n\n" + txtFile
		}
		err = ioutil.WriteFile(filepath.Join(cfg.Posts.Directory, collDir, postFilename), []byte(txtFile), 0644)
		if err != nil {
			Errorln("Error creating file %s: %s", postFilename, err)
		}
		Info(c, "Saved post "+postFilename)

		// Update mtime and atime on files
		modTime := p.Updated.Local()
		err = os.Chtimes(filepath.Join(cfg.Posts.Directory, collDir, postFilename), modTime, modTime)
		if err != nil {
			Errorln("Error setting time on %s: %s", postFilename, err)
		}
	}

	return nil
}

// TODO: move UserConfig to its own package, and this to sync package
func syncSetUp(cfg *UserConfig) error {
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

	// Create directory if needed
	if !fileutils.Exists(dir) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			if debug {
				Errorln("Error creating data directory: %s", err)
			}
			return err
		}
		fmt.Println("Created posts directory.")
	}

	// Save preference
	cfg.Posts.Directory = dir
	err = saveConfig(cfg)
	if err != nil {
		if debug {
			Errorln("Unable to save config: %s", err)
		}
		return err
	}
	fmt.Println("Saved config.")

	return nil
}
