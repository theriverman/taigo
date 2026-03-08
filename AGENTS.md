# Taigo Project Agent Notes

## 1) What this repo is
- Go client library for Taiga REST API v1 (`module github.com/theriverman/taigo`).
- Main package is a library (`package taigo`) with service-based API wrappers.
- Two extra modules:
  - `cli/`: standalone CLI demo/auth utility.
  - `contribute/`: integration/demo runner against a real Taiga instance.

## 2) Repository layout
- Root library:
  - `client.go`: client lifecycle, headers, auth state, service wiring, token refresh routine.
  - `requests.go`: HTTP transport layer (`RequestService`) and multipart attachment upload.
  - `auth.go` + `auth.models.go`: login/public registration/token refresh models.
  - `projects.go`, `epics.go`, `user_stories.go`, `tasks.go`, `issues.go`, `milestones.go`, `users.go`, `webhooks.go`, `stats.go`, `resolver.go`, `wiki.go`: service methods.
  - `*.models.go`: DTOs/query structs/conversion helpers for each resource.
  - `*_custom_attributes_values.go`: generic CAVD wrappers.
  - `30` placeholder files with only `TODO` comments (e.g. `application_tokens.go`, `contact.go`, `timelines.go`, etc.).
- Tests:
  - `tests/`: integration tests expecting Taiga at `http://localhost:9000`.
  - Includes Docker setup helper scripts and seed data.
- Tooling/docs:
  - `README.md`, `examples/README.MD`, `CONTRIBUTION.md`, `CHANGELOG.md`.
  - `.github/workflows/go.yml`.

## 3) Architecture and conventions
- Entry point is `taigo.Client`.
- `Client.Initialise()`:
  - Sets headers, base API URL (`/api/v1` by default), disables pagination by default.
  - Instantiates all services.
  - Starts optional token-refresh ticker goroutine.
- Services:
  - Thin wrappers over `RequestService` (`GET/POST/PATCH/DELETE`) and `MakeURL`.
  - Optional `defaultProjectID` for project-scoped mapped services (`Project.ConfigureMappedServices`).
- Model pattern:
  - Generic object types (`Epic`, `Task`, `Issue`, etc.) plus meta pointers to concrete response variants (`Detail`, `DetailGET`, `DetailLIST`).
  - Conversion mostly via JSON round-trip helpers.

## 4) Implemented vs missing surface (high level)
- Implemented core resources (partial): auth, projects, epics, user stories, tasks, issues, milestones, users, webhooks, resolver, stats, wiki attachments.
- Major gaps:
  - Many endpoints are placeholders (`30` stub files).
  - Several implemented resources are only partial (e.g. `WikiService` has only attachment creation, `IssueService` lacks delete/by_ref methods, `TaskService` lacks edit/delete).

## 5) Testing reality
- Root package has no unit tests; tests live under `tests/` and are integration-heavy.
- Integration tests require a running Taiga stack and network access.
- In this environment:
  - `GOWORK=off GOCACHE=/tmp/... go test ./...` for root reaches tests but fails to connect to `localhost:9000` (sandbox network restriction).
  - `cli/` compiles (`go test ./...` => no test files).
  - `contribute/` does not currently compile as a module (missing `require` for replaced root module).

## 6) Key technical risks to remember
- Pagination toggle logic and header semantics are incorrect for Taiga backend behavior.
- Some query parameter structs cannot represent valid Taiga filters due to type/tag choices.
- Several methods have endpoint/path mismatches with Taiga (`webhooks/test`, custom attribute value models, watched/liked decoding).
- API consistency is weak across services (method signatures vary significantly by resource).

## 7) Practical commands
- List all Go files quickly:
  - `rg --files -g '*.go'`
- Run root compile/tests without workspace interference:
  - `GOWORK=off GOCACHE=/tmp/taigo-gocache go test ./...`
- Run CLI module:
  - `cd cli && GOWORK=off GOCACHE=/tmp/taigo-gocache-cli go test ./...`
- Run contribute module:
  - `cd contribute && GOWORK=off GOCACHE=/tmp/taigo-gocache-contrib go test ./...`

## 8) External reference sources for future checks
- Taiga API docs: [https://docs.taiga.io/api.html](https://docs.taiga.io/api.html)
- Taiga backend source: [https://github.com/taigaio/taiga-back](https://github.com/taigaio/taiga-back)

## 9) Self-check checklist before future changes
- Does endpoint path/method exactly match Taiga docs + `taiga-back` viewset route?
- Do request/response DTO field names match serializer fields (especially custom attributes values)?
- Are query params encoded in the same format Taiga expects (comma strings vs repeated keys, bool semantics)?
- Are service method signatures consistent with rest of the client API?
- Are integration tests updated/expanded for changed resource behavior?

