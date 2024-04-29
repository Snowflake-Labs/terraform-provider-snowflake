package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Stages(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	storageIntegrationName := acc.TestClient().Ids.Alpha()
	stageName := acc.TestClient().Ids.Alpha()
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: stages(databaseName, schemaName, storageIntegrationName, stageName, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_stages.test", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_stages.test", "schema", schemaName),
					resource.TestCheckResourceAttr("data.snowflake_stages.test", "stages.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_stages.test", "stages.0.name", stageName),
					resource.TestCheckResourceAttr("data.snowflake_stages.test", "stages.0.storage_integration", storageIntegrationName),
					resource.TestCheckResourceAttr("data.snowflake_stages.test", "stages.0.comment", comment),
				),
			},
		},
	})
}

func stages(databaseName string, schemaName string, storageIntegrationName string, stageName string, comment string) string {
	return fmt.Sprintf(`
	resource "snowflake_database" "test" {
		name = "%s"
	}

	resource "snowflake_schema" "test"{
		name 	 = "%s"
		database = snowflake_database.test.name
	}

	resource "snowflake_storage_integration" "test" {
		name = "%s"
		storage_allowed_locations = ["s3://foo/"]
		storage_provider = "S3"

		storage_aws_role_arn = "arn:aws:iam::000000000001:/role/test"
	}

	resource "snowflake_stage" "test"{
		name 	 				= "%s"
		database 				= snowflake_schema.test.database
		schema 	 				= snowflake_schema.test.name
		url 					= "s3://foo/"
		storage_integration 	= snowflake_storage_integration.test.name
		comment  				= "%s"
	}

	data "snowflake_stages" "test" {
		depends_on = [snowflake_storage_integration.test, snowflake_stage.test]

		database = snowflake_stage.test.database
		schema = snowflake_stage.test.schema
	}
	`, databaseName, schemaName, storageIntegrationName, stageName, comment)
}
