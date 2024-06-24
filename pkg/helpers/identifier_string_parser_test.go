package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO [SNOW-999049]: add more fancy cases
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
}
