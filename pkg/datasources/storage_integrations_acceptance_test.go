package datasources_test

import (
	"fmt"
	"strconv"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_StorageIntegrations_basic(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigVariables: config.Variables{
					"name": config.StringVariable(name),
					"allowed_locations": config.SetVariable(
						config.StringVariable("gcs://foo/"),
						config.StringVariable("gcs://bar/"),
					),
					"blocked_locations": config.SetVariable(
						config.StringVariable("gcs://foo/"),
						config.StringVariable("gcs://bar/"),
					),
					"comment": config.StringVariable("some comment"),
				},
				ConfigDirectory: config.TestNameDirectory(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_storage_integrations.test", "storage_integrations.#"),
					containsStorageIntegration(name, true, "some comment"),
				),
			},
		},
	})
}

func containsStorageIntegration(name string, enabled bool, comment string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "snowflake_storage_integrations" {
				continue
			}
			iter, err := strconv.ParseInt(rs.Primary.Attributes["storage_integrations.#"], 10, 32)
			if err != nil {
				return err
			}

			for i := 0; i < int(iter); i++ {
				if rs.Primary.Attributes[fmt.Sprintf("storage_integrations.%d.name", i)] == name {
					actualEnabled, err := strconv.ParseBool(rs.Primary.Attributes[fmt.Sprintf("storage_integrations.%d.enabled", i)])
					if err != nil {
						return err
					}

					if actualEnabled != enabled {
						return fmt.Errorf("expected comment: %v, but got: %v", enabled, actualEnabled)
					}

					actualComment := rs.Primary.Attributes[fmt.Sprintf("storage_integrations.%d.comment", i)]
					if actualComment != comment {
						return fmt.Errorf("expected comment: %s, but got: %s", comment, actualComment)
					}

					return nil
				}
			}

			return fmt.Errorf("storage integration (%s) not found", name)
		}
		return nil
	}
}
