package main

import (
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"
	"unsafe"

	taiga "github.com/theriverman/taigo/v2"
)

const testHostURL string = "http://localhost:9000"
const testUsername string = "admin"
const testPassword string = "admin"
const testProjSlug string = "taigo-test"
const testProjID int = 2
const testUserID int = 5

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

	return *(*string)(unsafe.Pointer(&b))
}
