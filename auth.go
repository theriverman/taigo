package taigo

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
//	Generates a new pair of bearer and refresh token
//	If `selfUpdate` is true, `*Client` is refreshed with the returned token values
func (s *AuthService) RefreshAuthToken(selfUpdate bool) (RefreshResponse *RefreshToken, err error) {
	url := s.client.MakeURL(s.Endpoint, "refresh")
	authToken, refreshToken := s.client.currentTokens()
	data := RefreshToken{
		AuthToken: authToken,
		Refresh:   refreshToken,
	}
	response := &RefreshToken{}
	_, err = s.client.Request.Post(url, &data, response)
	if err != nil {
		return nil, err
	}
	if selfUpdate {
		s.client.setAuthTokens("", response.AuthToken, response.Refresh)
	}
	return response, nil
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
	if err := requireNonNil("credentials", credentials); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, "register")
	u := UserAuthenticationDetail{}
	payload := *credentials
	payload.Type = "public"
	payload.AcceptedTerms = true // Hardcoded for simplicity; otherwise this func would be useless
	_, err := s.client.Request.Post(url, &payload, &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// PrivateRegistry => https://taigaio.github.io/taiga-doc/dist/api.html#auth-private-registry
func (s *AuthService) PrivateRegistry(credentials *Credentials) (*UserAuthenticationDetail, error) {
	if err := requireNonNil("credentials", credentials); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, "register")
	u := UserAuthenticationDetail{}
	payload := *credentials
	payload.Type = "private"
	if !payload.AcceptedTerms {
		payload.AcceptedTerms = true
	}
	_, err := s.client.Request.Post(url, &payload, &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// login authenticates to Taiga
func (s *AuthService) login(credentials *Credentials) (*UserAuthenticationDetail, error) {
	if err := requireNonNil("credentials", credentials); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	u := UserAuthenticationDetail{}

	_, err := s.client.Request.Post(url, credentials, &u)
	if err != nil {
		return nil, err
	}
	s.client.setAuthTokens("", u.AuthToken, u.Refresh)
	return &u, nil
}
