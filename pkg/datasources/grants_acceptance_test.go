package datasources_test

import (
	"context"
	"regexp"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

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

// TODO [SNOW-1284382]: Implement after snowflake_application and snowflake_application_role resources are introduced.
func TestAcc_Grants_To_Application(t *testing.T) {
	t.Skip("Skipped until snowflake_application and snowflake_application_role resources are introduced. Currently, behavior tested in application_roles_gen_integration_test.go.")
}

// TODO [SNOW-1284382]: Implement after snowflake_application and snowflake_application_role resources are introduced.
func TestAcc_Grants_To_ApplicationRole(t *testing.T) {
	t.Skip("Skipped until snowflake_application and snowflake_application_role resources are introduced. Currently, behavior tested in application_roles_gen_integration_test.go.")
}

func TestAcc_Grants_To_AccountRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/AccountRole"),
				Check:           checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_To_DatabaseRole(t *testing.T) {
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/DatabaseRole"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneGrantPresent(),
			},
		},
	})
}

func TestAcc_Grants_To_User(t *testing.T) {
	user := acc.TestClient().Context.CurrentUser(t)
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
				Check:           checkAtLeastOneGrantPresentLimited(),
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

// TODO [SNOW-1284382]: Implement after SHOW GRANTS TO SHARE <share_name> IN APPLICATION PACKAGE <app_package_name> syntax starts working.
func TestAcc_Grants_To_ShareWithApplicationPackage(t *testing.T) {
	t.Skip("Skipped until SHOW GRANTS TO SHARE <share_name> IN APPLICATION PACKAGE <app_package_name> syntax starts working.")
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

func TestAcc_Grants_To_Invalid_DatabaseRoleIdInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/Invalid/DatabaseRoleIdInvalid"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_Grants_To_Invalid_ApplicationRoleIdInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/To/Invalid/ApplicationRoleIdInvalid"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_Grants_Of_AccountRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/Of/AccountRole"),
				Check:           checkAtLeastOneGrantPresentLimited(),
			},
		},
	})
}

func TestAcc_Grants_Of_DatabaseRole(t *testing.T) {
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/Of/DatabaseRole"),
				ConfigVariables: configVariables,
				Check:           checkAtLeastOneGrantPresentLimited(),
			},
		},
	})
}

// TODO [SNOW-1284382]: Implement after snowflake_application and snowflake_application_role resources are introduced.
func TestAcc_Grants_Of_ApplicationRole(t *testing.T) {
	t.Skip("Skipped until snowflake_application and snowflake_application_role resources are introduced. Currently, behavior tested in application_roles_gen_integration_test.go.")
}

// TODO [SNOW-1284394]: Unskip the test
func TestAcc_Grants_Of_Share(t *testing.T) {
	t.Skip("TestAcc_Share are skipped")
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	shareName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	getSecondaryAccountIdentifier := func(t *testing.T) *sdk.AccountIdentifier {
		t.Helper()

		client := acc.Client(t)
		secondaryClient := acc.SecondaryClient(t)
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

	accountId := getSecondaryAccountIdentifier(t)
	require.NotNil(t, accountId)

	configVariables := config.Variables{
		"database": config.StringVariable(databaseName),
		"share":    config.StringVariable(shareName),
		"account":  config.StringVariable(accountId.FullyQualifiedName()),
	}
	datasourceName := "data.snowflake_grants.test"

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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "grants.#"),
					resource.TestCheckNoResourceAttr(datasourceName, "grants.0.created_on"),
				),
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

func TestAcc_Grants_Of_Invalid_DatabaseRoleIdInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/Of/Invalid/DatabaseRoleIdInvalid"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_Grants_Of_Invalid_ApplicationRoleIdInvalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/Of/Invalid/ApplicationRoleIdInvalid"),
				PlanOnly:        true,
				ExpectError:     regexp.MustCompile("Error: Invalid identifier type"),
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

func TestAcc_Grants_FutureTo_AccountRole(t *testing.T) {
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Grants/FutureTo/AccountRole"),
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

func checkAtLeastOneGrantPresentLimited() resource.TestCheckFunc {
	datasourceName := "data.snowflake_grants.test"
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(datasourceName, "grants.#"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.created_on"),
		resource.TestCheckResourceAttrSet(datasourceName, "grants.0.granted_to"),
	)
}
