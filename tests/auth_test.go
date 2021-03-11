package main

import (
	"testing"

	taiga "github.com/theriverman/taigo"
)

func TestAuth(t *testing.T) {
	setupClient()
	t.Cleanup(teardownClient)

	// Test Public Registry
	randomString := randomString(12)
	username := "test_" + randomString
	credentials := taiga.Credentials{
		Type:          "public",
		Username:      username,
		Password:      "test1",
		Email:         username + "@taigo.com",
		FullName:      "Taigo Test User",
		AcceptedTerms: true,
	}
	userAuthDetail, err := Client.Auth.PublicRegistry(&credentials)
	if err != nil {
		t.Error(err)
	}
	if userAuthDetail.Username != username {
		t.Errorf("got %q, want %q", userAuthDetail.Username, username)
	}

}
