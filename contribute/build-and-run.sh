#!/usr/bin/env bash

set -e

go build -o taigo-dev -ldflags \
" \
  -X main.taigaUsername=$taigaUsername \
  -X main.taigaPassword=$taigaPassword \
  -X main.sandboxProjectSlug=$sandboxProjectSlug \
  -X main.sandboxEpicID=$sandboxEpicID \
  -X main.sandboxFileUploadPath=$sandboxFileUploadPath \
"

./taigo-dev
