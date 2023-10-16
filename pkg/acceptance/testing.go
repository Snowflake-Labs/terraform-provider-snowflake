package acceptance

import (
	"context"
	"sync"
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

const (
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
	})
}
