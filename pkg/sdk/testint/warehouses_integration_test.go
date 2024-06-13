package testint

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-1348102 - next PR]: more unset tests
// TODO [SNOW-1348102 - next PR]: test how suspension.resuming works for different states
// TODO [this PR]: show -> showbyid in multiple tests
func TestInt_Warehouses(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	prefix := random.StringN(6)
	precreatedWarehouseId := testClientHelper().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	precreatedWarehouseId2 := testClientHelper().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	// new warehouses created on purpose
	_, precreatedWarehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithOptions(t, precreatedWarehouseId, nil)
	t.Cleanup(precreatedWarehouseCleanup)
	_, precreatedWarehouse2Cleanup := testClientHelper().Warehouse.CreateWarehouseWithOptions(t, precreatedWarehouseId2, nil)
	t.Cleanup(precreatedWarehouse2Cleanup)

	tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)
	tag2, tag2Cleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tag2Cleanup)

	resourceMonitor, resourceMonitorCleanup := testClientHelper().ResourceMonitor.CreateResourceMonitor(t)
	t.Cleanup(resourceMonitorCleanup)

	t.Run("show: without options", func(t *testing.T) {
		warehouses, err := client.Warehouses.Show(ctx, nil)
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(warehouses))
	})

	t.Run("show: like", func(t *testing.T) {
		showOptions := &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.Pointer(prefix + "%"),
			},
		}
		warehouses, err := client.Warehouses.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Len(t, warehouses, 2)
	})

	t.Run("show: with options", func(t *testing.T) {
		showOptions := &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.Pointer(precreatedWarehouseId.Name()),
			},
		}
		warehouses, err := client.Warehouses.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		assert.Equal(t, precreatedWarehouseId.Name(), warehouses[0].Name)
		assert.Equal(t, sdk.WarehouseSizeXSmall, warehouses[0].Size)
		assert.Equal(t, "ROLE", warehouses[0].OwnerRoleType)
	})

	t.Run("show: when searching a non-existent warehouse", func(t *testing.T) {
		showOptions := &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		}
		warehouses, err := client.Warehouses.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Len(t, warehouses, 0)
	})

	t.Run("create: complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.Create(ctx, id, &sdk.CreateWarehouseOptions{
			OrReplace:                       sdk.Bool(true),
			WarehouseType:                   &sdk.WarehouseTypeStandard,
			WarehouseSize:                   &sdk.WarehouseSizeSmall,
			MaxClusterCount:                 sdk.Int(8),
			MinClusterCount:                 sdk.Int(2),
			ScalingPolicy:                   &sdk.ScalingPolicyEconomy,
			AutoSuspend:                     sdk.Int(1000),
			AutoResume:                      sdk.Bool(true),
			InitiallySuspended:              sdk.Bool(false),
			ResourceMonitor:                 sdk.Pointer(resourceMonitor.ID()),
			Comment:                         sdk.String("comment"),
			EnableQueryAcceleration:         sdk.Bool(true),
			QueryAccelerationMaxScaleFactor: sdk.Int(90),
			MaxConcurrencyLevel:             sdk.Int(10),
			StatementQueuedTimeoutInSeconds: sdk.Int(2000),
			StatementTimeoutInSeconds:       sdk.Int(3000),
			Tag: []sdk.TagAssociation{
				{
					Name:  tag.ID(),
					Value: "v1",
				},
				{
					Name:  tag2.ID(),
					Value: "v2",
				},
			},
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		warehouses, err := client.Warehouses.Show(ctx, &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(id.Name()),
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(warehouses))
		warehouse := warehouses[0]
		assert.Equal(t, id.Name(), warehouse.Name)
		assert.Equal(t, sdk.WarehouseTypeStandard, warehouse.Type)
		assert.Equal(t, sdk.WarehouseSizeSmall, warehouse.Size)
		assert.Equal(t, 8, warehouse.MaxClusterCount)
		assert.Equal(t, 2, warehouse.MinClusterCount)
		assert.Equal(t, sdk.ScalingPolicyEconomy, warehouse.ScalingPolicy)
		assert.Equal(t, 1000, warehouse.AutoSuspend)
		assert.Equal(t, true, warehouse.AutoResume)
		assert.Contains(t, []sdk.WarehouseState{sdk.WarehouseStateResuming, sdk.WarehouseStateStarted}, warehouse.State)
		assert.Equal(t, resourceMonitor.ID().Name(), warehouse.ResourceMonitor)
		assert.Equal(t, "comment", warehouse.Comment)
		assert.Equal(t, true, warehouse.EnableQueryAcceleration)
		assert.Equal(t, 90, warehouse.QueryAccelerationMaxScaleFactor)

		tag1Value, err := client.SystemFunctions.GetTag(ctx, tag.ID(), warehouse.ID(), sdk.ObjectTypeWarehouse)
		require.NoError(t, err)
		assert.Equal(t, "v1", tag1Value)
		tag2Value, err := client.SystemFunctions.GetTag(ctx, tag2.ID(), warehouse.ID(), sdk.ObjectTypeWarehouse)
		require.NoError(t, err)
		assert.Equal(t, "v2", tag2Value)
	})

	t.Run("create: no options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		warehouses, err := client.Warehouses.Show(ctx, &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(id.Name()),
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		assert.Equal(t, id.Name(), result.Name)
		assert.Equal(t, sdk.WarehouseTypeStandard, result.Type)
		assert.Equal(t, sdk.WarehouseSizeXSmall, result.Size)
		assert.Equal(t, 1, result.MaxClusterCount)
		assert.Equal(t, 1, result.MinClusterCount)
		assert.Equal(t, sdk.ScalingPolicyStandard, result.ScalingPolicy)
		assert.Equal(t, 600, result.AutoSuspend)
		assert.Equal(t, true, result.AutoResume)
		assert.Contains(t, []sdk.WarehouseState{sdk.WarehouseStateResuming, sdk.WarehouseStateStarted}, result.State)
		assert.Equal(t, "", result.Comment)
		assert.Equal(t, false, result.EnableQueryAcceleration)
		assert.Equal(t, 8, result.QueryAccelerationMaxScaleFactor)
	})

	t.Run("create: empty comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.Create(ctx, id, &sdk.CreateWarehouseOptions{Comment: sdk.String("")})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		warehouses, err := client.Warehouses.Show(ctx, &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(id.Name()),
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		assert.Equal(t, "", result.Comment)
	})

	t.Run("alter: set and unset", func(t *testing.T) {
		createOptions := &sdk.CreateWarehouseOptions{
			Comment:         sdk.String("test comment"),
			MaxClusterCount: sdk.Int(10),
		}
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithOptions(t, id, createOptions)
		t.Cleanup(warehouseCleanup)

		alterOptions := &sdk.AlterWarehouseOptions{
			Set: &sdk.WarehouseSet{
				ResourceMonitor:         resourceMonitor.ID(),
				WarehouseSize:           &sdk.WarehouseSizeMedium,
				AutoSuspend:             sdk.Int(1234),
				EnableQueryAcceleration: sdk.Bool(true),
			},
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		warehouses, err := client.Warehouses.Show(ctx, &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(warehouse.Name),
			},
		})
		assert.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseSizeMedium, result.Size)
		assert.Equal(t, true, result.EnableQueryAcceleration)
		assert.Equal(t, 1234, result.AutoSuspend)
		assert.Equal(t, "test comment", result.Comment)
		assert.Equal(t, 10, result.MaxClusterCount)

		alterOptions = &sdk.AlterWarehouseOptions{
			Unset: &sdk.WarehouseUnset{
				ResourceMonitor: sdk.Bool(true),
				Comment:         sdk.Bool(true),
				MaxClusterCount: sdk.Bool(true),
			},
		}
		err = client.Warehouses.Alter(ctx, id, alterOptions)
		require.NoError(t, err)

		warehouses, err = client.Warehouses.Show(ctx, &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result = warehouses[0]
		assert.Equal(t, warehouse.Name, result.Name)
		assert.Equal(t, "", result.Comment)
		assert.Equal(t, 1, result.MaxClusterCount)
		assert.Equal(t, sdk.WarehouseSizeMedium, result.Size)
		assert.Equal(t, true, result.EnableQueryAcceleration)
		assert.Equal(t, 1234, result.AutoSuspend)
	})

	t.Run("alter: prove problems with unset auto suspend", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		alterOptions := &sdk.AlterWarehouseOptions{
			Unset: &sdk.WarehouseUnset{AutoSuspend: sdk.Bool(true)},
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		returnedWarehouse, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		// TODO [SNOW-1473453]: change when UNSET starts working correctly (expecting to unset to default 600)
		// assert.Equal(t, "600", returnedWarehouse.AutoSuspend)
		assert.Equal(t, "0", returnedWarehouse.AutoSuspend)
	})

	t.Run("alter: prove problems with unset warehouse type", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		alterOptions := &sdk.AlterWarehouseOptions{
			Unset: &sdk.WarehouseUnset{WarehouseType: sdk.Bool(true)},
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		// TODO [SNOW-1473453]: change when UNSET starts working correctly (expecting to unset to default type STANDARD)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid type of property 'null' for 'WAREHOUSE_TYPE'")
	})

	t.Run("alter: rename", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		newID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		alterOptions := &sdk.AlterWarehouseOptions{
			NewName: &newID,
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, newID))

		result, err := client.Warehouses.Describe(ctx, newID)
		require.NoError(t, err)
		assert.Equal(t, newID.Name(), result.Name)
	})

	// This proves that we don't have to handle empty comment inside the resource.
	t.Run("alter: set empty comment versus unset", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.Create(ctx, id, &sdk.CreateWarehouseOptions{Comment: sdk.String("abc")})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		// can't use normal way, because of our SDK validation
		_, err = client.ExecForTests(ctx, fmt.Sprintf("ALTER WAREHOUSE %s SET COMMENT = ''", id.FullyQualifiedName()))
		require.NoError(t, err)

		warehouses, err := client.Warehouses.Show(ctx, &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(id.Name()),
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(warehouses))
		assert.Equal(t, "", warehouses[0].Comment)

		alterOptions := &sdk.AlterWarehouseOptions{
			Set: &sdk.WarehouseSet{
				Comment: sdk.String("abc"),
			},
		}
		err = client.Warehouses.Alter(ctx, id, alterOptions)
		require.NoError(t, err)

		alterOptions = &sdk.AlterWarehouseOptions{
			Unset: &sdk.WarehouseUnset{
				Comment: sdk.Bool(true),
			},
		}
		err = client.Warehouses.Alter(ctx, id, alterOptions)
		require.NoError(t, err)

		warehouses, err = client.Warehouses.Show(ctx, &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(id.Name()),
			},
		})
		require.NoError(t, err)
		require.Equal(t, 1, len(warehouses))
		assert.Equal(t, "", warehouses[0].Comment)
	})

	t.Run("alter: suspend and resume", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		alterOptions := &sdk.AlterWarehouseOptions{
			Suspend: sdk.Bool(true),
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		warehouses, err := client.Warehouses.Show(ctx, &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		assert.Contains(t, []sdk.WarehouseState{sdk.WarehouseStateSuspended, sdk.WarehouseStateSuspending}, result.State)

		alterOptions = &sdk.AlterWarehouseOptions{
			Resume: sdk.Bool(true),
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		warehouses, err = client.Warehouses.Show(ctx, &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result = warehouses[0]
		assert.Contains(t, []sdk.WarehouseState{sdk.WarehouseStateStarted, sdk.WarehouseStateResuming}, result.State)
	})

	t.Run("alter: resume without suspending", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		alterOptions := &sdk.AlterWarehouseOptions{
			Resume:      sdk.Bool(true),
			IfSuspended: sdk.Bool(true),
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		warehouses, err := client.Warehouses.Show(ctx, &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		assert.Contains(t, []sdk.WarehouseState{sdk.WarehouseStateStarted, sdk.WarehouseStateResuming}, result.State)
	})

	t.Run("alter: abort all queries", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		resetWarehouse := testClientHelper().Warehouse.UseWarehouse(t, warehouse.ID())
		t.Cleanup(resetWarehouse)

		// Start a long query
		go client.ExecForTests(ctx, "CALL SYSTEM$WAIT(30);") //nolint:errcheck // we don't care if this eventually errors, as long as it runs for a little while
		time.Sleep(5 * time.Second)

		// Check that query is running
		warehouses, err := client.Warehouses.Show(ctx, &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result := warehouses[0]
		assert.Equal(t, 1, result.Running)
		assert.Equal(t, 0, result.Queued)

		// Abort all queries
		alterOptions := &sdk.AlterWarehouseOptions{
			AbortAllQueries: sdk.Bool(true),
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		// Wait for abort to be effective
		time.Sleep(5 * time.Second)

		// Check no query is running
		warehouses, err = client.Warehouses.Show(ctx, &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(warehouse.Name),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		result = warehouses[0]
		assert.Equal(t, 0, result.Running)
		assert.Equal(t, 0, result.Queued)
	})

	t.Run("alter: set tags and unset tags", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		alterOptions := &sdk.AlterWarehouseOptions{
			SetTag: []sdk.TagAssociation{
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

		val, err := client.SystemFunctions.GetTag(ctx, tag.ID(), warehouse.ID(), sdk.ObjectTypeWarehouse)
		require.NoError(t, err)
		require.Equal(t, "val", val)
		val2, err := client.SystemFunctions.GetTag(ctx, tag2.ID(), warehouse.ID(), sdk.ObjectTypeWarehouse)
		require.NoError(t, err)
		require.Equal(t, "val2", val2)

		alterOptions = &sdk.AlterWarehouseOptions{
			UnsetTag: []sdk.ObjectIdentifier{
				tag.ID(),
				tag2.ID(),
			},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		val, err = client.SystemFunctions.GetTag(ctx, tag.ID(), warehouse.ID(), sdk.ObjectTypeWarehouse)
		require.Error(t, err)
		require.Equal(t, "", val)
		val2, err = client.SystemFunctions.GetTag(ctx, tag2.ID(), warehouse.ID(), sdk.ObjectTypeWarehouse)
		require.Error(t, err)
		require.Equal(t, "", val2)
	})

	t.Run("describe: when warehouse exists", func(t *testing.T) {
		result, err := client.Warehouses.Describe(ctx, precreatedWarehouseId)
		require.NoError(t, err)
		assert.Equal(t, precreatedWarehouseId.Name(), result.Name)
		assert.Equal(t, "WAREHOUSE", result.Kind)
		assert.WithinDuration(t, time.Now(), result.CreatedOn, 1*time.Minute)
	})

	t.Run("describe: when warehouse does not exist", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier
		_, err := client.Warehouses.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: when warehouse exists", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		err := client.Warehouses.Drop(ctx, warehouse.ID(), nil)
		require.NoError(t, err)
		_, err = client.Warehouses.Describe(ctx, warehouse.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: when warehouse does not exist", func(t *testing.T) {
		err := client.Warehouses.Drop(ctx, NonExistingAccountObjectIdentifier, nil)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
