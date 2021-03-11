package main

import (
	"math/rand"
	"net/http"
	"time"

	taiga "github.com/theriverman/taigo"
)

const testHostURL string = "http://localhost:9000"
const testUsername string = "admin"
const testPassword string = "admin"
const testProjSlug string = "taigo-test"
const testProjID int = 2
const testUserID int = 5

const pool = "abcdefghijklmnopqrstuvwxyzABCEFGHIJKLMNOPQRSTUVWXYZ"

// Client is the foundation for making requests against Taiga
var Client *taiga.Client = nil

func setupClient() {
	if Client != nil {
		return // client already set; skipping
	}

	// Create client
	client := taiga.Client{
		BaseURL:    testHostURL,
		HTTPClient: &http.Client{},
	}
	// Initialise client (authenticates to Taiga)
	err := client.Initialise()
	if err != nil {
		panic(err)
	}
	err = client.AuthByCredentials(&taiga.Credentials{
		Type:     "normal",
		Username: testUsername,
		Password: testPassword,
	})
	if err != nil {
		panic(err)
	}
	Client = &client

	// Random String Generation
	rand.Seed(time.Now().UnixNano())

}

func teardownClient() {
	Client = nil
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = pool[rand.Intn(len(pool))]
	}

	return string(bytes)
}
