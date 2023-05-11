package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrantPrivilegeToShare(t *testing.T) {
	builder := testBuilder(t)
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
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
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
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
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
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
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
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
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
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("GRANT USAGE ON VIEW %s TO SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}

func TestRevokePrivilegeFromShare(t *testing.T) {
	builder := testBuilder(t)
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
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
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
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
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
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
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
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
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
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
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
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("REVOKE USAGE ON ALL VIEWS IN SCHEMA %s FROM SHARE %s", otherID.FullyQualifiedName(), id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}

func TestGrantShow(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		builder := testBuilder(t)
		opts := &ShowGrantsOptions{}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := "SHOW GRANTS"
		assert.Equal(t, expected, actual)
	})

	t.Run("on account", func(t *testing.T) {
		builder := testBuilder(t)
		opts := &ShowGrantsOptions{
			On: &ShowGrantsOn{
				Account: Bool(true),
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := "SHOW GRANTS ON ACCOUNT"
		assert.Equal(t, expected, actual)
	})

	t.Run("on database", func(t *testing.T) {
		builder := testBuilder(t)
		dbID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantsOptions{
			On: &ShowGrantsOn{
				Object: &Object{
					ObjectType: ObjectTypeDatabase,
					Name:       dbID,
				},
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("SHOW GRANTS ON DATABASE %s", dbID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("to role", func(t *testing.T) {
		builder := testBuilder(t)
		roleID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantsOptions{
			To: &ShowGrantsTo{
				Role: roleID,
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("SHOW GRANTS TO ROLE %s", roleID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("to user", func(t *testing.T) {
		builder := testBuilder(t)
		userID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantsOptions{
			To: &ShowGrantsTo{
				User: userID,
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("SHOW GRANTS TO USER %s", userID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("to share", func(t *testing.T) {
		builder := testBuilder(t)
		shareID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantsOptions{
			To: &ShowGrantsTo{
				Share: shareID,
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("SHOW GRANTS TO SHARE %s", shareID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("of role", func(t *testing.T) {
		builder := testBuilder(t)
		roleID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantsOptions{
			Of: &ShowGrantsOf{
				Role: roleID,
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("SHOW GRANTS OF ROLE %s", roleID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("of share", func(t *testing.T) {
		builder := testBuilder(t)
		shareID := randomAccountObjectIdentifier(t)
		opts := &ShowGrantsOptions{
			Of: &ShowGrantsOf{
				Share: shareID,
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("SHOW GRANTS OF SHARE %s", shareID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}
