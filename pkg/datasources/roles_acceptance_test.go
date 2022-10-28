package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Roles(t *testing.T) {
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleName2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: roles(roleName, roleName2, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_roles.r", "roles.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_roles.r", "roles.0.name"),
					//resource.TestCheckTypeSetElemAttr("data.snowflake_roles.r", "roles.*", "name"),
					//TODO show roles also includes built in role such as ACCOUNTADMIN, SYSADMIN, etc.
				),
			},
		},
	})
}

func roles(roleName, roleName2, comment string) string {
	return fmt.Sprintf(`
		resource snowflake_role "test_role" {
			name = "%v"
			comment = "%v"
		}
		resource snowflake_role "test_role_2" {
			name = "%v"
			comment = "%v"
		}
		data snowflake_roles "r" {
			depends_on = [
				snowflake_role.test_role,
				snowflake_role.test_role_2,
			]
		}
	`, roleName, comment, roleName2, comment)
}
