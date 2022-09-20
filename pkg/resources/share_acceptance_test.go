package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	shareComment = "Created by a Terraform acceptance test"
)

func TestAcc_Share(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	account2 := os.Getenv("SNOWFLAKE_ACCOUNT_SECOND")
	account3 := os.Getenv("SNOWFLAKE_ACCOUNT_THIRD")

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: shareConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_share.test", "comment", shareComment),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.#", "0"),
				),
			},
			{
				Config: shareConfigTwoAccounts(name, account2, account3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.#", "2"),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.0", account2),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.1", account3),
				),
			},
			{
				Config: shareConfigOneAccount(name, account2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.0", account2),
				),
			},
			{
				Config: shareConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.#", "0"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_share.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func shareConfig(name string) string {
	return fmt.Sprintf(`
resource "snowflake_share" "test" {
	name           = "%v"
	comment        = "%v"
}
`, name, shareComment)
}

func shareConfigOneAccount(name string, account2 string) string {
	return fmt.Sprintf(`
resource "snowflake_share" "test" {
	name           = "%v"
	comment        = "%v"
	accounts       = ["%v"]
}
`, name, shareComment, account2)
}

func shareConfigTwoAccounts(name string, account2 string, account3 string) string {
	return fmt.Sprintf(`
resource "snowflake_share" "test" {
	name           = "%v"
	comment        = "%v"
	accounts       = ["%v", "%v"]
}
`, name, shareComment, account2, account3)
}
