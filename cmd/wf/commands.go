package main

import (
	"fmt"

	"github.com/writeas/writeas-cli/config"
	"github.com/writeas/writeas-cli/executable"
	cli "gopkg.in/urfave/cli.v1"
)

func requireAuth(f cli.ActionFunc, action string) cli.ActionFunc {
	return func(c *cli.Context) error {
		u, err := config.LoadUser(c)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("couldn't load config: %v", err), 1)
		}
		if u == nil {
			return cli.NewExitError("You must be authenticated to "+action+".\nLog in first with: "+executable.Name()+" auth <username>", 1)
		}

		return f(c)
	}
}
