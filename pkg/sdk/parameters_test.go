package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetObjectParameterOnObject(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	defaultOpts := func() *setParameterOnObject {
		return &setParameterOnObject{
			objectType:       ObjectTypeUser,
			objectIdentifier: id,
			parameterKey:     "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR",
			parameterValue:   "TRUE",
		}
	}

	t.Run("empty options", func(t *testing.T) {
		opts := defaultOpts()
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "ALTER USER " + id.FullyQualifiedName() + " SET ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = TRUE"
		assert.Equal(t, expected, actual)
	})
}
