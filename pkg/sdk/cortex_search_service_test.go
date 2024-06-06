package sdk

import (
	"testing"
)

func TestCortexSearchServiceCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *createCortexSearchServiceOptions {
		return &createCortexSearchServiceOptions{
			name:      id,
			on:        "searchable_text",
			targetLag: "1 minutes",
			warehouse: AccountObjectIdentifier{
				name: "warehouse_name",
			},
			query: "SELECT product_id, product_name, searchable_text FROM staging_table",
		}
	}
	t.Run("validation: nil options", func(t *testing.T) {
		var opts *createCortexSearchServiceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE CORTEX SEARCH SERVICE %s ON searchable_text WAREHOUSE = "warehouse_name" TARGET_LAG = '1 minutes' AS SELECT product_id, product_name, searchable_text FROM staging_table`, id.FullyQualifiedName())
	})

	t.Run("all optional", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.attributes = &Attributes{
			columns: []string{"product_id", "product_name"},
		}
		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE CORTEX SEARCH SERVICE %s ON searchable_text ATTRIBUTES product_id, product_name WAREHOUSE = "warehouse_name" TARGET_LAG = '1 minutes' COMMENT = 'comment' AS SELECT product_id, product_name, searchable_text FROM staging_table`, id.FullyQualifiedName())
	})
}

func TestCortexSearchServiceAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *alterCortexSearchServiceOptions {
		return &alterCortexSearchServiceOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *alterCortexSearchServiceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no alter action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("alterCortexSearchServiceOptions", "Set"))
	})

	t.Run("validation: no property to unset", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("alterCortexSearchServiceOptions", "Set"))
	})
	t.Run("set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &CortexSearchServiceSet{
			TargetLag: String("1 minutes"),
			Warehouse: &AccountObjectIdentifier{
				name: "warehouse_name",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER CORTEX SEARCH SERVICE %s SET TARGET_LAG = '1 minutes' WAREHOUSE = "warehouse_name"`, id.FullyQualifiedName())
	})
}

func TestCortexSearchServiceDrop(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *dropCortexSearchServiceOptions {
		return &dropCortexSearchServiceOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *dropCortexSearchServiceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP CORTEX SEARCH SERVICE %s`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP CORTEX SEARCH SERVICE IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestCortexSearchServiceShow(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *showCortexSearchServiceOptions {
		return &showCortexSearchServiceOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *showCortexSearchServiceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: empty like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{}
		assertOptsInvalidJoinedErrors(t, opts, ErrPatternRequiredForLikeKeyword)
	})

	t.Run("show with in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Database: NewAccountObjectIdentifier("database"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW CORTEX SEARCH SERVICES IN DATABASE "database"`)
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW CORTEX SEARCH SERVICES LIKE '%s'`, id.Name())
	})

	t.Run("show with like and in", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		opts.In = &In{
			Database: NewAccountObjectIdentifier("database"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW CORTEX SEARCH SERVICES LIKE '%s' IN DATABASE "database"`, id.Name())
	})
}

func TestCortexSearchServiceDescribe(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *describeCortexSearchServiceOptions {
		return &describeCortexSearchServiceOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *describeCortexSearchServiceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("describe", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE CORTEX SEARCH SERVICE %s`, id.FullyQualifiedName())
	})
}
