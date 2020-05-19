package taigo

import (
	"errors"
	"fmt"

	"github.com/google/go-querystring/query"
)

const endpointMilestones = "/milestones"

// MilestoneService is a handle to actions related to Milestones
//
// https://taigaio.github.io/taiga-doc/dist/api.html#milestones
type MilestoneService struct {
	client *Client
}

// List => https://taigaio.github.io/taiga-doc/dist/api.html#Milestones-list
func (s *MilestoneService) List(queryParams *MilestonesQueryParams) ([]Milestone, error) {
	// prepare url & parameters
	url := s.client.APIURL + endpointMilestones
	if queryParams != nil {
		paramValues, _ := query.Values(queryParams)
		url = fmt.Sprintf("%s?%s", url, paramValues.Encode())
	} else if s.client.HasDefaultProject() {
		url = url + s.client.GetDefaultProjectAsQueryParam()
	}
	// execute requests
	var Milestones []Milestone
	err := s.client.Request.GetRequest(url, &Milestones)
	if err != nil {
		return nil, err
	}

	return Milestones, nil
}

// Create => https://taigaio.github.io/taiga-doc/dist/api.html#milestones-create
//
// Mandatory fields: Project, Name, EstimatedStart, EstimatedFinish
func (s *MilestoneService) Create(milestone *Milestone) (*Milestone, error) {
	url := s.client.APIURL + endpointMilestones
	var respMilestone Milestone

	// Check for required fields
	// project, name, estimated_start, estimated_finish
	if (isEmpty(milestone.Project) ||
		isEmpty(milestone.Name)) ||
		isEmpty(milestone.EstimatedStart) ||
		isEmpty(milestone.EstimatedFinish) {
		return nil, errors.New("A mandatory field is missing. See API documentataion")
	}

	err := s.client.Request.PostRequest(url, &milestone, &respMilestone)
	if err != nil {
		return nil, err
	}

	return &respMilestone, nil
}

// Get => https://taigaio.github.io/taiga-doc/dist/api.html#Milestones-get
func (s *MilestoneService) Get(milestoneID int) (*Milestone, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointMilestones, milestoneID)
	var m Milestone
	err := s.client.Request.GetRequest(url, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// Edit edits an Milestone via a PATCH request => https://taigaio.github.io/taiga-doc/dist/api.html#milestones-edit
// Available Meta: MilestoneDetail
func (s *MilestoneService) Edit(milestone *Milestone) (*Milestone, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointMilestones, milestone.ID)
	var M Milestone

	if milestone.ID == 0 {
		return nil, errors.New("Passed Milestone does not have an ID yet. Does it exist?")
	}

	err := s.client.Request.PatchRequest(url, &milestone, &M)
	if err != nil {
		return nil, err
	}
	return &M, nil
}

// Delete => https://taigaio.github.io/taiga-doc/dist/api.html#milestones-delete
func (s *MilestoneService) Delete(milestoneID int) error {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointMilestones, milestoneID)
	return s.client.Request.DeleteRequest(url)
}

// Stats

// Watch a milestone

// Stop watching a milestone

// List milestone watchers
