package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// TaskService is a handle to actions related to Tasks
//
// https://taigaio.github.io/taiga-doc/dist/api.html#tasks
type TaskService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

type taskCreatePayload struct {
	AssignedTo        int      `json:"assigned_to,omitempty"`
	BlockedNote       string   `json:"blocked_note,omitempty"`
	Description       string   `json:"description,omitempty"`
	ExternalReference []string `json:"external_reference,omitempty"`
	IsBlocked         bool     `json:"is_blocked,omitempty"`
	IsIocaine         bool     `json:"is_iocaine,omitempty"`
	Milestone         int      `json:"milestone,omitempty"`
	Project           int      `json:"project"`
	Status            int      `json:"status,omitempty"`
	Subject           string   `json:"subject"`
	Tags              []string `json:"tags,omitempty"`
	TaskboardOrder    int      `json:"taskboard_order,omitempty"`
	UsOrder           int      `json:"us_order,omitempty"`
	UserStory         int      `json:"user_story,omitempty"`
	Watchers          []int    `json:"watchers,omitempty"`
}

// TaskPatch represents an explicit PATCH payload for tasks.
// Pointer fields allow intentionally setting zero-values (false, 0, "").
type TaskPatch struct {
	AssignedTo        *int      `json:"assigned_to,omitempty"`
	BlockedNote       *string   `json:"blocked_note,omitempty"`
	Description       *string   `json:"description,omitempty"`
	ExternalReference *[]string `json:"external_reference,omitempty"`
	IsBlocked         *bool     `json:"is_blocked,omitempty"`
	IsIocaine         *bool     `json:"is_iocaine,omitempty"`
	Milestone         *int      `json:"milestone,omitempty"`
	Project           *int      `json:"project,omitempty"`
	Status            *int      `json:"status,omitempty"`
	Subject           *string   `json:"subject,omitempty"`
	Tags              *[]string `json:"tags,omitempty"`
	TaskboardOrder    *int      `json:"taskboard_order,omitempty"`
	UsOrder           *int      `json:"us_order,omitempty"`
	UserStory         *int      `json:"user_story,omitempty"`
	Version           int       `json:"version"`
	Watchers          *[]int    `json:"watchers,omitempty"`
}

// List => https://taigaio.github.io/taiga-doc/dist/api.html#tasks-list
func (s *TaskService) List(queryParams *TasksQueryParams) ([]Task, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}
	var tasks TaskDetailLIST
	_, err = s.client.Request.Get(url, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks.AsTasks()
}

// Create creates a new Task | https://taigaio.github.io/taiga-doc/dist/api.html#tasks-create
// Meta Available: *TaskDetail
func (s *TaskService) Create(task *Task) (*Task, error) {
	if err := requireNonNil("task", task); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var t TaskDetail
	projectID, err := resolveProjectID(task.Project, s.defaultProjectID, "project")
	if err != nil {
		return nil, err
	}

	// Check for required fields
	// project, subject
	if isEmpty(task.Subject) {
		return nil, errors.New("a mandatory field is missing. See API documentataion")
	}

	payload := taskCreatePayload{
		AssignedTo:        task.AssignedTo,
		BlockedNote:       task.BlockedNote,
		Description:       task.Description,
		ExternalReference: task.ExternalReference,
		IsBlocked:         task.IsBlocked,
		IsIocaine:         task.IsIocaine,
		Milestone:         task.Milestone,
		Project:           projectID,
		Status:            task.Status,
		Subject:           task.Subject,
		Tags:              tagsToNames(task.Tags),
		TaskboardOrder:    task.TaskboardOrder,
		UsOrder:           task.UsOrder,
		UserStory:         task.UserStory,
		Watchers:          task.Watchers,
	}

	_, err = s.client.Request.Post(url, &payload, &t)
	if err != nil {
		return nil, err
	}
	return t.AsTask()
}

// Get => https://taigaio.github.io/taiga-doc/dist/api.html#tasks-get
func (s *TaskService) Get(taskID int) (*Task, error) {
	if err := requirePositiveID("taskID", taskID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(taskID))
	var t TaskDetailGET
	_, err := s.client.Request.Get(url, &t)
	if err != nil {
		return nil, err
	}
	return t.AsTask()
}

// GetByRef => https://taigaio.github.io/taiga-doc/dist/api.html#tasks-get-by-ref
func (s *TaskService) GetByRef(taskRef int, project *Project) (*Task, error) {
	if err := requirePositiveID("taskRef", taskRef); err != nil {
		return nil, err
	}
	var t TaskDetailGET
	var url string
	type byRefQueryParams struct {
		Ref         int    `url:"ref"`
		Project     int    `url:"project,omitempty"`
		ProjectSlug string `url:"project__slug,omitempty"`
	}
	queryParams := &byRefQueryParams{Ref: taskRef}
	if project != nil && project.ID != 0 {
		queryParams.Project = project.ID
	} else if project != nil && len(project.Slug) > 0 {
		queryParams.ProjectSlug = project.Slug
	} else if s.defaultProjectID > 0 {
		queryParams.Project = s.defaultProjectID
	} else {
		return nil, errors.New("no project ID/slug provided and no mapped default project ID set")
	}
	url, err := appendQueryParams(s.client.MakeURL(s.Endpoint, "by_ref"), queryParams)
	if err != nil {
		return nil, err
	}

	_, err = s.client.Request.Get(url, &t)
	if err != nil {
		return nil, err
	}
	return t.AsTask()
}

