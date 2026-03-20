package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// TaskCustomAttributeService is a handle to actions related to task custom attributes.
type TaskCustomAttributeService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#task-custom-attributes-list
func (s *TaskCustomAttributeService) List(queryParams *ProjectIDQueryParams) ([]TaskCustomAttribute, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var attrs []TaskCustomAttribute
	_, err = s.client.Request.Get(url, &attrs)
	if err != nil {
		return nil, err
	}
	return attrs, nil
}

// Get -> https://docs.taiga.io/api.html#task-custom-attributes-get
func (s *TaskCustomAttributeService) Get(customAttributeID int) (*TaskCustomAttribute, error) {
	if err := requirePositiveID("customAttributeID", customAttributeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	var attr TaskCustomAttribute
	_, err := s.client.Request.Get(url, &attr)
	if err != nil {
		return nil, err
	}
	return &attr, nil
}

// Create -> https://docs.taiga.io/api.html#task-custom-attributes-create
func (s *TaskCustomAttributeService) Create(request *TaskCustomAttributeCreateRequest) (*TaskCustomAttribute, error) {
	if err := requireNonNil("request", request); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responseAttr TaskCustomAttribute
	projectID, err := resolveProjectID(request.Project, s.defaultProjectID, "project")
	if err != nil {
		return nil, err
	}
	if isEmpty(request.Name) || isEmpty(request.Type) {
		return nil, errors.New("a mandatory field(project, name, type) is missing. See API documentation")
	}
	payload := *request
	payload.Project = projectID
	_, err = s.client.Request.Post(url, &payload, &responseAttr)
	if err != nil {
		return nil, err
	}
	return &responseAttr, nil
}

// Edit -> https://docs.taiga.io/api.html#task-custom-attributes-edit
func (s *TaskCustomAttributeService) Edit(customAttributeID int, request *TaskCustomAttributeEditRequest) (*TaskCustomAttribute, error) {
	if err := requireNonNil("request", request); err != nil {
		return nil, err
	}
	if err := requirePositiveID("customAttributeID", customAttributeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	var responseAttr TaskCustomAttribute
	payload, err := sparsePatchMapFromStruct(request)
	if err != nil {
		return nil, err
	}
	_, err = s.client.Request.Patch(url, &payload, &responseAttr)
	if err != nil {
		return nil, err
	}
	return &responseAttr, nil
}

// Patch sends an explicit PATCH payload to edit a task custom attribute.
func (s *TaskCustomAttributeService) Patch(customAttributeID int, patch *TaskCustomAttributePatch) (*TaskCustomAttribute, error) {
	if err := requireNonNil("patch", patch); err != nil {
		return nil, err
	}
	if err := requirePositiveID("customAttributeID", customAttributeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	var responseAttr TaskCustomAttribute
	_, err := s.client.Request.Patch(url, patch, &responseAttr)
	if err != nil {
		return nil, err
	}
	return &responseAttr, nil
}

// Delete -> https://docs.taiga.io/api.html#task-custom-attributes-delete
func (s *TaskCustomAttributeService) Delete(customAttributeID int) (*http.Response, error) {
	if err := requirePositiveID("customAttributeID", customAttributeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	return s.client.Request.Delete(url)
}

// Update is an alias for Edit.
func (s *TaskCustomAttributeService) Update(customAttributeID int, request *TaskCustomAttributeEditRequest) (*TaskCustomAttribute, error) {
	return s.Edit(customAttributeID, request)
}
