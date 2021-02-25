package taigo

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/google/go-querystring/query"
)

// IssueService is a handle to actions related to Issues
//
// https://taigaio.github.io/taiga-doc/dist/api.html#issues
type IssueService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// List => https://taigaio.github.io/taiga-doc/dist/api.html#issues-list
func (s *IssueService) List(queryParams *IssueQueryParams) ([]Issue, error) {
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
	var issues IssueDetailLIST
	_, err := s.client.Request.Get(url, &issues)
	if err != nil {
		return nil, err
	}

	return issues.AsIssues()
}

// CreateAttachment creates a new Issue attachment => https://taigaio.github.io/taiga-doc/dist/api.html#issues-create-attachment
func (s *IssueService) CreateAttachment(attachment *Attachment, task *Task) (*Attachment, error) {
	url := s.client.MakeURL(s.Endpoint, "attachments")
	return newfileUploadRequest(s.client, url, attachment, task)
}

// Get -> https://taigaio.github.io/taiga-doc/dist/api.html#issues-get
//
// Available Meta: *IssueDetailGET
func (s *IssueService) Get(issueID int) (*Issue, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(issueID))
	var issue IssueDetailGET
	_, err := s.client.Request.Get(url, &issue)
	if err != nil {
		return nil, err
	}
	return issue.AsIssue()
}

// Edit sends a PATCH request to edit a Issue -> https://taigaio.github.io/taiga-doc/dist/api.html#issues-edit
// Available Meta: IssueDetail
func (s *IssueService) Edit(issue *Issue) (*Issue, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(issue.ID))
	var responseIssue IssueDetail

	if issue.ID == 0 {
		return nil, errors.New("Passed Issue does not have an ID yet. Does it exist?")
	}

	// Taiga OCC
	remoteIssue, err := s.Get(issue.ID)
	if err != nil {
		return nil, err
	}

	issue.Version = remoteIssue.Version
	_, err = s.client.Request.Patch(url, &issue, &responseIssue)
	if err != nil {
		return nil, err
	}
	return responseIssue.AsIssue()
}
