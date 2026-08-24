package tspath

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/wow-look-at-my/testify/require"
)

func TestNormalizeSlashes(t *testing.T) {
	t.Parallel()
	require.Equal(t, NormalizeSlashes("a"), "a")
	require.Equal(t, NormalizeSlashes("a/b"), "a/b")
	require.Equal(t, NormalizeSlashes("a\\b"), "a/b")
	require.Equal(t, NormalizeSlashes("\\\\server\\path"), "//server/path")
}

func TestGetRootLength(t *testing.T) {
	t.Parallel()
	require.Equal(t, GetRootLength("a"), 0)
	require.Equal(t, GetRootLength("/"), 1)
	require.Equal(t, GetRootLength("/path"), 1)
	require.Equal(t, GetRootLength("c:"), 2)
	require.Equal(t, GetRootLength("c:d"), 0)
	require.Equal(t, GetRootLength("c:/"), 3)
	require.Equal(t, GetRootLength("c:\\"), 3)
	require.Equal(t, GetRootLength("//server"), 8)
	require.Equal(t, GetRootLength("//server/share"), 9)
	require.Equal(t, GetRootLength("\\\\server"), 8)
	require.Equal(t, GetRootLength("\\\\server\\share"), 9)
	require.Equal(t, GetRootLength("file:///"), 8)
	require.Equal(t, GetRootLength("file:///path"), 8)
	require.Equal(t, GetRootLength("file:///c:"), 10)
	require.Equal(t, GetRootLength("file:///c:d"), 8)
	require.Equal(t, GetRootLength("file:///c:/path"), 11)
	require.Equal(t, GetRootLength("file:///c%3a"), 12)
	require.Equal(t, GetRootLength("file:///c%3ad"), 8)
	require.Equal(t, GetRootLength("file:///c%3a/path"), 13)
	require.Equal(t, GetRootLength("file:///c%3A"), 12)
	require.Equal(t, GetRootLength("file:///c%3Ad"), 8)
	require.Equal(t, GetRootLength("file:///c%3A/path"), 13)
	require.Equal(t, GetRootLength("file://localhost"), 16)
	require.Equal(t, GetRootLength("file://localhost/"), 17)
	require.Equal(t, GetRootLength("file://localhost/path"), 17)
	require.Equal(t, GetRootLength("file://localhost/c:"), 19)
	require.Equal(t, GetRootLength("file://localhost/c:d"), 17)
	require.Equal(t, GetRootLength("file://localhost/c:/path"), 20)
	require.Equal(t, GetRootLength("file://localhost/c%3a"), 21)
	require.Equal(t, GetRootLength("file://localhost/c%3ad"), 17)
	require.Equal(t, GetRootLength("file://localhost/c%3a/path"), 22)
	require.Equal(t, GetRootLength("file://localhost/c%3A"), 21)
	require.Equal(t, GetRootLength("file://localhost/c%3Ad"), 17)
	require.Equal(t, GetRootLength("file://localhost/c%3A/path"), 22)
	require.Equal(t, GetRootLength("file://server"), 13)
	require.Equal(t, GetRootLength("file://server/"), 14)
	require.Equal(t, GetRootLength("file://server/path"), 14)
	require.Equal(t, GetRootLength("file://server/c:"), 14)
	require.Equal(t, GetRootLength("file://server/c:d"), 14)
	require.Equal(t, GetRootLength("file://server/c:/d"), 14)
	require.Equal(t, GetRootLength("file://server/c%3a"), 14)
	require.Equal(t, GetRootLength("file://server/c%3ad"), 14)
	require.Equal(t, GetRootLength("file://server/c%3a/d"), 14)
	require.Equal(t, GetRootLength("file://server/c%3A"), 14)
	require.Equal(t, GetRootLength("file://server/c%3Ad"), 14)
	require.Equal(t, GetRootLength("file://server/c%3A/d"), 14)
	require.Equal(t, GetRootLength("http://server"), 13)
	require.Equal(t, GetRootLength("http://server/path"), 14)
}

func TestPathIsAbsolute(t *testing.T) {
	t.Parallel()
	// POSIX
	require.Equal(t, PathIsAbsolute("/path/to/file.ext"), true)
	// DOS
	require.Equal(t, PathIsAbsolute("c:/path/to/file.ext"), true)
	// URL
	require.Equal(t, PathIsAbsolute("file:///path/to/file.ext"), true)
	// Non-absolute
	require.Equal(t, PathIsAbsolute("path/to/file.ext"), false)
	require.Equal(t, PathIsAbsolute("./path/to/file.ext"), false)
}

func TestIsUrl(t *testing.T) {
	t.Parallel()
	require.Equal(t, IsUrl("a"), false)
	require.Equal(t, IsUrl("/"), false)
	require.Equal(t, IsUrl("c:"), false)
	require.Equal(t, IsUrl("c:d"), false)
	require.Equal(t, IsUrl("c:/"), false)
	require.Equal(t, IsUrl("c:\\"), false)
	require.Equal(t, IsUrl("//server"), false)
	require.Equal(t, IsUrl("//server/share"), false)
	require.Equal(t, IsUrl("\\\\server"), false)
	require.Equal(t, IsUrl("\\\\server\\share"), false)

	require.Equal(t, IsUrl("file:///path"), true)
	require.Equal(t, IsUrl("file:///c:"), true)
	require.Equal(t, IsUrl("file:///c:d"), true)
	require.Equal(t, IsUrl("file:///c:/path"), true)
	require.Equal(t, IsUrl("file://server"), true)
	require.Equal(t, IsUrl("file://server/path"), true)
	require.Equal(t, IsUrl("http://server"), true)
	require.Equal(t, IsUrl("http://server/path"), true)
}

