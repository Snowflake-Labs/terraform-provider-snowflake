package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_ScimIntegration(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_SCIM_INTEGRATION_TESTS"); ok {
		t.Skip("Skipping TestAccScimIntegration")
	}

	scimIntName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	scimIntName2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: scimIntegrationConfig_azure(scimIntName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test_azure_int", "name", scimIntName),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test_azure_int", "scim_client", "AZURE"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test_azure_int", "run_as_role", "AAD_PROVISIONER"),
					resource.TestCheckResourceAttrSet("snowflake_scim_integration.test_azure_int", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_scim_integration.test_azure_int", "enabled"),
				),
			},
			{
				Config: scimIntegrationConfig_azure_np(scimIntName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_scim_integration.test_azure_int", "name", scimIntName2),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test_azure_int", "scim_client", "AZURE"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test_azure_int", "run_as_role", "AAD_PROVISIONER"),
					resource.TestCheckResourceAttr("snowflake_scim_integration.test_azure_int", "network_policy", "AAD_NETWORK_POLICY"),
					resource.TestCheckResourceAttrSet("snowflake_scim_integration.test_azure_int", "created_on"),
					resource.TestCheckResourceAttrSet("snowflake_scim_integration.test_azure_int", "enabled"),
				),
			},
		},
	})
}

func scimIntegrationConfig_azure(name string) string {
	return fmt.Sprintf(`
	resource "snowflake_scim_integration" "test_azure_int" {
		name = "%s"
		scim_client = "AZURE"
		run_as_role = "AAD_PROVISIONER"
		enabled = true
	}
	`, name)
}

func scimIntegrationConfig_azure_np(name string) string {
	return fmt.Sprintf(`
	resource "snowflake_scim_integration" "test_azure_int_np" {
		name = "%s"
		scim_client = "AZURE"
		run_as_role = "AAD_PROVISIONER"
		network_policy = "AAD_NETWORK_POLICY"
		enabled = true
	}
	`, name)
}
