package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// EpicService is a handle to actions related to Epics
//
// https://taigaio.github.io/taiga-doc/dist/api.html#epics
type EpicService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

type epicCreatePayload struct {
	AssignedTo        int      `json:"assigned_to,omitempty"`
	BlockedNote       string   `json:"blocked_note,omitempty"`
	ClientRequirement bool     `json:"client_requirement,omitempty"`
	Color             string   `json:"color,omitempty"`
	Description       string   `json:"description,omitempty"`
	EpicsOrder        int64    `json:"epics_order,omitempty"`
	IsBlocked         bool     `json:"is_blocked,omitempty"`
	Project           int      `json:"project"`
	Status            int      `json:"status,omitempty"`
	Subject           string   `json:"subject"`
	Tags              []string `json:"tags,omitempty"`
	TeamRequirement   bool     `json:"team_requirement,omitempty"`
	Watchers          []int    `json:"watchers,omitempty"`
}

// EpicPatch represents an explicit PATCH payload for epics.
// Pointer fields allow intentionally setting zero-values (false, 0, "").
type EpicPatch struct {
	AssignedTo        *int      `json:"assigned_to,omitempty"`
	BlockedNote       *string   `json:"blocked_note,omitempty"`
	ClientRequirement *bool     `json:"client_requirement,omitempty"`
	Color             *string   `json:"color,omitempty"`
	Description       *string   `json:"description,omitempty"`
	EpicsOrder        *int64    `json:"epics_order,omitempty"`
	IsBlocked         *bool     `json:"is_blocked,omitempty"`
	Project           *int      `json:"project,omitempty"`
	Status            *int      `json:"status,omitempty"`
	Subject           *string   `json:"subject,omitempty"`
	Tags              *[]string `json:"tags,omitempty"`
	TeamRequirement   *bool     `json:"team_requirement,omitempty"`
	Version           int       `json:"version"`
	Watchers          *[]int    `json:"watchers,omitempty"`
}

// EpicBulkCreatePayload is used by the bulk-create endpoint.
type EpicBulkCreatePayload struct {
	Project   int    `json:"project"`
	BulkEpics string `json:"bulk_epics"`
}

// List => https://taigaio.github.io/taiga-doc/dist/api.html#epics-list
//
// Available Meta: *EpicDetailLIST
func (s *EpicService) List(queryParams *EpicsQueryParams) ([]Epic, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var epics EpicDetailLIST
	_, err = s.client.Request.Get(url, &epics)
	if err != nil {
		return nil, err
	}

	return epics.AsEpics()
}

// Create => https://taigaio.github.io/taiga-doc/dist/api.html#epics-create
//
// Available Meta: *EpicDetail
func (s *EpicService) Create(epic *Epic) (*Epic, error) {
	if err := requireNonNil("epic", epic); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var e EpicDetail

	// Check for required fields
	// project, subject
	if isEmpty(epic.Project) || isEmpty(epic.Subject) {
		return nil, errors.New("a mandatory field(Project, Subject) is missing. See API documentataion")
	}

	payload := epicCreatePayload{
		AssignedTo:        epic.AssignedTo,
		BlockedNote:       epic.BlockedNote,
		ClientRequirement: epic.ClientRequirement,
		Color:             epic.Color,
		Description:       epic.Description,
		EpicsOrder:        epic.EpicsOrder,
		IsBlocked:         epic.IsBlocked,
		Project:           epic.Project,
		Status:            epic.Status,
		Subject:           epic.Subject,
		Tags:              tagsToNames(epic.Tags),
		TeamRequirement:   epic.TeamRequirement,
		Watchers:          epic.Watchers,
	}

	_, err := s.client.Request.Post(url, &payload, &e)
	if err != nil {
		return nil, err
	}

	return e.AsEpic()
}

