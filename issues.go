package taigo

import (
	"errors"
	"net/http"
	"strconv"
)

// IssueService is a handle to actions related to Issues
//
// https://taigaio.github.io/taiga-doc/dist/api.html#issues
type IssueService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

type issueCreatePayload struct {
	AssignedTo    int      `json:"assigned_to,omitempty"`
	BlockedNote   string   `json:"blocked_note,omitempty"`
	Description   string   `json:"description,omitempty"`
	IsBlocked     bool     `json:"is_blocked,omitempty"`
	Milestone     int      `json:"milestone,omitempty"`
	Owner         int      `json:"owner,omitempty"`
	Priority      int      `json:"priority,omitempty"`
	Project       int      `json:"project"`
	Severity      int      `json:"severity,omitempty"`
	Status        int      `json:"status,omitempty"`
	Subject       string   `json:"subject"`
	Tags          []string `json:"tags,omitempty"`
	Type          int      `json:"type,omitempty"`
	Watchers      []int    `json:"watchers,omitempty"`
	DueDate       string   `json:"due_date,omitempty"`
	DueDateReason string   `json:"due_date_reason,omitempty"`
	DueDateStatus string   `json:"due_date_status,omitempty"`
}

// IssuePatch represents an explicit PATCH payload for issues.
// Pointer fields allow intentionally setting zero-values (false, 0, "").
type IssuePatch struct {
	AssignedTo    *int      `json:"assigned_to,omitempty"`
	BlockedNote   *string   `json:"blocked_note,omitempty"`
	Description   *string   `json:"description,omitempty"`
	IsBlocked     *bool     `json:"is_blocked,omitempty"`
	Milestone     *int      `json:"milestone,omitempty"`
	Owner         *int      `json:"owner,omitempty"`
	Priority      *int      `json:"priority,omitempty"`
	Project       *int      `json:"project,omitempty"`
	Severity      *int      `json:"severity,omitempty"`
	Status        *int      `json:"status,omitempty"`
	Subject       *string   `json:"subject,omitempty"`
	Tags          *[]string `json:"tags,omitempty"`
	Type          *int      `json:"type,omitempty"`
	Version       int       `json:"version"`
	Watchers      *[]int    `json:"watchers,omitempty"`
	DueDate       *string   `json:"due_date,omitempty"`
	DueDateReason *string   `json:"due_date_reason,omitempty"`
	DueDateStatus *string   `json:"due_date_status,omitempty"`
}

