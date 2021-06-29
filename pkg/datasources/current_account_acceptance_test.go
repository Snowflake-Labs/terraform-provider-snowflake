package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCurrentAccount(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: currentAccount(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_current_account.p", "account"),
					resource.TestCheckResourceAttrSet("data.snowflake_current_account.p", "region"),
					resource.TestCheckResourceAttrSet("data.snowflake_current_account.p", "url"),
				),
			},
		},
	})
}

func currentAccount() string {
	s := `
	data snowflake_current_account p {}
	`
	return s
}
