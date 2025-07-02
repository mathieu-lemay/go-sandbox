package result

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
)

var fake = faker.NewWithSeed(rand.NewSource(time.Now().UnixNano()))

func TestFrom_ReturnsNewResultFromArgs(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		v := fake.RandomStringWithLength(8)

		f := func() (string, error) {
			return v, nil
		}

		res := From(f())

		assert.True(t, res.IsOk())
		assert.Equal(t, v, res.val)
		assert.Nil(t, res.err)
	})

	t.Run("err", func(t *testing.T) {
		v := fake.RandomStringWithLength(8)
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))

		f := func() (string, error) {
			return v, err
		}

		res := From(f())

		assert.False(t, res.IsOk())
		assert.Zero(t, res.val)
		assert.Equal(t, err, res.err)
	})
}

func TestOK_ReturnsNewOkResult(t *testing.T) {
	v := fake.Lorem().Word()

	res := Ok(v)

	expected := Result[string]{
		val: v,
		err: nil,
	}
	assert.Equal(t, expected, res)
}

func TestErr_ReturnsNewErrResult(t *testing.T) {
	err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))

	res := Err(err)

	expected := Result[any]{
		val: nil,
		err: err,
	}
	assert.Equal(t, expected, res)
}

func TestExpect(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		v := fake.RandomStringWithLength(8)
		msg := fake.RandomStringWithLength(8)

		res := Ok(v)

		assert.Equal(t, v, res.Expect(msg))
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		msg := fake.RandomStringWithLength(8)

		res := Err(err)

		expectedError := fmt.Errorf("%s: %w", msg, err)
		assert.PanicsWithError(t, expectedError.Error(), func() {
			res.Expect(msg)
		})
	})
}

func TestExpectErr(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		v := fake.RandomStringWithLength(8)
		msg := fake.RandomStringWithLength(8)

		res := Ok(v)

		expectedError := fmt.Errorf("%s: %s", msg, v)
		assert.PanicsWithError(t, expectedError.Error(), func() {
			res.ExpectErr(msg)
		})
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		msg := fake.RandomStringWithLength(8)

		res := Err(err)

		assert.Equal(t, err, res.ExpectErr(msg))
	})
}

func TestIsErr(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		v := fake.RandomStringWithLength(8)

		res := Ok(v)

		assert.False(t, res.IsErr())
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))

		res := Err(err)

		assert.True(t, res.IsErr())
	})
}

func TestIsErrAnd(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		v := fake.RandomStringWithLength(8)

		res := Ok(v)
		predicate := func(e error) bool {
			assert.Fail(t, "predicate should not have been called")
			return true
		}

		assert.False(t, res.IsErrAnd(predicate))
	})

	t.Run("err with non matching predicate", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))

		res := Err(err)
		predicate := func(e error) bool {
			assert.Equal(t, err, e)
			return false
		}

		assert.False(t, res.IsErrAnd(predicate))
	})

	t.Run("err with matching predicate", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))

		res := Err(err)
		predicate := func(e error) bool {
			assert.Equal(t, err, e)
			return true
		}

		assert.True(t, res.IsErrAnd(predicate))
	})
}

func TestIsOk(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		v := fake.RandomStringWithLength(8)

		res := Ok(v)

		assert.True(t, res.IsOk())
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))

		res := Err(err)

		assert.False(t, res.IsOk())
	})
}

func TestIsOkAnd(t *testing.T) {
	t.Run("ok with matching predicate", func(t *testing.T) {
		v := fake.RandomStringWithLength(8)

		res := Ok(v)
		predicate := func(s string) bool {
			assert.Equal(t, v, s)
			return true
		}

		assert.True(t, res.IsOkAnd(predicate))
	})

	t.Run("ok with non matching predicate", func(t *testing.T) {
		v := fake.RandomStringWithLength(8)

		res := Ok(v)
		predicate := func(s string) bool {
			assert.Equal(t, v, s)
			return false
		}

		assert.False(t, res.IsOkAnd(predicate))
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))

		res := Err(err)
		predicate := func(_ any) bool {
			assert.Fail(t, "predicate should not have been called")
			return true
		}

		assert.False(t, res.IsOkAnd(predicate))
	})
}
