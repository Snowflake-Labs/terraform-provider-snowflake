package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_AccountParameter(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountParameterBasic("ALLOW_ID_TOKEN", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "key", "ALLOW_ID_TOKEN"),
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "value", "true"),
				),
			},
		},
	})
}

func accountParameterBasic(key, value string) string {
	s := `
resource "snowflake_account_parameter" "p" {
	key = "%s"
	value = "%s"
}
`
	return fmt.Sprintf(s, key, value)
}
