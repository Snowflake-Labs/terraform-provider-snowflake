package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_suppressIdentifierQuoting(t *testing.T) {
	firstId := "a.b.c"
	firstIdQuoted := "\"a\".b.\"c\""
	secondId := "d.e.f"
	incorrectId := "a.b.c.d.e.f"

	t.Run("old identifier with too many parts", func(t *testing.T) {
		result := suppressIdentifierQuoting("", incorrectId, firstId, nil)
		require.False(t, result)
	})

	t.Run("new identifier with too many parts", func(t *testing.T) {
		result := suppressIdentifierQuoting("", firstId, incorrectId, nil)
		require.False(t, result)
	})

	t.Run("old identifier empty", func(t *testing.T) {
		result := suppressIdentifierQuoting("", "", firstId, nil)
		require.False(t, result)
	})

	t.Run("new identifier empty", func(t *testing.T) {
		result := suppressIdentifierQuoting("", firstId, "", nil)
		require.False(t, result)
	})

	t.Run("identifiers the same", func(t *testing.T) {
		result := suppressIdentifierQuoting("", firstId, firstId, nil)
		require.True(t, result)
	})

	t.Run("identifiers the same (but different quoting)", func(t *testing.T) {
		result := suppressIdentifierQuoting("", firstId, firstIdQuoted, nil)
		require.True(t, result)
	})

	t.Run("identifiers different", func(t *testing.T) {
		result := suppressIdentifierQuoting("", firstId, secondId, nil)
		require.False(t, result)
	})
}
