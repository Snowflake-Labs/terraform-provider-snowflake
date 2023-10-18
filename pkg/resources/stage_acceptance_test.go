package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlterStageWhenBothURLAndStorageIntegrationChange(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: stageIntegrationConfig(name, "si1", "s3://foo/", "ff1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stage.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_stage.test", "url", "s3://foo/"),
				),
				Destroy: false,
			},
			{
				Config: stageIntegrationConfig(name, "changed", "s3://changed/", "ff2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stage.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_stage.test", "url", "s3://changed/"),
				),
			},
		},
	})
}

func stageIntegrationConfig(name string, siNameSuffix string, url string, ffName string) string {
	resources := `
resource "snowflake_database" "test" {
	name = "%s"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test" {
	name = "%s"
	database = snowflake_database.test.name
	comment = "Terraform acceptance test"
}

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
	schema = snowflake_schema.test.name
	database = snowflake_database.test.name
	file_format = snowflake_file_format.test.name
}

resource "snowflake_file_format" "test" {
	database             = snowflake_database.test.name
	schema               = snowflake_schema.test.name
	name                 = "%s"
	format_type          = "JSON"
	compression          = "GZIP"
	skip_byte_order_mark = true
}
`

	return fmt.Sprintf(resources, name, name, name, siNameSuffix, url, name, url, ffName)
}
