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
	if err := requirePositiveID("statusID", statusID); err != nil {
		return nil, err
	}
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
	projectID, err := resolveProjectID(status.ProjectID, s.defaultProjectID, "project_id")
	if err != nil {
		return nil, err
	}
	if isEmpty(status.Name) {
		return nil, errors.New("a mandatory field(project, name) is missing. See API documentation")
	}
	payload := *status
	payload.ProjectID = projectID
	_, err = s.client.Request.Post(url, &payload, &responseStatus)
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
	if err := requirePositiveID("statusID", status.ID); err != nil {
		return nil, err
	}
	payload, err := sparsePatchMapFromStruct(status, "id", "slug")
	if err != nil {
		return nil, err
	}
	_, err = s.client.Request.Patch(url, &payload, &responseStatus)
	if err != nil {
		return nil, err
	}
	return &responseStatus, nil
}

// Delete -> https://docs.taiga.io/api.html#issue-statuses-delete
func (s *IssueStatusService) Delete(statusID int) (*http.Response, error) {
	if err := requirePositiveID("statusID", statusID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	return s.client.Request.Delete(url)
}
