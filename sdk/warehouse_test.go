package sdk

import (
	"context"
	"strings"
)

func (ts *testSuite) createWarehouse() (*Warehouse, error) {
	options := WarehouseCreateOptions{
		Name: "WAREHOUSE_TEST",
		WarehouseProperties: &WarehouseProperties{
			WarehouseType:   String("STANDARD"),
			WarehouseSize:   String("SMALL"),
			MaxClusterCount: Int32(5),
			MinClusterCount: Int32(5),
			ScalingPolicy:   String("STANDARD"),
			AutoSuspend:     Int32(600),
			AutoResume:      Bool(true),
			Comment:         String("test warehouse"),
		},
	}
	return ts.client.Warehouses.Create(context.Background(), options)
}

func (ts *testSuite) TestListWarehouse() {
	warehouse, err := ts.createWarehouse()
	ts.NoError(err)

	limit := 1
	warehouses, err := ts.client.Warehouses.List(context.Background(), WarehouseListOptions{
		Pattern: "WAREHOUSE%",
		Limit:   Int(limit),
	})
	ts.NoError(err)
	ts.Equal(limit, len(warehouses))

	ts.NoError(ts.client.Warehouses.Delete(context.Background(), warehouse.Name))
}

func (ts *testSuite) TestReadWarehouse() {
	warehouse, err := ts.createWarehouse()
	ts.NoError(err)

	entity, err := ts.client.Warehouses.Read(context.Background(), warehouse.Name)
	ts.NoError(err)
	ts.Equal(warehouse.Name, entity.Name)

	ts.NoError(ts.client.Warehouses.Delete(context.Background(), warehouse.Name))
}

func (ts *testSuite) TestCreateWarehouse() {
	warehouse, err := ts.createWarehouse()
	ts.NoError(err)
	ts.NoError(ts.client.Warehouses.Delete(context.Background(), warehouse.Name))
}

func (ts *testSuite) TestUpdateWarehouse() {
	warehouse, err := ts.createWarehouse()
	ts.NoError(err)

	options := WarehouseUpdateOptions{
		WarehouseProperties: &WarehouseProperties{
			WarehouseSize:   String("MEDIUM"),
			MaxClusterCount: Int32(10),
			MinClusterCount: Int32(10),
			ScalingPolicy:   String("ECONOMY"),
			AutoSuspend:     Int32(1000),
			Comment:         String("updated warehouse"),
		},
	}
	afterUpdate, err := ts.client.Warehouses.Update(context.Background(), warehouse.Name, options)
	ts.NoError(err)
	ts.Equal(*options.WarehouseProperties.WarehouseSize, strings.ToUpper(afterUpdate.WarehouseSize))
	ts.Equal(*options.WarehouseProperties.MaxClusterCount, afterUpdate.MaxClusterCount)
	ts.Equal(*options.WarehouseProperties.MinClusterCount, afterUpdate.MinClusterCount)
	ts.Equal(*options.WarehouseProperties.ScalingPolicy, afterUpdate.ScalingPolicy)
	ts.Equal(*options.WarehouseProperties.AutoSuspend, afterUpdate.AutoSuspend)
	ts.Equal(*options.WarehouseProperties.Comment, afterUpdate.Comment)

	ts.NoError(ts.client.Warehouses.Delete(context.Background(), warehouse.Name))
}

func (ts *testSuite) TestRenameWarehouse() {
	warehouse, err := ts.createWarehouse()
	ts.NoError(err)

	newWarehouse := "NEW_WAREHOUSE_TEST"
	ts.NoError(ts.client.Warehouses.Rename(context.Background(), warehouse.Name, newWarehouse))
	ts.NoError(ts.client.Warehouses.Delete(context.Background(), newWarehouse))
}

func (ts *testSuite) TestUseWarehouse() {
	warehouse, err := ts.createWarehouse()
	ts.NoError(err)

	ts.NoError(ts.client.Warehouses.Use(context.Background(), warehouse.Name))
	ts.NoError(ts.client.Warehouses.Delete(context.Background(), warehouse.Name))
}
