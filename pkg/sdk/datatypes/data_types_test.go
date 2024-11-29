package datatypes

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseDataType_Number(t *testing.T) {
	type test struct {
		input                  string
		expectedPrecision      int
		expectedScale          int
		expectedUnderlyingType string
	}
	defaults := func(input string) test {
		return test{
			input:                  input,
			expectedPrecision:      DefaultNumberPrecision,
			expectedScale:          DefaultNumberScale,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		{input: "NUMBER(30)", expectedPrecision: 30, expectedScale: DefaultNumberScale, expectedUnderlyingType: "NUMBER"},
		{input: "NUMBER(30, 2)", expectedPrecision: 30, expectedScale: 2, expectedUnderlyingType: "NUMBER"},
		{input: "dec(30)", expectedPrecision: 30, expectedScale: DefaultNumberScale, expectedUnderlyingType: "DEC"},
		{input: "dec(30, 2)", expectedPrecision: 30, expectedScale: 2, expectedUnderlyingType: "DEC"},
		{input: "decimal(30)", expectedPrecision: 30, expectedScale: DefaultNumberScale, expectedUnderlyingType: "DECIMAL"},
		{input: "decimal(30, 2)", expectedPrecision: 30, expectedScale: 2, expectedUnderlyingType: "DECIMAL"},
		{input: "NuMeRiC(30)", expectedPrecision: 30, expectedScale: DefaultNumberScale, expectedUnderlyingType: "NUMERIC"},
		{input: "NuMeRiC(30, 2)", expectedPrecision: 30, expectedScale: 2, expectedUnderlyingType: "NUMERIC"},
		{input: "NUMBER(   30   ,  2   )", expectedPrecision: 30, expectedScale: 2, expectedUnderlyingType: "NUMBER"},
		{input: "    NUMBER   (   30   ,  2   )    ", expectedPrecision: 30, expectedScale: 2, expectedUnderlyingType: "NUMBER"},
		{input: fmt.Sprintf("NUMBER(%d)", DefaultNumberPrecision), expectedPrecision: DefaultNumberPrecision, expectedScale: DefaultNumberScale, expectedUnderlyingType: "NUMBER"},
		{input: fmt.Sprintf("NUMBER(%d, %d)", DefaultNumberPrecision, DefaultNumberScale), expectedPrecision: DefaultNumberPrecision, expectedScale: DefaultNumberScale, expectedUnderlyingType: "NUMBER"},

		defaults("NUMBER"),
		defaults("DEC"),
		defaults("DECIMAL"),
		defaults("NUMERIC"),
		defaults("   NUMBER   "),

		defaults("INT"),
		defaults("INTEGER"),
		defaults("BIGINT"),
		defaults("SMALLINT"),
		defaults("TINYINT"),
		defaults("BYTEINT"),
		defaults("int"),
		defaults("integer"),
		defaults("bigint"),
		defaults("smallint"),
		defaults("tinyint"),
		defaults("byteint"),
	}

	negativeTestCases := []test{
		negative("other(1, 2)"),
		negative("other(1)"),
		negative("other"),
		negative("NUMBER()"),
		negative("NUMBER(x)"),
		negative(fmt.Sprintf("NUMBER(%d, x)", DefaultNumberPrecision)),
		negative(fmt.Sprintf("NUMBER(x, %d)", DefaultNumberScale)),
		negative("NUMBER(1, 2, 3)"),
		negative("NUMBER("),
		negative("NUMBER)"),
		negative("NUM BER"),
		negative("INT(30)"),
		negative("INT(30, 2)"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &NumberDataType{}, parsed)

			assert.Equal(t, tc.expectedPrecision, parsed.(*NumberDataType).precision)
			assert.Equal(t, tc.expectedScale, parsed.(*NumberDataType).scale)
			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*NumberDataType).underlyingType)
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_Float(t *testing.T) {
	type test struct {
		input                  string
		expectedUnderlyingType string
	}
	defaults := func(input string) test {
		return test{
			input:                  input,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		defaults("   FLOAT   "),
		defaults("FLOAT"),
		defaults("FLOAT4"),
		defaults("FLOAT8"),
		defaults("DOUBLE PRECISION"),
		defaults("DOUBLE"),
		defaults("REAL"),
		defaults("float"),
		defaults("float4"),
		defaults("float8"),
		defaults("double precision"),
		defaults("double"),
		defaults("real"),
	}

	negativeTestCases := []test{
		negative("FLOAT(38, 0)"),
		negative("FLOAT(38, 2)"),
		negative("FLOAT(38)"),
		negative("FLOAT()"),
		negative("F L O A T"),
		negative("other"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &FloatDataType{}, parsed)

			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*FloatDataType).underlyingType)
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_Text(t *testing.T) {
	type test struct {
		input                  string
		expectedLength         int
		expectedUnderlyingType string
	}
	defaultsVarchar := func(input string) test {
		return test{
			input:                  input,
			expectedLength:         DefaultVarcharLength,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	defaultsChar := func(input string) test {
		return test{
			input:                  input,
			expectedLength:         DefaultCharLength,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		{input: "VARCHAR(30)", expectedLength: 30, expectedUnderlyingType: "VARCHAR"},
		{input: "string(30)", expectedLength: 30, expectedUnderlyingType: "STRING"},
		{input: "VARCHAR(   30   )", expectedLength: 30, expectedUnderlyingType: "VARCHAR"},
		{input: "    VARCHAR   (   30   )    ", expectedLength: 30, expectedUnderlyingType: "VARCHAR"},
		{input: fmt.Sprintf("VARCHAR(%d)", DefaultVarcharLength), expectedLength: DefaultVarcharLength, expectedUnderlyingType: "VARCHAR"},

		{input: "CHAR(30)", expectedLength: 30, expectedUnderlyingType: "CHAR"},
		{input: "character(30)", expectedLength: 30, expectedUnderlyingType: "CHARACTER"},
		{input: "CHAR(   30   )", expectedLength: 30, expectedUnderlyingType: "CHAR"},
		{input: "    CHAR   (   30   )    ", expectedLength: 30, expectedUnderlyingType: "CHAR"},
		{input: fmt.Sprintf("CHAR(%d)", DefaultCharLength), expectedLength: DefaultCharLength, expectedUnderlyingType: "CHAR"},

		defaultsVarchar("   VARCHAR   "),
		defaultsVarchar("VARCHAR"),
		defaultsVarchar("STRING"),
		defaultsVarchar("TEXT"),
		defaultsVarchar("NVARCHAR"),
		defaultsVarchar("NVARCHAR2"),
		defaultsVarchar("CHAR VARYING"),
		defaultsVarchar("NCHAR VARYING"),
		defaultsVarchar("varchar"),
		defaultsVarchar("string"),
		defaultsVarchar("text"),
		defaultsVarchar("nvarchar"),
		defaultsVarchar("nvarchar2"),
		defaultsVarchar("char varying"),
		defaultsVarchar("nchar varying"),

		defaultsChar("   CHAR   "),
		defaultsChar("CHAR"),
		defaultsChar("CHARACTER"),
		defaultsChar("NCHAR"),
		defaultsChar("char"),
		defaultsChar("character"),
		defaultsChar("nchar"),
	}

	negativeTestCases := []test{
		negative("other(1, 2)"),
		negative("other(1)"),
		negative("other"),
		negative("VARCHAR()"),
		negative("VARCHAR(x)"),
		negative("VARCHAR(   )"),
		negative("CHAR()"),
		negative("CHAR(x)"),
		negative("CHAR(   )"),
		negative("VARCHAR(1, 2)"),
		negative("VARCHAR("),
		negative("VARCHAR)"),
		negative("VAR CHAR"),
		negative("CHAR(1, 2)"),
		negative("CHAR("),
		negative("CHAR)"),
		negative("CH AR"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &TextDataType{}, parsed)

			assert.Equal(t, tc.expectedLength, parsed.(*TextDataType).length)
			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*TextDataType).underlyingType)
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}
