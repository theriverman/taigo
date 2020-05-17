package gotaiga

import (
	"errors"
	"fmt"

	"github.com/google/go-querystring/query"
)

const endpointEpics = "/epics"

// EpicService is a handle to actions related to Epics
//
// https://taigaio.github.io/taiga-doc/dist/api.html#epics
type EpicService struct {
	client *Client
}

// List => https://taigaio.github.io/taiga-doc/dist/api.html#epics-list
//
// Available Meta: *EpicDetailLIST
func (s *EpicService) List(queryParams *EpicsQueryParams) ([]Epic, error) {
	url := s.client.APIURL + endpointEpics
	if queryParams != nil {
		paramValues, _ := query.Values(queryParams)
		url = fmt.Sprintf("%s?%s", url, paramValues.Encode())
	} else if s.client.HasDefaultProject() {
		url = url + s.client.GetDefaultProjectAsQueryParam()
	}
	var epics EpicDetailLIST
	err := getRequest(s.client, &epics, url)
	if err != nil {
		return nil, err
	}

	return epics.AsEpics()
}

// Create => https://taigaio.github.io/taiga-doc/dist/api.html#epics-create
//
// Available Meta: *EpicDetail
func (s *EpicService) Create(epic Epic) (*Epic, error) {
	url := s.client.APIURL + endpointEpics
	var respEpic EpicDetail

	// Check for required fields
	// project, subject
	if isEmpty(epic.Project) || isEmpty(epic.Subject) {
		return nil, errors.New("A mandatory field is missing. See API documentataion")
	}

	err := postRequest(s.client, &respEpic, url, epic)
	if err != nil {
		return nil, err
	}

	return respEpic.AsEpic()
}

// Get => https://taigaio.github.io/taiga-doc/dist/api.html#epics-get
//
// Available Meta: *EpicDetailGET
func (s *EpicService) Get(epicID int) (*Epic, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointEpics, epicID)
	var e EpicDetailGET
	err := getRequest(s.client, &e, url)
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
	case project.ID != 0:
		url = s.client.APIURL + fmt.Sprintf("%s/by_ref?ref=%d&project=%d", endpointEpics, epicRef, project.ID)
		break
	case len(project.Slug) > 0:
		url = s.client.APIURL + fmt.Sprintf("%s/by_ref?ref=%d&project__slug=%s", endpointEpics, epicRef, project.Slug)
		break
	default:
		return nil, errors.New("No ID or Ref defined in passed project struct")
	}

	err := getRequest(s.client, &e, url)
	if err != nil {
		return nil, err
	}
	return e.AsEpic()
}

// Edit edits an Epic via a PATCH request => https://taigaio.github.io/taiga-doc/dist/api.html#epics-edit
// Available Meta: EpicDetail
func (s *EpicService) Edit(epic Epic) (*Epic, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointEpics, epic.ID)
	var epicDetail EpicDetail

	if epic.ID == 0 {
		return nil, errors.New("Passed Epic does not have an ID yet. Does it exist?")
	}

	// Taiga OCC
	remoteEpic, err := s.Get(epic.ID)
	if err != nil {
		return nil, err
	}
	epic.Version = remoteEpic.Version
	err = patchRequest(s.client, &epicDetail, url, &epic)
	if err != nil {
		return nil, err
	}
	return epicDetail.AsEpic()
}

// Delete => https://taigaio.github.io/taiga-doc/dist/api.html#epics-delete
func (s *EpicService) Delete(epicID int) error {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointEpics, epicID)
	return deleteRequest(s.client, url)
}

// BulkCreation => https://taigaio.github.io/taiga-doc/dist/api.html#epics-bulk-create
/*
This is not yet implemented, placeholder only.
It seems to be pointless to implement this operation here. A for loop around `CreateEpic` is much more efficient and accurate.
*/
// func (s *EpicService) BulkCreation() {}

// EpicFiltersData => https://taigaio.github.io/taiga-doc/dist/api.html#epics-get-filters-data
/*
This is not yet implemented, placeholder only.
It seems to be pointless to implement this operation here. A for loop around `CreateEpic` is much more efficient and accurate.
*/
// func (s *EpicService) EpicFiltersData() {}

// ListRelatedUserStories => https://taigaio.github.io/taiga-doc/dist/api.html#epics-related-user-stories-list
func (s *EpicService) ListRelatedUserStories(epicID int) ([]EpicRelatedUserStoryDetail, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d/related_userstories", endpointEpics, epicID)
	var resp []EpicRelatedUserStoryDetail
	err := getRequest(s.client, &resp, url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateRelatedUserStory => https://taigaio.github.io/taiga-doc/dist/api.html#epics-related-user-stories-create
//
// Mandatory parameters: `EpicID`; `UserStoryID`
// Accepted UserStory values: `UserStory.ID`
func (s *EpicService) CreateRelatedUserStory(epicRelatedUSDetail *EpicRelatedUserStoryDetail) (*EpicRelatedUserStoryDetail, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d/related_userstories", endpointEpics, epicRelatedUSDetail.EpicID)
	var resp EpicRelatedUserStoryDetail

	err := postRequest(s.client, &resp, url, epicRelatedUSDetail)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateAttachment creates a new Epic attachment => https://taigaio.github.io/taiga-doc/dist/api.html#epics-create-attachment
//
// Mandatory parameters: `ObjectID`; `ProjectID`;
// **Note**: Set ObjectID to epic.ID
func (s *EpicService) CreateAttachment(attachment *Attachment, filePath string) (*Attachment, error) {
	url := s.client.APIURL + endpointEpics + "/attachments"
	attachment.filePath = filePath
	return newfileUploadRequest(s.client, url, attachment)
}

// Clone takes an *Epic struct with loaded properties and duplicates it
//
// Available Meta: *EpicDetail
func (e *Epic) Clone(s *EpicService) (*Epic, error) {
	// Clean up data
	e.ID = 0
	e.Version = 0
	e.Ref = 0
	return s.Create(*e)
}
