package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

const (
	shareComment = "Created by a Terraform acceptance test"
)

func TestAccShare(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: shareConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_share.test", "comment", shareComment),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_share.test",
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}

func shareConfig(name string) string {
	return fmt.Sprintf(`
resource "snowflake_share" "test" {
	name           = "%v"
	comment        = "%v"
}
`, name, shareComment)
}
