package resources_test

import (
	"context"
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
	name := acc.TestClient().Ids.Alpha()
	comment := random.Comment()

	newName := acc.TestClient().Ids.Alpha()
	newComment := random.Comment()

	configVariables := func(name string, shareName sdk.ExternalObjectIdentifier, comment string) config.Variables {
		return config.Variables{
			"name":       config.StringVariable(name),
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
				ConfigVariables: configVariables(name, shareExternalId, comment),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SharedDatabase/basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "from_share", shareExternalId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "external_volume", ""),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "catalog", ""),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "default_ddl_collation", ""),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "log_level", "OFF"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "trace_level", "OFF"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "comment", comment),
				),
			},
			{
				ConfigVariables: configVariables(newName, shareExternalId, newComment),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SharedDatabase/basic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "name", newName),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "from_share", shareExternalId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "external_volume", ""),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "catalog", ""),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "default_ddl_collation", ""),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "log_level", "OFF"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "trace_level", "OFF"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "comment", newComment),
				),
			},
		},
	})
}

func TestAcc_CreateSharedDatabase_complete(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()
	comment := random.Comment()
	externalShareId := createShareableDatabase(t)

	externalVolumeId, externalVolumeCleanup := acc.TestClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	catalogId, catalogCleanup := acc.TestClient().CatalogIntegration.Create(t)
	t.Cleanup(catalogCleanup)

	configVariables := func(
		name string,
		shareName sdk.ExternalObjectIdentifier,
		transient bool,
		externalVolume string,
		catalog string,
		defaultDdlCollation string,
		logLevel string,
		traceLevel string,
		comment string,
	) config.Variables {
		return config.Variables{
			"name":                  config.StringVariable(name),
			"from_share":            config.StringVariable(shareName.FullyQualifiedName()),
			"transient":             config.BoolVariable(transient),
			"external_volume":       config.StringVariable(externalVolume),
			"catalog":               config.StringVariable(catalog),
			"default_ddl_collation": config.StringVariable(defaultDdlCollation),
			"log_level":             config.StringVariable(logLevel),
			"trace_level":           config.StringVariable(traceLevel),
			"comment":               config.StringVariable(comment),
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
				ConfigVariables: configVariables(name, externalShareId, true, externalVolumeId.Name(), catalogId.Name(), "en_US", string(sdk.LogLevelInfo), string(sdk.TraceLevelOnEvent), comment),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_SharedDatabase/complete"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "external_volume", externalVolumeId.Name()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "catalog", catalogId.Name()),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "default_ddl_collation", "en_US"),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "log_level", string(sdk.LogLevelInfo)),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "trace_level", string(sdk.TraceLevelOnEvent)),
					resource.TestCheckResourceAttr("snowflake_shared_database.test", "comment", comment),
				),
			},
			// Import all values
			{
				ConfigVariables:   configVariables(name, externalShareId, true, externalVolumeId.Name(), catalogId.Name(), "en_US", string(sdk.LogLevelInfo), string(sdk.TraceLevelOnEvent), comment),
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_SharedDatabase/complete"),
				ResourceName:      "snowflake_shared_database.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

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

	return share.ExternalID()
}
