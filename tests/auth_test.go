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
	fullName := "Taigo Test User"
	credentials := taiga.Credentials{
		Type:          "public",
		Username:      username,
		Password:      randomString,
		Email:         username + "@taigo.com",
		FullName:      fullName,
		AcceptedTerms: true,
	}
	userAuthDetail, err := Client.Auth.PublicRegistry(&credentials)
	if err != nil {
		t.Error(err)
	}
	if userAuthDetail.FullName != fullName {
		t.Errorf("got %s, want %s", userAuthDetail.FullName, fullName)
	}

}
