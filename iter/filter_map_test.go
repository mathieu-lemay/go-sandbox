package betteriter

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterMap_MapsAndFiltersElements(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	iter := New(values)

	filterMap := func(i int) (int, bool, error) {
		if i%2 != 0 {
			return 0, false, nil
		}

		return i * i, true, nil
	}

	output, err := FilterMap(iter, filterMap).Collect()

	require.NoError(t, err)

	assert.Equal(t, []int{4, 16}, output)
}

func TestFilterMap_IsLazy(t *testing.T) {
	values := []int{1, 2, 3}
	iter := New(values)

	filterMap := func(i int) (int, bool, error) {
		assert.LessOrEqualf(t, i, 2, "filter was called with unexpected value: %d", i)

		return i, true, nil
	}

	for v := range FilterMap(iter, filterMap).it {
		if v == 2 {
			break
		}
	}
}

func TestFilterMap_StopsOnError(t *testing.T) {
	values := []int{1, 2, 3}
	iter := New(values)

	filterMap := func(i int) (int, bool, error) {
		assert.LessOrEqualf(t, i, 2, "filter was called with unexpected value: %d", i)

		if i == 2 {
			return 0, false, errors.New("Invalid value")
		}

		return i, true, nil
	}

	output, err := FilterMap(iter, filterMap).Collect()
	assert.Empty(t, output)
	assert.ErrorContains(t, err, "Invalid value")
}

func TestFilterMap_PropagatesError(t *testing.T) {
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

	filterMap := func(i int) (int, bool, error) {
		// We should only be called with value = 1
		assert.Equal(t, 1, i, "Filter was called with unexpected value: %d", i)

		return i, true, nil
	}

	output, err := FilterMap(iter, filterMap).Collect()
	assert.Empty(t, output)
	assert.ErrorContains(t, err, "some error")
}
