package option

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNone_ReturnsNewNoneOption(t *testing.T) {
	res := None()

	n, ok := res.(none[any])
	require.True(t, ok, "result should be a none: %#v", res)

	assert.Equal(t, none[any]{}, n)
}

func TestNone_IsNone(t *testing.T) {
	n := None()

	assert.True(t, n.IsNone())
}

func TestNone_IsNoneOr(t *testing.T) {
	n := None()

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %v", res)
		t.Run(name, func(t *testing.T) {
			f := func(_ any) bool {
				// The predicate should _not_ be called, there's just no point in doing so.
				assert.Fail(t, "predicate should not have been called")

				return res
			}

			assert.True(t, n.IsNoneOr(f))
		})
	}
}

func TestNone_IsSome(t *testing.T) {
	n := None()

	assert.False(t, n.IsSome())
}

func TestNone_IsSomeAnd(t *testing.T) {
	n := None()

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %v", res)
		t.Run(name, func(t *testing.T) {
			f := func(_ any) bool {
				// The predicate should _not_ be called, there's just no point in doing so.
				assert.Fail(t, "predicate should not have been called")

				return res
			}

			assert.False(t, n.IsSomeAnd(f))
		})
	}
}

func TestNone_Expect(t *testing.T) {
	s := None()

	msg := fake.RandomStringWithLength(8)
	assert.PanicsWithError(t, msg, func() {
		s.Expect(msg)
	})
}

func TestNone_Unwrap(t *testing.T) {
	s := None()

	assert.PanicsWithError(t, "called `Option.Unwrap()` on a `None` value", func() {
		s.Unwrap()
	})
}

func TestNone_UnwrapOr(t *testing.T) {
	s := None()
	def := fake.RandomStringWithLength(8)

	assert.Equal(t, def, s.UnwrapOr(def))
}

func TestNone_UnwrapOrElse(t *testing.T) {
	s := None()

	def := fake.RandomStringWithLength(8)
	f := func() any {
		return def
	}

	assert.Equal(t, def, s.UnwrapOrElse(f))
}

func TestNone_UnwrapOrDefault(t *testing.T) {
	noneStr := none[string]{}
	noneInt := none[int]{}
	noneFloat := none[float64]{}

	assert.Equal(t, "", noneStr.UnwrapOrDefault())
	assert.Equal(t, 0, noneInt.UnwrapOrDefault())
	assert.Equal(t, 0.0, noneFloat.UnwrapOrDefault())
}
