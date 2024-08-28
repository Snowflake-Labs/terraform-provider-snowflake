package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Users(t *testing.T) {
	userName := acc.TestClient().Ids.Alpha()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: users(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_users.u", "users.#"),
					resource.TestCheckResourceAttr("data.snowflake_users.u", "users.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_users.u", "users.0.name", userName),
					resource.TestCheckResourceAttr("data.snowflake_users.u", "users.0.disabled", "false"),
				),
			},
		},
	})
}

func users(userName string) string {
	return fmt.Sprintf(`
	resource "snowflake_user" "u" {
		name = "%s"
		comment = "test comment"
		login_name = "%s_login"
		display_name = "Display Name"
		first_name = "Alex"
		last_name = "Kita"
		email = "fake@email.com"
		disabled = false
		default_warehouse="foo"
		default_role="FOO"
		default_secondary_roles = ["ALL"]
		default_namespace="FOO"
	}

	data snowflake_users "u" {
		pattern = "%s"
		depends_on = [snowflake_user.u]
	}
	`, userName, userName, userName)
}
