package sdk

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserCreate(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &CreateUserOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "CREATE USER"
		assert.Equal(t, expected, actual)
	})

	t.Run("with complete options", func(t *testing.T) {
		tags := []TagAssociation{
			{
				Name:  NewSchemaObjectIdentifier("db", "schema", "tag1"),
				Value: "v1",
			},
		}
		password := random.String()
		loginName := random.String()

		opts := &CreateUserOptions{
			OrReplace:   Bool(true),
			name:        id,
			IfNotExists: Bool(true),
			ObjectProperties: &UserObjectProperties{
				Password:  &password,
				LoginName: &loginName,
			},
			ObjectParameters: &UserObjectParameters{
				EnableUnredactedQuerySyntaxError: Bool(true),
			},
			SessionParameters: &SessionParameters{
				Autocommit: Bool(true),
			},
			With: Bool(true),
			Tags: tags,
		}

		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`CREATE OR REPLACE USER IF NOT EXISTS %s PASSWORD = '%s' LOGIN_NAME = '%s' ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = true AUTOCOMMIT = true WITH TAG ("db"."schema"."tag1" = 'v1')`, id.FullyQualifiedName(), password, loginName)
		assert.Equal(t, expected, actual)
	})
}

func TestUserAlter(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &AlterUserOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "ALTER USER"
		assert.Equal(t, expected, actual)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER USER %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with setting a policy", func(t *testing.T) {
		passwordPolicy := "PASSWORD_POLICY1"
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				PasswordPolicy: String(passwordPolicy),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER USER %s SET PASSWORD POLICY = %s", id.FullyQualifiedName(), passwordPolicy)
		assert.Equal(t, expected, actual)
	})

	t.Run("with setting tags", func(t *testing.T) {
		tags := []TagAssociation{
			{
				Name:  NewSchemaObjectIdentifier("db", "schema", "tag1"),
				Value: "v1",
			},
			{
				Name:  NewSchemaObjectIdentifier("db", "schema", "tag2"),
				Value: "v2",
			},
		}
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				Tags: tags,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`ALTER USER %s SET TAG ("db"."schema"."tag1" = 'v1', "db"."schema"."tag2" = 'v2')`, id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("with setting properties and parameters", func(t *testing.T) {
		password := random.String()
		objectProperties := UserObjectProperties{
			Password:             &password,
			DefaultSeconaryRoles: &SecondaryRoles{Roles: []SecondaryRole{{Value: "ALL"}}},
		}
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				ObjectProperties: &objectProperties,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER USER %s SET PASSWORD = '%s' DEFAULT_SECONDARY_ROLES = ( 'ALL' )", id.FullyQualifiedName(), password)
		assert.Equal(t, expected, actual)

		objectParameters := UserObjectParameters{
			EnableUnredactedQuerySyntaxError: Bool(true),
		}

		opts = &AlterUserOptions{
			name: id,
			Set: &UserSet{
				ObjectParameters: &objectParameters,
			},
		}
		actual, err = structToSQL(opts)
		require.NoError(t, err)
		expected = fmt.Sprintf("ALTER USER %s SET ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = true", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)

		sessionParameters := SessionParameters{
			Autocommit: Bool(true),
		}
		opts = &AlterUserOptions{
			name: id,
			Set: &UserSet{
				SessionParameters: &sessionParameters,
			},
		}
		actual, err = structToSQL(opts)
		require.NoError(t, err)
		expected = fmt.Sprintf("ALTER USER %s SET AUTOCOMMIT = true", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("reset password", func(t *testing.T) {
		id := RandomAccountObjectIdentifier()
		opts := &AlterUserOptions{
			name:          id,
			ResetPassword: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER USER %s RESET PASSWORD", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("abort all queries", func(t *testing.T) {
		id := RandomAccountObjectIdentifier()
		opts := &AlterUserOptions{
			name:            id,
			AbortAllQueries: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER USER %s ABORT ALL QUERIES", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("rename", func(t *testing.T) {
		newID := NewAccountObjectIdentifier(random.String())
		opts := &AlterUserOptions{
			name:    id,
			NewName: newID,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER USER %s RENAME TO %s", id.FullyQualifiedName(), newID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("with adding delegated authorization of role", func(t *testing.T) {
		role := "ROLE1"
		integration := "INTEGRATION1"
		opts := &AlterUserOptions{
			name: id,
			AddDelegatedAuthorization: &AddDelegatedAuthorization{
				Role:        role,
				Integration: integration,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER USER %s ADD DELEGATED AUTHORIZATION OF ROLE %s TO SECURITY INTEGRATION %s", id.FullyQualifiedName(), role, integration)
		assert.Equal(t, expected, actual)
	})

	t.Run("with unsetting tags", func(t *testing.T) {
		tag1 := "USER_TAG1"
		tag2 := "USER_TAG2"
		opts := &AlterUserOptions{
			name: id,
			Unset: &UserUnset{
				Tags: &[]string{tag1, tag2},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER USER %s UNSET TAG %s, %s", id.FullyQualifiedName(), tag1, tag2)
		assert.Equal(t, expected, actual)
	})
	t.Run("with unsetting properties", func(t *testing.T) {
		objectProperties := UserObjectPropertiesUnset{
			Password: Bool(true),
		}
		opts := &AlterUserOptions{
			name: id,
			Unset: &UserUnset{
				ObjectProperties: &objectProperties,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER USER %s UNSET PASSWORD", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with unsetting a policy", func(t *testing.T) {
		sessionPolicy := "SESSION_POLICY1"
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				SessionPolicy: String(sessionPolicy),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER USER %s SET SESSION POLICY = %s", id.FullyQualifiedName(), sessionPolicy)
		assert.Equal(t, expected, actual)
	})

	t.Run("with removing delegated authorization of role", func(t *testing.T) {
		role := "ROLE1"
		integration := "INTEGRATION1"
		opts := &AlterUserOptions{
			name: id,
			RemoveDelegatedAuthorization: &RemoveDelegatedAuthorization{
				Role:        &role,
				Integration: integration,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER USER %s REMOVE DELEGATED AUTHORIZATION OF ROLE %s FROM SECURITY INTEGRATION %s", id.FullyQualifiedName(), role, integration)
		assert.Equal(t, expected, actual)
	})
}

func TestUserDrop(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &DropUserOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "DROP USER"
		assert.Equal(t, expected, actual)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &DropUserOptions{
			name: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("DROP USER %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}

func TestUserShow(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &ShowUserOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "SHOW USERS"
		assert.Equal(t, expected, actual)
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowUserOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW USERS LIKE '%s'", id.Name())
		assert.Equal(t, expected, actual)
	})

	t.Run("with like and from", func(t *testing.T) {
		fromPatern := random.String()
		opts := &ShowUserOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			From: &fromPatern,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW USERS LIKE '%s' FROM '%s'", id.Name(), fromPatern)
		assert.Equal(t, expected, actual)
	})

	t.Run("with like and limit", func(t *testing.T) {
		limit := 5
		opts := &ShowUserOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			Limit: &limit,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW USERS LIKE '%s' LIMIT %v", id.Name(), limit)
		assert.Equal(t, expected, actual)
	})

	t.Run("with starts with and from", func(t *testing.T) {
		fromPattern := random.String()
		startsWithPattern := random.String()

		opts := &ShowUserOptions{
			StartsWith: &startsWithPattern,
			From:       &fromPattern,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW USERS STARTS WITH '%s' FROM '%s'", startsWithPattern, fromPattern)
		assert.Equal(t, expected, actual)
	})
}

func TestUserDescribe(t *testing.T) {
	id := RandomAccountObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &describeUserOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "DESCRIBE USER"
		assert.Equal(t, expected, actual)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &describeUserOptions{
			name: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("DESCRIBE USER %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}
