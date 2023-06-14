package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestWarehouseCreate(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &CreateWarehouseOptions{
			name: AccountObjectIdentifier{
				name: "mywarehouse",
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE WAREHOUSE "mywarehouse"`
		assert.Equal(t, expected, actual)
	})

	t.Run("with complete options", func(t *testing.T) {
		opts := &CreateWarehouseOptions{
			OrReplace:   Bool(true),
			name:        NewAccountObjectIdentifier("completewarehouse"),
			IfNotExists: Bool(true),

			WarehouseType:                   &WarehouseTypeStandard,
			WarehouseSize:                   &WarehouseSizeX4Large,
			MaxClusterCount:                 Int(8),
			MinClusterCount:                 Int(3),
			ScalingPolicy:                   &ScalingPolicyEconomy,
			AutoSuspend:                     Int(1000),
			AutoResume:                      Bool(true),
			InitiallySuspended:              Bool(false),
			ResourceMonitor:                 String("myresmon"),
			Comment:                         String("hello"),
			EnableQueryAcceleration:         Bool(true),
			QueryAccelerationMaxScaleFactor: Int(62),

			MaxConcurrencyLevel:             Int(7),
			StatementQueuedTimeoutInSeconds: Int(29),
			StatementTimeoutInSeconds:       Int(89),
			Tag: []TagAssociation{
				{
					Name:  NewSchemaObjectIdentifier("db1", "schema1", "tag1"),
					Value: "v1",
				},
				{
					Name:  NewSchemaObjectIdentifier("db1", "schema1", "tag2"),
					Value: "v2",
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE OR REPLACE WAREHOUSE IF NOT EXISTS "completewarehouse" WAREHOUSE_TYPE = 'STANDARD' WAREHOUSE_SIZE = 'X4LARGE' MAX_CLUSTER_COUNT = 8 MIN_CLUSTER_COUNT = 3 SCALING_POLICY = 'ECONOMY' AUTO_SUSPEND = 1000 AUTO_RESUME = true INITIALLY_SUSPENDED = false RESOURCE_MONITOR = "myresmon" COMMENT = 'hello' ENABLE_QUERY_ACCELERATION = true QUERY_ACCELERATION_MAX_SCALE_FACTOR = 62 MAX_CONCURRENCY_LEVEL = 7 STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 29 STATEMENT_TIMEOUT_IN_SECONDS = 89 TAG ("db1"."schema1"."tag1" = 'v1', "db1"."schema1"."tag2" = 'v2')`
		assert.Equal(t, expected, actual)
	})
}

