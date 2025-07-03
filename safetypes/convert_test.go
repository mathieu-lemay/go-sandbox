package safetypes

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"

	"github.com/mathieu-lemay/go-sandbox/safetypes/option"
	"github.com/mathieu-lemay/go-sandbox/safetypes/result"
)

var fake = faker.NewWithSeed(rand.NewSource(time.Now().UnixNano()))

func TestAsOkOr(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		value := fake.Int()
		s := option.Some(value)

		expected := result.Ok(value)

		assert.Equal(t, expected, AsOkOr(s, errors.New(fake.RandomStringWithLength(8))))
	})

	t.Run("none", func(t *testing.T) {
		n := option.From(0)
		err := errors.New(fake.RandomStringWithLength(8))

		expected := result.From(0, err)

		assert.Equal(t, expected, AsOkOr(n, err))
	})
}

func TestAsOkOrElse(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		value := fake.Int()
		s := option.Some(value)

		f := func() error {
			assert.Fail(t, "should not be called")

			return errors.New(fake.RandomStringWithLength(8))
		}

		expected := result.Ok(value)

		assert.Equal(t, expected, AsOkOrElse(s, f))
	})

	t.Run("none", func(t *testing.T) {
		n := option.From(0)
		err := errors.New(fake.RandomStringWithLength(8))

		f := func() error {
			return err
		}

		expected := result.From(0, err)

		assert.Equal(t, expected, AsOkOrElse(n, f))
	})
}

func TestAsOptionValue(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := result.Ok(value)

		expected := option.Some(value)

		assert.Equal(t, expected, AsOptionValue(o))
	})

	t.Run("err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		e := result.From(0, err)

		expected := option.From(0)

		assert.Equal(t, expected, AsOptionValue(e))
	})
}

func TestAsOptionErr(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := result.Ok(value)

		expected := option.From((error)(nil))

		assert.Equal(t, expected, AsOptionErr(o))
	})

	t.Run("err", func(t *testing.T) {
		err := errors.New(fake.RandomStringWithLength(8))
		e := result.From(0, err)

		expected := option.From(err)

		assert.Equal(t, expected, AsOptionErr(e))
	})
}
