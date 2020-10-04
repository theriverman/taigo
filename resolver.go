package taigo

import (
	"fmt"

	"github.com/google/go-querystring/query"
)

/*
NOTES ON LOGIC:
	Received structs are not passed directly to `genericResolver` but instead they're manually assigned
	before sending to make sure there are no extra query parameters sent/picked up when only one is expected.
*/

// ResolverService is a handle to actions related to Resolver
//
// https://taigaio.github.io/taiga-doc/dist/api.html#resolver
type ResolverService struct {
	client           *Client
	defaultProjectID int
	Endpoint         string
}

// ResolveProject => https://taigaio.github.io/taiga-doc/dist/api.html#resolver-projects
func (s *ResolverService) ResolveProject(queryParameters *ResolverQueryParams) (*Resolver, error) {
	qp := ResolverQueryParams{
		Project: queryParameters.Project,
	}
	return s.genericResolver(&qp)
}

// ResolveUserStory => https://taigaio.github.io/taiga-doc/dist/api.html#resolver-user-stories
func (s *ResolverService) ResolveUserStory(queryParameters *ResolverQueryParams) (*Resolver, error) {
	qp := ResolverQueryParams{
		US:      queryParameters.US,
		Project: queryParameters.Project,
	}
	return s.genericResolver(&qp)
}

// ResolveIssue => https://taigaio.github.io/taiga-doc/dist/api.html#resolver-issues
func (s *ResolverService) ResolveIssue(queryParameters *ResolverQueryParams) (*Resolver, error) {
	qp := ResolverQueryParams{
		Issue:   queryParameters.Issue,
		Project: queryParameters.Project,
	}
	return s.genericResolver(&qp)
}

// ResolveTask => https://taigaio.github.io/taiga-doc/dist/api.html#resolver-tasks
func (s *ResolverService) ResolveTask(queryParameters *ResolverQueryParams) (*Resolver, error) {
	qp := ResolverQueryParams{
		Task:    queryParameters.Task,
		Project: queryParameters.Project,
	}
	return s.genericResolver(&qp)
}

// ResolveMilestone => https://taigaio.github.io/taiga-doc/dist/api.html#resolver-milestones
func (s *ResolverService) ResolveMilestone(queryParameters *ResolverQueryParams) (*Resolver, error) {
	qp := ResolverQueryParams{
		Milestone: queryParameters.Milestone,
		Project:   queryParameters.Project,
	}
	return s.genericResolver(&qp)
}

// ResolveWikiPage => https://taigaio.github.io/taiga-doc/dist/api.html#resolver-wiki-pages
func (s *ResolverService) ResolveWikiPage(queryParameters *ResolverQueryParams) (*Resolver, error) {
	qp := ResolverQueryParams{
		WikiPage: queryParameters.WikiPage,
		Project:  queryParameters.Project,
	}
	return s.genericResolver(&qp)
}

// ResolveMultipleObjects => https://taigaio.github.io/taiga-doc/dist/api.html#resolver-multiple-resolution
func (s *ResolverService) ResolveMultipleObjects(queryParameters *ResolverQueryParams) (*Resolver, error) {
	return s.genericResolver(queryParameters)
}

// ResolveByRefValue => https://taigaio.github.io/taiga-doc/dist/api.html#resolver-ref
// TODO: Not yet implemented. Considered for later.

// genericResolver acts as a common request execution middleware
func (s *ResolverService) genericResolver(queryParameters *ResolverQueryParams) (*Resolver, error) {
	paramValues, _ := query.Values(queryParameters)
	url := s.client.MakeURL(fmt.Sprintf("%s?%s", s.Endpoint, paramValues.Encode()))
	var respResolver Resolver
	_, err := s.client.Request.Get(url, &respResolver)
	if err != nil {
		return nil, err
	}
	return &respResolver, nil
}
