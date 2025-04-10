package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/ids"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_StageAlterWhenBothURLAndStorageIntegrationChange(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	roleArn := "arn:aws:iam::000000000001:/role/test"
	baseUrl := "s3://foo"
	firstUrl := baseUrl + "/allowed-location"
	secondUrl := baseUrl + "/allowed-location2"

	storageIntegration, storageIntegrationCleanup := acc.TestClient().StorageIntegration.CreateS3(t, baseUrl, roleArn)
	t.Cleanup(storageIntegrationCleanup)

	stageId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.Stage),
		Steps: []resource.TestStep{
			{
				Config: stageIntegrationConfig(stageId, storageIntegration.ID(), firstUrl),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stage.test", "name", stageId.Name()),
					resource.TestCheckResourceAttr("snowflake_stage.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_stage.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_stage.test", "fully_qualified_name", stageId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_stage.test", "url", firstUrl),
				),
			},
			{
				Config: stageIntegrationConfig(stageId, storageIntegration.ID(), secondUrl),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stage.test", "name", stageId.Name()),
					resource.TestCheckResourceAttr("snowflake_stage.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_stage.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_stage.test", "fully_qualified_name", stageId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_stage.test", "url", secondUrl),
				),
			},
		},
	})
}

func stageIntegrationConfig(stageId sdk.SchemaObjectIdentifier, storageIntegrationId sdk.AccountObjectIdentifier, url string) string {
	return fmt.Sprintf(`
resource "snowflake_stage" "test" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
	storage_integration = "%[4]s"
	url = "%[5]s"
}
`, stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(), storageIntegrationId.Name(), url)
}

func TestAcc_Stage_CreateAndAlter(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	awsBucketUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	url := "s3://foo/"
	comment := random.Comment()
	initialStorageIntegration := ""
	credentials := fmt.Sprintf("AWS_KEY_ID = '%s' AWS_SECRET_KEY = '%s'", awsKeyId, awsSecretKey)
	encryption := "TYPE = 'NONE'"
	copyOptionsWithQuotes := "ON_ERROR = 'CONTINUE'"

	changedUrl := awsBucketUrl + "/some-path"
	changedStorageIntegration := ids.PrecreatedS3StorageIntegration
	changedEncryption := "TYPE = 'AWS_SSE_S3'"
	changedFileFormat := "TYPE = JSON NULL_IF = []"
	changedFileFormatWithQuotes := "FIELD_DELIMITER = '|' PARSE_HEADER = true"
	changedFileFormatWithoutQuotes := "FIELD_DELIMITER = | PARSE_HEADER = true"
	changedComment := random.Comment()
	copyOptionsWithoutQuotes := "ON_ERROR = CONTINUE"

	configVariables := func(url string, storageIntegration string, credentials string, encryption string, fileFormat string, comment string, copyOptions string) config.Variables {
		return config.Variables{
			"database":            config.StringVariable(id.DatabaseName()),
			"schema":              config.StringVariable(id.SchemaName()),
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
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
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
				ConfigVariables: configVariables(changedUrl, changedStorageIntegration.Name(), credentials, changedEncryption, changedFileFormat, changedComment, copyOptionsWithoutQuotes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "storage_integration", changedStorageIntegration.Name()),
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
				ConfigVariables: configVariables(changedUrl, changedStorageIntegration.Name(), credentials, changedEncryption, changedFileFormatWithQuotes, changedComment, copyOptionsWithoutQuotes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "storage_integration", changedStorageIntegration.Name()),
					resource.TestCheckResourceAttr(resourceName, "credentials", credentials),
					resource.TestCheckResourceAttr(resourceName, "encryption", changedEncryption),
					resource.TestCheckResourceAttr(resourceName, "file_format", changedFileFormatWithoutQuotes),
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
				ConfigVariables: configVariables(changedUrl, changedStorageIntegration.Name(), credentials, changedEncryption, changedFileFormat, changedComment, copyOptionsWithoutQuotes),
				Destroy:         true,
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables(url, initialStorageIntegration, credentials, encryption, "", comment, copyOptionsWithoutQuotes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
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

func TestAcc_Stage_Issue2972(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	stageId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newId := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(stageId.SchemaId())

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
				Config: stageIssue2972Config(stageId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database", stageId.DatabaseName()),
					resource.TestCheckResourceAttr(resourceName, "schema", stageId.SchemaName()),
					resource.TestCheckResourceAttr(resourceName, "name", stageId.Name()),
				),
			},
			{
				PreConfig: func() {
					acc.TestClient().Stage.Rename(t, stageId, newId)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
					},
				},
				Config: stageIssue2972Config(stageId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "database", stageId.DatabaseName()),
					resource.TestCheckResourceAttr(resourceName, "schema", stageId.SchemaName()),
					resource.TestCheckResourceAttr(resourceName, "name", stageId.Name()),
				),
			},
		},
	})
}

func stageIssue2972Config(stageId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_stage" "test" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[3]s"
}
`, stageId.DatabaseName(), stageId.SchemaName(), stageId.Name())
}

// TODO [SNOW-1348110]: fix behavior with stage rework
func TestAcc_Stage_Issue2679(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	roleArn := "arn:aws:iam::000000000001:/role/test"
	baseUrl := "s3://foo"
	allowedUrl := baseUrl + "/allowed-location"

	storageIntegration, storageIntegrationCleanup := acc.TestClient().StorageIntegration.CreateS3(t, baseUrl, roleArn)
	t.Cleanup(storageIntegrationCleanup)

	storageIntegrationId := storageIntegration.ID()
	stageId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	fileFormatWithDefaultTypeCsv := "TYPE = CSV NULL_IF = []"
	fileFormatWithoutType := "NULL_IF = []"

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
				Config: stageIssue2679Config(stageId, storageIntegrationId, fileFormatWithDefaultTypeCsv, allowedUrl),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", stageId.Name()),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: stageIssue2679Config(stageId, storageIntegrationId, fileFormatWithoutType, allowedUrl),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", stageId.Name()),
					// TODO [SNOW-1348110]: use generated assertions after stage rework
					func(_ *terraform.State) error {
						properties, err := acc.TestClient().Stage.Describe(t, stageId)
						if err != nil {
							return err
						}
						typeProperty, err := collections.FindFirst(properties, func(property sdk.StageProperty) bool {
							return property.Parent == "STAGE_FILE_FORMAT" && property.Name == "TYPE"
						})
						if err != nil {
							return err
						}
						if typeProperty.Value != "CSV" {
							return fmt.Errorf("expected type property 'CSV', got '%s'", typeProperty.Value)
						}
						return nil
					},
				),
			},
		},
	})
}

func stageIssue2679Config(stageId sdk.SchemaObjectIdentifier, storageIntegrationId sdk.AccountObjectIdentifier, fileFormat string, url string) string {
	return fmt.Sprintf(`
resource "snowflake_stage" "test" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
	storage_integration = "%[4]s"
	file_format = "%[5]s"
	url = "%[6]s"
}
`, stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(), storageIntegrationId.Name(), fileFormat, url)
}
