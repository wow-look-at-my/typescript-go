package cachedvfs_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/cachedvfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfsmock"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"github.com/wow-look-at-my/testify/require"
)

func createMockFS() *vfsmock.FSMock {
	return vfsmock.Wrap(vfstest.FromMap(map[string]string{
		"/some/path/file.txt": "hello world",
	}, true))
}

func TestDirectoryExists(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	cached.DirectoryExists("/some/path")
	require.Equal(t, 1, len(underlying.DirectoryExistsCalls()))

	cached.DirectoryExists("/some/path")
	require.Equal(t, 1, len(underlying.DirectoryExistsCalls()))

	cached.ClearCache()
	cached.DirectoryExists("/some/path")
	require.Equal(t, 2, len(underlying.DirectoryExistsCalls()))

	cached.DirectoryExists("/other/path")
	require.Equal(t, 3, len(underlying.DirectoryExistsCalls()))

	cached.DisableAndClearCache()
	cached.DirectoryExists("/some/path")
	require.Equal(t, 4, len(underlying.DirectoryExistsCalls()))

	cached.DirectoryExists("/some/path")
	require.Equal(t, 5, len(underlying.DirectoryExistsCalls()))

	cached.Enable()
	cached.DirectoryExists("/some/path")
	require.Equal(t, 6, len(underlying.DirectoryExistsCalls()))

	cached.DirectoryExists("/some/path")
	require.Equal(t, 6, len(underlying.DirectoryExistsCalls()))
}

func TestFileExists(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	cached.FileExists("/some/path/file.txt")
	require.Equal(t, 1, len(underlying.FileExistsCalls()))

	cached.FileExists("/some/path/file.txt")
	require.Equal(t, 1, len(underlying.FileExistsCalls()))

	cached.ClearCache()
	cached.FileExists("/some/path/file.txt")
	require.Equal(t, 2, len(underlying.FileExistsCalls()))

	cached.FileExists("/other/path/file.txt")
	require.Equal(t, 3, len(underlying.FileExistsCalls()))

	cached.DisableAndClearCache()
	cached.FileExists("/some/path/file.txt")
	require.Equal(t, 4, len(underlying.FileExistsCalls()))

	cached.FileExists("/some/path/file.txt")
	require.Equal(t, 5, len(underlying.FileExistsCalls()))

	cached.Enable()
	cached.FileExists("/some/path/file.txt")
	require.Equal(t, 6, len(underlying.FileExistsCalls()))

	cached.FileExists("/some/path/file.txt")
	require.Equal(t, 6, len(underlying.FileExistsCalls()))
}

func TestGetAccessibleEntries(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	cached.GetAccessibleEntries("/some/path")
	require.Equal(t, 1, len(underlying.GetAccessibleEntriesCalls()))

	cached.GetAccessibleEntries("/some/path")
	require.Equal(t, 1, len(underlying.GetAccessibleEntriesCalls()))

	cached.ClearCache()
	cached.GetAccessibleEntries("/some/path")
	require.Equal(t, 2, len(underlying.GetAccessibleEntriesCalls()))

	cached.GetAccessibleEntries("/other/path")
	require.Equal(t, 3, len(underlying.GetAccessibleEntriesCalls()))

	cached.DisableAndClearCache()
	cached.GetAccessibleEntries("/some/path")
	require.Equal(t, 4, len(underlying.GetAccessibleEntriesCalls()))

	cached.GetAccessibleEntries("/some/path")
	require.Equal(t, 5, len(underlying.GetAccessibleEntriesCalls()))

	cached.Enable()
	cached.GetAccessibleEntries("/some/path")
	require.Equal(t, 6, len(underlying.GetAccessibleEntriesCalls()))

	cached.GetAccessibleEntries("/some/path")
	require.Equal(t, 6, len(underlying.GetAccessibleEntriesCalls()))
}

