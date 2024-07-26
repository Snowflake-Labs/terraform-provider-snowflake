package helpers

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func Test_ParseIdentifierString(t *testing.T) {
	containsAll := func(t *testing.T, parts, expectedParts []string) {
		t.Helper()
		require.Len(t, parts, len(expectedParts))
		for _, part := range expectedParts {
			require.Contains(t, parts, part)
		}
	}

	t.Run("returns read error", func(t *testing.T) {
		input := `ab"c`

		_, err := parseIdentifierString(input)

		require.ErrorContains(t, err, "unable to read identifier")
		require.ErrorContains(t, err, `bare " in non-quoted-field`)
	})

	t.Run("returns error for empty input", func(t *testing.T) {
		input := ""

		_, err := parseIdentifierString(input)

		require.ErrorContains(t, err, "incompatible identifier")
	})

	t.Run("returns error for multiple lines", func(t *testing.T) {
		input := "abc\ndef"

		_, err := parseIdentifierString(input)

		require.ErrorContains(t, err, "incompatible identifier")
	})

	t.Run("returns parts correctly without quoting", func(t *testing.T) {
		input := "abc.def"
		expected := []string{"abc", "def"}

		parts, err := parseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("returns parts correctly with quoting", func(t *testing.T) {
		input := `"abc"."def"`
		expected := []string{"abc", "def"}

		parts, err := parseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("returns parts correctly with mixed quoting", func(t *testing.T) {
		input := `"abc".def."ghi"`
		expected := []string{"abc", "def", "ghi"}

		parts, err := parseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	// Quote inside must have a preceding quote (https://docs.snowflake.com/en/sql-reference/identifiers-syntax).
	t.Run("returns parts correctly with quote inside", func(t *testing.T) {
		input := `"ab""c".def`
		expected := []string{`ab"c`, "def"}

		parts, err := parseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("returns parts correctly with dots inside", func(t *testing.T) {
		input := `"ab.c".def`
		expected := []string{`ab.c`, "def"}

		parts, err := parseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("empty identifier", func(t *testing.T) {
		input := `""`
		expected := []string{""}

		parts, err := parseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("handled correctly double quotes", func(t *testing.T) {
		input := `""."."".".".""."".""."".""."".""."""""`
		expected := []string{"", `.".`, `.".".".".".".".""`}

		parts, err := parseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})
}

func Test_ParseAccountObjectIdentifier(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected sdk.AccountObjectIdentifier
		Error    string
	}{
		{Input: ``, Error: "incompatible identifier: "},
		{Input: "a\nb", Error: "incompatible identifier: a\nb"},
		{Input: `a"b`, Error: "unable to read identifier: a\"b, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc.cde`, Error: `unexpected number of parts 2 in identifier abc.cde, expected 1 in a form of "<account_object_name>"`},
		{Input: `""`, Expected: sdk.NewAccountObjectIdentifier(``)},
		{Input: `""""`, Expected: sdk.NewAccountObjectIdentifier(`"`)},
		{Input: `abc`, Expected: sdk.NewAccountObjectIdentifier(`abc`)},
		{Input: `"abc"`, Expected: sdk.NewAccountObjectIdentifier(`abc`)},
		{Input: `"ab.c"`, Expected: sdk.NewAccountObjectIdentifier(`ab.c`)},
		{Input: `"a""bc"`, Expected: sdk.NewAccountObjectIdentifier(`a"bc`)},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Parsing account object identifier with input: "%s"`, testCase.Input), func(t *testing.T) {
			id, err := ParseAccountObjectIdentifier(testCase.Input)

			if testCase.Error != "" {
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.Equal(t, testCase.Expected, id)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_ParseDatabaseObjectIdentifier(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected sdk.DatabaseObjectIdentifier
		Error    string
	}{
		{Input: ``, Error: "incompatible identifier: "},
		{Input: "a\nb.cde", Error: "unable to read identifier: a\nb.cde, err = record on line 2: wrong number of fields"},
		{Input: `a"b.cde`, Error: "unable to read identifier: a\"b.cde, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc.cde.efg`, Error: `unexpected number of parts 3 in identifier abc.cde.efg, expected 2 in a form of "<database_name>.<database_object_name>"`},
		{Input: `abc`, Error: `unexpected number of parts 1 in identifier abc, expected 2 in a form of "<database_name>.<database_object_name>"`},
		{Input: `"".""`, Expected: sdk.NewDatabaseObjectIdentifier(``, ``)},
		{Input: `"""".""""`, Expected: sdk.NewDatabaseObjectIdentifier(`"`, `"`)},
		{Input: `abc.cde`, Expected: sdk.NewDatabaseObjectIdentifier(`abc`, `cde`)},
		{Input: `"abc"."cde"`, Expected: sdk.NewDatabaseObjectIdentifier(`abc`, `cde`)},
		{Input: `"ab.c"."cd.e"`, Expected: sdk.NewDatabaseObjectIdentifier(`ab.c`, `cd.e`)},
		{Input: `"a""bc"."cd""e"`, Expected: sdk.NewDatabaseObjectIdentifier(`a"bc`, `cd"e`)},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Parsing database object identifier with input: "%s"`, testCase.Input), func(t *testing.T) {
			id, err := ParseDatabaseObjectIdentifier(testCase.Input)

			if testCase.Error != "" {
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.Equal(t, testCase.Expected, id)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_ParseSchemaObjectIdentifier(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected sdk.SchemaObjectIdentifier
		Error    string
	}{
		{Input: ``, Error: "incompatible identifier: "},
		{Input: "a\nb.cde.efg", Error: "unable to read identifier: a\nb.cde.efg, err = record on line 2: wrong number of fields"},
		{Input: `a"b.cde.efg`, Error: "unable to read identifier: a\"b.cde.efg, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc.cde.efg.ghi`, Error: `unexpected number of parts 4 in identifier abc.cde.efg.ghi, expected 3 in a form of "<database_name>.<schema_name>.<schema_object_name>"`},
		{Input: `abc.cde`, Error: `unexpected number of parts 2 in identifier abc.cde, expected 3 in a form of "<database_name>.<schema_name>.<schema_object_name>"`},
		{Input: `""."".""`, Expected: sdk.NewSchemaObjectIdentifier(``, ``, ``)},
		{Input: `""""."""".""""`, Expected: sdk.NewSchemaObjectIdentifier(`"`, `"`, `"`)},
		{Input: `abc.cde.efg`, Expected: sdk.NewSchemaObjectIdentifier(`abc`, `cde`, `efg`)},
		{Input: `"abc"."cde"."efg"`, Expected: sdk.NewSchemaObjectIdentifier(`abc`, `cde`, `efg`)},
		{Input: `"ab.c"."cd.e"."ef.g"`, Expected: sdk.NewSchemaObjectIdentifier(`ab.c`, `cd.e`, `ef.g`)},
		{Input: `"a""bc"."cd""e"."ef""g"`, Expected: sdk.NewSchemaObjectIdentifier(`a"bc`, `cd"e`, `ef"g`)},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Parsing schema object identifier with input: "%s"`, testCase.Input), func(t *testing.T) {
			id, err := ParseSchemaObjectIdentifier(testCase.Input)

			if testCase.Error != "" {
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.Equal(t, testCase.Expected, id)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_ParseTableColumnIdentifier(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected sdk.TableColumnIdentifier
		Error    string
	}{
		{Input: ``, Error: "incompatible identifier: "},
		{Input: "a\nb.cde.efg.ghi", Error: "unable to read identifier: a\nb.cde.efg.ghi, err = record on line 2: wrong number of fields"},
		{Input: `a"b.cde.efg.ghi`, Error: "unable to read identifier: a\"b.cde.efg.ghi, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc.cde.efg.ghi.ijk`, Error: `unexpected number of parts 5 in identifier abc.cde.efg.ghi.ijk, expected 4 in a form of "<database_name>.<schema_name>.<table_name>.<table_column_name>"`},
		{Input: `abc.cde`, Error: `unexpected number of parts 2 in identifier abc.cde, expected 4 in a form of "<database_name>.<schema_name>.<table_name>.<table_column_name>"`},
		{Input: `"".""."".""`, Expected: sdk.NewTableColumnIdentifier(``, ``, ``, ``)},
		{Input: `"""".""""."""".""""`, Expected: sdk.NewTableColumnIdentifier(`"`, `"`, `"`, `"`)},
		{Input: `abc.cde.efg.ghi`, Expected: sdk.NewTableColumnIdentifier(`abc`, `cde`, `efg`, `ghi`)},
		{Input: `"abc"."cde"."efg"."ghi"`, Expected: sdk.NewTableColumnIdentifier(`abc`, `cde`, `efg`, `ghi`)},
		{Input: `"ab.c"."cd.e"."ef.g"."gh.i"`, Expected: sdk.NewTableColumnIdentifier(`ab.c`, `cd.e`, `ef.g`, `gh.i`)},
		{Input: `"a""bc"."cd""e"."ef""g"."gh""i"`, Expected: sdk.NewTableColumnIdentifier(`a"bc`, `cd"e`, `ef"g`, `gh"i`)},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Parsing table column identifier with input: "%s"`, testCase.Input), func(t *testing.T) {
			id, err := ParseTableColumnIdentifier(testCase.Input)

			if testCase.Error != "" {
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.Equal(t, testCase.Expected, id)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_ParseAccountIdentifier(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected sdk.AccountIdentifier
		Error    string
	}{
		{Input: ``, Error: "incompatible identifier: "},
		{Input: "a\nb.cde", Error: "unable to read identifier: a\nb.cde, err = record on line 2: wrong number of fields"},
		{Input: `a"b.cde`, Error: "unable to read identifier: a\"b.cde, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc.cde.efg`, Error: `unexpected number of parts 3 in identifier abc.cde.efg, expected 2 in a form of "<organization_name>.<account_name>"`},
		{Input: `abc`, Error: `unexpected number of parts 1 in identifier abc, expected 2 in a form of "<organization_name>.<account_name>"`},
		{Input: `"".""`, Expected: sdk.NewAccountIdentifier(``, ``)},
		{Input: `"""".""""`, Expected: sdk.NewAccountIdentifier(`"`, `"`)},
		{Input: `abc.cde`, Expected: sdk.NewAccountIdentifier(`abc`, `cde`)},
		{Input: `"abc"."cde"`, Expected: sdk.NewAccountIdentifier(`abc`, `cde`)},
		{Input: `"ab.c"."cd.e"`, Expected: sdk.NewAccountIdentifier(`ab.c`, `cd.e`)},
		{Input: `"a""bc"."cd""e"`, Expected: sdk.NewAccountIdentifier(`a"bc`, `cd"e`)},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Parsing account identifier with input: "%s"`, testCase.Input), func(t *testing.T) {
			id, err := ParseAccountIdentifier(testCase.Input)

			if testCase.Error != "" {
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.Equal(t, testCase.Expected, id)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_ParseExternalObjectIdentifier(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected sdk.ExternalObjectIdentifier
		Error    string
	}{
		{Input: ``, Error: "incompatible identifier: "},
		{Input: "a\nb.cde.efg", Error: "unable to read identifier: a\nb.cde.efg, err = record on line 2: wrong number of fields"},
		{Input: `a"b.cde.efg`, Error: "unable to read identifier: a\"b.cde.efg, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc.cde.efg.ghi`, Error: `unexpected number of parts 4 in identifier abc.cde.efg.ghi, expected 3 in a form of "<organization_name>.<account_name>.<external_object_name>"`},
		{Input: `abc.cde`, Error: `unexpected number of parts 2 in identifier abc.cde, expected 3 in a form of "<organization_name>.<account_name>.<external_object_name>"`},
		{Input: `""."".""`, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(``, ``), sdk.NewAccountObjectIdentifier(``))},
		{Input: `""""."""".""""`, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`"`, `"`), sdk.NewAccountObjectIdentifier(`"`))},
		{Input: `abc.cde.efg`, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`abc`, `cde`), sdk.NewAccountObjectIdentifier(`efg`))},
		{Input: `"abc"."cde"."efg"`, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`abc`, `cde`), sdk.NewAccountObjectIdentifier(`efg`))},
		{Input: `"ab.c"."cd.e"."ef.g"`, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`ab.c`, `cd.e`), sdk.NewAccountObjectIdentifier(`ef.g`))},
		{Input: `"a""bc"."cd""e"."ef""g"`, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`a"bc`, `cd"e`), sdk.NewAccountObjectIdentifier(`ef"g`))},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Parsing external object identifier with input: "%s"`, testCase.Input), func(t *testing.T) {
			id, err := ParseExternalObjectIdentifier(testCase.Input)

			if testCase.Error != "" {
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.Equal(t, testCase.Expected, id)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Encoding_And_Parsing_Of_ResourceIdentifier(t *testing.T) {
	testCases := []struct {
		Input                 []string
		Expected              string
		ExpectedAfterDecoding []string
	}{
		{Input: []string{sdk.NewSchemaObjectIdentifier("a", "b", "c").FullyQualifiedName(), "info"}, Expected: `"a"."b"."c"|info`},
		{Input: []string{}, Expected: ``},
		{Input: []string{"", "", ""}, Expected: `||`},
		{Input: []string{"a", "b", "c"}, Expected: `a|b|c`},
		// If one of the parts contains a separator sign (pipe in this case),
		// we can end up with more parts than we started with.
		{Input: []string{"a", "b", "c|d"}, Expected: `a|b|c|d`, ExpectedAfterDecoding: []string{"a", "b", "c", "d"}},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Encoding and parsing %s resource identifier`, testCase.Input), func(t *testing.T) {
			encodedIdentifier := EncodeResourceIdentifier(testCase.Input...)
			assert.Equal(t, testCase.Expected, encodedIdentifier)

			parsedIdentifier := ParseResourceIdentifier(encodedIdentifier)
			if testCase.ExpectedAfterDecoding != nil {
				assert.Equal(t, testCase.ExpectedAfterDecoding, parsedIdentifier)
			} else {
				assert.Equal(t, testCase.Input, parsedIdentifier)
			}
		})
	}
}
