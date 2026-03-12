package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// IssueType -> https://docs.taiga.io/api.html#issue-types
type IssueType struct {
	Color   string `json:"color,omitempty"`
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Order   int    `json:"order,omitempty"`
	Project int    `json:"project,omitempty"`
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
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(issueTypeID))
	var issueType IssueType
	_, err := s.client.Request.Get(url, &issueType)
	if err != nil {
		return nil, err
	}
	return &issueType, nil
}

// Create -> https://docs.taiga.io/api.html#issue-types-create
func (s *IssueTypeService) Create(issueType *IssueType) (*IssueType, error) {
	if err := requireNonNil("issueType", issueType); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responseIssueType IssueType
	if isEmpty(issueType.Project) || isEmpty(issueType.Name) {
		return nil, errors.New("a mandatory field(project, name) is missing. See API documentation")
	}
	_, err := s.client.Request.Post(url, issueType, &responseIssueType)
	if err != nil {
		return nil, err
	}
	return &responseIssueType, nil
}

// Edit -> https://docs.taiga.io/api.html#issue-types-edit
func (s *IssueTypeService) Edit(issueType *IssueType) (*IssueType, error) {
	if err := requireNonNil("issueType", issueType); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(issueType.ID))
	var responseIssueType IssueType
	if issueType.ID == 0 {
		return nil, errors.New("passed IssueType does not have an ID yet. Does it exist?")
	}
	_, err := s.client.Request.Patch(url, issueType, &responseIssueType)
	if err != nil {
		return nil, err
	}
	return &responseIssueType, nil
}

// Delete -> https://docs.taiga.io/api.html#issue-types-delete
func (s *IssueTypeService) Delete(issueTypeID int) (*http.Response, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(issueTypeID))
	return s.client.Request.Delete(url)
}
