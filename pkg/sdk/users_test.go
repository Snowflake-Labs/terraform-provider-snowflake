package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
)

func TestUserCreate(t *testing.T) {
	id := randomAccountObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &CreateUserOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("with complete options", func(t *testing.T) {
		tagId := randomSchemaObjectIdentifier()
		tags := []TagAssociation{
			{
				Name:  tagId,
				Value: "v1",
			},
		}
		password := random.Password()
		loginName := random.String()
		defaultRoleId := randomAccountObjectIdentifier()
		defaultWarehouseId := randomAccountObjectIdentifier()
		var defaultNamespaceId ObjectIdentifier = randomDatabaseObjectIdentifier()

		opts := &CreateUserOptions{
			OrReplace:   Bool(true),
			name:        id,
			IfNotExists: Bool(true),
			ObjectProperties: &UserObjectProperties{
				Password:         &password,
				LoginName:        &loginName,
				DefaultRole:      Pointer(defaultRoleId),
				DefaultNamespace: Pointer(defaultNamespaceId),
				DefaultWarehouse: Pointer(defaultWarehouseId),
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

		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE USER IF NOT EXISTS %s PASSWORD = '%s' LOGIN_NAME = '%s' DEFAULT_WAREHOUSE = %s DEFAULT_NAMESPACE = %s DEFAULT_ROLE = %s ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = true AUTOCOMMIT = true WITH TAG (%s = 'v1')`, id.FullyQualifiedName(), password, loginName, defaultWarehouseId.FullyQualifiedName(), defaultNamespaceId.FullyQualifiedName(), defaultRoleId.FullyQualifiedName(), tagId.FullyQualifiedName())
	})
}

func TestUserAlter(t *testing.T) {
	id := randomAccountObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &AlterUserOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterUserOptions", "NewName", "ResetPassword", "AbortAllQueries", "AddDelegatedAuthorization", "RemoveDelegatedAuthorization", "Set", "Unset", "SetTag", "UnsetTag"))
	})

	t.Run("validation: no set", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
			Set:  &UserSet{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("UserSet", "PasswordPolicy", "SessionPolicy", "ObjectProperties", "ObjectParameters", "SessionParameters"))
	})

	t.Run("two sets", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				SessionParameters: &SessionParameters{AbortDetachedQuery: Bool(true)},
				ObjectParameters:  &UserObjectParameters{EnableUnredactedQuerySyntaxError: Bool(true)},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = true ABORT_DETACHED_QUERY = true", id.FullyQualifiedName())
	})

	t.Run("validation: no unset", func(t *testing.T) {
		opts := &AlterUserOptions{
			name:  id,
			Unset: &UserUnset{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("UserUnset", "PasswordPolicy", "SessionPolicy", "ObjectProperties", "ObjectParameters", "SessionParameters"))
	})

	t.Run("validation: two incompatible unsets", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
			Unset: &UserUnset{
				SessionParameters: &SessionParametersUnset{BinaryOutputFormat: Bool(true)},
				ObjectParameters:  &UserObjectParametersUnset{EnableUnredactedQuerySyntaxError: Bool(true)},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("UserUnset", "PasswordPolicy", "SessionPolicy", "ObjectProperties", "ObjectParameters", "SessionParameters"))
	})

	t.Run("with setting a policy", func(t *testing.T) {
		passwordPolicy := randomSchemaObjectIdentifier()
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				PasswordPolicy: &passwordPolicy,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET PASSWORD POLICY %s", id.FullyQualifiedName(), passwordPolicy.FullyQualifiedName())
	})

	t.Run("with setting tags", func(t *testing.T) {
		tagId1 := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
		tags := []TagAssociation{
			{
				Name:  tagId1,
				Value: "v1",
			},
			{
				Name:  tagId2,
				Value: "v2",
			},
		}
		opts := &AlterUserOptions{
			name:   id,
			SetTag: tags,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER USER %s SET TAG %s = 'v1', %s = 'v2'`, id.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName())
	})

	t.Run("with setting properties and parameters", func(t *testing.T) {
		password := random.Password()
		objectProperties := UserObjectProperties{
			Password:              &password,
			DefaultSecondaryRoles: &SecondaryRoles{Roles: []SecondaryRole{{Value: "ALL"}}},
		}
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				ObjectProperties: &objectProperties,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET PASSWORD = '%s' DEFAULT_SECONDARY_ROLES = ( 'ALL' )", id.FullyQualifiedName(), password)

		objectParameters := UserObjectParameters{
			EnableUnredactedQuerySyntaxError: Bool(true),
		}

		opts = &AlterUserOptions{
			name: id,
			Set: &UserSet{
				ObjectParameters: &objectParameters,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = true", id.FullyQualifiedName())

		sessionParameters := SessionParameters{
			Autocommit: Bool(true),
		}
		opts = &AlterUserOptions{
			name: id,
			Set: &UserSet{
				SessionParameters: &sessionParameters,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET AUTOCOMMIT = true", id.FullyQualifiedName())
	})

	t.Run("reset password", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := &AlterUserOptions{
			name:          id,
			ResetPassword: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s RESET PASSWORD", id.FullyQualifiedName())
	})

	t.Run("abort all queries", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := &AlterUserOptions{
			name:            id,
			AbortAllQueries: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s ABORT ALL QUERIES", id.FullyQualifiedName())
	})

	t.Run("rename", func(t *testing.T) {
		newID := randomAccountObjectIdentifier()
		opts := &AlterUserOptions{
			name:    id,
			NewName: newID,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s RENAME TO %s", id.FullyQualifiedName(), newID.FullyQualifiedName())
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
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s ADD DELEGATED AUTHORIZATION OF ROLE %s TO SECURITY INTEGRATION %s", id.FullyQualifiedName(), role, integration)
	})

	t.Run("with unsetting tags", func(t *testing.T) {
		tag1 := randomSchemaObjectIdentifier()
		tag2 := randomSchemaObjectIdentifier()
		opts := &AlterUserOptions{
			name:     id,
			UnsetTag: []ObjectIdentifier{tag1, tag2},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s UNSET TAG %s, %s", id.FullyQualifiedName(), tag1.FullyQualifiedName(), tag2.FullyQualifiedName())
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
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s UNSET PASSWORD", id.FullyQualifiedName())
	})

	t.Run("with unsetting a policy", func(t *testing.T) {
		sessionPolicy := "SESSION_POLICY1"
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				SessionPolicy: String(sessionPolicy),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET SESSION POLICY = %s", id.FullyQualifiedName(), sessionPolicy)
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
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s REMOVE DELEGATED AUTHORIZATION OF ROLE %s FROM SECURITY INTEGRATION %s", id.FullyQualifiedName(), role, integration)
	})
}

func TestUserDrop(t *testing.T) {
	id := randomAccountObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &DropUserOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &DropUserOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "DROP USER %s", id.FullyQualifiedName())
	})
}

func TestUserShow(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &ShowUserOptions{}
		assertOptsValidAndSQLEquals(t, opts, "SHOW USERS")
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowUserOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW USERS LIKE '%s'", id.Name())
	})

	t.Run("with like and from", func(t *testing.T) {
		fromPatern := random.String()
		opts := &ShowUserOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			From: &fromPatern,
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW USERS LIKE '%s' FROM '%s'", id.Name(), fromPatern)
	})

	t.Run("with like and limit", func(t *testing.T) {
		limit := 5
		opts := &ShowUserOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			Limit: &limit,
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW USERS LIKE '%s' LIMIT %v", id.Name(), limit)
	})

	t.Run("with starts with and from", func(t *testing.T) {
		fromPattern := random.String()
		startsWithPattern := random.String()

		opts := &ShowUserOptions{
			StartsWith: &startsWithPattern,
			From:       &fromPattern,
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW USERS STARTS WITH '%s' FROM '%s'", startsWithPattern, fromPattern)
	})
}

func TestUserDescribe(t *testing.T) {
	id := randomAccountObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &describeUserOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &describeUserOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE USER %s", id.FullyQualifiedName())
	})
}