// Edit sends a PATCH request to edit a Task -> https://taigaio.github.io/taiga-doc/dist/api.html#tasks-edit
// Available Meta: TaskDetail
func (s *TaskService) Edit(task *Task) (*Task, error) {
	if err := requireNonNil("task", task); err != nil {
		return nil, err
	}

	if task.ID == 0 {
		return nil, errors.New("passed Task does not have an ID yet. Does it exist?")
	}
	if task.Version == 0 {
		return nil, errors.New("version is required for task edit")
	}

	patchPayload := map[string]any{
		"version": task.Version,
	}
	if task.AssignedTo != 0 {
		patchPayload["assigned_to"] = task.AssignedTo
	}
	if task.BlockedNote != "" {
		patchPayload["blocked_note"] = task.BlockedNote
	}
	if task.Description != "" {
		patchPayload["description"] = task.Description
	}
	if task.IsBlocked {
		patchPayload["is_blocked"] = task.IsBlocked
	}
	if task.IsIocaine {
		patchPayload["is_iocaine"] = task.IsIocaine
	}
	if task.Milestone != 0 {
		patchPayload["milestone"] = task.Milestone
	}
	if task.Project != 0 {
		patchPayload["project"] = task.Project
	}
	if task.Status != 0 {
		patchPayload["status"] = task.Status
	}
	if task.Subject != "" {
		patchPayload["subject"] = task.Subject
	}
	if task.TaskboardOrder != 0 {
		patchPayload["taskboard_order"] = task.TaskboardOrder
	}
	if task.UsOrder != 0 {
		patchPayload["us_order"] = task.UsOrder
	}
	if task.UserStory != 0 {
		patchPayload["user_story"] = task.UserStory
	}
	if task.ExternalReference != nil {
		externalRef := append([]string(nil), task.ExternalReference...)
		patchPayload["external_reference"] = externalRef
	}
	if task.Tags != nil {
		tags := tagsToNames(task.Tags)
		if tags == nil {
			tags = []string{}
		}
		patchPayload["tags"] = tags
	}
	if task.Watchers != nil {
		watchers := append([]int(nil), task.Watchers...)
		patchPayload["watchers"] = watchers
	}
	if len(patchPayload) == 1 {
		return nil, errors.New("no updatable task fields were provided; use Patch for explicit zero-value updates")
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(task.ID))
	var responseTask TaskDetail
	_, err := s.client.Request.Patch(url, &patchPayload, &responseTask)
	if err != nil {
		return nil, err
	}
	return responseTask.AsTask()
}

// Patch sends an explicit PATCH payload to edit a task.
func (s *TaskService) Patch(taskID int, patch *TaskPatch) (*Task, error) {
	if err := requireNonNil("patch", patch); err != nil {
		return nil, err
	}
	if err := requirePositiveID("taskID", taskID); err != nil {
		return nil, err
	}
	if patch.Version == 0 {
		return nil, errors.New("version is required for task patch")
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(taskID))
	var responseTask TaskDetail
	_, err := s.client.Request.Patch(url, patch, &responseTask)
	if err != nil {
		return nil, err
	}
	return responseTask.AsTask()
}

// Update is an alias for Edit.
func (s *TaskService) Update(task *Task) (*Task, error) {
	return s.Edit(task)
}

// Delete -> https://taigaio.github.io/taiga-doc/dist/api.html#tasks-delete
func (s *TaskService) Delete(taskID int) (*http.Response, error) {
	if err := requirePositiveID("taskID", taskID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(taskID))
	return s.client.Request.Delete(url)
}

// GetAttachment retrives a Task attachment by its ID => https://taigaio.github.io/taiga-doc/dist/api.html#tasks-get-attachment
func (s *TaskService) GetAttachment(attachmentID int) (*Attachment, error) {
	a, err := getAttachmentForEndpoint(s.client, attachmentID, s.Endpoint)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// ListAttachments returns a list of Task attachments => https://taigaio.github.io/taiga-doc/dist/api.html#tasks-list-attachments
func (s *TaskService) ListAttachments(task any) ([]Attachment, error) {
	if err := requireNonNil("task", task); err != nil {
		return nil, err
	}
	t := Task{}
	err := convertStructViaJSON(task, &t)
	if err != nil {
		return nil, err
	}
	if t.ID == 0 || t.Project == 0 {
		return nil, errors.New("task id and project are required to list attachments")
	}

	queryParams := attachmentsQueryParams{
		endpointURI: s.Endpoint,
		ObjectID:    t.ID,
		Project:     t.Project,
	}

	attachments, err := listAttachmentsForEndpoint(s.client, &queryParams)
	if err != nil {
		return nil, err
	}
	return attachments, nil
}

// CreateAttachment creates a new Task attachment => https://taigaio.github.io/taiga-doc/dist/api.html#tasks-create-attachment
func (s *TaskService) CreateAttachment(attachment *Attachment, task *Task) (*Attachment, error) {
	url := s.client.MakeURL(s.Endpoint, "attachments")
	return newfileUploadRequest(s.client, url, attachment, task)
}
