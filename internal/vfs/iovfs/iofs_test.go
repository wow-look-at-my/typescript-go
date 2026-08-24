package iovfs_test

import (
	"slices"
	"testing"
	"testing/fstest"

	"github.com/microsoft/typescript-go/internal/testutil"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/iovfs"
	"github.com/wow-look-at-my/testify/require"
)

func TestIOFS(t *testing.T) {
	t.Parallel()

	testfs := fstest.MapFS{
		"foo.ts": &fstest.MapFile{
			Data: []byte("hello, world"),
		},
		"dir1/file1.ts": &fstest.MapFile{
			Data: []byte("export const foo = 42;"),
		},
		"dir1/file2.ts": &fstest.MapFile{
			Data: []byte("export const foo = 42;"),
		},
		"dir2/file1.ts": &fstest.MapFile{
			Data: []byte("export const foo = 42;"),
		},
	}

	fs := iovfs.From(testfs, true)

	t.Run("ReadFile", func(t *testing.T) {
		t.Parallel()

		content, ok := fs.ReadFile("/foo.ts")
		require.True(t, ok)
		require.Equal(t, content, "hello, world")

		content, ok = fs.ReadFile("/does/not/exist.ts")
		require.True(t, !ok)
		require.Equal(t, content, "")
	})

	t.Run("ReadFileUnrooted", func(t *testing.T) {
		t.Parallel()

		testutil.AssertPanics(t, func() { fs.ReadFile("bar") }, `vfs: path "bar" is not absolute`)
	})

	t.Run("FileExists", func(t *testing.T) {
		t.Parallel()

		require.True(t, fs.FileExists("/foo.ts"))
		require.True(t, !fs.FileExists("/bar"))
	})

	t.Run("DirectoryExists", func(t *testing.T) {
		t.Parallel()

		require.True(t, fs.DirectoryExists("/"))
		require.True(t, fs.DirectoryExists("/dir1"))
		require.True(t, fs.DirectoryExists("/dir1/"))
		require.True(t, fs.DirectoryExists("/dir1/./"))
		require.True(t, !fs.DirectoryExists("/bar"))
	})

	t.Run("GetAccessibleEntries", func(t *testing.T) {
		t.Parallel()

		entries := fs.GetAccessibleEntries("/")
		require.Equal(t, entries.Directories, []string{"dir1", "dir2"})
		require.Equal(t, entries.Files, []string{"foo.ts"})
	})

	t.Run("WalkDir", func(t *testing.T) {
		t.Parallel()

		var files []string
		err := fs.WalkDir("/", func(path string, d vfs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		require.NoError(t, err)

		slices.Sort(files)

		require.Equal(t, files, []string{"/dir1/file1.ts", "/dir1/file2.ts", "/dir2/file1.ts", "/foo.ts"})
	})

	t.Run("WalkDirSkip", func(t *testing.T) {
		t.Parallel()

		var files []string
		err := fs.WalkDir("/", func(path string, d vfs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				files = append(files, path)
			}

			if path == "/" {
				return nil
			}

			return vfs.SkipDir
		})
		require.NoError(t, err)

		slices.Sort(files)

		require.Equal(t, files, []string{"/foo.ts"})
	})

	t.Run("Realpath", func(t *testing.T) {
		t.Parallel()

		realpath := fs.Realpath("/foo.ts")
		require.Equal(t, realpath, "/foo.ts")
	})

	t.Run("UseCaseSensitiveFileNames", func(t *testing.T) {
		t.Parallel()

		require.True(t, fs.UseCaseSensitiveFileNames())
	})
}
