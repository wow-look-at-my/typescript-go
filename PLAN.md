# Plan: Public Go library API for TypeScript compilation

## Context

All compiler code is under `internal/`, blocking external Go imports. The user wants to use this as a Go library: give it .ts files (from memory or disk), get a single compiled JS output (string or `io.Writer`). The bundle emit (`outFile`) infrastructure already exists.

## Approach

Create package `tsgo/` (import path `github.com/microsoft/typescript-go/tsgo`) — a single file that re-exports necessary internal types and provides two top-level functions.

## API

```go
package tsgo

import "io"

// Re-exported types (type aliases to internal packages)
type CompilerOptions = core.CompilerOptions
type Tristate = core.Tristate
type ScriptTarget = core.ScriptTarget
type ModuleKind = core.ModuleKind
type Diagnostic = ast.Diagnostic

// Re-exported constants
const (
	TSTrue    = core.TSTrue
	TSFalse   = core.TSFalse

	ES5    = core.ScriptTargetES5
	ES2015 = core.ScriptTargetES2015
	ES2020 = core.ScriptTargetES2020
	ES2022 = core.ScriptTargetES2022
	ESNext = core.ScriptTargetESNext
	// ... all targets

	ModuleNone     = core.ModuleKindNone
	ModuleCommonJS = core.ModuleKindCommonJS
	// ... etc
)

// Result holds compilation output.
type Result struct {
	JS          string
	SourceMap   string
	Diagnostics []*Diagnostic
}

// WriteTo writes the compiled JS to w.
func (r *Result) WriteTo(w io.Writer) (int64, error)

// Compile compiles in-memory TypeScript sources to a single JS output.
// sources maps filenames to source code.
// opts configures the compiler. OutFile and Module are set automatically.
// If opts is nil, defaults are used (Target: ESNext, Strict: true).
func Compile(sources map[string]string, opts *CompilerOptions) (*Result, error)

// CompileFiles reads .ts files from disk and compiles to a single JS output.
// If opts is nil, defaults are used.
func CompileFiles(filenames []string, opts *CompilerOptions) (*Result, error)
```

Callers use `core.CompilerOptions` directly (via the type alias) — no wrapper struct, no dumbing down. The functions set `OutFile` and `Module: ModuleKindNone` automatically since the whole point is single-file output.

## File to create

### `tsgo/tsgo.go`

**`Compile` implementation:**
1. Normalize source map keys to absolute paths under `/src/` (e.g. `"foo.ts"` → `"/src/foo.ts"`)
2. `vfstest.FromMap(normalizedSources, true)` → `bundled.WrapFS(fs)`
3. Clone `opts` (or create defaults), set `OutFile: "/output/bundle.js"`, `Module: ModuleKindNone`
4. `compiler.NewCompilerHost("/src", fs, bundled.LibPath(), nil, nil)`
5. `tsoptions.NewParsedCommandLine(opts, fileNames, comparePathsOpts)`
6. `compiler.NewProgram(compiler.ProgramOptions{Config: config, Host: host})`
7. Collect diagnostics: `program.GetSyntacticDiagnostics(ctx, nil)` + `program.GetSemanticDiagnostics(ctx, nil)`
8. `program.Emit(ctx, compiler.EmitOptions{WriteFile: capture})` — capture callback stores JS and `.map` text
9. Append emit diagnostics
10. Return `&Result{JS, SourceMap, Diagnostics}`

**`CompileFiles` implementation:**
1. Resolve filenames to absolute paths via `filepath.Abs`
2. Determine `currentDirectory` from common directory prefix
3. `osvfs.FS()` → `bundled.WrapFS(fs)`
4. Same steps 3-10 as `Compile`, but with real FS and real paths
5. `OutFile` set to a temp path like `/dev/null/bundle.js` (intercepted by WriteFile callback, never hits disk)

### `tsgo/tsgo_test.go`

- Single-file compile from memory
- Multi-file compile with references
- `CompileFiles` with temp files on disk
- `WriteTo` to a `bytes.Buffer`
- Type error produces diagnostics but still emits JS

## Internal functions reused

| Function | Location |
|----------|----------|
| `vfstest.FromMap` | `internal/vfs/vfstest/vfstest.go` |
| `osvfs.FS()` | `internal/vfs/osvfs/os.go` |
| `bundled.WrapFS`, `bundled.LibPath` | `internal/bundled/bundled.go` |
| `compiler.NewCompilerHost` | `internal/compiler/host.go` |
| `compiler.NewProgram` | `internal/compiler/program.go` |
| `compiler.EmitOptions`, `WriteFile` | `internal/compiler/program.go` |
| `tsoptions.NewParsedCommandLine` | `internal/tsoptions/parsedcommandline.go` |

## Verification

1. `go build ./tsgo/` — compiles
2. `go test ./tsgo/` — tests pass
3. `npx hereby build` — full project builds
4. `npx hereby test` — existing tests unaffected
