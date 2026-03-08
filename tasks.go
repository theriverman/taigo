package taigo

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/go-querystring/query"
)

// TaskService is a handle to actions related to Tasks
//
// https://taigaio.github.io/taiga-doc/dist/api.html#tasks
type TaskService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List => https://taigaio.github.io/taiga-doc/dist/api.html#tasks-list
func (s *TaskService) List(queryParams *TasksQueryParams) ([]Task, error) {
	url := s.client.MakeURL(s.Endpoint)
	switch {
	case queryParams != nil:
		paramValues, _ := query.Values(queryParams)
		url = fmt.Sprintf("%s?%s", url, paramValues.Encode())
	case s.defaultProjectID != 0:
		url = url + projectIDQueryParam(s.defaultProjectID)
	}
	var tasks TaskDetailLIST
	_, err := s.client.Request.Get(url, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks.AsTasks()
}

// Create creates a new Task | https://taigaio.github.io/taiga-doc/dist/api.html#tasks-create
// Meta Available: *TaskDetail
func (s *TaskService) Create(task *Task) (*Task, error) {
	url := s.client.MakeURL(s.Endpoint)
	var t TaskDetail

	// Check for required fields
	// project, subject
	if isEmpty(task.Project) || isEmpty(task.Subject) {
		return nil, errors.New("a mandatory field is missing. See API documentataion")
	}

	_, err := s.client.Request.Post(url, &task, &t)
	if err != nil {
		return nil, err
	}
	return t.AsTask()
}

// Get => https://taigaio.github.io/taiga-doc/dist/api.html#tasks-get
func (s *TaskService) Get(taskID int) (*Task, error) {
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
	var t TaskDetailGET
	var url string
	if project == nil {
		return nil, errors.New("project must not be nil")
	}
	if project.ID != 0 {
		url = s.client.MakeURL(fmt.Sprintf("%s/by_ref?ref=%d&project=%d", s.Endpoint, taskRef, project.ID))
	} else if len(project.Slug) > 0 {
		url = s.client.MakeURL(fmt.Sprintf("%s/by_ref?ref=%d&project__slug=%s", s.Endpoint, taskRef, project.Slug))
	} else {
		return nil, errors.New("no ID or Ref defined in passed project struct")
	}

	_, err := s.client.Request.Get(url, &t)
	if err != nil {
		return nil, err
	}
	return t.AsTask()
}

// Edit sends a PATCH request to edit a Task -> https://taigaio.github.io/taiga-doc/dist/api.html#tasks-edit
// Available Meta: TaskDetail
func (s *TaskService) Edit(task *Task) (*Task, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(task.ID))
	var responseTask TaskDetail

	if task.ID == 0 {
		return nil, errors.New("passed Task does not have an ID yet. Does it exist?")
	}

	// Taiga OCC
	remoteTask, err := s.Get(task.ID)
	if err != nil {
		return nil, err
	}
	task.Version = remoteTask.Version
	_, err = s.client.Request.Patch(url, &task, &responseTask)
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
func (s *TaskService) ListAttachments(task interface{}) ([]Attachment, error) {
	t := Task{}
	err := convertStructViaJSON(task, &t)
	if err != nil {
		return nil, err
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
