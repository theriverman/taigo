package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// EpicStatusService is a handle to actions related to epic statuses.
type EpicStatusService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#epic-statuses-list
func (s *EpicStatusService) List(queryParams *ProjectIDQueryParams) ([]EpicStatus, error) {
	url := s.client.MakeURL(s.Endpoint)
	url = urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	var statuses []EpicStatus
	_, err := s.client.Request.Get(url, &statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// Get -> https://docs.taiga.io/api.html#epic-statuses-get
func (s *EpicStatusService) Get(statusID int) (*EpicStatus, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	var status EpicStatus
	_, err := s.client.Request.Get(url, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// Create -> https://docs.taiga.io/api.html#epic-statuses-create
func (s *EpicStatusService) Create(status *EpicStatus) (*EpicStatus, error) {
	if err := requireNonNil("status", status); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responseStatus EpicStatus
	if isEmpty(status.Project) || isEmpty(status.Name) {
		return nil, errors.New("a mandatory field(project, name) is missing. See API documentation")
	}
	_, err := s.client.Request.Post(url, &status, &responseStatus)
	if err != nil {
		return nil, err
	}
	return &responseStatus, nil
}

// Edit -> https://docs.taiga.io/api.html#epic-statuses-edit
func (s *EpicStatusService) Edit(status *EpicStatus) (*EpicStatus, error) {
	if err := requireNonNil("status", status); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(status.ID))
	var responseStatus EpicStatus
	if status.ID == 0 {
		return nil, errors.New("passed EpicStatus does not have an ID yet. Does it exist?")
	}
	_, err := s.client.Request.Patch(url, &status, &responseStatus)
	if err != nil {
		return nil, err
	}
	return &responseStatus, nil
}

// Delete -> https://docs.taiga.io/api.html#epic-statuses-delete
func (s *EpicStatusService) Delete(statusID int) (*http.Response, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	return s.client.Request.Delete(url)
}