func TestIsRootedDiskPath(t *testing.T) {
	t.Parallel()
	require.Equal(t, IsRootedDiskPath("a"), false)
	require.Equal(t, IsRootedDiskPath("/"), true)
	require.Equal(t, IsRootedDiskPath("c:"), true)
	require.Equal(t, IsRootedDiskPath("c:d"), false)
	require.Equal(t, IsRootedDiskPath("c:/"), true)
	require.Equal(t, IsRootedDiskPath("c:\\"), true)
	require.Equal(t, IsRootedDiskPath("//server"), true)
	require.Equal(t, IsRootedDiskPath("//server/share"), true)
	require.Equal(t, IsRootedDiskPath("\\\\server"), true)
	require.Equal(t, IsRootedDiskPath("\\\\server\\share"), true)
	require.Equal(t, IsRootedDiskPath("file:///path"), false)
	require.Equal(t, IsRootedDiskPath("file:///c:"), false)
	require.Equal(t, IsRootedDiskPath("file:///c:d"), false)
	require.Equal(t, IsRootedDiskPath("file:///c:/path"), false)
	require.Equal(t, IsRootedDiskPath("file://server"), false)
	require.Equal(t, IsRootedDiskPath("file://server/path"), false)
	require.Equal(t, IsRootedDiskPath("http://server"), false)
	require.Equal(t, IsRootedDiskPath("http://server/path"), false)
}

func TestGetDirectoryPath(t *testing.T) {
	t.Parallel()
	require.Equal(t, GetDirectoryPath(""), "")
	require.Equal(t, GetDirectoryPath("a"), "")
	require.Equal(t, GetDirectoryPath("a/b"), "a")
	require.Equal(t, GetDirectoryPath("/"), "/")
	require.Equal(t, GetDirectoryPath("/a"), "/")
	require.Equal(t, GetDirectoryPath("/a/"), "/")
	require.Equal(t, GetDirectoryPath("/a/b"), "/a")
	require.Equal(t, GetDirectoryPath("/a/b/"), "/a")
	require.Equal(t, GetDirectoryPath("c:"), "c:")
	require.Equal(t, GetDirectoryPath("c:d"), "")
	require.Equal(t, GetDirectoryPath("c:/"), "c:/")
	require.Equal(t, GetDirectoryPath("c:/path"), "c:/")
	require.Equal(t, GetDirectoryPath("c:/path/"), "c:/")
	require.Equal(t, GetDirectoryPath("//server"), "//server")
	require.Equal(t, GetDirectoryPath("//server/"), "//server/")
	require.Equal(t, GetDirectoryPath("//server/share"), "//server/")
	require.Equal(t, GetDirectoryPath("//server/share/"), "//server/")
	require.Equal(t, GetDirectoryPath("\\\\server"), "//server")
	require.Equal(t, GetDirectoryPath("\\\\server\\"), "//server/")
	require.Equal(t, GetDirectoryPath("\\\\server\\share"), "//server/")
	require.Equal(t, GetDirectoryPath("\\\\server\\share\\"), "//server/")
	require.Equal(t, GetDirectoryPath("file:///"), "file:///")
	require.Equal(t, GetDirectoryPath("file:///path"), "file:///")
	require.Equal(t, GetDirectoryPath("file:///path/"), "file:///")
	require.Equal(t, GetDirectoryPath("file:///c:"), "file:///c:")
	require.Equal(t, GetDirectoryPath("file:///c:d"), "file:///")
	require.Equal(t, GetDirectoryPath("file:///c:/"), "file:///c:/")
	require.Equal(t, GetDirectoryPath("file:///c:/path"), "file:///c:/")
	require.Equal(t, GetDirectoryPath("file:///c:/path/"), "file:///c:/")
	require.Equal(t, GetDirectoryPath("file://server"), "file://server")
	require.Equal(t, GetDirectoryPath("file://server/"), "file://server/")
	require.Equal(t, GetDirectoryPath("file://server/path"), "file://server/")
	require.Equal(t, GetDirectoryPath("file://server/path/"), "file://server/")
	require.Equal(t, GetDirectoryPath("http://server"), "http://server")
	require.Equal(t, GetDirectoryPath("http://server/"), "http://server/")
	require.Equal(t, GetDirectoryPath("http://server/path"), "http://server/")
	require.Equal(t, GetDirectoryPath("http://server/path/"), "http://server/")
}

// !!!
// getBaseFileName
// getAnyExtensionFromPath

