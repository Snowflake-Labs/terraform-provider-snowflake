package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	shareComment = "Created by a Terraform acceptance test"
)

func TestAcc_Share(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_SHARE_TESTS"); ok {
		t.Skip("Skipping TestAccShare")
	}

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
				ResourceName:      "snowflake_share.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func shareConfig(name string) string {
	return fmt.Sprintf(`
resource "snowflake_share" "test" {
	name           = "%v"
	comment        = "%v"
	accounts       = ["PC37737"]
}
`, name, shareComment)
}
