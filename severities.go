package taigo

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/go-querystring/query"
)

// Severity -> https://docs.taiga.io/api.html#severities
type Severity struct {
	Color   string `json:"color,omitempty"`
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Order   int    `json:"order,omitempty"`
	Project int    `json:"project,omitempty"`
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
	switch {
	case queryParams != nil:
		paramValues, _ := query.Values(queryParams)
		url = fmt.Sprintf("%s?%s", url, paramValues.Encode())
	case s.defaultProjectID != 0:
		url = url + projectIDQueryParam(s.defaultProjectID)
	}
	var severities []Severity
	_, err := s.client.Request.Get(url, &severities)
	if err != nil {
		return nil, err
	}
	return severities, nil
}

// Get -> https://docs.taiga.io/api.html#severities-get
func (s *SeverityService) Get(severityID int) (*Severity, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(severityID))
	var severity Severity
	_, err := s.client.Request.Get(url, &severity)
	if err != nil {
		return nil, err
	}
	return &severity, nil
}

// Create -> https://docs.taiga.io/api.html#severities-create
func (s *SeverityService) Create(severity *Severity) (*Severity, error) {
	url := s.client.MakeURL(s.Endpoint)
	var responseSeverity Severity
	if isEmpty(severity.Project) || isEmpty(severity.Name) {
		return nil, errors.New("a mandatory field(project, name) is missing. See API documentation")
	}
	_, err := s.client.Request.Post(url, &severity, &responseSeverity)
	if err != nil {
		return nil, err
	}
	return &responseSeverity, nil
}

// Edit -> https://docs.taiga.io/api.html#severities-edit
func (s *SeverityService) Edit(severity *Severity) (*Severity, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(severity.ID))
	var responseSeverity Severity
	if severity.ID == 0 {
		return nil, errors.New("passed Severity does not have an ID yet. Does it exist?")
	}
	_, err := s.client.Request.Patch(url, &severity, &responseSeverity)
	if err != nil {
		return nil, err
	}
	return &responseSeverity, nil
}

// Delete -> https://docs.taiga.io/api.html#severities-delete
func (s *SeverityService) Delete(severityID int) (*http.Response, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(severityID))
	return s.client.Request.Delete(url)
}
