#!/bin/sh

cd $(dirname $0)

if [[ ! -z $1 ]]; then
    if [[ "$1" == "db-init" ]]; then
        echo "Initializing database schema..."
        exit 0
    fi
fi

echo "Please specify service name as the first argument."
