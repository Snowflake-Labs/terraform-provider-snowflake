package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDynamicTableCreate(t *testing.T) {
	t.Run("complete", func(t *testing.T) {
		opts := &CreateDynamicTableOptions{
			OrReplace: Bool(true),
			name: AccountObjectIdentifier{
				name: "dynamic-table",
			},
			TargetLag: "1 minutes",
			warehouse: AccountObjectIdentifier{
				name: "warehouse_name",
			},
			Query:   "SELECT product_id, product_name FROM staging_table",
			Comment: String("comment"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE OR REPLACE DYNAMIC TABLE "dynamic-table" TARGET_LAG = '1 minutes' WAREHOUSE = "warehouse_name" AS SELECT product_id, product_name FROM staging_table COMMENT = 'comment'`
		assert.Equal(t, expected, actual)
	})

	t.Run("validate-opts-target-lag-string", func(t *testing.T) {
		opts := &CreateDynamicTableOptions{
			OrReplace: Bool(true),
			name: AccountObjectIdentifier{
				name: "dynamic-table",
			},
			TargetLag: "1minutes",
			warehouse: AccountObjectIdentifier{
				name: "warehouse_name",
			},
			Query:   "SELECT product_id, product_name FROM staging_table",
			Comment: String("comment"),
		}
		err := opts.validate()
		expected := `The string format is invalid`
		assert.Equal(t, expected, err.Error())
	})

	t.Run("validate-opts-target-lag-number", func(t *testing.T) {
		opts := &CreateDynamicTableOptions{
			OrReplace: Bool(true),
			name: AccountObjectIdentifier{
				name: "dynamic-table",
			},
			TargetLag: "no minutes",
			warehouse: AccountObjectIdentifier{
				name: "warehouse_name",
			},
			Query:   "SELECT product_id, product_name FROM staging_table",
			Comment: String("comment"),
		}
		err := opts.validate()
		expected := `The number value is invalid`
		assert.Equal(t, expected, err.Error())
	})

	t.Run("validate-opts-target-lag-unit", func(t *testing.T) {
		opts := &CreateDynamicTableOptions{
			OrReplace: Bool(true),
			name: AccountObjectIdentifier{
				name: "dynamic-table",
			},
			TargetLag: "1 year",
			warehouse: AccountObjectIdentifier{
				name: "warehouse_name",
			},
			Query:   "SELECT product_id, product_name FROM staging_table",
			Comment: String("comment"),
		}
		err := opts.validate()
		expected := `The unit is invalid`
		assert.Equal(t, expected, err.Error())
	})
}

func TestDynamicTableAlter(t *testing.T) {
	t.Run("suspend", func(t *testing.T) {
		opts := &AlterDynamicTableOptions{
			name: AccountObjectIdentifier{
				name: "dynamic-table",
			},
			Suspend: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER DYNAMIC TABLE "dynamic-table" SUSPEND`
		assert.Equal(t, expected, actual)
	})

	t.Run("resume", func(t *testing.T) {
		opts := &AlterDynamicTableOptions{
			name: AccountObjectIdentifier{
				name: "dynamic-table",
			},
			Resume: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER DYNAMIC TABLE "dynamic-table" RESUME`
		assert.Equal(t, expected, actual)
	})

	t.Run("with set", func(t *testing.T) {
		lag := TargetLag("1 minutes")
		opts := &AlterDynamicTableOptions{
			name: AccountObjectIdentifier{
				name: "dynamic-table",
			},
			Set: &DynamicTableSet{
				TargetLag: &lag,
				Warehouse: &AccountObjectIdentifier{
					name: "warehouse_name",
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER DYNAMIC TABLE "dynamic-table" SET TARGET_LAG = '1 minutes' WAREHOUSE = "warehouse_name"`
		assert.Equal(t, expected, actual)
	})

	t.Run("validate-opts-exact-one-work", func(t *testing.T) {
		opts := &AlterDynamicTableOptions{
			name: AccountObjectIdentifier{
				name: "dynamic-table",
			},
			Suspend: Bool(true),
			Resume:  Bool(true),
		}
		err := opts.validate()
		expected := `exactly one of Suspend, Resume, Refresh, Set must be set`
		assert.Equal(t, expected, err.Error())
	})
}

func TestDynamicTableDrop(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &DropDynamicTableOptions{
			name: NewAccountObjectIdentifier("dynamic-table"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP DYNAMIC TABLE "dynamic-table"`
		assert.Equal(t, expected, actual)
	})
}

func TestDynamicTableShow(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		opts := &ShowDynamicTableOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW DYNAMIC TABLES`
		assert.Equal(t, expected, actual)
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowDynamicTableOptions{
			Like: &Like{
				Pattern: String("dynamic-table"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW DYNAMIC TABLES LIKE 'dynamic-table'`
		assert.Equal(t, expected, actual)
	})
}

func TestDynamicTableDescribe(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &dynamicTableDescribeOptions{
			name: NewAccountObjectIdentifier("dynamic-table"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DESCRIBE DYNAMIC TABLE "dynamic-table"`
		assert.Equal(t, expected, actual)
	})
}
