package datasources_test

import (
	"context"
	"regexp"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func checkAtLeastOneGrantPresent() resource.TestCheckFunc {
	datasourceName := "data.snowflake_grants.test"
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(datasourceName, "grants.#"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.created_on"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.privilege"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.granted_on"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.name"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.granted_to"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.grantee_name"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.grant_option"),
	)
}

func checkAtLeastOneFutureGrantPresent() resource.TestCheckFunc {
	datasourceName := "data.snowflake_grants.test"
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(datasourceName, "grants.#"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.created_on"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.privilege"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.name"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.grantee_name"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.grant_option"),
	)
}

func checkAtLeastOneGrantPresentWithoutValidations() resource.TestCheckFunc {
	datasourceName := "data.snowflake_grants.test"
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(datasourceName, "grants.#"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.created_on"),
	)
}

func checkNoGrantsPresent() resource.TestCheckFunc {
	datasourceName := "data.snowflake_grants.test"
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(datasourceName, "grants.#"),
		resource.TestCheckNoResourceAttr(datasourceName, "grants.0.created_on"),
	)
}

func getCurrentUser(t *testing.T) string {
	t.Helper()
	client, err := sdk.NewDefaultClient()
	if err != nil {
		t.Fatal(err)
	}
	user, err := client.ContextFunctions.CurrentUser(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func getSecondaryAccountIdentifier(t *testing.T) *sdk.AccountIdentifier {
	t.Helper()

	t.Helper()
	client, err := sdk.NewDefaultClient()
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := sdk.ProfileConfig(testprofiles.Secondary)
	if err != nil {
		t.Fatal(err)
	}
	secondaryClient, err := sdk.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	replicationAccounts, err := client.ReplicationFunctions.ShowReplicationAccounts(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, replicationAccount := range replicationAccounts {
		if replicationAccount.AccountLocator == secondaryClient.GetAccountLocator() {
			return sdk.Pointer(sdk.NewAccountIdentifier(replicationAccount.OrganizationName, replicationAccount.AccountName))
		}
	}
	return nil
}

// TODO: tests (examples from the correct ones):
// + on - account
// + on - account object
// + on - db object
// + on - schema object
// + on - invalid config - no attribute
// + on - invalid config - missing object type or name
// - to - application
// - to - application role
// + to - role
// + to - user
// + to - share
// +/- to - share with application package
// + to - invalid config - no attribute
// + to - invalid config - share name missing
// + of - role
// - of - application role
// + of - share
// + of - invalid config - no attribute
// + future in - database
// + future in - schema (both db and sc present)
// + future in - invalid config - no attribute
// + future in - invalid config - schema name not fully qualified
// + future to - role
// + future to - database role
// + future to - invalid config - no attribute
// + future to - invalid config - database role id invalid
func TestAcc_Grants_On_Account(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/On/Account"),
				Check:           checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_On_AccountObject(t *testing.T) {
	configVariables := config.Variables{
		"database": config.StringVariable(acc.TestDatabaseName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/On/AccountObject"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_On_DatabaseObject(t *testing.T) {
	configVariables := config.Variables{
		"database": config.StringVariable(acc.TestDatabaseName),
		"schema":   config.StringVariable(acc.TestSchemaName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/On/DatabaseObject"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_On_SchemaObject(t *testing.T) {
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := config.Variables{
		"database": config.StringVariable(acc.TestDatabaseName),
		"schema":   config.StringVariable(acc.TestSchemaName),
		"table":    config.StringVariable(tableName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/On/SchemaObject"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_On_Invalid_NoAttribute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/On/Invalid/NoAttribute"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Grants_On_Invalid_MissingObjectType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/On/Invalid/MissingObjectType"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Missing required argument"),
			},
		},
	})
}

func TestAcc_Grants_To_Role(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/Role"),
				Check:           checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_To_User(t *testing.T) {
	user := getCurrentUser(t)
	configVariables := config.Variables{
		"user": config.StringVariable(user),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/User"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneGrantPresentWithoutValidations(),
			},
		},
	})
}

func TestAcc_Grants_To_Share(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	shareName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := config.Variables{
		"database": config.StringVariable(databaseName),
		"share":    config.StringVariable(shareName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/Share"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_To_ShareWithApplicationPackage(t *testing.T) {
	t.Skip("No SDK support yet")
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	shareName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := config.Variables{
		"database": config.StringVariable(databaseName),
		"share":    config.StringVariable(shareName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/ShareWithApplicationPackage"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_To_Invalid_NoAttribute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/Invalid/NoAttribute"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Grants_To_Invalid_ShareNameMissing(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/Invalid/ShareNameMissing"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Missing required argument"),
			},
		},
	})
}

func TestAcc_Grants_Of_Role(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/Of/Role"),
				Check:           checkAtLeastOneGrantPresentWithoutValidations(),
			},
		},
	})
}

func TestAcc_Grants_Of_Share(t *testing.T) {
	t.Skip("TestAcc_Share are skipped")
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	shareName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountId := getSecondaryAccountIdentifier(t)
	require.NotNil(t, accountId)

	configVariables := config.Variables{
		"database": config.StringVariable(databaseName),
		"share":    config.StringVariable(shareName),
		"account":  config.StringVariable(accountId.FullyQualifiedName()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/Of/Share"),
				ConfigVariables: configVariables,
				Check:           checkNoGrantsPresent(),
			},
		},
	})
}

func TestAcc_Grants_Of_Invalid_NoAttribute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/Of/Invalid/NoAttribute"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Grants_FutureIn_Database(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := config.Variables{
		"database": config.StringVariable(databaseName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureIn/Database"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneFutureGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_FutureIn_Schema(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := config.Variables{
		"database": config.StringVariable(databaseName),
		"schema":   config.StringVariable(schemaName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureIn/Schema"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneFutureGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_FutureIn_Invalid_NoAttribute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureIn/Invalid/NoAttribute"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Grants_FutureIn_Invalid_SchemaNameNotFullyQualified(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureIn/Invalid/SchemaNameNotFullyQualified"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_Grants_FutureTo_Role(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := config.Variables{
		"database": config.StringVariable(databaseName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureTo/Role"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneFutureGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_FutureTo_DatabaseRole(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := config.Variables{
		"database":      config.StringVariable(databaseName),
		"database_role": config.StringVariable(databaseRoleName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureTo/DatabaseRole"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneFutureGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_FutureTo_Invalid_NoAttribute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureTo/Invalid/NoAttribute"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_Grants_FutureTo_Invalid_DatabaseRoleIdInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureTo/Invalid/DatabaseRoleIdInvalid"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}
