package result

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOK_ReturnsNewOkResult(t *testing.T) {
	v := fake.Lorem().Word()

	res := Ok(v)

	expected := ok[string, error]{
		val: v,
	}
	assert.Equal(t, expected, res)
}

//	func TestOk_Expect(t *testing.T) {
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
//	func TestOk_ExpectErr(t *testing.T) {
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
