package taigo

// ImporterService is a handle to provider importer actions.
type ImporterService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// ImporterAuthURLQueryParams holds query params for importer auth_url endpoints.
type ImporterAuthURLQueryParams struct {
	Project int `url:"project,omitempty"`
}

// TrelloAuthURL -> https://docs.taiga.io/api.html#importers-trello-auth-url
func (s *ImporterService) TrelloAuthURL(queryParams *ImporterAuthURLQueryParams) (*RawResource, error) {
	return getRawResourceAtPathWithQuery(s.client, queryParams, s.Endpoint, "trello", "auth_url")
}

// TrelloListProjects -> https://docs.taiga.io/api.html#importers-trello-list-projects
func (s *ImporterService) TrelloListProjects(queryParams any) (*RawResource, error) {
	return getRawResourceAtPathWithQuery(s.client, queryParams, s.Endpoint, "trello", "list_projects")
}

// TrelloImportProject -> https://docs.taiga.io/api.html#importers-trello-import-project
func (s *ImporterService) TrelloImportProject(payload any) (*RawResource, error) {
	return postRawResourceAtPath(s.client, payload, s.Endpoint, "trello", "import_project")
}

// TrelloLoadData is kept as a backward-compatible alias for TrelloImportProject.
func (s *ImporterService) TrelloLoadData(payload any) (*RawResource, error) {
	return s.TrelloImportProject(payload)
}

// TrelloAuthorize -> https://docs.taiga.io/api.html#importers-trello-authorize
func (s *ImporterService) TrelloAuthorize(payload any) (*RawResource, error) {
	return postRawResourceAtPath(s.client, payload, s.Endpoint, "trello", "authorize")
}

// GithubAuthURL -> https://docs.taiga.io/api.html#importers-github-auth-url
func (s *ImporterService) GithubAuthURL(queryParams *ImporterAuthURLQueryParams) (*RawResource, error) {
	return getRawResourceAtPathWithQuery(s.client, queryParams, s.Endpoint, "github", "auth_url")
}

// GithubListProjects -> https://docs.taiga.io/api.html#importers-github-list-projects
func (s *ImporterService) GithubListProjects(queryParams any) (*RawResource, error) {
	return getRawResourceAtPathWithQuery(s.client, queryParams, s.Endpoint, "github", "list_projects")
}

// GithubListUsers -> https://docs.taiga.io/api.html#importers-github-list-users
func (s *ImporterService) GithubListUsers(queryParams any) (*RawResource, error) {
	return getRawResourceAtPathWithQuery(s.client, queryParams, s.Endpoint, "github", "list_users")
}

// GithubImportProject -> https://docs.taiga.io/api.html#importers-github-import-project
func (s *ImporterService) GithubImportProject(payload any) (*RawResource, error) {
	return postRawResourceAtPath(s.client, payload, s.Endpoint, "github", "import_project")
}

// GithubLoadData is kept as a backward-compatible alias for GithubImportProject.
func (s *ImporterService) GithubLoadData(payload any) (*RawResource, error) {
	return s.GithubImportProject(payload)
}

// GithubAuthorize -> https://docs.taiga.io/api.html#importers-github-authorize
func (s *ImporterService) GithubAuthorize(payload any) (*RawResource, error) {
	return postRawResourceAtPath(s.client, payload, s.Endpoint, "github", "authorize")
}

// JiraAuthURL -> https://docs.taiga.io/api.html#importers-jira-auth-url
func (s *ImporterService) JiraAuthURL(queryParams *ImporterAuthURLQueryParams) (*RawResource, error) {
	return getRawResourceAtPathWithQuery(s.client, queryParams, s.Endpoint, "jira", "auth_url")
}

// JiraListProjects -> https://docs.taiga.io/api.html#importers-jira-list-projects
func (s *ImporterService) JiraListProjects(queryParams any) (*RawResource, error) {
	return getRawResourceAtPathWithQuery(s.client, queryParams, s.Endpoint, "jira", "list_projects")
}

// JiraImportProject -> https://docs.taiga.io/api.html#importers-jira-import-project
func (s *ImporterService) JiraImportProject(payload any) (*RawResource, error) {
	return postRawResourceAtPath(s.client, payload, s.Endpoint, "jira", "import_project")
}

// JiraLoadData is kept as a backward-compatible alias for JiraImportProject.
func (s *ImporterService) JiraLoadData(payload any) (*RawResource, error) {
	return s.JiraImportProject(payload)
}

// JiraAuthorize -> https://docs.taiga.io/api.html#importers-jira-authorize
func (s *ImporterService) JiraAuthorize(payload any) (*RawResource, error) {
	return postRawResourceAtPath(s.client, payload, s.Endpoint, "jira", "authorize")
}
