package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Role(t *testing.T) {
	name := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	name2 := "5tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
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

func roleBasicConfig(name, comment string) string {
	return fmt.Sprintf(`
		resource "snowflake_role" "role" {
			name = "%s"
			comment = "%s"
		}
	`, name, comment)
}
