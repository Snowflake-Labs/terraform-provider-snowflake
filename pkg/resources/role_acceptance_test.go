package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Role(t *testing.T) {
	prefix := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	prefix2 := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: rConfig(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_role.w", "name", prefix),
					resource.TestCheckResourceAttr("snowflake_role.w", "comment", "test comment"),
				),
			},
			// RENAME
			{
				Config: rConfig(prefix2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_role.w", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_role.w", "comment", "test comment"),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: rConfig2(prefix2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_role.w", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_role.w", "comment", "test comment 2"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_role.w",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func rConfig(prefix string) string {
	s := `
resource "snowflake_role" "w" {
	name = "%s"
	comment = "test comment"
}
`
	return fmt.Sprintf(s, prefix)
}

func rConfig2(prefix string) string {
	s := `
resource "snowflake_role" "w" {
	name = "%s"
	comment = "test comment 2"
}
`
	return fmt.Sprintf(s, prefix)
}
