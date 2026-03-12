package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// IssueCustomAttributeService is a handle to actions related to issue custom attributes.
type IssueCustomAttributeService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#issue-custom-attributes-list
func (s *IssueCustomAttributeService) List(queryParams *ProjectIDQueryParams) ([]IssueCustomAttribute, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var attrs []IssueCustomAttribute
	_, err = s.client.Request.Get(url, &attrs)
	if err != nil {
		return nil, err
	}
	return attrs, nil
}

// Get -> https://docs.taiga.io/api.html#issue-custom-attributes-get
func (s *IssueCustomAttributeService) Get(customAttributeID int) (*IssueCustomAttribute, error) {
	if err := requirePositiveID("customAttributeID", customAttributeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	var attr IssueCustomAttribute
	_, err := s.client.Request.Get(url, &attr)
	if err != nil {
		return nil, err
	}
	return &attr, nil
}

// Create -> https://docs.taiga.io/api.html#issue-custom-attributes-create
func (s *IssueCustomAttributeService) Create(customAttribute *IssueCustomAttribute) (*IssueCustomAttribute, error) {
	if err := requireNonNil("customAttribute", customAttribute); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responseAttr IssueCustomAttribute
	projectID, err := resolveProjectID(customAttribute.Project, s.defaultProjectID, "project")
	if err != nil {
		return nil, err
	}
	if isEmpty(customAttribute.Name) || isEmpty(customAttribute.Type) {
		return nil, errors.New("a mandatory field(project, name, type) is missing. See API documentation")
	}
	payload := *customAttribute
	payload.Project = projectID
	_, err = s.client.Request.Post(url, &payload, &responseAttr)
	if err != nil {
		return nil, err
	}
	return &responseAttr, nil
}

// Edit -> https://docs.taiga.io/api.html#issue-custom-attributes-edit
func (s *IssueCustomAttributeService) Edit(customAttribute *IssueCustomAttribute) (*IssueCustomAttribute, error) {
	if err := requireNonNil("customAttribute", customAttribute); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttribute.ID))
	var responseAttr IssueCustomAttribute
	if err := requirePositiveID("customAttributeID", customAttribute.ID); err != nil {
		return nil, err
	}
	payload, err := sparsePatchMapFromStruct(customAttribute, "id", "created_date", "modified_date")
	if err != nil {
		return nil, err
	}
	_, err = s.client.Request.Patch(url, &payload, &responseAttr)
	if err != nil {
		return nil, err
	}
	return &responseAttr, nil
}

// Delete -> https://docs.taiga.io/api.html#issue-custom-attributes-delete
func (s *IssueCustomAttributeService) Delete(customAttributeID int) (*http.Response, error) {
	if err := requirePositiveID("customAttributeID", customAttributeID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(customAttributeID))
	return s.client.Request.Delete(url)
}

// Update is an alias for Edit.
func (s *IssueCustomAttributeService) Update(customAttribute *IssueCustomAttribute) (*IssueCustomAttribute, error) {
	return s.Edit(customAttribute)
}
