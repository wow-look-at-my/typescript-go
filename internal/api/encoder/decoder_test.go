package encoder_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/microsoft/typescript-go/internal/api/encoder"
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/wow-look-at-my/testify/require"
)

func parseSourceFile(code string) *ast.SourceFile {
	return parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName:	"/test.ts",
		Path:		"/test.ts",
	}, code, core.ScriptKindTS)
}

func TestDecodeSourceFile_Basic(t *testing.T) {
	t.Parallel()
	sf := parseSourceFile("let x = 1;")
	buf, err := encoder.EncodeSourceFile(sf)
	require.NoError(t, err)

	decoded, err := encoder.DecodeSourceFile(buf)
	require.NoError(t, err)
	require.Equal(t, decoded.AsNode().Kind, ast.KindSourceFile)
	require.Equal(t, decoded.FileName(), "/test.ts")
	require.Equal(t, decoded.Text(), "let x = 1;")
	require.True(t, decoded.Statements != nil)
	require.True(t, decoded.EndOfFileToken != nil)
}

func TestDecodeSourceFile_Statements(t *testing.T) {
	t.Parallel()
	sf := parseSourceFile("let a = 1;\nlet b = 2;\nlet c = 3;")
	buf, err := encoder.EncodeSourceFile(sf)
	require.NoError(t, err)

	decoded, err := encoder.DecodeSourceFile(buf)
	require.NoError(t, err)
	require.Equal(t, len(decoded.Statements.Nodes), 3)
	for i, stmt := range decoded.Statements.Nodes {
		require.Equal(t, stmt.Kind, ast.KindVariableStatement, "statement %d", i)
	}
}

func TestDecodeSourceFile_VariableDeclaration(t *testing.T) {
	t.Parallel()
	sf := parseSourceFile("let x = 1;")
	buf, err := encoder.EncodeSourceFile(sf)
	require.NoError(t, err)

	decoded, err := encoder.DecodeSourceFile(buf)
	require.NoError(t, err)

	varStmt := decoded.Statements.Nodes[0].AsVariableStatement()
	require.True(t, varStmt.DeclarationList != nil)
	declList := varStmt.DeclarationList.AsVariableDeclarationList()
	require.True(t, declList.Declarations != nil)
	require.Equal(t, len(declList.Declarations.Nodes), 1)

	decl := declList.Declarations.Nodes[0].AsVariableDeclaration()
	require.Equal(t, decl.Name().Kind, ast.KindIdentifier)
	require.Equal(t, decl.Name().AsIdentifier().Text, "x")
	require.True(t, decl.Initializer != nil)
	require.Equal(t, decl.Initializer.Kind, ast.KindNumericLiteral)
	require.Equal(t, decl.Initializer.AsNumericLiteral().Text, "1")
}

func TestDecodeSourceFile_FunctionDeclaration(t *testing.T) {
	t.Parallel()
	sf := parseSourceFile("function add(a: number, b: number): number { return a + b; }")
	buf, err := encoder.EncodeSourceFile(sf)
	require.NoError(t, err)

	decoded, err := encoder.DecodeSourceFile(buf)
	require.NoError(t, err)

	funcDecl := decoded.Statements.Nodes[0].AsFunctionDeclaration()
	require.True(t, funcDecl.Name() != nil)
	require.Equal(t, funcDecl.Name().AsIdentifier().Text, "add")
	require.True(t, funcDecl.Parameters != nil)
	require.Equal(t, len(funcDecl.Parameters.Nodes), 2)
	require.True(t, funcDecl.Type != nil)
	require.True(t, funcDecl.Body != nil)

	param0 := funcDecl.Parameters.Nodes[0].AsParameterDeclaration()
	require.Equal(t, param0.Name().AsIdentifier().Text, "a")
	require.True(t, param0.Type != nil)
}

