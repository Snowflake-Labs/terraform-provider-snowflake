package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_AccountGrant_defaults(t *testing.T) {
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountGrantConfig(roleName, "MONITOR USAGE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_grant.test", "privilege", "MONITOR USAGE"),
				),
			},
			// UPDATE ALL PRIVILEGES
			{
				Config: accountGrantConfig(roleName, "ALL PRIVILEGES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_grant.test", "privilege", "ALL PRIVILEGES"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_account_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_AccountGrantManagedTask(t *testing.T) {
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountGrantConfig(roleName, "EXECUTE MANAGED TASK"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_grant.test", "privilege", "EXECUTE MANAGED TASK"),
				),
			},
		},
	})
}

func TestAcc_AccountGrantManageSupportCases(t *testing.T) {
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountGrantConfig(roleName, "MANAGE ACCOUNT SUPPORT CASES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_grant.test", "privilege", "MANAGE ACCOUNT SUPPORT CASES"),
				),
			},
		},
	})
}

func TestAcc_AccountGrantManageWarehouses(t *testing.T) {
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountGrantConfig(roleName, "MANAGE WAREHOUSES"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_grant.test", "privilege", "MANAGE WAREHOUSES"),
				),
			},
		},
	})
}

func accountGrantConfig(role, privilege string) string {
	return fmt.Sprintf(`

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_account_grant" "test" {
  roles     = [snowflake_role.test.name]
  privilege = "%s"
}
`, role, privilege)
}
