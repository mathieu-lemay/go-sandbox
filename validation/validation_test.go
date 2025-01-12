package validation

import (
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/go-playground/validator/v10"
	"github.com/jaswdr/faker/v2"
	"github.com/mathieu-lemay/go-sandbox/logging"
	"github.com/stretchr/testify/assert"
)

var (
	_    = logging.ConfigureLogger(logging.WithLevel(zerolog.InfoLevel))
	fake = faker.New()
)

type JsonMap map[string]interface{}

type T struct {
	I int    `json:"i" validate:"required"`
	S string `json:"s" validate:"required"`
	B bool   `json:"b" validate:"required"`
}

type PT struct {
	I *int    `json:"i" validate:"required"`
	S *string `json:"s" validate:"required"`
	B *bool   `json:"b" validate:"required"`
}

func TestDeserializeSimpleStruct(t *testing.T) {

	testCases := []struct {
		name     string
		data     string
		isValid  bool
		expected *T
	}{
		{
			"missing values",
			`{}`,
			false,
			nil,
		},
		{
			"null values",
			`{
				"i": null,
				"s": null,
				"b": null
			}`,
			false,
			nil,
		},
		{
			"zero values",
			`{
				"i": 0,
				"s": "",
				"b": false
			}`,
			true,
			&T{I: 0, S: "", B: false},
		},
		{
			"non-zero values",
			`{
				"i": 5,
				"s": "string",
				"b": true
			}`,
			true,
			&T{I: 5, S: "string", B: true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var st T
			st, err := Deserialize[T]([]byte(tc.data))
			if tc.isValid {
				assert.NoError(t, err, tc.name)
				assert.Equal(t, *tc.expected, st, tc.name)
			} else {
				assert.Error(t, err, tc.name)
				assert.ErrorAs(t, err, &validator.ValidationErrors{}, tc.name)
			}
		})
	}
}

func TestDeserializeStructWithPointerFields(t *testing.T) {
	testCases := []struct {
		name     string
		data     string
		isValid  bool
		expected *PT
	}{
		{
			"missing values",
			`{}`,
			false,
			nil,
		},
		{
			"null values",
			`{
				"i": null,
				"s": null,
				"b": null
			}`,
			false,
			nil,
		},
		{
			"zero values",
			`{
				"i": 0,
				"s": "",
				"b": false
			}`,
			true,
			&PT{I: ptr(0), S: ptr(""), B: ptr(false)},
		},
		{
			"non-zero values",
			`{
				"i": 5,
				"s": "string",
				"b": true
			}`,
			true,
			&PT{I: ptr(5), S: ptr("string"), B: ptr(true)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			st, err := Deserialize[PT]([]byte(tc.data))
			if tc.isValid {
				assert.NoError(t, err, tc.name)
				assert.Equal(t, *tc.expected, st, tc.name)
			} else {
				assert.Error(t, err, tc.name)
				assert.ErrorAs(t, err, &validator.ValidationErrors{}, tc.name)
			}
		})
	}
}

func TestDeserializeStructWithComplexFields(t *testing.T) {
	type T struct {
		I int `json:"i" validate:"required"`
	}

	type S struct {
		T T `json:"t" validate:"required"`
	}

	testCases := []struct {
		name     string
		data     string
		isValid  bool
		expected *S
	}{
		{
			"missing value",
			`{}`,
			false,
			nil,
		},
		{
			"null value",
			`{
				"t": null
			}`,
			false,
			nil,
		},
		{
			"missing sub value",
			`{
				"t": {
					"i": null
				}
			}`,
			false,
			nil,
		},
		{
			"zero sub value",
			`{
				"t": {
					"i": 0
				}
			}`,
			true,
			&S{T: T{I: 0}},
		},
		{
			"non zero sub value",
			`{
				"t": {
					"i": 5
				}
			}`,
			true,
			&S{T: T{I: 5}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			st, err := Deserialize[S]([]byte(tc.data))
			if tc.isValid {
				require.NoError(t, err, tc.name)
				assert.Equal(t, *tc.expected, st, tc.name)
			} else {
				require.Error(t, err, tc.name)
				assert.ErrorAs(t, err, &validator.ValidationErrors{}, tc.name)
			}
		})
	}
}

