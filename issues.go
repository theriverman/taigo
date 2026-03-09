package taigo

import (
	"errors"
	"fmt"
	"net/http"
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
	case s.defaultProjectID != 0:
		url = url + projectIDQueryParam(s.defaultProjectID)
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
func (s *IssueService) CreateAttachment(attachment *Attachment, issue *Issue) (*Attachment, error) {
	url := s.client.MakeURL(s.Endpoint, "attachments")
	return newfileUploadRequest(s.client, url, attachment, issue)
}

// GetAttachment retrives an Issue attachment by its ID => https://taigaio.github.io/taiga-doc/dist/api.html#issues-get-attachment
func (s *IssueService) GetAttachment(attachmentID int) (*Attachment, error) {
	a, err := getAttachmentForEndpoint(s.client, attachmentID, s.Endpoint)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// ListAttachments returns a list of Issue attachments => https://taigaio.github.io/taiga-doc/dist/api.html#issues-list-attachments
func (s *IssueService) ListAttachments(issue any) ([]Attachment, error) {
	i := Issue{}
	err := convertStructViaJSON(issue, &i)
	if err != nil {
		return nil, err
	}

	queryParams := attachmentsQueryParams{
		endpointURI: s.Endpoint,
		ObjectID:    i.ID,
		Project:     i.Project,
	}

	attachments, err := listAttachmentsForEndpoint(s.client, &queryParams)
	if err != nil {
		return nil, err
	}
	return attachments, nil
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

// GetByRef returns an Issue by Ref -> https://taigaio.github.io/taiga-doc/dist/api.html#issues-get-by-ref
func (s *IssueService) GetByRef(issueRef int, project *Project) (*Issue, error) {
	var issue IssueDetailGET
	var url string
	if project == nil {
		return nil, errors.New("project must not be nil")
	}

	switch {
	case project.ID != 0:
		url = s.client.MakeURL(fmt.Sprintf("%s/by_ref?ref=%d&project=%d", s.Endpoint, issueRef, project.ID))
	case len(project.Slug) > 0:
		url = s.client.MakeURL(fmt.Sprintf("%s/by_ref?ref=%d&project__slug=%s", s.Endpoint, issueRef, project.Slug))
	default:
		return nil, errors.New("no ID or Ref defined in passed project struct")
	}

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
		return nil, errors.New("passed Issue does not have an ID yet. Does it exist?")
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

// Update is an alias for Edit.
func (s *IssueService) Update(issue *Issue) (*Issue, error) {
	return s.Edit(issue)
}

// Create creates a new Issue | https://taigaio.github.io/taiga-doc/dist/api.html#issues-create
//
// Available Meta: *IssueDetail
func (s *IssueService) Create(issue *Issue) (*Issue, error) {
	url := s.client.MakeURL(s.Endpoint)
	var issueDetail IssueDetail

	// Check for required fields
	// project, subject
	if isEmpty(issue.Project) || isEmpty(issue.Subject) {
		return nil, errors.New("a mandatory field is missing. See API documentataion")
	}

	_, err := s.client.Request.Post(url, &issue, &issueDetail)
	if err != nil {
		return nil, err
	}

	return issueDetail.AsIssue()
}

// Delete -> https://taigaio.github.io/taiga-doc/dist/api.html#issues-delete
func (s *IssueService) Delete(issueID int) (*http.Response, error) {
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(issueID))
	return s.client.Request.Delete(url)
}
