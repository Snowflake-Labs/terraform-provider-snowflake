package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Account_complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	password := acc.TestClient().Ids.AlphaContaining("123ABC")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Account),
		// this errors with: Error running post-test destroy, there may be dangling resources: exit status 1
		// unless we change the resource to return nil on destroy then this is unavoidable
		Steps: []resource.TestStep{
			{
				Config: accountConfig(id.Name(), password, "Terraform acceptance test", 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_account.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_account.test", "admin_name", "someadmin"),
					resource.TestCheckResourceAttr("snowflake_account.test", "first_name", "Ad"),
					resource.TestCheckResourceAttr("snowflake_account.test", "last_name", "Min"),
					resource.TestCheckResourceAttr("snowflake_account.test", "email", "admin@example.com"),
					resource.TestCheckResourceAttr("snowflake_account.test", "must_change_password", "false"),
					resource.TestCheckResourceAttr("snowflake_account.test", "edition", "BUSINESS_CRITICAL"),
					resource.TestCheckResourceAttr("snowflake_account.test", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_account.test", "grace_period_in_days", "3"),
				),
				Destroy: false,
			},
			// Change Grace Period In Days
			{
				Config: accountConfig(id.Name(), password, "Terraform acceptance test", 4),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account.test", "grace_period_in_days", "4"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_account.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"admin_name",
					"admin_password",
					"admin_rsa_public_key",
					"email",
					"must_change_password",
					"first_name",
					"last_name",
					"grace_period_in_days",
				},
			},
		},
	})
}

func accountConfig(name string, password string, comment string, gracePeriodInDays int) string {
	return fmt.Sprintf(`
data "snowflake_current_account" "current" {}

resource "snowflake_account" "test" {
  name = "%s"
  admin_name = "someadmin"
  admin_password = "%s"
  first_name = "Ad"
  last_name = "Min"
  email = "admin@example.com"
  must_change_password = false
  edition = "BUSINESS_CRITICAL"
  comment = "%s"
  region = data.snowflake_current_account.current.region
  grace_period_in_days = %d
}
`, name, password, comment, gracePeriodInDays)
}
