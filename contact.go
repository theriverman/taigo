package taigo

// ContactService is a handle to project contact requests.
type ContactService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// Send -> https://docs.taiga.io/api.html#contact
func (s *ContactService) Send(payload any) (*RawResource, error) {
	return postRawResourceAtPath(s.client, payload, s.Endpoint)
}
