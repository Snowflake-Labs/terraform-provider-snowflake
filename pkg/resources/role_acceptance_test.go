package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Role(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()
	name2 := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.Role),
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
	name := acc.TestClient().Ids.Alpha()
	comment := random.Comment()
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
		CheckDestroy: acc.CheckDestroy(t, resources.Role),
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

	name := acc.TestClient().Ids.Alpha()
	newName := acc.TestClient().Ids.Alpha()
	comment := random.Comment()
	NewComment := "updated comment with 'single' quotes"
	resourceName := "snowflake_role.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Role),
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

func roleBasicConfig(name, comment string) string {
	s := `
resource "snowflake_role" "role" {
	name = "%s"
	comment = "%s"
}
`
	return fmt.Sprintf(s, name, comment)
}
