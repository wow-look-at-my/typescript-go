// Package tsgo provides a public Go API for compiling TypeScript sources.
//
// It wraps the internal compiler pipeline to offer two entry points:
//   - [Compile] compiles in-memory TypeScript sources to a single JS output.
//   - [CompileFiles] reads .ts files from disk and compiles them.
//
// Both functions produce a single bundled JS file (using outFile mode).
package tsgo

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
)

// Re-exported types from internal packages.
type CompilerOptions = core.CompilerOptions
type Tristate = core.Tristate
type ScriptTarget = core.ScriptTarget
type ModuleKind = core.ModuleKind
type Diagnostic = ast.Diagnostic

// Tristate values.
const (
	TSTrue    = core.TSTrue
	TSFalse   = core.TSFalse
	TSUnknown = core.TSUnknown
)

// ScriptTarget values.
const (
	ScriptTargetES5    = core.ScriptTargetES5
	ScriptTargetES2015 = core.ScriptTargetES2015
	ScriptTargetES2016 = core.ScriptTargetES2016
	ScriptTargetES2017 = core.ScriptTargetES2017
	ScriptTargetES2018 = core.ScriptTargetES2018
	ScriptTargetES2019 = core.ScriptTargetES2019
	ScriptTargetES2020 = core.ScriptTargetES2020
	ScriptTargetES2021 = core.ScriptTargetES2021
	ScriptTargetES2022 = core.ScriptTargetES2022
	ScriptTargetES2023 = core.ScriptTargetES2023
	ScriptTargetES2024 = core.ScriptTargetES2024
	ScriptTargetES2025 = core.ScriptTargetES2025
	ScriptTargetESNext = core.ScriptTargetESNext
)

// ModuleKind values.
const (
	ModuleKindNone     = core.ModuleKindNone
	ModuleKindCommonJS = core.ModuleKindCommonJS
	ModuleKindES2015   = core.ModuleKindES2015
	ModuleKindES2020   = core.ModuleKindES2020
	ModuleKindES2022   = core.ModuleKindES2022
	ModuleKindESNext   = core.ModuleKindESNext
	ModuleKindNode16   = core.ModuleKindNode16
	ModuleKindNodeNext = core.ModuleKindNodeNext
	ModuleKindPreserve = core.ModuleKindPreserve
)

// Result holds compilation output.
type Result struct {
	JS          string
	SourceMap   string
	Diagnostics []*Diagnostic
}

// WriteTo writes the compiled JS to w, implementing [io.WriterTo].
func (r *Result) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, r.JS)
	return int64(n), err
}

// Compile compiles in-memory TypeScript sources to a single JS output.
// sources maps filenames (e.g. "index.ts") to source code.
// Relative filenames are normalized under a virtual /src/ directory.
// opts configures the compiler; OutFile and Module are set automatically.
// If opts is nil, defaults are used (Target: ESNext, Strict: true).
func Compile(sources map[string]string, opts *CompilerOptions) (*Result, error) {
	if len(sources) == 0 {
		return nil, fmt.Errorf("tsgo: no source files provided")
	}

	// Normalize source paths to absolute paths under /src/.
	normalized := make(map[string]string, len(sources))
	var fileNames []string
	for name, content := range sources {
		absName := name
		if !strings.HasPrefix(name, "/") {
			absName = "/src/" + name
		}
		normalized[absName] = content
		fileNames = append(fileNames, absName)
	}

	// Create virtual filesystem with bundled lib files.
	fs := vfstest.FromMap(normalized, true)
	fs = bundled.WrapFS(fs)

	// Prepare compiler options.
	compilerOpts := prepareOptions(opts)
	compilerOpts.OutFile = "/output/bundle.js"

	host := compiler.NewCompilerHost("/src", fs, bundled.LibPath(), nil, nil)

	config := tsoptions.NewParsedCommandLine(
		compilerOpts,
		fileNames,
		tspath.ComparePathsOptions{
			UseCaseSensitiveFileNames: true,
			CurrentDirectory:          "/src",
		},
	)

	program := compiler.NewProgram(compiler.ProgramOptions{
		Config: config,
		Host:   host,
	})

	return emitResult(program)
}

