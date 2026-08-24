package project

import (
	"testing"

	"github.com/wow-look-at-my/testify/require"
)

func TestGetPathComponentsForWatching(t *testing.T) {
	t.Parallel()

	require.Equal(t, getPathComponentsForWatching("/project", ""), []string{"/", "project"})
	require.Equal(t, getPathComponentsForWatching("C:\\project", ""), []string{"C:/", "project"})
	require.Equal(t, getPathComponentsForWatching("//server/share/project/tsconfig.json", ""), []string{"//server/share", "project", "tsconfig.json"})
	require.Equal(t, getPathComponentsForWatching(`\\server\share\project\tsconfig.json`, ""), []string{"//server/share", "project", "tsconfig.json"})
	require.Equal(t, getPathComponentsForWatching("C:\\Users", ""), []string{"C:/Users"})
	require.Equal(t, getPathComponentsForWatching("C:\\Users\\andrew\\project", ""), []string{"C:/Users/andrew", "project"})
	require.Equal(t, getPathComponentsForWatching("/home", ""), []string{"/home"})
	require.Equal(t, getPathComponentsForWatching("/home/andrew/project", ""), []string{"/home/andrew", "project"})
}

func TestNilWatchedFilesClone(t *testing.T) {
	t.Parallel()

	var w *WatchedFiles[int]
	result := w.Clone(42)
	require.True(t, result == nil, "clone on a nil `WatchedFiles` should return nil")
}
