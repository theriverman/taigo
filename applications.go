package taigo

import "strconv"

// Application is a raw DTO for /applications endpoints.
type Application = RawResource

// ApplicationService is a handle to actions related to external applications.
type ApplicationService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// Get -> https://docs.taiga.io/api.html#applications-get
func (s *ApplicationService) Get(applicationID int) (*Application, error) {
	return getRawResource(s.client, s.Endpoint, applicationID)
}

// GetToken retrieves token data for a specific application.
// https://docs.taiga.io/api.html#applications-token
func (s *ApplicationService) GetToken(applicationID int) (*RawResource, error) {
	return getRawResourceAtPath(s.client, s.Endpoint, strconv.Itoa(applicationID), "token")
}
