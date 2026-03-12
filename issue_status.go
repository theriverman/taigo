package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// IssueStatusService is a handle to actions related to issue statuses.
type IssueStatusService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#issue-statuses-list
func (s *IssueStatusService) List(queryParams *ProjectIDQueryParams) ([]IssueStatus, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var statuses []IssueStatus
	_, err = s.client.Request.Get(url, &statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// Get -> https://docs.taiga.io/api.html#issue-statuses-get
func (s *IssueStatusService) Get(statusID int) (*IssueStatus, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	var status IssueStatus
	_, err := s.client.Request.Get(url, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// Create -> https://docs.taiga.io/api.html#issue-statuses-create
func (s *IssueStatusService) Create(status *IssueStatus) (*IssueStatus, error) {
	if err := requireNonNil("status", status); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responseStatus IssueStatus
	if isEmpty(status.ProjectID) || isEmpty(status.Name) {
		return nil, errors.New("a mandatory field(project, name) is missing. See API documentation")
	}
	_, err := s.client.Request.Post(url, status, &responseStatus)
	if err != nil {
		return nil, err
	}
	return &responseStatus, nil
}

// Edit -> https://docs.taiga.io/api.html#issue-statuses-edit
func (s *IssueStatusService) Edit(status *IssueStatus) (*IssueStatus, error) {
	if err := requireNonNil("status", status); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(status.ID))
	var responseStatus IssueStatus
	if status.ID == 0 {
		return nil, errors.New("passed IssueStatus does not have an ID yet. Does it exist?")
	}
	_, err := s.client.Request.Patch(url, status, &responseStatus)
	if err != nil {
		return nil, err
	}
	return &responseStatus, nil
}

// Delete -> https://docs.taiga.io/api.html#issue-statuses-delete
func (s *IssueStatusService) Delete(statusID int) (*http.Response, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	return s.client.Request.Delete(url)
}
