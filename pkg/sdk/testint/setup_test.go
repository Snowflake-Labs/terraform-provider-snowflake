package testint

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/snowflakedb/gosnowflake"
)

var (
	TestWarehouseName = "int_test_wh_" + random.UUID()
	TestDatabaseName  = "int_test_db_" + random.UUID()
	TestSchemaName    = "int_test_sc_" + random.UUID()
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
	if itc.secondaryDatabaseCleanup != nil {
		defer itc.secondaryDatabaseCleanup()
	}
	if itc.secondarySchemaCleanup != nil {
		defer itc.secondarySchemaCleanup()
	}
}

type integrationTestContext struct {
	config *gosnowflake.Config
	client *sdk.Client
	ctx    context.Context

	database         *sdk.Database
	databaseCleanup  func()
	schema           *sdk.Schema
	schemaCleanup    func()
	warehouse        *sdk.Warehouse
	warehouseCleanup func()

	secondaryClient *sdk.Client
	secondaryCtx    context.Context

	secondaryDatabase         *sdk.Database
	secondaryDatabaseCleanup  func()
	secondarySchema           *sdk.Schema
	secondarySchemaCleanup    func()
	secondaryWarehouse        *sdk.Warehouse
	secondaryWarehouseCleanup func()

	testClient          *helpers.TestClient
	secondaryTestClient *helpers.TestClient
}

func (itc *integrationTestContext) initialize() error {
	log.Println("Initializing integration test context")

	defaultConfig, err := sdk.ProfileConfig(testprofiles.Default)
	if err != nil {
		return err
	}
	if defaultConfig == nil {
		return errors.New("config is required to run integration tests")
	}
	itc.config = defaultConfig

	c, err := sdk.NewClient(defaultConfig)
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

	config, err := sdk.ProfileConfig(testprofiles.Secondary)
	if err != nil {
		return err
	}

	secondaryClient, err := sdk.NewClient(config)
	if err != nil {
		return err
	}
	itc.secondaryClient = secondaryClient
	itc.secondaryCtx = context.Background()

	secondaryDb, secondaryDbCleanup, err := createDb(itc.secondaryClient, itc.secondaryCtx)
	if err != nil {
		return err
	}
	itc.secondaryDatabase = secondaryDb
	itc.secondaryDatabaseCleanup = secondaryDbCleanup

	secondarySchema, secondarySchemaCleanup, err := createSc(itc.secondaryClient, itc.secondaryCtx, itc.database)
	if err != nil {
		return err
	}
	itc.secondarySchema = secondarySchema
	itc.secondarySchemaCleanup = secondarySchemaCleanup

	secondaryWarehouse, secondaryWarehouseCleanup, err := createWh(itc.secondaryClient, itc.secondaryCtx)
	if err != nil {
		return err
	}
	itc.secondaryWarehouse = secondaryWarehouse
	itc.secondaryWarehouseCleanup = secondaryWarehouseCleanup

	itc.testClient = helpers.NewTestClient(c, TestDatabaseName, TestSchemaName, TestWarehouseName)
	itc.secondaryTestClient = helpers.NewTestClient(secondaryClient, TestDatabaseName, TestSchemaName, TestWarehouseName)

	return nil
}

func createDb(client *sdk.Client, ctx context.Context) (*sdk.Database, func(), error) {
	id := sdk.NewAccountObjectIdentifier(TestDatabaseName)
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
	id := sdk.NewDatabaseObjectIdentifier(db.Name, TestSchemaName)
	err := client.Schemas.Create(ctx, id, nil)
	if err != nil {
		return nil, nil, err
	}
	schema, err := client.Schemas.ShowByID(ctx, sdk.NewDatabaseObjectIdentifier(db.Name, TestSchemaName))
	return schema, func() {
		_ = client.Schemas.Drop(ctx, id, nil)
	}, err
}

func createWh(client *sdk.Client, ctx context.Context) (*sdk.Warehouse, func(), error) {
	id := sdk.NewAccountObjectIdentifier(TestWarehouseName)
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

func testSecondaryClient(t *testing.T) *sdk.Client {
	t.Helper()
	return itc.secondaryClient
}

func testSecondaryContext(t *testing.T) context.Context {
	t.Helper()
	return itc.secondaryCtx
}

func testSecondaryDb(t *testing.T) *sdk.Database {
	t.Helper()
	return itc.secondaryDatabase
}

func testSecondarySchema(t *testing.T) *sdk.Schema {
	t.Helper()
	return itc.secondarySchema
}

func testSecondaryWarehouse(t *testing.T) *sdk.Warehouse {
	t.Helper()
	return itc.secondaryWarehouse
}

func testConfig(t *testing.T) *gosnowflake.Config {
	t.Helper()
	return itc.config
}

func testClientHelper() *helpers.TestClient {
	return itc.testClient
}

func secondaryTestClientHelper() *helpers.TestClient {
	return itc.secondaryTestClient
}
