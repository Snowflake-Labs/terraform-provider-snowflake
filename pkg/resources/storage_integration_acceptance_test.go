package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStorageIntegration_validation(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      storageIntegrationConfig(name, []string{}, false),
				ExpectError: regexp.MustCompile("Not enough list items"),
			},
		},
	})
}

func TestAccStorageIntegration_aws(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: storageIntegrationConfig(name, []string{"s3://foo/"}, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.i", "name", name),
					resource.TestCheckNoResourceAttr("snowflake_storage_integration.i", "storage_aws_object_acl"),
				),
			},
			{
				Config: storageIntegrationConfig(name, []string{"s3://foo/"}, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_storage_integration.i", "name", name),
					resource.TestCheckResourceAttr("snowflake_storage_integration.i", "storage_aws_object_acl", "bucket-owner-full-control"),
				),
			},
		},
	})
}

func storageIntegrationConfig(name string, locations []string, awsObjectACL bool) string {
	awsObjectACLConfig := ""
	if awsObjectACL {
		awsObjectACLConfig = "storage_aws_object_acl = \"bucket-owner-full-control\""
	}
	return fmt.Sprintf(`
resource snowflake_storage_integration i {
	name = "%s"
	storage_allowed_locations = %q
	storage_provider = "S3"

	storage_aws_role_arn = "arn:aws:iam::000000000001:/role/test"
	%s
}
`, name, locations, awsObjectACLConfig)
}
