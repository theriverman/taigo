package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// TaskStatusService is a handle to actions related to task statuses.
type TaskStatusService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#task-statuses-list
func (s *TaskStatusService) List(queryParams *ProjectIDQueryParams) ([]TaskStatus, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var statuses []TaskStatus
	_, err = s.client.Request.Get(url, &statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// Get -> https://docs.taiga.io/api.html#task-statuses-get
func (s *TaskStatusService) Get(statusID int) (*TaskStatus, error) {
	if err := requirePositiveID("statusID", statusID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	var status TaskStatus
	_, err := s.client.Request.Get(url, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// Create -> https://docs.taiga.io/api.html#task-statuses-create
func (s *TaskStatusService) Create(request *TaskStatusCreateRequest) (*TaskStatus, error) {
	if err := requireNonNil("request", request); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responseStatus TaskStatus
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

// Edit -> https://docs.taiga.io/api.html#task-statuses-edit
func (s *TaskStatusService) Edit(statusID int, request *TaskStatusEditRequest) (*TaskStatus, error) {
	if err := requireNonNil("request", request); err != nil {
		return nil, err
	}
	if err := requirePositiveID("statusID", statusID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	var responseStatus TaskStatus
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

// Patch sends an explicit PATCH payload to edit a task status.
func (s *TaskStatusService) Patch(statusID int, patch *TaskStatusPatch) (*TaskStatus, error) {
	if err := requireNonNil("patch", patch); err != nil {
		return nil, err
	}
	if err := requirePositiveID("statusID", statusID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	var responseStatus TaskStatus
	_, err := s.client.Request.Patch(url, patch, &responseStatus)
	if err != nil {
		return nil, err
	}
	return &responseStatus, nil
}

// Update is an alias for Edit.
func (s *TaskStatusService) Update(statusID int, request *TaskStatusEditRequest) (*TaskStatus, error) {
	return s.Edit(statusID, request)
}

// Delete -> https://docs.taiga.io/api.html#task-statuses-delete
func (s *TaskStatusService) Delete(statusID int) (*http.Response, error) {
	if err := requirePositiveID("statusID", statusID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	return s.client.Request.Delete(url)
}
