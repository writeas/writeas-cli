// Package executable holds utility functions that assist both CLI executables,
// writeas and wf.
package executable

import (
	"os"
	"path"
)

func Name() string {
	n := os.Args[0]
	return path.Base(n)
}
