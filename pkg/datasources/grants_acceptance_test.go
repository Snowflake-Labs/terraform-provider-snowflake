package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Grants(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantsAccount(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_grants.g", "grants.#"),
				),
			},
		},
	})
}

func grantsAccount() string {
	s := `
data "snowflake_grants" "g" {
	grants_on {
		account = true
	}
}
`
	return s
}
