package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_StageAlterWhenBothURLAndStorageIntegrationChange(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.Stage),
		Steps: []resource.TestStep{
			{
				Config: stageIntegrationConfig(name, "si1", "s3://foo/", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stage.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_stage.test", "url", "s3://foo/"),
				),
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

func TestAcc_Stage_CreateAndAlter(t *testing.T) {
	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	name := acc.TestClient().Ids.Alpha()
	url := "s3://foo/"
	comment := random.Comment()
	initialStorageIntegration := ""
	credentials := fmt.Sprintf("AWS_KEY_ID = '%s' AWS_SECRET_KEY = '%s'", awsKeyId, awsSecretKey)
	encryption := "TYPE = 'NONE'"
	copyOptionsWithQuotes := "ON_ERROR = 'CONTINUE'"

	changedUrl := awsBucketUrl + "/some-path"
	changedStorageIntegration := "S3_STORAGE_INTEGRATION"
	changedEncryption := "TYPE = 'AWS_SSE_S3'"
	changedFileFormat := "TYPE = JSON NULL_IF = []"
	changedComment := random.Comment()
	copyOptionsWithoutQuotes := "ON_ERROR = CONTINUE"

	configVariables := func(url string, storageIntegration string, credentials string, encryption string, fileFormat string, comment string, copyOptions string) config.Variables {
		return config.Variables{
			"database":            config.StringVariable(databaseName),
			"schema":              config.StringVariable(schemaName),
			"name":                config.StringVariable(name),
			"url":                 config.StringVariable(url),
			"storage_integration": config.StringVariable(storageIntegration),
			"credentials":         config.StringVariable(credentials),
			"encryption":          config.StringVariable(encryption),
			"file_format":         config.StringVariable(fileFormat),
			"comment":             config.StringVariable(comment),
			"copy_options":        config.StringVariable(copyOptions),
		}
	}

	resourceName := "snowflake_stage.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Stage),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables(url, initialStorageIntegration, credentials, encryption, "", comment, copyOptionsWithQuotes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", schemaName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "storage_integration", initialStorageIntegration),
					resource.TestCheckResourceAttr(resourceName, "credentials", credentials),
					resource.TestCheckResourceAttr(resourceName, "encryption", encryption),
					resource.TestCheckResourceAttr(resourceName, "file_format", ""),
					resource.TestCheckResourceAttr(resourceName, "copy_options", copyOptionsWithoutQuotes),
					resource.TestCheckResourceAttr(resourceName, "url", url),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables(changedUrl, changedStorageIntegration, credentials, changedEncryption, changedFileFormat, changedComment, copyOptionsWithoutQuotes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", schemaName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "storage_integration", changedStorageIntegration),
					resource.TestCheckResourceAttr(resourceName, "credentials", credentials),
					resource.TestCheckResourceAttr(resourceName, "encryption", changedEncryption),
					resource.TestCheckResourceAttr(resourceName, "file_format", changedFileFormat),
					resource.TestCheckResourceAttr(resourceName, "copy_options", copyOptionsWithoutQuotes),
					resource.TestCheckResourceAttr(resourceName, "url", changedUrl),
					resource.TestCheckResourceAttr(resourceName, "comment", changedComment),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables(changedUrl, changedStorageIntegration, credentials, changedEncryption, changedFileFormat, changedComment, copyOptionsWithoutQuotes),
				Destroy:         true,
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables(url, initialStorageIntegration, credentials, encryption, "", comment, copyOptionsWithoutQuotes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", schemaName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "storage_integration", initialStorageIntegration),
					resource.TestCheckResourceAttr(resourceName, "credentials", credentials),
					resource.TestCheckResourceAttr(resourceName, "encryption", encryption),
					resource.TestCheckResourceAttr(resourceName, "file_format", ""),
					resource.TestCheckResourceAttr(resourceName, "copy_options", copyOptionsWithoutQuotes),
					resource.TestCheckResourceAttr(resourceName, "url", url),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

func stageIntegrationConfig(name string, siNameSuffix string, url string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_storage_integration" "test" {
	name = "%[1]s%[2]s"
	storage_allowed_locations = ["%[3]s"]
	storage_provider = "S3"

  	storage_aws_role_arn = "arn:aws:iam::000000000001:/role/test"
}

resource "snowflake_stage" "test" {
	name = "%[1]s"
	url = "%[3]s"
	storage_integration = snowflake_storage_integration.test.name
	database = "%[4]s"
	schema = "%[5]s"
}
`, name, siNameSuffix, url, databaseName, schemaName)
}
