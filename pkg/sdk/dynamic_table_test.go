package sdk

import (
	"fmt"
	"testing"
)

func TestDynamicTableCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	defaultOpts := func() *createDynamicTableOptions {
		return &createDynamicTableOptions{
			name: id,
			targetLag: TargetLag{
				Lagtime: String("1 minutes"),
			},
			warehouse: AccountObjectIdentifier{
				name: "warehouse_name",
			},
			query: "SELECT product_id, product_name FROM staging_table",
		}
	}
	t.Run("validation: nil options", func(t *testing.T) {
		var opts *createDynamicTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, fmt.Sprintf(`CREATE OR REPLACE DYNAMIC TABLE %s TARGET_LAG = '1 minutes' WAREHOUSE = "warehouse_name" AS SELECT product_id, product_name FROM staging_table`, id.FullyQualifiedName()))
	})

	t.Run("all optional", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, fmt.Sprintf(`CREATE OR REPLACE DYNAMIC TABLE %s TARGET_LAG = '1 minutes' WAREHOUSE = "warehouse_name" COMMENT = 'comment' AS SELECT product_id, product_name FROM staging_table`, id.FullyQualifiedName()))
	})
}

func TestDynamicTableAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	defaultOpts := func() *alterDynamicTableOptions {
		return &alterDynamicTableOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *alterDynamicTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("validation: no alter action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errAlterNeedsExactlyOneAction)
	})

	t.Run("validation: multiple alter actions", func(t *testing.T) {
		opts := defaultOpts()
		opts.Resume = Bool(true)
		opts.Suspend = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errAlterNeedsExactlyOneAction)
	})

	t.Run("suspend", func(t *testing.T) {
		opts := defaultOpts()
		opts.Suspend = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, fmt.Sprintf(`ALTER DYNAMIC TABLE %s SUSPEND`, id.FullyQualifiedName()))
	})

	t.Run("resume", func(t *testing.T) {
		opts := defaultOpts()
		opts.Resume = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, fmt.Sprintf(`ALTER DYNAMIC TABLE %s RESUME`, id.FullyQualifiedName()))
	})

	t.Run("set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DynamicTableSet{
			TargetLag: &TargetLag{
				Lagtime: String("1 minutes"),
			},
			Warehouse: &AccountObjectIdentifier{
				name: "warehouse_name",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, fmt.Sprintf(`ALTER DYNAMIC TABLE %s SET TARGET_LAG = '1 minutes' WAREHOUSE = "warehouse_name"`, id.FullyQualifiedName()))
	})
}

func TestDynamicTableDrop(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	defaultOpts := func() *dropDynamicTableOptions {
		return &dropDynamicTableOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *dropDynamicTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP DYNAMIC TABLE %s`, id.FullyQualifiedName())
	})
}

func TestDynamicTableShow(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	defaultOpts := func() *showDynamicTableOptions {
		return &showDynamicTableOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			In: &In{
				Database: NewAccountObjectIdentifier("my_database"),
			},
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *showDynamicTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: empty like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{}
		assertOptsInvalidJoinedErrors(t, opts, errPatternRequiredForLikeKeyword)
	})

	t.Run("show", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = nil
		assertOptsValidAndSQLEquals(t, opts, `SHOW DYNAMIC TABLES IN DATABASE "my_database"`)
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW DYNAMIC TABLES LIKE '%s' IN DATABASE "my_database"`, id.Name())
	})
}

func TestDynamicTableDescribe(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	defaultOpts := func() *describeDynamicTableOptions {
		return &describeDynamicTableOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *describeDynamicTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, errInvalidObjectIdentifier)
	})

	t.Run("describe", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE DYNAMIC TABLE %s`, id.FullyQualifiedName())
	})
}
