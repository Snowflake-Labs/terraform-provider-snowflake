package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrantPrivilegeToShare(t *testing.T) {
	id := randomAccountObjectIdentifier(t)
	t.Run("on database", func(t *testing.T) {
		otherID := randomAccountObjectIdentifier(t)
		opts := &grantPrivilegeToShareOptions{
			objectPrivilege: PrivilegeUsage,
			On: &GrantPrivilegeToShareOn{
				Database: otherID,
			},
			to: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("GRANT USAGE ON DATABASE %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on schema", func(t *testing.T) {
		otherID := randomSchemaIdentifier(t)
		opts := &grantPrivilegeToShareOptions{
			objectPrivilege: PrivilegeUsage,
			On: &GrantPrivilegeToShareOn{
				Schema: otherID,
			},
			to: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("GRANT USAGE ON SCHEMA %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on table", func(t *testing.T) {
		otherID := randomSchemaObjectIdentifier(t)
		opts := &grantPrivilegeToShareOptions{
			objectPrivilege: PrivilegeUsage,
			On: &GrantPrivilegeToShareOn{
				Table: &OnTable{
					Name: otherID,
				},
			},
			to: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("GRANT USAGE ON TABLE %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on all tables", func(t *testing.T) {
		otherID := randomSchemaIdentifier(t)
		opts := &grantPrivilegeToShareOptions{
			objectPrivilege: PrivilegeUsage,
			On: &GrantPrivilegeToShareOn{
				Table: &OnTable{
					AllInSchema: otherID,
				},
			},
			to: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("GRANT USAGE ON ALL TABLES IN SCHEMA %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on view", func(t *testing.T) {
		otherID := randomSchemaObjectIdentifier(t)
		opts := &grantPrivilegeToShareOptions{
			objectPrivilege: PrivilegeUsage,
			On: &GrantPrivilegeToShareOn{
				View: otherID,
			},
			to: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("GRANT USAGE ON VIEW %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}

func TestRevokePrivilegeFromShare(t *testing.T) {
	id := randomAccountObjectIdentifier(t)
	t.Run("on database", func(t *testing.T) {
		otherID := randomAccountObjectIdentifier(t)
		opts := &revokePrivilegeFromShareOptions{
			objectPrivilege: PrivilegeUsage,
			On: &RevokePrivilegeFromShareOn{
				Database: otherID,
			},
			from: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("REVOKE USAGE ON DATABASE %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on schema", func(t *testing.T) {
		otherID := randomSchemaIdentifier(t)
		opts := &revokePrivilegeFromShareOptions{
			objectPrivilege: PrivilegeUsage,
			On: &RevokePrivilegeFromShareOn{
				Schema: otherID,
			},
			from: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("REVOKE USAGE ON SCHEMA %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on table", func(t *testing.T) {
		otherID := randomSchemaObjectIdentifier(t)
		opts := &revokePrivilegeFromShareOptions{
			objectPrivilege: PrivilegeUsage,
			On: &RevokePrivilegeFromShareOn{
				Table: &OnTable{
					Name: otherID,
				},
			},
			from: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("REVOKE USAGE ON TABLE %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on all tables", func(t *testing.T) {
		otherID := randomSchemaIdentifier(t)
		opts := &revokePrivilegeFromShareOptions{
			objectPrivilege: PrivilegeUsage,
			On: &RevokePrivilegeFromShareOn{
				Table: &OnTable{
					AllInSchema: otherID,
				},
			},
			from: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("REVOKE USAGE ON ALL TABLES IN SCHEMA %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on view", func(t *testing.T) {
		otherID := randomSchemaObjectIdentifier(t)
		opts := &revokePrivilegeFromShareOptions{
			objectPrivilege: PrivilegeUsage,
			On: &RevokePrivilegeFromShareOn{
				View: &OnView{
					Name: otherID,
				},
			},
			from: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("REVOKE USAGE ON VIEW %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("on all views", func(t *testing.T) {
		otherID := randomSchemaIdentifier(t)
		opts := &revokePrivilegeFromShareOptions{
			objectPrivilege: PrivilegeUsage,
			On: &RevokePrivilegeFromShareOn{
				View: &OnView{
					AllInSchema: otherID,
				},
			},
			from: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("REVOKE USAGE ON ALL VIEWS IN SCHEMA %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}

func TestGrantShow(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		opts := &ShowGrantOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "SHOW GRANTS"
		assert.Equal(t, expected, actual)
	})

	t.Run("on account", func(t *testing.T) {
		opts := &ShowGrantOptions{
			On: &ShowGrantsOn{
				Account: Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "SHOW GRANTS ON ACCOUNT"
		assert.Equal(t, expected, actual)
	})

	t.Run("on database", func(t *testing.T) {
		dbID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			On: &ShowGrantsOn{
				Object: &Object{
					ObjectType: ObjectTypeDatabase,
					Name:       dbID,
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW GRANTS ON DATABASE %s", dbID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("to role", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			To: &ShowGrantsTo{
				Role: roleID,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW GRANTS TO ROLE %s", roleID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("to user", func(t *testing.T) {
		userID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			To: &ShowGrantsTo{
				User: userID,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW GRANTS TO USER %s", userID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("to share", func(t *testing.T) {
		shareID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			To: &ShowGrantsTo{
				Share: shareID,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW GRANTS TO SHARE %s", shareID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("of role", func(t *testing.T) {
		roleID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			Of: &ShowGrantsOf{
				Role: roleID,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW GRANTS OF ROLE %s", roleID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("of share", func(t *testing.T) {
		shareID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantOptions{
			Of: &ShowGrantsOf{
				Share: shareID,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW GRANTS OF SHARE %s", shareID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}
