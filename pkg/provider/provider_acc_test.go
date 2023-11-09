package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

var testAccProvider *schema.Provider

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"snowflake": func() (tfprotov6.ProviderServer, error) {
		return tf5to6server.UpgradeServer(
			context.Background(),
			testAccProvider.GRPCProvider,
		)
	},
}

func init() {
	testAccProvider = Provider()
}

func TestAcc_ProviderUsernameAndPasswordAuth(t *testing.T) {
	user := os.Getenv("SNOWFLAKE_USER")
	if user == "" {
		t.Skip("SNOWFLAKE_USER must be set")
	}
	password := os.Getenv("SNOWFLAKE_PASSWORD")
	if password == "" {
		t.Skip("SNOWFLAKE_PASSWORD must be set")
	}
	account := os.Getenv("SNOWFLAKE_ACCOUNT")
	if account == "" {
		t.Skip("SNOWFLAKE_ACCOUNT must be set")
	}
	role := os.Getenv("SNOWFLAKE_ROLE")
	if role == "" {
		t.Skip("SNOWFLAKE_ROLE must be set")
	}

	configVars := map[string]config.Variable{
		"user":     config.StringVariable(user),
		"password": config.StringVariable(password),
		"account":  config.StringVariable(account),
		"role":     config.StringVariable(role),
	}
	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				ConfigDirectory:          config.TestNameDirectory(),
				ConfigVariables:          configVars,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCurrentAccount(t, account),
				),
				PlanOnly: true,
			},
		},
	})
}

func testAccCheckCurrentAccount(t *testing.T, n string) resource.TestCheckFunc {
	t.Helper()
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find resource: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("account Id resource ID not set.")
		}

		if rs.Primary.Attributes["account"] == "" {
			return fmt.Errorf("account expected to not be nil")
		}

		if rs.Primary.Attributes["region"] == "" {
			return fmt.Errorf("region expected to not be nil")
		}

		if rs.Primary.Attributes["url"] == "" {
			return fmt.Errorf("url expected to not be nil")
		}

		return nil
	}
}
