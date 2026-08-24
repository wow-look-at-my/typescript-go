package symlinks

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/wow-look-at-my/testify/assert"
	"github.com/wow-look-at-my/testify/require"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func TestNewKnownSymlink(t *testing.T) {
	t.Parallel()
	cache := NewKnownSymlink("/test/dir", true)
	require.NotNil(t, cache)
	assert.Equal(t, "/test/dir", cache.cwd)
	assert.True(t, cache.useCaseSensitiveFileNames)

}

func TestSetDirectory(t *testing.T) {
	t.Parallel()
	cache := NewKnownSymlink("/test/dir", true)
	symlinkPath := tspath.ToPath("/test/symlink", "/test/dir", true).EnsureTrailingDirectorySeparator()
	realDirectory := &KnownDirectoryLink{
		Real:		"/real/path/",
		RealPath:	tspath.ToPath("/real/path", "/test/dir", true).EnsureTrailingDirectorySeparator(),
	}

	cache.SetDirectory("/test/symlink", symlinkPath, realDirectory)

	// Check that directory was stored
	stored, ok := cache.Directories().Load(symlinkPath)
	require.True(t, ok)
	assert.Equal(t, realDirectory.Real, stored.Real)
	assert.Equal(t, realDirectory.RealPath, stored.RealPath)

	// Check that realpath mapping was created
	set, ok := cache.DirectoriesByRealpath().Load(realDirectory.RealPath)
	require.False(t, !ok || set.Size() == 0)
	assert.True(t, set.Has("/test/symlink"))

}

func TestSetFile(t *testing.T) {
	t.Parallel()
	cache := NewKnownSymlink("/test/dir", true)
	symlink := "/test/symlink/file.ts"
	symlinkPath := tspath.ToPath(symlink, "/test/dir", true)
	realpath := "/real/path/file.ts"

	cache.SetFile(symlink, symlinkPath, realpath)

	stored, ok := cache.Files().Load(symlinkPath)
	require.True(t, ok)
	assert.Equal(t, realpath, stored)

}

func TestProcessResolution(t *testing.T) {
	t.Parallel()
	cache := NewKnownSymlink("/test/dir", true)

	// Test with empty paths
	cache.ProcessResolution("", "")
	cache.ProcessResolution("original", "")
	cache.ProcessResolution("", "resolved")

	// Test with valid paths
	originalPath := "/test/original/file.ts"
	resolvedPath := "/test/resolved/file.ts"
	cache.ProcessResolution(originalPath, resolvedPath)

	// Check that file was stored
	symlinkPath := tspath.ToPath(originalPath, "/test/dir", true)
	stored, ok := cache.Files().Load(symlinkPath)
	require.True(t, ok)
	assert.Equal(t, resolvedPath, stored)

}

func TestGuessDirectorySymlink(t *testing.T) {
	t.Parallel()
	cache := NewKnownSymlink("/test/dir", true)

	tests := []struct {
		name		string
		a		string
		b		string
		cwd		string
		expected	[2]string	// [commonResolved, commonOriginal]
	}{
		{
			name:		"identical paths",
			a:		"/test/path/file.ts",
			b:		"/test/path/file.ts",
			cwd:		"/test/dir",
			expected:	[2]string{"/", "/"},
		},
		{
			name:		"different files same directory",
			a:		"/test/path/file1.ts",
			b:		"/test/path/file2.ts",
			cwd:		"/test/dir",
			expected:	[2]string{"", ""},
		},
		{
			name:		"different directories",
			a:		"/test/path1/file.ts",
			b:		"/test/path2/file.ts",
			cwd:		"/test/dir",
			expected:	[2]string{"/test/path1", "/test/path2"},
		},
		{
			name:		"node_modules paths",
			a:		"/test/node_modules/pkg/file.ts",
			b:		"/test/node_modules/pkg/file.ts",
			cwd:		"/test/dir",
			expected:	[2]string{"/test/node_modules/pkg", "/test/node_modules/pkg"},
		},
		{
			name:		"scoped package paths",
			a:		"/test/node_modules/@scope/pkg/file.ts",
			b:		"/test/node_modules/@scope/pkg/file.ts",
			cwd:		"/test/dir",
			expected:	[2]string{"/test/node_modules/@scope/pkg", "/test/node_modules/@scope/pkg"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			commonResolved, commonOriginal := cache.guessDirectorySymlink(tt.a, tt.b, tt.cwd)
			assert.Equal(t, tt.expected[0], commonResolved)
			assert.Equal(t, tt.expected[1], commonOriginal)

		})
	}
}