func TestRealpath(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	cached.Realpath("/some/path")
	require.Equal(t, 1, len(underlying.RealpathCalls()))

	cached.Realpath("/some/path")
	require.Equal(t, 1, len(underlying.RealpathCalls()))

	cached.ClearCache()
	cached.Realpath("/some/path")
	require.Equal(t, 2, len(underlying.RealpathCalls()))

	cached.Realpath("/other/path")
	require.Equal(t, 3, len(underlying.RealpathCalls()))

	cached.DisableAndClearCache()
	cached.Realpath("/some/path")
	require.Equal(t, 4, len(underlying.RealpathCalls()))

	cached.Realpath("/some/path")
	require.Equal(t, 5, len(underlying.RealpathCalls()))

	cached.Enable()
	cached.Realpath("/some/path")
	require.Equal(t, 6, len(underlying.RealpathCalls()))

	cached.Realpath("/some/path")
	require.Equal(t, 6, len(underlying.RealpathCalls()))
}

func TestStat(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	cached.Stat("/some/path")
	require.Equal(t, 1, len(underlying.StatCalls()))

	cached.Stat("/some/path")
	require.Equal(t, 1, len(underlying.StatCalls()))

	cached.ClearCache()
	cached.Stat("/some/path")
	require.Equal(t, 2, len(underlying.StatCalls()))

	cached.Stat("/other/path")
	require.Equal(t, 3, len(underlying.StatCalls()))

	cached.DisableAndClearCache()
	cached.Stat("/some/path")
	require.Equal(t, 4, len(underlying.StatCalls()))

	cached.Stat("/some/path")
	require.Equal(t, 5, len(underlying.StatCalls()))

	cached.Enable()
	cached.Stat("/some/path")
	require.Equal(t, 6, len(underlying.StatCalls()))

	cached.Stat("/some/path")
	require.Equal(t, 6, len(underlying.StatCalls()))
}

func TestReadFile(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	cached.ReadFile("/some/path/file.txt")
	require.Equal(t, 1, len(underlying.ReadFileCalls()))

	cached.ReadFile("/some/path/file.txt")
	require.Equal(t, 2, len(underlying.ReadFileCalls()))

	cached.ClearCache()
	cached.ReadFile("/some/path/file.txt")
	require.Equal(t, 3, len(underlying.ReadFileCalls()))

	cached.DisableAndClearCache()
	cached.ReadFile("/some/path/file.txt")
	require.Equal(t, 4, len(underlying.ReadFileCalls()))

	cached.ReadFile("/some/path/file.txt")
	require.Equal(t, 5, len(underlying.ReadFileCalls()))

	cached.Enable()
	cached.ReadFile("/some/path/file.txt")
	require.Equal(t, 6, len(underlying.ReadFileCalls()))

	cached.ReadFile("/some/path/file.txt")
	require.Equal(t, 7, len(underlying.ReadFileCalls()))
}

func TestUseCaseSensitiveFileNames(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	cached.UseCaseSensitiveFileNames()
	require.Equal(t, 1, len(underlying.UseCaseSensitiveFileNamesCalls()))

	cached.UseCaseSensitiveFileNames()
	require.Equal(t, 2, len(underlying.UseCaseSensitiveFileNamesCalls()))

	cached.ClearCache()
	cached.UseCaseSensitiveFileNames()
	require.Equal(t, 3, len(underlying.UseCaseSensitiveFileNamesCalls()))

	cached.DisableAndClearCache()
	cached.UseCaseSensitiveFileNames()
	require.Equal(t, 4, len(underlying.UseCaseSensitiveFileNamesCalls()))

	cached.UseCaseSensitiveFileNames()
	require.Equal(t, 5, len(underlying.UseCaseSensitiveFileNamesCalls()))

	cached.Enable()
	cached.UseCaseSensitiveFileNames()
	require.Equal(t, 6, len(underlying.UseCaseSensitiveFileNamesCalls()))

	cached.UseCaseSensitiveFileNames()
	require.Equal(t, 7, len(underlying.UseCaseSensitiveFileNamesCalls()))
}

