#!/usr/bin/env bash

printf "Running manager unit tests...\n"
go test -v ./manager/...
