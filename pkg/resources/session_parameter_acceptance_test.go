package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SessionParameterWithUser(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)

	user, userCleanup := acc.TestClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: sessionParameterWithUser(user.ID(), "BINARY_OUTPUT_FORMAT", "BASE64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_session_parameter.p", "key", "BINARY_OUTPUT_FORMAT"),
					resource.TestCheckResourceAttr("snowflake_session_parameter.p", "value", "BASE64"),
					resource.TestCheckResourceAttr("snowflake_session_parameter.p", "user", user.ID().Name()),
					resource.TestCheckResourceAttr("snowflake_session_parameter.p", "on_account", "false"),
				),
			},
		},
	})
}

func TestAcc_SessionParameterOnAccount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: sessionParameterOnAccount("AUTOCOMMIT", "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_session_parameter.p", "key", "AUTOCOMMIT"),
					resource.TestCheckResourceAttr("snowflake_session_parameter.p", "value", "false"),
					resource.TestCheckResourceAttr("snowflake_session_parameter.p", "on_account", "true"),
				),
			},
		},
	})
}

func sessionParameterWithUser(userId sdk.AccountObjectIdentifier, key, value string) string {
	return fmt.Sprintf(`
resource "snowflake_session_parameter" "p" {
	key = "%[2]s"
	value = "%[3]s"
	user = %[1]s
}
`, userId.FullyQualifiedName(), key, value)
}

func sessionParameterOnAccount(key, value string) string {
	s := `
resource "snowflake_session_parameter" "p" {
	key = "%s"
	value = "%s"
	on_account = true
}
`
	return fmt.Sprintf(s, key, value)
}
