// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package testint

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/sdk/internal/random"
)

var itc integrationTestContext

func TestMain(m *testing.M) {
	exitVal := execute(m)
	os.Exit(exitVal)
}

func execute(m *testing.M) int {
	defer timer("tests")()
	setup()
	exitVal := m.Run()
	cleanup()
	return exitVal
}

func setup() {
	log.Println("Running integration tests setup")

	err := itc.initialize()
	if err != nil {
		log.Printf("Integration test context initialisation failed with %s\n", err)
		cleanup()
		os.Exit(1)
	}
}

func cleanup() {
	log.Println("Running integration tests cleanup")
	if itc.databaseCleanup != nil {
		defer itc.databaseCleanup()
	}
	if itc.schemaCleanup != nil {
		defer itc.schemaCleanup()
	}
}

type integrationTestContext struct {
	client *sdk.Client
	ctx    context.Context

	database         *sdk.Database
	databaseCleanup  func()
	schema           *sdk.Schema
	schemaCleanup    func()
	warehouse        *sdk.Warehouse
	warehouseCleanup func()
}

func (itc *integrationTestContext) initialize() error {
	log.Println("Initializing integration test context")
	var err error
	c, err := sdk.NewDefaultClient()
	if err != nil {
		return err
	}
	itc.client = c
	itc.ctx = context.Background()

	db, dbCleanup, err := createDb(itc.client, itc.ctx)
	if err != nil {
		return err
	}
	itc.database = db
	itc.databaseCleanup = dbCleanup

	sc, scCleanup, err := createSc(itc.client, itc.ctx, itc.database)
	if err != nil {
		return err
	}
	itc.schema = sc
	itc.schemaCleanup = scCleanup

	wh, whCleanup, err := createWh(itc.client, itc.ctx)
	if err != nil {
		return err
	}
	itc.warehouse = wh
	itc.warehouseCleanup = whCleanup

	return nil
}

func createDb(client *sdk.Client, ctx context.Context) (*sdk.Database, func(), error) {
	name := "int_test_db_" + random.UUID()
	id := sdk.NewAccountObjectIdentifier(name)
	err := client.Databases.Create(ctx, id, nil)
	if err != nil {
		return nil, nil, err
	}
	database, err := client.Databases.ShowByID(ctx, id)
	return database, func() {
		_ = client.Databases.Drop(ctx, id, nil)
	}, err
}

func createSc(client *sdk.Client, ctx context.Context, db *sdk.Database) (*sdk.Schema, func(), error) {
	name := "int_test_sc_" + random.UUID()
	id := sdk.NewDatabaseObjectIdentifier(db.Name, name)
	err := client.Schemas.Create(ctx, id, nil)
	if err != nil {
		return nil, nil, err
	}
	schema, err := client.Schemas.ShowByID(ctx, sdk.NewDatabaseObjectIdentifier(db.Name, name))
	return schema, func() {
		_ = client.Schemas.Drop(ctx, id, nil)
	}, err
}

func createWh(client *sdk.Client, ctx context.Context) (*sdk.Warehouse, func(), error) {
	name := "int_test_wh_" + random.UUID()
	id := sdk.NewAccountObjectIdentifier(name)
	err := client.Warehouses.Create(ctx, id, nil)
	if err != nil {
		return nil, nil, err
	}
	warehouse, err := client.Warehouses.ShowByID(ctx, id)
	return warehouse, func() {
		_ = client.Warehouses.Drop(ctx, id, nil)
	}, err
}

// timer measures time from invocation point to the end of method.
// It's supposed to be used like:
//
//	defer timer("something to measure name")()
func timer(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("[DEBUG] %s took %v\n", name, time.Since(start))
	}
}

func testClient(t *testing.T) *sdk.Client {
	t.Helper()
	return itc.client
}

func testContext(t *testing.T) context.Context {
	t.Helper()
	return itc.ctx
}

func testDb(t *testing.T) *sdk.Database {
	t.Helper()
	return itc.database
}

func testSchema(t *testing.T) *sdk.Schema {
	t.Helper()
	return itc.schema
}

func testWarehouse(t *testing.T) *sdk.Warehouse {
	t.Helper()
	return itc.warehouse
}
