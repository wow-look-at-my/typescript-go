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
// The tool preserves the package declaration, adds necessary imports via goimports,
// and removes extracted declarations from the source file.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"sort"
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
	type region struct {
		start int // byte offset in content
		end   int // byte offset in content (exclusive)
		name  string
	}
	var regions []region

	for _, decl := range file.Decls {
		declStart := fset.Position(decl.Pos())
		declEnd := fset.Position(decl.End())

		// Include preceding doc comment if any
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

		// Also include any preceding blank line / comment that's
		// directly attached (scan backwards for whitespace)
		for byteStart > 0 && content[byteStart-1] == '\n' {
			// Include one preceding newline for spacing
			byteStart--
			break
		}

		shouldExtract := false
		declName := ""

		switch d := decl.(type) {
		case *ast.FuncDecl:
			declName = d.Name.Name
			if d.Recv != nil && len(d.Recv.List) > 0 {
				// Method — use "receiver.method" for display but match on method name
				declName = d.Name.Name
			}
		case *ast.GenDecl:
			// type, var, const — extract if any spec name matches
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					declName = s.Name.Name
				case *ast.ValueSpec:
					for _, n := range s.Names {
						if nameSet[n.Name] {
							declName = n.Name
							break
						}
					}
				}
			}
		}

		if len(nameSet) > 0 {
			shouldExtract = nameSet[declName]
		} else if *startLine > 0 || *endLine > 0 {
			line := declStart.Line
			if *startLine > 0 && line < *startLine {
				shouldExtract = false
			} else if *endLine > 0 && line > *endLine {
				shouldExtract = false
			} else {
				shouldExtract = true
			}
		}

		if shouldExtract {
			regions = append(regions, region{
				start: byteStart,
				end:   declEnd.Offset,
				name:  declName,
			})
		}
	}

	if len(regions) == 0 {
		fmt.Fprintln(os.Stderr, "No matching declarations found")
		os.Exit(1)
	}

	// Sort by position
	sort.Slice(regions, func(i, j int) bool {
		return regions[i].start < regions[j].start
	})

	if *dryRun {
		fmt.Printf("Would extract %d declarations from %s to %s:\n", len(regions), *src, *dst)
		for _, r := range regions {
			startPos := fset.Position(token.Pos(r.start + 1))
			fmt.Printf("  %s (lines %d-%d)\n", r.name, startPos.Line, fset.Position(token.Pos(r.end)).Line)
		}
		return
	}

	// Build the destination file content
	var dstBuf bytes.Buffer
	fmt.Fprintf(&dstBuf, "package %s\n\n", file.Name.Name)
	for _, r := range regions {
		dstBuf.Write(content[r.start:r.end])
		dstBuf.WriteByte('\n')
		dstBuf.WriteByte('\n')
	}

	// Build the source file with regions removed
	var srcBuf bytes.Buffer
	cursor := 0
	for _, r := range regions {
		if r.start > cursor {
			srcBuf.Write(content[cursor:r.start])
		}
		cursor = r.end
		// Skip trailing newlines after the removed region
		for cursor < len(content) && content[cursor] == '\n' {
			cursor++
		}
	}
	if cursor < len(content) {
		srcBuf.Write(content[cursor:])
	}

	// Write files
	if err := os.WriteFile(*dst, dstBuf.Bytes(), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", *dst, err)
		os.Exit(1)
	}
	if err := os.WriteFile(*src, srcBuf.Bytes(), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", *src, err)
		os.Exit(1)
	}

	// Run gofumpt on both files to fix formatting
	for _, f := range []string{*src, *dst} {
		cmd := exec.Command("go", "tool", "mvdan.cc/gofumpt", "-w", f)
		if out, err := cmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: gofumpt failed on %s: %v\n%s\n", f, err, out)
		}
	}

	// Run goimports on both files to fix imports
	for _, f := range []string{*src, *dst} {
		cmd := exec.Command("go", "tool", "golang.org/x/tools/cmd/goimports", "-w", f)
		if out, err := cmd.CombinedOutput(); err != nil {
			// goimports may not be available as a go tool, try PATH
			cmd2 := exec.Command("goimports", "-w", f)
			if out2, err2 := cmd2.CombinedOutput(); err2 != nil {
				fmt.Fprintf(os.Stderr, "Warning: goimports unavailable for %s: %v\n%s%s\n", f, err, out, out2)
			}
		}
	}

	// Count lines in resulting files
	srcLines := bytes.Count(srcBuf.Bytes(), []byte("\n"))
	dstLines := bytes.Count(dstBuf.Bytes(), []byte("\n"))
	fmt.Printf("Extracted %d declarations: %s (%d lines) -> %s (%d lines), src now %d lines\n",
		len(regions), *src, dstLines, *dst, dstLines, srcLines)
}
