package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// EpicCustomAttributeService is a handle to actions related to epic custom attributes.
type EpicCustomAttributeService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#epic-custom-attributes-list
func (s *EpicCustomAttributeService) List(queryParams *ProjectIDQueryParams) ([]EpicCustomAttribute, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var attrs []EpicCustomAttribute
	_, err = s.client.Request.Get(url, &attrs)
	if err != nil {
		return nil, err
	}
	return attrs, nil
}

// Get -> https://docs.taiga.io/api.html#epic-custom-attributes-get
func (s *EpicCustomAttributeService) Get(customAttributeID int) (*EpicCustomAttribute, error) {
	if err := requirePositiveID("customAttributeID", customAttributeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	var attr EpicCustomAttribute
	_, err := s.client.Request.Get(url, &attr)
	if err != nil {
		return nil, err
	}
	return &attr, nil
}

// Create -> https://docs.taiga.io/api.html#epic-custom-attributes-create
func (s *EpicCustomAttributeService) Create(request *EpicCustomAttributeCreateRequest) (*EpicCustomAttribute, error) {
	if err := requireNonNil("request", request); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responseAttr EpicCustomAttribute
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

// Edit -> https://docs.taiga.io/api.html#epic-custom-attributes-edit
func (s *EpicCustomAttributeService) Edit(customAttributeID int, request *EpicCustomAttributeEditRequest) (*EpicCustomAttribute, error) {
	if err := requireNonNil("request", request); err != nil {
		return nil, err
	}
	if err := requirePositiveID("customAttributeID", customAttributeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	var responseAttr EpicCustomAttribute
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

// Patch sends an explicit PATCH payload to edit an epic custom attribute.
func (s *EpicCustomAttributeService) Patch(customAttributeID int, patch *EpicCustomAttributePatch) (*EpicCustomAttribute, error) {
	if err := requireNonNil("patch", patch); err != nil {
		return nil, err
	}
	if err := requirePositiveID("customAttributeID", customAttributeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	var responseAttr EpicCustomAttribute
	_, err := s.client.Request.Patch(url, patch, &responseAttr)
	if err != nil {
		return nil, err
	}
	return &responseAttr, nil
}

// Delete -> https://docs.taiga.io/api.html#epic-custom-attributes-delete
func (s *EpicCustomAttributeService) Delete(customAttributeID int) (*http.Response, error) {
	if err := requirePositiveID("customAttributeID", customAttributeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	return s.client.Request.Delete(url)
}

// Update is an alias for Edit.
func (s *EpicCustomAttributeService) Update(customAttributeID int, request *EpicCustomAttributeEditRequest) (*EpicCustomAttribute, error) {
	return s.Edit(customAttributeID, request)
}
