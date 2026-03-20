package taigo

// SearchResult is a raw DTO for /search endpoints.
type SearchResult = RawResource

// SearchQueryParams holds fields to be used as URL query parameters for search.
type SearchQueryParams struct {
	Project int    `url:"project,omitempty"`
	Text    string `url:"text,omitempty"`
}

// SearchService is a handle to actions related to search.
type SearchService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// Search -> https://docs.taiga.io/api.html#searches
func (s *SearchService) Search(queryParams *SearchQueryParams) ([]SearchResult, error) {
	return listRawResources(s.client, s.Endpoint, s.defaultProjectID, queryParams)
}

// List is an alias for Search.
func (s *SearchService) List(queryParams *SearchQueryParams) ([]SearchResult, error) {
	return s.Search(queryParams)
}
