package main

import (
	"net/http"
	"os"
	"strings"
	"testing"

	taiga "github.com/theriverman/taigo/v2"
)

func TestAuthMatrixLive(t *testing.T) {
	setupClient(t)
	t.Cleanup(teardownClient)

	unauthenticated := &taiga.Client{
		BaseURL:             testHostURL,
		HTTPClient:          &http.Client{},
		AutoRefreshDisabled: true,
	}
	if err := unauthenticated.Initialise(); err != nil {
		t.Fatalf("initialise unauthenticated client failed: %v", err)
	}

	if _, err := unauthenticated.User.Me(); err == nil {
		t.Fatalf("expected unauthenticated /users/me call to fail")
	} else if err := assertAPIErrorStatusIn(err, http.StatusUnauthorized, http.StatusForbidden); err != nil {
		t.Fatalf("unexpected unauthenticated error: %v", err)
	}

	if _, err := Client.User.Me(); err != nil {
		t.Fatalf("authenticated /users/me call failed: %v", err)
	}
}

func TestRoleMatrixLive(t *testing.T) {
	setupClient(t)
	t.Cleanup(teardownClient)

	memberUsername := strings.TrimSpace(os.Getenv("TAIGO_MEMBER_USERNAME"))
	memberPassword := strings.TrimSpace(os.Getenv("TAIGO_MEMBER_PASSWORD"))
	if memberUsername == "" || memberPassword == "" {
		t.Skip("set TAIGO_MEMBER_USERNAME and TAIGO_MEMBER_PASSWORD to run role matrix tests")
	}

	memberClient := &taiga.Client{
		BaseURL:             testHostURL,
		HTTPClient:          &http.Client{},
		AutoRefreshDisabled: true,
	}
	if err := memberClient.Initialise(); err != nil {
		t.Fatalf("initialise member client failed: %v", err)
	}
	if err := memberClient.AuthByCredentials(&taiga.Credentials{
		Type:     "normal",
		Username: memberUsername,
		Password: memberPassword,
	}); err != nil {
		t.Fatalf("member authentication failed: %v", err)
	}

	if _, err := memberClient.User.Me(); err != nil {
		t.Fatalf("member /users/me call failed: %v", err)
	}

	writeExpectation := strings.ToLower(strings.TrimSpace(os.Getenv("TAIGO_MEMBER_WRITE_EXPECTATION")))
	if writeExpectation == "" {
		writeExpectation = "forbid"
	}
	if writeExpectation != "allow" && writeExpectation != "forbid" {
		t.Fatalf("invalid TAIGO_MEMBER_WRITE_EXPECTATION value %q (use allow|forbid)", writeExpectation)
	}

	webhook, err := memberClient.Webhook.Create(&taiga.Webhook{
		Project: testProjID,
		Name:    "role-matrix-" + strings.ToLower(RandStringBytesMaskImprSrcUnsafe(8)),
		Key:     "role-key-" + strings.ToLower(RandStringBytesMaskImprSrcUnsafe(8)),
		URL:     "https://example.com/role-matrix",
	})
	switch writeExpectation {
	case "allow":
		if err != nil {
			t.Fatalf("expected member write action to be allowed, got error: %v", err)
		}
		if webhook != nil {
			t.Cleanup(func() {
				_ = memberClient.Webhook.Delete(webhook.ID)
			})
		}
	case "forbid":
		if err == nil {
			if webhook != nil {
				_ = memberClient.Webhook.Delete(webhook.ID)
			}
			t.Fatalf("expected member write action to be forbidden, but webhook creation succeeded")
		}
		if err := assertAPIErrorStatusIn(err, http.StatusForbidden, http.StatusNotFound); err != nil {
			t.Fatalf("expected forbidden/not-found status for member write action: %v", err)
		}
	}
}
