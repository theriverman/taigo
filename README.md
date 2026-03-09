# Taigo

[![Go](https://github.com/theriverman/taigo/actions/workflows/go.yml/badge.svg)](https://github.com/theriverman/taigo/actions/workflows/go.yml)
[![GoDoc](https://godoc.org/github.com/theriverman/taigo/v2?status.svg)](https://pkg.go.dev/github.com/theriverman/taigo/v2?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/theriverman/taigo)](https://goreportcard.com/report/github.com/theriverman/taigo)

Taigo is a Go client library for the [Taiga](https://github.com/taigaio) REST API v1.

## Status

This release line contains the v2 API with breaking changes and expanded endpoint coverage.
See [MIGRATION.md](./MIGRATION.md) for upgrade guidance.

## Install

```bash
go get github.com/theriverman/taigo/v2
```

## Quick Start

```go
package main

import (
	"fmt"
	"net/http"

	taiga "github.com/theriverman/taigo/v2"
)

func main() {
	client := taiga.Client{
		BaseURL:    "https://api.taiga.io",
		HTTPClient: &http.Client{},
	}

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

## Client Behaviour

- API root defaults to `/api/v1`.
- Pagination is disabled by default by sending `x-disable-pagination`.
- Calling `DisablePagination(false)` removes that header (Taiga checks presence, not value).
- Request errors return a typed `*taigo.APIError` with status code and response body.

## Services

Core services available via `Client`:

- `Auth`, `Project`, `User`, `Resolver`, `Stats`
- `Epic`, `UserStory`, `Task`, `Issue`, `Milestone`, `Wiki`
- `Webhook`
- `Point`, `Priority`, `Severity`, `IssueType`
- `EpicStatus`, `IssueStatus`, `TaskStatus`, `UserStoryStatus`
- `EpicCustomAttribute`, `IssueCustomAttribute`, `TaskCustomAttribute`, `UserStoryCustomAttribute`
- `Application`, `ApplicationToken`, `Search`, `UserStorage`
- `ProjectTemplate`, `ProjectTemplateDetail`
- `MembershipInvitation`, `WikiLink`, `History`, `NotifyPolicy`
- `Contact`, `Feedback`, `ExportImport`, `Timeline`, `Locale`, `Importer`
- `ContribPlugin`, `ObjectsSummary`

Project-scoped mapped services are available via:

```go
client.Project.ConfigureMappedServices(projectID)
```

## Notable v2 API Changes

- `TaskService.Get` now takes `taskID int`.
- `TaskService.GetByRef` now takes `taskRef int, project *Project`.
- `IssueService` now includes `GetByRef`, `Delete`, `GetAttachment`, `ListAttachments`.
- `WikiService` now includes `List`, `Create`, `Get`, `GetBySlug`, `Edit`, `Delete`, `Render`.
- `WebhookService.GetWebhookLog` now takes `webhookLogID int`.
- `UserService.GetWatchedContent` and `GetLikedContent` return slices.

Many services now expose `Update(...)` as an alias for `Edit(...)`.

## Query Parameters

Query models were tightened to match Taiga semantics:

- Optional booleans use pointer-bools (`*bool`) for tri-state filtering.
- `TasksQueryParams.Tags` is a comma-delimited string; use `SetTags(...)` helper.
- Project list ordering now serialises correctly via exported `OrderBy`.

## Testing

Run unit tests and offline-safe checks:

```bash
GOWORK=off GOCACHE=/tmp/taigo-gocache go test ./...
```

Integration tests under `tests/` require a reachable Taiga instance and are opt-in:

```bash
TAIGO_RUN_INTEGRATION_TESTS=1 GOWORK=off GOCACHE=/tmp/taigo-gocache go test ./tests/...
```

The integration harness supports environment overrides:

- `TAIGO_BASE_URL` (default: `http://localhost:9000`)
- `TAIGO_USERNAME` (default: `admin`)
- `TAIGO_PASSWORD` (default: `admin`)
- `TAIGO_PROJECT_ID` (default: `2`)
- `TAIGO_PROJECT_SLUG` (default: `taigo-test`)
- `TAIGO_USER_ID` (default: `5`)

Run only the table-driven smoke matrix harness:

```bash
TAIGO_RUN_INTEGRATION_TESTS=1 TAIGO_PROJECT_ID=2 go test ./tests/... -run TestSmokeCRUDMatrix -v
```

## Related Modules

- `cli/`: CLI utility and auth/config workflow examples.
- `contribute/`: small executable module for live-instance contribution tests.

## References

- Taiga API docs: <https://docs.taiga.io/api.html>
- Taiga backend source: <https://github.com/taigaio/taiga-back>

## Contributing

- Please open issues or pull requests for bugs and endpoint gaps.
- Keep endpoint paths, query encoding, and model fields aligned with Taiga docs and backend serializers.
