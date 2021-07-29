package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStorageIntegrations(t *testing.T) {
	storageIntegrationName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: storageIntegrations(storageIntegrationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_storage_integrations.s", "storage_integrations.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_storage_integrations.s", "storage_integrations.0.name"),
				),
			},
		},
	})
}

func storageIntegrations(storageIntegrationName string) string {
	return fmt.Sprintf(`
	
	resource snowflake_storage_integration i {
		name = "%v"
		storage_allowed_locations = ["s3://foo/"]
		storage_provider = "S3"
		storage_aws_role_arn = "arn:aws:iam::000000000001:/role/test"
	}

	data snowflake_storage_integrations "s" {
		depends_on = [snowflake_storage_integration.i]
	}
	`, storageIntegrationName)
}
