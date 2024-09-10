package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AccountParameter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
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
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
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

func TestAcc_AccountParameter_REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountParameterBasic("REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "key", "REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION"),
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "value", "true"),
				),
			},
		},
	})
}

// TODO [SNOW-1528546]: unskip
func TestAcc_AccountParameter_Issue2573(t *testing.T) {
	t.Skipf("The cleanup for parameter is currently incorrect and this test messes with other ones. Skipping until SNOW-1528546 is resolved.")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountParameterBasic("TIMEZONE", "Etc/UTC"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "key", "TIMEZONE"),
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "value", "Etc/UTC"),
				),
			},
			{
				ResourceName:            "snowflake_account_parameter.p",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAcc_AccountParameter_Issue3025(t *testing.T) {
	t.Skipf("The cleanup for parameter is currently incorrect and this test messes with other ones. Skipping until SNOW-1528546 is resolved.")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: accountParameterBasic("OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "key", "OAUTH_ADD_PRIVILEGED_ROLES_TO_BLOCKED_LIST"),
					resource.TestCheckResourceAttr("snowflake_account_parameter.p", "value", "true"),
				),
			},
			{
				ResourceName:            "snowflake_account_parameter.p",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}
