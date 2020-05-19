package taigo

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

const endpointIssues = "/issues"

// IssueService is a handle to actions related to Issues
//
// https://taigaio.github.io/taiga-doc/dist/api.html#issues
type IssueService struct {
	client *Client
}

// List => https://taigaio.github.io/taiga-doc/dist/api.html#issues-list
func (s *IssueService) List(queryParams *IssueQueryParams) ([]Issue, error) {
	// prepare url & parameters
	url := s.client.APIURL + endpointIssues
	if queryParams != nil {
		paramValues, _ := query.Values(queryParams)
		url = fmt.Sprintf("%s?%s", url, paramValues.Encode())
	}
	// execute requests
	var issues IssueDetailLIST
	err := s.client.Request.Get(url, &issues)
	if err != nil {
		return nil, err
	}

	return issues.AsIssues()
}

// CreateAttachment creates a new Issue attachment => https://taigaio.github.io/taiga-doc/dist/api.html#issues-create-attachment
func (s *IssueService) CreateAttachment(attachment *Attachment, task *Task) (*Attachment, error) {
	url := s.client.APIURL + endpointTasks + "/attachments"
	return newfileUploadRequest(s.client, url, attachment, task)
}
