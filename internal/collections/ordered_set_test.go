package collections_test

import (
	"slices"
	"testing"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/wow-look-at-my/testify/require"
)

func TestOrderedSet(t *testing.T) {
	t.Parallel()

	s := &collections.OrderedSet[int]{}

	s.Add(1)
	s.Add(2)
	s.Add(3)

	require.True(t, s.Has(1))
	require.True(t, s.Has(2))
	require.True(t, s.Has(3))

	require.True(t, s.Delete(2))

	values := slices.Collect(s.Values())
	require.Equal(t, len(values), 2)
	require.True(t, slices.IsSorted(values))

	s.Clear()

	require.Equal(t, s.Size(), 0)
	require.True(t, !s.Has(1))
	require.True(t, !s.Has(2))
	require.True(t, !s.Has(3))

	s2 := s.Clone()
	require.True(t, s != s2)
	require.Equal(t, s2.Size(), 0)
}

func TestOrderedSetWithSizeHint(t *testing.T) {	//nolint:paralleltest
	const N = 1024

	allocs := testing.AllocsPerRun(10, func() {
		m := collections.NewOrderedSetWithSizeHint[int](N)
		for i := range N {
			m.Add(i)
		}
	})

	require.True(t, allocs < 10, "allocs = %v", allocs)
}
