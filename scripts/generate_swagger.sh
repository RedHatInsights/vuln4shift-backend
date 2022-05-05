#!/usr/bin/env bash

CONVERT_URL="https://converter.swagger.io/api/convert"
DOCS_TMP="/tmp"

which $GOPATH/bin/swag &> /dev/null

if [ "$?" != "0" ]; then
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate 2.0 swagger
$GOPATH/bin/swag init -g ./manager/manager.go --output $DOCS_TMP

# Convert 2.0 -> 3.0
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" \
     -d @$DOCS_TMP/swagger.json $CONVERT_URL > ./manager/docs/swagger.json
