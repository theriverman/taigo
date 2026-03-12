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
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	var attr TaskCustomAttribute
	_, err := s.client.Request.Get(url, &attr)
	if err != nil {
		return nil, err
	}
	return &attr, nil
}

// Create -> https://docs.taiga.io/api.html#task-custom-attributes-create
func (s *TaskCustomAttributeService) Create(customAttribute *TaskCustomAttribute) (*TaskCustomAttribute, error) {
	if err := requireNonNil("customAttribute", customAttribute); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responseAttr TaskCustomAttribute
	if isEmpty(customAttribute.Project) || isEmpty(customAttribute.Name) || isEmpty(customAttribute.Type) {
		return nil, errors.New("a mandatory field(project, name, type) is missing. See API documentation")
	}
	_, err := s.client.Request.Post(url, customAttribute, &responseAttr)
	if err != nil {
		return nil, err
	}
	return &responseAttr, nil
}

// Edit -> https://docs.taiga.io/api.html#task-custom-attributes-edit
func (s *TaskCustomAttributeService) Edit(customAttribute *TaskCustomAttribute) (*TaskCustomAttribute, error) {
	if err := requireNonNil("customAttribute", customAttribute); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttribute.ID))
	var responseAttr TaskCustomAttribute
	if customAttribute.ID == 0 {
		return nil, errors.New("passed TaskCustomAttribute does not have an ID yet. Does it exist?")
	}
	_, err := s.client.Request.Patch(url, customAttribute, &responseAttr)
	if err != nil {
		return nil, err
	}
	return &responseAttr, nil
}

// Delete -> https://docs.taiga.io/api.html#task-custom-attributes-delete
func (s *TaskCustomAttributeService) Delete(customAttributeID int) (*http.Response, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	return s.client.Request.Delete(url)
}

// Update is an alias for Edit.
func (s *TaskCustomAttributeService) Update(customAttribute *TaskCustomAttribute) (*TaskCustomAttribute, error) {
	return s.Edit(customAttribute)
}