// List => https://taigaio.github.io/taiga-doc/dist/api.html#issues-list
func (s *IssueService) List(queryParams *IssueQueryParams) ([]Issue, error) {
	url := s.client.MakeURL(s.Endpoint)
	url, err := urlWithQueryOrDefaultProject(url, queryParams, s.defaultProjectID)
	if err != nil {
		return nil, err
	}

	// execute requests
	var issues IssueDetailLIST
	_, err = s.client.Request.Get(url, &issues)
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
	if err := requireNonNil("issue", issue); err != nil {
		return nil, err
	}
	i := Issue{}
	err := convertStructViaJSON(issue, &i)
	if err != nil {
		return nil, err
	}
	if i.ID == 0 || i.Project == 0 {
		return nil, errors.New("issue id and project are required to list attachments")
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
	if err := requirePositiveID("issueID", issueID); err != nil {
		return nil, err
	}
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
	if err := requirePositiveID("issueRef", issueRef); err != nil {
		return nil, err
	}
	var issue IssueDetailGET
	var url string

	type byRefQueryParams struct {
		Ref         int    `url:"ref"`
		Project     int    `url:"project,omitempty"`
		ProjectSlug string `url:"project__slug,omitempty"`
	}
	queryParams := &byRefQueryParams{Ref: issueRef}
	switch {
	case project != nil && project.ID != 0:
		queryParams.Project = project.ID
	case project != nil && len(project.Slug) > 0:
		queryParams.ProjectSlug = project.Slug
	case s.defaultProjectID > 0:
		queryParams.Project = s.defaultProjectID
	default:
		return nil, errors.New("no project ID/slug provided and no mapped default project ID set")
	}
	url, err := appendQueryParams(s.client.MakeURL(s.Endpoint, "by_ref"), queryParams)
	if err != nil {
		return nil, err
	}

	_, err = s.client.Request.Get(url, &issue)
	if err != nil {
		return nil, err
	}
	return issue.AsIssue()
}

// Edit sends a PATCH request to edit a Issue -> https://taigaio.github.io/taiga-doc/dist/api.html#issues-edit
// Available Meta: IssueDetail
func (s *IssueService) Edit(issue *Issue) (*Issue, error) {
	if err := requireNonNil("issue", issue); err != nil {
		return nil, err
	}

	if issue.ID == 0 {
		return nil, errors.New("passed Issue does not have an ID yet. Does it exist?")
	}
	if issue.Version == 0 {
		return nil, errors.New("version is required for issue edit")
	}

	patchPayload := map[string]any{
		"version": issue.Version,
	}
	if issue.AssignedTo != 0 {
		patchPayload["assigned_to"] = issue.AssignedTo
	}
	if issue.BlockedNote != "" {
		patchPayload["blocked_note"] = issue.BlockedNote
	}
	if issue.Description != "" {
		patchPayload["description"] = issue.Description
	}
	if issue.IsBlocked {
		patchPayload["is_blocked"] = issue.IsBlocked
	}
	if issue.Milestone != 0 {
		patchPayload["milestone"] = issue.Milestone
	}
	if issue.Owner != 0 {
		patchPayload["owner"] = issue.Owner
	}
	if issue.Priority != 0 {
		patchPayload["priority"] = issue.Priority
	}
	if issue.Project != 0 {
		patchPayload["project"] = issue.Project
	}
	if issue.Severity != 0 {
		patchPayload["severity"] = issue.Severity
	}
	if issue.Status != 0 {
		patchPayload["status"] = issue.Status
	}
	if issue.Subject != "" {
		patchPayload["subject"] = issue.Subject
	}
	if issue.Type != 0 {
		patchPayload["type"] = issue.Type
	}
	if issue.DueDate != "" {
		patchPayload["due_date"] = issue.DueDate
	}
	if issue.DueDateReason != "" {
		patchPayload["due_date_reason"] = issue.DueDateReason
	}
	if issue.DueDateStatus != "" {
		patchPayload["due_date_status"] = issue.DueDateStatus
	}
	if issue.Tags != nil {
		tags := tagsToNames(issue.Tags)
		if tags == nil {
			tags = []string{}
		}
		patchPayload["tags"] = tags
	}
	if issue.Watchers != nil {
		watchers := append([]int(nil), issue.Watchers...)
		patchPayload["watchers"] = watchers
	}
	if len(patchPayload) == 1 {
		return nil, errors.New("no updatable issue fields were provided; use Patch for explicit zero-value updates")
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(issue.ID))
	var responseIssue IssueDetail
	_, err := s.client.Request.Patch(url, &patchPayload, &responseIssue)
	if err != nil {
		return nil, err
	}
	return responseIssue.AsIssue()
}

// Patch sends an explicit PATCH payload to edit an issue.
func (s *IssueService) Patch(issueID int, patch *IssuePatch) (*Issue, error) {
	if err := requireNonNil("patch", patch); err != nil {
		return nil, err
	}
	if err := requirePositiveID("issueID", issueID); err != nil {
		return nil, err
	}
	if patch.Version == 0 {
		return nil, errors.New("version is required for issue patch")
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(issueID))
	var responseIssue IssueDetail
	_, err := s.client.Request.Patch(url, patch, &responseIssue)
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
	if err := requireNonNil("issue", issue); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint)
	var issueDetail IssueDetail
	projectID, err := resolveProjectID(issue.Project, s.defaultProjectID, "project")
	if err != nil {
		return nil, err
	}

	// Check for required fields
	// project, subject
	if isEmpty(issue.Subject) {
		return nil, errors.New("a mandatory field is missing. See API documentataion")
	}

	payload := issueCreatePayload{
		AssignedTo:    issue.AssignedTo,
		BlockedNote:   issue.BlockedNote,
		Description:   issue.Description,
		IsBlocked:     issue.IsBlocked,
		Milestone:     issue.Milestone,
		Owner:         issue.Owner,
		Priority:      issue.Priority,
		Project:       projectID,
		Severity:      issue.Severity,
		Status:        issue.Status,
		Subject:       issue.Subject,
		Tags:          tagsToNames(issue.Tags),
		Type:          issue.Type,
		Watchers:      issue.Watchers,
		DueDate:       issue.DueDate,
		DueDateReason: issue.DueDateReason,
		DueDateStatus: issue.DueDateStatus,
	}

	_, err = s.client.Request.Post(url, &payload, &issueDetail)
	if err != nil {
		return nil, err
	}

	return issueDetail.AsIssue()
}

// Delete -> https://taigaio.github.io/taiga-doc/dist/api.html#issues-delete
func (s *IssueService) Delete(issueID int) (*http.Response, error) {
	if err := requirePositiveID("issueID", issueID); err != nil {
		return nil, err
	}
	url := s.client.MakeURL(s.Endpoint, strconv.Itoa(issueID))
	return s.client.Request.Delete(url)
}
