package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// Severity -> https://docs.taiga.io/api.html#severities
type Severity struct {
	Color     string `json:"color,omitempty"`
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Order     int    `json:"order,omitempty"`
	ProjectID int    `json:"project_id,omitempty"`
}

// SeverityCreateRequest represents payload for creating severities.
type SeverityCreateRequest struct {
	Color   string `json:"color,omitempty"`
	Name    string `json:"name"`
	Order   int    `json:"order,omitempty"`
	Project int    `json:"project"`
}

// SeverityEditRequest represents sparse non-destructive updates for severities.
type SeverityEditRequest struct {
	Color   string `json:"color,omitempty"`
	Name    string `json:"name,omitempty"`
	Order   int    `json:"order,omitempty"`
	Project int    `json:"project,omitempty"`
}

// SeverityPatch represents explicit PATCH payload for severities.
type SeverityPatch struct {
	Color   *string `json:"color,omitempty"`
	Name    *string `json:"name,omitempty"`
	Order   *int    `json:"order,omitempty"`
	Project *int    `json:"project,omitempty"`
}

// SeverityService is a handle to actions related to severities.
type SeverityService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List -> https://docs.taiga.io/api.html#severities-list
func (s *SeverityService) List(queryParams *ProjectIDQueryParams) ([]Severity, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var severities []Severity
	_, err = s.client.Request.Get(url, &severities)
	if err != nil {
		return nil, err
	}
	return severities, nil
}

// Get -> https://docs.taiga.io/api.html#severities-get
func (s *SeverityService) Get(severityID int) (*Severity, error) {
	if err := requirePositiveID("severityID", severityID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(severityID))
	var severity Severity
	_, err := s.client.Request.Get(url, &severity)
	if err != nil {
		return nil, err
	}
	return &severity, nil
}

// Create -> https://docs.taiga.io/api.html#severities-create
func (s *SeverityService) Create(request *SeverityCreateRequest) (*Severity, error) {
	if err := requireNonNil("request", request); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var responseSeverity Severity
	projectID, err := resolveProjectID(request.Project, s.defaultProjectID, "project")
	if err != nil {
		return nil, err
	}
	if isEmpty(request.Name) {
		return nil, errors.New("a mandatory field(project, name) is missing. See API documentation")
	}
	payload := *request
	payload.Project = projectID
	_, err = s.client.Request.Post(url, &payload, &responseSeverity)
	if err != nil {
		return nil, err
	}
	return &responseSeverity, nil
}

// Edit -> https://docs.taiga.io/api.html#severities-edit
func (s *SeverityService) Edit(severityID int, request *SeverityEditRequest) (*Severity, error) {
	if err := requireNonNil("request", request); err != nil {
		return nil, err
	}
	if err := requirePositiveID("severityID", severityID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(severityID))
	var responseSeverity Severity
	payload, err := sparsePatchMapFromStruct(request)
	if err != nil {
		return nil, err
	}
	_, err = s.client.Request.Patch(url, &payload, &responseSeverity)
	if err != nil {
		return nil, err
	}
	return &responseSeverity, nil
}

// Update is an alias for Edit.
func (s *SeverityService) Update(severityID int, request *SeverityEditRequest) (*Severity, error) {
	return s.Edit(severityID, request)
}

// Patch sends an explicit PATCH payload to edit a severity.
func (s *SeverityService) Patch(severityID int, patch *SeverityPatch) (*Severity, error) {
	if err := requireNonNil("patch", patch); err != nil {
		return nil, err
	}
	if err := requirePositiveID("severityID", severityID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(severityID))
	var responseSeverity Severity
	_, err := s.client.Request.Patch(url, patch, &responseSeverity)
	if err != nil {
		return nil, err
	}
	return &responseSeverity, nil
}

// Delete -> https://docs.taiga.io/api.html#severities-delete
func (s *SeverityService) Delete(severityID int) (*http.Response, error) {
	if err := requirePositiveID("severityID", severityID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(severityID))
	return s.client.Request.Delete(url)
}
