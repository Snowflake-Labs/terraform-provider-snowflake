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

func TestUnSetObjectParameterNetworkPolicyOnAccount(t *testing.T) {
	opts := &AlterAccountOptions{
		Unset: &AccountUnset{
			Parameters: &AccountLevelParametersUnset{
				ObjectParameters: &ObjectParametersUnset{
					NetworkPolicy: Bool(true),
				},
			},
		},
	}
	t.Run("Unset Account Network Policy", func(t *testing.T) {
		assertOptsValidAndSQLEquals(t, opts, "ALTER ACCOUNT UNSET NETWORK_POLICY")
	})
}

func TestUnSetObjectParameterNetworkPolicyOnUser(t *testing.T) {
	opts := &AlterUserOptions{
		name: NewAccountObjectIdentifierFromFullyQualifiedName("TEST_USER"),
		Unset: &UserUnset{
			ObjectParameters: &UserObjectParametersUnset{
				NetworkPolicy: Bool(true),
			},
		},
	}
	t.Run("Unset User Network Policy", func(t *testing.T) {
		assertOptsValidAndSQLEquals(t, opts, `ALTER USER "TEST_USER" UNSET NETWORK_POLICY`)
	})
}
