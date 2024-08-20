#!/bin/bash
set -euo pipefail

go test ./lambdas/... -coverprofile=./cover.out -covermode=atomic -coverpkg=./lambdas/...

${GOPATH}/bin/go-test-coverage -config=./.testcoverage.yml