package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
)

func TestMaskingPolicyCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier()
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

	t.Run("validation: no body", func(t *testing.T) {
		opts := &CreateMaskingPolicyOptions{
			name:      id,
			signature: signature,
			returns:   DataTypeVARCHAR,
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateMaskingPolicyOptions", "body"))
	})

	t.Run("validation: no signature", func(t *testing.T) {
		opts := &CreateMaskingPolicyOptions{
			name:    id,
			body:    expression,
			returns: DataTypeVARCHAR,
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateMaskingPolicyOptions", "signature"))
	})

	t.Run("validation: no returns", func(t *testing.T) {
		opts := &CreateMaskingPolicyOptions{
			name:      id,
			signature: signature,
			body:      expression,
		}
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateMaskingPolicyOptions", "returns"))
	})

	t.Run("only required options", func(t *testing.T) {
		opts := &CreateMaskingPolicyOptions{
			name:      id,
			signature: signature,
			body:      expression,
			returns:   DataTypeVARCHAR,
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE MASKING POLICY %s AS ("col1" VARCHAR, "col2" VARCHAR) RETURNS %s -> %s`, id.FullyQualifiedName(), DataTypeVARCHAR, expression)
	})

	t.Run("with complete options", func(t *testing.T) {
		comment := random.Comment()

		opts := &CreateMaskingPolicyOptions{
			OrReplace:           Bool(true),
			name:                id,
			IfNotExists:         Bool(true),
			signature:           signature,
			body:                expression,
			returns:             DataTypeVARCHAR,
			Comment:             String(comment),
			ExemptOtherPolicies: Bool(true),
		}

		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE MASKING POLICY IF NOT EXISTS %s AS ("col1" VARCHAR, "col2" VARCHAR) RETURNS %s -> %s COMMENT = '%s' EXEMPT_OTHER_POLICIES = %t`, id.FullyQualifiedName(), DataTypeVARCHAR, expression, comment, true)
	})
}

// TODO: add tests for body and tags
func TestMaskingPolicyAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &AlterMaskingPolicyOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no option", func(t *testing.T) {
		opts := &AlterMaskingPolicyOptions{
			name: id,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterMaskingPolicyOptions", "Set", "Unset", "SetTag", "UnsetTag", "NewName"))
	})

	t.Run("with set", func(t *testing.T) {
		newComment := random.Comment()
		opts := &AlterMaskingPolicyOptions{
			name: id,
			Set: &MaskingPolicySet{
				Comment: String(newComment),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER MASKING POLICY %s SET COMMENT = '%s'", id.FullyQualifiedName(), newComment)
	})

	t.Run("with unset", func(t *testing.T) {
		opts := &AlterMaskingPolicyOptions{
			name: id,
			Unset: &MaskingPolicyUnset{
				Comment: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER MASKING POLICY %s UNSET COMMENT", id.FullyQualifiedName())
	})

	t.Run("rename", func(t *testing.T) {
		newID := randomSchemaObjectIdentifierInSchema(id.SchemaId())
		opts := &AlterMaskingPolicyOptions{
			name:    id,
			NewName: &newID,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER MASKING POLICY %s RENAME TO %s", id.FullyQualifiedName(), newID.FullyQualifiedName())
	})
}

func TestMaskingPolicyDrop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &DropMaskingPolicyOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &DropMaskingPolicyOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "DROP MASKING POLICY %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := &DropMaskingPolicyOptions{
			name:     id,
			IfExists: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "DROP MASKING POLICY IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestMaskingPolicyShow(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &ShowMaskingPolicyOptions{}
		assertOptsValidAndSQLEquals(t, opts, "SHOW MASKING POLICIES")
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowMaskingPolicyOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW MASKING POLICIES LIKE '%s'", id.Name())
	})

	t.Run("with like and in account", func(t *testing.T) {
		opts := &ShowMaskingPolicyOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &In{
				Account: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW MASKING POLICIES LIKE '%s' IN ACCOUNT", id.Name())
	})

	t.Run("with like and in database", func(t *testing.T) {
		databaseIdentifier := NewAccountObjectIdentifier(id.DatabaseName())
		opts := &ShowMaskingPolicyOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &In{
				Database: databaseIdentifier,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW MASKING POLICIES LIKE '%s' IN DATABASE %s", id.Name(), databaseIdentifier.FullyQualifiedName())
	})

	t.Run("with like and in schema", func(t *testing.T) {
		schemaIdentifier := NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName())
		opts := &ShowMaskingPolicyOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &In{
				Schema: schemaIdentifier,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW MASKING POLICIES LIKE '%s' IN SCHEMA %s", id.Name(), schemaIdentifier.FullyQualifiedName())
	})

	t.Run("with limit", func(t *testing.T) {
		opts := &ShowMaskingPolicyOptions{
			Limit: Int(10),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW MASKING POLICIES LIMIT 10")
	})
}

func TestMaskingPolicyDescribe(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &describeMaskingPolicyOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &describeMaskingPolicyOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE MASKING POLICY %s", id.FullyQualifiedName())
	})
}
