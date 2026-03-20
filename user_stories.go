package taigo

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// UserStoryService is a handle to actions related to UserStories
//
// https://taigaio.github.io/taiga-doc/dist/api.html#user-stories
type UserStoryService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

type userStoryCreatePayload struct {
	AssignedTo        int         `json:"assigned_to,omitempty"`
	BacklogOrder      int64       `json:"backlog_order,omitempty"`
	BlockedNote       string      `json:"blocked_note,omitempty"`
	ClientRequirement bool        `json:"client_requirement,omitempty"`
	Description       string      `json:"description,omitempty"`
	ExternalReference []string    `json:"external_reference,omitempty"`
	IsBlocked         bool        `json:"is_blocked,omitempty"`
	KanbanOrder       int64       `json:"kanban_order,omitempty"`
	Milestone         int         `json:"milestone,omitempty"`
	Points            AgilePoints `json:"points,omitempty"`
	Project           int         `json:"project"`
	SprintOrder       int         `json:"sprint_order,omitempty"`
	Status            int         `json:"status,omitempty"`
	Subject           string      `json:"subject"`
	Tags              []string    `json:"tags,omitempty"`
	TeamRequirement   bool        `json:"team_requirement,omitempty"`
	Watchers          []int       `json:"watchers,omitempty"`
}

// UserStoryPatch represents an explicit PATCH payload for user stories.
// Pointer fields allow intentionally setting zero-values (false, 0, "").
type UserStoryPatch struct {
	AssignedTo        *int         `json:"assigned_to,omitempty"`
	BacklogOrder      *int64       `json:"backlog_order,omitempty"`
	BlockedNote       *string      `json:"blocked_note,omitempty"`
	ClientRequirement *bool        `json:"client_requirement,omitempty"`
	Description       *string      `json:"description,omitempty"`
	ExternalReference *[]string    `json:"external_reference,omitempty"`
	IsBlocked         *bool        `json:"is_blocked,omitempty"`
	KanbanOrder       *int64       `json:"kanban_order,omitempty"`
	Milestone         *int         `json:"milestone,omitempty"`
	Points            *AgilePoints `json:"points,omitempty"`
	Project           *int         `json:"project,omitempty"`
	SprintOrder       *int         `json:"sprint_order,omitempty"`
	Status            *int         `json:"status,omitempty"`
	Subject           *string      `json:"subject,omitempty"`
	Tags              *[]string    `json:"tags,omitempty"`
	TeamRequirement   *bool        `json:"team_requirement,omitempty"`
	Version           int          `json:"version"`
	Watchers          *[]int       `json:"watchers,omitempty"`
}

// List returns all User Stories | https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-list
// Available Meta: *[]UserStoryDetailLIST
func (s *UserStoryService) List(queryParams *UserStoryQueryParams) ([]UserStory, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var userstories UserStoryDetailLIST
	_, err = s.client.Request.Get(url, &userstories)
	if err != nil {
		return nil, err
	}
	return userstories.AsUserStory()
}

// Create creates a new User Story | https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-create
//
// Available Meta: *UserStoryDetail
func (s *UserStoryService) Create(userStory *UserStory) (*UserStory, error) {
	if err := requireNonNil("userStory", userStory); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var us UserStoryDetail
	projectID, err := resolveProjectID(userStory.Project, s.defaultProjectID, "project")
	if err != nil {
		return nil, err
	}

	// Check for required fields
	// project, subject
	if isEmpty(userStory.Subject) {
		return nil, errors.New("a mandatory field is missing. See API documentataion")
	}

	payload := userStoryCreatePayload{
		AssignedTo:        userStory.AssignedTo,
		BacklogOrder:      userStory.BacklogOrder,
		BlockedNote:       userStory.BlockedNote,
		ClientRequirement: userStory.ClientRequirement,
		Description:       userStory.Description,
		ExternalReference: userStory.ExternalReference,
		IsBlocked:         userStory.IsBlocked,
		KanbanOrder:       userStory.KanbanOrder,
		Milestone:         userStory.Milestone,
		Points:            userStory.Points,
		Project:           projectID,
		SprintOrder:       userStory.SprintOrder,
		Status:            userStory.Status,
		Subject:           userStory.Subject,
		Tags:              tagsToNames(userStory.Tags),
		TeamRequirement:   userStory.TeamRequirement,
		Watchers:          userStory.Watchers,
	}

	_, err = s.client.Request.Post(url, &payload, &us)
	if err != nil {
		return nil, err
	}

	return us.AsUserStory()
}

