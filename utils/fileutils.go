package fileutils

import (
	"os"
)

// Exists returns whether or not the given file exists
func Exists(p string) bool {
	if _, err := os.Stat(p); err == nil {
		return true
	}
	return false
}
