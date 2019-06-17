package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccView(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: viewConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "Terraform test resource"),
					checkBool("snowflake_view.test", "is_secure", true), // this is from user_acceptance_test.go
					resource.TestCheckResourceAttr("snowflake_view.test", "view_text", "SELECT a, b FROM MY_TABLE WHERE this = something"),
				),
			},
		},
	})
}

func viewConfig(n string) string {
	return fmt.Sprintf(`
resource "snowflake_view" "test" {
	name                = "%v"
	comment             = "Terraform test resource"
	is_secure           = true
	statement           = "SELECT a, b FROM MY_TABLE WHERE this = ?"
	statement_arguments = ["something"]
}
`, n)
}
