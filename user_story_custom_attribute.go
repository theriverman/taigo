package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// UserStoryCustomAttributeService is a handle to actions related to user story custom attributes.
type UserStoryCustomAttributeService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#user-story-custom-attributes-list
func (s *UserStoryCustomAttributeService) List(queryParams *ProjectIDQueryParams) ([]UserStoryCustomAttribute, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var attrs []UserStoryCustomAttribute
	_, err = s.client.Request.Get(url, &attrs)
	if err != nil {
		return nil, err
	}
	return attrs, nil
}

// Get -> https://docs.taiga.io/api.html#user-story-custom-attributes-get
func (s *UserStoryCustomAttributeService) Get(customAttributeID int) (*UserStoryCustomAttribute, error) {
	if err := requirePositiveID("customAttributeID", customAttributeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	var attr UserStoryCustomAttribute
	_, err := s.client.Request.Get(url, &attr)
	if err != nil {
		return nil, err
	}
	return &attr, nil
}

// Create -> https://docs.taiga.io/api.html#user-story-custom-attributes-create
func (s *UserStoryCustomAttributeService) Create(request *UserStoryCustomAttributeCreateRequest) (*UserStoryCustomAttribute, error) {
	if err := requireNonNil("request", request); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responseAttr UserStoryCustomAttribute
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

// Edit -> https://docs.taiga.io/api.html#user-story-custom-attributes-edit
func (s *UserStoryCustomAttributeService) Edit(customAttributeID int, request *UserStoryCustomAttributeEditRequest) (*UserStoryCustomAttribute, error) {
	if err := requireNonNil("request", request); err != nil {
		return nil, err
	}
	if err := requirePositiveID("customAttributeID", customAttributeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	var responseAttr UserStoryCustomAttribute
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

// Patch sends an explicit PATCH payload to edit a user story custom attribute.
func (s *UserStoryCustomAttributeService) Patch(customAttributeID int, patch *UserStoryCustomAttributePatch) (*UserStoryCustomAttribute, error) {
	if err := requireNonNil("patch", patch); err != nil {
		return nil, err
	}
	if err := requirePositiveID("customAttributeID", customAttributeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	var responseAttr UserStoryCustomAttribute
	_, err := s.client.Request.Patch(url, patch, &responseAttr)
	if err != nil {
		return nil, err
	}
	return &responseAttr, nil
}

// Delete -> https://docs.taiga.io/api.html#user-story-custom-attributes-delete
func (s *UserStoryCustomAttributeService) Delete(customAttributeID int) (*http.Response, error) {
	if err := requirePositiveID("customAttributeID", customAttributeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	return s.client.Request.Delete(url)
}

// Update is an alias for Edit.
func (s *UserStoryCustomAttributeService) Update(customAttributeID int, request *UserStoryCustomAttributeEditRequest) (*UserStoryCustomAttribute, error) {
	return s.Edit(customAttributeID, request)
}
