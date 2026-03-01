# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

TypeScript 7 ("Corsa") — a native Go port of the TypeScript compiler, type checker, and language server. The original JS-based TypeScript compiler is referred to as "Strada". The reference implementation lives in `_submodules/TypeScript/`.

**This is a personal fork** — not intended for upstream merge. Upstream made design decisions that overcomplicate simple use cases. This fork restores and extends functionality to make the compiler work in a straightforward way.

### Fork changes: `outFile` bundle emit

Upstream Corsa removed `outFile` and marked it deprecated. This fork re-enables it:
- `outFile` works with `module: none` — concatenates all input files into a single output
- Source files are topologically sorted via `/// <reference path>` directives (Kahn's algorithm, falls back to original order on cycles)
- Source map support (inline and external) works across the bundle
- Duplicate `"use strict"` prologues are stripped from non-first files
- Key files: `emitter.go` (bundleEmitter), `program.go` (emitBundled path), `printer.go` (WriteBundleSourceFile), `outputpaths.go` (GetBundleOutputPaths)
- Declaration emit paths are computed but not yet implemented

## Build & Test Commands

This project uses `hereby` (a JS task runner). All tasks are defined in `Herebyfile.mjs`.

```sh
npx hereby build            # Build tsgo binary to built/local/tsgo
npx hereby test             # Run all tests
npx hereby lint             # Run linters (Go + custom lint rules)
npx hereby format           # Format code (gofumpt for Go, dprint for TS/JS)
npx hereby check:format     # Check formatting without modifying
npx hereby baseline-accept  # Accept new test baselines after test output changes
npx hereby generate         # Run all code generation (diagnostics, etc.)
npx hereby generate:enums   # Generate enum stringer methods
npx hereby install-tools    # Install additional tools (linters, etc.)
```

Standard `go` tooling also works (`go build`, `go test ./...`), but `hereby` is preferred for CI parity.

### Running specific tests

```sh
# Submodule tests (from _submodules/TypeScript/tests/cases/)
go test -run='TestSubmodule/<test name>' ./internal/testrunner

# Local tests (from testdata/tests/cases/compiler/)
go test -run='TestLocal/<test name>' ./internal/testrunner
```

### Environment variables for test variants

- `TSGO_HEREBY_RACE` — enable Go race detector
- `TSGO_HEREBY_NOEMBED` — build without embedded lib files
- `TSGO_HEREBY_CONCURRENT_TEST_PROGRAMS` — run test programs concurrently

## Architecture

### Entry point

`cmd/tsgo/main.go` — dispatches to:
- Compiler CLI (default) via `execute.CommandLine()`
- LSP server (`--lsp --stdio`)
- API server (`--api`)

### Core compiler pipeline (`internal/`)

| Package | Role |
|---------|------|
| `ast` | AST node definitions, kinds, flags, symbols, visitor |
| `parser` | Parsing and scanning |
| `binder` | Symbol binding and name resolution |
| `checker` | Type checker — the main type system (~21 files, largest package) |
| `compiler` | Program creation, file loading, emit orchestration |
| `printer` | AST-to-text printing, emit text writer |
| `transformers/` | Emit transforms: TS erasure, module systems, JSX, decorators, ES downlevel, declarations |
| `execute` | tsc-like CLI execution |
| `diagnostics` | Diagnostic message definitions (code-generated) |
| `core` | Shared utilities, compiler options, fundamental types |

### Language service & LSP

| Package | Role |
|---------|------|
| `ls` | Language service features: completions, hover, go-to-def, find-all-refs, formatting, organize imports |
| `lsp` | LSP server and JSON-RPC protocol layer (`lsproto/`) |
| `project` | Project/workspace management for the language service |
| `api` | Binary API server (MessagePack + JSON-RPC) for programmatic access |

### Key supporting packages

| Package | Role |
|---------|------|
| `bundled` | Embeds TypeScript lib files; `noembed` build tag disables embedding |
| `vfs` | Virtual filesystem abstraction (OS, cached, mock implementations) |
| `module` | Module resolution |
| `tspath` | TypeScript-compatible path utilities |
| `tsoptions` | tsconfig.json parsing |
| `outputpaths` | Output file path computation |
| `fourslash` | Fourslash test format implementation (language service integration tests) |
| `testrunner` | Compiler test runner (`TestSubmodule`, `TestLocal`) |

### npm packages (`_packages/`)

- `@typescript/native-preview` — ships the `tsgo` binary
- `@typescript/api` — programmatic API (JS/TS)
- `@typescript/ast` — AST types and node factory

## Test System

This project uses **snapshot/baseline/golden tests**, not primarily unit tests.

### Compiler tests

Test files in `testdata/tests/cases/compiler/*.ts` use comment-based compiler options:

```ts
// @strict: true
// @target: esnext
// @module: preserve

export const x: string = "hello";
```

**New tests must enable `@strict: true`** unless testing non-strict behavior. Use `@filename:` for multi-file tests. Comma-separated option values test multiple configurations.

### Baseline workflow

1. Tests produce output in `testdata/baselines/local/`
2. Reference baselines live in `testdata/baselines/reference/`
3. `.diff` files indicate behavioral divergence from Strada — reducing/eliminating them is progress
4. Run `npx hereby baseline-accept` to promote local baselines to reference
5. Use `git diff` to review what changed

### Fourslash tests

Language server integration tests live in `internal/fourslash/tests/`. Convert from TypeScript submodule via `npm run convertfourslash`.

## Code Porting Reference

Code in `internal/` is ported from `_submodules/TypeScript/`. When implementing features or fixing bugs, search the submodule for the original TypeScript function — it's the reference implementation. Some TypeScript files have been split into multiple Go files.

## Formatting

- Go: `gofumpt` (stricter than `gofmt`)
- TypeScript/JS: `dprint`

## Important Rules

- Do NOT add/change dependencies without being asked
- Do NOT remove debug assertions or `panic` calls — existing assertions are intentional
- Do NOT use `timeout` command when running tests
- Go 1.26+ required (uses `tool` and `ignore` directives in go.mod)
