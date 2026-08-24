package bundled_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/osvfs"
	"github.com/wow-look-at-my/testify/require"
)

func TestTestingLibPath(t *testing.T) {
	t.Parallel()

	p := bundled.TestingLibPath()

	_, err := os.Stat(p)
	require.NoError(t, err)

	libdts := filepath.Join(p, "lib.d.ts")

	_, err = os.Stat(libdts)
	require.NoError(t, err)
}

func TestEmbeddedLibs(t *testing.T) {
	t.Parallel()

	fs := bundled.WrapFS(osvfs.FS())

	var files []string

	err := fs.WalkDir(bundled.LibPath(), func(path string, d vfs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, tspath.GetBaseFileName(path))
		}
		return nil
	})
	require.NoError(t, err)

	require.Equal(t, files, bundled.LibNames)
}