func TestGetPathComponents(t *testing.T) {
	t.Parallel()
	require.Equal(t, GetPathComponents("", ""), []string{""})
	require.Equal(t, GetPathComponents("a", ""), []string{"", "a"})
	require.Equal(t, GetPathComponents("./a", ""), []string{"", ".", "a"})
	require.Equal(t, GetPathComponents("/", ""), []string{"/"})
	require.Equal(t, GetPathComponents("/a", ""), []string{"/", "a"})
	require.Equal(t, GetPathComponents("/a/", ""), []string{"/", "a"})
	require.Equal(t, GetPathComponents("c:", ""), []string{"c:"})
	require.Equal(t, GetPathComponents("c:d", ""), []string{"", "c:d"})
	require.Equal(t, GetPathComponents("c:/", ""), []string{"c:/"})
	require.Equal(t, GetPathComponents("c:/path", ""), []string{"c:/", "path"})
	require.Equal(t, GetPathComponents("//server", ""), []string{"//server"})
	require.Equal(t, GetPathComponents("//server/", ""), []string{"//server/"})
	require.Equal(t, GetPathComponents("//server/share", ""), []string{"//server/", "share"})
	require.Equal(t, GetPathComponents("file:///", ""), []string{"file:///"})
	require.Equal(t, GetPathComponents("file:///path", ""), []string{"file:///", "path"})
	require.Equal(t, GetPathComponents("file:///c:", ""), []string{"file:///c:"})
	require.Equal(t, GetPathComponents("file:///c:d", ""), []string{"file:///", "c:d"})
	require.Equal(t, GetPathComponents("file:///c:/", ""), []string{"file:///c:/"})
	require.Equal(t, GetPathComponents("file:///c:/path", ""), []string{"file:///c:/", "path"})
	require.Equal(t, GetPathComponents("file://server", ""), []string{"file://server"})
	require.Equal(t, GetPathComponents("file://server/", ""), []string{"file://server/"})
	require.Equal(t, GetPathComponents("file://server/path", ""), []string{"file://server/", "path"})
	require.Equal(t, GetPathComponents("http://server", ""), []string{"http://server"})
	require.Equal(t, GetPathComponents("http://server/", ""), []string{"http://server/"})
	require.Equal(t, GetPathComponents("http://server/path", ""), []string{"http://server/", "path"})
}

func TestReducePathComponents(t *testing.T) {
	t.Parallel()
	require.Equal(t, reducePathComponents([]string{""}), []string{""})
	require.Equal(t, reducePathComponents([]string{"", "."}), []string{""})
	require.Equal(t, reducePathComponents([]string{"", ".", "a"}), []string{"", "a"})
	require.Equal(t, reducePathComponents([]string{"", "a", "."}), []string{"", "a"})
	require.Equal(t, reducePathComponents([]string{"", ".."}), []string{"", ".."})
	require.Equal(t, reducePathComponents([]string{"", "..", ".."}), []string{"", "..", ".."})
	require.Equal(t, reducePathComponents([]string{"", "..", ".", ".."}), []string{"", "..", ".."})
	require.Equal(t, reducePathComponents([]string{"", "a", ".."}), []string{""})
	require.Equal(t, reducePathComponents([]string{"", "..", "a"}), []string{"", "..", "a"})
	require.Equal(t, reducePathComponents([]string{"/"}), []string{"/"})
	require.Equal(t, reducePathComponents([]string{"/", "."}), []string{"/"})
	require.Equal(t, reducePathComponents([]string{"/", ".."}), []string{"/"})
	require.Equal(t, reducePathComponents([]string{"/", "a", ".."}), []string{"/"})
}

func TestCombinePaths(t *testing.T) {
	t.Parallel()
	// Non-rooted
	require.Equal(t, CombinePaths("path", "to", "file.ext"), "path/to/file.ext")
	require.Equal(t, CombinePaths("path", "dir", "..", "to", "file.ext"), "path/dir/../to/file.ext")
	// POSIX
	require.Equal(t, CombinePaths("/path", "to", "file.ext"), "/path/to/file.ext")
	require.Equal(t, CombinePaths("/path", "/to", "file.ext"), "/to/file.ext")
	// DOS
	require.Equal(t, CombinePaths("c:/path", "to", "file.ext"), "c:/path/to/file.ext")
	require.Equal(t, CombinePaths("c:/path", "c:/to", "file.ext"), "c:/to/file.ext")
	// URL
	require.Equal(t, CombinePaths("file:///path", "to", "file.ext"), "file:///path/to/file.ext")
	require.Equal(t, CombinePaths("file:///path", "file:///to", "file.ext"), "file:///to/file.ext")

	require.Equal(t, CombinePaths("/", "/node_modules/@types"), "/node_modules/@types")
	require.Equal(t, CombinePaths("/a/..", ""), "/a/..")
	require.Equal(t, CombinePaths("/a/..", "b"), "/a/../b")
	require.Equal(t, CombinePaths("/a/..", "b/"), "/a/../b/")
	require.Equal(t, CombinePaths("/a/..", "/"), "/")
	require.Equal(t, CombinePaths("/a/..", "/b"), "/b")
}

func BenchmarkCombinePaths(b *testing.B) {
	tests := [][]string{
		{"path", "to", "file.ext"},
		{"path", "dir", "..", "to", "file.ext"},
		{"/path", "to", "file.ext"},
		{"/path", "/to", "file.ext"},
		{"c:/path", "to", "file.ext"},
		{"file:///path", "to", "file.ext"},
	}

	for _, test := range tests {
		name := shortenName(strings.Join(test, "/"))
		b.Run(name, func(b *testing.B) {
			first, rest := test[0], test[1:]
			b.ReportAllocs()
			for b.Loop() {
				CombinePaths(first, rest...)
			}
		})
	}
}

