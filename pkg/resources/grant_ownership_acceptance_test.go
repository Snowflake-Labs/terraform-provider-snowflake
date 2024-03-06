package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantOwnership_OnObject(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseFullyQualifiedName := sdk.NewAccountObjectIdentifier(databaseName).FullyQualifiedName()
	accountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountRoleFullyQualifiedName := sdk.NewAccountObjectIdentifier(accountRoleName).FullyQualifiedName()
	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil, // TODO
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnObject"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

//func TestAcc_GrantOwnership_OnAll(t *testing.T) {
//	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
//	databaseFullyQualifiedName := sdk.NewAccountObjectIdentifier(databaseName).FullyQualifiedName()
//	accountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
//	accountRoleFullyQualifiedName := sdk.NewAccountObjectIdentifier(accountRoleName).FullyQualifiedName()
//	configVariables := config.Variables{
//		"account_role_name": config.StringVariable(accountRoleName),
//		"database_name":     config.StringVariable(databaseName),
//	}
//	resourceName := "snowflake_grant_ownership.test"
//
//	resource.Test(t, resource.TestCase{
//		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
//		PreCheck:                 func() { acc.TestAccPreCheck(t) },
//		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
//			tfversion.RequireAbove(tfversion.Version1_5_0),
//		},
//		CheckDestroy: nil, // TODO
//		Steps: []resource.TestStep{
//			{
//				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnAll"),
//				ConfigVariables: configVariables,
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
//					resource.TestCheckResourceAttr(resourceName, "on.0.all.0.object_type_plural", "TABLES"),
//					resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseName),
//					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
//				),
//			},
//			{
//				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnAll"),
//				ConfigVariables:   configVariables,
//				ResourceName:      resourceName,
//				ImportState:       true,
//				ImportStateVerify: true,
//			},
//		},
//	})
//}

func TestAcc_GrantOwnership_OnFuture(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseFullyQualifiedName := sdk.NewAccountObjectIdentifier(databaseName).FullyQualifiedName()
	accountRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountRoleFullyQualifiedName := sdk.NewAccountObjectIdentifier(accountRoleName).FullyQualifiedName()
	configVariables := config.Variables{
		"account_role_name": config.StringVariable(accountRoleName),
		"database_name":     config.StringVariable(databaseName),
	}
	resourceName := "snowflake_grant_ownership.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil, // TODO
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnFuture"),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleName),
					resource.TestCheckResourceAttr(resourceName, "on.0.future.0.object_type_plural", "VIEWS"),
					resource.TestCheckResourceAttr(resourceName, "on.0.future.0.in_database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnFuture|VIEWS|InDatabase|%s", accountRoleFullyQualifiedName, databaseFullyQualifiedName)),
				),
			},
			{
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_GrantOwnership/OnFuture"),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
