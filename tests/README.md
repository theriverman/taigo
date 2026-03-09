# Table-driven Real-Instance Harness for CRUD Smoke Matrix

```bash
export TAIGO_RUN_INTEGRATION_TESTS=1
export TAIGO_BASE_URL=http://127.0.0.1:9000
export TAIGO_USERNAME=admin
export TAIGO_PASSWORD=123123
export TAIGO_PROJECT_ID=1
go test ./tests/... -run TestSmokeCRUDMatrix -v
```
