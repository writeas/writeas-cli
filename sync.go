package writeascli

import (
	//"github.com/writeas/writeas-cli/sync"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/writeas/writeas-cli/config"
	"github.com/writeas/writeas-cli/fileutils"
	cli "gopkg.in/urfave/cli.v1"
)

const (
	postFileExt  = ".txt"
	userFilename = "writeas_user"
)

func CmdPull(c *cli.Context) error {
	cfg, err := config.LoadConfig(userDataDir())
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
func syncSetUp(cfg *config.UserConfig) error {
	// Get user information and fail early (before we make the user do
	// anything), if we're going to
	u, err := LoadUser(userDataDir())
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
			if debug {
				Errorln("Error creating data directory: %s", err)
			}
			return err
		}
		// Create username file in directory
		err = ioutil.WriteFile(userFilepath, []byte(u.User.Username), 0644)
		fmt.Println("Created posts directory.")
	}

	// Save preference
	cfg.Posts.Directory = dir
	err = config.SaveConfig(userDataDir(), cfg)
	if err != nil {
		if debug {
			Errorln("Unable to save config: %s", err)
		}
		return err
	}
	fmt.Println("Saved config.")

	return nil
}
