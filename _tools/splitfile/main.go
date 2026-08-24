// splitfile extracts top-level declarations from a Go source file into a new file.
//
// Usage:
//
//	go run _tools/splitfile/main.go -src file.go -dst newfile.go -names "Foo,Bar,Baz"
//	go run _tools/splitfile/main.go -src file.go -dst newfile.go -startline 100 -endline 500
//
// Modes:
//   - -names: extract specific top-level declarations (functions, methods, types, vars, consts) by name
//   - -startline/-endline: extract all top-level declarations whose start line falls in [startline, endline]
//
// The tool copies the full import block from source to destination, then uses
// go/ast to determine which imports are actually needed in each file and removes
// unused ones. No external tools (goimports, gofumpt) are required.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	src := flag.String("src", "", "Source .go file to split")
	dst := flag.String("dst", "", "Destination .go file for extracted declarations")
	names := flag.String("names", "", "Comma-separated list of declaration names to extract")
	startLine := flag.Int("startline", 0, "Extract declarations starting at or after this line")
	endLine := flag.Int("endline", 0, "Extract declarations starting at or before this line")
	dryRun := flag.Bool("dry-run", false, "Print what would be extracted without modifying files")
	flag.Parse()

	if *src == "" || *dst == "" {
		fmt.Fprintln(os.Stderr, "Usage: splitfile -src <file.go> -dst <newfile.go> [-names Foo,Bar] [-startline N -endline M]")
		os.Exit(1)
	}

	if *names == "" && (*startLine == 0 && *endLine == 0) {
		fmt.Fprintln(os.Stderr, "Error: must specify either -names or -startline/-endline")
		os.Exit(1)
	}

	content, err := os.ReadFile(*src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", *src, err)
		os.Exit(1)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, *src, content, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", *src, err)
		os.Exit(1)
	}

	// Build the set of names to extract
	nameSet := make(map[string]bool)
	if *names != "" {
		for _, n := range strings.Split(*names, ",") {
			n = strings.TrimSpace(n)
			if n != "" {
				nameSet[n] = true
			}
		}
	}

	// Collect byte ranges to extract (sorted by position)
	var extractRegions []region
	var keepRegions []region

	for _, decl := range file.Decls {
		// Skip import declarations
		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.IMPORT {
			continue
		}

		declStart := fset.Position(decl.Pos())
		declEnd := fset.Position(decl.End())

		// Include preceding doc comment
		var docComment *ast.CommentGroup
		switch d := decl.(type) {
		case *ast.FuncDecl:
			docComment = d.Doc
		case *ast.GenDecl:
			docComment = d.Doc
		}

		byteStart := declStart.Offset
		if docComment != nil {
			commentStart := fset.Position(docComment.Pos())
			byteStart = commentStart.Offset
		}
		if byteStart > 0 && content[byteStart-1] == '\n' {
			byteStart--
		}

		shouldExtract := false
		declName := getDeclName(decl, nameSet)

		if len(nameSet) > 0 {
			shouldExtract = nameSet[declName]
		} else if *startLine > 0 || *endLine > 0 {
			line := declStart.Line
			shouldExtract = (*startLine <= 0 || line >= *startLine) && (*endLine <= 0 || line <= *endLine)
		}

		r := region{start: byteStart, end: declEnd.Offset, name: declName, decl: decl}
		if shouldExtract {
			extractRegions = append(extractRegions, r)
		} else {
			keepRegions = append(keepRegions, r)
		}
	}

	if len(extractRegions) == 0 {
		fmt.Fprintln(os.Stderr, "No matching declarations found")
		os.Exit(1)
	}

	sort.Slice(extractRegions, func(i, j int) bool {
		return extractRegions[i].start < extractRegions[j].start
	})

	if *dryRun {
		fmt.Printf("Would extract %d declarations from %s to %s:\n", len(extractRegions), *src, *dst)
		for _, r := range extractRegions {
			startPos := fset.Position(token.Pos(r.start + 1))
			fmt.Printf("  %s (lines %d-%d)\n", r.name, startPos.Line, fset.Position(token.Pos(r.end)).Line)
		}
		return
	}

	// Collect imports from the original file
	var imports []*ast.ImportSpec
	for _, imp := range file.Imports {
		imports = append(imports, imp)
	}

	// Determine which imports are needed for extracted vs kept declarations
	extractedIdents := collectIdents(extractRegions)
	keptIdents := collectIdents(keepRegions)

	// Build destination file
	var dstBuf bytes.Buffer
	fmt.Fprintf(&dstBuf, "package %s\n", file.Name.Name)
	writeImports(&dstBuf, imports, extractedIdents, fset)
	for _, r := range extractRegions {
		dstBuf.WriteByte('\n')
		dstBuf.Write(content[r.start:r.end])
		dstBuf.WriteByte('\n')
	}

	// Build source file with regions removed
	var srcBuf bytes.Buffer
	cursor := 0
	for _, r := range extractRegions {
		if r.start > cursor {
			srcBuf.Write(content[cursor:r.start])
		}
		cursor = r.end
		for cursor < len(content) && content[cursor] == '\n' {
			cursor++
		}
	}
	if cursor < len(content) {
		srcBuf.Write(content[cursor:])
	}

	// Rewrite source file imports to only include what's needed
	srcContent := rewriteImports(srcBuf.Bytes(), imports, keptIdents, fset)

	if err := os.WriteFile(*dst, dstBuf.Bytes(), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", *dst, err)
		os.Exit(1)
	}
	if err := os.WriteFile(*src, srcContent, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", *src, err)
		os.Exit(1)
	}

	srcLines := bytes.Count(srcContent, []byte("\n"))
	dstLines := bytes.Count(dstBuf.Bytes(), []byte("\n"))
	fmt.Printf("Split %d decls: %s -> %s (%d lines), src now %d lines\n",
		len(extractRegions), *src, *dst, dstLines, srcLines)
}

