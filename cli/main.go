package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	taiga "github.com/theriverman/taigo"
)

const TaigoCfgFolder string = "TaigoCLI"
const UserCfgFilename string = "taigocli.cfg.yaml"
const DefaultTaigaHost string = "https://api.taiga.io"
const DefaultTaigaHostAuthType string = "normal"

var cfg *UserConfig = &UserConfig{}
var client taiga.Client

var AppName string = "Taigo CLI"
var username, password string // if set, overrides cfg.Username|Password

func init() {
	/*
		setup userspace
	*/
	// get home
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	// check/create TaigoCli directory existence
	TaigoCliHomeDir := filepath.Join(userHomeDir, TaigoCfgFolder)
	TaigoCliCfgPath := filepath.Join(TaigoCliHomeDir, UserCfgFilename)
	if err = os.MkdirAll(TaigoCliHomeDir, 0644); err != nil {
		log.Fatal(err)
	}
	// get/create config file
	if cfg, err = readConfigFile(TaigoCliCfgPath); err != nil {
		log.Printf("config file not found (%s). creating a new one", err)
		if err = writeConfigFile(cfg, TaigoCliCfgPath); err != nil {
			log.Fatalln("writeConfigFile:", err)
		}
	}
	/*
		setup Taigo client
	*/
	// new client
	client = taiga.Client{
		BaseURL:    DefaultTaigaHost,
		HTTPClient: &http.Client{},
	}
	// store the password in an encrypted fashion
	if password != "" {
		cfg.Password = password
	}
	if cfg.PasswordEncrypted != "" {
		cfg.Password = PasswordDecrypt(cfg.PasswordEncrypted)
	} else if cfg.Password != "" {
		cfg.PasswordEncrypted = PasswordEncrypt(cfg.Password)
	}
	// set default values if empty/missing from config file
	if cfg.Host != "" {
		client.BaseURL = cfg.Host
	}
	if username != "" {
		cfg.Username = username
	}
	if err = writeConfigFile(cfg, TaigoCliCfgPath); err != nil {
		log.Fatalln("writeConfigFile:", err)
	}
}

func main() {
	// init client
	if err := client.Initialise(); err != nil {
		log.Fatal(err)
	}
	// create credentials
	credentials := &taiga.Credentials{
		Type:     DefaultTaigaHostAuthType,
		Username: cfg.Username,
		Password: cfg.Password,
	}
	// set host auth type
	if cfg.HostAuthType != "" {
		credentials.Type = cfg.HostAuthType
	}
	// authenticate (get/set token)
	if err := client.AuthByCredentials(credentials); err != nil {
		log.Fatal(err)
	}
	// get /users/me
	me, _ := client.User.Me()
	fmt.Println("Me: (ID, Username, FullName)", me.ID, me.Username, me.FullName)

	// Get Project (by its slug)
	// slug := "therivermantaigo-taigo-public-test"
	// fmt.Printf("Getting Project (slug=%s)..\n", slug)
	// project, err := client.Project.GetBySlug(slug)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("Project name: %s \n\n", project.Name)
}
