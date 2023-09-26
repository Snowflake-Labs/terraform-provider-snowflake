package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDynamicTableCreate(t *testing.T) {
	t.Run("complete", func(t *testing.T) {
		opts := &createDynamicTableOptions{
			OrReplace: Bool(true),
			name: AccountObjectIdentifier{
				name: "dynamic-table",
			},
			targetLag: TargetLag{
				Lagtime: String("1 minutes"),
			},
			warehouse: AccountObjectIdentifier{
				name: "warehouse_name",
			},
			query:   "SELECT product_id, product_name FROM staging_table",
			Comment: String("comment"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE OR REPLACE DYNAMIC TABLE "dynamic-table" TARGET_LAG = '1 minutes' WAREHOUSE = "warehouse_name" COMMENT = 'comment' AS SELECT product_id, product_name FROM staging_table`
		assert.Equal(t, expected, actual)
	})

	t.Run("complete with target lag", func(t *testing.T) {
		opts := &createDynamicTableOptions{
			OrReplace: Bool(true),
			name: AccountObjectIdentifier{
				name: "dynamic-table",
			},
			targetLag: TargetLag{
				Downstream: Bool(true),
			},
			warehouse: AccountObjectIdentifier{
				name: "warehouse_name",
			},
			query:   "SELECT product_id, product_name FROM staging_table",
			Comment: String("comment"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE OR REPLACE DYNAMIC TABLE "dynamic-table" TARGET_LAG = DOWNSTREAM WAREHOUSE = "warehouse_name" COMMENT = 'comment' AS SELECT product_id, product_name FROM staging_table`
		assert.Equal(t, expected, actual)
	})
}

func TestDynamicTableAlter(t *testing.T) {
	t.Run("suspend", func(t *testing.T) {
		opts := &alterDynamicTableOptions{
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
		opts := &alterDynamicTableOptions{
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
		lag := TargetLag{
			Lagtime: String("1 minutes"),
		}
		opts := &alterDynamicTableOptions{
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
		opts := &alterDynamicTableOptions{
			name: AccountObjectIdentifier{
				name: "dynamic-table",
			},
			Suspend: Bool(true),
			Resume:  Bool(true),
		}
		err := opts.validate()
		require.Error(t, err)
	})
}

func TestDynamicTableDrop(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &dropDynamicTableOptions{
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
		opts := &describeDynamicTableOptions{
			name: NewAccountObjectIdentifier("dynamic-table"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DESCRIBE DYNAMIC TABLE "dynamic-table"`
		assert.Equal(t, expected, actual)
	})
}
