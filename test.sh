#!/bin/bash

set -e
set -o pipefail

if [ "$1" = "base" ]; then
    go test -v ./parser ./tokenizer ./query 2>&1 | go-junit-report > report.xml
elif [ "$1" = "new" ]; then
    go test -v -run TestParseLimit ./parser 2>&1 | go-junit-report > report.xml
else
    echo "Error: Invalid mode. Use 'base' or 'new'."
    exit 1
fi