func TestResolvePath(t *testing.T) {
	t.Parallel()
	require.Equal(t, ResolvePath(""), "")
	require.Equal(t, ResolvePath("."), "")
	require.Equal(t, ResolvePath("./"), "")
	require.Equal(t, ResolvePath(".."), "..")
	require.Equal(t, ResolvePath("../"), "../")
	require.Equal(t, ResolvePath("/"), "/")
	require.Equal(t, ResolvePath("/."), "/")
	require.Equal(t, ResolvePath("/./"), "/")
	require.Equal(t, ResolvePath("/../"), "/")
	require.Equal(t, ResolvePath("/a"), "/a")
	require.Equal(t, ResolvePath("/a/"), "/a/")
	require.Equal(t, ResolvePath("/a/."), "/a")
	require.Equal(t, ResolvePath("/a/./"), "/a/")
	require.Equal(t, ResolvePath("/a/./b"), "/a/b")
	require.Equal(t, ResolvePath("/a/./b/"), "/a/b/")
	require.Equal(t, ResolvePath("/a/.."), "/")
	require.Equal(t, ResolvePath("/a/../"), "/")
	require.Equal(t, ResolvePath("/a/../b"), "/b")
	require.Equal(t, ResolvePath("/a/../b/"), "/b/")
	require.Equal(t, ResolvePath("/a/..", "b"), "/b")
	require.Equal(t, ResolvePath("/a/..", "/"), "/")
	require.Equal(t, ResolvePath("/a/..", "b/"), "/b/")
	require.Equal(t, ResolvePath("/a/..", "/b"), "/b")
	require.Equal(t, ResolvePath("/a/.", "b"), "/a/b")
	require.Equal(t, ResolvePath("/a/.", "."), "/a")
	require.Equal(t, ResolvePath("a", "b", "c"), "a/b/c")
	require.Equal(t, ResolvePath("a", "b", "/c"), "/c")
	require.Equal(t, ResolvePath("a", "b", "../c"), "a/c")
}

