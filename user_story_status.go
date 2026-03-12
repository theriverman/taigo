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
	if err := requirePositiveID("statusID", statusID); err != nil {
		return nil, err
	}
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
	projectID, err := resolveProjectID(status.ProjectID, s.defaultProjectID, "project_id")
	if err != nil {
		return nil, err
	}
	if isEmpty(status.Name) {
		return nil, errors.New("a mandatory field(project_id, name) is missing. See API documentation")
	}
	payload := *status
	payload.ProjectID = projectID
	_, err = s.client.Request.Post(url, &payload, &responseStatus)
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

// Delete -> https://docs.taiga.io/api.html#user-story-statuses-delete
func (s *UserStoryStatusService) Delete(statusID int) (*http.Response, error) {
	if err := requirePositiveID("statusID", statusID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	return s.client.Request.Delete(url)
}

// Update is an alias for Edit.
func (s *UserStoryStatusService) Update(status *UserStoryStatus) (*UserStoryStatus, error) {
	return s.Edit(status)
}
