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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random/acceptancetests"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/snowflakedb/gosnowflake"
	"github.com/stretchr/testify/require"
)

const AcceptanceTestPrefix = "acc_test_"

var (
	TestDatabaseName  = fmt.Sprintf("%sdb_%s", AcceptanceTestPrefix, acceptancetests.ObjectsSuffix)
	TestSchemaName    = fmt.Sprintf("%ssc_%s", AcceptanceTestPrefix, acceptancetests.ObjectsSuffix)
	TestWarehouseName = fmt.Sprintf("%swh_%s", AcceptanceTestPrefix, acceptancetests.ObjectsSuffix)

	NonExistingAccountObjectIdentifier  = sdk.NewAccountObjectIdentifier("does_not_exist")
	NonExistingDatabaseObjectIdentifier = sdk.NewDatabaseObjectIdentifier(TestDatabaseName, "does_not_exist")
	NonExistingSchemaObjectIdentifier   = sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, "does_not_exist")
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

	defaultConfig, err := sdk.ProfileConfig(testprofiles.Default, true)
	if err != nil {
		log.Panicf("Could not read configuration from profile: %v", err)
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

	cfg, err := sdk.ProfileConfig(testprofiles.Secondary, true)
	if err != nil {
		log.Panicf("Config for the secondary client is needed to run acceptance tests, err: %v", err)
	}
	secondaryClient, err := sdk.NewClient(cfg)
	if err != nil {
		log.Panicf("Cannot instantiate new secondary client, err: %v", err)
	}
	atc.secondaryClient = secondaryClient

	atc.testClient = helpers.NewTestClient(client, TestDatabaseName, TestSchemaName, TestWarehouseName, acceptancetests.ObjectsSuffix)
	atc.secondaryTestClient = helpers.NewTestClient(secondaryClient, TestDatabaseName, TestSchemaName, TestWarehouseName, acceptancetests.ObjectsSuffix)
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

		dbId := TestClient().Ids.DatabaseId()
		schemaId := TestClient().Ids.SchemaId()
		warehouseId := TestClient().Ids.WarehouseId()

		if err := atc.client.Databases.Create(ctx, dbId, &sdk.CreateDatabaseOptions{IfNotExists: sdk.Bool(true)}); err != nil {
			t.Fatal(err)
		}

		if err := atc.client.Schemas.Create(ctx, schemaId, &sdk.CreateSchemaOptions{IfNotExists: sdk.Bool(true)}); err != nil {
			t.Fatal(err)
		}

		if err := atc.client.Warehouses.Create(ctx, warehouseId, &sdk.CreateWarehouseOptions{IfNotExists: sdk.Bool(true)}); err != nil {
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

		if err := helpers.EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(atc.client, ctx); err != nil {
			t.Fatal(err)
		}

		if err := helpers.EnsureQuotedIdentifiersIgnoreCaseIsSetToFalse(atc.secondaryClient, ctx); err != nil {
			t.Fatal(err)
		}

		if err := helpers.EnsureScimProvisionerRolesExist(atc.client, ctx); err != nil {
			t.Fatal(err)
		}

		if err := helpers.EnsureScimProvisionerRolesExist(atc.secondaryClient, ctx); err != nil {
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

// ExternalProviderWithExactVersion returns a map of external providers with an exact version constraint
func ExternalProviderWithExactVersion(version string) map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"snowflake": {
			VersionConstraint: fmt.Sprintf("=%s", version),
			Source:            "Snowflake-Labs/snowflake",
		},
	}
}

// SetV097CompatibleConfigPathEnv sets a new config path in a relevant env variable for a file that is compatible with v0.97.
func SetV097CompatibleConfigPathEnv(t *testing.T) {
	t.Helper()
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	configPath := filepath.Join(home, ".snowflake", "config_v097_compatible")
	t.Setenv(snowflakeenvs.ConfigPath, configPath)
}

// UnsetConfigPathEnv unsets a config path env
func UnsetConfigPathEnv(t *testing.T) {
	t.Helper()
	t.Setenv(snowflakeenvs.ConfigPath, "")
}
