package datasources_test

import (
	"fmt"
	"maps"
	"strconv"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Roles_Complete(t *testing.T) {
	accountRoleNamePrefix := random.AlphaN(10)
	accountRoleName1 := acc.TestClient().Ids.AlphaWithPrefix(accountRoleNamePrefix + "1")
	accountRoleName2 := acc.TestClient().Ids.AlphaWithPrefix(accountRoleNamePrefix + "2")
	accountRoleName3 := acc.TestClient().Ids.Alpha()
	comment := random.Comment()

	commonVariables := config.Variables{
		"account_role_name_1": config.StringVariable(accountRoleName1),
		"account_role_name_2": config.StringVariable(accountRoleName2),
		"account_role_name_3": config.StringVariable(accountRoleName3),
		"comment":             config.StringVariable(comment),
	}

	likeVariables := maps.Clone(commonVariables)
	likeVariables["like"] = config.StringVariable(accountRoleNamePrefix + "%")

	// TODO(SNOW-1353303): Add test case for instance classes after they're available in the provider
	// inClassVariables := maps.Clone(commonVariables)
	// inClassVariables["in_class"] = config.StringVariable("<class name>")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: likeVariables,
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
		if err != nil && err.Error() == fmt.Sprintf("role %s not found", name) {
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
				if rs.Primary.Attributes[fmt.Sprintf("roles.%d.show_output.0.name", i)] == name {
					actualComment := rs.Primary.Attributes[fmt.Sprintf("roles.%d.show_output.0.comment", i)]
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