// Get -> https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-get
//
// Available Meta: *UserStoryDetailGET
func (s *UserStoryService) Get(userStoryID int) (*UserStory, error) {
	if err := requirePositiveID("userStoryID", userStoryID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(userStoryID))
	var us UserStoryDetailGET
	_, err := s.client.Request.Get(url, &us)
	if err != nil {
		return nil, err
	}
	return us.AsUserStory()
}

// GetByRef returns a User Story by Ref -> https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-get-by-ref
//
// The passed userStoryRef should be an int taken from the UserStory's URL
// The passed *Project struct should have at least one of the following fields set:
//
//	ID 	 (int)
//	Slug (string)
//
// If none of the above fields are set, an error is returned.
// If both fields are set, *Project.ID will be preferred.
//
// Available Meta: UserStoryDetailGET
func (s *UserStoryService) GetByRef(userStoryRef int, project *Project) (*UserStory, error) {
	if err := requirePositiveID("userStoryRef", userStoryRef); err != nil {
		return nil, err
	}
	var us UserStoryDetailGET
	var url string

	type byRefQueryParams struct {
		Ref         int    `url:"ref"`
		Project     int    `url:"project,omitempty"`
		ProjectSlug string `url:"project__slug,omitempty"`
	}
	queryParams := &byRefQueryParams{Ref: userStoryRef}
	switch {
	case project != nil && project.ID != 0:
		queryParams.Project = project.ID
	case project != nil && len(project.Slug) > 0:
		queryParams.ProjectSlug = project.Slug
	case s.defaultProjectID > 0:
		queryParams.Project = s.defaultProjectID
	default:
		return nil, errors.New("no project ID/slug provided and no mapped default project ID set")
	}
	url, err := appendQueryParams(s.client.MakeURL(s.Endpoint, "by_ref"), queryParams)
	if err != nil {
		return nil, err
	}

	_, err = s.client.Request.Get(url, &us)
	if err != nil {
		return nil, err
	}
	return us.AsUserStory()
}

// Edit sends a PATCH request to edit a User Story -> https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-edit
// Available Meta: UserStoryDetail
func (s *UserStoryService) Edit(us *UserStory) (*UserStory, error) {
	if err := requireNonNil("userStory", us); err != nil {
		return nil, err
	}

	if us.ID == 0 {
		return nil, errors.New("passed UserStory does not have an ID yet. Does it exist?")
	}
	if us.Version == 0 {
		return nil, errors.New("version is required for user story edit")
	}

	patchPayload := map[string]any{
		"version": us.Version,
	}
	if us.AssignedTo != 0 {
		patchPayload["assigned_to"] = us.AssignedTo
	}
	if us.BacklogOrder != 0 {
		patchPayload["backlog_order"] = us.BacklogOrder
	}
	if us.BlockedNote != "" {
		patchPayload["blocked_note"] = us.BlockedNote
	}
	if us.ClientRequirement {
		patchPayload["client_requirement"] = us.ClientRequirement
	}
	if us.Description != "" {
		patchPayload["description"] = us.Description
	}
	if us.IsBlocked {
		patchPayload["is_blocked"] = us.IsBlocked
	}
	if us.KanbanOrder != 0 {
		patchPayload["kanban_order"] = us.KanbanOrder
	}
	if us.Milestone != 0 {
		patchPayload["milestone"] = us.Milestone
	}
	if len(us.Points) > 0 {
		patchPayload["points"] = us.Points
	}
	if us.Project != 0 {
		patchPayload["project"] = us.Project
	}
	if us.SprintOrder != 0 {
		patchPayload["sprint_order"] = us.SprintOrder
	}
	if us.Status != 0 {
		patchPayload["status"] = us.Status
	}
	if us.Subject != "" {
		patchPayload["subject"] = us.Subject
	}
	if us.TeamRequirement {
		patchPayload["team_requirement"] = us.TeamRequirement
	}
	if us.ExternalReference != nil {
		externalRef := append([]string(nil), us.ExternalReference...)
		patchPayload["external_reference"] = externalRef
	}
	if us.Tags != nil {
		tags := tagsToNames(us.Tags)
		if tags == nil {
			tags = []string{}
		}
		patchPayload["tags"] = tags
	}
	if us.Watchers != nil {
		watchers := append([]int(nil), us.Watchers...)
		patchPayload["watchers"] = watchers
	}
	if len(patchPayload) == 1 {
		return nil, errors.New("no updatable user story fields were provided; use Patch for explicit zero-value updates")
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(us.ID))
	var responseUS UserStoryDetail
	_, err := s.client.Request.Patch(url, &patchPayload, &responseUS)
	if err != nil {
		return nil, err
	}
	return responseUS.AsUserStory()
}

