package testint

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/snowflakedb/gosnowflake"
)

const IntegrationTestPrefix = "int_test_"

var (
	TestWarehouseName = fmt.Sprintf("%swh_%s", IntegrationTestPrefix, random.IntegrationTestsSuffix)
	TestDatabaseName  = fmt.Sprintf("%sdb_%s", IntegrationTestPrefix, random.IntegrationTestsSuffix)
	TestSchemaName    = fmt.Sprintf("%ssc_%s", IntegrationTestPrefix, random.IntegrationTestsSuffix)

	NonExistingAccountObjectIdentifier  = sdk.NewAccountObjectIdentifier("does_not_exist")
	NonExistingDatabaseObjectIdentifier = sdk.NewDatabaseObjectIdentifier(TestDatabaseName, "does_not_exist")
	NonExistingSchemaObjectIdentifier   = sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, "does_not_exist")
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
		log.Printf("Integration test context initialisation failed with: `%s`\n", err)
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
	if itc.warehouseCleanup != nil {
		defer itc.warehouseCleanup()
	}
	if itc.secondaryDatabaseCleanup != nil {
		defer itc.secondaryDatabaseCleanup()
	}
	if itc.secondarySchemaCleanup != nil {
		defer itc.secondarySchemaCleanup()
	}
	if itc.secondaryWarehouseCleanup != nil {
		defer itc.secondaryWarehouseCleanup()
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

	testObjectSuffix := os.Getenv(fmt.Sprintf("%v", testenvs.TestObjectsSuffix))
	requireTestObjectSuffix := os.Getenv(fmt.Sprintf("%v", testenvs.RequireTestObjectsSuffix))
	if requireTestObjectSuffix != "" && testObjectSuffix == "" {
		return errors.New("test object suffix is required for this test run")
	}

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

	db, dbCleanup, err := createDb(itc.client, itc.ctx, false)
	itc.databaseCleanup = dbCleanup
	if err != nil {
		return err
	}
	itc.database = db

	sc, scCleanup, err := createSc(itc.client, itc.ctx, itc.database, false)
	itc.schemaCleanup = scCleanup
	if err != nil {
		return err
	}
	itc.schema = sc

	wh, whCleanup, err := createWh(itc.client, itc.ctx, false)
	itc.warehouseCleanup = whCleanup
	if err != nil {
		return err
	}
	itc.warehouse = wh

	config, err := sdk.ProfileConfig(testprofiles.Secondary)
	if err != nil {
		return err
	}

	if config.Account == defaultConfig.Account {
		log.Println("[WARN] default and secondary configs are set to the same account; it may cause problems in tests requiring multiple accounts")
	}

	secondaryClient, err := sdk.NewClient(config)
	if err != nil {
		return err
	}
	itc.secondaryClient = secondaryClient
	itc.secondaryCtx = context.Background()

	secondaryDb, secondaryDbCleanup, err := createDb(itc.secondaryClient, itc.secondaryCtx, true)
	itc.secondaryDatabaseCleanup = secondaryDbCleanup
	if err != nil {
		return err
	}
	itc.secondaryDatabase = secondaryDb

	secondarySchema, secondarySchemaCleanup, err := createSc(itc.secondaryClient, itc.secondaryCtx, itc.database, true)
	itc.secondarySchemaCleanup = secondarySchemaCleanup
	if err != nil {
		return err
	}
	itc.secondarySchema = secondarySchema

	secondaryWarehouse, secondaryWarehouseCleanup, err := createWh(itc.secondaryClient, itc.secondaryCtx, true)
	itc.secondaryWarehouseCleanup = secondaryWarehouseCleanup
	if err != nil {
		return err
	}
	itc.secondaryWarehouse = secondaryWarehouse

	itc.testClient = helpers.NewTestClient(c, TestDatabaseName, TestSchemaName, TestWarehouseName, random.IntegrationTestsSuffix)
	itc.secondaryTestClient = helpers.NewTestClient(secondaryClient, TestDatabaseName, TestSchemaName, TestWarehouseName, random.IntegrationTestsSuffix)

	err = helpers.EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(itc.client, itc.ctx)
	if err != nil {
		return err
	}
	err = helpers.EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(itc.secondaryClient, itc.secondaryCtx)
	if err != nil {
		return err
	}

	err = helpers.EnsureScimProvisionerRolesExist(itc.client, itc.ctx)
	if err != nil {
		return err
	}
	err = helpers.EnsureScimProvisionerRolesExist(itc.secondaryClient, itc.secondaryCtx)
	if err != nil {
		return err
	}

	return nil
}

func createDb(client *sdk.Client, ctx context.Context, ifNotExists bool) (*sdk.Database, func(), error) {
	id := sdk.NewAccountObjectIdentifier(TestDatabaseName)
	cleanup := func() {
		_ = client.Databases.Drop(ctx, id, &sdk.DropDatabaseOptions{IfExists: sdk.Bool(true)})
	}
	err := client.Databases.Create(ctx, id, &sdk.CreateDatabaseOptions{IfNotExists: sdk.Bool(ifNotExists)})
	if err != nil {
		return nil, cleanup, err
	}
	database, err := client.Databases.ShowByID(ctx, id)
	return database, cleanup, err
}

func createSc(client *sdk.Client, ctx context.Context, db *sdk.Database, ifNotExists bool) (*sdk.Schema, func(), error) {
	id := sdk.NewDatabaseObjectIdentifier(db.Name, TestSchemaName)
	cleanup := func() {
		_ = client.Schemas.Drop(ctx, id, &sdk.DropSchemaOptions{IfExists: sdk.Bool(true)})
	}
	err := client.Schemas.Create(ctx, id, &sdk.CreateSchemaOptions{IfNotExists: sdk.Bool(ifNotExists)})
	if err != nil {
		return nil, cleanup, err
	}
	schema, err := client.Schemas.ShowByID(ctx, sdk.NewDatabaseObjectIdentifier(db.Name, TestSchemaName))
	return schema, cleanup, err
}

func createWh(client *sdk.Client, ctx context.Context, ifNotExists bool) (*sdk.Warehouse, func(), error) {
	id := sdk.NewAccountObjectIdentifier(TestWarehouseName)
	cleanup := func() {
		_ = client.Warehouses.Drop(ctx, id, &sdk.DropWarehouseOptions{IfExists: sdk.Bool(true)})
	}
	err := client.Warehouses.Create(ctx, id, &sdk.CreateWarehouseOptions{IfNotExists: sdk.Bool(ifNotExists)})
	if err != nil {
		return nil, cleanup, err
	}
	warehouse, err := client.Warehouses.ShowByID(ctx, id)
	return warehouse, cleanup, err
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

// TODO [SNOW-1653619]: distinguish between testClient and testClientHelper in tests (one is tested, the second helps with the tests, they should be the different ones)
func testClient(t *testing.T) *sdk.Client {
	t.Helper()
	return itc.client
}

func testContext(t *testing.T) context.Context {
	t.Helper()
	return itc.ctx
}

func testSecondaryClient(t *testing.T) *sdk.Client {
	t.Helper()
	return itc.secondaryClient
}

func testSecondaryContext(t *testing.T) context.Context {
	t.Helper()
	return itc.secondaryCtx
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
