package taigo

// Credentials is the payload for normal authentication
type Credentials struct {
	Type          string `json:"type,omitempty"` // normal;github;ldap;public;private
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Code          string `json:"code,omitempty"` // GitHub Authentication Code
	Email         string `json:"email,omitempty"`
	Existing      bool   `json:"existing,omitempty"`
	FullName      string `json:"full_name,omitempty"`
	Token         string `json:"token,omitempty"`
	AcceptedTerms bool   `json:"accepted_terms,omitempty"` // Required for registration only
}

// UserAuthenticationDetail is a superset of User extended by an AuthToken field
type UserAuthenticationDetail struct {
	AuthToken string `json:"auth_token"`
	Refresh   string `json:"refresh"`
	User             // Embedding type User struct
}

// AsUser returns a *User from *UserAuthenticationDetail
// 	The AuthToken can be accessed from `User` via `.GetToken()`
func (u *UserAuthenticationDetail) AsUser() *User {
	user := &User{}
	err := convertStructViaJSON(u, user)
	if err != nil {
		return nil
	}
	user.authToken = u.AuthToken
	return user
}

type RefreshToken struct {
	AuthToken string `json:"auth_token,omitempty"`
	Refresh   string `json:"refresh,omitempty"`
}
