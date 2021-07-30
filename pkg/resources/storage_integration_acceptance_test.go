package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStorageIntegration_validation(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config:      storageIntegrationConfig(name, []string{}),
				ExpectError: regexp.MustCompile("1 item minimum, but config has only 0 declared"),
			},
		},
	})
}

func TestAccStorageIntegration_aws(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: storageIntegrationConfig(name, []string{"s3://foo/"}),
			},
		},
	})
}

func storageIntegrationConfig(name string, locations []string) string {
	return fmt.Sprintf(`
resource snowflake_storage_integration i {
	name = "%s"
	storage_allowed_locations = %q
	storage_provider = "S3"

	storage_aws_role_arn = "arn:aws:iam::000000000001:/role/test"
}
`, name, locations)
}
