package taigo

import "log"

// AuthService is a handle to actions related to Auths
//
// https://taigaio.github.io/taiga-doc/dist/api.html#auths
type AuthService struct {
	client   *Client
	Endpoint string
}

// PublicRegistry => https://taigaio.github.io/taiga-doc/dist/api.html#auth-public-registry
func (s *AuthService) PublicRegistry(credentials *Credentials) (*UserAuthenticationDetail, error) {
	url := s.client.MakeURL(s.Endpoint, "register")
	u := UserAuthenticationDetail{}

	credentials.Type = "public"
	credentials.AcceptedTerms = true // Hardcoded for simplicity; otherwise this func would be useless
	_, err := s.client.Request.Post(url, &credentials, &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// PrivateRegistry => https://taigaio.github.io/taiga-doc/dist/api.html#auth-private-registry
// TODO: TO BE IMPLEMENTED
// func (s *AuthService) PrivateRegistry(credentials *Credentials) {}

// login authenticates to Taiga
func (s *AuthService) login(credentials *Credentials) (*UserAuthenticationDetail, error) {
	url := s.client.MakeURL(s.Endpoint)
	u := UserAuthenticationDetail{}

	_, err := s.client.Request.Post(url, &credentials, &u)
	if err != nil {
		log.Println("Failed to authenticate to Taiga.")
		return nil, err
	}
	s.client.Token = u.AuthToken
	s.client.setToken()
	return &u, nil
}
