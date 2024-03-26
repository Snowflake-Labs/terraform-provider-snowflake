package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWarehouseCreate(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &CreateWarehouseOptions{
			name: AccountObjectIdentifier{
				name: "mywarehouse",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE WAREHOUSE "mywarehouse"`)
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
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE WAREHOUSE IF NOT EXISTS "completewarehouse" WAREHOUSE_TYPE = 'STANDARD' WAREHOUSE_SIZE = 'X4LARGE' MAX_CLUSTER_COUNT = 8 MIN_CLUSTER_COUNT = 3 SCALING_POLICY = 'ECONOMY' AUTO_SUSPEND = 1000 AUTO_RESUME = true INITIALLY_SUSPENDED = false RESOURCE_MONITOR = "myresmon" COMMENT = 'hello' ENABLE_QUERY_ACCELERATION = true QUERY_ACCELERATION_MAX_SCALE_FACTOR = 62 MAX_CONCURRENCY_LEVEL = 7 STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 29 STATEMENT_TIMEOUT_IN_SECONDS = 89 TAG ("db1"."schema1"."tag1" = 'v1', "db1"."schema1"."tag2" = 'v2')`)
	})
}

func TestWarehouseSizing(t *testing.T) {
	t.Run("validation: Min bigger than Max", func(t *testing.T) {
		opts := &CreateWarehouseOptions{
			name:            NewAccountObjectIdentifier("mywarehouse"),
			MaxClusterCount: Int(1),
			MinClusterCount: Int(2),
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("MinClusterCount must be less than or equal to MaxClusterCount"))
	})

	t.Run("validation: only min set", func(t *testing.T) {
		opts := &CreateWarehouseOptions{
			name:            NewAccountObjectIdentifier("mywarehouse"),
			MinClusterCount: Int(2),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE WAREHOUSE \"mywarehouse\" MIN_CLUSTER_COUNT = 2")
	})

	t.Run("Max equal Min", func(t *testing.T) {
		opts := &CreateWarehouseOptions{
			name:            NewAccountObjectIdentifier("mywarehouse"),
			MaxClusterCount: Int(2),
			MinClusterCount: Int(2),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE WAREHOUSE \"mywarehouse\" MAX_CLUSTER_COUNT = 2 MIN_CLUSTER_COUNT = 2")
	})

	t.Run("Max greater than Min", func(t *testing.T) {
		opts := &CreateWarehouseOptions{
			name:            NewAccountObjectIdentifier("mywarehouse"),
			MaxClusterCount: Int(2),
			MinClusterCount: Int(1),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE WAREHOUSE \"mywarehouse\" MAX_CLUSTER_COUNT = 2 MIN_CLUSTER_COUNT = 1")
	})

	t.Run("Allow large Min Max Values", func(t *testing.T) {
		opts := &CreateWarehouseOptions{
			name:            NewAccountObjectIdentifier("mywarehouse"),
			MaxClusterCount: Int(100),
			MinClusterCount: Int(11),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE WAREHOUSE \"mywarehouse\" MAX_CLUSTER_COUNT = 100 MIN_CLUSTER_COUNT = 11")
	})
}

// TODO: add validation tests
func TestWarehouseAlter(t *testing.T) {
	t.Run("with set params", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			Set: &WarehouseSet{
				WarehouseType:                   &WarehouseTypeSnowparkOptimized,
				WaitForCompletion:               Bool(false),
				MinClusterCount:                 Int(4),
				MaxClusterCount:                 Int(5),
				AutoSuspend:                     Int(200),
				ResourceMonitor:                 NewAccountObjectIdentifier("resmon"),
				EnableQueryAcceleration:         Bool(false),
				StatementQueuedTimeoutInSeconds: Int(1200),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" SET WAREHOUSE_TYPE = 'SNOWPARK-OPTIMIZED' WAIT_FOR_COMPLETION = false MAX_CLUSTER_COUNT = 5 MIN_CLUSTER_COUNT = 4 AUTO_SUSPEND = 200 RESOURCE_MONITOR = "resmon" ENABLE_QUERY_ACCELERATION = false STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 1200`)
	})

	t.Run("with set tag", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			SetTag: []TagAssociation{
				{
					Name:  NewSchemaObjectIdentifier("db", "schema", "tag1"),
					Value: "v1",
				},
				{
					Name:  NewSchemaObjectIdentifier("db", "schema", "tag2"),
					Value: "v2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" SET TAG "db"."schema"."tag1" = 'v1', "db"."schema"."tag2" = 'v2'`)
	})

	t.Run("with unset tag", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			UnsetTag: []ObjectIdentifier{
				NewSchemaObjectIdentifier("db", "schema", "tag1"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" UNSET TAG "db"."schema"."tag1"`)
	})

	t.Run("with unset params", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			Unset: &WarehouseUnset{
				MaxClusterCount: Bool(true),
				AutoResume:      Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" UNSET MAX_CLUSTER_COUNT, AUTO_RESUME`)
	})

	t.Run("rename", func(t *testing.T) {
		newName := NewAccountObjectIdentifier("newName")
		opts := &AlterWarehouseOptions{
			name:    NewAccountObjectIdentifier("oldName"),
			NewName: &newName,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "oldName" RENAME TO "newName"`)
	})

	t.Run("suspend", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name:    NewAccountObjectIdentifier("mywarehouse"),
			Suspend: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" SUSPEND`)
	})

	t.Run("resume", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name:        NewAccountObjectIdentifier("mywarehouse"),
			Resume:      Bool(true),
			IfSuspended: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" RESUME IF SUSPENDED`)
	})

	t.Run("abort all queries", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name:            NewAccountObjectIdentifier("mywarehouse"),
			AbortAllQueries: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" ABORT ALL QUERIES`)
	})

	t.Run("with set tag", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			SetTag: []TagAssociation{
				{
					Name:  NewSchemaObjectIdentifier("db1", "schema1", "tag1"),
					Value: "v1",
				},
				{
					Name:  NewSchemaObjectIdentifier("db2", "schema2", "tag2"),
					Value: "v2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" SET TAG "db1"."schema1"."tag1" = 'v1', "db2"."schema2"."tag2" = 'v2'`)
	})

	t.Run("with unset tag", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			UnsetTag: []ObjectIdentifier{
				NewSchemaObjectIdentifier("db1", "schema1", "tag1"),
				NewSchemaObjectIdentifier("db2", "schema2", "tag2"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" UNSET TAG "db1"."schema1"."tag1", "db2"."schema2"."tag2"`)
	})
}

func TestWarehouseDrop(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &DropWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP WAREHOUSE "mywarehouse"`)
	})

	t.Run("with if exists", func(t *testing.T) {
		opts := &DropWarehouseOptions{
			name:     NewAccountObjectIdentifier("mywarehouse"),
			IfExists: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP WAREHOUSE IF EXISTS "mywarehouse"`)
	})
}

func TestWarehouseShow(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		opts := &ShowWarehouseOptions{}
		assertOptsValidAndSQLEquals(t, opts, `SHOW WAREHOUSES`)
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowWarehouseOptions{
			Like: &Like{
				Pattern: String("mywarehouse"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW WAREHOUSES LIKE 'mywarehouse'`)
	})
}

func TestWarehouseDescribe(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &describeWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
		}
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE WAREHOUSE "mywarehouse"`)
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
