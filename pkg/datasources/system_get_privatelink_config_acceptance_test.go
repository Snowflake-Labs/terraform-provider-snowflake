package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSystemGetPrivateLinkConfig_aws(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: privateLinkConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "account_name"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "account_url"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "ocsp_url"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "aws_vpce_id"),
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
