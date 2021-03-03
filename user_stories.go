package taigo

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/go-querystring/query"
)

// UserStoryService is a handle to actions related to UserStories
//
// https://taigaio.github.io/taiga-doc/dist/api.html#user-stories
type UserStoryService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List returns all User Stories | https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-list
// Available Meta: *[]UserStoryDetailLIST
func (s *UserStoryService) List(queryParams *UserStoryQueryParams) ([]UserStory, error) {
	url := s.client.MakeURL(s.Endpoint)
	switch {
	case queryParams != nil:
		paramValues, _ := query.Values(queryParams)
		url = fmt.Sprintf("%s?%s", url, paramValues.Encode())
		break
	case s.defaultProjectID != 0:
		url = url + projectIDQueryParam(s.defaultProjectID)
		break
	}
	var userstories UserStoryDetailLIST
	_, err := s.client.Request.Get(url, &userstories)
	if err != nil {
		return nil, err
	}
	return userstories.AsUserStory()
}

// Create creates a new User Story | https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-create
//
// Available Meta: *UserStoryDetail
func (s *UserStoryService) Create(userStory *UserStory) (*UserStory, error) {
	url := s.client.MakeURL(s.Endpoint)
	var us UserStoryDetail

	// Check for required fields
	// project, subject
	if isEmpty(userStory.Project) || isEmpty(userStory.Subject) {
		return nil, errors.New("A mandatory field is missing. See API documentataion")
	}

	_, err := s.client.Request.Post(url, &userStory, &us)
	if err != nil {
		return nil, err
	}

	return us.AsUserStory()
}

// Get -> https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-get
//
// Available Meta: *UserStoryDetailGET
func (s *UserStoryService) Get(userStoryID int) (*UserStory, error) {
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
//		ID 	 (int)
//		Slug (string)
// If none of the above fields are set, an error is returned.
// If both fields are set, *Project.ID will be preferred.
//
// Available Meta: UserStoryDetailGET
func (s *UserStoryService) GetByRef(userStoryRef int, project *Project) (*UserStory, error) {
	var us UserStoryDetailGET
	var url string

	switch {
	case project.ID != 0:
		url = s.client.MakeURL(fmt.Sprintf("%s/by_ref?ref=%d&project=%d", s.Endpoint, userStoryRef, project.ID))
		break
	case len(project.Slug) > 0:
		url = s.client.MakeURL(fmt.Sprintf("%s/by_ref?ref=%d&project__slug=%s", s.Endpoint, userStoryRef, project.Slug))
		break
	default:
		return nil, errors.New("No ID or Ref defined in passed project struct")
	}

	_, err := s.client.Request.Get(url, &us)
	if err != nil {
		return nil, err
	}
	return us.AsUserStory()
}

// Edit sends a PATCH request to edit a User Story -> https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-edit
// Available Meta: UserStoryDetail
func (s *UserStoryService) Edit(us *UserStory) (*UserStory, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(us.ID))
	var responseUS UserStoryDetail

	if us.ID == 0 {
		return nil, errors.New("Passed UserStory does not have an ID yet. Does it exist?")
	}

	// Taiga OCC
	remoteUS, err := s.Get(us.ID)
	if err != nil {
		return nil, err
	}
	us.Version = remoteUS.Version
	_, err = s.client.Request.Patch(url, &us, &responseUS)
	if err != nil {
		return nil, err
	}
	return responseUS.AsUserStory()
}

// Delete -> https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-delete
func (s *UserStoryService) Delete(usID int) (*http.Response, error) {
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
	if us.ID == 0 {
		return nil, fmt.Errorf("UserStory must be created before relating it to an Epic. UserStory.ID was 0")
	}
	return client.Epic.CreateRelatedUserStory(epicID, us.ID)
}

// Clone clones an existing UserStory with most fields
//
// Available Meta: UserStoryDetail
func (s *UserStoryService) Clone(srcUS *UserStory) (*UserStory, error) {
	srcUS.ID = 0
	srcUS.Ref = 0
	srcUS.Version = 0
	return s.Create(srcUS)
}

// ListRelatedTasks returns all Tasks related to this UserStory
func (us *UserStory) ListRelatedTasks(client *Client, userStoryID int) ([]Task, error) {
	return client.Task.List(&TasksQueryParams{UserStory: userStoryID})
}

// CreateRelatedTask creates a Task related to a UserStory
// Available Meta: *TaskDetail
func (us *UserStory) CreateRelatedTask(client *Client, task Task) (*Task, error) {
	task.UserStory = us.ID
	task.Project = us.Project
	return client.Task.Create(&task)
}

// CloneUserStory clones an existing UserStory with most fields
// Available Meta: UserStoryDetail
func (us *UserStory) CloneUserStory(client *Client) (*UserStory, error) {
	us.ID = 0
	us.Ref = 0
	us.Version = 0
	return client.UserStory.Create(us)
}

// GetRelatedTasks returns all Tasks related to this UserStory
func (us *UserStory) GetRelatedTasks(client *Client) ([]Task, error) {
	return client.Task.List(&TasksQueryParams{UserStory: us.ID})
}
