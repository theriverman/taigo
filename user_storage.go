package taigo

import "net/http"

// UserStorageEntry is a raw DTO for /user-storage endpoints.
type UserStorageEntry = RawResource

// UserStorageService is a handle to actions related to user storage.
type UserStorageService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#user-storage-list
func (s *UserStorageService) List() ([]UserStorageEntry, error) {
	return listRawResources(s.client, s.Endpoint, 0, nil)
}

// Get -> https://docs.taiga.io/api.html#user-storage-get
func (s *UserStorageService) Get(key string) (*UserStorageEntry, error) {
	return getRawResource(s.client, s.Endpoint, key)
}

// Create -> https://docs.taiga.io/api.html#user-storage-create
func (s *UserStorageService) Create(payload any) (*UserStorageEntry, error) {
	return createRawResource(s.client, s.Endpoint, payload)
}

// Edit -> https://docs.taiga.io/api.html#user-storage-edit
func (s *UserStorageService) Edit(key string, payload any) (*UserStorageEntry, error) {
	return patchRawResource(s.client, s.Endpoint, key, payload)
}

// Replace updates a key using PUT.
func (s *UserStorageService) Replace(key string, payload any) (*UserStorageEntry, error) {
	return putRawResource(s.client, s.Endpoint, key, payload)
}

// Update is an alias for Edit.
func (s *UserStorageService) Update(key string, payload any) (*UserStorageEntry, error) {
	return s.Edit(key, payload)
}

// Delete -> https://docs.taiga.io/api.html#user-storage-delete
func (s *UserStorageService) Delete(key string) (*http.Response, error) {
	return deleteRawResource(s.client, s.Endpoint, key)
}