func TestGetNormalizedAbsolutePath(t *testing.T) {
	t.Parallel()

	require.Equal(t, GetNormalizedAbsolutePath("/", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("/.", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("/./", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("/../", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("/a", ""), "/a")
	require.Equal(t, GetNormalizedAbsolutePath("/a/", ""), "/a")
	require.Equal(t, GetNormalizedAbsolutePath("/a/.", ""), "/a")
	require.Equal(t, GetNormalizedAbsolutePath("/a/foo.", ""), "/a/foo.")
	require.Equal(t, GetNormalizedAbsolutePath("/a/./", ""), "/a")
	require.Equal(t, GetNormalizedAbsolutePath("/a/./b", ""), "/a/b")
	require.Equal(t, GetNormalizedAbsolutePath("/a/./b/", ""), "/a/b")
	require.Equal(t, GetNormalizedAbsolutePath("/a/..", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("/a/../", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("/a/../", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("/a/../b", ""), "/b")
	require.Equal(t, GetNormalizedAbsolutePath("/a/../b/", ""), "/b")
	require.Equal(t, GetNormalizedAbsolutePath("/a/..", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("/a/..", "/"), "/")
	require.Equal(t, GetNormalizedAbsolutePath("/a/..", "b/"), "/")
	require.Equal(t, GetNormalizedAbsolutePath("/a/..", "/b"), "/")
	require.Equal(t, GetNormalizedAbsolutePath("/a/.", "b"), "/a")
	require.Equal(t, GetNormalizedAbsolutePath("/a/.", "."), "/a")

	// Tests as above, but with backslashes.
	require.Equal(t, GetNormalizedAbsolutePath("\\", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("\\.", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("\\.\\", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("\\..\\", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("\\a\\.\\", ""), "/a")
	require.Equal(t, GetNormalizedAbsolutePath("\\a\\.\\b", ""), "/a/b")
	require.Equal(t, GetNormalizedAbsolutePath("\\a\\.\\b\\", ""), "/a/b")
	require.Equal(t, GetNormalizedAbsolutePath("\\a\\..", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("\\a\\..\\", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("\\a\\..\\", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("\\a\\..\\b", ""), "/b")
	require.Equal(t, GetNormalizedAbsolutePath("\\a\\..\\b\\", ""), "/b")
	require.Equal(t, GetNormalizedAbsolutePath("\\a\\..", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("\\a\\..", "\\"), "/")
	require.Equal(t, GetNormalizedAbsolutePath("\\a\\..", "b\\"), "/")
	require.Equal(t, GetNormalizedAbsolutePath("\\a\\..", "\\b"), "/")
	require.Equal(t, GetNormalizedAbsolutePath("\\a\\.", "b"), "/a")
	require.Equal(t, GetNormalizedAbsolutePath("\\a\\.", "."), "/a")

	// Relative paths on an empty currentDirectory.
	require.Equal(t, GetNormalizedAbsolutePath("", ""), "")
	require.Equal(t, GetNormalizedAbsolutePath(".", ""), "")
	require.Equal(t, GetNormalizedAbsolutePath("./", ""), "")
	// Strangely, these do not normalize to the empty string.
	require.Equal(t, GetNormalizedAbsolutePath("..", ""), "..")
	require.Equal(t, GetNormalizedAbsolutePath("../", ""), "..")

	// Interaction between relative paths and currentDirectory.
	require.Equal(t, GetNormalizedAbsolutePath("", "/home"), "/home")
	require.Equal(t, GetNormalizedAbsolutePath(".", "/home"), "/home")
	require.Equal(t, GetNormalizedAbsolutePath("./", "/home"), "/home")
	require.Equal(t, GetNormalizedAbsolutePath("..", "/home"), "/")
	require.Equal(t, GetNormalizedAbsolutePath("../", "/home"), "/")
	require.Equal(t, GetNormalizedAbsolutePath("a", "b"), "b/a")
	require.Equal(t, GetNormalizedAbsolutePath("a", "b/c"), "b/c/a")

	// Base names starting or ending with a dot do not affect normalization.
	require.Equal(t, GetNormalizedAbsolutePath(".a", ""), ".a")
	require.Equal(t, GetNormalizedAbsolutePath("..a", ""), "..a")
	require.Equal(t, GetNormalizedAbsolutePath("a.", ""), "a.")
	require.Equal(t, GetNormalizedAbsolutePath("a..", ""), "a..")

	require.Equal(t, GetNormalizedAbsolutePath("/base/./.a", ""), "/base/.a")
	require.Equal(t, GetNormalizedAbsolutePath("/base/../.a", ""), "/.a")
	require.Equal(t, GetNormalizedAbsolutePath("/base/./..a", ""), "/base/..a")
	require.Equal(t, GetNormalizedAbsolutePath("/base/../..a", ""), "/..a")
	require.Equal(t, GetNormalizedAbsolutePath("/base/./..a/b", ""), "/base/..a/b")
	require.Equal(t, GetNormalizedAbsolutePath("/base/../..a/b", ""), "/..a/b")

	require.Equal(t, GetNormalizedAbsolutePath("/base/./a.", ""), "/base/a.")
	require.Equal(t, GetNormalizedAbsolutePath("/base/../a.", ""), "/a.")
	require.Equal(t, GetNormalizedAbsolutePath("/base/./a..", ""), "/base/a..")
	require.Equal(t, GetNormalizedAbsolutePath("/base/../a..", ""), "/a..")
	require.Equal(t, GetNormalizedAbsolutePath("/base/./a../b", ""), "/base/a../b")
	require.Equal(t, GetNormalizedAbsolutePath("/base/../a../b", ""), "/a../b")

	require.Equal(t, GetNormalizedAbsolutePath("a/..", ""), "")
	require.Equal(t, GetNormalizedAbsolutePath("/a//", ""), "/a")
	require.Equal(t, GetNormalizedAbsolutePath("//a", "a"), "//a/")
	require.Equal(t, GetNormalizedAbsolutePath("/\\", ""), "//")
	require.Equal(t, GetNormalizedAbsolutePath("a///", "a"), "a/a")
	require.Equal(t, GetNormalizedAbsolutePath("/.//", ""), "/")
	require.Equal(t, GetNormalizedAbsolutePath("//\\\\", ""), "///")
	require.Equal(t, GetNormalizedAbsolutePath(".//a", "."), "a")
	require.Equal(t, GetNormalizedAbsolutePath("a/../..", ""), "..")
	require.Equal(t, GetNormalizedAbsolutePath("../..", "\\a"), "/")
	require.Equal(t, GetNormalizedAbsolutePath("a:", "b"), "a:/")
	require.Equal(t, GetNormalizedAbsolutePath("a/../..", ".."), "../..")
	require.Equal(t, GetNormalizedAbsolutePath("a/../..", "b"), "")
	require.Equal(t, GetNormalizedAbsolutePath("a//../..", ".."), "../..")

	// Consecutive intermediate slashes are normalized to a single slash.
	require.Equal(t, GetNormalizedAbsolutePath("a//b", ""), "a/b")
	require.Equal(t, GetNormalizedAbsolutePath("a///b", ""), "a/b")
	require.Equal(t, GetNormalizedAbsolutePath("a/b//c", ""), "a/b/c")
	require.Equal(t, GetNormalizedAbsolutePath("/a/b//c", ""), "/a/b/c")
	require.Equal(t, GetNormalizedAbsolutePath("//a/b//c", ""), "//a/b/c")

	// Backslashes are converted to slashes,
	// and then consecutive intermediate slashes are normalized to a single slash
	require.Equal(t, GetNormalizedAbsolutePath("a\\\\b", ""), "a/b")
	require.Equal(t, GetNormalizedAbsolutePath("a\\\\\\b", ""), "a/b")
	require.Equal(t, GetNormalizedAbsolutePath("a\\b\\\\c", ""), "a/b/c")
	require.Equal(t, GetNormalizedAbsolutePath("\\a\\b\\\\c", ""), "/a/b/c")
	require.Equal(t, GetNormalizedAbsolutePath("\\\\a\\b\\\\c", ""), "//a/b/c")

	// The same occurs for mixed slashes.
	require.Equal(t, GetNormalizedAbsolutePath("a/\\b", ""), "a/b")
	require.Equal(t, GetNormalizedAbsolutePath("a\\/b", ""), "a/b")
	require.Equal(t, GetNormalizedAbsolutePath("a\\/\\b", ""), "a/b")
	require.Equal(t, GetNormalizedAbsolutePath("a\\b//c", ""), "a/b/c")
	require.Equal(t, GetNormalizedAbsolutePath("\\a\\b\\\\c", ""), "/a/b/c")
	require.Equal(t, GetNormalizedAbsolutePath("\\\\a\\b\\\\c", ""), "//a/b/c")
}

func TestGetNormalizedAbsolutePathWithoutRoot(t *testing.T) {
	t.Parallel()

	require.Equal(t, GetNormalizedAbsolutePathWithoutRoot("/a/b/c.txt", "/a/b"), "a/b/c.txt")
	require.Equal(t, GetNormalizedAbsolutePathWithoutRoot("c:/work/hello.txt", "c:/work"), "work/hello.txt")
	require.Equal(t, GetNormalizedAbsolutePathWithoutRoot("c:/work/hello.txt", "d:/worspaces"), "work/hello.txt")
}

var getNormalizedAbsolutePathTests = map[string][][]string{
	"non-normalized inputs": {
		{"/.", ""},
		{"/./", ""},
		{"/../", ""},
		{"/a/", ""},
		{"/a/.", ""},
		{"/a/foo.", ""},
		{"/a/./", ""},
		{"/a/./b", ""},
		{"/a/./b/", ""},
		{"/a/..", ""},
		{"/a/../", ""},
		{"/a/../", ""},
		{"/a/../b", ""},
		{"/a/../b/", ""},
		{"/a/..", ""},
		{"/a/..", "/"},
		{"/a/..", "b/"},
		{"/a/..", "/b"},
		{"/a/.", "b"},
		{"/a/.", "."},
	},
	"normalized inputs": {
		{"/a/b", ""},
		{"/one/two/three", ""},
		{"/users/root/project/src/foo.ts", ""},
	},
	"normalized inputs (long)": {
		{"/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t/u/v/w/x/y/z", ""},
		{"/one/two/three/four/five/six/seven/eight/nine/ten/eleven/twelve/thirteen/fourteen/fifteen/sixteen/seventeen/eighteen/nineteen/twenty", ""},
		{"/users/root/project/src/foo/bar/baz/qux/quux/corge/grault/garply/waldo/fred/plugh/xyzzy/thud", ""},
		{"/lorem/ipsum/dolor/sit/amet/consectetur/adipiscing/elit/sed/do/eiusmod/tempor/incididunt/ut/labore/et/dolore/magna/aliqua/ut/enim/ad/minim/veniam", ""},
	},
}

func BenchmarkGetNormalizedAbsolutePath(b *testing.B) {
	funcs := map[string]func(string, string) string{
		"GetNormalizedAbsolutePath":		GetNormalizedAbsolutePath,
		"GetNormalizedAbsolutePath (old)":	getNormalizedAbsolutePath_old,
	}
	for name, tests := range getNormalizedAbsolutePathTests {
		b.Run(name, func(b *testing.B) {
			for fnName, fn := range funcs {
				b.Run(fnName, func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						for _, test := range tests {
							fn(test[0], test[1])
						}
					}
				})
			}
		})
	}
}

func FuzzGetNormalizedAbsolutePath(f *testing.F) {
	for _, tests := range getNormalizedAbsolutePathTests {
		for _, test := range tests {
			f.Add(test[0], test[1])
		}
	}

	f.Fuzz(func(t *testing.T, p string, dir string) {
		require.Equal(t, GetNormalizedAbsolutePath(p, dir), getNormalizedAbsolutePath_old(p, dir), fmt.Sprintf("p=%q, dir=%q", p, dir))
	})
}

func TestGetRelativePathToDirectoryOrUrl(t *testing.T) {
	t.Parallel()
	// !!!
	// Based on tests for `getRelativePathFromDirectory`.

	require.Equal(t, GetRelativePathToDirectoryOrUrl("/", "/", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("/a", "/a", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("/a/", "/a", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("/a", "/", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "..")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("/a", "/b", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "../b")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("/a/b", "/b", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "../../b")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("/a/b/c", "/b", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "../../../b")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("/a/b/c", "/b/c", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "../../../b/c")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("/a/b/c", "/a/b", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "..")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("c:", "d:", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "d:/")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("file:///", "file:///", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("file:///a", "file:///a", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("file:///a/", "file:///a", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("file:///a", "file:///", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "..")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("file:///a", "file:///b", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "../b")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("file:///a/b", "file:///b", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "../../b")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("file:///a/b/c", "file:///b", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "../../../b")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("file:///a/b/c", "file:///b/c", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "../../../b/c")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("file:///a/b/c", "file:///a/b", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "..")
	require.Equal(t, GetRelativePathToDirectoryOrUrl("file:///c:", "file:///d:", false /*isAbsolutePathAnUrl*/, ComparePathsOptions{}), "file:///d:/")
}

func TestToFileNameLowerCase(t *testing.T) {
	t.Parallel()
	require.Equal(t, ToFileNameLowerCase("/user/UserName/projects/Project/file.ts"), "/user/username/projects/project/file.ts")
	require.Equal(t, ToFileNameLowerCase("/user/UserName/projects/projectß/file.ts"), "/user/username/projects/projectß/file.ts")
	require.Equal(t, ToFileNameLowerCase("/user/UserName/projects/İproject/file.ts"), "/user/username/projects/İproject/file.ts")
	require.Equal(t, ToFileNameLowerCase("/user/UserName/projects/ı/file.ts"), "/user/username/projects/ı/file.ts")
}

var toFileNameLowerCaseTests = []string{
	"/path/to/file.ext",
	"/PATH/TO/FILE.EXT",
	"/path/to/FILE.EXT",
	"/user/UserName/projects/Project/file.ts",
	"/user/UserName/projects/projectß/file.ts",
	"/user/UserName/projects/İproject/file.ts",
	"/user/UserName/projects/ı/file.ts",
	strings.Repeat("FoO/", 100),
}

// See [toFileNameLowerCase] for more info.
//
// To avoid having to do string building for most common cases, also ignore
// a-z, 0-9, \u0131, \u00DF, \, /, ., : and space
var fileNameLowerCaseRegExp = regexp.MustCompile(`[^\x{0130}\x{0131}\x{00DF}a-z0-9\\/:\-_. ]+`)

func oldToFileNameLowerCase(fileName string) string {
	return fileNameLowerCaseRegExp.ReplaceAllStringFunc(fileName, strings.ToLower)
}

func BenchmarkToFileNameLowerCase(b *testing.B) {
	for _, test := range toFileNameLowerCaseTests {
		name := shortenName(test)
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				ToFileNameLowerCase(test)
			}
		})
	}
}

func FuzzToFileNameLowerCase(f *testing.F) {
	for _, test := range toFileNameLowerCaseTests {
		f.Add(test)
	}

	f.Fuzz(func(t *testing.T, p string) {
		require.Equal(t, oldToFileNameLowerCase(p), ToFileNameLowerCase(p))
	})
}

func TestToPath(t *testing.T) {
	t.Parallel()
	require.Equal(t, string(ToPath("file.ext", "path/to", false /*useCaseSensitiveFileNames*/)), "path/to/file.ext")
	require.Equal(t, string(ToPath("file.ext", "/path/to", true /*useCaseSensitiveFileNames*/)), "/path/to/file.ext")
	require.Equal(t, string(ToPath("/path/to/../file.ext", "path/to", true /*useCaseSensitiveFileNames*/)), "/path/file.ext")
}

var relativePathSegmentRegExp = regexp.MustCompile(`//|(?:^|/)\.\.?(?:$|/)`)

func oldHasRelativePathSegment(p string) bool {
	return relativePathSegmentRegExp.MatchString(p)
}

var hasRelativePathSegmentTests = []struct {
	p	string
	bench	bool
}{
	{"//", false},
	{"foo/bar/baz", true},
	{"foo/./baz", false},
	{"foo/../baz", false},
	{"foo/bar/baz/.", false},
	{"./some/path", true},
	{"/foo//bar/", false},
	{"/foo/./bar/../../.", true},
	{strings.Repeat("foo/", 100) + "..", true},
}

func BenchmarkHasRelativePathSegment(b *testing.B) {
	for _, tt := range hasRelativePathSegmentTests {
		if !tt.bench {
			continue
		}
		name := shortenName(tt.p)
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				hasRelativePathSegment(tt.p)
			}
		})
	}
}

func FuzzHasRelativePathSegment(f *testing.F) {
	for _, tt := range hasRelativePathSegmentTests {
		f.Add(tt.p)
	}

	f.Fuzz(func(t *testing.T, p string) {
		require.Equal(t, oldHasRelativePathSegment(p), hasRelativePathSegment(p))
	})
}

var pathIsRelativeTests = []struct {
	p		string
	isRelative	bool
	benchmark	bool
}{
	// relative
	{".", true, false},
	{"..", true, false},
	{"./", true, false},
	{"../", true, false},
	{"./foo/bar", true, true},
	{"../foo/bar", true, true},
	{"../" + strings.Repeat("foo/", 100), true, true},
	// non-relative
	{"", false, false},
	{"foo", false, false},
	{"foo/bar", false, false},
	{"/foo/bar", false, false},
	{"c:/foo/bar", false, false},
}

func init() {
	old := pathIsRelativeTests

	for _, t := range old {
		t.p = strings.ReplaceAll(t.p, "/", "\\")
		pathIsRelativeTests = append(pathIsRelativeTests, t)
	}
}

func TestPathIsRelative(t *testing.T) {
	t.Parallel()
	for _, tt := range pathIsRelativeTests {
		name := shortenName(tt.p)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, PathIsRelative(tt.p), tt.isRelative)
		})
	}
}

func BenchmarkPathIsRelative(b *testing.B) {
	for _, tt := range pathIsRelativeTests {
		if !tt.benchmark {
			continue
		}
		name := shortenName(tt.p)
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				PathIsRelative(tt.p)
			}
		})
	}
}

func shortenName(name string) string {
	if len(name) > 20 {
		return name[:20] + "...etc"
	}
	return name
}

func normalizePath_old(path string) string {
	path = NormalizeSlashes(path)
	// Most paths don't require normalization
	if !hasRelativePathSegment(path) {
		return path
	}
	// Some paths only require cleanup of `/./` or leading `./`
	simplified := strings.ReplaceAll(path, "/./", "/")
	simplified = strings.TrimPrefix(simplified, "./")
	if simplified != path && !hasRelativePathSegment(simplified) {
		path = simplified
		return path
	}
	// Other paths require full normalization
	normalized := GetPathFromPathComponents(reducePathComponents(GetPathComponents(path, "")))
	if normalized != "" && HasTrailingDirectorySeparator(path) {
		normalized = EnsureTrailingDirectorySeparator(normalized)
	}
	return normalized
}

func getNormalizedAbsolutePath_old(fileName string, currentDirectory string) string {
	return GetPathFromPathComponents(GetNormalizedPathComponents(fileName, currentDirectory))
}

func TestGetCommonParents(t *testing.T) {
	t.Parallel()

	opts := ComparePathsOptions{}

	t.Run("empty input", func(t *testing.T) {
		t.Parallel()
		var paths []string
		got, ignored := GetCommonParents(paths, 1, GetPathComponents, opts)
		require.Equal(t, len(ignored), 0)
		require.Equal(t, got, ([]string)(nil))
	})

	t.Run("single path returns itself", func(t *testing.T) {
		t.Parallel()
		paths := []string{"/a/b/c/d"}
		got, ignored := GetCommonParents(paths, 1, GetPathComponents, opts)
		require.Equal(t, len(ignored), 0)
		expected := []string{paths[0]}
		require.Equal(t, got, expected)
	})

	t.Run("paths shorter than minComponents are ignored", func(t *testing.T) {
		t.Parallel()
		paths := []string{"/a/b/c/d", "/a/b/c/e", "/a/b/f/g", "/x/y"}
		got, ignored := GetCommonParents(paths, 4, GetPathComponents, opts)
		require.Equal(t, ignored, map[string]struct{}{"/x/y": {}})
		expected := []string{"/a/b/c", "/a/b/f/g"}
		require.Equal(t, got, expected)
	})

	t.Run("three paths share /a/b", func(t *testing.T) {
		t.Parallel()
		paths := []string{"/a/b/c/d", "/a/b/c/e", "/a/b/f/g"}
		got, ignored := GetCommonParents(paths, 1, GetPathComponents, opts)
		require.Equal(t, len(ignored), 0)
		expected := []string{"/a/b"}
		require.Equal(t, got, expected)
	})

	t.Run("mixed with short path collapses to root when minComponents=1", func(t *testing.T) {
		t.Parallel()
		paths := []string{"/a/b/c/d", "/a/b/c/e", "/a/b/f/g", "/x/y/z"}
		got, ignored := GetCommonParents(paths, 1, GetPathComponents, opts)
		require.Equal(t, len(ignored), 0)
		expected := []string{"/"}
		require.Equal(t, got, expected)
	})

	t.Run("mixed with short path preserves both when minComponents=3", func(t *testing.T) {
		t.Parallel()
		paths := []string{"/a/b/c/d", "/a/b/c/e", "/a/b/f/g", "/x/y/z"}
		got, ignored := GetCommonParents(paths, 3, GetPathComponents, opts)
		require.Equal(t, len(ignored), 0)
		expected := []string{"/a/b", "/x/y/z"}
		require.Equal(t, got, expected)
	})

	t.Run("different volumes are returned individually", func(t *testing.T) {
		t.Parallel()
		paths := []string{"c:/a/b/c/d", "d:/a/b/c/d"}
		got, ignored := GetCommonParents(paths, 1, GetPathComponents, opts)
		require.Equal(t, len(ignored), 0)
		expected := []string{paths[0], paths[1]}
		require.Equal(t, got, expected)
	})

	t.Run("duplicate paths deduplicate result", func(t *testing.T) {
		t.Parallel()
		paths := []string{"/a/b/c/d", "/a/b/c/d"}
		got, ignored := GetCommonParents(paths, 1, GetPathComponents, opts)
		require.Equal(t, len(ignored), 0)
		expected := []string{paths[0]}
		require.Equal(t, got, expected)
	})

	t.Run("paths with few components are returned as-is when minComponents met", func(t *testing.T) {
		t.Parallel()
		paths := []string{"/a/b/c/d", "/x/y"}
		got, ignored := GetCommonParents(paths, 2, GetPathComponents, opts)
		require.Equal(t, len(ignored), 0)
		expected := []string{"/a/b/c/d", "/x/y"}
		require.Equal(t, got, expected)
	})

	t.Run("minComponents=2", func(t *testing.T) {
		t.Parallel()
		paths := []string{"/a/b/c/d", "/a/z/c/e", "/a/aaa/f/g", "/x/y/z"}
		got, ignored := GetCommonParents(paths, 2, GetPathComponents, opts)
		require.Equal(t, len(ignored), 0)
		expected := []string{"/a", "/x/y/z"}
		require.Equal(t, got, expected)
	})

	t.Run("trailing separators are handled", func(t *testing.T) {
		t.Parallel()
		paths := []string{"/a/b/", "/a/b/c"}
		got, ignored := GetCommonParents(paths, 1, GetPathComponents, opts)
		require.Equal(t, len(ignored), 0)
		expected := []string{"/a/b"}
		require.Equal(t, got, expected)
	})
}
