package autoimport

import (
	"testing"

	"github.com/wow-look-at-my/testify/require"
)

type testEntry struct {
	name		string
	package_	string
}

func (e *testEntry) Name() string	{ return e.name }

func TestIndexClone(t *testing.T) {
	t.Parallel()

	t.Run("filters entries by package", func(t *testing.T) {
		t.Parallel()

		idx := &Index[*testEntry]{}
		idx.insertAsWords(&testEntry{name: "fooBar", package_: "pkg-a"})
		idx.insertAsWords(&testEntry{name: "bazQux", package_: "pkg-b"})
		idx.insertAsWords(&testEntry{name: "fooQux", package_: "pkg-a"})

		// Clone excluding pkg-b
		cloned := idx.Clone(func(e *testEntry) bool {
			return e.package_ != "pkg-b"
		})

		// Original should have all 3 entries
		require.Equal(t, len(idx.entries), 3)

		// Cloned should have 2 entries (only pkg-a)
		require.Equal(t, len(cloned.entries), 2)

		// Search should work on cloned index
		results := cloned.Find("fooBar", true)
		require.Equal(t, len(results), 1)
		require.Equal(t, results[0].name, "fooBar")

		// bazQux should not be in cloned index
		results = cloned.Find("bazQux", true)
		require.Equal(t, len(results), 0)

		// Word prefix search should work
		results = cloned.SearchWordPrefix("foo")
		require.Equal(t, len(results), 2)
	})

	t.Run("handles nil index", func(t *testing.T) {
		t.Parallel()

		var idx *Index[*testEntry]
		cloned := idx.Clone(func(e *testEntry) bool { return true })
		require.True(t, cloned == nil)
	})

	t.Run("handles empty index", func(t *testing.T) {
		t.Parallel()

		idx := &Index[*testEntry]{}
		cloned := idx.Clone(func(e *testEntry) bool { return true })
		require.Equal(t, len(cloned.entries), 0)
	})

	t.Run("filters all entries", func(t *testing.T) {
		t.Parallel()

		idx := &Index[*testEntry]{}
		idx.insertAsWords(&testEntry{name: "fooBar", package_: "pkg-a"})
		idx.insertAsWords(&testEntry{name: "bazQux", package_: "pkg-b"})

		cloned := idx.Clone(func(e *testEntry) bool { return false })
		require.Equal(t, len(cloned.entries), 0)
		require.Equal(t, len(cloned.index), 0)
	})
}
