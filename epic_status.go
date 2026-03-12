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
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var statuses []EpicStatus
	_, err = s.client.Request.Get(url, &statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// Get -> https://docs.taiga.io/api.html#epic-statuses-get
func (s *EpicStatusService) Get(statusID int) (*EpicStatus, error) {
	if err := requirePositiveID("statusID", statusID); err != nil {
		return nil, err
	}
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
	projectID, err := resolveProjectID(status.Project, s.defaultProjectID, "project")
	if err != nil {
		return nil, err
	}
	if isEmpty(status.Name) {
		return nil, errors.New("a mandatory field(project, name) is missing. See API documentation")
	}
	payload := *status
	payload.Project = projectID
	_, err = s.client.Request.Post(url, &payload, &responseStatus)
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

// Delete -> https://docs.taiga.io/api.html#epic-statuses-delete
func (s *EpicStatusService) Delete(statusID int) (*http.Response, error) {
	if err := requirePositiveID("statusID", statusID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(statusID))
	return s.client.Request.Delete(url)
}
