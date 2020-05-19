package taigo

import (
	"errors"
	"fmt"

	"github.com/google/go-querystring/query"
)

const endpointUserStories = "/userstories"

// UserStoryService is a handle to actions related to UserStories
//
// https://taigaio.github.io/taiga-doc/dist/api.html#user-stories
type UserStoryService struct {
	client *Client
}

// List returns all User Stories | https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-list
// Available Meta: *[]UserStoryDetailLIST
func (s *UserStoryService) List(queryParameters *UserStoryQueryParams) ([]UserStory, error) {
	url := s.client.APIURL + endpointUserStories
	if queryParameters != nil {
		paramValues, _ := query.Values(queryParameters)
		url = fmt.Sprintf("%s?%s", url, paramValues.Encode())
	} else if s.client.HasDefaultProject() {
		url = url + s.client.GetDefaultProjectAsQueryParam()
	}
	var UserStoryDetailList UserStoryDetailLIST
	err := s.client.Request.Get(url, &UserStoryDetailList)
	if err != nil {
		return nil, err
	}
	return UserStoryDetailList.AsUserStory()
}

// Create creates a new User Story | https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-create
//
// Available Meta: *UserStoryDetail
func (s *UserStoryService) Create(userStory *UserStory) (*UserStory, error) {
	url := s.client.APIURL + endpointUserStories
	var newUserStory UserStoryDetail

	// Check for required fields
	// project, subject
	if isEmpty(userStory.Project) || isEmpty(userStory.Subject) {
		return nil, errors.New("A mandatory field is missing. See API documentataion")
	}

	err := s.client.Request.Post(url, &userStory, &newUserStory)
	if err != nil {
		return nil, err
	}

	return newUserStory.AsUserStory()
}

// Get -> https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-get
//
// Available Meta: *UserStoryDetailGET
func (s *UserStoryService) Get(userStoryID int) (*UserStory, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointUserStories, userStoryID)
	var us UserStoryDetailGET
	err := s.client.Request.Get(url, &us)
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
		url = s.client.APIURL + fmt.Sprintf("%s/by_ref?ref=%d&project=%d", endpointEpics, userStoryRef, project.ID)
		break
	case len(project.Slug) > 0:
		url = s.client.APIURL + fmt.Sprintf("%s/by_ref?ref=%d&project__slug=%s", endpointEpics, userStoryRef, project.Slug)
		break
	default:
		return nil, errors.New("No ID or Ref defined in passed project struct")
	}

	err := s.client.Request.Get(url, &us)
	if err != nil {
		return nil, err
	}
	return us.AsUserStory()
}

// Edit sends a PATCH request to edit a User Story -> https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-edit
// Available Meta: UserStoryDetail
func (s *UserStoryService) Edit(userStory *UserStory) (*UserStory, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointUserStories, userStory.ID)
	var responseUS UserStoryDetail

	if userStory.ID == 0 {
		return nil, errors.New("Passed UserStory does not have an ID yet. Does it exist?")
	}

	// Taiga OCC
	remoteUS, err := s.Get(userStory.ID)
	if err != nil {
		return nil, err
	}
	userStory.Version = remoteUS.Version
	err = s.client.Request.Patch(url, &userStory, &responseUS)
	if err != nil {
		return nil, err
	}
	return responseUS.AsUserStory()
}

// Delete -> https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-delete
func (s *UserStoryService) Delete(userStoryID int) error {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointUserStories, userStoryID)
	return s.client.Request.Delete(url)
}

// CreateAttachment creates a new UserStory attachment => https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-create-attachment
func (s *UserStoryService) CreateAttachment(attachment *Attachment, task *Task) (*Attachment, error) {
	url := s.client.APIURL + endpointTasks + "/attachments"
	return newfileUploadRequest(s.client, url, attachment, task)
}

/*
	Advanced Operations
*/

// RelateToEpic relates the UserStory to an Epic via an EpicID
//
// TaigaClient must be a pointer to taiga.Client
// EpicID must be an int to desired Epic
func (us *UserStory) RelateToEpic(TaigaClient *Client, EpicID int) (*EpicRelatedUserStoryDetail, error) {
	if us.ID == 0 {
		return nil, fmt.Errorf("UserStory must be created before relating it to an Epic. UserStory.ID was 0")
	}
	return TaigaClient.Epic.CreateRelatedUserStory(EpicID, us.ID)
}

// Clone clones an existing UserStory with most fields
//
// Available Meta: UserStoryDetail
func (s *UserStoryService) Clone(srcUserStory *UserStory) (*UserStory, error) {
	srcUserStory.ID = 0
	srcUserStory.Ref = 0
	srcUserStory.Version = 0
	return s.Create(srcUserStory)
}

// ListRelatedTasks returns all Tasks related to this UserStory
func (us *UserStory) ListRelatedTasks(TaigaClient *Client, userStoryID int) ([]Task, error) {
	return TaigaClient.Task.List(&TasksQueryParams{UserStory: userStoryID})
}

// CreateRelatedTask creates a Task related to a UserStory
// Available Meta: *TaskDetail
func (us *UserStory) CreateRelatedTask(TaigaClient *Client, task Task) (*Task, error) {
	task.UserStory = us.ID
	task.Project = us.Project
	return TaigaClient.Task.Create(&task)
}

// CloneUserStory clones an existing UserStory with most fields
// Available Meta: UserStoryDetail
func (us *UserStory) CloneUserStory(TaigaClient *Client) (*UserStory, error) {
	us.ID = 0
	us.Ref = 0
	us.Version = 0
	return TaigaClient.UserStory.Create(us)
}

// GetRelatedTasks returns all Tasks related to this UserStory
func (us *UserStory) GetRelatedTasks(TaigaClient *Client) ([]Task, error) {
	return TaigaClient.Task.List(&TasksQueryParams{UserStory: us.ID})
}
