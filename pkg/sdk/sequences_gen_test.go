package sdk

import (
	"testing"
)

func TestSequences_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateSequenceOptions {
		return &CreateSequenceOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateSequenceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSequenceOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Start = Int(1)
		opts.Increment = Int(1)
		opts.ValuesBehavior = ValuesBehaviorPointer(ValuesBehaviorOrder)
		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE SEQUENCE %s START = 1 INCREMENT = 1 ORDER COMMENT = 'comment'`, id.FullyQualifiedName())
	})
}

func TestSequences_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *AlterSequenceOptions {
		return &AlterSequenceOptions{
			name:     id,
			IfExists: Bool(true),
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*AlterSequenceOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterSequenceOptions", "RenameTo", "SetIncrement", "Set", "UnsetComment"))
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetIncrement = Int(1)
		opts.UnsetComment = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterSequenceOptions", "RenameTo", "SetIncrement", "Set", "UnsetComment"))
	})

	t.Run("alter: rename to", func(t *testing.T) {
		opts := defaultOpts()
		target := randomSchemaObjectIdentifierInSchema(id.SchemaId())
		opts.RenameTo = &target
		assertOptsValidAndSQLEquals(t, opts, `ALTER SEQUENCE IF EXISTS %s RENAME TO %s`, id.FullyQualifiedName(), opts.RenameTo.FullyQualifiedName())
	})

	t.Run("alter: set options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &SequenceSet{
			Comment:        String("comment"),
			ValuesBehavior: ValuesBehaviorPointer(ValuesBehaviorOrder),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SEQUENCE IF EXISTS %s SET ORDER COMMENT = 'comment'`, id.FullyQualifiedName())
	})

	t.Run("alter: unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetComment = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER SEQUENCE IF EXISTS %s UNSET COMMENT`, id.FullyQualifiedName())
	})

	t.Run("alter: set increment", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetIncrement = Int(1)
		assertOptsValidAndSQLEquals(t, opts, `ALTER SEQUENCE IF EXISTS %s SET INCREMENT = 1`, id.FullyQualifiedName())
	})
}

func TestSequences_Show(t *testing.T) {
	defaultOpts := func() *ShowSequenceOptions {
		return &ShowSequenceOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*ShowSequenceOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("show with empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW SEQUENCES`)
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("pattern"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW SEQUENCES LIKE 'pattern'`)
	})

	t.Run("show with in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Account: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW SEQUENCES IN ACCOUNT`)
	})
}

func TestSequences_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *DescribeSequenceOptions {
		return &DescribeSequenceOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DescribeSequenceOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE SEQUENCE %s`, id.FullyQualifiedName())
	})
}

func TestSequences_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *DropSequenceOptions {
		return &DropSequenceOptions{
			name: id,
		}
	}
	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropSequenceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Constraint = &SequenceConstraint{
			Cascade: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP SEQUENCE IF EXISTS %s CASCADE`, id.FullyQualifiedName())
	})
}
