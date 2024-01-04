package sdk

import (
	"log"
	"testing"
)

// TODO: add more tests
func TestSetObjectParameterOnObject(t *testing.T) {
	id := RandomAccountObjectIdentifier()

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

func (o ObjectType) Check() bool {
	var m map[ObjectType]bool
	if _, ok := m[o]; ok {
		return true
	}
	return false
}

func Test(t *testing.T) {
	a := "abc"
	b := ObjectTypeDatabase

	log.Println(ObjectType(a).Check())
	log.Println(b.Check())
}
