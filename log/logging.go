package log

import (
	"fmt"
	"os"

	cli "gopkg.in/urfave/cli.v1"
)

// Info logs general diagnostic messages, shown only when the -v or --verbose
// flag is provided.
func Info(c *cli.Context, s string, p ...interface{}) {
	if c.Bool("v") || c.Bool("verbose") {
		fmt.Fprintf(os.Stderr, s+"\n", p...)
	}
}

// InfolnQuit logs the message to stderr and exits normally (without an error).
func InfolnQuit(s string, p ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", p...)
	os.Exit(0)
}

// Errorln logs the message to stderr
func Errorln(s string, p ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", p...)
}

// ErrorlnQuit logs the message to stderr and exits with an error code.
func ErrorlnQuit(s string, p ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", p...)
	os.Exit(1)
}
