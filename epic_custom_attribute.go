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
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	var attr EpicCustomAttribute
	_, err := s.client.Request.Get(url, &attr)
	if err != nil {
		return nil, err
	}
	return &attr, nil
}

// Create -> https://docs.taiga.io/api.html#epic-custom-attributes-create
func (s *EpicCustomAttributeService) Create(customAttribute *EpicCustomAttribute) (*EpicCustomAttribute, error) {
	if err := requireNonNil("customAttribute", customAttribute); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responseAttr EpicCustomAttribute
	if isEmpty(customAttribute.Project) || isEmpty(customAttribute.Name) || isEmpty(customAttribute.Type) {
		return nil, errors.New("a mandatory field(project, name, type) is missing. See API documentation")
	}
	_, err := s.client.Request.Post(url, &customAttribute, &responseAttr)
	if err != nil {
		return nil, err
	}
	return &responseAttr, nil
}

// Edit -> https://docs.taiga.io/api.html#epic-custom-attributes-edit
func (s *EpicCustomAttributeService) Edit(customAttribute *EpicCustomAttribute) (*EpicCustomAttribute, error) {
	if err := requireNonNil("customAttribute", customAttribute); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttribute.ID))
	var responseAttr EpicCustomAttribute
	if customAttribute.ID == 0 {
		return nil, errors.New("passed EpicCustomAttribute does not have an ID yet. Does it exist?")
	}
	_, err := s.client.Request.Patch(url, &customAttribute, &responseAttr)
	if err != nil {
		return nil, err
	}
	return &responseAttr, nil
}

// Delete -> https://docs.taiga.io/api.html#epic-custom-attributes-delete
func (s *EpicCustomAttributeService) Delete(customAttributeID int) (*http.Response, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	return s.client.Request.Delete(url)
}

// Update is an alias for Edit.
func (s *EpicCustomAttributeService) Update(customAttribute *EpicCustomAttribute) (*EpicCustomAttribute, error) {
	return s.Edit(customAttribute)
}
