package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestPasswordPolicyCreate(t *testing.T) {
	builder := testBuilder(t)
	id := randomSchemaObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &PasswordPolicyCreateOptions{}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := "CREATE PASSWORD POLICY"
		assert.Equal(t, expected, actual)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &PasswordPolicyCreateOptions{
			name: id,
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("CREATE PASSWORD POLICY %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with complete options", func(t *testing.T) {
		opts := &PasswordPolicyCreateOptions{
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
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf(`CREATE OR REPLACE PASSWORD POLICY IF NOT EXISTS %s PASSWORD_MIN_LENGTH = 10 PASSWORD_MAX_LENGTH = 20 PASSWORD_MIN_UPPER_CASE_CHARS = 1 PASSWORD_MIN_LOWER_CASE_CHARS = 1 PASSWORD_MIN_NUMERIC_CHARS = 1 PASSWORD_MIN_SPECIAL_CHARS = 1 PASSWORD_MAX_AGE_DAYS = 30 PASSWORD_MAX_RETRIES = 5 PASSWORD_LOCKOUT_TIME_MINS = 30 COMMENT = 'test comment'`, id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}

func TestPasswordPolicyAlter(t *testing.T) {
	builder := testBuilder(t)
	id := randomSchemaObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &PasswordPolicyAlterOptions{}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := "ALTER PASSWORD POLICY"
		assert.Equal(t, expected, actual)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &PasswordPolicyAlterOptions{
			name: id,
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("ALTER PASSWORD POLICY %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with set", func(t *testing.T) {
		opts := &PasswordPolicyAlterOptions{
			name: id,
			Set: &PasswordPolicySet{
				PasswordMinLength:         Int(10),
				PasswordMaxLength:         Int(20),
				PasswordMinUpperCaseChars: Int(1),
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("ALTER PASSWORD POLICY %s SET PASSWORD_MIN_LENGTH = 10 PASSWORD_MAX_LENGTH = 20 PASSWORD_MIN_UPPER_CASE_CHARS = 1", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with unset", func(t *testing.T) {
		opts := &PasswordPolicyAlterOptions{
			name: id,
			Unset: &PasswordPolicyUnset{
				PasswordMinLength: Bool(true),
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("ALTER PASSWORD POLICY %s UNSET PASSWORD_MIN_LENGTH", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("rename", func(t *testing.T) {
		newID := NewSchemaObjectIdentifier(id.databaseName, id.schemaName, randomUUID(t))
		opts := &PasswordPolicyAlterOptions{
			name:    id,
			NewName: newID,
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("ALTER PASSWORD POLICY %s RENAME TO %s", id.FullyQualifiedName(), newID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}

func TestPasswordPolicyDrop(t *testing.T) {
	builder := testBuilder(t)
	id := randomSchemaObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &PasswordPolicyDropOptions{}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := "DROP PASSWORD POLICY"
		assert.Equal(t, expected, actual)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &PasswordPolicyDropOptions{
			name: id,
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("DROP PASSWORD POLICY %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with if exists", func(t *testing.T) {
		opts := &PasswordPolicyDropOptions{
			name:     id,
			IfExists: Bool(true),
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("DROP PASSWORD POLICY IF EXISTS %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}

func TestPasswordPolicyShow(t *testing.T) {
	builder := testBuilder(t)
	id := randomSchemaObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &PasswordPolicyShowOptions{}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := "SHOW PASSWORD POLICIES"
		assert.Equal(t, expected, actual)
	})

	t.Run("with like", func(t *testing.T) {
		opts := &PasswordPolicyShowOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("SHOW PASSWORD POLICIES LIKE '%s'", id.Name())
		assert.Equal(t, expected, actual)
	})

	t.Run("with like and in account", func(t *testing.T) {
		opts := &PasswordPolicyShowOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &In{
				Account: Bool(true),
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("SHOW PASSWORD POLICIES LIKE '%s' IN ACCOUNT", id.Name())
		assert.Equal(t, expected, actual)
	})

	t.Run("with like and in database", func(t *testing.T) {
		databaseIdentifier := NewAccountObjectIdentifier(id.DatabaseName())
		opts := &PasswordPolicyShowOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &In{
				Database: databaseIdentifier,
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("SHOW PASSWORD POLICIES LIKE '%s' IN DATABASE %s", id.Name(), databaseIdentifier.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with like and in schema", func(t *testing.T) {
		schemaIdentifier := NewSchemaIdentifier(id.DatabaseName(), id.SchemaName())
		opts := &PasswordPolicyShowOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &In{
				Schema: schemaIdentifier,
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("SHOW PASSWORD POLICIES LIKE '%s' IN SCHEMA %s", id.Name(), schemaIdentifier.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with limit", func(t *testing.T) {
		opts := &PasswordPolicyShowOptions{
			Limit: Int(10),
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := "SHOW PASSWORD POLICIES LIMIT 10"
		assert.Equal(t, expected, actual)
	})
}

func TestPasswordPolicyDescribe(t *testing.T) {
	builder := testBuilder(t)
	id := randomSchemaObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &passwordPolicyDescribeOptions{}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := "DESCRIBE PASSWORD POLICY"
		assert.Equal(t, expected, actual)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &passwordPolicyDescribeOptions{
			name: id,
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("DESCRIBE PASSWORD POLICY %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}
