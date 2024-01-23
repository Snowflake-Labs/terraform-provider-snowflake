package sdk

import "testing"

func TestStreamlits_Create(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *CreateStreamlitOptions {
		return &CreateStreamlitOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateStreamlitOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateStreamlitOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("all options", func(t *testing.T) {
		warehouse := NewAccountObjectIdentifier("test_warehouse")
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.RootLocation = "@test"
		opts.MainFile = "manifest.yml"
		opts.Warehouse = &warehouse
		opts.Comment = String("test")
		assertOptsValidAndSQLEquals(t, opts, `CREATE STREAMLIT IF NOT EXISTS %s ROOT_LOCATION = '@test' MAIN_FILE = 'manifest.yml' QUERY_WAREHOUSE = %s COMMENT = 'test'`, id.FullyQualifiedName(), warehouse.FullyQualifiedName())
	})
}

func TestStreamlits_Alter(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *AlterStreamlitOptions {
		return &AlterStreamlitOptions{
			IfExists: Bool(true),
			name:     id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterStreamlitOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterStreamlitOptions", "RenameTo", "Set"))
	})

	t.Run("alter: set options", func(t *testing.T) {
		warehouse := NewAccountObjectIdentifier("test_warehouse")

		opts := defaultOpts()
		opts.Set = &StreamlitSet{
			RootLocation: String("@test"),
			MainFile:     String("manifest.yml"),
			Warehouse:    &warehouse,
			Comment:      String("test"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER STREAMLIT IF EXISTS %s SET ROOT_LOCATION = '@test' MAIN_FILE = 'manifest.yml' QUERY_WAREHOUSE = %s COMMENT = 'test'`, id.FullyQualifiedName(), warehouse.FullyQualifiedName())
	})
}

func TestStreamlits_Drop(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *DropStreamlitOptions {
		return &DropStreamlitOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropStreamlitOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP STREAMLIT %s`, id.FullyQualifiedName())
	})
}

func TestStreamlits_Show(t *testing.T) {
	defaultOpts := func() *ShowStreamlitOptions {
		return &ShowStreamlitOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowStreamlitOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("show with empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW STREAMLITS`)
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("pattern"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW STREAMLITS LIKE 'pattern'`)
	})

	t.Run("show with in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Account: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW STREAMLITS IN ACCOUNT`)
	})
}

func TestStreamlits_Describe(t *testing.T) {
	id := RandomSchemaObjectIdentifier()

	defaultOpts := func() *DescribeStreamlitOptions {
		return &DescribeStreamlitOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DescribeStreamlitOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE STREAMLIT %s`, id.FullyQualifiedName())
	})
}
