package taigo

// StatsService is a handle to Stats operations
// -> https://taigaio.github.io/taiga-doc/dist/api.html#stats
type StatsService struct {
	client   *Client
	Endpoint string
}

// GetDiscoverStats => https://taigaio.github.io/taiga-doc/dist/api.html#discover-stats
func (s *StatsService) GetDiscoverStats() (*DiscoverStats, error) {
	url := s.client.MakeURL(s.Endpoint, "discover")
	var respDiscoverStats DiscoverStats
	_, err := s.client.Request.Get(url, &respDiscoverStats)
	if err != nil {
		return nil, err
	}
	return &respDiscoverStats, nil
}

// GetSystemStats => https://taigaio.github.io/taiga-doc/dist/api.html#system-stats
func (s *StatsService) GetSystemStats() (*SystemStats, error) {
	url := s.client.MakeURL(s.Endpoint, "system")
	var respSystemStats SystemStats
	_, err := s.client.Request.Get(url, &respSystemStats)
	if err != nil {
		return nil, err
	}
	return &respSystemStats, nil
}
