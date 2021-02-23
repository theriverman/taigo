package main

// import (
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"testing"

// 	taiga "github.com/theriverman/taigo"
// )

// var (
// 	comTestClient  *taiga.Client = nil
// 	dummyProjectID int
// )

// func setupClient() {
// 	if comTestClient != nil {
// 		return // client already set; skipping
// 	}

// 	url, ok := os.LookupEnv("CI_URL")
// 	if !ok {
// 		panic("Missing Environment Variable: CI_URL")
// 	}
// 	username, ok := os.LookupEnv("CI_USERNAME")
// 	if !ok {
// 		panic("Missing Environment Variable: CI_USERNAME")
// 	}
// 	password, ok := os.LookupEnv("CI_PASSWORD")
// 	if !ok {
// 		panic("Missing Environment Variable: CI_PASSWORD")
// 	}
// 	loginType := os.Getenv("CI_LOGIN_TYPE")
// 	if loginType == "" {
// 		loginType = "normal"
// 	}
// 	dProjID, ok := os.LookupEnv("CI_DUMMY_PROJECT_ID")
// 	if !ok {
// 		panic("Missing Environment Variable: CI_PASSWORD")
// 	}
// 	pid, err := strconv.Atoi(dProjID)
// 	if err != nil {
// 		panic("CI_DUMMY_PROJECT_ID: Invalid Project ID integer")
// 	} else {
// 		dummyProjectID = pid
// 	}

// 	// Create client
// 	client := taiga.Client{
// 		BaseURL:    url,
// 		HTTPClient: &http.Client{},
// 	}
// 	// Initialise client (authenticates to Taiga)
// 	err = client.Initialise()
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = client.AuthByCredentials(&taiga.Credentials{
// 		Type:     "normal",
// 		Username: username,
// 		Password: password,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	comTestClient = &client
// }

// func teardownClient() {
// 	comTestClient = nil
// }

// func TestClient(t *testing.T) {
// 	setupClient()

// 	var makeurltests = []struct {
// 		in  []string
// 		out string
// 	}{
// 		{[]string{"epics"}, "https://api.taiga.io/api/v1/epics"},
// 		{[]string{"epics", "5"}, "https://api.taiga.io/api/v1/epics/5"},
// 		{[]string{"epics", "bulk_create"}, "https://api.taiga.io/api/v1/epics/bulk_create"},
// 		{[]string{"epics", "attachments", "5"}, "https://api.taiga.io/api/v1/epics/attachments/5"},
// 	}

// 	for _, tt := range makeurltests {
// 		s := comTestClient.MakeURL(tt.in...)
// 		if s != tt.out {
// 			t.Errorf("got %q, want %q", s, tt.out)
// 		}
// 	}

// }
