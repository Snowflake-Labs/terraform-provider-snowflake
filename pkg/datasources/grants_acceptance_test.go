package datasources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO: tests (examples from the correct ones):
// - on - account
// - on - db object
// - on - schema object
// - on - invalid config - no attribute
// - on - invalid config - missing object type or name
// - to - application
// - to - application role
// - to - role
// - to - user
// - to - share
// - to - share with application package
// - to - invalid config - no attribute
// - to - invalid config - share name missing
// - of - role
// - of - application role
// - of - share
// - of - invalid config - no attribute
// - future in - database
// - future in - schema (both db and sc present)
// - future in - schema (only sc present)
// - future in - invalid config - no attribute
// - future in - invalid config - schema with no schema name
// - future to - role
// - future to - database role
// - future to - invalid config - no attribute
// - future to - invalid config - database role id invalid
func TestAcc_Grants(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
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
