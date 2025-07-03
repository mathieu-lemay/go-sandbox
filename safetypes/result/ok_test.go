package result

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOK_ReturnsNewOkResult(t *testing.T) {
	value := fake.Lorem().Word()

	res := Ok(value)

	expected := ok[string, error]{
		val: value,
	}
	assert.Equal(t, expected, res)
}

func TestOk_IsOk(t *testing.T) {
	o := Ok(fake.Int())

	assert.True(t, o.IsOk())
}

func TestOk_IsOkAnd(t *testing.T) {
	value := fake.Int()
	o := Ok(value)

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %v", res)
		t.Run(name, func(t *testing.T) {
			called := false

			f := func(v int) bool {
				called = true

				assert.Equal(t, value, v)

				return res
			}

			assert.Equal(t, res, o.IsOkAnd(f))
			assert.True(t, called, "predicate should have been called")
		})
	}
}

func TestOk_IsErr(t *testing.T) {
	o := Ok(fake.Int())

	assert.False(t, o.IsErr())
}

func TestOk_IsErrAnd(t *testing.T) {
	o := Ok(fake.Int())

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %v", res)
		t.Run(name, func(t *testing.T) {
			f := func(e error) bool {
				assert.Fail(t, "predicate should not have been called")

				return res
			}

			assert.False(t, o.IsErrAnd(f))
		})
	}
}

func TestOk_Expect(t *testing.T) {
	value := fake.Int()
	o := Ok(value)

	msg := fake.RandomStringWithLength(8)
	assert.Equal(t, value, o.Expect(msg))
}

func TestOk_ExpectErr(t *testing.T) {
	value := fake.Int()
	o := Ok(value)

	msg := fake.RandomStringWithLength(8)
	expectedError := fmt.Errorf("%s: %v", msg, value)
	assert.PanicsWithError(t, expectedError.Error(), func() {
		_ = o.ExpectErr(msg)
	})
}

func TestOk_Unwrap(t *testing.T) {
	value := fake.Int()
	s := Ok(value)

	assert.Equal(t, value, s.Unwrap())
}

func TestOk_UnwrapOr(t *testing.T) {
	value := fake.Int()
	s := Ok(value)

	def := fake.Int()

	assert.Equal(t, value, s.UnwrapOr(def))
}

func TestOk_UnwrapOrElse(t *testing.T) {
	value := fake.Int()
	s := Ok(value)

	def := fake.Int()
	f := func() int {
		assert.Fail(t, "should not be called")
		return def
	}

	assert.Equal(t, value, s.UnwrapOrElse(f))
}

func TestOk_UnwrapOrDefault(t *testing.T) {
	valStr := fake.RandomStringWithLength(8)
	valInt := fake.IntBetween(1, 100)
	valFloat := fake.Float64(2, 1, 100)

	assert.Equal(t, valStr, Ok(valStr).UnwrapOrDefault())
	assert.Equal(t, valInt, Ok(valInt).UnwrapOrDefault())
	assert.Equal(t, valFloat, Ok(valFloat).UnwrapOrDefault())
}

func TestOk_UnwrapErr(t *testing.T) {
	value := fake.Int()
	o := Ok(value)

	expected := fmt.Errorf("called `Result.UnwrapErr()` on an `Ok` value: %v", value)
	assert.PanicsWithError(t, expected.Error(), func() {
		_ = o.UnwrapErr()
	})
}
