package main

import (
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
)

/*
This cli.app.go file contains everything related to the main CLI application
All subsequent or lower-level commands are defined in cli.commands.go file.

Refer to the documentation of urfave/cli at https://github.com/urfave/cli
*/

// application behaviour
var appVerboseMode bool = false
var ignoreSavedState bool = false

// NewCLIApplication returns a *cli.App pointer which is the primary entrypoint to our CLI application.
// NewCLIApplication is the first call from main().
func NewCLIApplication() *cli.App {
	return &cli.App{
		Name:    AppName,
		Usage:   fmt.Sprintf("%s is a small application to interact with Taiga", AppName),
		Version: AppSemVersion,
		// application-level flags can be define below. these are applicable during the whole runtime
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Destination: &appVerboseMode,
				Usage:       fmt.Sprintf("Runs %s in verbose mode", AppName),
			},
			&cli.BoolFlag{
				Name:        "ignore-state",
				Destination: &ignoreSavedState,
				Usage:       "Ignore saved state (username,password,server,projects,etc.)",
			},
		},
		Before: func(c *cli.Context) error {
			// CLI flags are processed at this point. Consider configuring your logging level
			if appVerboseMode {
				fmt.Println("verbose mode has been enabled")
			}
			if err := autoLogin(); err != nil {
				log.Println("auto-authentication failed: " + err.Error())
			}
			return nil
		},
		Commands:  getApplicationCommands(),
		Copyright: AppCopyrightText,
		// see the urfave/cli documentation for all possible options: https://github.com/urfave/cli/blob/master/docs/v2/manual.md
	}
}