func TestDecodeSourceFile_ImportDeclaration(t *testing.T) {
	t.Parallel()
	sf := parseSourceFile(`import { bar } from "bar";`)
	buf, err := encoder.EncodeSourceFile(sf)
	require.NoError(t, err)

	decoded, err := encoder.DecodeSourceFile(buf)
	require.NoError(t, err)

	imp := decoded.Statements.Nodes[0].AsImportDeclaration()
	require.True(t, imp.ImportClause != nil)
	require.True(t, imp.ModuleSpecifier != nil)
	require.Equal(t, imp.ModuleSpecifier.AsStringLiteral().Text, "bar")

	clause := imp.ImportClause.AsImportClause()
	require.True(t, clause.NamedBindings != nil)
	namedImports := clause.NamedBindings.AsNamedImports()
	require.True(t, namedImports.Elements != nil)
	require.Equal(t, len(namedImports.Elements.Nodes), 1)
	spec := namedImports.Elements.Nodes[0].AsImportSpecifier()
	require.Equal(t, spec.Name().AsIdentifier().Text, "bar")
}

func TestDecodeSourceFile_IfStatement(t *testing.T) {
	t.Parallel()
	sf := parseSourceFile("if (true) { } else { }")
	buf, err := encoder.EncodeSourceFile(sf)
	require.NoError(t, err)

	decoded, err := encoder.DecodeSourceFile(buf)
	require.NoError(t, err)

	ifStmt := decoded.Statements.Nodes[0].AsIfStatement()
	require.True(t, ifStmt.Expression != nil)
	require.True(t, ifStmt.ThenStatement != nil)
	require.True(t, ifStmt.ElseStatement != nil)
	require.Equal(t, ifStmt.ThenStatement.Kind, ast.KindBlock)
	require.Equal(t, ifStmt.ElseStatement.Kind, ast.KindBlock)
}

func TestDecodeSourceFile_TemplateExpression(t *testing.T) {
	t.Parallel()
	sf := parseSourceFile("let x = `hello ${name} world`;")
	buf, err := encoder.EncodeSourceFile(sf)
	require.NoError(t, err)

	decoded, err := encoder.DecodeSourceFile(buf)
	require.NoError(t, err)

	varDecl := decoded.Statements.Nodes[0].AsVariableStatement().DeclarationList.AsVariableDeclarationList().Declarations.Nodes[0].AsVariableDeclaration()
	tmplExpr := varDecl.Initializer.AsTemplateExpression()
	require.True(t, tmplExpr.Head != nil)
	require.Equal(t, tmplExpr.Head.AsTemplateHead().Text, "hello ")
	require.True(t, tmplExpr.TemplateSpans != nil)
	require.Equal(t, len(tmplExpr.TemplateSpans.Nodes), 1)

	span := tmplExpr.TemplateSpans.Nodes[0].AsTemplateSpan()
	require.True(t, span.Expression != nil)
	require.Equal(t, span.Expression.Kind, ast.KindIdentifier)
	require.True(t, span.Literal != nil)
	require.Equal(t, span.Literal.AsTemplateTail().Text, " world")
}

func TestDecodeSourceFile_ExportModifier(t *testing.T) {
	t.Parallel()
	sf := parseSourceFile("export function foo() {}")
	buf, err := encoder.EncodeSourceFile(sf)
	require.NoError(t, err)

	decoded, err := encoder.DecodeSourceFile(buf)
	require.NoError(t, err)

	funcDecl := decoded.Statements.Nodes[0].AsFunctionDeclaration()
	require.True(t, funcDecl.Modifiers() != nil)
	require.Equal(t, len(funcDecl.Modifiers().Nodes), 1)
	require.Equal(t, funcDecl.Modifiers().Nodes[0].Kind, ast.KindExportKeyword)
}

func TestDecodeSourceFile_Positions(t *testing.T) {
	t.Parallel()
	code := "let x = 1;"
	sf := parseSourceFile(code)
	buf, err := encoder.EncodeSourceFile(sf)
	require.NoError(t, err)

	decoded, err := encoder.DecodeSourceFile(buf)
	require.NoError(t, err)

	require.Equal(t, decoded.AsNode().Pos(), 0)
	require.Equal(t, decoded.AsNode().End(), len(code))
}

