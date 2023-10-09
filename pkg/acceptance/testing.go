package acceptance

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

var (
	TestDatabaseName  = "terraform_test_database"
	TestSchemaName    = "terraform_test_schema"
	TestWarehouseName = "terraform_test_warehouse"
)

var TestAccProvider *schema.Provider

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"snowflake": func() (tfprotov6.ProviderServer, error) {
		return tf5to6server.UpgradeServer(
			context.Background(),
			TestAccProvider.GRPCProvider,
		)
	},
}

func init() {
	TestAccProvider = provider.Provider()
}

func TestAccProviders() map[string]*schema.Provider {
	p := provider.Provider()
	return map[string]*schema.Provider{
		"snowflake": p,
	}
}

func TestAccPreCheck(t *testing.T) {
	client, err := sdk.NewDefaultClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	// create test database
	dbId := sdk.NewAccountObjectIdentifier("terraform_test_database")
	_, err = client.Databases.ShowByID(ctx, dbId)
	if err != nil {
		if err := client.Databases.Create(ctx, dbId, &sdk.CreateDatabaseOptions{
			OrReplace: sdk.Bool(true),
		}); err != nil {
			t.Fatal(err)
		}
	}

	// create test schema
	schemaId := sdk.NewDatabaseObjectIdentifier("terraform_test_database", "terraform_test_schema")
	_, err = client.Schemas.ShowByID(ctx, schemaId)
	if err != nil {
		if err := client.Schemas.Create(ctx, schemaId, &sdk.CreateSchemaOptions{
			OrReplace: sdk.Bool(true),
		}); err != nil {
			t.Fatal(err)
		}
	}

	// create test warehouse
	warehouseId := sdk.NewAccountObjectIdentifier("terraform_test_warehouse")
	_, err = client.Warehouses.ShowByID(ctx, warehouseId)
	if err != nil {
		if err := client.Warehouses.Create(ctx, warehouseId, &sdk.CreateWarehouseOptions{
			OrReplace: sdk.Bool(true),
		}); err != nil {
			t.Fatal(err)
		}
	}
}
