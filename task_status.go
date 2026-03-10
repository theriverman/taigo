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
	url = urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	var statuses []TaskStatus
	_, err := s.client.Request.Get(url, &statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// Get -> https://docs.taiga.io/api.html#task-statuses-get
func (s *TaskStatusService) Get(statusID int) (*TaskStatus, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	var status TaskStatus
	_, err := s.client.Request.Get(url, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// Create -> https://docs.taiga.io/api.html#task-statuses-create
func (s *TaskStatusService) Create(status *TaskStatus) (*TaskStatus, error) {
	if err := requireNonNil("status", status); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responseStatus TaskStatus
	if isEmpty(status.ProjectID) || isEmpty(status.Name) {
		return nil, errors.New("a mandatory field(project, name) is missing. See API documentation")
	}
	_, err := s.client.Request.Post(url, &status, &responseStatus)
	if err != nil {
		return nil, err
	}
	return &responseStatus, nil
}

// Edit -> https://docs.taiga.io/api.html#task-statuses-edit
func (s *TaskStatusService) Edit(status *TaskStatus) (*TaskStatus, error) {
	if err := requireNonNil("status", status); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(status.ID))
	var responseStatus TaskStatus
	if status.ID == 0 {
		return nil, errors.New("passed TaskStatus does not have an ID yet. Does it exist?")
	}
	_, err := s.client.Request.Patch(url, &status, &responseStatus)
	if err != nil {
		return nil, err
	}
	return &responseStatus, nil
}

// Delete -> https://docs.taiga.io/api.html#task-statuses-delete
func (s *TaskStatusService) Delete(statusID int) (*http.Response, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	return s.client.Request.Delete(url)
}
