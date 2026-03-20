package taigo

// ObjectSummary is a generic summary payload.
type ObjectSummary = RawResource

// ObjectsSummaryQueryParams holds optional filters for object summary listing.
type ObjectsSummaryQueryParams struct {
	Project int `url:"project,omitempty"`
}

// ObjectsSummaryService is a handle to object summary endpoints.
type ObjectsSummaryService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List returns object summaries.
func (s *ObjectsSummaryService) List(queryParams *ObjectsSummaryQueryParams) ([]ObjectSummary, error) {
	return listRawResources(s.client, s.Endpoint, s.defaultProjectID, queryParams)
}
