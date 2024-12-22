package betteriter

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilter_FiltersElements(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	iter := New(values)

	filter := func(i int) bool {
		return i%2 == 1
	}

	output, err := Filter(iter, filter).Collect()

	require.NoError(t, err)

	assert.Equal(t, []int{1, 3, 5}, output)
}

func TestFilter_IsLazy(t *testing.T) {
	values := []int{1, 2, 3}
	iter := New(values)

	filter := func(i int) bool {
		assert.LessOrEqualf(t, i, 2, "filter was called with unexpected value: %d", i)

		return true
	}

	for v := range Filter(iter, filter).it {
		if v == 2 {
			break
		}
	}
}

func TestFilter_PropagatesError(t *testing.T) {
	iter := &Iterator[int]{
		it: func(yield func(int, error) bool) {
			if !yield(1, nil) {
				return
			}
			if !yield(0, errors.New("some error")) {
				return
			}
			require.Fail(t, "Should not reach this point")
		},
	}

	filter := func(i int) bool {
		// We should only be called with value = 1
		assert.Equal(t, 1, i, "Filter was called with unexpected value: %d", i)

		return true
	}

	output, err := Filter(iter, filter).Collect()
	assert.Empty(t, output)
	assert.ErrorContains(t, err, "some error")
}

func Benchmark_Filter(b *testing.B) {
	type S struct {
		id  int
		i   int
		s   string
		b   bool
		i32 int32
		i64 int64
		u32 uint32
		u64 uint64
		f32 float32
		f64 float64
	}

	values := make([]S, 100)
	for i := range values {
		fake.Struct().Fill(&values[i])
		(&values[i]).id = i
	}

	iter := New(values)
	filter := Filter(iter, func(s S) bool { return s.id%2 == 0 })

	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		for range filter.it {

		}
	}
}
