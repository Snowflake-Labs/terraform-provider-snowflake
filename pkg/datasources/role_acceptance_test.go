package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestRole(t *testing.T) {
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{

		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: role(roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_role.t", "name", roleName),
				),
			},
		},
	})
}

func role(roleName string) string {
	return fmt.Sprintf(`
		resource snowflake_role "test_role" {
			name = "%v"
		}
		data snowflake_role "t" {
			depends_on = [snowflake_role.test_role]
			name = "%v"
		}
	`, roleName, roleName)
}
