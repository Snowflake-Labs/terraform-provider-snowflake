package resources_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Role(t *testing.T) {
	name := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	name2 := "5tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: roleBasicConfig(name, "test comment"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_role.role", "name", name),
					resource.TestCheckResourceAttr("snowflake_role.role", "comment", "test comment"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_role.role",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// RENAME
			{
				Config: roleBasicConfig(name2, "test comment"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_role.role", "name", name2),
					resource.TestCheckResourceAttr("snowflake_role.role", "comment", "test comment"),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: roleBasicConfig(name2, "test comment 2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_role.role", "name", name2),
					resource.TestCheckResourceAttr("snowflake_role.role", "comment", "test comment 2"),
				),
			},
		},
	})
}

func TestAcc_AccountRole_basic(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	configVariables := map[string]config.Variable{
		"name":    config.StringVariable(name),
		"comment": config.StringVariable(comment),
	}
	resourceName := "snowflake_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckAccountRoleDestroy(name),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
					resource.TestCheckResourceAttr(resourceName, "id", name),
				),
			},
			// test import
			{
				ConfigDirectory:   config.TestNameDirectory(),
				ConfigVariables:   configVariables,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_AccountRole_updates(t *testing.T) {
	configVariables := func(name string, comment string) config.Variables {
		return config.Variables{
			"name":    config.StringVariable(name),
			"comment": config.StringVariable(comment),
		}
	}

	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	newName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	NewComment := "updated comment with 'single' quotes"
	resourceName := "snowflake_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckAccountRoleDestroy(name),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables(name, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
					resource.TestCheckResourceAttr(resourceName, "id", name),
				),
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables(newName, NewComment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", newName),
					resource.TestCheckResourceAttr(resourceName, "comment", NewComment),
					resource.TestCheckResourceAttr(resourceName, "id", newName),
				),
			},
			// test import
			{
				ConfigDirectory:   config.TestNameDirectory(),
				ConfigVariables:   configVariables(newName, NewComment),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAccountRoleDestroy(accountRoleName string) func(state *terraform.State) error {
	return func(state *terraform.State) error {
		db := acc.TestAccProvider.Meta().(*sql.DB)
		client := sdk.NewClientFromDB(db)
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "snowflake_role" {
				continue
			}
			ctx := context.Background()
			id := sdk.NewAccountObjectIdentifier(rs.Primary.Attributes["name"])
			_, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(id))
			if err == nil {
				return fmt.Errorf("account role %v still exists", accountRoleName)
			}
		}
		return nil
	}
}

func roleBasicConfig(name, comment string) string {
	s := `
resource "snowflake_role" "role" {
	name = "%s"
	comment = "%s"
}
`
	return fmt.Sprintf(s, name, comment)
}
