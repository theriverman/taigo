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
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	var attr UserStoryCustomAttribute
	_, err := s.client.Request.Get(url, &attr)
	if err != nil {
		return nil, err
	}
	return &attr, nil
}

// Create -> https://docs.taiga.io/api.html#user-story-custom-attributes-create
func (s *UserStoryCustomAttributeService) Create(customAttribute *UserStoryCustomAttribute) (*UserStoryCustomAttribute, error) {
	if err := requireNonNil("customAttribute", customAttribute); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responseAttr UserStoryCustomAttribute
	if isEmpty(customAttribute.Project) || isEmpty(customAttribute.Name) || isEmpty(customAttribute.Type) {
		return nil, errors.New("a mandatory field(project, name, type) is missing. See API documentation")
	}
	_, err := s.client.Request.Post(url, customAttribute, &responseAttr)
	if err != nil {
		return nil, err
	}
	return &responseAttr, nil
}

// Edit -> https://docs.taiga.io/api.html#user-story-custom-attributes-edit
func (s *UserStoryCustomAttributeService) Edit(customAttribute *UserStoryCustomAttribute) (*UserStoryCustomAttribute, error) {
	if err := requireNonNil("customAttribute", customAttribute); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttribute.ID))
	var responseAttr UserStoryCustomAttribute
	if customAttribute.ID == 0 {
		return nil, errors.New("passed UserStoryCustomAttribute does not have an ID yet. Does it exist?")
	}
	_, err := s.client.Request.Patch(url, customAttribute, &responseAttr)
	if err != nil {
		return nil, err
	}
	return &responseAttr, nil
}

// Delete -> https://docs.taiga.io/api.html#user-story-custom-attributes-delete
func (s *UserStoryCustomAttributeService) Delete(customAttributeID int) (*http.Response, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	return s.client.Request.Delete(url)
}

// Update is an alias for Edit.
func (s *UserStoryCustomAttributeService) Update(customAttribute *UserStoryCustomAttribute) (*UserStoryCustomAttribute, error) {
	return s.Edit(customAttribute)
}
