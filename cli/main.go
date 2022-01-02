package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/json"
	taiga "github.com/theriverman/taigo"
	"github.com/urfave/cli/v2"
)

// application details injected at build time
var (
	AppName          string = "app" // pretty-formatted
	AppBuildType     string = "unreleased/internal"
	AppBuildDate     string = time.Now().Format("02 Jan 2006 15:04:05") // equals to date '+%c'
	AppSemVersion    string = "no-version"
	AppCopyrightText string = "no copyright"
	GitCommit        string = "commit-id-could-not-be-retrieved"
	configPath       string
	client           taiga.Client
	// taigoAuthToken   *string = &client.Token
)

func init() {
	// create a taigo client
	client = taiga.Client{
		HTTPClient: &http.Client{},
	}

	// get $HOME
	homeDirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// load/dump config
	configPath = path.Join(homeDirname, ".taigo.json")
	config.WithOptions(config.ParseEnv)
	config.AddDriver(json.Driver)
	err = config.LoadFiles(configPath)
	if err != nil {
		log.Println(err)
		dumpConfigToFile(configPath) // dump default config b/c it doesn't exist
	} else {
		log.Printf("config loaded from %s\n", configPath)
	}
}

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "Prints version information of go-socks5-cli and quit",
	}

	app := NewCLIApplication()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
