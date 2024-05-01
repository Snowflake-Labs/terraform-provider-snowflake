package sdk

import (
	"testing"
)

// TODO: add more tests
func TestSetObjectParameterOnObject(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOpts := func() *setParameterOnObject {
		return &setParameterOnObject{
			objectType:       ObjectTypeUser,
			objectIdentifier: id,
			parameterKey:     "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR",
			parameterValue:   "TRUE",
		}
	}

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = TRUE", id.FullyQualifiedName())
	})
}
