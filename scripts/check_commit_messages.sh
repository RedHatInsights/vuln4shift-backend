#!/usr/bin/env bash

commitTitle="$1"

# ignore merge requests
if echo "$commitTitle" | grep -qE "^Merge branch \'"; then
	echo "Commit hook: ignoring branch merge"
	exit 0
fi

# check semantic versioning scheme
if ! echo "$commitTitle" | grep -qE '^(feat|fix|docs|style|refactor|perf|test|chore|build|ci)\(?.*\)?:\s\w+'; then
	echo "Your commit title did not follow semantic versioning: $commitTitle"
	echo "Please see https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#commit-message-format"
	exit 1
fi
