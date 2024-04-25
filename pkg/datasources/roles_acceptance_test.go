package datasources_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const (
	accountAdmin = "ACCOUNTADMIN"
)

func TestAcc_Roles(t *testing.T) {
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleName2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: roles(roleName, roleName2, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_roles.r", "roles.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_roles.r", "roles.0.name"),
					// resource.TestCheckTypeSetElemAttr("data.snowflake_roles.r", "roles.*", "name"),
					// TODO SHOW ROLES output also includes built in roles, i.e. ACCOUNTADMIN, SYSADMIN, etc.
				),
			},
			{
				Config: rolesPattern(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_roles.r", "roles.#"),
					// resource.TestCheckResourceAttrSet("data.snowflake_roles.r", "roles.0.name"),
					resource.TestCheckResourceAttr("data.snowflake_roles.r", "roles.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_roles.r", "roles.0.name", accountAdmin),
				),
			},
		},
	})
}

func TestAcc_AccountRoles_basic(t *testing.T) {
	accountRoleNamePrefix := "account_roles_test_prefix_"
	accountRoleName1 := acc.TestClient().Ids.AlphaWithPrefix(accountRoleNamePrefix)
	accountRoleName2 := acc.TestClient().Ids.AlphaWithPrefix(accountRoleNamePrefix)
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
					resource.TestCheckResourceAttr("data.snowflake_roles.test", "roles.#", "2"),
					containsAccountRole(accountRoleName1, comment),
					containsAccountRole(accountRoleName2, comment),
					doesntContainAccountRole(accountRoleName3, comment),
				),
			},
		},
	})
}

func doesntContainAccountRole(name string, comment string) func(s *terraform.State) error {
	return func(state *terraform.State) error {
		err := containsAccountRole(name, comment)(state)
		if err.Error() == fmt.Sprintf("role %s not found", name) {
			return nil
		}
		return fmt.Errorf("expected %s not to be present", name)
	}
}

func containsAccountRole(name string, comment string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "snowflake_roles" {
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

func rolesPattern() string {
	return fmt.Sprintf(`
		data snowflake_roles "r" {
			pattern = "%v"
		}
	`, accountAdmin)
}
