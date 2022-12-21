package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_SessionParameter(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: sessionParameterBasic("AUTOCOMMIT", "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_session_parameter.p", "key", "AUTOCOMMIT"),
					resource.TestCheckResourceAttr("snowflake_session_parameter.p", "value", "false"),
				),
			},
		},
	})
}

func sessionParameterBasic(key, value string) string {
	s := `
resource "snowflake_session_parameter" "p" {
	key = "%s"
	value = "%s"
}
`
	return fmt.Sprintf(s, key, value)
}
