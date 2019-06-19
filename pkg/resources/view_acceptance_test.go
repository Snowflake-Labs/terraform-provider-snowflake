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
					resource.TestCheckResourceAttr("snowflake_view.test", "database", "DEMO_DB"),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "Terraform test resource"),
					checkBool("snowflake_view.test", "is_secure", true), // this is from user_acceptance_test.go
				),
			},
		},
	})
}

func viewConfig(n string) string {
	return fmt.Sprintf(`
resource "snowflake_view" "test" {
	name      = "%v"
	comment   = "Terraform test resource"
	database  = "DEMO_DB"
	is_secure = true
	statement = "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
}
`, n)
}
