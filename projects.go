package taigo

import (
	"errors"
	"fmt"

	"github.com/google/go-querystring/query"
)

const endpointProjects = "/projects"

// ProjectService is a handle to actions related to Projects
//
// https://taigaio.github.io/taiga-doc/dist/api.html#projects
type ProjectService struct {
	client *Client
}

// List -> https://taigaio.github.io/taiga-doc/dist/api.html#projects-list
//
// The results can be filtered by passing in a ProjectListQueryFilter struct
func (s *ProjectService) List(queryParameters *ProjectsQueryParameters) (*ProjectsList, error) {
	/*
		The results can be filtered using the following parameters:
		  * Member
		  * Members
		  * IsLookingForPeople
		  * IsFeatured
		  * IsBacklogActivated
		  * IsKanbanActivated

		The results can be ordered using the order_by parameter with the values:
		  * memberships__user_order
		  * total_fans
		  * total_fans_last_week
		  * total_fans_last_month
		  * total_fans_last_year
		  * total_activity
		  * total_activity_last_week
		  * total_activity_last_month
		  * total_activity_last_year
	*/

	url := s.client.APIURL + endpointProjects
	if queryParameters != nil {
		paramValues, _ := query.Values(queryParameters)
		url = fmt.Sprintf("%s?%s", url, paramValues.Encode())
	}
	var projects ProjectsList

	err := getRequest(s.client, &projects, url)
	if err != nil {
		return nil, err
	}
	return &projects, nil
}

// Create -> https://taigaio.github.io/taiga-doc/dist/api.html#projects-create
func (s *ProjectService) Create(project *Project) (*Project, error) {
	url := s.client.APIURL + endpointProjects
	var p Project

	// Check for required fields
	// name, description
	if isEmpty(project.Name) || isEmpty(project.Description) {
		return nil, errors.New("A mandatory field is missing. See API documentataion")
	}

	err := postRequest(s.client, &p, url, project)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Get -> https://taigaio.github.io/taiga-doc/dist/api.html#projects-get
func (s *ProjectService) Get(projectID int) (*Project, error) {
	url := s.client.APIURL + fmt.Sprintf("/%s/%d", endpointProjects, projectID)
	var p ProjectDetail

	err := getRequest(s.client, &p, url)
	if err != nil {
		return nil, err
	}
	return p.AsProject()
}

// GetBySlug -> https://taigaio.github.io/taiga-doc/dist/api.html#projects-get-by-slug
func (s *ProjectService) GetBySlug(slug string) (*Project, error) {
	url := s.client.APIURL + fmt.Sprintf("/%s/by_slug?slug=%s", endpointProjects, slug)
	var p ProjectDetail

	err := getRequest(s.client, &p, url)
	if err != nil {
		return nil, err
	}
	return p.AsProject()
}

// Edit edits an Project via a PATCH request => https://taigaio.github.io/taiga-doc/dist/api.html#projects-edit
// Available Meta: ProjectDetail
func (s *ProjectService) Edit(project Project) (*Project, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointProjects, project.ID)
	var projectDetail ProjectDetail

	if project.ID == 0 {
		return nil, errors.New("Passed Project does not have an ID yet. Does it exist?")
	}

	err := patchRequest(s.client, &projectDetail, url, &project)
	if err != nil {
		return nil, err
	}
	return projectDetail.AsProject()
}

// Delete => https://taigaio.github.io/taiga-doc/dist/api.html#projects-delete
func (s *ProjectService) Delete(projectID int) error {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointProjects, projectID)
	return deleteRequest(s.client, url)
}
