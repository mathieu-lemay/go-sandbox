package option

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
)

var fake = faker.NewWithSeed(rand.NewSource(time.Now().UnixNano()))

func TestFrom_ReturnsNewOptionFromArgs(t *testing.T) {
	type S struct {
		value int
	}

	t.Run("some", func(t *testing.T) {
		res1 := From(fake.IntBetween(1, 100))
		assert.True(t, res1.IsSome())

		res2 := From(fake.Float(2, 1, 100))
		assert.True(t, res2.IsSome())

		res3 := From(fake.RandomStringWithLength(8))
		assert.True(t, res3.IsSome())

		res4 := From(true)
		assert.True(t, res4.IsSome())

		res5 := From(S{value: 1})
		assert.True(t, res5.IsSome())

		zeroInt := 0
		res6 := From(&zeroInt)
		assert.True(t, res6.IsSome())

		zeroStr := ""
		res7 := From(&zeroStr)
		assert.True(t, res7.IsSome())

		zeroS := S{} //nolint:exhaustruct  // We want the zero value
		res8 := From(&zeroS)
		assert.True(t, res8.IsSome())
	})

	t.Run("none", func(t *testing.T) {
		res1 := From(0)
		assert.True(t, res1.IsNone())

		res2 := From(0.0)
		assert.True(t, res2.IsNone())

		res3 := From("")
		assert.True(t, res3.IsNone())

		res4 := From(false)
		assert.True(t, res4.IsNone())

		res5 := From(S{}) //nolint:exhaustruct  // We want the zero value
		assert.True(t, res5.IsNone())

		res6 := From((*string)(nil))
		assert.True(t, res6.IsNone())
	})
}

func TestMap_ReturnsANewOptionWithMappedValue(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		s := Some(val)

		expected := Some(strconv.Itoa(val))

		assert.Equal(t, expected, Map(s, strconv.Itoa))
	})

	t.Run("none", func(t *testing.T) {
		n := None()

		f := func(any) int {
			assert.Fail(t, "mapper should not have been called")

			return 0
		}

		assert.Equal(t, none[int]{}, Map(n, f))
	})
}

func TestMapOr_ReturnsTheMappedValueOrDefault(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		s := Some(val)

		def := fake.RandomStringWithLength(9)

		expected := strconv.Itoa(val)

		assert.Equal(t, expected, MapOr(s, def, strconv.Itoa))
	})

	t.Run("none", func(t *testing.T) {
		n := None()

		def := fake.RandomStringWithLength(9)
		f := func(any) string {
			assert.Fail(t, "mapper should not have been called")

			return ""
		}

		assert.Equal(t, def, MapOr(n, def, f))
	})
}

func TestMapOrElse_ReturnsTheMappedValueOrCallsDefaultFactory(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		s := Some(val)

		factory := func() string {
			assert.Fail(t, "factory should not have been called")

			return fake.RandomStringWithLength(9)
		}

		expected := strconv.Itoa(val)

		assert.Equal(t, expected, MapOrElse(s, factory, strconv.Itoa))
	})

	t.Run("none", func(t *testing.T) {
		n := None()

		mapper := func(any) string {
			assert.Fail(t, "mapper should not have been called")

			return ""
		}
		def := fake.RandomStringWithLength(9)
		factory := func() string {
			return def
		}

		assert.Equal(t, def, MapOrElse(n, factory, mapper))
	})
}

func TestAnd_ReturnsOtherOrNone(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		s := Some(val)

		other := Some(fake.RandomStringWithLength(8))

		assert.Equal(t, other, And(s, other))
	})

	t.Run("none", func(t *testing.T) {
		n := None()

		other := Some(fake.RandomStringWithLength(8))

		assert.Equal(t, none[string]{}, And(n, other))
	})
}

func TestAndThen_ReturnsMappedOrNone(t *testing.T) {
	t.Run("some", func(t *testing.T) {
		val := fake.Int()
		s := Some(val)

		other := Some(fake.RandomStringWithLength(8))
		f := func(o int) Option[string] {
			assert.Equal(t, val, o)

			return other
		}

		assert.Equal(t, other, AndThen(s, f))
	})

	t.Run("none", func(t *testing.T) {
		n := None()

		f := func(_ any) Option[string] {
			assert.Fail(t, "mapper should not have been called")

			return none[string]{}
		}

		assert.Equal(t, none[string]{}, AndThen(n, f))
	})
}
