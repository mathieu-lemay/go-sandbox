package betteriter

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFold_AppliesTheFolderFunctionOnAllValues(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}

	iter := New(values)

	res1, err := Fold(iter, 0, func(cur int, v *int) int { return cur + *v })
	require.NoError(t, err)

	assert.Equal(t, 15, res1)

	res2, err := Fold(
		iter,
		"result: ",
		func(cur string, v *int) string { return cur + strconv.Itoa(*v) },
	)
	require.NoError(t, err)

	assert.Equal(t, "result: 12345", res2)
}

func TestFold_PropagatesError(t *testing.T) {
	iter := &Iterator[int]{
		it: func(yield func(int, error) bool) {
			if !yield(1, nil) {
				return
			}
			if !yield(42, errors.New("some error")) {
				return
			}
			require.Fail(t, "Should not reach this point")
		},
	}

	output, err := Fold(iter, 0, func(_ int, v int) int { return v })
	assert.ErrorContains(t, err, "some error")
	// folding should have stopped at 1
	assert.Equal(t, 1, output)
}

func TestStringJoin_JoinsAllValuesInOneString(t *testing.T) {
	values := []string{"a", "b", "c", "d", "e"}

	testCases := []struct{ joiner, expected string }{
		{
			"",
			"abcde",
		},
		{
			",",
			"a,b,c,d,e",
		},
		{
			", ",
			"a, b, c, d, e",
		},
		{
			"-",
			"a-b-c-d-e",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.joiner, func(t *testing.T) {
			iter := New(values)

			output, err := StringJoin(iter, tc.joiner)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, output)

		})
	}

}

func TestStringJoin_PropagatesError(t *testing.T) {
	iter := &Iterator[*string]{
		it: func(yield func(*string, error) bool) {
			s := "a"
			if !yield(&s, nil) {
				return
			}
			s = "b"
			if !yield(&s, errors.New("some error")) {
				return
			}
			require.Fail(t, "Should not reach this point")
		},
	}

	output, err := StringJoin(iter, ",")
	assert.ErrorContains(t, err, "some error")
	assert.Empty(t, output)
}
