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
func (s *IssueStatusService) Create(request *IssueStatusCreateRequest) (*IssueStatus, error) {
	if err := requireNonNil("request", request); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responseStatus IssueStatus
	projectID, err := resolveProjectID(request.Project, s.defaultProjectID, "project")
	if err != nil {
		return nil, err
	}
	if isEmpty(request.Name) {
		return nil, errors.New("a mandatory field(project, name) is missing. See API documentation")
	}
	payload := *request
	payload.Project = projectID
	_, err = s.client.Request.Post(url, &payload, &responseStatus)
	if err != nil {
		return nil, err
	}
	return &responseStatus, nil
}

// Edit -> https://docs.taiga.io/api.html#issue-statuses-edit
func (s *IssueStatusService) Edit(statusID int, request *IssueStatusEditRequest) (*IssueStatus, error) {
	if err := requireNonNil("request", request); err != nil {
		return nil, err
	}
	if err := requirePositiveID("statusID", statusID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	var responseStatus IssueStatus
	payload, err := sparsePatchMapFromStruct(request)
	if err != nil {
		return nil, err
	}
	_, err = s.client.Request.Patch(url, &payload, &responseStatus)
	if err != nil {
		return nil, err
	}
	return &responseStatus, nil
}

// Patch sends an explicit PATCH payload to edit an issue status.
func (s *IssueStatusService) Patch(statusID int, patch *IssueStatusPatch) (*IssueStatus, error) {
	if err := requireNonNil("patch", patch); err != nil {
		return nil, err
	}
	if err := requirePositiveID("statusID", statusID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	var responseStatus IssueStatus
	_, err := s.client.Request.Patch(url, patch, &responseStatus)
	if err != nil {
		return nil, err
	}
	return &responseStatus, nil
}

// Update is an alias for Edit.
func (s *IssueStatusService) Update(statusID int, request *IssueStatusEditRequest) (*IssueStatus, error) {
	return s.Edit(statusID, request)
}

// Delete -> https://docs.taiga.io/api.html#issue-statuses-delete
func (s *IssueStatusService) Delete(statusID int) (*http.Response, error) {
	if err := requirePositiveID("statusID", statusID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	return s.client.Request.Delete(url)
}
