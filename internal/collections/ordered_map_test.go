package collections_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/json"
	"github.com/wow-look-at-my/testify/require"
)

func TestOrderedMap(t *testing.T) {
	t.Parallel()

	var m collections.OrderedMap[int, string]

	require.True(t, !m.Has(1))

	const (
		N	= 1000
		start	= 1
		end	= start + N
	)

	// Seed the map with ascending keys and values for easier testing.
	for i := start; i < end; i++ {
		m.Set(i, padInt(i))
	}

	require.Equal(t, m.Size(), N)

	// Attempt to overwrite existing keys in reverse order.
	for i := end - 1; i >= start; i-- {
		m.Set(i, padInt(i))
	}

	require.Equal(t, m.Size(), N)

	for i := start; i < end; i++ {
		v, ok := m.Get(i)
		require.True(t, ok)
		require.Equal(t, v, padInt(i))
	}

	for k, v := range m.Entries() {
		require.Equal(t, v, padInt(k))
	}

	keys := slices.Collect(m.Keys())
	require.Equal(t, len(keys), N)
	require.True(t, slices.IsSorted(keys))

	values := slices.Collect(m.Values())
	require.Equal(t, len(values), N)
	require.True(t, slices.IsSorted(values))

	var firstKey int
	for k := range m.Keys() {
		firstKey = k
		break
	}
	require.Equal(t, firstKey, start)

	var firstValue string
	for v := range m.Values() {
		firstValue = v
		break
	}
	require.Equal(t, firstValue, padInt(start))

	for k, v := range m.Entries() {
		firstKey = k
		firstValue = v
		break
	}

	require.Equal(t, firstKey, start)
	require.Equal(t, firstValue, padInt(start))

	for i := start + 1; i < end; i++ {
		v, ok := m.Delete(i)
		require.True(t, ok)
		require.Equal(t, v, padInt(i))
		require.True(t, !m.Has(i))

		v, ok = m.Get(i)
		require.True(t, !ok)
		require.Equal(t, v, "")

		v, ok = m.Delete(i)
		require.True(t, !ok)
		require.Equal(t, v, "")
	}

	require.Equal(t, m.Size(), 1)
	require.True(t, m.Has(start))

	v, ok := m.Delete(start)
	require.True(t, ok)
	require.Equal(t, v, padInt(start))

	require.Equal(t, m.Size(), 0)
}

func TestOrderedMapClone(t *testing.T) {
	t.Parallel()

	m := &collections.OrderedMap[int, string]{}
	m.Set(1, "one")
	m.Set(2, "two")

	clone := m.Clone()

	require.True(t, clone != m)
	require.Equal(t, clone.Size(), 2)
	require.Equal(t, slices.Collect(clone.Keys()), []int{1, 2})
	require.Equal(t, slices.Collect(clone.Values()), []string{"one", "two"})

	v, ok := clone.Get(1)
	require.True(t, ok)
	require.Equal(t, v, "one")

	m.Delete(1)

	require.Equal(t, m.Size(), 1)
	require.Equal(t, clone.Size(), 2)
	require.Equal(t, slices.Collect(clone.Keys()), []int{1, 2})
	require.Equal(t, slices.Collect(clone.Values()), []string{"one", "two"})
}

func TestOrderedMapClear(t *testing.T) {
	t.Parallel()

	var m collections.OrderedMap[int, string]
	m.Set(1, "one")
	m.Set(2, "two")

	m.Clear()

	require.Equal(t, m.Size(), 0)
}

func padInt(n int) string {
	return fmt.Sprintf("%10d", n)
}

func TestOrderedMapWithSizeHint(t *testing.T) {	//nolint:paralleltest
	const N = 1024

	allocs := testing.AllocsPerRun(10, func() {
		m := collections.NewOrderedMapWithSizeHint[int, int](N)
		for i := range N {
			m.Set(i, i)
		}
	})

	require.True(t, allocs < 10, "allocs = %v", allocs)
}

func TestOrderedMapUnmarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("UnmarshalJSONV2", func(t *testing.T) {
		t.Parallel()
		testOrderedMapUnmarshalJSON(t, func(in []byte, out any) error { return json.Unmarshal(in, out) })
	})
}

func testOrderedMapUnmarshalJSON(t *testing.T, unmarshal func([]byte, any) error) {
	var m collections.OrderedMap[string, any]
	err := unmarshal([]byte(`{"a": 1, "b": "two", "c": { "d": 4 } }`), &m)
	require.NoError(t, err)

	require.Equal(t, m.Size(), 3)
	require.Equal(t, m.GetOrZero("a"), float64(1))

	err = unmarshal([]byte(`null`), &m)
	require.NoError(t, err)

	err = unmarshal([]byte(`"foo"`), &m)
	require.ErrorContains(t, err, "cannot unmarshal non-object JSON value into Map")

	var invalidMap collections.OrderedMap[int, any]
	err = unmarshal([]byte(`{"a": 1, "b": "two"}`), &invalidMap)
	require.ErrorContains(t, err, "unmarshal")
}
