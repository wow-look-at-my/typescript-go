#!/usr/bin/env bash
set -euo pipefail

# split.sh - Split a Go source file by extracting declarations to a new file.
#
# Usage:
#   ./_tools/splitfile/split.sh <src> <dst> [-names "Foo,Bar"] [-startline N -endline M]
#
# Runs the splitfile Go tool then gofumpt on both files.

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$REPO_ROOT/_tools"
go run ./splitfile/ "$@"

# Fix formatting on both files
SRC=""
DST=""
while [[ $# -gt 0 ]]; do
	case "$1" in
		-src) SRC="$2"; shift 2 ;;
		-dst) DST="$2"; shift 2 ;;
		*) shift ;;
	esac
done

cd "$REPO_ROOT"
if [[ -n "$SRC" && -f "$SRC" ]]; then
	go tool mvdan.cc/gofumpt -w "$SRC"
fi
if [[ -n "$DST" && -f "$DST" ]]; then
	go tool mvdan.cc/gofumpt -w "$DST"
fi
