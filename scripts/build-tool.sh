#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

echo "building tool $1"

INPUT_FILE="./tool/$1"
OUTPUT_FILE="./bin/$1"

echo "go build -o $OUTPUT_FILE $INPUT_FILE"

go build -o $OUTPUT_FILE -ldflags \
  "\
 " $INPUT_FILE
