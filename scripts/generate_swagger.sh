#!/usr/bin/env bash

set -e

CONVERT_URL="https://converter.swagger.io/api/convert"
DOCS_TMP="/tmp"

GOPATH=$(go env GOPATH) # Not exported in centos:stream8 image

if ! which "$GOPATH"/bin/swag &> /dev/null; then
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate 2.0 swagger
"$GOPATH"/bin/swag init -g ./manager/manager.go --output $DOCS_TMP

if ! ls ./manager/docs &> /dev/null; then
  mkdir ./manager/docs
fi

# Convert 2.0 -> 3.0
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" \
     -d @$DOCS_TMP/swagger.json $CONVERT_URL > ./manager/docs/swagger.json
