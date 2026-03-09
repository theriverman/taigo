# Taigo Project Agent Notes

## 1) What this repo is
- Go client library for Taiga REST API v1 with **v2 module versioning** (`module github.com/theriverman/taigo/v2`).
- Main package is a library (`package taigo`) with service-based API wrappers.
- Two extra modules:
  - `cli/`: standalone CLI utility.
  - `contribute/`: integration/demo runner against a real Taiga instance.

## 2) Repository layout
- Root library:
  - `client.go`: client lifecycle, headers, auth state, service wiring, token refresh routine.
  - `requests.go`: HTTP transport (`RequestService`) and multipart attachment uploads.
  - Core resource services: `projects.go`, `epics.go`, `user_stories.go`, `tasks.go`, `issues.go`, `milestones.go`, `users.go`, `webhooks.go`, `wiki.go`, `stats.go`, `resolver.go`.
  - Classification and custom-attribute services: `points.go`, `priorities.go`, `severities.go`, `issue_types.go`, `*_status.go`, `*_custom_attribute.go`.
  - Extended surface services: `applications.go`, `application_tokens.go`, `searches.go`, `user_storage.go`, `project_templates.go`, `project_templates_detail.go`, `memberships_invitations.go`, `wiki_links.go`, `history.go`, `notify_policies.go`, `contact.go`, `feedback.go`, `export_import.go`, `timelines.go`, `locales.go`, `importers.go`, `contrib_plugins.go`, `objects_summary.go`.
  - `raw_resource.go`: generic helpers for endpoints that are still represented by raw map-based DTOs.
  - `*.models.go`: DTO/query structs/conversion helpers.

## 3) Architecture and conventions
- Entry point is `taigo.Client`.
- `Client.Initialise()`:
  - Sets headers, base API URL (`/api/v1` by default), disables pagination by default.
  - Instantiates all services.
  - Starts optional token-refresh ticker goroutine.
- Services:
  - Thin wrappers over `RequestService` (`GET/POST/PUT/PATCH/DELETE`) and `MakeURL`.
  - Optional `defaultProjectID` for project-scoped mapped services (`Project.ConfigureMappedServices`).
- Models:
  - Strongly typed models are used for core resources.
  - Raw map DTOs (`RawResource`) are used for less stable or less modelled endpoint groups.

## 4) Testing reality
- Root has unit tests (`v2_proposals_test.go`) for transport/query/header behaviour.
- Root has matrix-style offline suites:
  - `contract_matrix_test.go` (method/path/query/body contracts).
  - `query_filter_matrix_test.go` (query encoding semantics).
  - `negative_matrix_test.go` (API error and validation guard behaviour).
- `tests/` contains integration tests for core resources.
- `tests/smoke_matrix_test.go` is a table-driven real-instance harness for CRUD smoke coverage.
- `tests/workflow_matrix_test.go` covers cross-resource end-to-end lifecycle flows.
- `tests/negative_matrix_test.go` covers live negative paths.
- `tests/auth_role_matrix_test.go` covers auth and optional role expectations.
- CI gates in `.github/workflows/go.yml`:
  - PR/push: unit+contract+query+negative offline checks.
  - PR/push: live smoke subset on Docker Taiga.
  - Nightly/manual: full live test suite.
  - Tag `v*`: pre-release full live suite across Taiga refs (`master`, `v2`).
- Integration suite is opt-in and skipped by default unless:
  - `TAIGO_RUN_INTEGRATION_TESTS=1`
- Default integration target in tests: `http://localhost:9000` (override via `TAIGO_BASE_URL` and related env vars).

## 5) Practical commands
- List all Go files quickly:
  - `rg --files -g '*.go'`
- Run root tests:
  - `GOWORK=off GOCACHE=/tmp/taigo-gocache go test ./...`
- Run CLI module:
  - `cd cli && GOWORK=off GOCACHE=/tmp/taigo-gocache-cli go test ./...`
- Run contribute module:
  - `cd contribute && GOWORK=off GOCACHE=/tmp/taigo-gocache-contrib go test ./...`
- Run integration tests explicitly:
  - `TAIGO_RUN_INTEGRATION_TESTS=1 GOWORK=off GOCACHE=/tmp/taigo-gocache go test ./tests/...`
  - `TAIGO_RUN_INTEGRATION_TESTS=1 TAIGO_PROJECT_ID=2 go test ./tests/... -run TestSmokeCRUDMatrix -v`
  - `TAIGO_RUN_INTEGRATION_TESTS=1 TAIGO_PROJECT_ID=2 go test ./tests/... -run TestWorkflowMatrixLive -v`
  - `TAIGO_RUN_INTEGRATION_TESTS=1 TAIGO_PROJECT_ID=2 go test ./tests/... -run TestNegativeMatrixLive -v`
  - `TAIGO_RUN_INTEGRATION_TESTS=1 TAIGO_MEMBER_USERNAME=<u> TAIGO_MEMBER_PASSWORD=<p> TAIGO_MEMBER_WRITE_EXPECTATION=forbid go test ./tests/... -run TestRoleMatrixLive -v`

## 6) External references
- Taiga API docs: [https://docs.taiga.io/api.html](https://docs.taiga.io/api.html)
- Taiga backend source: [https://github.com/taigaio/taiga-back](https://github.com/taigaio/taiga-back)

## 7) Self-check checklist before future changes
- Does endpoint path/method match Taiga docs and backend route naming?
- Do request/response DTO field names match serializer fields?
- Are query params encoded exactly as Taiga expects (including comma-joined filters and pointer-bool semantics)?
- Are method signatures consistent across resource services (`Create/Get/GetByRef/Edit|Update/Delete/List` where relevant)?
- Are unit tests added for path/query/transport behaviour, and integration tests updated when behaviour changes?
- Is the code written using the latest Go version's capabilities?
