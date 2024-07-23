package helpers

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"testing"

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

		_, err := ParseIdentifierString(input)

		require.ErrorContains(t, err, "unable to read identifier")
		require.ErrorContains(t, err, `bare " in non-quoted-field`)
	})

	t.Run("returns error for empty input", func(t *testing.T) {
		input := ""

		_, err := ParseIdentifierString(input)

		require.ErrorContains(t, err, "incompatible identifier")
	})

	t.Run("returns error for multiple lines", func(t *testing.T) {
		input := "abc\ndef"

		_, err := ParseIdentifierString(input)

		require.ErrorContains(t, err, "incompatible identifier")
	})

	t.Run("returns parts correctly without quoting", func(t *testing.T) {
		input := "abc.def"
		expected := []string{"abc", "def"}

		parts, err := ParseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("returns parts correctly with quoting", func(t *testing.T) {
		input := `"abc"."def"`
		expected := []string{"abc", "def"}

		parts, err := ParseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("returns parts correctly with mixed quoting", func(t *testing.T) {
		input := `"abc".def."ghi"`
		expected := []string{"abc", "def", "ghi"}

		parts, err := ParseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	// Quote inside must have a preceding quote (https://docs.snowflake.com/en/sql-reference/identifiers-syntax).
	t.Run("returns parts correctly with quote inside", func(t *testing.T) {
		input := `"ab""c".def`
		expected := []string{`ab"c`, "def"}

		parts, err := ParseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("returns parts correctly with dots inside", func(t *testing.T) {
		input := `"ab.c".def`
		expected := []string{`ab.c`, "def"}

		parts, err := ParseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("empty identifier", func(t *testing.T) {
		input := `""`
		expected := []string{""}

		parts, err := ParseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("handled correctly double quotes", func(t *testing.T) {
		input := `""."."".".".""."".""."".""."".""."""""`
		expected := []string{"", `.".`, `.".".".".".".".""`}

		parts, err := ParseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})
}

func Test_ParseAccountObjectIdentifier(t *testing.T) {
	testCases := []struct {
		Input    string
		Delim    rune
		Expected sdk.AccountObjectIdentifier
		Error    string
	}{
		{Input: ``, Delim: IdDelimiter, Error: "incompatible identifier: "},
		{Input: "a\nb", Delim: IdDelimiter, Error: "incompatible identifier: a\nb"},
		{Input: `a"b`, Delim: IdDelimiter, Error: "unable to read identifier: a\"b, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc.cde`, Delim: IdDelimiter, Error: `unexpected number of parts 2 in identifier abc.cde, expected 1 in a form of "<account_object_name>"`},
		{Input: `""`, Delim: IdDelimiter, Expected: sdk.NewAccountObjectIdentifier(``)},
		{Input: `""""`, Delim: IdDelimiter, Expected: sdk.NewAccountObjectIdentifier(`"`)},
		{Input: `abc`, Delim: IdDelimiter, Expected: sdk.NewAccountObjectIdentifier(`abc`)},
		{Input: `"abc"`, Delim: IdDelimiter, Expected: sdk.NewAccountObjectIdentifier(`abc`)},
		{Input: `"ab.c"`, Delim: IdDelimiter, Expected: sdk.NewAccountObjectIdentifier(`ab.c`)},
		{Input: `"a""bc"`, Delim: IdDelimiter, Expected: sdk.NewAccountObjectIdentifier(`a"bc`)},

		{Input: ``, Delim: ResourceIdDelimiter, Error: "incompatible identifier: "},
		{Input: "a\nb", Delim: ResourceIdDelimiter, Error: "incompatible identifier: a\nb"},
		{Input: `a"b`, Delim: ResourceIdDelimiter, Error: "unable to read identifier: a\"b, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc|cde`, Delim: ResourceIdDelimiter, Error: `unexpected number of parts 2 in identifier abc|cde, expected 1 in a form of "<account_object_name>"`},
		{Input: `""`, Delim: ResourceIdDelimiter, Expected: sdk.NewAccountObjectIdentifier(``)},
		{Input: `""""`, Delim: ResourceIdDelimiter, Expected: sdk.NewAccountObjectIdentifier(`"`)},
		{Input: `abc`, Delim: ResourceIdDelimiter, Expected: sdk.NewAccountObjectIdentifier(`abc`)},
		{Input: `"abc"`, Delim: ResourceIdDelimiter, Expected: sdk.NewAccountObjectIdentifier(`abc`)},
		{Input: `"ab|c"`, Delim: IdDelimiter, Expected: sdk.NewAccountObjectIdentifier(`ab|c`)},
		{Input: `"a""bc"`, Delim: ResourceIdDelimiter, Expected: sdk.NewAccountObjectIdentifier(`a"bc`)},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Parsing account object identifier with input: "%s" and delimiter %c`, testCase.Input, testCase.Delim), func(t *testing.T) {
			var id sdk.AccountObjectIdentifier
			var err error

			switch testCase.Delim {
			case IdDelimiter:
				id, err = ParseAccountObjectIdentifier(testCase.Input)
			case ResourceIdDelimiter:
				id, err = ParseAccountObjectResourceIdentifier(testCase.Input)
			default:
				t.Errorf("unexpected delimiter: %c", testCase.Delim)
			}

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
		Delim    rune
		Expected sdk.DatabaseObjectIdentifier
		Error    string
	}{
		{Input: ``, Delim: IdDelimiter, Error: "incompatible identifier: "},
		{Input: "a\nb.cde", Delim: IdDelimiter, Error: "unable to read identifier: a\nb.cde, err = record on line 2: wrong number of fields"},
		{Input: `a"b.cde`, Delim: IdDelimiter, Error: "unable to read identifier: a\"b.cde, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc.cde.efg`, Delim: IdDelimiter, Error: `unexpected number of parts 3 in identifier abc.cde.efg, expected 2 in a form of "<database_name>.<database_object_name>"`},
		{Input: `abc`, Delim: IdDelimiter, Error: `unexpected number of parts 1 in identifier abc, expected 2 in a form of "<database_name>.<database_object_name>"`},
		{Input: `"".""`, Delim: IdDelimiter, Expected: sdk.NewDatabaseObjectIdentifier(``, ``)},
		{Input: `"""".""""`, Delim: IdDelimiter, Expected: sdk.NewDatabaseObjectIdentifier(`"`, `"`)},
		{Input: `abc.cde`, Delim: IdDelimiter, Expected: sdk.NewDatabaseObjectIdentifier(`abc`, `cde`)},
		{Input: `"abc"."cde"`, Delim: IdDelimiter, Expected: sdk.NewDatabaseObjectIdentifier(`abc`, `cde`)},
		{Input: `"ab.c"."cd.e"`, Delim: IdDelimiter, Expected: sdk.NewDatabaseObjectIdentifier(`ab.c`, `cd.e`)},
		{Input: `"a""bc"."cd""e"`, Delim: IdDelimiter, Expected: sdk.NewDatabaseObjectIdentifier(`a"bc`, `cd"e`)},

		{Input: ``, Delim: ResourceIdDelimiter, Error: "incompatible identifier: "},
		{Input: "a\nb|cde", Delim: ResourceIdDelimiter, Error: "unable to read identifier: a\nb|cde, err = record on line 2: wrong number of fields"},
		{Input: `a"b|cde`, Delim: ResourceIdDelimiter, Error: "unable to read identifier: a\"b|cde, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc|cde|efg`, Delim: ResourceIdDelimiter, Error: `unexpected number of parts 3 in identifier abc|cde|efg, expected 2 in a form of "<database_name>|<database_object_name>"`},
		{Input: `abc`, Delim: ResourceIdDelimiter, Error: `unexpected number of parts 1 in identifier abc, expected 2 in a form of "<database_name>|<database_object_name>"`},
		{Input: `""|""`, Delim: ResourceIdDelimiter, Expected: sdk.NewDatabaseObjectIdentifier(``, ``)},
		{Input: `""""|""""`, Delim: ResourceIdDelimiter, Expected: sdk.NewDatabaseObjectIdentifier(`"`, `"`)},
		{Input: `abc|cde`, Delim: ResourceIdDelimiter, Expected: sdk.NewDatabaseObjectIdentifier(`abc`, `cde`)},
		{Input: `"abc"|"cde"`, Delim: ResourceIdDelimiter, Expected: sdk.NewDatabaseObjectIdentifier(`abc`, `cde`)},
		{Input: `"ab|c"|"cd|e"`, Delim: ResourceIdDelimiter, Expected: sdk.NewDatabaseObjectIdentifier(`ab|c`, `cd|e`)},
		{Input: `"a""bc"|"cd""e"`, Delim: ResourceIdDelimiter, Expected: sdk.NewDatabaseObjectIdentifier(`a"bc`, `cd"e`)},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Parsing database object identifier with input: "%s" and delimiter %c`, testCase.Input, testCase.Delim), func(t *testing.T) {
			var id sdk.DatabaseObjectIdentifier
			var err error

			switch testCase.Delim {
			case IdDelimiter:
				id, err = ParseDatabaseObjectIdentifier(testCase.Input)
			case ResourceIdDelimiter:
				id, err = ParseDatabaseObjectResourceIdentifier(testCase.Input)
			default:
				t.Errorf("unexpected delimiter: %c", testCase.Delim)
			}

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
		Delim    rune
		Expected sdk.SchemaObjectIdentifier
		Error    string
	}{
		{Input: ``, Delim: IdDelimiter, Error: "incompatible identifier: "},
		{Input: "a\nb.cde.efg", Delim: IdDelimiter, Error: "unable to read identifier: a\nb.cde.efg, err = record on line 2: wrong number of fields"},
		{Input: `a"b.cde.efg`, Delim: IdDelimiter, Error: "unable to read identifier: a\"b.cde.efg, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc.cde.efg.ghi`, Delim: IdDelimiter, Error: `unexpected number of parts 4 in identifier abc.cde.efg.ghi, expected 3 in a form of "<database_name>.<schema_name>.<schema_object_name>"`},
		{Input: `abc.cde`, Delim: IdDelimiter, Error: `unexpected number of parts 2 in identifier abc.cde, expected 3 in a form of "<database_name>.<schema_name>.<schema_object_name>"`},
		{Input: `""."".""`, Delim: IdDelimiter, Expected: sdk.NewSchemaObjectIdentifier(``, ``, ``)},
		{Input: `""""."""".""""`, Delim: IdDelimiter, Expected: sdk.NewSchemaObjectIdentifier(`"`, `"`, `"`)},
		{Input: `abc.cde.efg`, Delim: IdDelimiter, Expected: sdk.NewSchemaObjectIdentifier(`abc`, `cde`, `efg`)},
		{Input: `"abc"."cde"."efg"`, Delim: IdDelimiter, Expected: sdk.NewSchemaObjectIdentifier(`abc`, `cde`, `efg`)},
		{Input: `"ab.c"."cd.e"."ef.g"`, Delim: IdDelimiter, Expected: sdk.NewSchemaObjectIdentifier(`ab.c`, `cd.e`, `ef.g`)},
		{Input: `"a""bc"."cd""e"."ef""g"`, Delim: IdDelimiter, Expected: sdk.NewSchemaObjectIdentifier(`a"bc`, `cd"e`, `ef"g`)},

		{Input: ``, Delim: ResourceIdDelimiter, Error: "incompatible identifier: "},
		{Input: "a\nb|cde", Delim: ResourceIdDelimiter, Error: "unable to read identifier: a\nb|cde, err = record on line 2: wrong number of fields"},
		{Input: `a"b|cde`, Delim: ResourceIdDelimiter, Error: "unable to read identifier: a\"b|cde, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc|cde|efg|ghi`, Delim: ResourceIdDelimiter, Error: `unexpected number of parts 4 in identifier abc|cde|efg|ghi, expected 3 in a form of "<database_name>|<schema_name>|<schema_object_name>"`},
		{Input: `abc|cde`, Delim: ResourceIdDelimiter, Error: `unexpected number of parts 2 in identifier abc|cde, expected 3 in a form of "<database_name>|<schema_name>|<schema_object_name>"`},
		{Input: `""|""|""`, Delim: ResourceIdDelimiter, Expected: sdk.NewSchemaObjectIdentifier(``, ``, ``)},
		{Input: `""""|""""|""""`, Delim: ResourceIdDelimiter, Expected: sdk.NewSchemaObjectIdentifier(`"`, `"`, `"`)},
		{Input: `abc|cde|efg`, Delim: ResourceIdDelimiter, Expected: sdk.NewSchemaObjectIdentifier(`abc`, `cde`, `efg`)},
		{Input: `"abc"|"cde"|"efg"`, Delim: ResourceIdDelimiter, Expected: sdk.NewSchemaObjectIdentifier(`abc`, `cde`, `efg`)},
		{Input: `"ab|c"|"cd|e"|"ef|g"`, Delim: ResourceIdDelimiter, Expected: sdk.NewSchemaObjectIdentifier(`ab|c`, `cd|e`, `ef|g`)},
		{Input: `"a""bc"|"cd""e"|"ef""g"`, Delim: ResourceIdDelimiter, Expected: sdk.NewSchemaObjectIdentifier(`a"bc`, `cd"e`, `ef"g`)},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Parsing schema object identifier with input: "%s" and delimiter %c`, testCase.Input, testCase.Delim), func(t *testing.T) {
			var id sdk.SchemaObjectIdentifier
			var err error

			switch testCase.Delim {
			case IdDelimiter:
				id, err = ParseSchemaObjectIdentifier(testCase.Input)
			case ResourceIdDelimiter:
				id, err = ParseSchemaObjectResourceIdentifier(testCase.Input)
			default:
				t.Errorf("unexpected delimiter: %c", testCase.Delim)
			}

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
		Delim    rune
		Expected sdk.TableColumnIdentifier
		Error    string
	}{
		{Input: ``, Delim: IdDelimiter, Error: "incompatible identifier: "},
		{Input: "a\nb.cde.efg.ghi", Delim: IdDelimiter, Error: "unable to read identifier: a\nb.cde.efg.ghi, err = record on line 2: wrong number of fields"},
		{Input: `a"b.cde.efg.ghi`, Delim: IdDelimiter, Error: "unable to read identifier: a\"b.cde.efg.ghi, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc.cde.efg.ghi.ijk`, Delim: IdDelimiter, Error: `unexpected number of parts 5 in identifier abc.cde.efg.ghi.ijk, expected 4 in a form of "<database_name>.<schema_name>.<table_name>.<table_column_name>"`},
		{Input: `abc.cde`, Delim: IdDelimiter, Error: `unexpected number of parts 2 in identifier abc.cde, expected 4 in a form of "<database_name>.<schema_name>.<table_name>.<table_column_name>"`},
		{Input: `"".""."".""`, Delim: IdDelimiter, Expected: sdk.NewTableColumnIdentifier(``, ``, ``, ``)},
		{Input: `"""".""""."""".""""`, Delim: IdDelimiter, Expected: sdk.NewTableColumnIdentifier(`"`, `"`, `"`, `"`)},
		{Input: `abc.cde.efg.ghi`, Delim: IdDelimiter, Expected: sdk.NewTableColumnIdentifier(`abc`, `cde`, `efg`, `ghi`)},
		{Input: `"abc"."cde"."efg"."ghi"`, Delim: IdDelimiter, Expected: sdk.NewTableColumnIdentifier(`abc`, `cde`, `efg`, `ghi`)},
		{Input: `"ab.c"."cd.e"."ef.g"."gh.i"`, Delim: IdDelimiter, Expected: sdk.NewTableColumnIdentifier(`ab.c`, `cd.e`, `ef.g`, `gh.i`)},
		{Input: `"a""bc"."cd""e"."ef""g"."gh""i"`, Delim: IdDelimiter, Expected: sdk.NewTableColumnIdentifier(`a"bc`, `cd"e`, `ef"g`, `gh"i`)},

		{Input: ``, Delim: ResourceIdDelimiter, Error: "incompatible identifier: "},
		{Input: "a\nb|cde|efg|ghi", Delim: ResourceIdDelimiter, Error: "unable to read identifier: a\nb|cde|efg|ghi, err = record on line 2: wrong number of fields"},
		{Input: `a"b|cde|efg|ghi`, Delim: ResourceIdDelimiter, Error: "unable to read identifier: a\"b|cde|efg|ghi, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc|cde|efg|ghi|ijk`, Delim: ResourceIdDelimiter, Error: `unexpected number of parts 5 in identifier abc|cde|efg|ghi|ijk, expected 4 in a form of "<database_name>|<schema_name>|<table_name>|<table_column_name>"`},
		{Input: `abc|cde`, Delim: ResourceIdDelimiter, Error: `unexpected number of parts 2 in identifier abc|cde, expected 4 in a form of "<database_name>|<schema_name>|<table_name>|<table_column_name>"`},
		{Input: `""|""|""|""`, Delim: ResourceIdDelimiter, Expected: sdk.NewTableColumnIdentifier(``, ``, ``, ``)},
		{Input: `""""|""""|""""|""""`, Delim: ResourceIdDelimiter, Expected: sdk.NewTableColumnIdentifier(`"`, `"`, `"`, `"`)},
		{Input: `abc|cde|efg|ghi`, Delim: ResourceIdDelimiter, Expected: sdk.NewTableColumnIdentifier(`abc`, `cde`, `efg`, `ghi`)},
		{Input: `"abc"|"cde"|"efg"|"ghi"`, Delim: ResourceIdDelimiter, Expected: sdk.NewTableColumnIdentifier(`abc`, `cde`, `efg`, `ghi`)},
		{Input: `"ab|c"|"cd|e"|"ef|g"|"gh|i"`, Delim: ResourceIdDelimiter, Expected: sdk.NewTableColumnIdentifier(`ab|c`, `cd|e`, `ef|g`, `gh|i`)},
		{Input: `"a""bc"|"cd""e"|"ef""g"|"gh""i"`, Delim: ResourceIdDelimiter, Expected: sdk.NewTableColumnIdentifier(`a"bc`, `cd"e`, `ef"g`, `gh"i`)},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Parsing table column identifier with input: "%s" and delimiter %c`, testCase.Input, testCase.Delim), func(t *testing.T) {
			var id sdk.TableColumnIdentifier
			var err error

			switch testCase.Delim {
			case IdDelimiter:
				id, err = ParseTableColumnIdentifier(testCase.Input)
			case ResourceIdDelimiter:
				id, err = ParseTableColumnResourceIdentifier(testCase.Input)
			default:
				t.Errorf("unexpected delimiter: %c", testCase.Delim)
			}

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
		Delim    rune
		Expected sdk.AccountIdentifier
		Error    string
	}{
		{Input: ``, Delim: IdDelimiter, Error: "incompatible identifier: "},
		{Input: "a\nb.cde", Delim: IdDelimiter, Error: "unable to read identifier: a\nb.cde, err = record on line 2: wrong number of fields"},
		{Input: `a"b.cde`, Delim: IdDelimiter, Error: "unable to read identifier: a\"b.cde, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc.cde.efg`, Delim: IdDelimiter, Error: `unexpected number of parts 3 in identifier abc.cde.efg, expected 2 in a form of "<organization_name>.<account_name>"`},
		{Input: `abc`, Delim: IdDelimiter, Error: `unexpected number of parts 1 in identifier abc, expected 2 in a form of "<organization_name>.<account_name>"`},
		{Input: `"".""`, Delim: IdDelimiter, Expected: sdk.NewAccountIdentifier(``, ``)},
		{Input: `"""".""""`, Delim: IdDelimiter, Expected: sdk.NewAccountIdentifier(`"`, `"`)},
		{Input: `abc.cde`, Delim: IdDelimiter, Expected: sdk.NewAccountIdentifier(`abc`, `cde`)},
		{Input: `"abc"."cde"`, Delim: IdDelimiter, Expected: sdk.NewAccountIdentifier(`abc`, `cde`)},
		{Input: `"ab.c"."cd.e"`, Delim: IdDelimiter, Expected: sdk.NewAccountIdentifier(`ab.c`, `cd.e`)},
		{Input: `"a""bc"."cd""e"`, Delim: IdDelimiter, Expected: sdk.NewAccountIdentifier(`a"bc`, `cd"e`)},

		{Input: ``, Delim: ResourceIdDelimiter, Error: "incompatible identifier: "},
		{Input: "a\nb|cde", Delim: ResourceIdDelimiter, Error: "unable to read identifier: a\nb|cde, err = record on line 2: wrong number of fields"},
		{Input: `a"b|cde`, Delim: ResourceIdDelimiter, Error: "unable to read identifier: a\"b|cde, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc|cde|efg`, Delim: ResourceIdDelimiter, Error: `unexpected number of parts 3 in identifier abc|cde|efg, expected 2 in a form of "<organization_name>|<account_name>"`},
		{Input: `abc`, Delim: ResourceIdDelimiter, Error: `unexpected number of parts 1 in identifier abc, expected 2 in a form of "<organization_name>|<account_name>"`},
		{Input: `""|""`, Delim: ResourceIdDelimiter, Expected: sdk.NewAccountIdentifier(``, ``)},
		{Input: `""""|""""`, Delim: ResourceIdDelimiter, Expected: sdk.NewAccountIdentifier(`"`, `"`)},
		{Input: `abc|cde`, Delim: ResourceIdDelimiter, Expected: sdk.NewAccountIdentifier(`abc`, `cde`)},
		{Input: `"abc"|"cde"`, Delim: ResourceIdDelimiter, Expected: sdk.NewAccountIdentifier(`abc`, `cde`)},
		{Input: `"ab|c"|"cd|e"`, Delim: ResourceIdDelimiter, Expected: sdk.NewAccountIdentifier(`ab|c`, `cd|e`)},
		{Input: `"a""bc"|"cd""e"`, Delim: ResourceIdDelimiter, Expected: sdk.NewAccountIdentifier(`a"bc`, `cd"e`)},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Parsing account identifier with input: "%s" and delimiter %c`, testCase.Input, testCase.Delim), func(t *testing.T) {
			var id sdk.AccountIdentifier
			var err error

			switch testCase.Delim {
			case IdDelimiter:
				id, err = ParseAccountIdentifier(testCase.Input)
			case ResourceIdDelimiter:
				id, err = ParseAccountResourceIdentifier(testCase.Input)
			default:
				t.Errorf("unexpected delimiter: %c", testCase.Delim)
			}

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
		Delim    rune
		Expected sdk.ExternalObjectIdentifier
		Error    string
	}{
		{Input: ``, Delim: IdDelimiter, Error: "incompatible identifier: "},
		{Input: "a\nb.cde.efg", Delim: IdDelimiter, Error: "unable to read identifier: a\nb.cde.efg, err = record on line 2: wrong number of fields"},
		{Input: `a"b.cde.efg`, Delim: IdDelimiter, Error: "unable to read identifier: a\"b.cde.efg, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc.cde.efg.ghi`, Delim: IdDelimiter, Error: `unexpected number of parts 4 in identifier abc.cde.efg.ghi, expected 3 in a form of "<organization_name>.<account_name>.<external_object_name>"`},
		{Input: `abc.cde`, Delim: IdDelimiter, Error: `unexpected number of parts 2 in identifier abc.cde, expected 3 in a form of "<organization_name>.<account_name>.<external_object_name>"`},
		{Input: `""."".""`, Delim: IdDelimiter, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(``, ``), sdk.NewAccountObjectIdentifier(``))},
		{Input: `""""."""".""""`, Delim: IdDelimiter, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`"`, `"`), sdk.NewAccountObjectIdentifier(`"`))},
		{Input: `abc.cde.efg`, Delim: IdDelimiter, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`abc`, `cde`), sdk.NewAccountObjectIdentifier(`efg`))},
		{Input: `"abc"."cde"."efg"`, Delim: IdDelimiter, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`abc`, `cde`), sdk.NewAccountObjectIdentifier(`efg`))},
		{Input: `"ab.c"."cd.e"."ef.g"`, Delim: IdDelimiter, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`ab.c`, `cd.e`), sdk.NewAccountObjectIdentifier(`ef.g`))},
		{Input: `"a""bc"."cd""e"."ef""g"`, Delim: IdDelimiter, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`a"bc`, `cd"e`), sdk.NewAccountObjectIdentifier(`ef"g`))},

		{Input: ``, Delim: ResourceIdDelimiter, Error: "incompatible identifier: "},
		{Input: "a\nb|cde", Delim: ResourceIdDelimiter, Error: "unable to read identifier: a\nb|cde, err = record on line 2: wrong number of fields"},
		{Input: `a"b|cde`, Delim: ResourceIdDelimiter, Error: "unable to read identifier: a\"b|cde, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: `abc|cde|efg|ghi`, Delim: ResourceIdDelimiter, Error: `unexpected number of parts 4 in identifier abc|cde|efg|ghi, expected 3 in a form of "<organization_name>|<account_name>|<external_object_name>"`},
		{Input: `abc|cde`, Delim: ResourceIdDelimiter, Error: `unexpected number of parts 2 in identifier abc|cde, expected 3 in a form of "<organization_name>|<account_name>|<external_object_name>"`},
		{Input: `""|""|""`, Delim: ResourceIdDelimiter, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(``, ``), sdk.NewAccountObjectIdentifier(``))},
		{Input: `""""|""""|""""`, Delim: ResourceIdDelimiter, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`"`, `"`), sdk.NewAccountObjectIdentifier(`"`))},
		{Input: `abc|cde|efg`, Delim: ResourceIdDelimiter, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`abc`, `cde`), sdk.NewAccountObjectIdentifier(`efg`))},
		{Input: `"abc"|"cde"|"efg"`, Delim: ResourceIdDelimiter, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`abc`, `cde`), sdk.NewAccountObjectIdentifier(`efg`))},
		{Input: `"ab|c"|"cd|e"|"ef|g"`, Delim: ResourceIdDelimiter, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`ab|c`, `cd|e`), sdk.NewAccountObjectIdentifier(`ef|g`))},
		{Input: `"a""bc"|"cd""e"|"ef""g"`, Delim: ResourceIdDelimiter, Expected: sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier(`a"bc`, `cd"e`), sdk.NewAccountObjectIdentifier(`ef"g`))},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Parsing external object identifier with input: "%s" and delimiter %c`, testCase.Input, testCase.Delim), func(t *testing.T) {
			var id sdk.ExternalObjectIdentifier
			var err error

			switch testCase.Delim {
			case IdDelimiter:
				id, err = ParseExternalObjectIdentifier(testCase.Input)
			case ResourceIdDelimiter:
				id, err = ParseExternalObjectResourceIdentifier(testCase.Input)
			default:
				t.Errorf("unexpected delimiter: %c", testCase.Delim)
			}

			if testCase.Error != "" {
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.Equal(t, testCase.Expected, id)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_EncodeResourceIdentifier(t *testing.T) {
	testCases := []struct {
		Input    []string
		Expected string
	}{
		{[]string{sdk.NewSchemaObjectIdentifier("a", "b", "c").FullyQualifiedName(), "info"}, `"a"."b"."c"|info`},
		{[]string{}, ``},
		{[]string{"", "", ""}, `||`},
		{[]string{"a", "b", "c"}, `a|b|c`},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Encoding %s to resource identifier`, testCase.Input), func(t *testing.T) {
			assert.Equal(t, testCase.Expected, EncodeResourceIdentifier(testCase.Input...))
		})
	}
}
