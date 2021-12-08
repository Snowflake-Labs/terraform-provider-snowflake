package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAccountGrant_defaults(t *testing.T) {
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: accountGrantConfig(roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_grant.test", "privilege", "MONITOR USAGE"),
				),
			},
		},
	})
}

func accountGrantConfig(role string) string {
	return fmt.Sprintf(`

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_account_grant" "test" {
  roles          = [snowflake_role.test.name]
  privilege = "MONITOR USAGE"
}
`, role)
}

func TestAcc_AccountGrantManagedTask(t *testing.T) {
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: accountGrantManagedTaskConfig(roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_grant.test", "privilege", "EXECUTE MANAGED TASK"),
				),
			},
		},
	})
}

func accountGrantManagedTaskConfig(role string) string {
	return fmt.Sprintf(`

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_account_grant" "test" {
  roles     = [snowflake_role.test.name]
  privilege = "EXECUTE MANAGED TASK"
}
`, role)
}