func TestDeserializeStructSliceField(t *testing.T) {
	type T struct {
		I int `json:"i" validate:"required"`
	}

	type S struct {
		Ts []T `json:"ts" validate:"required,dive"`
	}

	testCases := []struct {
		name     string
		data     string
		isValid  bool
		expected *S
	}{
		{
			"missing value",
			`{}`,
			false,
			nil,
		},
		{
			"null value",
			`{
				"t": null
			}`,
			false,
			nil,
		},
		{
			"missing sub value",
			`{
				"t": [
					{
						"i": null
					}
				]
			}`,
			false,
			nil,
		},
		{
			"zero sub value",
			`{
				"t": [
					{
						"i": 0
					}
				]
			}`,
			true,
			&S{
				Ts: []T{
					{I: 0},
				},
			},
		},
		{
			"non zero sub value",
			`{
				"t": [
					{
						"i": 5
					}
				]
			}`,
			true,
			&S{
				Ts: []T{
					{I: 5},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			st, err := Deserialize[S]([]byte(tc.data))
			if tc.isValid {
				require.NoError(t, err, tc.name)
				assert.Equal(t, *tc.expected, st, tc.name)
			} else {
				require.Error(t, err, tc.name)
				assert.ErrorAs(t, err, &validator.ValidationErrors{}, tc.name)
			}
		})
	}
}

func TestDeserialize_EnsuresNonPointerFieldsAreRequired(t *testing.T) {
	type S struct {
		Name string `validate:"required"`
	}

	_, err := Deserialize[S]([]byte("{}"))

	assert.ErrorAs(t, err, &validator.ValidationErrors{})
}

func TestParse(t *testing.T) {
	testCases := []struct {
		msg      string
		data     string
		isValid  bool
		expected JsonMap
	}{
		{
			"missing values",
			`{}`,
			false,
			JsonMap{
				"i": nil,
				"s": nil,
				"b": nil,
			},
		},
		{
			"null values",
			`{
				"i": null,
				"s": null,
				"b": null
			}`,
			false,
			JsonMap{
				"i": nil,
				"s": nil,
				"b": nil,
			},
		},
		{
			"zero values",
			`{
				"i": 0,
				"s": "",
				"b": false
			}`,
			true,
			JsonMap{
				"i": float64(0),
				"s": "",
				"b": false,
			},
		},
		{
			"non-zero values",
			`{
				"i": 5,
				"s": "string",
				"b": true
			}`,
			true,
			JsonMap{
				"i": float64(5),
				"s": "string",
				"b": true,
			},
		},
	}

	for _, tc := range testCases {
		var st T
		res, err := parse([]byte(tc.data), &st)
		assert.Equal(t, tc.expected, must(toMap(res)), tc.msg)
		assert.NoError(t, err, tc.msg)
	}
}

func TestCopyStruct(t *testing.T) {
	src := PT{
		I: ptr(42),
		S: ptr("str"),
		B: ptr(true),
	}

	var dst T

	err := copyStruct(&src, &dst)

	expected := T{
		I: 42,
		S: "str",
		B: true,
	}

	assert.Equal(t, expected, dst)
	assert.NoError(t, err)
}

func TestCopyStructPointers(t *testing.T) {
	i := fake.Int()
	s := fake.RandomStringWithLength(8)
	b := fake.Bool()

	src := PT{
		I: &i,
		S: &s,
		B: &b,
	}

	var dst PT

	err := copyStruct(&src, &dst)
	require.NoError(t, err)

	assert.Equal(t, src, dst)
	assert.Equal(t, i, *dst.I)
	assert.Equal(t, s, *dst.S)
	assert.Equal(t, b, *dst.B)
}

func TestCopyStruct_Errors(t *testing.T) {
	src := PT{
		I: nil,
		S: ptr("str"),
		B: ptr(true),
	}

	var dst T

	err := copyStruct(&src, &dst)

	assert.ErrorContains(t, err, "can't set nil value to non ptr field: I")
}

func BenchmarkReference(b *testing.B) {
	_ = logging.ConfigureLogger(logging.WithLevel(zerolog.ErrorLevel))
	data := []byte(`{"i": 45, "s": "some-str", "b": true}`)
	for i := 0; i < b.N; i++ {
		var st T
		err := json.Unmarshal(data, &st)
		assert.NoError(b, err)

		err = defaultValidator.Struct(&st)
		assert.NoError(b, err)
	}
}

func BenchmarkDeserializer(b *testing.B) {
	_ = logging.ConfigureLogger(logging.WithLevel(zerolog.ErrorLevel))
	data := []byte(`{"i": 45, "s": "some-str", "b": true}`)
	for i := 0; i < b.N; i++ {
		_, err := Deserialize[T](data)
		assert.NoError(b, err)
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

func ptr[T any](v T) *T {
	return &v
}

func toMap(t interface{}) (JsonMap, error) {
	a, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	var res JsonMap
	json.Unmarshal(a, &res)

	return res, nil
}