func TestWalkDir(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	walkFn := vfs.WalkDirFunc(func(path string, info vfs.DirEntry, err error) error {
		return nil
	})

	_ = cached.WalkDir("/some/path", walkFn)
	require.Equal(t, 1, len(underlying.WalkDirCalls()))

	_ = cached.WalkDir("/some/path", walkFn)
	require.Equal(t, 2, len(underlying.WalkDirCalls()))

	cached.ClearCache()
	_ = cached.WalkDir("/some/path", walkFn)
	require.Equal(t, 3, len(underlying.WalkDirCalls()))

	cached.DisableAndClearCache()
	_ = cached.WalkDir("/some/path", walkFn)
	require.Equal(t, 4, len(underlying.WalkDirCalls()))

	_ = cached.WalkDir("/some/path", walkFn)
	require.Equal(t, 5, len(underlying.WalkDirCalls()))

	cached.Enable()
	_ = cached.WalkDir("/some/path", walkFn)
	require.Equal(t, 6, len(underlying.WalkDirCalls()))

	_ = cached.WalkDir("/some/path", walkFn)
	require.Equal(t, 7, len(underlying.WalkDirCalls()))
}

func TestRemove(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	_ = cached.Remove("/some/path/file.txt")
	require.Equal(t, 1, len(underlying.RemoveCalls()))

	_ = cached.Remove("/some/path/file.txt")
	require.Equal(t, 2, len(underlying.RemoveCalls()))

	cached.ClearCache()
	_ = cached.Remove("/some/path/file.txt")
	require.Equal(t, 3, len(underlying.RemoveCalls()))

	cached.DisableAndClearCache()
	_ = cached.Remove("/some/path/file.txt")
	require.Equal(t, 4, len(underlying.RemoveCalls()))

	_ = cached.Remove("/some/path/file.txt")
	require.Equal(t, 5, len(underlying.RemoveCalls()))

	cached.Enable()
	_ = cached.Remove("/some/path/file.txt")
	require.Equal(t, 6, len(underlying.RemoveCalls()))

	_ = cached.Remove("/some/path/file.txt")
	require.Equal(t, 7, len(underlying.RemoveCalls()))
}

func TestWriteFile(t *testing.T) {
	t.Parallel()

	underlying := createMockFS()
	cached := cachedvfs.From(underlying)

	_ = cached.WriteFile("/some/path/file.txt", "new content", false)
	require.Equal(t, 1, len(underlying.WriteFileCalls()))

	_ = cached.WriteFile("/some/path/file.txt", "another content", true)
	require.Equal(t, 2, len(underlying.WriteFileCalls()))

	cached.ClearCache()
	_ = cached.WriteFile("/some/path/file.txt", "third content", false)
	require.Equal(t, 3, len(underlying.WriteFileCalls()))

	call := underlying.WriteFileCalls()[2]
	require.Equal(t, "/some/path/file.txt", call.Path)
	require.Equal(t, "third content", call.Data)
	require.Equal(t, false, call.WriteByteOrderMark)

	cached.DisableAndClearCache()
	_ = cached.WriteFile("/some/path/file.txt", "fourth content", false)
	require.Equal(t, 4, len(underlying.WriteFileCalls()))

	_ = cached.WriteFile("/some/path/file.txt", "fifth content", true)
	require.Equal(t, 5, len(underlying.WriteFileCalls()))

	cached.Enable()
	_ = cached.WriteFile("/some/path/file.txt", "sixth content", false)
	require.Equal(t, 6, len(underlying.WriteFileCalls()))

	_ = cached.WriteFile("/some/path/file.txt", "seventh content", true)
	require.Equal(t, 7, len(underlying.WriteFileCalls()))
}
