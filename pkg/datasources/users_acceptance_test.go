package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Users(t *testing.T) {
	userName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.Test(t, resource.TestCase{
		Providers: providers(),
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
		default_role="foo"
		default_secondary_roles = ["ALL"]
		default_namespace="foo"
	}

	data snowflake_users "u" {
		pattern = "%s"
		depends_on = [snowflake_user.u]
	}
	`, userName, userName, userName)
}
