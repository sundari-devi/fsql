#!/bin/bash
set -e
set -o pipefail

OUTPUT_PATH="report.xml"
MODE=""

while [[ "$#" -gt 0 ]]; do
    case $1 in
        --output_path) OUTPUT_PATH="$2"; shift ;;
        base|new) MODE="$1" ;;
        *) echo "Unknown parameter: $1"; exit 1 ;;
    esac
    shift
done

if [ "$MODE" = "base" ]; then
    go test -v ./... 2>&1 | go-junit-report > "$OUTPUT_PATH"
elif [ "$MODE" = "new" ]; then
    go test -v -run TestParseLimit ./parser 2>&1 | go-junit-report > "$OUTPUT_PATH"
else
    echo "Error: Invalid mode. Use 'base' or 'new'."
    exit 1
fi
