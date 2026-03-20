package taigo

// Locale is a raw DTO for /locales endpoint.
type Locale = RawResource

// LocaleService is a handle to locale actions.
type LocaleService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#locales
func (s *LocaleService) List() ([]Locale, error) {
	return listRawResources(s.client, s.Endpoint, 0, nil)
}
