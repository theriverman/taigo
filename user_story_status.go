package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// UserStoryStatusService is a handle to actions related to user story statuses.
type UserStoryStatusService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#user-story-statuses-list
func (s *UserStoryStatusService) List(queryParams *ProjectIDQueryParams) ([]UserStoryStatus, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var statuses []UserStoryStatus
	_, err = s.client.Request.Get(url, &statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// Get -> https://docs.taiga.io/api.html#user-story-statuses-get
func (s *UserStoryStatusService) Get(statusID int) (*UserStoryStatus, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	var status UserStoryStatus
	_, err := s.client.Request.Get(url, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

// Create -> https://docs.taiga.io/api.html#user-story-statuses-create
func (s *UserStoryStatusService) Create(status *UserStoryStatus) (*UserStoryStatus, error) {
	if err := requireNonNil("status", status); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responseStatus UserStoryStatus
	if isEmpty(status.ProjectID) || isEmpty(status.Name) {
		return nil, errors.New("a mandatory field(project_id, name) is missing. See API documentation")
	}
	_, err := s.client.Request.Post(url, status, &responseStatus)
	if err != nil {
		return nil, err
	}
	return &responseStatus, nil
}

// Edit -> https://docs.taiga.io/api.html#user-story-statuses-edit
func (s *UserStoryStatusService) Edit(status *UserStoryStatus) (*UserStoryStatus, error) {
	if err := requireNonNil("status", status); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(status.ID))
	var responseStatus UserStoryStatus
	if status.ID == 0 {
		return nil, errors.New("passed UserStoryStatus does not have an ID yet. Does it exist?")
	}
	_, err := s.client.Request.Patch(url, status, &responseStatus)
	if err != nil {
		return nil, err
	}
	return &responseStatus, nil
}

// Delete -> https://docs.taiga.io/api.html#user-story-statuses-delete
func (s *UserStoryStatusService) Delete(statusID int) (*http.Response, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	return s.client.Request.Delete(url)
}

// Update is an alias for Edit.
func (s *UserStoryStatusService) Update(status *UserStoryStatus) (*UserStoryStatus, error) {
	return s.Edit(status)
}
