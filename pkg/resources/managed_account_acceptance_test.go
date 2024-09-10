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

const (
	managedAccountComment = "Created by a Terraform acceptance test"
)

func TestAcc_ManagedAccount(t *testing.T) {
	// TODO [SNOW-1011985]: unskip the tests
	testenvs.SkipTestIfSet(t, testenvs.SkipManagedAccountTest, "error: 090337 (23001): Number of managed accounts allowed exceeded the limit. Please contact Snowflake support")

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	adminName := acc.TestClient().Ids.Alpha()
	adminPass := acc.TestClient().Ids.AlphaWithPrefix("A1")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ManagedAccount),
		Steps: []resource.TestStep{
			{
				Config: managedAccountConfig(id.Name(), adminName, adminPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_managed_account.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_managed_account.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_managed_account.test", "admin_name", adminName),
					resource.TestCheckResourceAttr("snowflake_managed_account.test", "admin_password", adminPass),
					resource.TestCheckResourceAttr("snowflake_managed_account.test", "comment", managedAccountComment),
					resource.TestCheckResourceAttr("snowflake_managed_account.test", "type", "READER"),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_managed_account.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"admin_name", "admin_password"},
			},
		},
	})
}

func managedAccountConfig(accName, aName, aPass string) string {
	return fmt.Sprintf(`
resource "snowflake_managed_account" "test" {
	name           = "%v"
	admin_name     = "%v"
	admin_password = "%v"
	comment        = "%v"
}
`, accName, aName, aPass, managedAccountComment)
}
