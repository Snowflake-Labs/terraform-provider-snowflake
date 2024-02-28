package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Share(t *testing.T) {
	t.Skip("second and third account must be set for Share acceptance tests")
	var account2 string
	var account3 string

	shareComment := "Created by a Terraform acceptance test"
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: shareConfig(name, shareComment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_share.test", "comment", shareComment),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.#", "0"),
				),
			},
			{
				Config: shareConfigTwoAccounts(name, shareComment, account2, account3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.#", "2"),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.0", account2),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.1", account3),
				),
			},
			{
				Config: shareConfigOneAccount(name, shareComment, account2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.#", "1"),
					resource.TestCheckResourceAttr("snowflake_share.test", "accounts.0", account2),
				),
			},
			{
				Config: shareConfig(name, shareComment),
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

func shareConfig(name string, comment string) string {
	return fmt.Sprintf(`
resource "snowflake_share" "test" {
	name           = "%v"
	comment        = "%v"
}
`, name, comment)
}

func shareConfigOneAccount(name string, comment string, account2 string) string {
	return fmt.Sprintf(`
resource "snowflake_share" "test" {
	name           = "%v"
	comment        = "%v"
	accounts       = ["%v"]
}
`, name, comment, account2)
}

func shareConfigTwoAccounts(name string, comment string, account2 string, account3 string) string {
	return fmt.Sprintf(`
resource "snowflake_share" "test" {
	name           = "%v"
	comment        = "%v"
	accounts       = ["%v", "%v"]
}
`, name, comment, account2, account3)
}
