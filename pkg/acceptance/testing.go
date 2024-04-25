package acceptance

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/snowflakedb/gosnowflake"
)

var (
	TestDatabaseName   = "acc_test_db_" + random.AcceptanceTestsSuffix
	TestSchemaName     = "acc_test_sc_" + random.AcceptanceTestsSuffix
	TestWarehouseName  = "acc_test_wh_" + random.AcceptanceTestsSuffix
	TestWarehouseName2 = "acc_test_wh2_" + random.AcceptanceTestsSuffix
)

var (
	TestAccProvider *schema.Provider
	v5Server        tfprotov5.ProviderServer
	v6Server        tfprotov6.ProviderServer
	atc             acceptanceTestContext
)

func init() {
	testObjectSuffix := os.Getenv(fmt.Sprintf("%v", testenvs.TestObjectsSuffix))
	requireTestObjectSuffix := os.Getenv(fmt.Sprintf("%v", testenvs.RequireTestObjectsSuffix))
	if requireTestObjectSuffix != "" && testObjectSuffix == "" {
		log.Println("test object suffix is required for this test run")
		os.Exit(1)
	}

	TestAccProvider = provider.Provider()
	v5Server = TestAccProvider.GRPCProvider()
	var err error
	v6Server, err = tf5to6server.UpgradeServer(
		context.Background(),
		func() tfprotov5.ProviderServer {
			return v5Server
		},
	)
	if err != nil {
		log.Panicf("Cannot upgrade server from proto v5 to proto v6, failing, err: %v", err)
	}
	_ = testAccProtoV6ProviderFactoriesNew

	defaultConfig, err := sdk.ProfileConfig(testprofiles.Default)
	if err != nil {
		log.Panicf("Cannot load default config, err: %v", err)
	}
	if defaultConfig == nil {
		log.Panic("Config is required to run acceptance tests")
	}
	atc.config = defaultConfig

	client, err := sdk.NewClient(defaultConfig)
	if err != nil {
		log.Panicf("Cannot instantiate new client, err: %v", err)
	}
	atc.client = client

	cfg, err := sdk.ProfileConfig(testprofiles.Secondary)
	if err != nil {
		log.Panicf("Config for the secondary client is needed to run acceptance tests, err: %v", err)
	}
	secondaryClient, err := sdk.NewClient(cfg)
	if err != nil {
		log.Panicf("Cannot instantiate new secondary client, err: %v", err)
	}
	atc.secondaryClient = secondaryClient

	atc.testClient = helpers.NewTestClient(client, TestDatabaseName, TestSchemaName, TestWarehouseName, random.AcceptanceTestsSuffix)
	atc.secondaryTestClient = helpers.NewTestClient(secondaryClient, TestDatabaseName, TestSchemaName, TestWarehouseName, random.AcceptanceTestsSuffix)
}

type acceptanceTestContext struct {
	config              *gosnowflake.Config
	client              *sdk.Client
	secondaryClient     *sdk.Client
	testClient          *helpers.TestClient
	secondaryTestClient *helpers.TestClient
}

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"snowflake": func() (tfprotov6.ProviderServer, error) {
		return v6Server, nil
	},
}

// if we do not reuse the created objects there is no `Previously configured provider being re-configured.` warning
// currently left for possible usage after other improvements
var testAccProtoV6ProviderFactoriesNew = map[string]func() (tfprotov6.ProviderServer, error){
	"snowflake": func() (tfprotov6.ProviderServer, error) {
		return tf5to6server.UpgradeServer(
			context.Background(),
			provider.Provider().GRPCProvider,
		)
	},
}

var once sync.Once

func TestAccPreCheck(t *testing.T) {
	// use singleton design pattern to ensure we only create these resources once
	// there is no cleanup currently, sweepers take care of it
	once.Do(func() {
		ctx := context.Background()

		dbId := sdk.NewAccountObjectIdentifier(TestDatabaseName)
		schemaId := sdk.NewDatabaseObjectIdentifier(TestDatabaseName, TestSchemaName)
		warehouseId := sdk.NewAccountObjectIdentifier(TestWarehouseName)
		warehouseId2 := sdk.NewAccountObjectIdentifier(TestWarehouseName2)

		if err := atc.client.Databases.Create(ctx, dbId, &sdk.CreateDatabaseOptions{IfNotExists: sdk.Bool(true)}); err != nil {
			t.Fatal(err)
		}

		if err := atc.client.Schemas.Create(ctx, schemaId, &sdk.CreateSchemaOptions{IfNotExists: sdk.Bool(true)}); err != nil {
			t.Fatal(err)
		}

		if err := atc.client.Warehouses.Create(ctx, warehouseId, &sdk.CreateWarehouseOptions{IfNotExists: sdk.Bool(true)}); err != nil {
			t.Fatal(err)
		}

		if err := atc.client.Warehouses.Create(ctx, warehouseId2, &sdk.CreateWarehouseOptions{IfNotExists: sdk.Bool(true)}); err != nil {
			t.Fatal(err)
		}

		if err := atc.secondaryClient.Databases.Create(ctx, dbId, &sdk.CreateDatabaseOptions{IfNotExists: sdk.Bool(true)}); err != nil {
			t.Fatal(err)
		}

		if err := atc.secondaryClient.Schemas.Create(ctx, schemaId, &sdk.CreateSchemaOptions{IfNotExists: sdk.Bool(true)}); err != nil {
			t.Fatal(err)
		}

		if err := atc.secondaryClient.Warehouses.Create(ctx, warehouseId, &sdk.CreateWarehouseOptions{IfNotExists: sdk.Bool(true)}); err != nil {
			t.Fatal(err)
		}

		if err := atc.secondaryClient.Warehouses.Create(ctx, warehouseId2, &sdk.CreateWarehouseOptions{IfNotExists: sdk.Bool(true)}); err != nil {
			t.Fatal(err)
		}
	})
}

// ConfigurationSameAsStepN should be used to obtain configuration for one of the previous steps to avoid duplication of configuration and var files.
// Based on config.TestStepDirectory.
func ConfigurationSameAsStepN(step int) func(config.TestStepConfigRequest) string {
	return func(req config.TestStepConfigRequest) string {
		return filepath.Join("testdata", req.TestName, strconv.Itoa(step))
	}
}

// ConfigurationDirectory should be used to obtain configuration if the same can be shared between multiple tests to avoid duplication of configuration and var files.
// Based on config.TestNameDirectory. Similar to config.StaticDirectory but prefixed provided directory with `testdata`.
func ConfigurationDirectory(directory string) func(config.TestStepConfigRequest) string {
	return func(req config.TestStepConfigRequest) string {
		return filepath.Join("testdata", directory)
	}
}

func Client(t *testing.T) *sdk.Client {
	t.Helper()
	return atc.client
}

func SecondaryClient(t *testing.T) *sdk.Client {
	t.Helper()
	return atc.secondaryClient
}

func DefaultConfig(t *testing.T) *gosnowflake.Config {
	t.Helper()
	return atc.config
}

func TestClient() *helpers.TestClient {
	return atc.testClient
}

func SecondaryTestClient() *helpers.TestClient {
	return atc.secondaryTestClient
}
