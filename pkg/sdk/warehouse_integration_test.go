package sdk

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_WarehousesShow(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	testWarehouse, warehouseCleanup := createWarehouseWithOptions(t, client, &WarehouseCreateOptions{
		WarehouseSize: &WarehouseSizeSmall,
	})
	t.Cleanup(warehouseCleanup)
	_, warehouse2Cleanup := createWarehouse(t, client)
	t.Cleanup(warehouse2Cleanup)

	t.Run("show without options", func(t *testing.T) {
		warehouses, err := client.Warehouses.Show(ctx, nil)
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(warehouses))
	})

	t.Run("show with options", func(t *testing.T) {
		showOptions := &WarehouseShowOptions{
			Like: &Like{
				Pattern: &testWarehouse.Name,
			},
		}
		warehouses, err := client.Warehouses.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		assert.Equal(t, testWarehouse.Name, warehouses[0].Name)
		assert.Equal(t, WarehouseSizeSmall, warehouses[0].Size)
	})

	t.Run("when searching a non-existent password policy", func(t *testing.T) {
		showOptions := &WarehouseShowOptions{
			Like: &Like{
				Pattern: String("non-existent"),
			},
		}
		warehouses, err := client.Warehouses.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 0, len(warehouses))
	})
}

func TestInt_WarehouseCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	database, dbCleanup := createDatabase(t, client)
	t.Cleanup(dbCleanup)
	schema, schemaCleanup := createSchema(t, client, database)
	t.Cleanup(schemaCleanup)
	tag, tagCleanup := createTag(t, client, database, schema)
	t.Cleanup(tagCleanup)

	t.Run("test complete", func(t *testing.T) {
		name := randomUUID(t)
		id := NewAccountObjectIdentifier(name)
		err := client.Warehouses.Create(ctx, id, &WarehouseCreateOptions{
			OrReplace:                       Bool(true),
			WarehouseType:                   &WarehouseTypeStandard,
			WarehouseSize:                   &WarehouseSizeSmall,
			MaxClusterCount:                 Uint8(8),
			MinClusterCount:                 Uint8(2),
			ScalingPolicy:                   &ScalingPolicyEconomy,
			AutoSuspend:                     Uint(1000),
			AutoResume:                      Bool(true),
			InitiallySuspended:              Bool(false),
			Comment:                         String("comment"),
			EnableQueryAcceleration:         Bool(true),
			QueryAccelerationMaxScaleFactor: Uint8(90),
			MaxConcurrencyLevel:             Uint(10),
			StatementQueuedTimeoutInSeconds: Uint(2000),
			StatementTimeoutInSeconds:       Uint(3000),
			Tags: []TagAssociation{
				{
					Name:  tag.ID(),
					Value: "myval",
				},
			},
		})
		require.NoError(t, err)
		warehouses, err := client.Warehouses.Show(ctx, &WarehouseShowOptions{
			Like: &Like{
				Pattern: String(name),
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		assert.Equal(t, name, result.Name)
		assert.Equal(t, WarehouseTypeStandard, result.Type)
		assert.Equal(t, WarehouseSizeSmall, result.Size)
		assert.Equal(t, uint8(8), result.MaxClusterCount)
		assert.Equal(t, uint8(2), result.MinClusterCount)
		assert.Equal(t, ScalingPolicyEconomy, result.ScalingPolicy)
		assert.Equal(t, uint(1000), result.AutoSuspend)
		assert.Equal(t, true, result.AutoResume)
		assert.Contains(t, []string{"RESUMING", "STARTED"}, result.State)
		assert.Equal(t, "comment", result.Comment)
		assert.Equal(t, true, result.EnableQueryAcceleration)
		assert.Equal(t, uint8(90), result.QueryAccelerationMaxScaleFactor)

		val, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, ObjectTypeWarehouse)
		require.NoError(t, err)
		require.Equal(t, "myval", val)
	})

	t.Run("test no options", func(t *testing.T) {
		name := randomUUID(t)
		id := NewAccountObjectIdentifier(name)
		err := client.Warehouses.Create(ctx, id, nil)
		require.NoError(t, err)
		warehouses, err := client.Warehouses.Show(ctx, &WarehouseShowOptions{
			Like: &Like{
				Pattern: String(name),
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		assert.Equal(t, name, result.Name)
		assert.Equal(t, WarehouseTypeStandard, result.Type)
		assert.Equal(t, WarehouseSizeXSmall, result.Size)
		assert.Equal(t, uint8(1), result.MaxClusterCount)
		assert.Equal(t, uint8(1), result.MinClusterCount)
		assert.Equal(t, ScalingPolicyStandard, result.ScalingPolicy)
		assert.Equal(t, uint(600), result.AutoSuspend)
		assert.Equal(t, true, result.AutoResume)
		assert.Contains(t, []string{"RESUMING", "STARTED"}, result.State)
		assert.Equal(t, "", result.Comment)
		assert.Equal(t, false, result.EnableQueryAcceleration)
		assert.Equal(t, uint8(8), result.QueryAccelerationMaxScaleFactor)
	})
}

func TestInt_WarehouseDescribe(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	warehouse, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)

	t.Run("when warehouse exists", func(t *testing.T) {
		result, err := client.Warehouses.Describe(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, warehouse.Name, result.Name)
		assert.Equal(t, "WAREHOUSE", result.Kind)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 5*time.Second)
	})

	t.Run("when warehouse does not exist", func(t *testing.T) {
		id := NewAccountObjectIdentifier("does_not_exist")
		_, err := client.Warehouses.Describe(ctx, id)
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_WarehouseAlter(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	database, dbCleanup := createDatabase(t, client)
	t.Cleanup(dbCleanup)
	schema, schemaCleanup := createSchema(t, client, database)
	t.Cleanup(schemaCleanup)
	tag, tagCleanup := createTag(t, client, database, schema)
	t.Cleanup(tagCleanup)
	tag2, tagCleanup2 := createTag(t, client, database, schema)
	t.Cleanup(tagCleanup2)

	t.Run("set", func(t *testing.T) {
		warehouse, warehouseCleanup := createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		alterOptions := &WarehouseAlterOptions{
			Set: &WarehouseSetOptions{
				WarehouseSize:           &WarehouseSizeMedium,
				AutoSuspend:             Uint(1234),
				EnableQueryAcceleration: Bool(true),
			},
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		warehouses, err := client.Warehouses.Show(ctx, &WarehouseShowOptions{
			Like: &Like{
				Pattern: String(warehouse.Name),
			},
		})
		assert.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		require.NoError(t, err)
		assert.Equal(t, WarehouseSizeMedium, result.Size)
		assert.Equal(t, true, result.EnableQueryAcceleration)
		assert.Equal(t, uint(1234), result.AutoSuspend)
	})

	t.Run("rename", func(t *testing.T) {
		warehouse, warehouseCleanup := createWarehouse(t, client)
		oldID := warehouse.ID()
		t.Cleanup(warehouseCleanup)

		newName := randomUUID(t)
		newID := NewAccountObjectIdentifier(newName)
		alterOptions := &WarehouseAlterOptions{
			NewName: &newID,
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		result, err := client.Warehouses.Describe(ctx, newID)
		require.NoError(t, err)
		assert.Equal(t, newName, result.Name)

		// rename back to original name so it can be cleaned up
		alterOptions = &WarehouseAlterOptions{
			NewName: &oldID,
		}
		err = client.Warehouses.Alter(ctx, newID, alterOptions)
		require.NoError(t, err)
	})

	t.Run("unset", func(t *testing.T) {
		createOptions := &WarehouseCreateOptions{
			Comment:         String("test comment"),
			MaxClusterCount: Uint8(10),
		}
		warehouse, warehouseCleanup := createWarehouseWithOptions(t, client, createOptions)
		t.Cleanup(warehouseCleanup)
		id := warehouse.ID()

		alterOptions := &WarehouseAlterOptions{
			Unset: &[]WarehouseUnsetField{
				CommentField,
				MaxClusterCountField,
			},
		}
		err := client.Warehouses.Alter(ctx, id, alterOptions)
		require.NoError(t, err)
		warehouses, err := client.Warehouses.Show(ctx, &WarehouseShowOptions{
			Like: &Like{
				Pattern: String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		assert.Equal(t, warehouse.Name, result.Name)
		assert.Equal(t, "", result.Comment)
		assert.Equal(t, uint8(1), result.MaxClusterCount)
	})

	t.Run("suspend & resume", func(t *testing.T) {
		warehouse, warehouseCleanup := createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		alterOptions := &WarehouseAlterOptions{
			Suspend: Bool(true),
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		warehouses, err := client.Warehouses.Show(ctx, &WarehouseShowOptions{
			Like: &Like{
				Pattern: String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		assert.Contains(t, []string{"SUSPENDING", "SUSPENDED"}, result.State)

		alterOptions = &WarehouseAlterOptions{
			Resume: Bool(true),
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		warehouses, err = client.Warehouses.Show(ctx, &WarehouseShowOptions{
			Like: &Like{
				Pattern: String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result = warehouses[0]
		assert.Contains(t, []string{"RESUMING", "STARTED"}, result.State)
	})

	t.Run("resume without suspending", func(t *testing.T) {
		warehouse, warehouseCleanup := createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		alterOptions := &WarehouseAlterOptions{
			Resume:      Bool(true),
			IfSuspended: Bool(true),
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		warehouses, err := client.Warehouses.Show(ctx, &WarehouseShowOptions{
			Like: &Like{
				Pattern: String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		assert.Contains(t, []string{"STARTED", "RESUMING"}, result.State)
	})

	t.Run("abort all queries", func(t *testing.T) {
		warehouse, warehouseCleanup := createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		resetWarehouse := useWarehouse(t, client, warehouse.ID())
		t.Cleanup(resetWarehouse)

		// Start a long query
		go client.exec(ctx, "CALL SYSTEM$WAIT(30);")
		time.Sleep(5 * time.Second)

		// Check that query is running
		warehouses, err := client.Warehouses.Show(ctx, &WarehouseShowOptions{
			Like: &Like{
				Pattern: String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		assert.Equal(t, uint(1), result.Running)
		assert.Equal(t, uint(0), result.Queued)

		// Abort all queries
		alterOptions := &WarehouseAlterOptions{
			AbortAllQueries: Bool(true),
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		// Wait for abort to be effective
		time.Sleep(5 * time.Second)

		// Check no query is running
		warehouses, err = client.Warehouses.Show(ctx, &WarehouseShowOptions{
			Like: &Like{
				Pattern: String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result = warehouses[0]
		assert.Equal(t, uint(0), result.Running)
		assert.Equal(t, uint(0), result.Queued)
	})

	t.Run("set tags", func(t *testing.T) {
		warehouse, warehouseCleanup := createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		alterOptions := &WarehouseAlterOptions{
			SetTags: &[]TagAssociation{
				{
					Name:  tag.ID(),
					Value: "val",
				},
				{
					Name:  tag2.ID(),
					Value: "val2",
				},
			},
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		val, err := client.SystemFunctions.GetTag(ctx, tag.ID(), warehouse.ID(), ObjectTypeWarehouse)
		require.NoError(t, err)
		require.Equal(t, "val", val)
		val, err = client.SystemFunctions.GetTag(ctx, tag2.ID(), warehouse.ID(), ObjectTypeWarehouse)
		require.NoError(t, err)
		require.Equal(t, "val2", val)
	})

	t.Run("unset tags", func(t *testing.T) {
		warehouse, warehouseCleanup := createWarehouseWithOptions(t, client, &WarehouseCreateOptions{
			Tags: []TagAssociation{
				{
					Name:  tag.ID(),
					Value: "value",
				},
				{
					Name:  tag2.ID(),
					Value: "value2",
				},
			},
		})
		t.Cleanup(warehouseCleanup)

		alterOptions := &WarehouseAlterOptions{
			UnsetTags: &[]ObjectIdentifier{
				tag.ID(),
				tag2.ID(),
			},
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		val, err := client.SystemFunctions.GetTag(ctx, tag.ID(), warehouse.ID(), ObjectTypeWarehouse)
		require.Error(t, err)
		require.Equal(t, "", val)
		val, err = client.SystemFunctions.GetTag(ctx, tag2.ID(), warehouse.ID(), ObjectTypeWarehouse)
		require.Error(t, err)
		require.Equal(t, "", val)
	})
}

func TestInt_WarehouseDrop(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	t.Run("when warehouse exists", func(t *testing.T) {
		warehouse, _ := createWarehouse(t, client)

		err := client.Warehouses.Drop(ctx, warehouse.ID(), nil)
		require.NoError(t, err)
		_, err = client.Warehouses.Describe(ctx, warehouse.ID())
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})

	t.Run("when warehouse does not exist", func(t *testing.T) {
		id := NewAccountObjectIdentifier("does_not_exist")
		err := client.Warehouses.Drop(ctx, id, nil)
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})

	t.Run("when warehouse exists and if exists is true", func(t *testing.T) {
		warehouse, _ := createWarehouse(t, client)

		dropOptions := &WarehouseDropOptions{IfExists: Bool(true)}
		err := client.Warehouses.Drop(ctx, warehouse.ID(), dropOptions)
		require.NoError(t, err)
		_, err = client.Warehouses.Describe(ctx, warehouse.ID())
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})
}
