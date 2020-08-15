package myjson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumber(t *testing.T) {
	assertParse(t, 123, "123")
	assertParse(t, -123, "-123")
	assertParse(t, 456.7, "456.7")
	assertParse(t, -456.7, "-456.7")
	assertParse(t, 3.14e-3, "3.14e-3")
	assertParse(t, 3.14e+3, "3.14e+3")
	assertParse(t, 3.14e+3, "3.14e3")
	assertParse(t, 3.14E-3, "3.14E-3")
	assertParse(t, -3.14e-3, "-3.14e-3")
	assertParse(t, 314e-3, "314e-3")
	assertParse(t, -314.e-3, "-314e-3")

	assertParse(t, 0, "0")
	assertParse(t, 0, "-0")
	assertParse(t, 100000, "100000")

	assertParseError(t, "0123")
	assertParseError(t, "+123")
	assertParseError(t, "3.14e-")
}

func TestLiteral(t *testing.T) {
	assertParse(t, true, "true")
	assertParse(t, false, "false")
	assertParse(t, nil, "null")

	assertParseError(t, "True")
	assertParseError(t, "False")
	assertParseError(t, "nil")
}

func TestArray(t *testing.T) {
	assertParse(t, []interface{}{}, "[]")
	assertParse(t, []interface{}{123, 456, 789}, "[123, 456, 789]")
	assertParse(
		t,
		[]interface{}{
			123,
			3.14,
			"Makise Kurisu",
			[]interface{}{333, 444},
			true,
			false,
			nil,
			map[string]interface{}{"C204": "Time Machine"},
		},
		`[123, 3.14, "Makise Kurisu", [333, 444], true, false, null, {"C204" : "Time Machine"}]`,
	)
}

func assertParse(t *testing.T, expected interface{}, json string) {
	t.Helper()

	actual, err := Parse(json)
	assert.Nil(t, err, json)
	assert.Equal(t, expected, actual, json)
}

func assertParseError(t *testing.T, json string) {
	t.Helper()

	_, err := Parse(json)
	assert.NotNil(t, err, json)
}