func TestWarehouseAlter(t *testing.T) {
	t.Run("with set params", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			Set: &WarehouseSet{
				WarehouseType:                   &WarehouseTypeSnowparkOptimized,
				WaitForCompletion:               Bool(false),
				MinClusterCount:                 Int(4),
				AutoSuspend:                     Int(200),
				ResourceMonitor:                 NewAccountObjectIdentifier("resmon"),
				EnableQueryAcceleration:         Bool(false),
				StatementQueuedTimeoutInSeconds: Int(1200),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER WAREHOUSE "mywarehouse" SET WAREHOUSE_TYPE = 'SNOWPARK-OPTIMIZED' WAIT_FOR_COMPLETION = false MIN_CLUSTER_COUNT = 4 AUTO_SUSPEND = 200 RESOURCE_MONITOR = "resmon" ENABLE_QUERY_ACCELERATION = false STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 1200`
		assert.Equal(t, expected, actual)
	})

	t.Run("with set tag", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			Set: &WarehouseSet{
				Tag: []TagAssociation{
					{
						Name:  NewSchemaObjectIdentifier("db", "schema", "tag1"),
						Value: "v1",
					},
					{
						Name:  NewSchemaObjectIdentifier("db", "schema", "tag2"),
						Value: "v2",
					},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER WAREHOUSE SET TAG "db"."schema"."tag1" = 'v1', "db"."schema"."tag2" = 'v2'`
		assert.Equal(t, expected, actual)
	})

	t.Run("with unset tag", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			Unset: &WarehouseUnset{
				Tag: []ObjectIdentifier{
					NewSchemaObjectIdentifier("db", "schema", "tag1"),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER WAREHOUSE "mywarehouse" UNSET TAG "db"."schema"."tag1"`
		assert.Equal(t, expected, actual)
	})

	t.Run("with unset params", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			Unset: &WarehouseUnset{
				WarehouseSize:   Bool(true),
				MaxClusterCount: Bool(true),
				AutoResume:      Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER WAREHOUSE "mywarehouse" UNSET WAREHOUSE_SIZE, MAX_CLUSTER_COUNT, AUTO_RESUME`
		assert.Equal(t, expected, actual)
	})

	t.Run("rename", func(t *testing.T) {
		newname := NewAccountObjectIdentifier("newname")
		opts := &AlterWarehouseOptions{
			name:    NewAccountObjectIdentifier("oldname"),
			NewName: newname,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER WAREHOUSE "oldname" RENAME TO "newname"`
		assert.Equal(t, expected, actual)
	})

	t.Run("suspend", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name:    NewAccountObjectIdentifier("mywarehouse"),
			Suspend: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER WAREHOUSE "mywarehouse" SUSPEND`
		assert.Equal(t, expected, actual)
	})

	t.Run("resume", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name:        NewAccountObjectIdentifier("mywarehouse"),
			Resume:      Bool(true),
			IfSuspended: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER WAREHOUSE "mywarehouse" RESUME IF SUSPENDED`
		assert.Equal(t, expected, actual)
	})

	t.Run("abort all queries", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name:            NewAccountObjectIdentifier("mywarehouse"),
			AbortAllQueries: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER WAREHOUSE "mywarehouse" ABORT ALL QUERIES`
		assert.Equal(t, expected, actual)
	})

	t.Run("with set tag", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			Set: &WarehouseSet{
				Tag: []TagAssociation{
					{
						Name:  NewSchemaObjectIdentifier("db1", "schema1", "tag1"),
						Value: "v1",
					},
					{
						Name:  NewSchemaObjectIdentifier("db2", "schema2", "tag2"),
						Value: "v2",
					},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER WAREHOUSE "mywarehouse" SET TAG "db1"."schema1"."tag1" = 'v1', "db2"."schema2"."tag2" = 'v2'`
		assert.Equal(t, expected, actual)
	})

	t.Run("with unset tag", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			Unset: &WarehouseUnset{
				Tag: []ObjectIdentifier{
					NewSchemaObjectIdentifier("db1", "schema1", "tag1"),
					NewSchemaObjectIdentifier("db2", "schema2", "tag2"),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER WAREHOUSE "mywarehouse" UNSET TAG "db1"."schema1"."tag1", "db2"."schema2"."tag2"`
		assert.Equal(t, expected, actual)
	})
}

func TestWarehouseDrop(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &DropWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP WAREHOUSE "mywarehouse"`
		assert.Equal(t, expected, actual)
	})

	t.Run("with if exists", func(t *testing.T) {
		opts := &DropWarehouseOptions{
			name:     NewAccountObjectIdentifier("mywarehouse"),
			IfExists: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP WAREHOUSE IF EXISTS "mywarehouse"`
		assert.Equal(t, expected, actual)
	})
}

func TestWarehouseShow(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		opts := &ShowWarehouseOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW WAREHOUSES`
		assert.Equal(t, expected, actual)
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowWarehouseOptions{
			Like: &Like{
				Pattern: String("mywarehouse"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW WAREHOUSES LIKE 'mywarehouse'`
		assert.Equal(t, expected, actual)
	})
}

func TestWarehouseDescribe(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &warehouseDescribeOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DESCRIBE WAREHOUSE "mywarehouse"`
		assert.Equal(t, expected, actual)
	})
}

func TestToWarehouseSize(t *testing.T) {
	type test struct {
		input string
		want  WarehouseSize
	}

	tests := []test{
		// case insensitive.
		{input: "XSMALL", want: WarehouseSizeXSmall},
		{input: "xsmall", want: WarehouseSizeXSmall},

		// Supported Values
		{input: "XSMALL", want: WarehouseSizeXSmall},
		{input: "SMALL", want: WarehouseSizeSmall},
		{input: "MEDIUM", want: WarehouseSizeMedium},
		{input: "LARGE", want: WarehouseSizeLarge},
		{input: "XLARGE", want: WarehouseSizeXLarge},
		{input: "XXLARGE", want: WarehouseSizeXXLarge},
		{input: "XXXLARGE", want: WarehouseSizeXXXLarge},
		{input: "X4LARGE", want: WarehouseSizeX4Large},
		{input: "X5LARGE", want: WarehouseSizeX5Large},
		{input: "X6LARGE", want: WarehouseSizeX6Large},

		// Synonyms
		{input: "X-SMALL", want: WarehouseSizeXSmall},
		{input: "X-LARGE", want: WarehouseSizeXLarge},
		{input: "X2LARGE", want: WarehouseSizeXXLarge},
		{input: "2X-LARGE", want: WarehouseSizeXXLarge},
		{input: "X3LARGE", want: WarehouseSizeXXXLarge},
		{input: "3X-LARGE", want: WarehouseSizeXXXLarge},
		{input: "4X-LARGE", want: WarehouseSizeX4Large},
		{input: "5X-LARGE", want: WarehouseSizeX5Large},
		{input: "6X-LARGE", want: WarehouseSizeX6Large},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToWarehouseSize(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})

		t.Run("invalid warehouse size", func(t *testing.T) {
			_, err := ToWarehouseSize("foo")
			require.Error(t, err)
		})
	}
}
