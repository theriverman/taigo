package taigo

import (
	"github.com/google/go-querystring/query"
)

var resolverURI = "/resolver"

/*
NOTES ON LOGIC:
	Received structs are not passed directly to `genericResolver` but instead they're manually assigned
	before sending to make sure there are no extra query parameters sent/picked up when only one is expected.
*/

// ResolveProject => https://taigaio.github.io/taiga-doc/dist/api.html#resolver-projects
func ResolveProject(c *Client, queryParameters *ResolverQueryParams) (*Resolver, error) {
	qp := ResolverQueryParams{
		Project: queryParameters.Project,
	}
	return genericResolver(c, &qp)
}

// ResolveUserStory => https://taigaio.github.io/taiga-doc/dist/api.html#resolver-user-stories
func ResolveUserStory(c *Client, queryParameters *ResolverQueryParams) (*Resolver, error) {
	qp := ResolverQueryParams{
		US:      queryParameters.US,
		Project: queryParameters.Project,
	}
	return genericResolver(c, &qp)
}

// ResolveIssue => https://taigaio.github.io/taiga-doc/dist/api.html#resolver-issues
func ResolveIssue(c *Client, queryParameters *ResolverQueryParams) (*Resolver, error) {
	qp := ResolverQueryParams{
		Issue:   queryParameters.Issue,
		Project: queryParameters.Project,
	}
	return genericResolver(c, &qp)
}

// ResolveTask => https://taigaio.github.io/taiga-doc/dist/api.html#resolver-tasks
func ResolveTask(c *Client, queryParameters *ResolverQueryParams) (*Resolver, error) {
	qp := ResolverQueryParams{
		Task:    queryParameters.Task,
		Project: queryParameters.Project,
	}
	return genericResolver(c, &qp)
}

// ResolveMilestone => https://taigaio.github.io/taiga-doc/dist/api.html#resolver-milestones
func ResolveMilestone(c *Client, queryParameters *ResolverQueryParams) (*Resolver, error) {
	qp := ResolverQueryParams{
		Milestone: queryParameters.Milestone,
		Project:   queryParameters.Project,
	}
	return genericResolver(c, &qp)
}

// ResolveWikiPage => https://taigaio.github.io/taiga-doc/dist/api.html#resolver-wiki-pages
func ResolveWikiPage(c *Client, queryParameters *ResolverQueryParams) (*Resolver, error) {
	qp := ResolverQueryParams{
		WikiPage: queryParameters.WikiPage,
		Project:  queryParameters.Project,
	}
	return genericResolver(c, &qp)
}

// ResolveMultipleObjects => https://taigaio.github.io/taiga-doc/dist/api.html#resolver-multiple-resolution
func ResolveMultipleObjects(c *Client, queryParameters *ResolverQueryParams) (*Resolver, error) {
	return genericResolver(c, queryParameters)
}

// ResolveByRefValue => https://taigaio.github.io/taiga-doc/dist/api.html#resolver-ref
// TODO: Not yet implemented. Considered for later.

// genericResolver acts as a common request execution middleware
func genericResolver(c *Client, queryParameters *ResolverQueryParams) (*Resolver, error) {
	paramValues, _ := query.Values(queryParameters)
	url := c.APIURL + resolverURI + "?" + paramValues.Encode()
	var respResolver Resolver
	err := c.Request.Get(url, &respResolver)
	if err != nil {
		return nil, err
	}
	return &respResolver, nil
}
