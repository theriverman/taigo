package gotaiga

import "log"

var endpointAuth = "/auth"

// AuthService is a handle to actions related to Auths
//
// https://taigaio.github.io/taiga-doc/dist/api.html#auths
type AuthService struct {
	client *Client
}

// PublicRegistry => https://taigaio.github.io/taiga-doc/dist/api.html#auth-public-registry
func (s *AuthService) PublicRegistry(credentials *Credentials) (*UserAuthenticationDetail, error) {
	url := s.client.APIURL + endpointAuth + "/register"
	response := UserAuthenticationDetail{}

	credentials.Type = "public"
	credentials.AcceptedTerms = true // Hardcoded for simplicity; otherwise this func would be useless
	err := postRequest(s.client, &response, url, &credentials)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PrivateRegistry => https://taigaio.github.io/taiga-doc/dist/api.html#auth-private-registry
// TODO: TO BE IMPLEMENTED
// func (s *AuthService) PrivateRegistry(credentials *Credentials) {}

// login authenticates to Taiga
func (s *AuthService) login(credentials *Credentials) (*UserAuthenticationDetail, error) {
	url := s.client.APIURL + endpointAuth
	response := UserAuthenticationDetail{}

	err := postRequest(s.client, &response, url, &credentials)
	if err != nil {
		log.Println("Failed to authenticate to Taiga.")
		return nil, err
	}
	s.client.Token = response.AuthToken
	s.client.setToken()
	s.client.IsLoggedIn = true
	return &response, nil
}
