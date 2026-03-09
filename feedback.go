package taigo

// FeedbackService is a handle to feedback actions.
type FeedbackService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// Create -> https://docs.taiga.io/api.html#feedback
func (s *FeedbackService) Create(payload any) (*RawResource, error) {
	return postRawResourceAtPath(s.client, payload, s.Endpoint)
}
