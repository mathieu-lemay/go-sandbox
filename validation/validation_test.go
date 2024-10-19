package main

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

type T struct {
	I int    `json:"i" validate:"required"`
	S string `json:"s" validate:"required"`
	B bool   `json:"b" validate:"required"`
}

func TestParseAndValidateOptionalValues(t *testing.T) {
	log.Logger = zerolog.Logger{}
	testCases := []struct {
		msg      string
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
		var st T
		err := parseAndValidate(&st, []byte(tc.data))
		if tc.isValid {
			assert.NoError(t, err, tc.msg)
			assert.Equal(t, *tc.expected, st, tc.msg)
		} else {
			assert.Error(t, err, tc.msg)
			assert.ErrorAs(t, err, &validator.ValidationErrors{}, tc.msg)
		}
	}
}

func BenchmarkParseAndValidate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var st T
		err := parseAndValidate(&st, []byte(`{}`))
		assert.Error(b, err)
	}
}
