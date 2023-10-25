// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/sdk/internal/random"
)

func TestPasswordPolicyCreate(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &CreatePasswordPolicyOptions{}
		assertOptsInvalid(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &CreatePasswordPolicyOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE PASSWORD POLICY %s", id.FullyQualifiedName())
	})

	t.Run("with complete options", func(t *testing.T) {
		opts := &CreatePasswordPolicyOptions{
			OrReplace:                 Bool(true),
			name:                      id,
			IfNotExists:               Bool(true),
			PasswordMinLength:         Int(10),
			PasswordMaxLength:         Int(20),
			PasswordMinUpperCaseChars: Int(1),
			PasswordMinLowerCaseChars: Int(1),
			PasswordMinNumericChars:   Int(1),
			PasswordMinSpecialChars:   Int(1),
			PasswordMaxAgeDays:        Int(30),
			PasswordMaxRetries:        Int(5),
			PasswordLockoutTimeMins:   Int(30),
			Comment:                   String("test comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE PASSWORD POLICY IF NOT EXISTS %s PASSWORD_MIN_LENGTH = 10 PASSWORD_MAX_LENGTH = 20 PASSWORD_MIN_UPPER_CASE_CHARS = 1 PASSWORD_MIN_LOWER_CASE_CHARS = 1 PASSWORD_MIN_NUMERIC_CHARS = 1 PASSWORD_MIN_SPECIAL_CHARS = 1 PASSWORD_MAX_AGE_DAYS = 30 PASSWORD_MAX_RETRIES = 5 PASSWORD_LOCKOUT_TIME_MINS = 30 COMMENT = 'test comment'`, id.FullyQualifiedName())
	})
}

func TestPasswordPolicyAlter(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &AlterPasswordPolicyOptions{}
		assertOptsInvalid(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &AlterPasswordPolicyOptions{
			name: id,
		}
		assertOptsInvalid(t, opts, errExactlyOneOf("Set", "Unset", "NewName"))
	})

	t.Run("with set", func(t *testing.T) {
		opts := &AlterPasswordPolicyOptions{
			name: id,
			Set: &PasswordPolicySet{
				PasswordMinLength:         Int(10),
				PasswordMaxLength:         Int(20),
				PasswordMinUpperCaseChars: Int(1),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER PASSWORD POLICY %s SET PASSWORD_MIN_LENGTH = 10 PASSWORD_MAX_LENGTH = 20 PASSWORD_MIN_UPPER_CASE_CHARS = 1", id.FullyQualifiedName())
	})

	t.Run("with unset", func(t *testing.T) {
		opts := &AlterPasswordPolicyOptions{
			name: id,
			Unset: &PasswordPolicyUnset{
				PasswordMinLength: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER PASSWORD POLICY %s UNSET PASSWORD_MIN_LENGTH", id.FullyQualifiedName())
	})

	t.Run("rename", func(t *testing.T) {
		newID := NewSchemaObjectIdentifier(id.databaseName, id.schemaName, random.UUID())
		opts := &AlterPasswordPolicyOptions{
			name:    id,
			NewName: &newID,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER PASSWORD POLICY %s RENAME TO %s", id.FullyQualifiedName(), newID.FullyQualifiedName())
	})
}

func TestPasswordPolicyDrop(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &DropPasswordPolicyOptions{}
		assertOptsInvalid(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &DropPasswordPolicyOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "DROP PASSWORD POLICY %s", id.FullyQualifiedName())
	})

	t.Run("with if exists", func(t *testing.T) {
		opts := &DropPasswordPolicyOptions{
			name:     id,
			IfExists: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "DROP PASSWORD POLICY IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestPasswordPolicyShow(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &ShowPasswordPolicyOptions{}
		assertOptsValidAndSQLEquals(t, opts, "SHOW PASSWORD POLICIES")
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowPasswordPolicyOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW PASSWORD POLICIES LIKE '%s'", id.Name())
	})

	t.Run("with like and in account", func(t *testing.T) {
		opts := &ShowPasswordPolicyOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &In{
				Account: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW PASSWORD POLICIES LIKE '%s' IN ACCOUNT", id.Name())
	})

	t.Run("with like and in database", func(t *testing.T) {
		databaseIdentifier := NewAccountObjectIdentifier(id.DatabaseName())
		opts := &ShowPasswordPolicyOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &In{
				Database: databaseIdentifier,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW PASSWORD POLICIES LIKE '%s' IN DATABASE %s", id.Name(), databaseIdentifier.FullyQualifiedName())
	})

	t.Run("with like and in schema", func(t *testing.T) {
		schemaIdentifier := NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName())
		opts := &ShowPasswordPolicyOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &In{
				Schema: schemaIdentifier,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW PASSWORD POLICIES LIKE '%s' IN SCHEMA %s", id.Name(), schemaIdentifier.FullyQualifiedName())
	})

	t.Run("with limit", func(t *testing.T) {
		opts := &ShowPasswordPolicyOptions{
			Limit: Int(10),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW PASSWORD POLICIES LIMIT 10")
	})
}

func TestPasswordPolicyDescribe(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &describePasswordPolicyOptions{}
		assertOptsInvalid(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &describePasswordPolicyOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE PASSWORD POLICY %s", id.FullyQualifiedName())
	})
}
