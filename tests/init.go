package main

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	taiga "github.com/theriverman/taigo/v2"
)

const defaultTestHostURL string = "http://localhost:9000"
const defaultTestUsername string = "admin"
const defaultTestPassword string = "123123"
const defaultTestProjSlug string = "taigo-test"
const defaultTestProjID int = 2
const defaultTestUserID int = 5

var testHostURL string = defaultTestHostURL
var testUsername string = defaultTestUsername
var testPassword string = defaultTestPassword
var testProjSlug string = defaultTestProjSlug
var testProjID int = defaultTestProjID
var testUserID int = defaultTestUserID

var integrationEnvLoaded bool

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// Client is the foundation for making requests against Taiga
var Client *taiga.Client = nil

func setupClient(t *testing.T) {
	t.Helper()
	if os.Getenv("TAIGO_RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("set TAIGO_RUN_INTEGRATION_TESTS=1 to run integration tests against a live Taiga instance")
	}

	if Client != nil {
		return // client already set; skipping
	}
	loadIntegrationEnv(t)

	// Create client
	client := taiga.Client{
		BaseURL:    testHostURL,
		HTTPClient: &http.Client{},
	}
	// Initialise client (authenticates to Taiga)
	err := client.Initialise()
	if err != nil {
		t.Skipf("skipping integration tests: could not initialise client: %v", err)
	}
	err = client.AuthByCredentials(&taiga.Credentials{
		Type:     "normal",
		Username: testUsername,
		Password: testPassword,
	})
	if err != nil {
		t.Skipf("skipping integration tests: Taiga backend unavailable at %s: %v", testHostURL, err)
	}
	Client = &client
}

func teardownClient() {
	Client = nil
}

func loadIntegrationEnv(t *testing.T) {
	t.Helper()
	if integrationEnvLoaded {
		return
	}

	if baseURL := strings.TrimSpace(os.Getenv("TAIGO_BASE_URL")); baseURL != "" {
		testHostURL = strings.TrimSuffix(baseURL, "/")
	}
	if username := strings.TrimSpace(os.Getenv("TAIGO_USERNAME")); username != "" {
		testUsername = username
	}
	if password := strings.TrimSpace(os.Getenv("TAIGO_PASSWORD")); password != "" {
		testPassword = password
	}
	if projectSlug := strings.TrimSpace(os.Getenv("TAIGO_PROJECT_SLUG")); projectSlug != "" {
		testProjSlug = projectSlug
	}
	if projectID := strings.TrimSpace(os.Getenv("TAIGO_PROJECT_ID")); projectID != "" {
		parsedID, err := strconv.Atoi(projectID)
		if err != nil {
			t.Fatalf("invalid TAIGO_PROJECT_ID value %q: %v", projectID, err)
		}
		testProjID = parsedID
	}
	if userID := strings.TrimSpace(os.Getenv("TAIGO_USER_ID")); userID != "" {
		parsedID, err := strconv.Atoi(userID)
		if err != nil {
			t.Fatalf("invalid TAIGO_USER_ID value %q: %v", userID, err)
		}
		testUserID = parsedID
	}

	integrationEnvLoaded = true
}

func RandStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	var src = rand.NewSource(time.Now().UnixNano())
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
