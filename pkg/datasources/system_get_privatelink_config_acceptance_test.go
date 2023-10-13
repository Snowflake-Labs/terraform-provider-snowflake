package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_SystemGetPrivateLinkConfig_aws(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: privateLinkConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "account_name"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "account_url"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "ocsp_url"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "aws_vpce_id"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "regionless_account_url"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "regionless_snowsight_url"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "snowsight_url"),
				),
			},
		},
	})
}

func privateLinkConfig() string {
	s := `
	data snowflake_system_get_privatelink_config p {}
	`
	return s
}
