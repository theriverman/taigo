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
	projectID   int
	projectSlug string
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
							Usage:       "Project unique slug from the URL (e.g.: riverman-taigo-demo)",
							Destination: &projectSlug,
						},
						&cli.IntFlag{
							Name:        "id",
							Usage:       "Project internal ID",
							Destination: &projectID,
						},
					},
				},
				{
					Name:   "remove",
					Usage:  "Remove a project from your favourites",
					Action: actionFavouritesRemove,
					Flags: []cli.Flag{
						&cli.IntFlag{
							Name:        "id",
							Usage:       "Project internal ID",
							Destination: &projectID,
							Required:    true,
						},
					},
				},
				{
					Name:   "select",
					Usage:  "Select a project from your favourites as default",
					Action: actionFavouritesSelect,
					Flags: []cli.Flag{
						&cli.IntFlag{
							Name:        "id",
							Usage:       "Project internal ID",
							Destination: &projectID,
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:    "project",
			Usage:   "Read/Write project details",
			Aliases: []string{"p", "proj"},
		},
		{
			Name:    "epic",
			Usage:   "Read/Write epic details",
			Aliases: []string{"e"},
			Subcommands: []*cli.Command{
				{
					Name:   "list",
					Usage:  "List all epics in a project",
					Action: actionEpicList,
				},
			},
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:        "id",
					Usage:       "Project internal ID",
					Destination: &projectID,
					Value:       *defaultProjectID,
					Required:    !hasDefaultProject(),
				},
			},
		},
		{
			Name:    "sprint",
			Usage:   "Read/Write sprint details",
			Aliases: []string{"s"},
		},
		{
			Name:    "userstory",
			Usage:   "Read/Write userstory details",
			Aliases: []string{"us"},
		},
		{
			Name:    "task",
			Usage:   "Read/Write task details",
			Aliases: []string{"t", "subtask"},
		},
		{
			Name:   "version",
			Usage:  fmt.Sprintf("Show the %s version information (detailed)", AppName),
			Action: actionVersion,
		},
	}
}
