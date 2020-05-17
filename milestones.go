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

// ListMilestones => https://taigaio.github.io/taiga-doc/dist/api.html#Milestones-list
func (s *MilestoneService) ListMilestones(queryParams *MilestonesQueryParams) ([]Milestone, error) {
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
	err := getRequest(s.client, &Milestones, url)
	if err != nil {
		return nil, err
	}

	return Milestones, nil
}

// CreateMilestone => https://taigaio.github.io/taiga-doc/dist/api.html#milestones-create
func (s *MilestoneService) CreateMilestone(milestone Milestone) (*Milestone, error) {
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

	err := postRequest(s.client, &respMilestone, url, milestone)
	if err != nil {
		return nil, err
	}

	return &respMilestone, nil
}

// GetMilestone => https://taigaio.github.io/taiga-doc/dist/api.html#Milestones-get
func (s *MilestoneService) GetMilestone(milestoneID int) (*Milestone, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointMilestones, milestoneID)
	var m Milestone
	err := getRequest(s.client, &m, url)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// EditMilestone edits an Milestone via a PATCH request => https://taigaio.github.io/taiga-doc/dist/api.html#milestones-edit
// Available Meta: MilestoneDetail
func (s *MilestoneService) EditMilestone(milestone Milestone) (*Milestone, error) {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointMilestones, milestone.ID)
	var M Milestone

	if milestone.ID == 0 {
		return nil, errors.New("Passed Milestone does not have an ID yet. Does it exist?")
	}

	err := patchRequest(s.client, &M, url, &milestone)
	if err != nil {
		return nil, err
	}
	return &M, nil
}

// DeleteMilestone => https://taigaio.github.io/taiga-doc/dist/api.html#milestones-delete
func (s *MilestoneService) DeleteMilestone(milestoneID int) error {
	url := s.client.APIURL + fmt.Sprintf("%s/%d", endpointMilestones, milestoneID)
	return deleteRequest(s.client, url)
}

// Stats

// Watch a milestone

// Stop watching a milestone

// List milestone watchers