func TestIsNodeModulesOrScopedPackageDirectory(t *testing.T) {
	t.Parallel()
	cache := NewKnownSymlink("/test/dir", true)

	tests := []struct {
		name		string
		dir		string
		expected	bool
	}{
		{"node_modules", "node_modules", true},
		{"scoped package", "@scope", true},
		{"regular directory", "src", false},
		{"empty string", "", false},
		{"case insensitive node_modules", "NODE_MODULES", false},	// The function is case sensitive
		{"case insensitive scoped", "@SCOPE", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := cache.isNodeModulesOrScopedPackageDirectory(tt.dir)
			assert.Equal(t, tt.expected, result)

		})
	}
}

func TestSetSymlinksFromResolutions(t *testing.T) {
	t.Parallel()
	cache := NewKnownSymlink("/test/dir", true)

	// Mock resolution data
	resolvedModules := []struct {
		originalPath	string
		resolvedPath	string
		moduleName	string
		mode		core.ResolutionMode
		filePath	tspath.Path
	}{
		{
			originalPath:	"/test/original/file1.ts",
			resolvedPath:	"/test/resolved/file1.ts",
			moduleName:	"module1",
			mode:		core.ResolutionModeNone,
			filePath:	tspath.ToPath("/test/source.ts", "/test/dir", true),
		},
		{
			originalPath:	"/test/original/file2.ts",
			resolvedPath:	"/test/resolved/file2.ts",
			moduleName:	"module2",
			mode:		core.ResolutionModeNone,
			filePath:	tspath.ToPath("/test/source.ts", "/test/dir", true),
		},
	}

	// Mock callbacks
	forEachResolvedModule := func(callback func(resolution *module.ResolvedModule, moduleName string, mode core.ResolutionMode, filePath tspath.Path), file *ast.SourceFile) {
		for _, res := range resolvedModules {
			resolution := &module.ResolvedModule{
				OriginalPath:		res.originalPath,
				ResolvedFileName:	res.resolvedPath,
			}
			callback(resolution, res.moduleName, res.mode, res.filePath)
		}
	}

	forEachResolvedTypeReferenceDirective := func(callback func(resolution *module.ResolvedTypeReferenceDirective, moduleName string, mode core.ResolutionMode, filePath tspath.Path), file *ast.SourceFile) {
		// No type reference directives for this test
	}

	cache.SetSymlinksFromResolutions(forEachResolvedModule, forEachResolvedTypeReferenceDirective)

	// Check that files were stored
	for _, res := range resolvedModules {
		symlinkPath := tspath.ToPath(res.originalPath, "/test/dir", true)
		stored, ok := cache.Files().Load(symlinkPath)
		assert.True(t, ok)
		assert.Equal(t, res.resolvedPath, stored)

	}
}

func TestKnownSymlinksThreadSafety(t *testing.T) {
	t.Parallel()
	cache := NewKnownSymlink("/test/dir", true)

	// Test concurrent access
	done := make(chan bool, 10)

	for i := range 10 {
		go func(id int) {
			defer func() { done <- true }()

			symlinkPath := tspath.ToPath("/test/symlink"+string(rune(id)), "/test/dir", true).EnsureTrailingDirectorySeparator()
			realDirectory := &KnownDirectoryLink{
				Real:		"/real/path" + string(rune(id)) + "/",
				RealPath:	tspath.ToPath("/real/path"+string(rune(id)), "/test/dir", true).EnsureTrailingDirectorySeparator(),
			}

			cache.SetDirectory("/test/symlink"+string(rune(id)), symlinkPath, realDirectory)

			// Read back
			stored, ok := cache.Directories().Load(symlinkPath)
			assert.True(t, ok)
			assert.Equal(t, realDirectory.Real, stored.Real)

		}(i)
	}

	// Wait for all goroutines to complete
	for range 10 {
		<-done
	}
	assert.

	// Verify all directories were stored
	Equal(t, 10, cache.Directories().Size())

}
