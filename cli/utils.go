package main

import (
	"fmt"
	"log"

	"github.com/gookit/config/v2"
	"github.com/theriverman/taigo"
	"github.com/zalando/go-keyring"
)

// autoLogin attempts to authenticate the client against selected Taiga server.
// on first try, it prefers the user-provided credentials all credentials [host,username,password] must provided
// on second try, it retrieves the saved Token from the config file and does an attempt
// on third try, it retrieves the saved user credentisl from the host OS's keychain
// finally it exits the application stating that authentication has failed
func autoLogin() (err error) {
	// set host URL
	// client.BaseURL = taigoHostAddr

	// 1st attempt - try with user credentials
	if taigoUsername != "" && taigoPassword != "" && taigoHostAddr != "" {
		client.BaseURL = taigoHostAddr
		if err = client.AuthByCredentials(&taigo.Credentials{
			Type:     taigoLoginType,
			Username: taigoUsername,
			Password: taigoPassword,
		}); err == nil {
			if appVerboseMode {
				log.Println("authenticated with manually provided credentials")
			}
			return // OK
		}
	}

	// 2nd attempt - try with saved token
	client.BaseURL = config.String("TAIGO_HOST")
	if err = client.AuthByToken("Bearer", config.String("TAIGO_TOKEN")); err == nil {
		if appVerboseMode {
			log.Println("authenticated with token saved in .taigo.json")
		}
		return // OK
	}

	// 3rd attempt - try with saved credentials
	client.BaseURL = config.String("TAIGO_HOST")
	taigoUsername = config.String("TAIGO_USERNAME")
	taigoPassword, err = keyring.Get(AppName, taigoUsername)
	if err != nil {
		return err
	}
	if err = client.AuthByCredentials(&taigo.Credentials{
		Type:     taigoLoginType,
		Username: taigoUsername,
		Password: taigoPassword,
	}); err == nil {
		if appVerboseMode {
			log.Println("authenticated with credentials saved in .taigo.json")
		}
		return // OK
	}
	return fmt.Errorf("authentication has failed with all possible methods")
}

func hasDefaultProject() bool {
	return *defaultProjectID > 0
}
