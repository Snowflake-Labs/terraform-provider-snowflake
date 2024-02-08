package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Role(t *testing.T) {
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: role(roleName, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_role.t", "name", roleName),
					resource.TestCheckResourceAttr("data.snowflake_role.t", "comment", comment),
				),
			},
		},
	})
}

func role(roleName, comment string) string {
	return fmt.Sprintf(`
		resource snowflake_role "test_role" {
			name = "%v"
			comment = "%v"
		}
		data snowflake_role "t" {
			depends_on = [snowflake_role.test_role]
			name = "%v"
		}
	`, roleName, comment, roleName)
}
