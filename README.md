# Taigo

[![CI](https://github.com/theriverman/taigo/actions/workflows/go.yml/badge.svg)](https://github.com/theriverman/taigo/actions/workflows/go.yml)
[![CodeQL](https://github.com/theriverman/taigo/actions/workflows/codeql.yml/badge.svg)](https://github.com/theriverman/taigo/actions/workflows/codeql.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/theriverman/taigo/v2.svg)](https://pkg.go.dev/github.com/theriverman/taigo/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/theriverman/taigo)](https://goreportcard.com/report/github.com/theriverman/taigo)

Taigo is a Go client library for the [Taiga](https://github.com/taigaio) REST API v1.

The current stable major version is `v2`, published at:

```text
github.com/theriverman/taigo/v2
```

## Requirements

- Go 1.25 or newer
- A reachable Taiga instance

## Install

```bash
go get github.com/theriverman/taigo/v2@latest
```

Pin the v2.0.0 release explicitly when you need reproducible installs:

```bash
go get github.com/theriverman/taigo/v2@v2.0.0
```

## Quick Start

```go
package main

import (
	"fmt"

	taiga "github.com/theriverman/taigo/v2"
)

func main() {
	client := taiga.Client{
		BaseURL: "https://api.taiga.io",
	}
	defer client.Close()

	if err := client.AuthByCredentials(&taiga.Credentials{
		Type:     "normal",
		Username: "your-user",
		Password: "your-password",
	}); err != nil {
		panic(err)
	}

	me, err := client.User.Me()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Logged in as %s (id=%d)\n", me.Username, me.ID)
}
```

`AuthByCredentials` and `AuthByToken` initialise the client automatically. If you only need public or manually authenticated calls, call `Initialise` yourself after setting `BaseURL`.

## Authentication

Use credentials when you want Taigo to issue bearer and refresh tokens:

```go
err := client.AuthByCredentials(&taiga.Credentials{
	Type:     "normal",
	Username: "your-user",
	Password: "your-password",
})
```

Use existing tokens when your application stores them elsewhere:

```go
err := client.AuthByToken(taiga.TokenBearer, authToken, refreshToken)
```

By default, Taigo refreshes stored tokens every 12 hours after authentication. Call `DisableAutomaticTokenRefresh` when your application owns token refresh, and call `Close` when a long-lived client is no longer needed.

## Common Usage

### Projects

```go
project, err := client.Project.GetBySlug("my-project")
if err != nil {
	panic(err)
}

client.Project.ConfigureMappedServices(project.ID)
```

After configuring mapped services, project-scoped calls can omit the project ID when the query or payload supports a project default:

```go
stories, err := client.Project.UserStory.List(&taiga.UserStoryQueryParams{
	StatusIsClosed: taiga.BoolPtr(false),
})
```

Explicit project values still win over the mapped default.

### Query Filters

List methods accept typed query structs. Optional boolean filters use pointers so both `true` and `false` can be sent to Taiga:

```go
tasks, err := client.Task.List(&taiga.TasksQueryParams{
	Project:            project.ID,
	StatusIsClosed:     taiga.BoolPtr(false),
	IncludeAttachments: taiga.BoolPtr(true),
})
```

Task tag filters use Taiga's comma-separated tag format:

```go
query := taiga.TasksQueryParams{Project: project.ID}
query.SetTags("backend", "api")

tasks, err := client.Task.List(&query)
```

### Create and Edit

Core resources such as projects, epics, user stories, tasks, issues, milestones, wiki pages, and webhooks use their resource models for create/edit operations.

Classification, status, and custom-attribute services use write-specific DTOs:

```go
priority, err := client.Project.Priority.Create(&taiga.PriorityCreateRequest{
	Name:  "High",
	Color: "#d9534f",
})
if err != nil {
	panic(err)
}

priority, err = client.Project.Priority.Edit(priority.ID, &taiga.PriorityEditRequest{
	Name: "Very High",
})
```

Use `Patch` with pointer fields when you need to send explicit zero-values such as `false`, `0`, or an empty string:

```go
name := ""
priority, err = client.Project.Priority.Patch(priority.ID, &taiga.PriorityPatch{
	Name: &name,
})
```

### Attachments

Attachments are supported for epics, issues, tasks, user stories, and wiki pages:

```go
attachment := &taiga.Attachment{Description: "Build log"}
attachment.SetFilePath("/path/to/log.txt")

created, err := client.Task.CreateAttachment(attachment, task)
```

### Raw Endpoints

Some less stable Taiga endpoints are exposed through `taiga.RawResource`, a `map[string]any` JSON object. This keeps the endpoint reachable while avoiding premature DTOs for surfaces that vary between deployments:

```go
results, err := client.Search.Search(&taiga.SearchQueryParams{
	Project: project.ID,
	Text:    "login",
})
```

## Client Behaviour

- `BaseURL` is the Taiga base URL, for example `https://api.taiga.io`; the client appends `/api/v1`.
- If `HTTPClient` is nil, the client uses a default `http.Client` with a 30-second timeout.
- Pagination is disabled by default by sending the `x-disable-pagination` header.
- `DisablePagination(false)` removes that header because Taiga checks header presence, not header value.
- `GetPagination` returns pagination details captured from the last HTTP response.
- Non-2xx responses return `*taiga.APIError` with the Taiga status code and response body.
- `RequestService` exposes lower-level HTTP helpers, including context-aware variants such as `GetCtx`, `PostCtx`, `PatchCtx`, and `DeleteCtx`.

Example API error handling with the standard `errors` package:

```go
project, err := client.Project.Get(123)
if err != nil {
	var apiErr *taiga.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("Taiga returned %d: %s\n", apiErr.StatusCode, apiErr.Body)
		return
	}
	panic(err)
}

_ = project
```

## Endpoint Coverage

Taigo includes typed services for the main Taiga resources:

- Authentication and users
- Projects, epics, user stories, tasks, issues, and milestones
- Wiki pages, wiki links, webhooks, resolver, search, stats, and object summaries
- Points, priorities, severities, issue types, and status services
- Epic, issue, task, and user-story custom attributes

It also exposes service coverage for application tokens, user storage, project templates, memberships and invitations, history, notification policies, contact and feedback, export/import, timelines, locales, importers, and contributed plugins. Where these surfaces are represented as `RawResource`, callers can still use typed request methods and decode the returned map into application-specific structs if needed.

## More Documentation

- [Package documentation](https://pkg.go.dev/github.com/theriverman/taigo/v2)
- [Examples](./examples/README.md)
- [v1 to v2 migration guide](./MIGRATION.md)
- [Changelog](./CHANGELOG.md)
- [Contribution guide](./CONTRIBUTION.md)