// CompileFiles reads .ts files from disk and compiles to a single JS output.
// Filenames are resolved to absolute paths. If opts is nil, defaults are used.
func CompileFiles(filenames []string, opts *CompilerOptions) (*Result, error) {
	if len(filenames) == 0 {
		return nil, fmt.Errorf("tsgo: no source files provided")
	}

	// Resolve to absolute paths.
	absNames := make([]string, len(filenames))
	for i, name := range filenames {
		abs, err := filepath.Abs(name)
		if err != nil {
			return nil, fmt.Errorf("tsgo: resolving path %q: %w", name, err)
		}
		absNames[i] = filepath.ToSlash(abs)
	}

	// Determine working directory from common prefix.
	currentDir := commonDir(absNames)

	// Create OS filesystem with bundled lib files.
	fs := osvfs.FS()
	fs = bundled.WrapFS(fs)

	// Prepare compiler options.
	compilerOpts := prepareOptions(opts)
	// Use a synthetic OutFile path — the WriteFile callback intercepts output before it hits disk.
	compilerOpts.OutFile = currentDir + "/__tsgo_bundle.js"

	host := compiler.NewCompilerHost(currentDir, fs, bundled.LibPath(), nil, nil)

	config := tsoptions.NewParsedCommandLine(
		compilerOpts,
		absNames,
		tspath.ComparePathsOptions{
			UseCaseSensitiveFileNames: true,
			CurrentDirectory:          currentDir,
		},
	)

	program := compiler.NewProgram(compiler.ProgramOptions{
		Config: config,
		Host:   host,
	})

	return emitResult(program)
}

// prepareOptions clones or creates default compiler options, ensuring
// Module is set to None (required for outFile bundling).
func prepareOptions(opts *CompilerOptions) *CompilerOptions {
	if opts == nil {
		return &CompilerOptions{
			Target: core.ScriptTargetESNext,
			Strict: core.TSTrue,
			Module: core.ModuleKindNone,
		}
	}
	// Clone so we don't mutate the caller's struct.
	clone := opts.Clone()
	clone.Module = core.ModuleKindNone
	return clone
}

// emitResult runs the compiler pipeline and collects output.
func emitResult(program *compiler.Program) (*Result, error) {
	ctx := context.Background()

	var diags []*ast.Diagnostic
	diags = append(diags, program.GetSyntacticDiagnostics(ctx, nil)...)
	diags = append(diags, program.GetSemanticDiagnostics(ctx, nil)...)

	var js, sourceMap string

	emitResult := program.Emit(ctx, compiler.EmitOptions{
		WriteFile: func(fileName string, text string, writeByteOrderMark bool, data *compiler.WriteFileData) error {
			if strings.HasSuffix(fileName, ".map") {
				sourceMap = text
			} else {
				js = text
			}
			return nil
		},
	})

	diags = append(diags, emitResult.Diagnostics...)

	return &Result{
		JS:          js,
		SourceMap:   sourceMap,
		Diagnostics: diags,
	}, nil
}

// commonDir returns the longest common directory prefix of the given absolute paths.
func commonDir(paths []string) string {
	if len(paths) == 0 {
		return "/"
	}

	// Start with the directory of the first path.
	common := paths[0]
	if idx := strings.LastIndex(common, "/"); idx > 0 {
		common = common[:idx]
	}

	// Narrow down by each subsequent path.
	for _, p := range paths[1:] {
		for !strings.HasPrefix(p, common+"/") && common != "/" {
			if idx := strings.LastIndex(common, "/"); idx > 0 {
				common = common[:idx]
			} else {
				common = "/"
			}
		}
	}

	return common
}
