package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// MilestoneService is a handle to actions related to Milestones
//
// https://taigaio.github.io/taiga-doc/dist/api.html#milestones
type MilestoneService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

type milestoneCreatePayload struct {
	EstimatedFinish string `json:"estimated_finish"`
	EstimatedStart  string `json:"estimated_start"`
	Name            string `json:"name"`
	Project         int    `json:"project"`
}

type milestoneEditPayload struct {
	Closed          *bool   `json:"closed,omitempty"`
	EstimatedFinish *string `json:"estimated_finish,omitempty"`
	EstimatedStart  *string `json:"estimated_start,omitempty"`
	Name            *string `json:"name,omitempty"`
	Project         *int    `json:"project,omitempty"`
}

// List => https://taigaio.github.io/taiga-doc/dist/api.html#Milestones-list
func (s *MilestoneService) List(queryParams *MilestonesQueryParams) ([]Milestone, *MilestoneTotalInfo, error) {
	// prepare url & parameters
	url := s.client.MakeURL(s.Endpoint)
	var err error
	url, err = urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, nil, err
	}
	// execute requests
	var Milestones []Milestone
	httpResponse, err := s.client.Request.Get(url, &Milestones)
	if err != nil {
		return nil, nil, err
	}
	mti := &MilestoneTotalInfo{}
	mti.LoadFromHeaders(httpResponse)

	return Milestones, mti, nil
}

// Create => https://taigaio.github.io/taiga-doc/dist/api.html#milestones-create
//
// Mandatory fields: Project, Name, EstimatedStart, EstimatedFinish
func (s *MilestoneService) Create(milestone *Milestone) (*Milestone, error) {
	if err := requireNonNil("milestone", milestone); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var respMilestone Milestone
	// Check for required fields
	if (isEmpty(milestone.Project) ||
		isEmpty(milestone.Name)) ||
		isEmpty(milestone.EstimatedStart) ||
		isEmpty(milestone.EstimatedFinish) {
		return nil, errors.New("a mandatory field is missing. See API documentataion")
	}
	payload := milestoneCreatePayload{
		EstimatedFinish: milestone.EstimatedFinish,
		EstimatedStart:  milestone.EstimatedStart,
		Name:            milestone.Name,
		Project:         milestone.Project,
	}
	_, err := s.client.Request.Post(url, &payload, &respMilestone)
	if err != nil {
		return nil, err
	}
	return &respMilestone, nil
}

// Get => https://taigaio.github.io/taiga-doc/dist/api.html#Milestones-get
func (s *MilestoneService) Get(milestoneID int) (*Milestone, error) {
	if err := requirePositiveID("milestoneID", milestoneID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(milestoneID))
	var m Milestone
	_, err := s.client.Request.Get(url, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// Edit edits an Milestone via a PATCH request => https://taigaio.github.io/taiga-doc/dist/api.html#milestones-edit
// Available Meta: MilestoneDetail
func (s *MilestoneService) Edit(milestone *Milestone) (*Milestone, error) {
	if err := requireNonNil("milestone", milestone); err != nil {
		return nil, err
	}
	if err := requirePositiveID("milestoneID", milestone.ID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(milestone.ID))

	var m Milestone
	payload := milestoneEditPayload{
		Closed:          ptr(milestone.Closed),
		EstimatedFinish: ptr(milestone.EstimatedFinish),
		EstimatedStart:  ptr(milestone.EstimatedStart),
		Name:            ptr(milestone.Name),
		Project:         ptr(milestone.Project),
	}
	_, err := s.client.Request.Patch(url, &payload, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// Update is an alias for Edit.
func (s *MilestoneService) Update(milestone *Milestone) (*Milestone, error) {
	return s.Edit(milestone)
}

// Delete => https://taigaio.github.io/taiga-doc/dist/api.html#milestones-delete
func (s *MilestoneService) Delete(milestoneID int) (*http.Response, error) {
	if err := requirePositiveID("milestoneID", milestoneID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(milestoneID))
	return s.client.Request.Delete(url)
}

// Stats

// Watch a milestone

// Stop watching a milestone

// List milestone watchers
