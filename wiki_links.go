package taigo

import "net/http"

// WikiLink is a raw DTO for /wiki-links endpoints.
type WikiLink = RawResource

// WikiLinksQueryParams holds list filters for wiki links.
type WikiLinksQueryParams struct {
	Project  int `url:"project,omitempty"`
	WikiPage int `url:"wiki_page,omitempty"`
}

// WikiLinkService is a handle to actions related to wiki links.
type WikiLinkService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#wiki-links-list
func (s *WikiLinkService) List(queryParams *WikiLinksQueryParams) ([]WikiLink, error) {
	return listRawResources(s.client, s.Endpoint, s.defaultProjectID, queryParams)
}

// Get -> https://docs.taiga.io/api.html#wiki-links-get
func (s *WikiLinkService) Get(wikiLinkID int) (*WikiLink, error) {
	return getRawResource(s.client, s.Endpoint, wikiLinkID)
}

// Create -> https://docs.taiga.io/api.html#wiki-links-create
func (s *WikiLinkService) Create(payload any) (*WikiLink, error) {
	return createRawResource(s.client, s.Endpoint, payload)
}

// Edit -> https://docs.taiga.io/api.html#wiki-links-edit
func (s *WikiLinkService) Edit(wikiLinkID int, payload any) (*WikiLink, error) {
	return patchRawResource(s.client, s.Endpoint, wikiLinkID, payload)
}

// Update is an alias for Edit.
func (s *WikiLinkService) Update(wikiLinkID int, payload any) (*WikiLink, error) {
	return s.Edit(wikiLinkID, payload)
}

// Delete -> https://docs.taiga.io/api.html#wiki-links-delete
func (s *WikiLinkService) Delete(wikiLinkID int) (*http.Response, error) {
	return deleteRawResource(s.client, s.Endpoint, wikiLinkID)
}
