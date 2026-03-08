# Migration Guide: v1 to v2

This guide describes the breaking changes introduced in Taigo v2 and how to update existing code.

## 1) Behavioural Changes

### Pagination toggle

v1 behaviour was inconsistent because Taiga checks header presence, not value.

- `DisablePagination(true)` sets `x-disable-pagination`.
- `DisablePagination(false)` now removes `x-disable-pagination`.

No action is required unless you relied on the old broken behaviour.

### Typed API errors

Transport failures for non-2xx responses now return `*taigo.APIError`.

```go
resp, err := client.Request.Get(url, &out)
if err != nil {
	if apiErr, ok := err.(*taigo.APIError); ok {
		fmt.Println(apiErr.StatusCode, apiErr.Body)
	}
}
_ = resp
```

## 2) Signature Changes

### Tasks

Before:

```go
task, err := client.Task.Get(&taiga.Task{ID: 10})
```

After:

```go
task, err := client.Task.Get(10)
```

Before:

```go
task, err := client.Task.GetByRef(&taiga.Task{Ref: 5}, project)
```

After:

```go
task, err := client.Task.GetByRef(5, project)
```

### Webhooks

Before:

```go
log, err := client.Webhook.GetWebhookLog(&taiga.WebhookLog{ID: 22})
```

After:

```go
log, err := client.Webhook.GetWebhookLog(22)
```

### Users watched/liked

Before (single object):

```go
watched, err := client.User.GetWatchedContent(userID)
liked, err := client.User.GetLikedContent(userID)
```

After (slice result + optional filters):

```go
watched, err := client.User.GetWatchedContent(userID, nil)
liked, err := client.User.GetLikedContent(userID, &taigo.UsersHighlightedQueryParams{Type: "task"})
```

## 3) Query Struct Changes

### Optional bool filters

Several query structs now use `*bool` fields so `false` can be sent explicitly.

Before:

```go
q := taiga.MilestonesQueryParams{Closed: false} // false omitted due omitempty
```

After:

```go
q := taiga.MilestonesQueryParams{Closed: taiga.BoolPtr(false)}
```

### Task tags filter

Before:

```go
q := taiga.TasksQueryParams{Tags: []string{"backend", "api"}}
```

After:

```go
q := taiga.TasksQueryParams{}
q.SetTags("backend", "api")
```

## 4) New Services and Added Coverage

The following service families are now available from `Client` and project mapping:

- `Point`, `Priority`, `Severity`, `IssueType`
- `EpicStatus`, `IssueStatus`, `TaskStatus`, `UserStoryStatus`
- `EpicCustomAttribute`, `IssueCustomAttribute`, `TaskCustomAttribute`, `UserStoryCustomAttribute`

Expanded existing resources include:

- `IssueService`: `GetByRef`, `Delete`, attachment retrieval/listing
- `TaskService`: `Edit`, `Delete`, improved by-ref API
- `WikiService`: full CRUD + render + attachments
- `WebhookService`: corrected test endpoint and log access helpers

## 5) Naming Consistency

Where practical, services now expose `Update(...)` aliases for `Edit(...)`.
Existing `Edit(...)` calls continue to work.

## 6) Integration Tests

Integration tests now skip by default unless explicitly enabled:

```bash
TAIGO_RUN_INTEGRATION_TESTS=1 go test ./tests/...
```

This prevents false failures when Taiga is not running locally.

## 7) Recommended Upgrade Path

1. Update dependency to the v2 release tag.
2. Fix method signature changes listed above.
3. Update query struct initialisation for pointer-bool fields.
4. Run `go test ./...` and address compiler feedback.
5. If you use integration tests, run them with `TAIGO_RUN_INTEGRATION_TESTS=1`.
