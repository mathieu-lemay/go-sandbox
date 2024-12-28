package betteriter

import (
	"fmt"
	"math"
	"testing"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fake = faker.New()

func ptr[T any](v T) *T { return &v }

func TestNew_ReturnsAnIteratorOverTheValues(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	iter := New(values)

	idx := 0
	for v := range iter.it {
		assert.Equal(t, &values[idx], v)
		idx += 1
	}
}

func TestNew_DoesntAllocate(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	iter := New(values)

	idx := 0
	for v := range iter.it {
		assert.True(t, &values[idx] == v, "%p != %p", &values[idx], v)
		idx += 1
	}
}

func TestRepeat_ReturnsAnInfiniteIterator(t *testing.T) {
	val := 4 // Random value chosen by a fair dice roll

	repeat := Repeat(val)

	c := 0

	for v := range repeat.it {
		require.Equal(t, val, v)
		c += 1

		// Let's stop right there
		if c == 1_000 {
			break
		}
	}
}

func TestRepeatN_ReturnsAnInteratorOfNValues(t *testing.T) {
	val := 4 // Random value chosen by a fair dice roll
	n := fake.IntBetween(1, 10_000)

	repeat := RepeatN(val, n)

	c := 0

	for v := range repeat.it {
		require.Equal(t, val, v)
		c += 1

		// Safety stop
		if c == n+1 {
			require.Fail(t, "Iterator didn't stop", "It should have stopped at n=%d", n)
		}
	}

	assert.Equal(t, n, c)
}

func TestRepeatN_ReturnsNothingIfNIsZeroOrNegative(t *testing.T) {
	val := 4 // Random value chosen by a fair dice roll

	for n := range []int{0, -1, -2, math.MinInt} {
		t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
			repeat := RepeatN(val, 0)

			for range repeat.it {
				require.Fail(t, "The iterator returned something")
			}
		})
	}
}

func TestIncr_ReturnsAnIteratorThatIncreasesInValue(t *testing.T) {
	c := 0

	incr := Incr()

	for i := range incr.it {
		require.Equal(t, c, i)

		c += 1
		if c == 1_000 {
			break
		}
	}
}

func TestIncrN_ReturnsAnIteratorThatIncreasesByNInValue(t *testing.T) {
	c := 0
	n := fake.IntBetween(1, 10)

	incr := IncrN(n)

	for i := range incr.it {
		require.Equal(t, c*n, i)

		c += 1
		if c == 1_000 {
			break
		}
	}
}

func TestIncrFrom_ReturnsAnIteratorThatIncreasesInValueFromStart(t *testing.T) {
	for start := range []int{-100, -1, 0, 1, 100} {
		t.Run(fmt.Sprintf("start=%d", start), func(t *testing.T) {
			c := 0
			incr := IncrFrom(start)

			for i := range incr.it {
				require.Equal(t, start+c, i)

				c += 1
				if c == 1_000 {
					break
				}
			}
		})
	}
}

func TestIncrNFrom_ReturnsAnIteratorThatIncreasesInValueFromStart(t *testing.T) {
	n := fake.IntBetween(1, 10)

	for start := range []int{-100, -1, 0, 1, 100} {
		t.Run(fmt.Sprintf("start=%d", start), func(t *testing.T) {
			c := 0
			incr := IncrNFrom(start, n)

			for i := range incr.it {
				require.Equal(t, start+c*n, i)

				c += 1
				if c == 1_000 {
					break
				}
			}
		})
	}
}

func TestRange_ReturnsAnIteratorThatIncreasesInValueFromStartToEnd(t *testing.T) {
	start := fake.IntBetween(0, 10)
	end := fake.IntBetween(100, 1000)

	c := 0
	rangeIter := Range(start, end)

	for i := range rangeIter.it {
		require.Equal(t, start+c, i)

		c += 1

		if i >= end {
			require.Fail(
				t,
				"Range iter didn't stop at end value",
				"Expected end value: %d: current_value: %d",
				end,
				i,
			)
		}
	}
}

func TestCycle_IteratesOverMultipleSlices(t *testing.T) {
	sl1 := []int{0, 1, 2}
	sl2 := []int{5, 4, 3}
	sl3 := []int{7, 8, 6}

	expected := []int{0, 1, 2, 5, 4, 3, 7, 8, 6}

	idx := 0

	for v := range Chain(sl1, sl2, sl3).it {
		require.Equal(t, &expected[idx], v)

		idx += 1
	}

	// Make sure we covered all expected values
	assert.Len(t, expected, idx)
}

func TestProduct_ReturnsTheCartesianProductOfTwoIterables(t *testing.T) {
	p := []string{"A", "B", "C", "D"}
	q := []int{0, 1, 2}

	expected := []Tuple[string, int]{
		{"A", 0},
		{"A", 1},
		{"A", 2},
		{"B", 0},
		{"B", 1},
		{"B", 2},
		{"C", 0},
		{"C", 1},
		{"C", 2},
		{"D", 0},
		{"D", 1},
		{"D", 2},
	}

	idx := 0

	for v := range Product(p, q).it {
		require.Equal(t, expected[idx].A, *v.A)
		require.Equal(t, expected[idx].B, *v.B)

		idx += 1
	}

	// Make sure we covered all expected values
	assert.Len(t, expected, idx)
}

func TestZip_ReturnsValuesFromTwoSlices(t *testing.T) {
	iterA := []int{1, 2, 3, 4, 5}
	iterB := []string{"one", "two", "three", "four", "five"}

	idx := 0
	for v := range Zip(iterA, iterB).it {
		assert.Equal(t, iterA[idx], v.A)
		assert.Equal(t, iterB[idx], v.B)

		idx += 1
	}
}

func TestZip_StopsAtTheShortestOfTwoSlices(t *testing.T) {
	iterA := []int{1, 2, 3, 4, 5}
	iterB := []string{"one", "two", "three", "four", "five"}

	// Slice A is shorter
	idx := 0
	for range Zip(iterA[:3], iterB).it {
		idx += 1
	}
	assert.Equal(t, 3, idx)

	// Slice B is shorter
	idx = 0
	for range Zip(iterA, iterB[:3]).it {
		idx += 1
	}
	assert.Equal(t, 3, idx)
}

func TestZipEq_ReturnsValuesFromTwoSlices(t *testing.T) {
	iterA := []int{1, 2, 3, 4, 5}
	iterB := []string{"one", "two", "three", "four", "five"}

	idx := 0
	for v := range ZipEq(iterA, iterB).it {
		assert.Equal(t, iterA[idx], v.A)
		assert.Equal(t, iterB[idx], v.B)

		idx += 1
	}
}

func TestZipEq_ReturnsAnErrorIfSlicesAreNotSameLength(t *testing.T) {
	iterA := []int{1, 2, 3, 4, 5}
	iterB := []string{"one", "two", "three", "four", "five"}

	// Slice A is shorter
	output, err := ZipEq(iterA[:3], iterB).Collect()
	assert.ErrorContains(t, err, "slices are not the same length")
	assert.Empty(t, output)

	// Slice B is shorter
	output, err = ZipEq(iterA, iterB[:3]).Collect()
	assert.ErrorContains(t, err, "slices are not the same length")
	assert.Empty(t, output)
}
