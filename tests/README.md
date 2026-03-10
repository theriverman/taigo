# Testing

Run unit tests and offline-safe checks:

```bash
GOWORK=off GOCACHE=/tmp/taigo-gocache go test ./...
```

Run focused offline matrices:

```bash
go test ./... -run 'TestContractMatrixSingleRequestRoutes|TestQueryFilterMatrixEncoding|TestNegativeMatrixOfflineAPIErrors' -v
```

Integration tests under `tests/` require a reachable Taiga instance and are opt-in:

```bash
TAIGO_RUN_INTEGRATION_TESTS=1 GOWORK=off GOCACHE=/tmp/taigo-gocache go test ./tests/...
```

The integration harness supports environment overrides:

- `TAIGO_BASE_URL` (default: `http://localhost:9000`)
- `TAIGO_USERNAME` (default: `admin`)
- `TAIGO_PASSWORD` (default: `123123`)
- `TAIGO_PROJECT_ID` (default: `2`)
- `TAIGO_PROJECT_SLUG` (default: `taigo-test`)
- `TAIGO_USER_ID` (default: `5`)
- `TAIGO_MEMBER_USERNAME` / `TAIGO_MEMBER_PASSWORD` (optional role-matrix actor)
- `TAIGO_MEMBER_WRITE_EXPECTATION` (`forbid` or `allow`, default: `forbid`)

Run only the table-driven smoke matrix harness:

```bash
TAIGO_RUN_INTEGRATION_TESTS=1 TAIGO_PROJECT_ID=2 go test ./tests/... -run TestSmokeCRUDMatrix -v
```

Run the end-to-end workflow matrix:

```bash
TAIGO_RUN_INTEGRATION_TESTS=1 TAIGO_PROJECT_ID=2 go test ./tests/... -run TestWorkflowMatrixLive -v
```

Run negative-path checks against a live instance:

```bash
TAIGO_RUN_INTEGRATION_TESTS=1 TAIGO_PROJECT_ID=2 go test ./tests/... -run TestNegativeMatrixLive -v
```

Run auth/role matrix checks:

```bash
TAIGO_RUN_INTEGRATION_TESTS=1 go test ./tests/... -run TestAuthMatrixLive -v
TAIGO_RUN_INTEGRATION_TESTS=1 TAIGO_MEMBER_USERNAME=<user> TAIGO_MEMBER_PASSWORD=<pass> TAIGO_MEMBER_WRITE_EXPECTATION=forbid go test ./tests/... -run TestRoleMatrixLive -v
```

# Table-driven Real-Instance Harness for CRUD Smoke Matrix

```bash
export TAIGO_RUN_INTEGRATION_TESTS=1
export TAIGO_BASE_URL=http://127.0.0.1:9000
export TAIGO_USERNAME=admin
export TAIGO_PASSWORD=123123
export TAIGO_PROJECT_ID=1
go test ./tests/... -run 'TestSmokeCRUDMatrix|TestWorkflowMatrixLive|TestNegativeMatrixLive|TestAuthMatrixLive' -v
```

## Role Matrix #1
```bash
export TAIGO_RUN_INTEGRATION_TESTS=1
export TAIGO_BASE_URL=http://127.0.0.1:9000
export TAIGO_USERNAME=admin
export TAIGO_PASSWORD=123123
export TAIGO_PROJECT_ID=1
export TAIGO_MEMBER_USERNAME=demo1
export TAIGO_MEMBER_PASSWORD=123123
export TAIGO_MEMBER_WRITE_EXPECTATION=forbid  # allow|forbid
go test ./tests/... -run TestRoleMatrixLive -v
```

## Role Matrix #2
```bash
export TAIGO_RUN_INTEGRATION_TESTS=1
export TAIGO_BASE_URL=http://127.0.0.1:9000
export TAIGO_USERNAME=admin
export TAIGO_PASSWORD=123123
export TAIGO_PROJECT_ID=1
export TAIGO_MEMBER_USERNAME=admin
export TAIGO_MEMBER_PASSWORD=123123
export TAIGO_MEMBER_WRITE_EXPECTATION=allow  # allow|forbid
go test ./tests/... -run TestRoleMatrixLive -v
```