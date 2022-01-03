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
	hasNoDefaultSelected := true
	fmt.Println("Your saved favourite projects:")
	for pid, pn := range configStruct.FavouriteProjects {
		if pid == *defaultProjectID {
			fmt.Printf("  * %d\t%s [selected as default]\n", pid, pn)
			hasNoDefaultSelected = false
		} else {
			fmt.Printf("  * %d\t%s\n", pid, pn)
		}
	}
	if hasNoDefaultSelected {
		fmt.Println("---------------------------------------------")
		fmt.Println("  You have no project selected as default. See ./taigo favourites select -h")
	}
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
		return fmt.Errorf("--id or --slug must be provided to add a project to favourites")
	}

	if configStruct.FavouriteProjects == nil {
		configStruct.FavouriteProjects = make(FavouriteProjects)
	}
	configStruct.FavouriteProjects[project.ID] = project.Name
	config.Set("FavouriteProjects", configStruct.FavouriteProjects)
	return nil
}

func actionFavouritesRemove(c *cli.Context) error {
	delete(configStruct.FavouriteProjects, projectID)
	config.Set("FavouriteProjects", configStruct.FavouriteProjects)
	return nil
}

func actionFavouritesSelect(c *cli.Context) (err error) {
	if err = actionFavouritesAdd(c); err != nil {
		return err
	}
	err = config.Set("FavouriteProjectID", projectID)
	return
}

/*
	EPICS
*/
func actionEpicList(c *cli.Context) (err error) {
	qp := taiga.EpicsQueryParams{}
	// defaultProjectID

	if projectID > 0 {
		qp.Project = projectID
	} else if *defaultProjectID > 0 {
		qp.Project = *defaultProjectID
	} else {
		return fmt.Errorf("add flag --id to define which project to use")
	}

	epics, err := client.Epic.List(&qp)
	for _, epic := range epics {
		fmt.Printf("  * ID: %d\tRef: %d\tSubject: %s\n", epic.ID, epic.Ref, epic.Subject)
	}
	return
}
