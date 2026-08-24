#!/usr/bin/env bash
set -euo pipefail

# batch_split.sh - Split multiple Go files in one go.
# Run from repo root: bash _tools/splitfile/batch_split.sh

TOOL_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$TOOL_DIR/../.." && pwd)"

# Track all modified files for goimports at the end
MODIFIED_FILES=()

split() {
	local src="$1" dst="$2"
	shift 2
	echo "  $src -> $dst"
	cd "$REPO_ROOT/_tools"
	go run ./splitfile/ -src "$REPO_ROOT/$src" -dst "$REPO_ROOT/$dst" "$@"
	cd "$REPO_ROOT"
	MODIFIED_FILES+=("$src" "$dst")
}

echo "==> Batch 4: Files just over 750 lines"

# internal/core/core.go (774 lines)
split internal/core/core.go internal/core/core_text.go -startline 399 -endline 9999

# internal/astnav/tokens.go (758 lines)
split internal/astnav/tokens.go internal/astnav/tokens_classify.go -startline 380 -endline 9999

# internal/vfs/vfstest/vfstest_test.go (751 lines)
split internal/vfs/vfstest/vfstest_test.go internal/vfs/vfstest/vfstest_map_test.go -startline 380 -endline 9999

# internal/transformers/estransforms/using.go (789 lines)
split internal/transformers/estransforms/using.go internal/transformers/estransforms/using_dispose.go -startline 395 -endline 9999

# internal/checker/symbolaccessibility.go (786 lines)
split internal/checker/symbolaccessibility.go internal/checker/symbolaccessibility_checks.go -startline 393 -endline 9999

# internal/ls/hover.go (780 lines)
split internal/ls/hover.go internal/ls/hover_docs.go -startline 390 -endline 9999

# internal/tspath/path_test.go (818 lines)
split internal/tspath/path_test.go internal/tspath/path_relative_test.go -startline 410 -endline 9999

# internal/execute/build/buildtask.go (854 lines)
split internal/execute/build/buildtask.go internal/execute/build/buildtask_exec.go -startline 430 -endline 9999

# internal/api/encoder/encoder.go (882 lines)
split internal/api/encoder/encoder.go internal/api/encoder/encoder_types.go -startline 445 -endline 9999

# internal/project/snapshotfs_test.go (880 lines)
split internal/project/snapshotfs_test.go internal/project/snapshotfs_overlay_test.go -startline 440 -endline 9999

# internal/semver/version_range_test.go (959 lines) — split around line 300 to balance
split internal/semver/version_range_test.go internal/semver/version_range_match_test.go -startline 300 -endline 700

# internal/project/session_test.go (956 lines) — split around line 300
split internal/project/session_test.go internal/project/session_refs_test.go -startline 300 -endline 700

echo ""
echo "==> Running goimports on all modified files..."
cd "$REPO_ROOT"
for f in "${MODIFIED_FILES[@]}"; do
	if [[ -f "$f" ]]; then
		go tool golang.org/x/tools/cmd/goimports -w "$f" 2>/dev/null || true
	fi
done

echo "==> Running gofumpt on all modified files..."
for f in "${MODIFIED_FILES[@]}"; do
	if [[ -f "$f" ]]; then
		go tool mvdan.cc/gofumpt -w "$f" 2>/dev/null || true
	fi
done

echo ""
echo "==> Batch 4 complete. File counts:"
for f in "${MODIFIED_FILES[@]}"; do
	if [[ -f "$f" ]]; then
		echo "  $(wc -l < "$f") $f"
	fi
done
