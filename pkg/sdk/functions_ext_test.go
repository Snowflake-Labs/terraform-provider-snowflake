package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO [next PR]: test parsing single
func Test_parseFunctionDetailsImport(t *testing.T) {
	inputs := []struct {
		rawInput string
		expected []NormalizedPath
	}{
		{"", []NormalizedPath{}},
		{`[]`, []NormalizedPath{}},
		{`[@~/abc]`, []NormalizedPath{{"~", "abc"}}},
		{`[@~/abc/def]`, []NormalizedPath{{"~", "abc/def"}}},
		{`[@"db"."sc"."st"/abc/def]`, []NormalizedPath{{`"db"."sc"."st"`, "abc/def"}}},
		{`[@db.sc.st/abc/def]`, []NormalizedPath{{`"db"."sc"."st"`, "abc/def"}}},
		{`[db.sc.st/abc/def]`, []NormalizedPath{{`"db"."sc"."st"`, "abc/def"}}},
		{`[@"db"."sc".st/abc/def]`, []NormalizedPath{{`"db"."sc"."st"`, "abc/def"}}},
		{`[@"db"."sc".st/abc/def, db."sc".st/abc]`, []NormalizedPath{{`"db"."sc"."st"`, "abc/def"}, {`"db"."sc"."st"`, "abc"}}},
	}

	badInputs := []struct {
		rawInput          string
		expectedErrorPart string
	}{
		{"[", "brackets not find"},
		{"]", "brackets not find"},
		{`[@~/]`, "contains empty path"},
		{`[@~]`, "cannot be split into stage and path"},
		{`[@"db"."sc"/abc]`, "contains incorrect stage location"},
		{`[@"db"/abc]`, "contains incorrect stage location"},
		{`[@"db"."sc"."st"."smth"/abc]`, "contains incorrect stage location"},
		{`[@"db/a"."sc"."st"/abc]`, "contains incorrect stage location"},
		{`[@"db"."sc"."st"/abc], @"db"."sc"/abc]`, "contains incorrect stage location"},
	}

	for _, tc := range inputs {
		tc := tc
		t.Run(fmt.Sprintf("Snowflake raw imports: %s", tc.rawInput), func(t *testing.T) {
			details := FunctionDetails{Imports: &tc.rawInput}

			results, err := parseFunctionDetailsImport(details)
			require.NoError(t, err)
			require.Equal(t, tc.expected, results)
		})
	}

	for _, tc := range badInputs {
		tc := tc
		t.Run(fmt.Sprintf("incorrect Snowflake input: %s, expecting error with: %s", tc.rawInput, tc.expectedErrorPart), func(t *testing.T) {
			details := FunctionDetails{Imports: &tc.rawInput}

			_, err := parseFunctionDetailsImport(details)
			require.Error(t, err)
			require.ErrorContains(t, err, "could not parse imports from Snowflake")
			require.ErrorContains(t, err, tc.expectedErrorPart)
		})
	}

	t.Run("Snowflake raw imports nil", func(t *testing.T) {
		details := FunctionDetails{Imports: nil}

		results, err := parseFunctionDetailsImport(details)
		require.NoError(t, err)
		require.Equal(t, []NormalizedPath{}, results)
	})
}

func Test_parseFunctionAndProcedureReturns(t *testing.T) {
	inputs := []struct {
		rawInput              string
		expectedRawDataType   string
		expectedReturnNotNull bool
	}{
		{"CHAR", "CHAR(1)", false},
		{"CHAR(1)", "CHAR(1)", false},
		{"CHAR NOT NULL", "CHAR(1)", true},
		{"  CHAR   NOT NULL  ", "CHAR(1)", true},
		{"OBJECT", "OBJECT", false},
		{"OBJECT NOT NULL", "OBJECT", true},
	}

	badInputs := []struct {
		rawInput          string
		expectedErrorPart string
	}{
		{"", "invalid data type"},
		{"NOT NULL", "invalid data type"},
		{"CHA NOT NULL", "invalid data type"},
		{"CHA NOT NULLS", "invalid data type"},
	}

	for _, tc := range inputs {
		tc := tc
		t.Run(fmt.Sprintf("return data type raw: %s", tc.rawInput), func(t *testing.T) {
			dt, returnNotNull, err := parseFunctionAndProcedureReturns(tc.rawInput)
			require.NoError(t, err)
			require.Equal(t, tc.expectedRawDataType, dt.ToSql())
			require.Equal(t, tc.expectedReturnNotNull, returnNotNull)
		})
	}

	for _, tc := range badInputs {
		tc := tc
		t.Run(fmt.Sprintf("incorrect return data type raw: %s, expecting error with: %s", tc.rawInput, tc.expectedErrorPart), func(t *testing.T) {
			_, _, err := parseFunctionAndProcedureReturns(tc.rawInput)
			require.Error(t, err)
			require.ErrorContains(t, err, tc.expectedErrorPart)
		})
	}
}
