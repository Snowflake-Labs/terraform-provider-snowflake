package acceptance

import (
	"context"
	"log"
	"path/filepath"
	"strconv"
	"sync"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/config"
)

const (
	TestDatabaseName   = "terraform_test_database"
	TestSchemaName     = "terraform_test_schema"
	TestWarehouseName  = "terraform_test_warehouse"
	TestWarehouseName2 = "terraform_test_warehouse_2"
)

var (
	TestAccProvider *schema.Provider
	v5Server        tfprotov5.ProviderServer
	v6Server        tfprotov6.ProviderServer
)

func init() {
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
	once.Do(func() {
		client, err := sdk.NewDefaultClient()
		if err != nil {
			t.Fatal(err)
		}
		ctx := context.Background()

		dbId := sdk.NewAccountObjectIdentifier(TestDatabaseName)
		if err := client.Databases.Create(ctx, dbId, &sdk.CreateDatabaseOptions{
			IfNotExists: sdk.Bool(true),
		}); err != nil {
			t.Fatal(err)
		}

		schemaId := sdk.NewDatabaseObjectIdentifier(TestDatabaseName, TestSchemaName)
		if err := client.Schemas.Create(ctx, schemaId, &sdk.CreateSchemaOptions{
			IfNotExists: sdk.Bool(true),
		}); err != nil {
			t.Fatal(err)
		}

		warehouseId := sdk.NewAccountObjectIdentifier(TestWarehouseName)
		if err := client.Warehouses.Create(ctx, warehouseId, &sdk.CreateWarehouseOptions{
			IfNotExists: sdk.Bool(true),
		}); err != nil {
			t.Fatal(err)
		}

		warehouseId2 := sdk.NewAccountObjectIdentifier(TestWarehouseName2)
		if err := client.Warehouses.Create(ctx, warehouseId2, &sdk.CreateWarehouseOptions{
			IfNotExists: sdk.Bool(true),
		}); err != nil {
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
