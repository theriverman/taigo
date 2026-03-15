# Contributing to Taigo

Taigo is a Go client library for the Taiga REST API v1, published as the `v2` Go module path:

```text
github.com/theriverman/taigo/v2
```

This document is for contributors and maintainers.
User-facing installation and quick-start guidance lives in [`README.md`](./README.md).
Breaking-change upgrade notes live in [`MIGRATION.md`](./MIGRATION.md).
Runnable usage snippets live in [`examples/README.md`](./examples/README.md).

## Sources of truth

When making changes, always verify behavior against the real Taiga contract:

- Repository code and tests
- Taiga API docs: <https://docs.taiga.io/api.html>
- Taiga backend source: <https://github.com/taigaio/taiga-back>

If the code, docs, and backend disagree, prefer the backend behavior and update the client and tests accordingly.

## Project layout

- Root module: Go client library in `package taigo`
- `tests/`: opt-in integration test suite for a real Taiga instance
- `examples/README.md`: usage snippets
- `.github/workflows/go.yml`: CI definition
- `.github/workflows/codeql.yml`: CodeQL security analysis
- `.github/dependabot.yml`: dependency update automation

Core implementation areas:

- `client.go`: client lifecycle, auth state, headers, service wiring, token refresh
- `requests.go`: HTTP transport and multipart attachment upload
- `projects.go`, `epics.go`, `user_stories.go`, `tasks.go`, `issues.go`, `milestones.go`, `users.go`, `webhooks.go`, `wiki.go`: primary resource services
- `points.go`, `priorities.go`, `severities.go`, `issue_types.go`, `*_status.go`, `*_custom_attribute.go`: taxonomy/status/custom-attribute services
- `raw_resource.go`: generic wrappers for endpoints that are not yet modeled with dedicated DTOs
- `*.models.go`: response DTOs, query types, and conversion helpers

## Tooling and environment

This is a Go project.
There is no Node.js, no npm, and no JavaScript build pipeline.
Do not add frontend toolchain dependencies to contributor workflows or CI unless the project direction explicitly changes.

Use the Go toolchain declared in `go.mod`.
At the time of writing, the repository targets Go `1.25`.

## Development principles

### Keep the client contract-correct

- Endpoint paths, HTTP methods, query parameters, and JSON field names must match Taiga.
- Do not add convenience behavior that silently changes Taiga semantics.
- If a behavior is intentionally opinionated, document it clearly.

### Separate write DTOs from response DTOs

For write-oriented resource families that use dedicated request models:

- `Create` accepts a `*...CreateRequest`
- `Edit` accepts `(resourceID int, *...EditRequest)` and performs sparse, non-destructive updates
- `Patch` accepts `(resourceID int, *...Patch)` for explicit zero-value updates

Do not reuse response DTOs as write payloads for these families.

### Preserve update semantics

- `Edit` means safe partial update
- `Patch` means explicit control, including `0`, `false`, empty string, or field clearing when the backend supports it

Do not blur these two behaviors.

### Respect mapped project services

`client.Project.ConfigureMappedServices(projectID)` is the canonical way to bind project-scoped services.
Mapped services must:

- use the configured project when the caller omits it
- not override explicit project values from the caller
- behave consistently across all resource families

### Validate early

Public methods should reject invalid IDs, missing required inputs, and obviously unusable arguments before performing network calls.

### Keep background behavior quiet

Do not emit internal background logging unless it is gated behind `Client.Verbose`.

## Formatting and code style

- Run `gofmt` on every touched Go file
- Keep imports standard and minimal
- Prefer small, explicit helpers over hidden magic
- Prefer typed DTOs over `map[string]any` when the endpoint contract is stable enough to model
- Use `RawResource` only where the API surface is still intentionally generic

## Recommended commands

Run these from the repository root before opening a pull request:

```bash
GOWORK=off go test ./...
GOWORK=off go vet ./...
GOWORK=off go test -race ./...
```

If you want an isolated cache, set `GOCACHE` explicitly:

```bash
GOWORK=off GOCACHE=/tmp/taigo-gocache go test ./...
```

Useful local commands:

```bash
rg --files -g '*.go'
gofmt -w path/to/file.go
```

## Integration tests

The live integration suite is opt-in and skipped by default.
It expects a reachable Taiga instance and valid credentials.

Run the full live suite:

```bash
TAIGO_RUN_INTEGRATION_TESTS=1 GOWORK=off go test ./tests/...
```

Useful focused runs:

```bash
TAIGO_RUN_INTEGRATION_TESTS=1 go test ./tests/... -run TestSmokeCRUDMatrix -v
TAIGO_RUN_INTEGRATION_TESTS=1 go test ./tests/... -run TestWorkflowMatrixLive -v
TAIGO_RUN_INTEGRATION_TESTS=1 go test ./tests/... -run TestNegativeMatrixLive -v
```

Environment variables used by the live suite include:

- `TAIGO_BASE_URL`
- `TAIGO_USERNAME`
- `TAIGO_PASSWORD`
- `TAIGO_PROJECT_ID`
- `TAIGO_PROJECT_SLUG`
- `TAIGO_USER_ID`
- `TAIGO_MEMBER_USERNAME`
- `TAIGO_MEMBER_PASSWORD`
- `TAIGO_MEMBER_WRITE_EXPECTATION`

## Testing expectations for changes

If you change endpoint behavior, signatures, DTO fields, or query encoding:

- add or update offline unit/contract tests in the root module
- update live tests in `tests/` when the change affects real-instance behavior
- update examples and docs when the public API changed

At minimum:

- transport and path/query/body contracts belong in `contract_matrix_test.go` or `v2_proposals_test.go`
- validation behavior belongs in `negative_matrix_test.go`
- live CRUD behavior belongs in `tests/`

## Documentation expectations

Keep documentation split by audience:

- `README.md`: user-facing overview, install, quick start, high-level behavior
- `MIGRATION.md`: breaking changes and upgrade guidance
- `examples/README.md`: usage snippets
- `CONTRIBUTION.md`: contributor workflow, design rules, testing expectations

When the public API changes, update the relevant documents in the same pull request.

## Pull request checklist

Before opening a PR:

1. Branch from `master`.
2. Keep the change scoped and coherent.
3. Run `gofmt`, tests, vet, and race checks.
4. Update docs for any public API change.
5. Add tests for behavior changes.
6. Call out any breaking changes clearly.
7. Reference the Taiga docs or backend route/serializer when fixing or adding endpoints.

## Reporting bugs and proposing changes

Good issues include:

- Taigo version or commit
- Go version
- Taiga version or deployment details if known
- Minimal reproduction code
- Expected behavior
- Actual behavior
- Relevant request/response payloads or status codes

For endpoint mismatches, include the Taiga doc link and, if possible, the corresponding `taiga-back` implementation link.

## Licensing

By contributing, you agree that your contributions are licensed under the repository's MIT license.
