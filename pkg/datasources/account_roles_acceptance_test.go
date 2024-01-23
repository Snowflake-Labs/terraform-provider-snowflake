package datasources_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_AccountRoles_basic(t *testing.T) {
	accountRoleNamePrefix := "account_roles_test_prefix_"
	accountRoleName1 := accountRoleNamePrefix + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountRoleName2 := accountRoleNamePrefix + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accountRoleName3 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	configVariables := config.Variables{
		"account_role_name_1": config.StringVariable(accountRoleName1),
		"account_role_name_2": config.StringVariable(accountRoleName2),
		"account_role_name_3": config.StringVariable(accountRoleName3),
		"pattern":             config.StringVariable(accountRoleNamePrefix + "%"),
		"comment":             config.StringVariable(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_account_roles.test", "roles.#", "2"),
					containsAccountRole(accountRoleName1, comment),
					containsAccountRole(accountRoleName2, comment),
					func(state *terraform.State) error {
						err := containsAccountRole(accountRoleName3, comment)(state)
						if err.Error() == fmt.Sprintf("role %s not found", accountRoleName3) {
							return nil
						}
						return fmt.Errorf("expected %s not to be present", accountRoleName3)
					},
				),
			},
		},
	})
}

func containsAccountRole(name string, comment string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "snowflake_account_roles" {
				continue
			}

			iter, err := strconv.ParseInt(rs.Primary.Attributes["roles.#"], 10, 32)
			if err != nil {
				return err
			}

			for i := 0; i < int(iter); i++ {
				if rs.Primary.Attributes[fmt.Sprintf("roles.%d.name", i)] == name {
					actualComment := rs.Primary.Attributes[fmt.Sprintf("roles.%d.comment", i)]
					if actualComment != comment {
						return fmt.Errorf("expected comment: %s, but got: %s", comment, actualComment)
					}

					return nil
				}
			}
		}

		return fmt.Errorf("role %s not found", name)
	}
}
