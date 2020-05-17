package taigo

// Resolver represents all possible response body key/value pairs
type Resolver struct {
	Project   int `json:"project,omitempty"`
	UserStory int `json:"us,omitempty"`
	Issue     int `json:"issue,omitempty"`
	Task      int `json:"task,omitempty"`
	Milestone int `json:"milestone,omitempty"`
	WikiPage  int `json:"wikipage,omitempty"`
}

// ResolverQueryParams holds fields to be used as URL query parameters to filter the queried objects
type ResolverQueryParams struct {
	Project   string `url:"project,omitempty"`
	Issue     int    `url:"issue,omitempty"`
	Task      int    `url:"task,omitempty"`
	Milestone string `url:"milestone,omitempty"`
	WikiPage  string `url:"wikipage,omitempty"`
	US        int    `url:"us,omitempty"` // UserStory
}
