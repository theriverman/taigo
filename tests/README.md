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
export TAIGO_MEMBER_WRITE_EXPECTATION=forbid  # allow|forbid
go test ./tests/... -run TestRoleMatrixLive -v
```