package datasources_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationAccounts(t *testing.T) {
	accountName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: organizationAccounts(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_organization_accounts.t", "accounts.#"),
					resource.TestCheckResourceAttr("data.snowflake_organization_accounts.t", "accounts.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_organization_accounts.t", "accounts.0.name", accountName),
				),
			},
		},
	})
}

func organizationAccounts() string {
	return "data snowflake_organization_accounts p {}"
}
