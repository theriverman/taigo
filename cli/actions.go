package main

/*
	The application's business logic shall be implemented here in cli.actions.go in a function similar to `actionGreet` and `actionVersion` below.
	They must be added as a parameter in cli.commands.go where all possible application commands are defined.
*/

import (
	"fmt"
	"log"
	"runtime"

	"github.com/gookit/config/v2"
	taiga "github.com/theriverman/taigo"
	"github.com/urfave/cli/v2"
	"github.com/zalando/go-keyring"
)

func loginToTaiga(c *cli.Context) error {
	client.BaseURL = taigoHostAddr

	if err := client.AuthByCredentials(&taiga.Credentials{
		Type:     taigoLoginType,
		Username: taigoUsername,
		Password: taigoPassword,
	}); err != nil {
		log.Fatalln(err)
	}

	log.Printf("Hello, %s!\n", client.Self.FullName)

	if taigoNoStateSave {
		log.Println("Login credentials will not be saved locally")
	} else {
		// set password in keyring
		err := keyring.Set(AppName, taigoUsername, taigoPassword)
		if err != nil {
			log.Fatal(err)
		}
		// store state details in $HOME/.taigo.json
		if !ignoreSavedState {
			config.Set("TAIGO_HOST", taigoHostAddr)
			config.Set("TAIGO_USERNAME", taigoUsername)
			config.Set("TAIGO_TOKEN", client.Token)
			dumpConfigToFile(configPath)
			log.Println("login credentials have been saved locally")
		}
	}
	return nil
}

func actionVersion(c *cli.Context) error {
	fmt.Println(AppName + ":")
	fmt.Printf("  Version: %s\n", AppSemVersion)
	fmt.Printf("  Go version: %s\n", runtime.Version())
	fmt.Printf("  Git commit: %s\n", GitCommit)
	fmt.Printf("  Built: %s\n", AppBuildDate)
	fmt.Printf("  OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("  Build type: %s\n", AppBuildType)
	return nil
}

/*
	FAVOURITES
*/
func actionFavouritesList(c *cli.Context) error {
	return nil
}

func actionFavouritesAdd(c *cli.Context) (err error) {
	var project *taiga.Project
	switch {
	case projectID != 0:
		if project, err = client.Project.Get(projectID); err != nil {
			return err
		}
	case projectSlug != "":
		if project, err = client.Project.GetBySlug(projectSlug); err != nil {
			return err
		}
	default:
		return fmt.Errorf("--id or --ref must be provided to add a project to favourites")
	}
	fmt.Printf("Project name: %s \n\n", project.Name)
	fmt.Printf("Project ID: %d \n\n", project.ID)
	// favouriteProjects[project.Name] = project.ID
	// config.Set("TAIGO_HOST", taigoHostAddr)
	return nil
}

func actionFavouritesRemove(c *cli.Context) error {
	return nil
}