// Patch sends an explicit PATCH payload to edit a user story.
func (s *UserStoryService) Patch(userStoryID int, patch *UserStoryPatch) (*UserStory, error) {
	if err := requireNonNil("patch", patch); err != nil {
		return nil, err
	}
	if err := requirePositiveID("userStoryID", userStoryID); err != nil {
		return nil, err
	}
	if patch.Version == 0 {
		return nil, errors.New("version is required for user story patch")
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(userStoryID))
	var responseUS UserStoryDetail
	_, err := s.client.Request.Patch(url, patch, &responseUS)
	if err != nil {
		return nil, err
	}
	return responseUS.AsUserStory()
}

// Update is an alias for Edit.
func (s *UserStoryService) Update(us *UserStory) (*UserStory, error) {
	return s.Edit(us)
}

// Delete -> https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-delete
func (s *UserStoryService) Delete(usID int) (*http.Response, error) {
	if err := requirePositiveID("userStoryID", usID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(usID))
	return s.client.Request.Delete(url)
}

// CreateAttachment creates a new UserStory attachment => https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-create-attachment
func (s *UserStoryService) CreateAttachment(attachment *Attachment, userstory *UserStory) (*Attachment, error) {
	url := s.client.MakeURL(s.Endpoint, "attachments")
	return newfileUploadRequest(s.client, url, attachment, userstory)
}

/*
	Advanced Operations
*/

// RelateToEpic relates the UserStory to an Epic via an EpicID
//
// TaigaClient must be a pointer to taiga.Client
// EpicID must be an int to desired Epic
func (us *UserStory) RelateToEpic(client *Client, epicID int) (*EpicRelatedUserStoryDetail, error) {
	if err := requireNonNil("client", client); err != nil {
		return nil, err
	}
	if err := requireNonNil("userStory", us); err != nil {
		return nil, err
	}
	if us.ID == 0 {
		return nil, fmt.Errorf("UserStory must be created before relating it to an Epic. UserStory.ID was 0")
	}
	return client.Epic.CreateRelatedUserStory(epicID, us.ID)
}

// Clone clones an existing UserStory with most fields
//
// Available Meta: UserStoryDetail
func (s *UserStoryService) Clone(srcUS *UserStory) (*UserStory, error) {
	if err := requireNonNil("userStory", srcUS); err != nil {
		return nil, err
	}
	clone := *srcUS
	clone.ID = 0
	clone.Ref = 0
	clone.Version = 0
	return s.Create(&clone)
}

// ListRelatedTasks returns all Tasks related to this UserStory
func (us *UserStory) ListRelatedTasks(client *Client, userStoryID int) ([]Task, error) {
	if err := requireNonNil("client", client); err != nil {
		return nil, err
	}
	return client.Task.List(&TasksQueryParams{UserStory: userStoryID})
}

// CreateRelatedTask creates a Task related to a UserStory
// Available Meta: *TaskDetail
func (us *UserStory) CreateRelatedTask(client *Client, task Task) (*Task, error) {
	if err := requireNonNil("client", client); err != nil {
		return nil, err
	}
	if err := requireNonNil("userStory", us); err != nil {
		return nil, err
	}
	task.UserStory = us.ID
	task.Project = us.Project
	return client.Task.Create(&task)
}

// CloneUserStory clones an existing UserStory with most fields
// Available Meta: UserStoryDetail
func (us *UserStory) CloneUserStory(client *Client) (*UserStory, error) {
	if err := requireNonNil("client", client); err != nil {
		return nil, err
	}
	if err := requireNonNil("userStory", us); err != nil {
		return nil, err
	}
	clone := *us
	clone.ID = 0
	clone.Ref = 0
	clone.Version = 0
	return client.UserStory.Create(&clone)
}

// GetRelatedTasks returns all Tasks related to this UserStory
func (us *UserStory) GetRelatedTasks(client *Client) ([]Task, error) {
	if err := requireNonNil("client", client); err != nil {
		return nil, err
	}
	if err := requireNonNil("userStory", us); err != nil {
		return nil, err
	}
	return client.Task.List(&TasksQueryParams{UserStory: us.ID})
}
