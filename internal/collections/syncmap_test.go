package collections_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/wow-look-at-my/testify/require"
)

func TestSyncMapWithNil(t *testing.T) {
	t.Parallel()

	var m collections.SyncMap[string, any]

	got1, ok := m.Load("foo")
	require.True(t, !ok)
	require.Equal(t, got1, nil)

	m.Store("foo", nil)

	got2, ok := m.Load("foo")
	require.True(t, ok)
	require.Equal(t, got2, nil)

	too, loaded := m.LoadOrStore("too", nil)
	require.True(t, !loaded)
	require.Equal(t, too, nil)

	m.Range(func(k string, v any) bool {
		return true
	})
}
