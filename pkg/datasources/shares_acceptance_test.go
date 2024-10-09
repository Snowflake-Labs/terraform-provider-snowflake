package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Shares(t *testing.T) {
	shareName := acc.TestClient().Ids.Alpha()
	shareName2 := acc.TestClient().Ids.Alpha()
	comment := random.Comment()
	pattern := shareName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: shares(shareName, shareName2, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_shares.r", "shares.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_shares.r", "shares.0.name"),
				),
			},
			{
				Config: sharesPattern(shareName, pattern, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_shares.r", "shares.#"),
					resource.TestCheckResourceAttr("data.snowflake_shares.r", "shares.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_shares.r", "shares.0.kind", "OUTBOUND"),
					resource.TestCheckResourceAttr("data.snowflake_shares.r", "shares.0.comment", comment),
				),
			},
		},
	})
}

func shares(shareName, shareName2, comment string) string {
	return fmt.Sprintf(`
		resource snowflake_share "test_share" {
			name = "%v"
			comment = "%v"
		}
		resource snowflake_share "test_share_2" {
			name = "%v"
			comment = "%v"
		}
		data snowflake_shares "r" {
			depends_on = [
				snowflake_share.test_share,
				snowflake_share.test_share_2,
			]
		}
	`, shareName, comment, shareName2, comment)
}

func sharesPattern(shareName, pattern, comment string) string {
	return fmt.Sprintf(`
		resource snowflake_share "test_share" {
			name = "%v"
			comment = "%v"
		}

		data snowflake_shares "r" {
			pattern = "%v"
			depends_on = [
				snowflake_share.test_share,
			]
		}
	`, shareName, comment, pattern)
}
