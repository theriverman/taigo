package taigo

import "net/http"

// ApplicationToken is a raw DTO for /application-tokens endpoints.
type ApplicationToken = RawResource

// ApplicationTokenService is a handle to actions related to application tokens.
type ApplicationTokenService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#application-tokens-list
func (s *ApplicationTokenService) List() ([]ApplicationToken, error) {
	return listRawResources(s.client, s.Endpoint, 0, nil)
}

// Get -> https://docs.taiga.io/api.html#application-tokens-get
func (s *ApplicationTokenService) Get(tokenID int) (*ApplicationToken, error) {
	return getRawResource(s.client, s.Endpoint, tokenID)
}

// Delete -> https://docs.taiga.io/api.html#application-tokens-delete
func (s *ApplicationTokenService) Delete(tokenID int) (*http.Response, error) {
	return deleteRawResource(s.client, s.Endpoint, tokenID)
}

// Authorize -> https://docs.taiga.io/api.html#application-tokens-authorize
func (s *ApplicationTokenService) Authorize(payload any) (*RawResource, error) {
	return postRawResourceAtPath(s.client, payload, s.Endpoint, "authorize")
}

// Validate -> https://docs.taiga.io/api.html#application-tokens-validate
func (s *ApplicationTokenService) Validate(payload any) (*RawResource, error) {
	return postRawResourceAtPath(s.client, payload, s.Endpoint, "validate")
}
