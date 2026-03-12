package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// Priority -> https://docs.taiga.io/api.html#priorities
type Priority struct {
	Color   string `json:"color,omitempty"`
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Order   int    `json:"order,omitempty"`
	Project int    `json:"project,omitempty"`
}

// PriorityService is a handle to actions related to priorities.
type PriorityService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#priorities-list
func (s *PriorityService) List(queryParams *ProjectIDQueryParams) ([]Priority, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var priorities []Priority
	_, err = s.client.Request.Get(url, &priorities)
	if err != nil {
		return nil, err
	}
	return priorities, nil
}

// Get -> https://docs.taiga.io/api.html#priorities-get
func (s *PriorityService) Get(priorityID int) (*Priority, error) {
	if err := requirePositiveID("priorityID", priorityID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(priorityID))
	var priority Priority
	_, err := s.client.Request.Get(url, &priority)
	if err != nil {
		return nil, err
	}
	return &priority, nil
}

// Create -> https://docs.taiga.io/api.html#priorities-create
func (s *PriorityService) Create(priority *Priority) (*Priority, error) {
	if err := requireNonNil("priority", priority); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responsePriority Priority
	projectID, err := resolveProjectID(priority.Project, s.defaultProjectID, "project")
	if err != nil {
		return nil, err
	}
	if isEmpty(priority.Name) {
		return nil, errors.New("a mandatory field(project, name) is missing. See API documentation")
	}
	payload := *priority
	payload.Project = projectID
	_, err = s.client.Request.Post(url, &payload, &responsePriority)
	if err != nil {
		return nil, err
	}
	return &responsePriority, nil
}

// Edit -> https://docs.taiga.io/api.html#priorities-edit
func (s *PriorityService) Edit(priority *Priority) (*Priority, error) {
	if err := requireNonNil("priority", priority); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(priority.ID))
	var responsePriority Priority
	if err := requirePositiveID("priorityID", priority.ID); err != nil {
		return nil, err
	}
	payload, err := sparsePatchMapFromStruct(priority, "id")
	if err != nil {
		return nil, err
	}
	_, err = s.client.Request.Patch(url, &payload, &responsePriority)
	if err != nil {
		return nil, err
	}
	return &responsePriority, nil
}

// Delete -> https://docs.taiga.io/api.html#priorities-delete
func (s *PriorityService) Delete(priorityID int) (*http.Response, error) {
	if err := requirePositiveID("priorityID", priorityID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(priorityID))
	return s.client.Request.Delete(url)
}