// collectIdents returns all identifier names used in the given regions' declarations.
func collectIdents(regions []region) map[string]bool {
	idents := make(map[string]bool)
	for _, r := range regions {
		ast.Inspect(r.decl, func(n ast.Node) bool {
			if id, ok := n.(*ast.Ident); ok {
				idents[id.Name] = true
			}
			if sel, ok := n.(*ast.SelectorExpr); ok {
				if id, ok := sel.X.(*ast.Ident); ok {
					idents[id.Name] = true
				}
			}
			return true
		})
	}
	return idents
}

// writeImports writes an import block containing only the imports that are referenced
// by the given set of identifiers.
func writeImports(buf *bytes.Buffer, imports []*ast.ImportSpec, usedIdents map[string]bool, fset *token.FileSet) {
	var needed []*ast.ImportSpec
	for _, imp := range imports {
		name := importLocalName(imp)
		if name == "_" || name == "." || usedIdents[name] {
			needed = append(needed, imp)
		}
	}
	if len(needed) == 0 {
		return
	}
	buf.WriteString("\nimport (\n")
	for _, imp := range needed {
		buf.WriteByte('\t')
		if imp.Name != nil {
			buf.WriteString(imp.Name.Name)
			buf.WriteByte(' ')
		}
		buf.WriteString(imp.Path.Value)
		buf.WriteByte('\n')
	}
	buf.WriteString(")\n")
}

// rewriteImports parses srcContent and rewrites its import block to only include
// imports that are actually used.
func rewriteImports(srcContent []byte, originalImports []*ast.ImportSpec, usedIdents map[string]bool, origFset *token.FileSet) []byte {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", srcContent, parser.ParseComments)
	if err != nil {
		return srcContent // bail out, don't break the file
	}

	// Find the import GenDecl and rewrite it
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.IMPORT {
			continue
		}
		var kept []ast.Spec
		for _, spec := range gd.Specs {
			imp, ok := spec.(*ast.ImportSpec)
			if !ok {
				continue
			}
			name := importLocalName(imp)
			if name == "_" || name == "." || usedIdents[name] {
				kept = append(kept, spec)
			}
		}
		gd.Specs = kept
		break
	}

	var buf bytes.Buffer
	printer.Fprint(&buf, fset, f)
	return buf.Bytes()
}

// importLocalName returns the local name that an import would be referenced by.
func importLocalName(imp *ast.ImportSpec) string {
	if imp.Name != nil {
		return imp.Name.Name
	}
	path, _ := strconv.Unquote(imp.Path.Value)
	parts := strings.Split(path, "/")
	name := parts[len(parts)-1]
	// Handle hyphenated package names (use last segment after hyphen)
	if idx := strings.LastIndex(name, "-"); idx >= 0 {
		name = name[idx+1:]
	}
	return name
}

func getDeclName(decl ast.Decl, nameSet map[string]bool) string {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return d.Name.Name
	case *ast.GenDecl:
		for _, spec := range d.Specs {
			switch s := spec.(type) {
			case *ast.TypeSpec:
				return s.Name.Name
			case *ast.ValueSpec:
				for _, n := range s.Names {
					if len(nameSet) == 0 || nameSet[n.Name] {
						return n.Name
					}
				}
				if len(s.Names) > 0 {
					return s.Names[0].Name
				}
			}
		}
	}
	return ""
}
