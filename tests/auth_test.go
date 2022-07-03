package main

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	taiga "github.com/theriverman/taigo"
)

func TestAuth(t *testing.T) {
	setupClient()
	t.Cleanup(teardownClient)

	// Test Public Registry
	randomString := RandStringBytesMaskImprSrcUnsafe(12)
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

func TestAuthService_RefreshAuthToken(t *testing.T) {
	setupClient()
	t.Cleanup(teardownClient)

	type fields struct {
		client           *taiga.Client
		defaultProjectID int
		Endpoint         string
	}
	type args struct {
		selfUpdate bool
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		oldClientToken string
		wantErr        bool
	}{
		{
			name: "Succesful manual token refresh",
			fields: fields{
				client:           Client,
				defaultProjectID: 0,
				Endpoint:         "auth",
			},
			args:           args{selfUpdate: true},
			oldClientToken: Client.Token,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Client.Auth.RefreshAuthToken(tt.args.selfUpdate)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.RefreshAuthToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.DeepEqual(tt.oldClientToken, Client.Token) {
				t.Errorf("Value of Client.Token did not change which means that .RefreshAuthToken(true) has failed")
				t.Errorf("AuthService.RefreshAuthToken().AuthToken = %v, want different than %v", Client.Token, tt.oldClientToken)
			}
		})
	}
}

func TestTokenRefreshRoutine(t *testing.T) {
	// we need a custom client here to set `AutoRefreshTickerDuration` to 5 seconds
	// otherwise the the test would fail b/c the default ticker duration is 12hrs
	client := &taiga.Client{
		BaseURL:                   testHostURL,
		HTTPClient:                &http.Client{},
		AutoRefreshTickerDuration: 5 * time.Second,
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

	// Let's loop for 35 seconds to observe the routine working
	testLoopLength := 5
	for i := 0; i < testLoopLength; i++ {
		if i == testLoopLength {
			break
		}
		t.Logf("Loop %d | client.Token  : %s\n", i, client.Token)
		t.Logf("Loop %d | client.Refresh: %s\n", i, client.RefreshToken)
		t.Logf("----------------------------\n")
		time.Sleep(5 * time.Second)
	}
	client = nil

	// TODO: Add test evaluation

}
