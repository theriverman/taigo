package main

import (
	"fmt"

	"github.com/gookit/config/v2"
	"github.com/urfave/cli/v2"
	"github.com/zalando/go-keyring"
)

/*
Refer to the documentation of urfave/cli at https://github.com/urfave/cli
*/

var (
	taigoHostAddr    string
	taigoUsername    string
	taigoPassword    string
	taigoLoginType   string
	taigoNoStateSave bool
)

var (
	projectID, projectIndex int
	projectSlug             string
)

func getApplicationCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "login",
			Usage:  "Authenticates you to Taiga",
			Action: loginToTaiga,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "host",
					Usage:       "Your Taiga server's API FQDN",
					Value:       "https://api.taiga.io",
					Destination: &taigoHostAddr,
					Required:    false,
					EnvVars:     []string{"TAIGO_HOST"},
				},
				&cli.StringFlag{
					Name:        "username",
					Usage:       "Your Taiga username or email (case sensitive)",
					Destination: &taigoUsername,
					Value:       taigoUsername,
					EnvVars:     []string{"TAIGO_USERNAME"},
					Required: func() bool {
						username := config.String("TAIGO_USERNAME")
						if username == "" {
							return true
						}
						taigoUsername = username
						return false
					}(),
				},
				&cli.StringFlag{
					Name:        "password",
					Usage:       "Your Taiga password",
					Destination: &taigoPassword,
					Value:       taigoPassword,
					EnvVars:     []string{"TAIGO_PASSWORD"},
					Required: func() bool {
						secret, err := keyring.Get(AppName, taigoUsername)
						if err != nil {
							return true
						}
						taigoPassword = secret
						return false
					}(),
				},
				&cli.StringFlag{
					Name:        "type",
					Usage:       "Used login type used for authentication",
					Value:       "normal",
					Destination: &taigoLoginType,
					Required:    false,
					EnvVars:     []string{"TAIGO_LOGIN_TYPE"},
				},
				&cli.BoolFlag{
					Name:        "no-save",
					Usage:       "Don't save login credentials locally following a successful authentication",
					Value:       false,
					Destination: &taigoNoStateSave,
				},
			},
		},
		{
			Name:  "favourites",
			Usage: "Manage your favourite projects",
			Subcommands: []*cli.Command{
				{
					Name:   "list",
					Usage:  "List your saved projects",
					Action: actionFavouritesList,
				},
				{
					Name:   "add",
					Usage:  "Save a project to your favourites",
					Action: actionFavouritesAdd,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "slug",
							Usage:       "Project's unique slug from the URL (e.g.: riverman-taigo-demo)",
							Destination: &projectSlug,
						},
						&cli.IntFlag{
							Name:        "id",
							Usage:       "Project's ID from the URL",
							Destination: &projectID,
						},
					},
				},
				{
					Name:   "remove",
					Usage:  "Remove a project from your favourites",
					Action: actionFavouritesRemove,
				},
				{
					Name:   "select",
					Usage:  "Select a project from your favourites as default",
					Action: actionFavouritesRemove,
					Flags: []cli.Flag{
						&cli.IntFlag{
							Name:        "index",
							Usage:       "Project's index from the list",
							Destination: &projectIndex,
						},
					},
				},
			},
		},
		{
			Name:   "version",
			Usage:  fmt.Sprintf("Show the %s version information (detailed)", AppName),
			Action: actionVersion,
		},
	}
}
