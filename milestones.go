package taigo

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/go-querystring/query"
)

// MilestoneService is a handle to actions related to Milestones
//
// https://taigaio.github.io/taiga-doc/dist/api.html#milestones
type MilestoneService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List => https://taigaio.github.io/taiga-doc/dist/api.html#Milestones-list
func (s *MilestoneService) List(queryParams *MilestonesQueryParams) ([]Milestone, *MilestoneTotalInfo, error) {
	// prepare url & parameters
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
	url := s.client.MakeURL(s.Endpoint)
	var respMilestone Milestone
	// Check for required fields
	if (isEmpty(milestone.Project) ||
		isEmpty(milestone.Name)) ||
		isEmpty(milestone.EstimatedStart) ||
		isEmpty(milestone.EstimatedFinish) {
		return nil, errors.New("A mandatory field is missing. See API documentataion")
	}
	_, err := s.client.Request.Post(url, &milestone, &respMilestone)
	if err != nil {
		return nil, err
	}
	return &respMilestone, nil
}

// Get => https://taigaio.github.io/taiga-doc/dist/api.html#Milestones-get
func (s *MilestoneService) Get(milestoneID int) (*Milestone, error) {
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
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(milestone.ID))

	var m Milestone
	if milestone.ID == 0 {
		return nil, errors.New("Passed Milestone does not have an ID yet. Does it exist?")
	}
	_, err := s.client.Request.Patch(url, &milestone, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// Delete => https://taigaio.github.io/taiga-doc/dist/api.html#milestones-delete
func (s *MilestoneService) Delete(milestoneID int) (*http.Response, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(milestoneID))
	return s.client.Request.Delete(url)
}

// Stats

// Watch a milestone

// Stop watching a milestone

// List milestone watchers
