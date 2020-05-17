package taigo

import "fmt"

var statsURI = "/stats"

// GetDiscoverStats => https://taigaio.github.io/taiga-doc/dist/api.html#discover-stats
func (c *Client) GetDiscoverStats() (*DiscoverStats, error) {
	url := c.APIURL + fmt.Sprintf("%s/discover", statsURI)
	var respDiscoverStats DiscoverStats
	err := getRequest(c, &respDiscoverStats, url)
	if err != nil {
		return nil, err
	}
	return &respDiscoverStats, nil
}

// GetSystemStats => https://taigaio.github.io/taiga-doc/dist/api.html#system-stats
func (c *Client) GetSystemStats() (*SystemStats, error) {
	url := c.APIURL + fmt.Sprintf("%s/system", statsURI)
	var respSystemStats SystemStats
	err := getRequest(c, &respSystemStats, url)
	if err != nil {
		return nil, err
	}
	return &respSystemStats, nil
}
