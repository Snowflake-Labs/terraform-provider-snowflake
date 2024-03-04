package resources_test

import (
	"fmt"
	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"strings"
	"testing"
)

func TestAcc_GrantOwnership_OnObject(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountRoleName := sdk.NewAccountObjectIdentifier(strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))).FullyQualifiedName()
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
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE|%s", accountRoleName, databaseName)),
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
