#!/bin/bash
set -e
if [ "$1" = "base" ]; then
    go test ./parser ./tokenizer ./query
elif [ "$1" = "new" ]; then
    go test -run TestParseLimit ./parser
fi