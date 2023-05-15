package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestMaskingPolicyCreate(t *testing.T) {
	builder := testBuilder(t)
	id := randomSchemaObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &MaskingPolicyCreateOptions{}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := "CREATE MASKING POLICY RETURNS ->"
		assert.Equal(t, expected, actual)
	})

	t.Run("with complete options", func(t *testing.T) {
		signature := []TableColumnSignature{
			{
				Name: "col1",
				Type: DataTypeVARCHAR,
			},
			{
				Name: "col2",
				Type: DataTypeVARCHAR,
			},
		}
		expression := "REPLACE('X', 1, 2)"
		comment := randomString(t)

		opts := &MaskingPolicyCreateOptions{
			OrReplace:           Bool(true),
			name:                id,
			IfNotExists:         Bool(true),
			signature:           signature,
			body:                expression,
			returns:             DataTypeVARCHAR,
			Comment:             String(comment),
			ExemptOtherPolicies: Bool(true),
		}

		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf(`CREATE OR REPLACE MASKING POLICY IF NOT EXISTS %s AS ("col1" VARCHAR,"col2" VARCHAR) RETURNS %s -> %s COMMENT = '%s' EXEMPT_OTHER_POLICIES = %t`, id.FullyQualifiedName(), DataTypeVARCHAR, expression, comment, true)
		assert.Equal(t, expected, actual)
	})
}

func TestMaskingPolicyAlter(t *testing.T) {
	builder := testBuilder(t)
	id := randomSchemaObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &MaskingPolicyAlterOptions{}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := "ALTER MASKING POLICY"
		assert.Equal(t, expected, actual)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &MaskingPolicyAlterOptions{
			name: id,
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("ALTER MASKING POLICY %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with set", func(t *testing.T) {
		newComment := randomString(t)
		opts := &MaskingPolicyAlterOptions{
			name: id,
			Set: &MaskingPolicySet{
				Comment: String(newComment),
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("ALTER MASKING POLICY %s SET COMMENT = '%s'", id.FullyQualifiedName(), newComment)
		assert.Equal(t, expected, actual)
	})

	t.Run("with unset", func(t *testing.T) {
		opts := &MaskingPolicyAlterOptions{
			name: id,
			Unset: &MaskingPolicyUnset{
				Comment: Bool(true),
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("ALTER MASKING POLICY %s UNSET COMMENT", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("rename", func(t *testing.T) {
		newID := NewSchemaObjectIdentifier(id.databaseName, id.schemaName, randomUUID(t))
		opts := &MaskingPolicyAlterOptions{
			name:    id,
			NewName: newID,
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("ALTER MASKING POLICY %s RENAME TO %s", id.FullyQualifiedName(), newID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}

func TestMaskingPolicyDrop(t *testing.T) {
	builder := testBuilder(t)
	id := randomSchemaObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &MaskingPolicyDropOptions{}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := "DROP MASKING POLICY"
		assert.Equal(t, expected, actual)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &MaskingPolicyDropOptions{
			name: id,
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("DROP MASKING POLICY %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}

func TestMaskingPolicyShow(t *testing.T) {
	builder := testBuilder(t)
	id := randomSchemaObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &MaskingPolicyShowOptions{}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := "SHOW MASKING POLICIES"
		assert.Equal(t, expected, actual)
	})

	t.Run("with like", func(t *testing.T) {
		opts := &MaskingPolicyShowOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("SHOW MASKING POLICIES LIKE '%s'", id.Name())
		assert.Equal(t, expected, actual)
	})

	t.Run("with like and in account", func(t *testing.T) {
		opts := &MaskingPolicyShowOptions{
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
		expected := fmt.Sprintf("SHOW MASKING POLICIES LIKE '%s' IN ACCOUNT", id.Name())
		assert.Equal(t, expected, actual)
	})

	t.Run("with like and in database", func(t *testing.T) {
		databaseIdentifier := NewAccountObjectIdentifier(id.DatabaseName())
		opts := &MaskingPolicyShowOptions{
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
		expected := fmt.Sprintf("SHOW MASKING POLICIES LIKE '%s' IN DATABASE %s", id.Name(), databaseIdentifier.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with like and in schema", func(t *testing.T) {
		schemaIdentifier := NewSchemaIdentifier(id.DatabaseName(), id.SchemaName())
		opts := &MaskingPolicyShowOptions{
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
		expected := fmt.Sprintf("SHOW MASKING POLICIES LIKE '%s' IN SCHEMA %s", id.Name(), schemaIdentifier.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with limit", func(t *testing.T) {
		opts := &MaskingPolicyShowOptions{
			Limit: Int(10),
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := "SHOW MASKING POLICIES LIMIT 10"
		assert.Equal(t, expected, actual)
	})
}

func TestMaskingPolicyDescribe(t *testing.T) {
	builder := testBuilder(t)
	id := randomSchemaObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &maskingPolicyDescribeOptions{}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := "DESCRIBE MASKING POLICY"
		assert.Equal(t, expected, actual)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &maskingPolicyDescribeOptions{
			name: id,
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		actual := builder.sql(clauses...)
		expected := fmt.Sprintf("DESCRIBE MASKING POLICY %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}
