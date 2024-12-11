package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseFunctionDetailsImport(t *testing.T) {
	inputs := []struct {
		rawInput string
		expected []FunctionDetailsImport
	}{
		{"", []FunctionDetailsImport{}},
		{`[]`, []FunctionDetailsImport{}},
		{`[@~/abc]`, []FunctionDetailsImport{{"~", "abc"}}},
		{`[@~/abc/def]`, []FunctionDetailsImport{{"~", "abc/def"}}},
		{`[@"db"."sc"."st"/abc/def]`, []FunctionDetailsImport{{`"db"."sc"."st"`, "abc/def"}}},
		{`[@db.sc.st/abc/def]`, []FunctionDetailsImport{{`"db"."sc"."st"`, "abc/def"}}},
		{`[db.sc.st/abc/def]`, []FunctionDetailsImport{{`"db"."sc"."st"`, "abc/def"}}},
		{`[@"db"."sc".st/abc/def]`, []FunctionDetailsImport{{`"db"."sc"."st"`, "abc/def"}}},
		{`[@"db"."sc".st/abc/def, db."sc".st/abc]`, []FunctionDetailsImport{{`"db"."sc"."st"`, "abc/def"}, {`"db"."sc"."st"`, "abc"}}},
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
}
