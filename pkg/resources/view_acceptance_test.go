package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccView(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: viewConfig(accName, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "Terraform test resource"),
					checkBool("snowflake_view.test", "is_secure", true), // this is from user_acceptance_test.go
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: viewConfig(accName, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES where ROLE_OWNER like 'foo%%';"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "Terraform test resource"),
					checkBool("snowflake_view.test", "is_secure", true), // this is from user_acceptance_test.go
				),
			},
		},
	})
}

func viewConfig(n string, q string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%v"
}

resource "snowflake_view" "test" {
	name      = "%v"
	comment   = "Terraform test resource"
	database  = snowflake_database.test.name
	is_secure = true
	statement = "%s"
}
`, n, n, q)
}
