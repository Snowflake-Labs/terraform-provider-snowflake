package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToDataType(t *testing.T) {
	type test struct {
		input string
		want  DataType
	}

	tests := []test{
		// case insensitive.
		{input: "STRING", want: DataTypeVARCHAR},
		{input: "string", want: DataTypeVARCHAR},
		{input: "String", want: DataTypeVARCHAR},

		// number types.
		{input: "number", want: DataTypeNumber},
		{input: "decimal", want: DataTypeNumber},
		{input: "numeric", want: DataTypeNumber},
		{input: "int", want: DataTypeNumber},
		{input: "integer", want: DataTypeNumber},
		{input: "bigint", want: DataTypeNumber},
		{input: "smallint", want: DataTypeNumber},
		{input: "tinyint", want: DataTypeNumber},
		{input: "byteint", want: DataTypeNumber},

		// float types.
		{input: "float", want: DataTypeFloat},
		{input: "float4", want: DataTypeFloat},
		{input: "float8", want: DataTypeFloat},
		{input: "double", want: DataTypeFloat},
		{input: "double precision", want: DataTypeFloat},
		{input: "real", want: DataTypeFloat},

		// varchar types.
		{input: "varchar", want: DataTypeVARCHAR},
		{input: "char", want: DataTypeVARCHAR},
		{input: "character", want: DataTypeVARCHAR},
		{input: "string", want: DataTypeVARCHAR},
		{input: "text", want: DataTypeVARCHAR},

		// binary types.
		{input: "binary", want: DataTypeBinary},
		{input: "varbinary", want: DataTypeBinary},
		{input: "boolean", want: DataTypeBoolean},

		// boolean types.
		{input: "boolean", want: DataTypeBoolean},
		{input: "bool", want: DataTypeBoolean},

		// timestamp ntz types.
		{input: "datetime", want: DataTypeTimestampNTZ},
		{input: "timestamp", want: DataTypeTimestampNTZ},
		{input: "timestamp_ntz", want: DataTypeTimestampNTZ},

		// timestamp tz types.
		{input: "timestamp_tz", want: DataTypeTimestampTZ},
		{input: "timestamp_tz(9)", want: DataTypeTimestampTZ},

		// timestamp ltz types.
		{input: "timestamp_ltz", want: DataTypeTimestampLTZ},
		{input: "timestamp_ltz(9)", want: DataTypeTimestampLTZ},

		// time types.
		{input: "time", want: DataTypeTime},
		{input: "time(9)", want: DataTypeTime},

		// all othertypes
		{input: "date", want: DataTypeDate},
		{input: "variant", want: DataTypeVariant},
		{input: "object", want: DataTypeObject},
		{input: "array", want: DataTypeArray},
		{input: "geography", want: DataTypeGeography},
		{input: "geometry", want: DataTypeGeometry},
		{input: "VECTOR(INT, 10)", want: "VECTOR(INT, 10)"},
		{input: "VECTOR(INT, 20)", want: "VECTOR(INT, 20)"},
		{input: "VECTOR(FLOAT, 10)", want: "VECTOR(FLOAT, 10)"},
		{input: "VECTOR(FLOAT, 20)", want: "VECTOR(FLOAT, 20)"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToDataType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

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
