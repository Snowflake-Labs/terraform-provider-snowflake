package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

func TestAcc_AccountParameter_PREVENT_LOAD_FROM_INLINE_URL(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountParameterBasic("PREVENT_LOAD_FROM_INLINE_URL", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "key", "PREVENT_LOAD_FROM_INLINE_URL"),
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "value", "true"),
				),
			},
		},
	})
}

func TestAcc_AccountParameter_CLIENT_METADATA_USE_SESSION_DATABASE(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountParameterBasic("CLIENT_METADATA_USE_SESSION_DATABASE", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "key", "CLIENT_METADATA_USE_SESSION_DATABASE"),
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "value", "true"),
				),
			},
		},
	})
}
