package taigo

// ContribPlugin is a raw DTO for /contrib-plugins endpoints.
type ContribPlugin = RawResource

// ContribPluginService is a handle to contrib plugins discovery.
type ContribPluginService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> endpoint depends on backend settings.
func (s *ContribPluginService) List() ([]ContribPlugin, error) {
	return listRawResources(s.client, s.Endpoint, 0, nil)
}
