package taigo

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/go-querystring/query"
)

// EpicService is a handle to actions related to Epics
//
// https://taigaio.github.io/taiga-doc/dist/api.html#epics
type EpicService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List => https://taigaio.github.io/taiga-doc/dist/api.html#epics-list
//
// Available Meta: *EpicDetailLIST
func (s *EpicService) List(queryParams *EpicsQueryParams) ([]Epic, error) {
	url := s.client.MakeURL(s.Endpoint)
	switch {
	case queryParams != nil:
		paramValues, _ := query.Values(queryParams)
		url = fmt.Sprintf("%s?%s", url, paramValues.Encode())
	case s.defaultProjectID != 0:
		url = url + projectIDQueryParam(s.defaultProjectID)
	}
	var epics EpicDetailLIST
	_, err := s.client.Request.Get(url, &epics)
	if err != nil {
		return nil, err
	}

	return epics.AsEpics()
}

// Create => https://taigaio.github.io/taiga-doc/dist/api.html#epics-create
//
// Available Meta: *EpicDetail
func (s *EpicService) Create(epic *Epic) (*Epic, error) {
	url := s.client.MakeURL(s.Endpoint)
	var e EpicDetail

	// Check for required fields
	// project, subject
	if isEmpty(epic.Project) || isEmpty(epic.Subject) {
		return nil, errors.New("a mandatory field(Project, Subject) is missing. See API documentataion")
	}

	_, err := s.client.Request.Post(url, &epic, &e)
	if err != nil {
		return nil, err
	}

	return e.AsEpic()
}

// Get => https://taigaio.github.io/taiga-doc/dist/api.html#epics-get
//
// Available Meta: *EpicDetailGET
func (s *EpicService) Get(epicID int) (*Epic, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(epicID))
	var e EpicDetailGET
	_, err := s.client.Request.Get(url, &e)
	if err != nil {
		return nil, err
	}
	return e.AsEpic()
}

// GetByRef => https://taigaio.github.io/taiga-doc/dist/api.html#epics-get-by-ref
//
// The passed epicRef should be an int taken from the Epic's URL
// The passed *Project struct should have at least one of the following fields set:
//		ID 	 (int)
//		Slug (string)
// If none of the above fields are set, an error is returned.
// If both fields are set, *Project.ID will be preferred.
//
// Available Meta: *EpicDetailGET
func (s *EpicService) GetByRef(epicRef int, project *Project) (*Epic, error) {
	var e EpicDetailGET
	var url string

	switch {
	case project.ID > 0:
		url = s.client.MakeURL(fmt.Sprintf("%s/by_ref?ref=%d&project=%d", s.Endpoint, (epicRef), project.ID))
	case len(project.Slug) > 0:
		url = s.client.MakeURL(fmt.Sprintf("%s/by_ref?ref=%d&project__slug=%s", s.Endpoint, epicRef, project.Slug))
	default:
		return nil, errors.New("no ID or Ref defined in passed project struct")
	}

	_, err := s.client.Request.Get(url, &e)
	if err != nil {
		return nil, err
	}
	return e.AsEpic()
}

// Edit edits an Epic via a PATCH request => https://taigaio.github.io/taiga-doc/dist/api.html#epics-edit
// Available Meta: EpicDetail
func (s *EpicService) Edit(epic *Epic) (*Epic, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(epic.ID))
	var e EpicDetail

	if epic.ID == 0 {
		return nil, errors.New("passed Epic does not have an ID yet. Does it exist?")
	}

	// Taiga OCC
	remoteEpic, err := s.Get(epic.ID)
	if err != nil {
		return nil, err
	}
	epic.Version = remoteEpic.Version
	_, err = s.client.Request.Patch(url, &epic, &e)
	if err != nil {
		return nil, err
	}
	return e.AsEpic()
}

// Delete => https://taigaio.github.io/taiga-doc/dist/api.html#epics-delete
func (s *EpicService) Delete(epicID int) (*http.Response, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(epicID))
	return s.client.Request.Delete(url)
}

// BulkCreation => https://taigaio.github.io/taiga-doc/dist/api.html#epics-bulk-create
/*
This is not yet implemented, placeholder only.
It seems to be pointless to implement this operation here. A for loop around `Create` is much more efficient and accurate.
*/
// func (s *EpicService) BulkCreation() {}

// EpicFiltersData => https://taigaio.github.io/taiga-doc/dist/api.html#epics-get-filters-data
/*
This is not yet implemented, placeholder only.
It seems to be pointless to implement this operation here. A for loop around `Create` is much more efficient and accurate.
*/
// func (s *EpicService) EpicFiltersData() {}

// ListRelatedUserStories => https://taigaio.github.io/taiga-doc/dist/api.html#epics-related-user-stories-list
func (s *EpicService) ListRelatedUserStories(epicID int) ([]EpicRelatedUserStoryDetail, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(epicID), "related_userstories")
	var e []EpicRelatedUserStoryDetail
	_, err := s.client.Request.Get(url, &e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// CreateRelatedUserStory => https://taigaio.github.io/taiga-doc/dist/api.html#epics-related-user-stories-create
//
// Mandatory parameters: `EpicID`; `UserStoryID`
// Accepted UserStory values: `UserStory.ID`
func (s *EpicService) CreateRelatedUserStory(EpicID int, UserStoryID int) (*EpicRelatedUserStoryDetail, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(EpicID), "related_userstories")
	e := EpicRelatedUserStoryDetail{EpicID: EpicID, UserStoryID: UserStoryID}
	_, err := s.client.Request.Post(url, &e, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// CreateAttachment creates a new Epic attachment => https://taigaio.github.io/taiga-doc/dist/api.html#epics-create-attachment
//
func (s *EpicService) CreateAttachment(attachment *Attachment, epic *Epic) (*Attachment, error) {
	url := s.client.MakeURL(s.Endpoint, "attachments")
	return newfileUploadRequest(s.client, url, attachment, epic)
}

// Clone takes an *Epic struct with loaded properties and duplicates it
//
// Available Meta: *EpicDetail
func (e *Epic) Clone(s *EpicService) (*Epic, error) {
	// Clean up data
	e.ID = 0
	e.Version = 0
	e.Ref = 0
	return s.Create(e)
}
