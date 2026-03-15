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

## Coverage

The library covers Taiga authentication, projects, epics, user stories, tasks, issues, milestones, wiki pages, webhooks, users, and related taxonomy/status/custom-attribute resources.

Project-scoped mapped services are available via `client.Project.ConfigureMappedServices(projectID)`.

## Documentation

- [MIGRATION.md](./MIGRATION.md): breaking changes and upgrade notes for `v2`
- [examples/README.md](./examples/README.md): runnable usage snippets and patterns
- [CONTRIBUTION.md](./CONTRIBUTION.md): repository layout, contributor workflow, design rules, and test expectations

## Contributing

Contributor and maintainer guidance lives in [CONTRIBUTION.md](./CONTRIBUTION.md).
