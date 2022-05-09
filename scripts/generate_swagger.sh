#!/usr/bin/env bash

set -e

CONVERT_URL="https://converter.swagger.io/api/convert"
DOCS_DIR="./manager/docs"

if ! which "$GOPATH"/bin/swag &> /dev/null; then
    go install github.com/swaggo/swag/cmd/swag@latest
fi

if ! ls "$DOCS_DIR" &> /dev/null; then
  mkdir "$DOCS_DIR"
fi

# Generate 2.0 swagger
"$GOPATH"/bin/swag init -g ./manager/manager.go --output $DOCS_DIR
