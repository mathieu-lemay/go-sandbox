package result

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
)

var fake = faker.NewWithSeed(rand.NewSource(time.Now().UnixNano()))

type TestErrorT struct {
	e string
}

func (t *TestErrorT) Error() string {
	return t.e
}

func TestFrom_ReturnsNewResultFromArgs(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		v := fake.RandomStringWithLength(8)

		f := func() (string, error) {
			return v, nil
		}

		res := From(f())

		assert.True(t, res.IsOk())
	})

	t.Run("err", func(t *testing.T) {
		v := fake.RandomStringWithLength(8)
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))

		f := func() (string, error) {
			return v, err
		}

		res := From(f())

		assert.True(t, res.IsErr())
	})
}

func TestMap_ReturnsANewResultWithMappedValue(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok(value)

		expected := Ok(strconv.Itoa(value))

		assert.Equal(t, expected, Map(o, strconv.Itoa))
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		e := Err(err)

		f := func(any) int {
			assert.Fail(t, "mapper should not have been called")
			return 0
		}

		expected := errT[int, error]{err}

		assert.Equal(t, expected, Map(e, f))
	})
}

func TestMapOr_ReturnsTheMappedValueOrDefault(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok(value)

		def := fake.RandomStringWithLength(9)

		expected := strconv.Itoa(value)

		assert.Equal(t, expected, MapOr(o, def, strconv.Itoa))
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		e := Err(err)

		def := fake.RandomStringWithLength(9)
		f := func(any) string {
			assert.Fail(t, "mapper should not have been called")
			return ""
		}

		assert.Equal(t, def, MapOr(e, def, f))
	})
}

func TestMapOrElse_ReturnsTheMappedValueOrCallsDefaultFactory(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok(value)

		factory := func() string {
			assert.Fail(t, "factory should not have been called")
			return fake.RandomStringWithLength(9)
		}

		expected := strconv.Itoa(value)

		assert.Equal(t, expected, MapOrElse(o, factory, strconv.Itoa))
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		e := Err(err)

		mapper := func(any) string {
			assert.Fail(t, "mapper should not have been called")
			return ""
		}
		def := fake.RandomStringWithLength(9)
		factory := func() string {
			return def
		}

		assert.Equal(t, def, MapOrElse(e, factory, mapper))
	})
}

func TestMapErr_ReturnsANewResultWithMappedValue(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		value := fake.Int()
		o := Ok(value)

		f := func(error) *TestErrorT {
			assert.Fail(t, "mapper should not have been called")
			return nil
		}

		expected := ok[int, *TestErrorT]{value}

		assert.Equal(t, expected, MapErr(o, f))
	})

	t.Run("err", func(t *testing.T) {
		err := fmt.Errorf("some error: %s", fake.RandomStringWithLength(8))
		e := Err(err)

		newE := &TestErrorT{e: fmt.Errorf("mapped error: %w", err).Error()}
		f := func(e error) *TestErrorT {
			assert.Equal(t, e, err)

			return newE
		}

		expected := errT[any, *TestErrorT]{err: newE}

		assert.Equal(t, expected, MapErr(e, f))
	})
}
