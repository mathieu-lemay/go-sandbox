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

//	func TestErr_Expect(t *testing.T) {
//		t.Run("ok", func(t *testing.T) {
//			v := fake.RandomStringWithLength(8)
//			msg := fake.RandomStringWithLength(8)
//
//			res := Ok(v)
//
//			assert.Equal(t, v, res.Expect(msg))
//		})
//
//		t.Run("err", func(t *testing.T) {
//			err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
//			msg := fake.RandomStringWithLength(8)
//
//			res := Err(err)
//
//			expectedError := fmt.Errorf("%s: %w", msg, err)
//			assert.PanicsWithError(t, expectedError.Error(), func() {
//				res.Expect(msg)
//			})
//		})
//	}
//
//	func TestErr_ExpectErr(t *testing.T) {
//		t.Run("ok", func(t *testing.T) {
//			v := fake.RandomStringWithLength(8)
//			msg := fake.RandomStringWithLength(8)
//
//			res := Ok(v)
//
//			expectedError := fmt.Errorf("%s: %s", msg, v)
//			assert.PanicsWithError(t, expectedError.Error(), func() {
//				res.ExpectErr(msg)
//			})
//		})
//
//		t.Run("err", func(t *testing.T) {
//			err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
//			msg := fake.RandomStringWithLength(8)
//
//			res := Err(err)
//
//			assert.Equal(t, err, res.ExpectErr(msg))
//		})
//	}

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
