package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSystemGetSnowflakePlatformInfo(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: snowflakePlatformInfo(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_snowflake_platform_info.p", "aws_vpc_ids.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_snowflake_platform_info.p", "azure_vnet_subnet_ids.#"),
				),
			},
		},
	})
}

func snowflakePlatformInfo() string {
	s := `
	data snowflake_system_get_snowflake_platform_info "p" {}
	`
	return s
}
