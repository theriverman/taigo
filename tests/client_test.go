package main

import (
	"net/http"
	"testing"

	taiga "github.com/theriverman/taigo"
)

const testHostURL string = "http://localhost:9000"
const testUsername string = "admin"
const testPassword string = "admin"
const testProjSlug string = "taigo-test"
const testProjID int = 2
const testUserID int = 5

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
}

func teardownClient() {
	Client = nil
}

func TestClient(t *testing.T) {
	setupClient()

	var makeurltests = []struct {
		in  []string
		out string
	}{
		{[]string{"epics"}, "http://localhost:9000/api/v1/epics"},
		{[]string{"epics", "5"}, "http://localhost:9000/api/v1/epics/5"},
		{[]string{"epics", "bulk_create"}, "http://localhost:9000/api/v1/epics/bulk_create"},
		{[]string{"epics", "attachments", "5"}, "http://localhost:9000/api/v1/epics/attachments/5"},
	}

	for _, tt := range makeurltests {
		s := Client.MakeURL(tt.in...)
		if s != tt.out {
			t.Errorf("got %q, want %q", s, tt.out)
		}
	}

}
