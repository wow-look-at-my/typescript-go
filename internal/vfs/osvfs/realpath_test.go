package osvfs

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/wow-look-at-my/testify/require"
	"github.com/wow-look-at-my/testify/assert"
)

func TestSymlinkRealpath(t *testing.T) {
	t.Parallel()

	targetFile, linkFile := setupSymlinks(t)

	gotContents, err := os.ReadFile(linkFile)
	require.NoError(t, err)
	require.Equal(t, string(gotContents), "hello")

	fs := FS()

	targetRealpath := fs.Realpath(tspath.NormalizePath(targetFile))
	linkRealpath := fs.Realpath(tspath.NormalizePath(linkFile))

	if !assert.Equal(t, targetRealpath, linkRealpath) {
		cmd := exec.Command("node", "-e", `console.log({ native: fs.realpathSync.native(process.argv[1]), node: fs.realpathSync(process.argv[1]) })`, linkFile)
		out, err := cmd.CombinedOutput()
		require.NoError(t, err)
		t.Logf("node: %s", out)
	}
}

func setupSymlinks(tb testing.TB) (targetFile, linkFile string) {
	tb.Helper()

	tmp := tb.TempDir()

	target := filepath.Join(tmp, "target")
	targetFile = filepath.Join(target, "file")

	link := filepath.Join(tmp, "link")
	linkFile = filepath.Join(link, "file")

	require.NoError(tb, os.MkdirAll(target, 0o777))
	require.NoError(tb, os.WriteFile(targetFile, []byte("hello"), 0o666))

	mklink(tb, target, link, true)

	return targetFile, linkFile
}

func BenchmarkRealpath(b *testing.B) {
	targetFile, linkFile := setupSymlinks(b)

	fs := FS()
	normalizedTargetFile := tspath.NormalizePath(targetFile)
	normalizedLinkFile := tspath.NormalizePath(linkFile)

	b.Run("target", func(b *testing.B) {
		b.ReportAllocs()

		for b.Loop() {
			fs.Realpath(normalizedTargetFile)
		}
	})

	b.Run("link", func(b *testing.B) {
		b.ReportAllocs()

		for b.Loop() {
			fs.Realpath(normalizedLinkFile)
		}
	})
}

func TestGetAccessibleEntries(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	target := filepath.Join(tmp, "target")
	link := filepath.Join(tmp, "link")

	require.NoError(t, os.MkdirAll(target, 0o777))
	require.NoError(t, os.MkdirAll(link, 0o777))

	targetFile1 := filepath.Join(target, "file1")
	targetFile2 := filepath.Join(target, "file2")

	require.NoError(t, os.WriteFile(targetFile1, []byte("hello"), 0o666))
	require.NoError(t, os.WriteFile(targetFile2, []byte("world"), 0o666))

	targetDir1 := filepath.Join(target, "dir1")
	targetDir2 := filepath.Join(target, "dir2")

	require.NoError(t, os.MkdirAll(targetDir1, 0o777))
	require.NoError(t, os.MkdirAll(targetDir2, 0o777))

	mklink(t, targetFile1, filepath.Join(link, "file1"), false)
	mklink(t, targetFile2, filepath.Join(link, "file2"), false)
	mklink(t, targetDir1, filepath.Join(link, "dir1"), true)
	mklink(t, targetDir2, filepath.Join(link, "dir2"), true)

	fs := FS()

	entries := fs.GetAccessibleEntries(tspath.NormalizePath(link))

	require.Equal(t, entries.Directories, []string{"dir1", "dir2"})
	require.Equal(t, entries.Files, []string{"file1", "file2"})
}
