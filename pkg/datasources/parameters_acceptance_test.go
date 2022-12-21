package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Parameters(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: parameters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.0.key"),
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.0.value"),
				),
			},
		},
	})
}

func parameters() string {
	return `data "snowflake_parameters" "p" {}`
}
