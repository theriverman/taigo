package gotaiga

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

// ListUserStories returns all User Stories | https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-list
// Available Meta: *[]UserStoryDetailLIST
func (s *UserStoryService) ListUserStories(queryParameters *UserStoryQueryParams) ([]UserStory, error) {
	url := s.client.APIURL + endpointUserStories
	if queryParameters != nil {
		paramValues, _ := query.Values(queryParameters)
		url = fmt.Sprintf("%s?%s", url, paramValues.Encode())
	} else if s.client.HasDefaultProject() {
		url = url + s.client.GetDefaultProjectAsQueryParam()
	}
	var UserStoryDetailList UserStoryDetailLIST
	err := getRequest(s.client, &UserStoryDetailList, url)
	if err != nil {
		return nil, err
	}
	return UserStoryDetailList.AsUserStory()
}

// CreateUserStory creates a new User Story | https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-create
//
// Available Meta: *UserStoryDetail
func (s *UserStoryService) CreateUserStory(userStory UserStory) (*UserStory, error) {
	url := s.client.APIURL + endpointUserStories
	var newUserStory UserStoryDetail

	// Check for required fields
	// project, subject
	if isEmpty(userStory.Project) || isEmpty(userStory.Subject) {
		return nil, errors.New("A mandatory field is missing. See API documentataion")
	}

	err := postRequest(s.client, &newUserStory, url, userStory)
	if err != nil {
		return nil, err
	}

	return newUserStory.AsUserStory()
}

// GetUserStory -> https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-get
//
// Available Meta: *UserStoryDetailGET
func (s *UserStoryService) GetUserStory(userStoryID int) (*UserStory, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointUserStories, userStoryID)
	var us UserStoryDetailGET
	err := getRequest(s.client, &us, url)
	if err != nil {
		return nil, err
	}
	return us.AsUserStory()
}

// GetUserStoryByRef returns a User Story by Ref -> https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-get-by-ref
//
// The passed userStoryRef should be an int taken from the UserStory's URL
// The passed *Project struct should have at least one of the following fields set:
//		ID 	 (int)
//		Slug (string)
// If none of the above fields are set, an error is returned.
// If both fields are set, *Project.ID will be preferred.
//
// Available Meta: UserStoryDetailGET
func (s *UserStoryService) GetUserStoryByRef(userStoryRef int, project *Project) (*UserStory, error) {
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

	err := getRequest(s.client, &us, url)
	if err != nil {
		return nil, err
	}
	return us.AsUserStory()
}

// EditUserStory sends a PATCH request to edit a User Story -> https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-edit
// Available Meta: UserStoryDetail
func (s *UserStoryService) EditUserStory(userStory UserStory) (*UserStory, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointUserStories, userStory.ID)
	var us UserStoryDetail

	if userStory.ID == 0 {
		return nil, errors.New("Passed UserStory does not have an ID yet. Does it exist?")
	}

	// Taiga OCC
	remoteUS, err := s.GetUserStory(userStory.ID)
	if err != nil {
		return nil, err
	}
	userStory.Version = remoteUS.Version
	err = patchRequest(s.client, &us, url, &userStory)
	if err != nil {
		return nil, err
	}
	return us.AsUserStory()
}

// DeleteUserStory -> https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-delete
func (s *UserStoryService) DeleteUserStory(userStoryID int) error {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointUserStories, userStoryID)
	return deleteRequest(s.client, url)
}

// CloneUserStory clones an existing UserStory with most fields
//
// Available Meta: UserStoryDetail
func (s *UserStoryService) CloneUserStory(srcUserStory UserStory) (*UserStory, error) {
	srcUserStory.ID = 0
	srcUserStory.Ref = 0
	srcUserStory.Version = 0
	return s.CreateUserStory(srcUserStory)
}

// GetRelatedTasks returns all Tasks related to this UserStory
func (s *UserStoryService) GetRelatedTasks(userStoryID int) ([]Task, error) {
	return s.client.Task.List(&TasksQueryParams{UserStory: userStoryID})
}

// CreateRelatedTask creates a Task related to a UserStory
// Available Meta: *TaskDetail
func (s *UserStoryService) CreateRelatedTask(task Task, userStoryID, projectID int) (*Task, error) {
	task.UserStory = userStoryID
	task.Project = projectID
	return s.client.Task.Create(&task)
}

// UserStoryCreateAttachment creates a new UserStory attachment => https://taigaio.github.io/taiga-doc/dist/api.html#user-stories-create-attachment
func (s *UserStoryService) UserStoryCreateAttachment(attachment *Attachment, filePath string) (*Attachment, error) {
	url := s.client.APIURL + endpointUserStories + "/attachments"
	attachment.filePath = filePath
	attachment, err := newfileUploadRequest(s.client, url, attachment)
	if err != nil {
		return nil, err
	}
	return attachment, nil
}
