package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Schema(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: schemaConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", accName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", "Terraform acceptance test"),
					checkBool("snowflake_schema.test", "is_transient", false), // this is from user_acceptance_test.go
					checkBool("snowflake_schema.test", "is_managed", false),
				),
			},
		},
	})
}

func schemaConfig(n string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%v"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test" {
	name = "%v"
	database = snowflake_database.test.name
	comment = "Terraform acceptance test"
}
`, n, n)
}
