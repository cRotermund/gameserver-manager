#!/bin/sh
ROOT="$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd -P)"
(cd "$ROOT/harness" && exec go run . "$@")