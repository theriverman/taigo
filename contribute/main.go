package main

import (
	"fmt"
	"net/http"

	taiga "github.com/theriverman/taigo"
)

func main() {
	// Create client
	client := taiga.Client{
		BaseURL:    "https://api.taiga.io",
		HTTPClient: &http.Client{},
	}
	// Initialise client (authenticates to Taiga)
	err := client.Initialise()
	if err != nil {
		panic(err)
	}

	// Authenticate (get/set Token)
	client.AuthByCredentials(&taiga.Credentials{
		Type:     "normal",
		Username: "admin",
		Password: "123123",
	})
	me, err := client.User.Me()
	if err != nil {
		panic(err)
	}
	fmt.Println("Me: (ID, Username, FullName)", me.ID, me.Username, me.FullName)
}