func TestDecodeSourceFile_ClassDeclaration(t *testing.T) {
	t.Parallel()
	sf := parseSourceFile("class Foo { bar(): void {} }")
	buf, err := encoder.EncodeSourceFile(sf)
	require.NoError(t, err)

	decoded, err := encoder.DecodeSourceFile(buf)
	require.NoError(t, err)

	classDecl := decoded.Statements.Nodes[0].AsClassDeclaration()
	require.True(t, classDecl.Name() != nil)
	require.Equal(t, classDecl.Name().AsIdentifier().Text, "Foo")
	require.True(t, classDecl.Members != nil)
	require.Equal(t, len(classDecl.Members.Nodes), 1)
	require.Equal(t, classDecl.Members.Nodes[0].Kind, ast.KindMethodDeclaration)
}

func TestDecodeNodes_SubtreeRoundTrip(t *testing.T) {
	t.Parallel()
	sf := parseSourceFile("function greet(name: string) { return `Hello, ${name}!`; }")

	var funcNode *ast.Node
	visitor := &ast.NodeVisitor{}
	visitor.Visit = func(node *ast.Node) *ast.Node {
		if node.Kind == ast.KindFunctionDeclaration && funcNode == nil {
			funcNode = node
		}
		return node
	}
	visitor.VisitEachChild(sf.AsNode())
	require.True(t, funcNode != nil)

	buf, err := encoder.EncodeNode(funcNode, sf)
	require.NoError(t, err)

	decoded, err := encoder.DecodeNodes(buf)
	require.NoError(t, err)

	require.Equal(t, decoded.Kind, ast.KindFunctionDeclaration)
	funcDecl := decoded.AsFunctionDeclaration()
	require.True(t, funcDecl.Name() != nil)
	require.Equal(t, funcDecl.Name().AsIdentifier().Text, "greet")
	require.True(t, funcDecl.Parameters != nil)
	require.Equal(t, len(funcDecl.Parameters.Nodes), 1)
	require.True(t, funcDecl.Body != nil)
}

func TestDecodeSourceFile_BinaryExpression(t *testing.T) {
	t.Parallel()
	sf := parseSourceFile("let x = 1 + 2;")
	buf, err := encoder.EncodeSourceFile(sf)
	require.NoError(t, err)

	decoded, err := encoder.DecodeSourceFile(buf)
	require.NoError(t, err)

	decl := decoded.Statements.Nodes[0].AsVariableStatement().DeclarationList.AsVariableDeclarationList().Declarations.Nodes[0].AsVariableDeclaration()
	binExpr := decl.Initializer.AsBinaryExpression()
	require.True(t, binExpr.Left != nil)
	require.True(t, binExpr.Right != nil)
	require.True(t, binExpr.OperatorToken != nil)
	require.Equal(t, binExpr.Left.Kind, ast.KindNumericLiteral)
	require.Equal(t, binExpr.Right.Kind, ast.KindNumericLiteral)
}

func BenchmarkDecodeSourceFile(b *testing.B) {
	repo.SkipIfNoTypeScriptSubmodule(b)
	filePath := filepath.Join(repo.TypeScriptSubmodulePath(), "src/compiler/checker.ts")
	fileContent, err := os.ReadFile(filePath)
	require.NoError(b, err)
	code := string(fileContent)
	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName:	"/checker.ts",
		Path:		"/checker.ts",
	}, code, core.ScriptKindTS)

	buf, err := encoder.EncodeSourceFile(sourceFile)
	require.NoError(b, err)

	b.Run("parse", func(b *testing.B) {
		for b.Loop() {
			parser.ParseSourceFile(ast.SourceFileParseOptions{
				FileName:	"/checker.ts",
				Path:		"/checker.ts",
			}, code, core.ScriptKindTS)
		}
	})

	b.Run("decode", func(b *testing.B) {
		for b.Loop() {
			_, decodeErr := encoder.DecodeSourceFile(buf)
			require.NoError(b, decodeErr)
		}
	})
}
