package taigo

import "net/http"

// ProjectTemplate is a raw DTO for /project-templates endpoints.
type ProjectTemplate = RawResource

// ProjectTemplateService is a handle to actions related to project templates.
type ProjectTemplateService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#project-templates-list
func (s *ProjectTemplateService) List() ([]ProjectTemplate, error) {
	return listRawResources(s.client, s.Endpoint, 0, nil)
}

// Get -> https://docs.taiga.io/api.html#project-templates-get
func (s *ProjectTemplateService) Get(projectTemplateID int) (*ProjectTemplate, error) {
	return getRawResource(s.client, s.Endpoint, projectTemplateID)
}

// Create -> https://docs.taiga.io/api.html#project-templates-create
func (s *ProjectTemplateService) Create(payload any) (*ProjectTemplate, error) {
	return createRawResource(s.client, s.Endpoint, payload)
}

// Edit -> https://docs.taiga.io/api.html#project-templates-edit
func (s *ProjectTemplateService) Edit(projectTemplateID int, payload any) (*ProjectTemplate, error) {
	return patchRawResource(s.client, s.Endpoint, projectTemplateID, payload)
}

// Update is an alias for Edit.
func (s *ProjectTemplateService) Update(projectTemplateID int, payload any) (*ProjectTemplate, error) {
	return s.Edit(projectTemplateID, payload)
}

// Delete -> https://docs.taiga.io/api.html#project-templates-delete
func (s *ProjectTemplateService) Delete(projectTemplateID int) (*http.Response, error) {
	return deleteRawResource(s.client, s.Endpoint, projectTemplateID)
}
