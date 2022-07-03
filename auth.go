package taigo

import (
	"log"
)

// AuthService is a handle to actions related to Auths
//
// https://taigaio.github.io/taiga-doc/dist/api.html#auths
type AuthService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// RefreshAuthToken => https://docs.taiga.io/api.html#auth-refresh
//
//   Generates a new pair of bearer and refresh token
//   If `selfUpdate` is true, `*Client` is refreshed with the returned token values
func (s *AuthService) RefreshAuthToken(selfUpdate bool) (RefreshResponse *RefreshToken, err error) {
	url := s.client.MakeURL(s.Endpoint, "refresh")
	data := RefreshToken{
		AuthToken: s.client.Token,
		Refresh:   s.client.RefreshToken,
	}
	_, err = s.client.Request.Post(url, &data, &RefreshResponse)
	if err != nil {
		return nil, err
	}
	if selfUpdate {
		s.client.Token = RefreshResponse.AuthToken
		s.client.RefreshToken = RefreshResponse.Refresh
	}
	return
}

// PublicRegistry => https://taigaio.github.io/taiga-doc/dist/api.html#auth-public-registry
/*
	type with value "public"
	username (required)
	password (required)
	email (required)
	full_name (required)
	accepted_terms (required): boolean
*/
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
	s.client.RefreshToken = u.Refresh
	s.client.setToken()
	return &u, nil
}
