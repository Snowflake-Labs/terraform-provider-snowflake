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

	testWarehouse, warehouseCleanup := createWarehouseWithOptions(t, client, &CreateWarehouseOptions{
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
		showOptions := &ShowWarehouseOptions{
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
		showOptions := &ShowWarehouseOptions{
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
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)
	tagTest, tagCleanup := createTag(t, client, databaseTest, schemaTest)
	t.Cleanup(tagCleanup)
	tag2Test, tag2Cleanup := createTag(t, client, databaseTest, schemaTest)
	t.Cleanup(tag2Cleanup)

	t.Run("test complete", func(t *testing.T) {
		id := randomAccountObjectIdentifier(t)
		err := client.Warehouses.Create(ctx, id, &CreateWarehouseOptions{
			OrReplace:                       Bool(true),
			WarehouseType:                   &WarehouseTypeStandard,
			WarehouseSize:                   &WarehouseSizeSmall,
			MaxClusterCount:                 Int(8),
			MinClusterCount:                 Int(2),
			ScalingPolicy:                   &ScalingPolicyEconomy,
			AutoSuspend:                     Int(1000),
			AutoResume:                      Bool(true),
			InitiallySuspended:              Bool(false),
			Comment:                         String("comment"),
			EnableQueryAcceleration:         Bool(true),
			QueryAccelerationMaxScaleFactor: Int(90),
			MaxConcurrencyLevel:             Int(10),
			StatementQueuedTimeoutInSeconds: Int(2000),
			StatementTimeoutInSeconds:       Int(3000),
			Tag: []TagAssociation{
				{
					Name:  tagTest.ID(),
					Value: "v1",
				},
				{
					Name:  tag2Test.ID(),
					Value: "v2",
				},
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Warehouses.Drop(ctx, id, &DropWarehouseOptions{
				IfExists: Bool(true),
			})
			require.NoError(t, err)
		})
		warehouses, err := client.Warehouses.Show(ctx, &ShowWarehouseOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(warehouses))
		warehouse := warehouses[0]
		assert.Equal(t, id.Name(), warehouse.Name)
		assert.Equal(t, WarehouseTypeStandard, warehouse.Type)
		assert.Equal(t, WarehouseSizeSmall, warehouse.Size)
		assert.Equal(t, 8, warehouse.MaxClusterCount)
		assert.Equal(t, 2, warehouse.MinClusterCount)
		assert.Equal(t, ScalingPolicyEconomy, warehouse.ScalingPolicy)
		assert.Equal(t, 1000, warehouse.AutoSuspend)
		assert.Equal(t, true, warehouse.AutoResume)
		assert.Contains(t, []WarehouseState{WarehouseStateResuming, WarehouseStateStarted}, warehouse.State)
		assert.Equal(t, "comment", warehouse.Comment)
		assert.Equal(t, true, warehouse.EnableQueryAcceleration)
		assert.Equal(t, 90, warehouse.QueryAccelerationMaxScaleFactor)

		tag1Value, err := client.SystemFunctions.GetTag(ctx, tagTest.ID(), warehouse.ID(), ObjectTypeWarehouse)
		require.NoError(t, err)
		assert.Equal(t, "v1", tag1Value)
		tag2Value, err := client.SystemFunctions.GetTag(ctx, tag2Test.ID(), warehouse.ID(), ObjectTypeWarehouse)
		require.NoError(t, err)
		assert.Equal(t, "v2", tag2Value)
	})

	t.Run("test no options", func(t *testing.T) {
		id := randomAccountObjectIdentifier(t)
		err := client.Warehouses.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Warehouses.Drop(ctx, id, &DropWarehouseOptions{
				IfExists: Bool(true),
			})
			require.NoError(t, err)
		})
		warehouses, err := client.Warehouses.Show(ctx, &ShowWarehouseOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		assert.Equal(t, id.Name(), result.Name)
		assert.Equal(t, WarehouseTypeStandard, result.Type)
		assert.Equal(t, WarehouseSizeXSmall, result.Size)
		assert.Equal(t, 1, result.MaxClusterCount)
		assert.Equal(t, 1, result.MinClusterCount)
		assert.Equal(t, ScalingPolicyStandard, result.ScalingPolicy)
		assert.Equal(t, 600, result.AutoSuspend)
		assert.Equal(t, true, result.AutoResume)
		assert.Contains(t, []WarehouseState{WarehouseStateResuming, WarehouseStateStarted}, result.State)
		assert.Equal(t, "", result.Comment)
		assert.Equal(t, false, result.EnableQueryAcceleration)
		assert.Equal(t, 8, result.QueryAccelerationMaxScaleFactor)
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

	t.Run("terraform acc test", func(t *testing.T) {
		id := randomAccountObjectIdentifier(t)
		opts := &CreateWarehouseOptions{
			Comment:            String("test comment"),
			WarehouseSize:      &WarehouseSizeXSmall,
			AutoSuspend:        Int(60),
			MaxClusterCount:    Int(1),
			MinClusterCount:    Int(1),
			ScalingPolicy:      &ScalingPolicyStandard,
			AutoResume:         Bool(true),
			InitiallySuspended: Bool(true),
		}
		err := client.Warehouses.Create(ctx, id, opts)
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Warehouses.Drop(ctx, id, &DropWarehouseOptions{
				IfExists: Bool(true),
			})
			require.NoError(t, err)
		})
		warehouse, err := client.Warehouses.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, 1, warehouse.MaxClusterCount)
		assert.Equal(t, 1, warehouse.MinClusterCount)
		assert.Equal(t, ScalingPolicyStandard, warehouse.ScalingPolicy)
		assert.Equal(t, 60, warehouse.AutoSuspend)
		assert.Equal(t, true, warehouse.AutoResume)
		assert.Equal(t, "test comment", warehouse.Comment)
		assert.Equal(t, WarehouseStateSuspended, warehouse.State)
		assert.Equal(t, WarehouseSizeXSmall, warehouse.Size)

		// rename
		newID := randomAccountObjectIdentifier(t)
		alterOptions := &AlterWarehouseOptions{
			NewName: newID,
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		warehouse, err = client.Warehouses.ShowByID(ctx, newID)
		require.NoError(t, err)
		assert.Equal(t, newID.Name(), warehouse.Name)

		// change props
		alterOptions = &AlterWarehouseOptions{
			Set: &WarehouseSet{
				WarehouseSize: &WarehouseSizeSmall,
				Comment:       String("test comment2"),
			},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		warehouse, err = client.Warehouses.ShowByID(ctx, newID)
		require.NoError(t, err)
		assert.Equal(t, "test comment2", warehouse.Comment)
		assert.Equal(t, WarehouseSizeSmall, warehouse.Size)
	})

	t.Run("set", func(t *testing.T) {
		warehouse, warehouseCleanup := createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		alterOptions := &AlterWarehouseOptions{
			Set: &WarehouseSet{
				WarehouseSize:           &WarehouseSizeMedium,
				AutoSuspend:             Int(1234),
				EnableQueryAcceleration: Bool(true),
			},
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		warehouses, err := client.Warehouses.Show(ctx, &ShowWarehouseOptions{
			Like: &Like{
				Pattern: String(warehouse.Name),
			},
		})
		assert.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		require.NoError(t, err)
		assert.Equal(t, WarehouseSizeMedium, result.Size)
		assert.Equal(t, true, result.EnableQueryAcceleration)
		assert.Equal(t, 1234, result.AutoSuspend)
	})

	t.Run("rename", func(t *testing.T) {
		warehouse, warehouseCleanup := createWarehouse(t, client)
		oldID := warehouse.ID()
		t.Cleanup(warehouseCleanup)

		newID := randomAccountObjectIdentifier(t)
		alterOptions := &AlterWarehouseOptions{
			NewName: newID,
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		result, err := client.Warehouses.Describe(ctx, newID)
		require.NoError(t, err)
		assert.Equal(t, newID.Name(), result.Name)

		// rename back to original name so it can be cleaned up
		alterOptions = &AlterWarehouseOptions{
			NewName: oldID,
		}
		err = client.Warehouses.Alter(ctx, newID, alterOptions)
		require.NoError(t, err)
	})

	t.Run("unset", func(t *testing.T) {
		createOptions := &CreateWarehouseOptions{
			Comment:         String("test comment"),
			MaxClusterCount: Int(10),
		}
		warehouse, warehouseCleanup := createWarehouseWithOptions(t, client, createOptions)
		t.Cleanup(warehouseCleanup)
		id := warehouse.ID()

		alterOptions := &AlterWarehouseOptions{
			Unset: &WarehouseUnset{
				Comment:         Bool(true),
				MaxClusterCount: Bool(true),
			},
		}
		err := client.Warehouses.Alter(ctx, id, alterOptions)
		require.NoError(t, err)
		warehouses, err := client.Warehouses.Show(ctx, &ShowWarehouseOptions{
			Like: &Like{
				Pattern: String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		assert.Equal(t, warehouse.Name, result.Name)
		assert.Equal(t, "", result.Comment)
		assert.Equal(t, 1, result.MaxClusterCount)
	})

	t.Run("suspend & resume", func(t *testing.T) {
		warehouse, warehouseCleanup := createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		alterOptions := &AlterWarehouseOptions{
			Suspend: Bool(true),
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		warehouses, err := client.Warehouses.Show(ctx, &ShowWarehouseOptions{
			Like: &Like{
				Pattern: String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		assert.Contains(t, []WarehouseState{WarehouseStateSuspended, WarehouseStateSuspending}, result.State)

		alterOptions = &AlterWarehouseOptions{
			Resume: Bool(true),
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		warehouses, err = client.Warehouses.Show(ctx, &ShowWarehouseOptions{
			Like: &Like{
				Pattern: String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result = warehouses[0]
		assert.Contains(t, []WarehouseState{WarehouseStateStarted, WarehouseStateResuming}, result.State)
	})

	t.Run("resume without suspending", func(t *testing.T) {
		warehouse, warehouseCleanup := createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		alterOptions := &AlterWarehouseOptions{
			Resume:      Bool(true),
			IfSuspended: Bool(true),
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		warehouses, err := client.Warehouses.Show(ctx, &ShowWarehouseOptions{
			Like: &Like{
				Pattern: String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		assert.Contains(t, []WarehouseState{WarehouseStateStarted, WarehouseStateResuming}, result.State)
	})

	t.Run("abort all queries", func(t *testing.T) {
		warehouse, warehouseCleanup := createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		resetWarehouse := useWarehouse(t, client, warehouse.ID())
		t.Cleanup(resetWarehouse)

		// Start a long query
		go client.exec(ctx, "CALL SYSTEM$WAIT(30);") //nolint:errcheck // we don't care if this eventually errors, as long as it runs for a little while
		time.Sleep(5 * time.Second)

		// Check that query is running
		warehouses, err := client.Warehouses.Show(ctx, &ShowWarehouseOptions{
			Like: &Like{
				Pattern: String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		assert.Equal(t, 1, result.Running)
		assert.Equal(t, 0, result.Queued)

		// Abort all queries
		alterOptions := &AlterWarehouseOptions{
			AbortAllQueries: Bool(true),
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		// Wait for abort to be effective
		time.Sleep(5 * time.Second)

		// Check no query is running
		warehouses, err = client.Warehouses.Show(ctx, &ShowWarehouseOptions{
			Like: &Like{
				Pattern: String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result = warehouses[0]
		assert.Equal(t, 0, result.Running)
		assert.Equal(t, 0, result.Queued)
	})

	t.Run("set tags", func(t *testing.T) {
		warehouse, warehouseCleanup := createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		alterOptions := &AlterWarehouseOptions{
			Set: &WarehouseSet{
				Tag: []TagAssociation{
					{
						Name:  tag.ID(),
						Value: "val",
					},
					{
						Name:  tag2.ID(),
						Value: "val2",
					},
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
		warehouse, warehouseCleanup := createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		alterOptions := &AlterWarehouseOptions{
			Set: &WarehouseSet{
				Tag: []TagAssociation{
					{
						Name:  tag.ID(),
						Value: "val1",
					},
					{
						Name:  tag2.ID(),
						Value: "val2",
					},
				},
			},
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		val, err := client.SystemFunctions.GetTag(ctx, tag.ID(), warehouse.ID(), ObjectTypeWarehouse)
		require.NoError(t, err)
		require.Equal(t, "val1", val)
		val2, err := client.SystemFunctions.GetTag(ctx, tag2.ID(), warehouse.ID(), ObjectTypeWarehouse)
		require.NoError(t, err)
		require.Equal(t, "val2", val2)

		alterOptions = &AlterWarehouseOptions{
			Unset: &WarehouseUnset{
				Tag: []ObjectIdentifier{
					tag.ID(),
					tag2.ID(),
				},
			},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		val, err = client.SystemFunctions.GetTag(ctx, tag.ID(), warehouse.ID(), ObjectTypeWarehouse)
		require.Error(t, err)
		require.Equal(t, "", val)
		val2, err = client.SystemFunctions.GetTag(ctx, tag2.ID(), warehouse.ID(), ObjectTypeWarehouse)
		require.Error(t, err)
		require.Equal(t, "", val2)
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

		dropOptions := &DropWarehouseOptions{IfExists: Bool(true)}
		err := client.Warehouses.Drop(ctx, warehouse.ID(), dropOptions)
		require.NoError(t, err)
		_, err = client.Warehouses.Describe(ctx, warehouse.ID())
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})
}
