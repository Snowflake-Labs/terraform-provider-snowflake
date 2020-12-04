package resources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	managedAccountComment = "Created by a Terraform acceptance test"
)

func TestAcc_ManagedAccount(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_MANAGED_ACCOUNT_TEST"); ok {
		t.Skip("Skipping TestAccManagedAccount")
	}

	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	adminName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	adminPass := fmt.Sprintf("A1%v", acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: managedAccountConfig(accName, adminName, adminPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_managed_account.test", "name", accName),
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
