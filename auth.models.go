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
