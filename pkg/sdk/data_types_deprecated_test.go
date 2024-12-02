package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsStringType(t *testing.T) {
	type test struct {
		input string
		want  bool
	}

	tests := []test{
		// case insensitive.
		{input: "STRING", want: true},
		{input: "string", want: true},
		{input: "String", want: true},

		// varchar types.
		{input: "VARCHAR", want: true},
		{input: "NVARCHAR", want: true},
		{input: "NVARCHAR2", want: true},
		{input: "CHAR", want: true},
		{input: "NCHAR", want: true},
		{input: "CHAR VARYING", want: true},
		{input: "NCHAR VARYING", want: true},
		{input: "TEXT", want: true},

		// with length
		{input: "VARCHAR(100)", want: true},
		{input: "NVARCHAR(100)", want: true},
		{input: "NVARCHAR2(100)", want: true},
		{input: "CHAR(100)", want: true},
		{input: "NCHAR(100)", want: true},
		{input: "CHAR VARYING(100)", want: true},
		{input: "NCHAR VARYING(100)", want: true},
		{input: "TEXT(100)", want: true},

		// binary is not string types.
		{input: "binary", want: false},
		{input: "varbinary", want: false},

		// other types
		{input: "boolean", want: false},
		{input: "number", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := IsStringType(tc.input)
			require.Equal(t, tc.want, got)
		})
	}
}

func Test_ParseNumberDataTypeRaw(t *testing.T) {
	type test struct {
		input             string
		expectedPrecision int
		expectedScale     int
	}
	defaults := func(input string) test {
		return test{input: input, expectedPrecision: DefaultNumberPrecision, expectedScale: DefaultNumberScale}
	}

	tests := []test{
		{input: "NUMBER(30)", expectedPrecision: 30, expectedScale: DefaultNumberScale},
		{input: "NUMBER(30, 2)", expectedPrecision: 30, expectedScale: 2},
		{input: "decimal(30, 2)", expectedPrecision: 30, expectedScale: 2},
		{input: "NUMBER(   30   ,  2   )", expectedPrecision: 30, expectedScale: 2},
		{input: "    NUMBER   (   30   ,  2   )    ", expectedPrecision: 30, expectedScale: 2},

		// returns defaults if it can't parse arguments, data type is different, or no arguments were provided
		defaults("VARCHAR(1, 2)"),
		defaults("VARCHAR(1)"),
		defaults("VARCHAR"),
		defaults("NUMBER"),
		defaults("NUMBER()"),
		defaults("NUMBER(x)"),
		defaults(fmt.Sprintf("NUMBER(%d)", DefaultNumberPrecision)),
		defaults(fmt.Sprintf("NUMBER(%d, x)", DefaultNumberPrecision)),
		defaults(fmt.Sprintf("NUMBER(x, %d)", DefaultNumberScale)),
		defaults("NUMBER(1, 2, 3)"),
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			precision, scale := ParseNumberDataTypeRaw(tc.input)
			assert.Equal(t, tc.expectedPrecision, precision)
			assert.Equal(t, tc.expectedScale, scale)
		})
	}
}

func Test_ParseVarcharDataTypeRaw(t *testing.T) {
	type test struct {
		input          string
		expectedLength int
	}
	defaults := func(input string) test {
		return test{input: input, expectedLength: DefaultVarcharLength}
	}

	tests := []test{
		{input: "VARCHAR(30)", expectedLength: 30},
		{input: "text(30)", expectedLength: 30},
		{input: "VARCHAR(   30   )", expectedLength: 30},
		{input: "    VARCHAR   (   30   )    ", expectedLength: 30},

		// returns defaults if it can't parse arguments, data type is different, or no arguments were provided
		defaults("VARCHAR(1, 2)"),
		defaults("VARCHAR(x)"),
		defaults("VARCHAR"),
		defaults("NUMBER"),
		defaults("NUMBER()"),
		defaults("NUMBER(x)"),
		defaults(fmt.Sprintf("VARCHAR(%d)", DefaultVarcharLength)),
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			length := ParseVarcharDataTypeRaw(tc.input)
			assert.Equal(t, tc.expectedLength, length)
		})
	}
}
