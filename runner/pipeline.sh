#!/bin/bash
set -e

# Assign positional arguments to readable variables
RUN_TEST=$1

if [ "$RUN_TEST" = "true" ]; then
    echo ">> Running Tests..."
    go test -v ./...
fi


