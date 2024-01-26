package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_StageAlterWhenBothURLAndStorageIntegrationChange(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
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
	if !hasExternalEnvironmentVariablesSet {
		t.Skip("Skipping TestAcc_Stages_CreateOnS3 because external environment variables are not set")
	}

	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	url := "s3://foo/"
	comment := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	storageIntegration := ""
	credentials := fmt.Sprintf("AWS_KEY_ID = '%s' AWS_SECRET_KEY = '%s'", awsKeyId, awsSecretKey)
	encryption := "TYPE = 'NONE'"
	fileFormat := "TYPE = JSON NULL_IF = []"

	changedUrl := "s3://bar/"
	changedStorageIntegration := "s3_storage_integration"
	changedCredentials := ""
	changedEncryption := "TYPE = 'AWS_SSE_S3'"
	changedFileFormat := "TYPE = CSV"
	changedComment := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	configVariables := func(url string, storageIntegration string, credentials string, encryption string, fileFormat string, comment string) config.Variables {
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
		}
	}

	resourceName := "snowflake_stage.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables(url, storageIntegration, credentials, encryption, fileFormat, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", schemaName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "storage_integration", storageIntegration),
					resource.TestCheckResourceAttr(resourceName, "credentials", credentials),
					resource.TestCheckResourceAttr(resourceName, "encryption", encryption),
					resource.TestCheckResourceAttr(resourceName, "file_format", fileFormat),
					resource.TestCheckResourceAttr(resourceName, "url", url),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables(changedUrl, changedStorageIntegration, changedCredentials, changedEncryption, changedFileFormat, changedComment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database", databaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", schemaName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "storage_integration", changedStorageIntegration),
					resource.TestCheckResourceAttr(resourceName, "credentials", changedCredentials),
					resource.TestCheckResourceAttr(resourceName, "encryption", changedEncryption),
					resource.TestCheckResourceAttr(resourceName, "file_format", changedFileFormat),
					resource.TestCheckResourceAttr(resourceName, "url", changedUrl),
					resource.TestCheckResourceAttr(resourceName, "comment", changedComment),
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
