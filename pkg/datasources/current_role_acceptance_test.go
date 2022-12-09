package datasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCurrentRole(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: currentRole(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_current_role.p", "current_role"),
				),
			},
		},
	})
}

func currentRole() string {
	s := `
	data snowflake_current_role p {}
	`
	return s
}
