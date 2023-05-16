package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestWarehouseCreate(t *testing.T) {
	builder := testBuilder(t)

	t.Run("only name", func(t *testing.T) {
		opts := &WarehouseCreateOptions{
			name: AccountObjectIdentifier{
				name: "mywarehouse",
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`CREATE WAREHOUSE "mywarehouse"`,
			builder.sql(clauses...),
		)
	})

	t.Run("with complete options", func(t *testing.T) {
		tag1 := randomSchemaObjectIdentifier(t)
		tag2 := randomSchemaObjectIdentifier(t)
		opts := &WarehouseCreateOptions{
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
					Name:  tag1,
					Value: "v1",
				},
				{
					Name:  tag2,
					Value: "v2",
				},
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`CREATE OR REPLACE WAREHOUSE IF NOT EXISTS "completewarehouse" WAREHOUSE_TYPE = 'STANDARD' WAREHOUSE_SIZE = 'X4LARGE' MAX_CLUSTER_COUNT = 8 MIN_CLUSTER_COUNT = 3 SCALING_POLICY = 'ECONOMY' AUTO_SUSPEND = 1000 AUTO_RESUME = true INITIALLY_SUSPENDED = false RESOURCE_MONITOR = "myresmon" COMMENT = 'hello' ENABLE_QUERY_ACCELERATION = true QUERY_ACCELERATION_MAX_SCALE_FACTOR = 62 MAX_CONCURRENCY_LEVEL = 7 STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 29 STATEMENT_TIMEOUT_IN_SECONDS = 89 TAG (`+tag1.FullyQualifiedName()+` = 'v1',`+tag2.FullyQualifiedName()+` = 'v2')`,
			builder.sql(clauses...),
		)
	})
}

func TestWarehouseAlter(t *testing.T) {
	builder := testBuilder(t)

	t.Run("with set", func(t *testing.T) {
		opts := &WarehouseAlterOptions{
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
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`ALTER WAREHOUSE "mywarehouse" SET WAREHOUSE_TYPE = 'SNOWPARK-OPTIMIZED' WAIT_FOR_COMPLETION = false MIN_CLUSTER_COUNT = 4 AUTO_SUSPEND = 200 RESOURCE_MONITOR = "resmon" ENABLE_QUERY_ACCELERATION = false STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 1200`,
			builder.sql(clauses...),
		)
	})

	t.Run("with unset", func(t *testing.T) {
		opts := &WarehouseAlterOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			Unset: &WarehouseUnset{
				WarehouseSize:   Bool(true),
				MaxClusterCount: Bool(true),
				AutoResume:      Bool(true),
				Tag: []ObjectIdentifier{
					NewSchemaObjectIdentifier("db", "schema", "tag1"),
				},
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`ALTER WAREHOUSE "mywarehouse" UNSET WAREHOUSE_SIZE,MAX_CLUSTER_COUNT,AUTO_RESUME,TAG "db"."schema"."tag1"`,
			builder.sql(clauses...),
		)
	})

	t.Run("rename", func(t *testing.T) {
		newname := NewAccountObjectIdentifier("newname")
		opts := &WarehouseAlterOptions{
			name:    NewAccountObjectIdentifier("oldname"),
			NewName: &newname,
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)

		assert.Equal(t,
			`ALTER WAREHOUSE "oldname" RENAME TO "newname"`,
			builder.sql(clauses...),
		)
	})

	t.Run("suspend", func(t *testing.T) {
		opts := &WarehouseAlterOptions{
			name:    NewAccountObjectIdentifier("mywarehouse"),
			Suspend: Bool(true),
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)

		assert.Equal(t,
			`ALTER WAREHOUSE "mywarehouse" SUSPEND`,
			builder.sql(clauses...),
		)
	})

	t.Run("resume", func(t *testing.T) {
		opts := &WarehouseAlterOptions{
			name:        NewAccountObjectIdentifier("mywarehouse"),
			Resume:      Bool(true),
			IfSuspended: Bool(true),
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)

		assert.Equal(t,
			`ALTER WAREHOUSE "mywarehouse" RESUME IF SUSPENDED`,
			builder.sql(clauses...),
		)
	})

	t.Run("abort all queries", func(t *testing.T) {
		opts := &WarehouseAlterOptions{
			name:            NewAccountObjectIdentifier("mywarehouse"),
			AbortAllQueries: Bool(true),
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)

		assert.Equal(t,
			`ALTER WAREHOUSE "mywarehouse" ABORT ALL QUERIES`,
			builder.sql(clauses...),
		)
	})

	t.Run("with set tag", func(t *testing.T) {
		opts := &WarehouseAlterOptions{
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
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`ALTER WAREHOUSE "mywarehouse" SET TAG "db1"."schema1"."tag1" = 'v1',"db2"."schema2"."tag2" = 'v2'`,
			builder.sql(clauses...),
		)
	})

	t.Run("with unset tag", func(t *testing.T) {
		opts := &WarehouseAlterOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			Unset: &WarehouseUnset{
				Tag: []ObjectIdentifier{
					NewSchemaObjectIdentifier("db1", "schema1", "tag1"),
					NewSchemaObjectIdentifier("db2", "schema2", "tag2"),
				},
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`ALTER WAREHOUSE "mywarehouse" UNSET TAG "db1"."schema1"."tag1","db2"."schema2"."tag2"`,
			builder.sql(clauses...),
		)
	})
}

func TestWarehouseDrop(t *testing.T) {
	builder := testBuilder(t)

	t.Run("only name", func(t *testing.T) {
		opts := &WarehouseDropOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)

		assert.Equal(t,
			`DROP WAREHOUSE "mywarehouse"`,
			builder.sql(clauses...),
		)
	})

	t.Run("with if exists", func(t *testing.T) {
		opts := &WarehouseDropOptions{
			name:     NewAccountObjectIdentifier("mywarehouse"),
			IfExists: Bool(true),
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)

		assert.Equal(t,
			`DROP WAREHOUSE IF EXISTS "mywarehouse"`,
			builder.sql(clauses...),
		)
	})
}

func TestWarehouseShow(t *testing.T) {
	builder := testBuilder(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &WarehouseShowOptions{}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			"SHOW WAREHOUSES",
			builder.sql(clauses...),
		)
	})

	t.Run("with like", func(t *testing.T) {
		opts := &WarehouseShowOptions{
			Like: &Like{
				Pattern: String("mywarehouse"),
			},
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			"SHOW WAREHOUSES LIKE 'mywarehouse'",
			builder.sql(clauses...),
		)
	})
}

func TestWarehouseDescribe(t *testing.T) {
	builder := testBuilder(t)

	t.Run("only name", func(t *testing.T) {
		opts := &warehouseDescribeOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
		}
		clauses, err := builder.parseStruct(opts)
		require.NoError(t, err)
		assert.Equal(t,
			`DESCRIBE WAREHOUSE "mywarehouse"`,
			builder.sql(clauses...),
		)
	})
}
