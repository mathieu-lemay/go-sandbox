package option

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSome_ReturnsNewSomeOption(t *testing.T) {
	v := fake.Lorem().Word()

	res := Some(v)

	s, ok := res.(some[string])
	require.True(t, ok, "result should be a some[string]: %#v", res)

	expected := some[string]{
		val: v,
	}
	assert.Equal(t, expected, s)
}

func TestSome_IsNone(t *testing.T) {
	s := Some(fake.Int())

	assert.False(t, s.IsNone())
}

func TestSome_IsNoneOr(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %T", res)
		t.Run(name, func(t *testing.T) {
			called := false

			f := func(v int) bool {
				called = true

				assert.Equal(t, value, v)

				return res
			}

			assert.Equal(t, res, s.IsNoneOr(f))
			assert.True(t, called, "predicate should have been called")
		})
	}
}

func TestSome_IsSome(t *testing.T) {
	s := Some(fake.Int())

	assert.True(t, s.IsSome())
}

func TestSome_IsSomeAnd(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %T", res)
		t.Run(name, func(t *testing.T) {
			called := false

			f := func(v int) bool {
				called = true

				assert.Equal(t, value, v)

				return res
			}

			assert.Equal(t, res, s.IsSomeAnd(f))
			assert.True(t, called, "predicate should have been called")
		})
	}
}

func TestSome_Expect(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	assert.Equal(t, value, s.Expect(fake.RandomStringWithLength(8)))
}

func TestSome_Unwrap(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	assert.Equal(t, value, s.Unwrap())
}

func TestSome_UnwrapOr(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	def := fake.Int()

	assert.Equal(t, value, s.UnwrapOr(def))
}

func TestSome_UnwrapOrElse(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	def := fake.Int()
	f := func() int {
		assert.Fail(t, "should not be called")
		return def
	}

	assert.Equal(t, value, s.UnwrapOrElse(f))
}

func TestSome_UnwrapOrDefault(t *testing.T) {
	valStr := fake.RandomStringWithLength(8)
	valInt := fake.IntBetween(1, 100)
	valFloat := fake.Float64(2, 1, 100)

	assert.Equal(t, valStr, Some(valStr).UnwrapOrDefault())
	assert.Equal(t, valInt, Some(valInt).UnwrapOrDefault())
	assert.Equal(t, valFloat, Some(valFloat).UnwrapOrDefault())
}

func TestSome_Inspect(t *testing.T) {
	value := fake.Int()
	s := Some(value)

	called := false
	f := func(v int) {
		called = true
		assert.Equal(t, value, v)
	}

	assert.Equal(t, s, s.Inspect(f))
	assert.True(t, called, "predicate should have been called")
}
