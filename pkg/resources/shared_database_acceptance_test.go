package resources_test

import (
	"context"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_CreateSharedDatabase_minimal(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	newId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newComment := random.Comment()

	configVariables := func(id sdk.AccountObjectIdentifier, shareName sdk.ExternalObjectIdentifier, comment string) config.Variables {
		return config.Variables{
			"name":       config.StringVariable(id.Name()),
			"from_share": config.StringVariable(shareName.FullyQualifiedName()),
			"comment":    config.StringVariable(comment),
		}
	}

	shareExternalId := createShareableDatabase(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.SharedDatabase),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables(id, shareExternalId, comment),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SharedDatabase/basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "from_share", shareExternalId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "external_volume", ""),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "catalog", ""),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "replace_invalid_characters", "false"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "default_ddl_collation", ""),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "storage_serialization_policy", "OPTIMIZED"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "log_level", "OFF"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "trace_level", "OFF"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "comment", comment),
				),
			},
			{
				ConfigVariables: configVariables(newId, shareExternalId, newComment),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SharedDatabase/basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "from_share", shareExternalId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "external_volume", ""),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "catalog", ""),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "replace_invalid_characters", "false"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "default_ddl_collation", ""),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "storage_serialization_policy", "OPTIMIZED"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "log_level", "OFF"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "trace_level", "OFF"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "comment", newComment),
				),
			},
			// Import all values
			{
				ConfigVariables:   configVariables(newId, shareExternalId, newComment),
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_SharedDatabase/basic"),
				ResourceName:      "snowflake_shared_database.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_CreateSharedDatabase_complete(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	externalShareId := createShareableDatabase(t)

	externalVolumeId, externalVolumeCleanup := acc.TestClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := acc.TestClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	configVariables := func(
		id sdk.AccountObjectIdentifier,
		shareName sdk.ExternalObjectIdentifier,
		externalVolume sdk.AccountObjectIdentifier,
		catalog sdk.AccountObjectIdentifier,
		replaceInvalidCharacters bool,
		defaultDdlCollation string,
		storageSerializationPolicy sdk.StorageSerializationPolicy,
		logLevel sdk.LogLevel,
		traceLevel sdk.TraceLevel,
		comment string,
	) config.Variables {
		return config.Variables{
			"name":                         config.StringVariable(id.Name()),
			"from_share":                   config.StringVariable(shareName.FullyQualifiedName()),
			"external_volume":              config.StringVariable(externalVolume.Name()),
			"catalog":                      config.StringVariable(catalog.Name()),
			"replace_invalid_characters":   config.BoolVariable(replaceInvalidCharacters),
			"default_ddl_collation":        config.StringVariable(defaultDdlCollation),
			"storage_serialization_policy": config.StringVariable(string(storageSerializationPolicy)),
			"log_level":                    config.StringVariable(string(logLevel)),
			"trace_level":                  config.StringVariable(string(traceLevel)),
			"comment":                      config.StringVariable(comment),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.SharedDatabase),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables(
					id,
					externalShareId,
					externalVolumeId,
					catalogId,
					true,
					"en_US",
					sdk.StorageSerializationPolicyOptimized,
					sdk.LogLevelInfo,
					sdk.TraceLevelOnEvent,
					comment,
				),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SharedDatabase/complete"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "replace_invalid_characters", "true"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "storage_serialization_policy", string(sdk.StorageSerializationPolicyOptimized)),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "trace_level", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "comment", comment),
				),
			},
			// Import all values
			{
				ConfigVariables: configVariables(
					id,
					externalShareId,
					externalVolumeId,
					catalogId,
					true,
					"en_US",
					sdk.StorageSerializationPolicyOptimized,
					sdk.LogLevelInfo,
					sdk.TraceLevelOnEvent,
					comment,
				),
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_SharedDatabase/complete"),
				ResourceName:      "snowflake_shared_database.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_CreateSharedDatabase_InvalidValues(t *testing.T) {
	comment := random.Comment()

	configVariables := func(
		replaceInvalidCharacters bool,
		defaultDdlCollation string,
		storageSerializationPolicy string,
		logLevel string,
		traceLevel string,
		comment string,
	) config.Variables {
		return config.Variables{
			"name":                         config.StringVariable(""),
			"from_share":                   config.StringVariable(""),
			"external_volume":              config.StringVariable(""),
			"catalog":                      config.StringVariable(""),
			"replace_invalid_characters":   config.BoolVariable(replaceInvalidCharacters),
			"default_ddl_collation":        config.StringVariable(defaultDdlCollation),
			"storage_serialization_policy": config.StringVariable(storageSerializationPolicy),
			"log_level":                    config.StringVariable(logLevel),
			"trace_level":                  config.StringVariable(traceLevel),
			"comment":                      config.StringVariable(comment),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.SharedDatabase),
		Steps: []resource.TestStep{
			{
				ConfigVariables: configVariables(
					true,
					"en_US",
					"invalid_value",
					"invalid_value",
					"invalid_value",
					comment,
				),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SharedDatabase/complete"),
				ExpectError: regexp.MustCompile(`(expected \[{{} log_level}\] to be one of \[\"TRACE\" \"DEBUG\" \"INFO\" \"WARN\" \"ERROR\" \"FATAL\" \"OFF\"\], got invalid_value)|` +
					`(expected \[{{} trace_level}\] to be one of \[\"ALWAYS\" \"ON_EVENT\" \"OFF\"\], got invalid_value)|` +
					`(expected \[{{} storage_serialization_policy}\] to be one of \[\"COMPATIBLE\" \"OPTIMIZED\"\], got invalid_value)`),
			},
		},
	})
}

// createShareableDatabase creates a database on the secondary account and enables database sharing on the primary account.
// TODO(SNOW-1431726): Later on, this function should be moved to more sophisticated helpers.
func createShareableDatabase(t *testing.T) sdk.ExternalObjectIdentifier {
	t.Helper()

	ctx := context.Background()

	share, shareCleanup := acc.SecondaryTestClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	sharedDatabase, sharedDatabaseCleanup := acc.SecondaryTestClient().Database.CreateDatabase(t)
	t.Cleanup(sharedDatabaseCleanup)

	err := acc.SecondaryClient(t).Grants.GrantPrivilegeToShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
		Database: sharedDatabase.ID(),
	}, share.ID())
	require.NoError(t, err)
	t.Cleanup(func() {
		err := acc.SecondaryClient(t).Grants.RevokePrivilegeFromShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
			Database: sharedDatabase.ID(),
		}, share.ID())
		require.NoError(t, err)
	})

	err = acc.SecondaryClient(t).Shares.Alter(ctx, share.ID(), &sdk.AlterShareOptions{
		IfExists: sdk.Bool(true),
		Set: &sdk.ShareSet{
			Accounts: []sdk.AccountIdentifier{
				acc.TestClient().Account.GetAccountIdentifier(t),
			},
		},
	})
	require.NoError(t, err)

	return sdk.NewExternalObjectIdentifier(acc.SecondaryTestClient().Account.GetAccountIdentifier(t), share.ID())
}
