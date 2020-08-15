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

func TestString(t *testing.T) {
	assertParse(t, "Future Gadget Laboratory", `"Future Gadget Laboratory"`)
	assertParse(t, "", `""`)
	assertParse(t, "Time\n\r\t\b\\/Leap", `"Time\n\r\t\b\\\/Leap"`)
	assertParse(t, "\\n", `"\\n"`)
	assertParse(t, "Time \u23F1", `"Time \u23F1"`)
	assertParse(t, "Phone \u260E", `"Phone \u260e"`)

	assertParse(t, `/`, `"/"`)
	assertParse(t, `/`, `"\/"`)
	assertParse(t, `/`, `"\u002F"`)
	assertParse(t, `/`, `"\u002f"`)

	assertParseError(
		t,
		`
		"Beta World Line
		Alpha World Line"
		`,
	)
	assertParseError(t, `"John`)
	assertParseError(t, `'John Titor'`)
	assertParseError(t, `"\uGGGG"`)
	assertParseError(t, `"\u123"`)
	assertParseError(t, `"\u"`)
	assertParseError(t, `"\a"`)
	assertParseError(t, `"\ n"`)
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

	assertParseError(t, "[123,]")
}

func TestObject(t *testing.T) {
	assertParse(t, map[string]interface{}{}, "{}")
	assertParse(
		t,
		map[string]interface{}{"No.001": "Okabe Rintaro", "No.004": "Makise Kurisu"},
		`{ "No.001" : "Okabe Rintaro", "No.004" : "Makise Kurisu" }`,
	)
	assertParse(
		t,
		map[string]interface{}{
			"number int":   123,
			"number float": 3.14,
			"string":       "Makise Kurisu",
			"array":        []interface{}{333, 444},
			"true":         true,
			"false":        false,
			"null":         nil,
			"object":       map[string]interface{}{"C204": "Time Machine"},
		},
		`{
			"number int": 123,
			"number float":3.14,
			"string": "Makise Kurisu",
			"array": [333, 444],
			"true": true,
			"false": false,
			"null": null,
			"object": {"C204": "Time Machine"}
		}`,
	)

	assertParseError(t, "{123,}")
	assertParseError(t, `{123: "key should be string"}`)
	assertParseError(t, `{true: "key should be string"}`)
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