// Get => https://taigaio.github.io/taiga-doc/dist/api.html#epics-get
//
// Available Meta: *EpicDetailGET
func (s *EpicService) Get(epicID int) (*Epic, error) {
	if err := requirePositiveID("epicID", epicID); err != nil {
		return nil, err
	}
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
//
//	ID 	 (int)
//	Slug (string)
//
// If none of the above fields are set, an error is returned.
// If both fields are set, *Project.ID will be preferred.
//
// Available Meta: *EpicDetailGET
func (s *EpicService) GetByRef(epicRef int, project *Project) (*Epic, error) {
	if err := requirePositiveID("epicRef", epicRef); err != nil {
		return nil, err
	}
	var e EpicDetailGET
	var url string
	if project == nil {
		return nil, errors.New("project must not be nil")
	}

	type byRefQueryParams struct {
		Ref         int    `url:"ref"`
		Project     int    `url:"project,omitempty"`
		ProjectSlug string `url:"project__slug,omitempty"`
	}
	queryParams := &byRefQueryParams{Ref: epicRef}
	switch {
	case project.ID > 0:
		queryParams.Project = project.ID
	case len(project.Slug) > 0:
		queryParams.ProjectSlug = project.Slug
	default:
		return nil, errors.New("no ID or Ref defined in passed project struct")
	}
	url, err := appendQueryParams(s.client.MakeURL(s.Endpoint, "by_ref"), queryParams)
	if err != nil {
		return nil, err
	}

	_, err = s.client.Request.Get(url, &e)
	if err != nil {
		return nil, err
	}
	return e.AsEpic()
}

// Edit edits an Epic via a PATCH request => https://taigaio.github.io/taiga-doc/dist/api.html#epics-edit
// Available Meta: EpicDetail
func (s *EpicService) Edit(epic *Epic) (*Epic, error) {
	if err := requireNonNil("epic", epic); err != nil {
		return nil, err
	}

	if epic.ID == 0 {
		return nil, errors.New("passed Epic does not have an ID yet. Does it exist?")
	}
	if epic.Version == 0 {
		return nil, errors.New("version is required for epic edit")
	}

	patch := &EpicPatch{
		AssignedTo:        ptr(epic.AssignedTo),
		BlockedNote:       ptr(epic.BlockedNote),
		ClientRequirement: ptr(epic.ClientRequirement),
		Color:             ptr(epic.Color),
		Description:       ptr(epic.Description),
		EpicsOrder:        ptr(epic.EpicsOrder),
		IsBlocked:         ptr(epic.IsBlocked),
		Project:           ptr(epic.Project),
		Status:            ptr(epic.Status),
		Subject:           ptr(epic.Subject),
		TeamRequirement:   ptr(epic.TeamRequirement),
		Version:           epic.Version,
	}
	if epic.Tags != nil {
		tags := tagsToNames(epic.Tags)
		if tags == nil {
			tags = []string{}
		}
		patch.Tags = &tags
	}
	if epic.Watchers != nil {
		watchers := append([]int(nil), epic.Watchers...)
		patch.Watchers = &watchers
	}
	return s.Patch(epic.ID, patch)
}

// Patch sends an explicit PATCH payload to edit an epic.
func (s *EpicService) Patch(epicID int, patch *EpicPatch) (*Epic, error) {
	if err := requireNonNil("patch", patch); err != nil {
		return nil, err
	}
	if err := requirePositiveID("epicID", epicID); err != nil {
		return nil, err
	}
	if patch.Version == 0 {
		return nil, errors.New("version is required for epic patch")
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(epicID))
	var e EpicDetail
	_, err := s.client.Request.Patch(url, patch, &e)
	if err != nil {
		return nil, err
	}
	return e.AsEpic()
}

// Update is an alias for Edit.
func (s *EpicService) Update(epic *Epic) (*Epic, error) {
	return s.Edit(epic)
}

// Delete => https://taigaio.github.io/taiga-doc/dist/api.html#epics-delete
func (s *EpicService) Delete(epicID int) (*http.Response, error) {
	if err := requirePositiveID("epicID", epicID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(epicID))
	return s.client.Request.Delete(url)
}

// BulkCreate => https://docs.taiga.io/api.html#epics-bulk-create
func (s *EpicService) BulkCreate(payload *EpicBulkCreatePayload) ([]RawResource, error) {
	if err := requireNonNil("payload", payload); err != nil {
		return nil, err
	}
	copyPayload := *payload
	if copyPayload.Project == 0 {
		copyPayload.Project = s.defaultProjectID
	}
	if err := requirePositiveID("project", copyPayload.Project); err != nil {
		return nil, err
	}
	if copyPayload.BulkEpics == "" {
		return nil, errors.New("bulk_epics is required")
	}
	return postRawResourceListAtPath(s.client, &copyPayload, s.Endpoint, "bulk_create")
}

// GetFiltersData => https://docs.taiga.io/api.html#epics-get-filters-data
func (s *EpicService) GetFiltersData(projectID int) (*EpicFiltersDataDetail, error) {
	if projectID == 0 {
		projectID = s.defaultProjectID
	}
	if err := requirePositiveID("projectID", projectID); err != nil {
		return nil, err
	}
	queryParams := struct {
		Project int `url:"project"`
	}{Project: projectID}
	url, err := appendQueryParams(s.client.MakeURL(s.Endpoint, "filters_data"), &queryParams)
	if err != nil {
		return nil, err
	}
	var filtersData EpicFiltersDataDetail
	_, err = s.client.Request.Get(url, &filtersData)
	if err != nil {
		return nil, err
	}
	return &filtersData, nil
}

// ListRelatedUserStories => https://taigaio.github.io/taiga-doc/dist/api.html#epics-related-user-stories-list
func (s *EpicService) ListRelatedUserStories(epicID int) ([]EpicRelatedUserStoryDetail, error) {
	if err := requirePositiveID("epicID", epicID); err != nil {
		return nil, err
	}
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
	if err := requirePositiveID("epicID", EpicID); err != nil {
		return nil, err
	}
	if err := requirePositiveID("userStoryID", UserStoryID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(EpicID), "related_userstories")
	e := EpicRelatedUserStoryDetail{EpicID: EpicID, UserStoryID: UserStoryID}
	_, err := s.client.Request.Post(url, &e, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// CreateAttachment creates a new Epic attachment => https://taigaio.github.io/taiga-doc/dist/api.html#epics-create-attachment
func (s *EpicService) CreateAttachment(attachment *Attachment, epic *Epic) (*Attachment, error) {
	url := s.client.MakeURL(s.Endpoint, "attachments")
	return newfileUploadRequest(s.client, url, attachment, epic)
}

// Clone takes an *Epic struct with loaded properties and duplicates it
//
// Available Meta: *EpicDetail
func (e *Epic) Clone(s *EpicService) (*Epic, error) {
	if err := requireNonNil("epic", e); err != nil {
		return nil, err
	}
	// Clean up data
	clone := *e
	clone.ID = 0
	clone.Version = 0
	clone.Ref = 0
	return s.Create(&clone)
}
