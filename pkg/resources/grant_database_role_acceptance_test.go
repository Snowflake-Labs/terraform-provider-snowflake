package resources_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_GrantDatabaseRole_databaseRole(t *testing.T) {
	databaseRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	parentDatabaseRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_grant_database_role.g"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                  config.StringVariable(acc.TestDatabaseName),
			"database_role_name":        config.StringVariable(databaseRoleName),
			"parent_database_role_name": config.StringVariable(parentDatabaseRoleName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckGrantDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/database_role"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", fmt.Sprintf(`"%v"."%v"`, acc.TestDatabaseName, databaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "parent_database_role_name", fmt.Sprintf(`"%v"."%v"`, acc.TestDatabaseName, parentDatabaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"."%v"|DATABASE ROLE|"%v"."%v"`, acc.TestDatabaseName, databaseRoleName, acc.TestDatabaseName, parentDatabaseRoleName)),
				),
			},
			// test import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/database_role"),
				ConfigVariables:   m(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_GrantDatabaseRole_issue2402(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	parentDatabaseRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_grant_database_role.g"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":                  config.StringVariable(databaseName),
			"database_role_name":        config.StringVariable(databaseRoleName),
			"parent_database_role_name": config.StringVariable(parentDatabaseRoleName),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckGrantDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_GrantDatabaseRole/issue2402"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", fmt.Sprintf(`"%v"."%v"`, databaseName, databaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "parent_database_role_name", fmt.Sprintf(`"%v"."%v"`, databaseName, parentDatabaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"."%v"|DATABASE ROLE|"%v"."%v"`, databaseName, databaseRoleName, databaseName, parentDatabaseRoleName)),
				),
			},
		},
	})
}

func TestAcc_GrantDatabaseRole_accountRole(t *testing.T) {
	databaseRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	parentRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_grant_database_role.g"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"database":           config.StringVariable(acc.TestDatabaseName),
			"database_role_name": config.StringVariable(databaseRoleName),
			"parent_role_name":   config.StringVariable(parentRoleName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckGrantDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/account_role"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", fmt.Sprintf(`"%v"."%v"`, acc.TestDatabaseName, databaseRoleName)),
					resource.TestCheckResourceAttr(resourceName, "parent_role_name", fmt.Sprintf("%v", parentRoleName)),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`"%v"."%v"|ROLE|"%v"`, acc.TestDatabaseName, databaseRoleName, parentRoleName)),
				),
			},
			// test import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/account_role"),
				ConfigVariables:   m(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2410 is fixed
func TestAcc_GrantDatabaseRole_share(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseRoleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseRoleId := sdk.NewDatabaseObjectIdentifier(databaseName, databaseRoleName).FullyQualifiedName()
	shareName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	shareId := sdk.NewAccountObjectIdentifier(shareName).FullyQualifiedName()
	resourceName := "snowflake_grant_database_role.test"
	configVariables := func() config.Variables {
		return config.Variables{
			"database":           config.StringVariable(databaseName),
			"database_role_name": config.StringVariable(databaseRoleName),
			"share_name":         config.StringVariable(shareName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckGrantDatabaseRoleDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/share"),
				ConfigVariables: configVariables(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database_role_name", databaseRoleId),
					resource.TestCheckResourceAttr(resourceName, "share_name", shareName),
					resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf(`%v|%v|%v`, databaseRoleId, "SHARE", shareId)),
				),
			},
			// test import
			{
				ConfigDirectory:   config.StaticDirectory("testdata/TestAcc_GrantDatabaseRole/share"),
				ConfigVariables:   configVariables(),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGrantDatabaseRoleDestroy(s *terraform.State) error {
	client := acc.TestAccProvider.Meta().(*provider.Context).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_grant_database_role" {
			continue
		}
		ctx := context.Background()
		id := rs.Primary.ID
		ids := strings.Split(id, "|")
		databaseRoleName := ids[0]
		objectType := ids[1]
		parentRoleName := ids[2]
		grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
			Of: &sdk.ShowGrantsOf{
				DatabaseRole: sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(databaseRoleName),
			},
		})
		if err != nil {
			continue
		}
		for _, grant := range grants {
			if grant.GrantedTo == sdk.ObjectType(objectType) {
				if grant.GranteeName.FullyQualifiedName() == parentRoleName {
					return fmt.Errorf("database role grant %v still exists", grant)
				}
			}
		}
	}
	return nil
}
