# Taigo v2 Review

## Scope
- Reviewed all repository Go source (`root`, `cli`, `contribute`, `tests`).
- Cross-checked behavior against:
  - Taiga REST docs: [docs.taiga.io/api.html](https://docs.taiga.io/api.html)
  - Taiga backend source: [taigaio/taiga-back](https://github.com/taigaio/taiga-back)
- Verified build/test reality locally where possible.

## Executive summary
- The project has a solid base client structure, but API correctness is currently mixed.
- There are several hard mismatches with Taiga behavior (not just missing coverage).
- Current surface area is incomplete (many endpoints stubbed, multiple services partial).
- For v2, a breaking redesign is justified and recommended.

---

## Critical findings (fix first)

### 1) Pagination toggle is functionally broken
- Local code:
  - [`client.go:197`](/Users/kristofdaja/Developer/taigo/client.go:197)
  - [`client.go:216`](/Users/kristofdaja/Developer/taigo/client.go:216)
- Problem:
  - `DisablePagination(false)` still sets `x-disable-pagination`.
  - Taiga backend disables pagination by header presence, not header value.
  - Header insertion uses `Add`, so repeated calls accumulate duplicate values.
- Backend proof:
  - `if "x-disable-pagination" in self.request.headers` in [pagination.py](https://github.com/taigaio/taiga-back/blob/main/taiga/base/api/pagination.py#L151-L176)
- Impact:
  - Library cannot reliably re-enable pagination once header is set.

### 2) Webhook test endpoint is wrong
- Local code:
  - [`webhooks.go:90`](/Users/kristofdaja/Developer/taigo/webhooks.go:90)
- Problem:
  - Calls `POST /webhooks/{id}`.
  - Taiga defines webhook test as detail action `POST /webhooks/{id}/test`.
- Backend proof:
  - `def test(...)` in [webhooks/api.py](https://github.com/taigaio/taiga-back/blob/main/taiga/webhooks/api.py#L35-L44)
- Impact:
  - `TestWebhook` likely fails (or hits wrong route).

### 3) Custom attribute value DTOs are incorrect for issue/task/user story
- Local code:
  - [`issue_custom_attributes_values.go:5`](/Users/kristofdaja/Developer/taigo/issue_custom_attributes_values.go:5)
  - [`task_custom_attributes_values.go:5`](/Users/kristofdaja/Developer/taigo/task_custom_attributes_values.go:5)
  - [`user_story_custom_attributes_values.go:5`](/Users/kristofdaja/Developer/taigo/user_story_custom_attributes_values.go:5)
- Problem:
  - All three use `Epic int 'json:"epic"'`.
  - Should be `issue`, `task`, `user_story` respectively.
- Backend proof:
  - [custom_attributes/serializers.py](https://github.com/taigaio/taiga-back/blob/main/taiga/projects/custom_attributes/serializers.py#L52-L65)
- Impact:
  - Serialization/deserialization mismatches; wrong payload contracts.

### 4) Project `order_by` never serializes
- Local code:
  - [`projects.models.go:458`](/Users/kristofdaja/Developer/taigo/projects.models.go:458)
- Problem:
  - `orderBy` is unexported (`lowercase`), so `go-querystring` does not encode it.
  - Sorting helper methods silently do nothing.
- Impact:
  - `ProjectsQueryParameters.TotalFans*()` etc. are ineffective.

### 5) `users/{id}/watched` and `users/{id}/liked` decode shape is wrong
- Local code:
  - [`users.go:73`](/Users/kristofdaja/Developer/taigo/users.go:73)
  - [`users.go:86`](/Users/kristofdaja/Developer/taigo/users.go:86)
- Problem:
  - Methods decode into single object (`UserWatched` / `UserLiked`).
  - Endpoint returns arrays/lists.
- Backend proof:
  - `response.Ok(response_data)` list in [users/api.py](https://github.com/taigaio/taiga-back/blob/main/taiga/users/api.py#L375-L429)
  - Docs sections: [users-watched](https://docs.taiga.io/api.html#users-watched), [users-liked](https://docs.taiga.io/api.html#users-liked)
- Impact:
  - Runtime decode errors or silently wrong behavior.

### 6) Required-field validation for numeric IDs is ineffective
- Local code:
  - [`common.functions.go:57`](/Users/kristofdaja/Developer/taigo/common.functions.go:57)
  - Examples: [`epics.go:51`](/Users/kristofdaja/Developer/taigo/epics.go:51), [`tasks.go:46`](/Users/kristofdaja/Developer/taigo/tasks.go:46), [`issues.go:93`](/Users/kristofdaja/Developer/taigo/issues.go:93), [`milestones.go:51`](/Users/kristofdaja/Developer/taigo/milestones.go:51)
- Problem:
  - `isEmpty()` does not treat `0` as empty, so required `project` ID checks do not work.
- Impact:
  - Missing required numeric fields are not caught client-side.

---

## High-priority issues

### 7) Bool query params with `omitempty` cannot express explicit `false`
- Local code examples:
  - [`milestones.models.go:42`](/Users/kristofdaja/Developer/taigo/milestones.models.go:42)
  - [`user_stories.models.go:333`](/Users/kristofdaja/Developer/taigo/user_stories.models.go:333)
  - [`issues.models.go:234`](/Users/kristofdaja/Developer/taigo/issues.models.go:234)
- Impact:
  - Cannot intentionally send false-valued filters where API semantics depend on explicit false.
- v2 fix:
  - Use `*bool` query fields for tri-state semantics.

### 8) Task tags query format mismatches Taiga filter expectations
- Local code:
  - [`tasks.models.go:250`](/Users/kristofdaja/Developer/taigo/tasks.models.go:250)
- Problem:
  - `[]string` encodes repeated query keys.
  - Backend parses `tags` as one comma-delimited string (`split(",")`).
- Backend proof:
  - [base/filters.py TagsFilter](https://github.com/taigaio/taiga-back/blob/main/taiga/base/filters.py#L478-L483)

### 9) Token refresh shutdown can panic/block
- Local code:
  - [`client.go:152`](/Users/kristofdaja/Developer/taigo/client.go:152)
- Problem:
  - `DisableAutomaticTokenRefresh` assumes ticker/channel are initialized.
  - If auto-refresh was disabled before init, ticker/channel may be nil.

### 10) Unimplemented HTTP methods panic instead of returning errors
- Local code:
  - [`requests.go:50`](/Users/kristofdaja/Developer/taigo/requests.go:50)
  - [`requests.go:89`](/Users/kristofdaja/Developer/taigo/requests.go:89)
  - [`requests.go:94`](/Users/kristofdaja/Developer/taigo/requests.go:94)
  - [`requests.go:99`](/Users/kristofdaja/Developer/taigo/requests.go:99)
- Impact:
  - Unexpected process crashes in library usage.

### 11) Service APIs are inconsistent and partial
- Examples:
  - `TaskService.Get(task *Task)` takes object instead of ID: [`tasks.go:58`](/Users/kristofdaja/Developer/taigo/tasks.go:58)
  - `IssueService` lacks delete/by_ref methods despite backend support: [`issues.go`](/Users/kristofdaja/Developer/taigo/issues.go)
  - `WikiService` only implements attachment creation: [`wiki.go:13`](/Users/kristofdaja/Developer/taigo/wiki.go:13)
- Impact:
  - Confusing developer experience and uneven capability.

---

## Medium-priority issues

### 12) `contribute` module does not build cleanly as-is
- Local code:
  - [`contribute/go.mod`](/Users/kristofdaja/Developer/taigo/contribute/go.mod)
- Problem:
  - Replaces root module but does not `require` it.

### 13) `go.work` includes only `./cli`
- Local code:
  - [`go.work`](/Users/kristofdaja/Developer/taigo/go.work)
- Impact:
  - Running `go test ./...` at repo root can behave unexpectedly in workspace mode.

### 14) CLI creates config directory with wrong permissions
- Local code:
  - [`cli/main.go:36`](/Users/kristofdaja/Developer/taigo/cli/main.go:36)
- Problem:
  - `os.MkdirAll(..., 0644)` for directory (should include execute bit, e.g., `0700`/`0755`).

### 15) CLI encryption design is weak
- Local code:
  - [`cli/passwordbasedencryption/pbewithmd5anddes.go`](/Users/kristofdaja/Developer/taigo/cli/passwordbasedencryption/pbewithmd5anddes.go)
- Problem:
  - MD5 + DES legacy scheme and low iteration count in CLI call sites.

### 16) Large surface still unimplemented
- Fact:
  - 30 root files are `TODO` stubs.
- Examples:
  - `applications.go`, `application_tokens.go`, `contact.go`, `timelines.go`, etc.

### 17) Test coverage holes
- Fact:
  - 12 test files, 4 marked `TODO` (`tasks`, `user_stories`, `wiki`, `webhooks`).
  - Root package has no unit tests.

---

## Design mistakes / failures

### A) “Meta pointer on generic struct” approach is costly and fragile
- Every item can carry pointers to full list/meta variants.
- Increases memory coupling and encourages hidden type-casts via JSON conversion.
- Better for v2: explicit typed responses per endpoint + shared field embeddings.

### B) API ergonomics are not coherent
- Resource method signatures differ arbitrarily (ID vs object pointer vs custom combinations).
- Naming also inconsistent (`CreateWebhook` vs `Create`, etc.).
- Better for v2: consistent resource interfaces and request options.

### C) Error model is underspecified
- Transport layer returns `error` as raw body string, often without status-rich typed errors.
- Better for v2: structured error type (`StatusCode`, `Body`, `RequestID`, parsed Taiga message).

---

## v2 proposal (breaking changes encouraged)

## 1) API redesign (public surface)
- Make all resource operations consistent:
  - `Create(ctx, input)`
  - `Get(ctx, id)`
  - `GetByRef(ctx, ref, projectSelector)`
  - `Update(ctx, id, patch)`
  - `Delete(ctx, id)`
  - `List(ctx, query)`
- Add `context.Context` to all network methods.
- Remove panics from public API.

## 2) Type system and models
- Replace `isEmpty(interface{})` with typed validation per request DTO.
- Use endpoint-specific request/response types; avoid JSON round-trip conversion helpers for core flow.
- Query structs:
  - `*bool` for optional booleans.
  - Comma-encoded string helpers for filters that require it.

## 3) Transport layer hardening
- Header handling should use `Set` where appropriate.
- Pagination toggle:
  - send `x-disable-pagination` only when disabling.
  - remove header when enabling.
- Introduce structured typed errors.

## 4) Coverage completion strategy
- Implement remaining Taiga routes in phases:
  1. Existing partially implemented resources (`tasks`, `issues`, `wiki`, `users`, `webhooks`)
  2. Status/type/priority/points/severity/custom-attribute endpoints
  3. Memberships, timelines, history, imports/exports, applications

## 5) Testing strategy
- Add fast unit tests for:
  - query encoding,
  - endpoint path construction,
  - model JSON mapping,
  - header behavior and pagination.
- Keep integration suite, but gate via env flags and avoid hard panics when backend unavailable.

---

## 6) Documentation update
- Current README.md must be updated as the last step once the above listed items were corrected
- Documentation must meet the rules of British English
- A `MIGRATION.md` file shall be added explaining present users of the library how to migrate from v1 to v2 in the future

## Verification notes from this review
- `go test ./...` (root) reaches integration tests but cannot connect to `localhost:9000` in this sandbox.
- `cli/` compiles (`no test files`).
- `contribute/` currently fails module setup as noted above.

---

## References used
- Taiga docs:
  - [API root](https://docs.taiga.io/api.html)
  - [Pagination](https://docs.taiga.io/api.html#_pagination)
  - [Users watched](https://docs.taiga.io/api.html#users-watched)
  - [Users liked](https://docs.taiga.io/api.html#users-liked)
  - [Webhooks](https://docs.taiga.io/api.html#webhooks)
  - [Webhook logs](https://docs.taiga.io/api.html#webhooklogs)
- Taiga backend:
  - [webhooks/api.py](https://github.com/taigaio/taiga-back/blob/main/taiga/webhooks/api.py)
  - [projects/custom_attributes/serializers.py](https://github.com/taigaio/taiga-back/blob/main/taiga/projects/custom_attributes/serializers.py)
  - [base/api/pagination.py](https://github.com/taigaio/taiga-back/blob/main/taiga/base/api/pagination.py)
  - [users/api.py](https://github.com/taigaio/taiga-back/blob/main/taiga/users/api.py)
  - [base/filters.py](https://github.com/taigaio/taiga-back/blob/main/taiga/base/filters.py)
  - [projects/api.py](https://github.com/taigaio/taiga-back/blob/main/taiga/projects/api.py)

