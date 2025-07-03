package result

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErr_ReturnsNewErrResult(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))

	res := Err(err)

	expected := errT[any, error]{
		err: err,
	}
	assert.Equal(t, expected, res)
}

func TestErr_IsOk(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err(err)

	assert.False(t, e.IsOk())
}

func TestErr_IsOkAnd(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err(err)

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %v", res)
		t.Run(name, func(t *testing.T) {
			f := func(v any) bool {
				assert.Fail(t, "predicate should not have been called")

				return res
			}

			assert.False(t, e.IsOkAnd(f))
		})
	}
}

func TestErr_IsErr(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err(err)

	assert.True(t, e.IsErr())
}

func TestErr_IsErrAnd(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err(err)

	for _, res := range []bool{true, false} {
		name := fmt.Sprintf("predicate returns %v", res)
		t.Run(name, func(t *testing.T) {
			called := false

			f := func(e error) bool {
				called = true

				assert.Equal(t, err, e)

				return res
			}

			assert.Equal(t, res, e.IsErrAnd(f))
			assert.True(t, called, "predicate should have been called")
		})
	}
}

func TestErr_Expect(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err(err)

	msg := fake.RandomStringWithLength(8)
	expectedError := fmt.Errorf("%s: %w", msg, err)
	assert.PanicsWithError(t, expectedError.Error(), func() {
		e.Expect(msg)
	})
}

func TestErr_ExpectErr(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err(err)

	msg := fake.RandomStringWithLength(8)

	assert.Equal(t, err, e.ExpectErr(msg))
}

func TestErr_Unwrap(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err(err)

	expected := fmt.Errorf("called `Result.Unwrap()` on an `Err` value: %w", err)
	assert.PanicsWithError(t, expected.Error(), func() {
		e.Unwrap()
	})
}

func TestErr_UnwrapOr(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err(err)

	def := fake.RandomStringWithLength(8)

	assert.Equal(t, def, e.UnwrapOr(def))
}

func TestErr_UnwrapOrElse(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err(err)

	def := fake.RandomStringWithLength(8)
	f := func() any {
		return def
	}

	assert.Equal(t, def, e.UnwrapOrElse(f))
}

func TestErr_UnwrapOrDefault(t *testing.T) {
	noneStr := errT[string, error]{}
	noneInt := errT[int, error]{}
	noneFloat := errT[float64, error]{}

	assert.Equal(t, "", noneStr.UnwrapOrDefault())
	assert.Equal(t, 0, noneInt.UnwrapOrDefault())
	assert.Equal(t, 0.0, noneFloat.UnwrapOrDefault())
}

func TestErr_UnwrapErr(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
	e := Err(err)

	assert.Equal(t, err, e.UnwrapErr())
}
