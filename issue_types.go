package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// IssueType -> https://docs.taiga.io/api.html#issue-types
type IssueType struct {
	Color     string `json:"color,omitempty"`
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Order     int    `json:"order,omitempty"`
	ProjectID int    `json:"project_id,omitempty"`
}

// IssueTypeCreateRequest represents payload for creating issue types.
type IssueTypeCreateRequest struct {
	Color   string `json:"color,omitempty"`
	Name    string `json:"name"`
	Order   int    `json:"order,omitempty"`
	Project int    `json:"project"`
}

// IssueTypeEditRequest represents sparse non-destructive updates for issue types.
type IssueTypeEditRequest struct {
	Color   string `json:"color,omitempty"`
	Name    string `json:"name,omitempty"`
	Order   int    `json:"order,omitempty"`
	Project int    `json:"project,omitempty"`
}

// IssueTypePatch represents explicit PATCH payload for issue types.
type IssueTypePatch struct {
	Color   *string `json:"color,omitempty"`
	Name    *string `json:"name,omitempty"`
	Order   *int    `json:"order,omitempty"`
	Project *int    `json:"project,omitempty"`
}

// IssueTypeService is a handle to actions related to issue types.
type IssueTypeService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#issue-types-list
func (s *IssueTypeService) List(queryParams *ProjectIDQueryParams) ([]IssueType, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var issueTypes []IssueType
	_, err = s.client.Request.Get(url, &issueTypes)
	if err != nil {
		return nil, err
	}
	return issueTypes, nil
}

// Get -> https://docs.taiga.io/api.html#issue-types-get
func (s *IssueTypeService) Get(issueTypeID int) (*IssueType, error) {
	if err := requirePositiveID("issueTypeID", issueTypeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(issueTypeID))
	var issueType IssueType
	_, err := s.client.Request.Get(url, &issueType)
	if err != nil {
		return nil, err
	}
	return &issueType, nil
}

// Create -> https://docs.taiga.io/api.html#issue-types-create
func (s *IssueTypeService) Create(request *IssueTypeCreateRequest) (*IssueType, error) {
	if err := requireNonNil("request", request); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responseIssueType IssueType
	projectID, err := resolveProjectID(request.Project, s.defaultProjectID, "project")
	if err != nil {
		return nil, err
	}
	if isEmpty(request.Name) {
		return nil, errors.New("a mandatory field(project, name) is missing. See API documentation")
	}
	payload := *request
	payload.Project = projectID
	_, err = s.client.Request.Post(url, &payload, &responseIssueType)
	if err != nil {
		return nil, err
	}
	return &responseIssueType, nil
}

// Edit -> https://docs.taiga.io/api.html#issue-types-edit
func (s *IssueTypeService) Edit(issueTypeID int, request *IssueTypeEditRequest) (*IssueType, error) {
	if err := requireNonNil("request", request); err != nil {
		return nil, err
	}
	if err := requirePositiveID("issueTypeID", issueTypeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(issueTypeID))
	var responseIssueType IssueType
	payload, err := sparsePatchMapFromStruct(request)
	if err != nil {
		return nil, err
	}
	_, err = s.client.Request.Patch(url, &payload, &responseIssueType)
	if err != nil {
		return nil, err
	}
	return &responseIssueType, nil
}

// Update is an alias for Edit.
func (s *IssueTypeService) Update(issueTypeID int, request *IssueTypeEditRequest) (*IssueType, error) {
	return s.Edit(issueTypeID, request)
}

// Patch sends an explicit PATCH payload to edit an issue type.
func (s *IssueTypeService) Patch(issueTypeID int, patch *IssueTypePatch) (*IssueType, error) {
	if err := requireNonNil("patch", patch); err != nil {
		return nil, err
	}
	if err := requirePositiveID("issueTypeID", issueTypeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(issueTypeID))
	var responseIssueType IssueType
	_, err := s.client.Request.Patch(url, patch, &responseIssueType)
	if err != nil {
		return nil, err
	}
	return &responseIssueType, nil
}

// Delete -> https://docs.taiga.io/api.html#issue-types-delete
func (s *IssueTypeService) Delete(issueTypeID int) (*http.Response, error) {
	if err := requirePositiveID("issueTypeID", issueTypeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(issueTypeID))
	return s.client.Request.Delete(url)
}
