package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_StageAlterWhenBothURLAndStorageIntegrationChange(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: stageIntegrationConfig(name, "si1", "s3://foo/", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stage.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_stage.test", "url", "s3://foo/"),
				),
				Destroy: false,
			},
			{
				Config: stageIntegrationConfig(name, "changed", "s3://changed/", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stage.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_stage.test", "url", "s3://changed/"),
				),
			},
		},
	})
}

func stageIntegrationConfig(name string, siNameSuffix string, url string, databaseName string, schemaName string) string {
	resources := `
resource "snowflake_storage_integration" "test" {
	name = "%s%s"
	storage_allowed_locations = ["%s"]
	storage_provider = "S3"

  	storage_aws_role_arn = "arn:aws:iam::000000000001:/role/test"
}

resource "snowflake_stage" "test" {
	name = "%s"
	url = "%s"
	storage_integration = snowflake_storage_integration.test.name
	database = "%s"
	schema = "%s"
}
`

	return fmt.Sprintf(resources, name, siNameSuffix, url, name, url, databaseName, schemaName)
}
