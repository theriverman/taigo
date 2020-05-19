package taigo

var statsURI = "/stats"

// StatsService is a handle to Stats operations
// -> https://taigaio.github.io/taiga-doc/dist/api.html#stats
type StatsService struct {
	client *Client
}

// GetDiscoverStats => https://taigaio.github.io/taiga-doc/dist/api.html#discover-stats
func (s *StatsService) GetDiscoverStats() (*DiscoverStats, error) {
	url := s.client.Request.MakeURL(statsURI, "discover")
	var respDiscoverStats DiscoverStats
	err := s.client.Request.GetRequest(url, &respDiscoverStats)
	if err != nil {
		return nil, err
	}
	return &respDiscoverStats, nil
}

// GetSystemStats => https://taigaio.github.io/taiga-doc/dist/api.html#system-stats
func (s *StatsService) GetSystemStats() (*SystemStats, error) {
	url := s.client.Request.MakeURL(statsURI, "system")
	var respSystemStats SystemStats
	err := s.client.Request.GetRequest(url, &respSystemStats)
	if err != nil {
		return nil, err
	}
	return &respSystemStats, nil
}
