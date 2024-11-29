package datatypes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseDataType_Number(t *testing.T) {
	type test struct {
		input             string
		expectedPrecision int
		expectedScale     int
	}
	defaults := func(input string) test {
		return test{input: input, expectedPrecision: DefaultNumberPrecision, expectedScale: DefaultNumberScale}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		{input: "NUMBER(30)", expectedPrecision: 30, expectedScale: DefaultNumberScale},
		{input: "NUMBER(30, 2)", expectedPrecision: 30, expectedScale: 2},
		{input: "dec(30)", expectedPrecision: 30, expectedScale: DefaultNumberScale},
		{input: "dec(30, 2)", expectedPrecision: 30, expectedScale: 2},
		{input: "decimal(30)", expectedPrecision: 30, expectedScale: DefaultNumberScale},
		{input: "decimal(30, 2)", expectedPrecision: 30, expectedScale: 2},
		{input: "NuMeRiC(30)", expectedPrecision: 30, expectedScale: DefaultNumberScale},
		{input: "NuMeRiC(30, 2)", expectedPrecision: 30, expectedScale: 2},
		{input: "NUMBER(   30   ,  2   )", expectedPrecision: 30, expectedScale: 2},
		{input: "    NUMBER   (   30   ,  2   )    ", expectedPrecision: 30, expectedScale: 2},

		defaults("NUMBER"),
		defaults("DEC"),
		defaults("DECIMAL"),
		defaults("NUMERIC"),
		defaults("   NUMBER   "),
		defaults(fmt.Sprintf("NUMBER(%d)", DefaultNumberPrecision)),

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
		negative("VARCHAR(1, 2)"),
		negative("VARCHAR(1)"),
		negative("VARCHAR"),
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
