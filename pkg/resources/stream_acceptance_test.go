package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_StreamCreateOnStageWithoutDirectoryEnabled(t *testing.T) {
	accName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.Stream),
		Steps: []resource.TestStep{
			{
				Config:      stageStreamConfig(accName, false),
				ExpectError: regexp.MustCompile("directory must be enabled on stage"),
			},
		},
	})
}

func TestAcc_StreamCreateOnStage(t *testing.T) {
	accName := acc.TestClient().Ids.Alpha()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.Stream),
		Steps: []resource.TestStep{
			{
				Config: stageStreamConfig(accName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
					checkBool("snowflake_stream.test_stream", "append_only", false),
					checkBool("snowflake_stream.test_stream", "insert_only", false),
					checkBool("snowflake_stream.test_stream", "show_initial_rows", false),
				),
			},
		},
	})
}

// proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2672
func TestAcc_Stream_OnTable(t *testing.T) {
	tableName := acc.TestClient().Ids.Alpha()
	tableName2 := acc.TestClient().Ids.Alpha()
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.Stream),
		Steps: []resource.TestStep{
			{
				Config: streamConfigOnTable(acc.TestDatabaseName, acc.TestSchemaName, tableName, id.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_table", fmt.Sprintf("\"%s\".\"%s\".%s", acc.TestDatabaseName, acc.TestSchemaName, tableName)),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				Config: streamConfigOnTable(acc.TestDatabaseName, acc.TestSchemaName, tableName2, id.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_table", fmt.Sprintf("\"%s\".\"%s\".%s", acc.TestDatabaseName, acc.TestSchemaName, tableName2)),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// proves issue https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2672
func TestAcc_Stream_OnView(t *testing.T) {
	// TODO(SNOW-1423486): Fix using warehouse in all tests and remove unsetting testenvs.ConfigureClientOnce
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	tableName := acc.TestClient().Ids.Alpha()
	viewName := acc.TestClient().Ids.Alpha()
	name := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.Stream),
		Steps: []resource.TestStep{
			{
				Config: streamConfigOnView(acc.TestDatabaseName, acc.TestSchemaName, tableName, viewName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", name),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_view", fmt.Sprintf("\"%s\".\"%s\".%s", acc.TestDatabaseName, acc.TestSchemaName, viewName)),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

func TestAcc_Stream(t *testing.T) {
	// Current error is User: <redacted> is not authorized to perform: sts:AssumeRole on resource: <redacted> duration 1.162414333s args {}] ()
	t.Skip("Skipping TestAcc_Stream")

	accName := acc.TestClient().Ids.Alpha()
	accNameExternalTable := acc.TestClient().Ids.Alpha()
	bucketURL := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	roleName := testenvs.GetOrSkipTest(t, testenvs.AwsExternalRoleArn)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Stream),
		Steps: []resource.TestStep{
			{
				Config: streamConfig(accName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_table", fmt.Sprintf("%s.%s.%s", accName, accName, "STREAM_ON_TABLE")),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
					checkBool("snowflake_stream.test_stream", "append_only", false),
					checkBool("snowflake_stream.test_stream", "insert_only", false),
					checkBool("snowflake_stream.test_stream", "show_initial_rows", false),
				),
			},
			{
				Config: streamConfig(accName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_table", fmt.Sprintf("%s.%s.%s", accName, accName, "STREAM_ON_TABLE")),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
					checkBool("snowflake_stream.test_stream", "append_only", true),
					checkBool("snowflake_stream.test_stream", "insert_only", false),
					checkBool("snowflake_stream.test_stream", "show_initial_rows", false),
				),
			},
			{
				Config: externalTableStreamConfig(accNameExternalTable, false, bucketURL, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", accNameExternalTable),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", accNameExternalTable),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", accNameExternalTable),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_table", fmt.Sprintf("%s.%s.%s", accNameExternalTable, accNameExternalTable, "STREAM_ON_EXTERNAL_TABLE")),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
					checkBool("snowflake_stream.test_stream", "append_only", false),
					checkBool("snowflake_stream.test_stream", "insert_only", false),
					checkBool("snowflake_stream.test_stream", "show_initial_rows", false),
				),
			},
			{
				Config: externalTableStreamConfig(accNameExternalTable, true, bucketURL, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", accNameExternalTable),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", accNameExternalTable),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", accNameExternalTable),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_table", fmt.Sprintf("%s.%s.%s", accNameExternalTable, accNameExternalTable, "STREAM_ON_EXTERNAL_TABLE")),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
					checkBool("snowflake_stream.test_stream", "append_only", false),
					checkBool("snowflake_stream.test_stream", "insert_only", true),
					checkBool("snowflake_stream.test_stream", "show_initial_rows", false),
				),
			},
			{
				Config: viewStreamConfig(accName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_view", fmt.Sprintf("%s.%s.%s", accName, accName, "STREAM_ON_VIEW")),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
					checkBool("snowflake_stream.test_stream", "append_only", false),
					checkBool("snowflake_stream.test_stream", "insert_only", false),
					checkBool("snowflake_stream.test_stream", "show_initial_rows", false),
				),
			},
			{
				Config: stageStreamConfig(accName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "name", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "database", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "on_stage", fmt.Sprintf("%s.%s.%s", accName, accName, "STREAM_ON_STAGE")),
					resource.TestCheckResourceAttr("snowflake_stream.test_stream", "comment", "Terraform acceptance test"),
				),
			},
			{
				ResourceName:      "snowflake_stream.test_stream",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func streamConfigOnTable(databaseName string, schemaName string, tableName string, name string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_stream_on_table" {
	database        = "%[1]s"
	schema          = "%[2]s"
	name            = "%[3]s"
	comment         = "Terraform acceptance test"
	change_tracking = true

	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR"
	}
}

resource "snowflake_stream" "test_stream" {
	database    = "%[1]s"
	schema      = "%[2]s"
	name        = "%[4]s"
	comment     = "Terraform acceptance test"
	on_table    = "\"%[1]s\".\"%[2]s\".\"${snowflake_table.test_stream_on_table.name}\""
}
`, databaseName, schemaName, tableName, name)
}

func streamConfigOnView(databaseName string, schemaName string, tableName string, viewName string, name string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test" {
	database        = "%[1]s"
	schema          = "%[2]s"
	name            = "%[3]s"
	comment         = "Terraform acceptance test"
	change_tracking = true

	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR"
	}
}

resource "snowflake_view" "test" {
	database = "%[1]s"
	schema   = "%[2]s"
	name     = "%[4]s"
	change_tracking = true

	statement = "select * from \"${snowflake_table.test.name}\""
}

resource "snowflake_stream" "test_stream" {
	database    = "%[1]s"
	schema      = "%[2]s"
	name        = "%[5]s"
	comment     = "Terraform acceptance test"
	on_view     = "\"%[1]s\".\"%[2]s\".\"${snowflake_view.test.name}\""
}
`, databaseName, schemaName, tableName, viewName, name)
}

func streamConfig(name string, appendOnly bool) string {
	appendOnlyConfig := ""
	if appendOnly {
		appendOnlyConfig = "append_only = true"
	}

	s := `
resource "snowflake_database" "test_database" {
	name    = "%s"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test_schema" {
	name     = "%s"
	database = snowflake_database.test_database.name
	comment  = "Terraform acceptance test"
}

resource "snowflake_table" "test_stream_on_table" {
	database        = snowflake_database.test_database.name
	schema          = snowflake_schema.test_schema.name
	name            = "STREAM_ON_TABLE"
	comment         = "Terraform acceptance test"
	change_tracking = true

	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR"
	}
}

resource "snowflake_stream" "test_stream" {
	database    = snowflake_database.test_database.name
	schema      = snowflake_schema.test_schema.name
	name        = "%s"
	comment     = "Terraform acceptance test"
	on_table    = "${snowflake_database.test_database.name}.${snowflake_schema.test_schema.name}.${snowflake_table.test_stream_on_table.name}"
	%s
}
`
	return fmt.Sprintf(s, name, name, name, appendOnlyConfig)
}

func externalTableStreamConfig(name string, insertOnly bool, bucketURL string, roleName string) string {
	// Refer to external_table_acceptance_test.go for the original source on
	// external table resources and dependents (modified slightly here).
	insertOnlyConfig := ""
	if insertOnly {
		insertOnlyConfig = "insert_only = true"
	}

	s := `
resource "snowflake_database" "test" {
	name = "%v"
	comment = "Terraform acceptance test"
}
resource "snowflake_schema" "test" {
	name = "%v"
	database = snowflake_database.test.name
	comment = "Terraform acceptance test"
}
resource "snowflake_stage" "test" {
	name = "%v"
	url = "%s"
	database = snowflake_database.test.name
	schema = snowflake_schema.test.name
	comment = "Terraform acceptance test"
	storage_integration = snowflake_storage_integration.external_table_stream_integration.name
}
resource "snowflake_storage_integration" "external_table_stream_integration" {
	name = "%v"
	storage_allowed_locations = ["%s"]
	storage_provider = "S3"
	storage_aws_role_arn = "%s"
}
resource "snowflake_external_table" "test_external_stream_table" {
	database = snowflake_database.test.name
	schema   = snowflake_schema.test.name
	name     = "STREAM_ON_EXTERNAL_TABLE"
	comment  = "Terraform acceptance test"
	column {
		name = "column1"
		type = "STRING"
		as   = "TO_VARCHAR(TO_TIMESTAMP_NTZ(value:unix_timestamp_property::NUMBER, 3), 'yyyy-mm-dd-hh')"
	}
	column {
		name = "column2"
		type = "TIMESTAMP_NTZ(9)"
		as   = "($1:\"CreatedDate\"::timestamp)"
	}
  file_format = "TYPE = CSV"
  location = "@${snowflake_database.test.name}.${snowflake_schema.test.name}.${snowflake_stage.test.name}"
}
resource "snowflake_stream" "test_external_table_stream" {
	database = snowflake_database.test.name
	schema   = snowflake_schema.test.name
	name     = "%s"
	comment  = "Terraform acceptance test"
	on_table = "${snowflake_database.test.name}.${snowflake_schema.test.name}.${snowflake_external_table.test_external_stream_table.name}"
	%s
}
`

	return fmt.Sprintf(s, name, name, name, bucketURL, name, bucketURL, roleName, name, insertOnlyConfig)
}

func viewStreamConfig(name string, appendOnly bool) string {
	appendOnlyConfig := ""
	if appendOnly {
		appendOnlyConfig = "append_only = true"
	}

	s := `
resource "snowflake_database" "test_database" {
	name    = "%s"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test_schema" {
	name     = "%s"
	database = snowflake_database.test_database.name
	comment  = "Terraform acceptance test"
}

resource "snowflake_table" "test_stream_on_view" {
	database        = snowflake_database.test_database.name
	schema          = snowflake_schema.test_schema.name
	name            = "STREAM_ON_VIEW_TABLE"
	comment         = "Terraform acceptance test"
	change_tracking = true

	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16777216)"
	}
}

resource "snowflake_view" "test_stream_on_view" {
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	name     = "STREAM_ON_VIEW"

	statement = "select * from ${snowflake_table.test_stream_on_view.name}"
}

resource "snowflake_stream" "test_stream" {
	database    = snowflake_database.test_database.name
	schema      = snowflake_schema.test_schema.name
	name        = "%s"
	comment     = "Terraform acceptance test"
	on_view    = "${snowflake_database.test_database.name}.${snowflake_schema.test_schema.name}.${snowflake_view.test_stream_on_view.name}"
	%s
}
`
	return fmt.Sprintf(s, name, name, name, appendOnlyConfig)
}

func stageStreamConfig(name string, directory bool) string {
	s := `
resource "snowflake_database" "test_database" {
	name    = "%s"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test_schema" {
	name     = "%s"
	database = snowflake_database.test_database.name
	comment  = "Terraform acceptance test"
}

resource "snowflake_stage" "test_stage" {
	name	 = "%s"
	database = snowflake_database.test_database.name
	schema	 = snowflake_schema.test_schema.name
	directory = "ENABLE = %t"
}

resource "snowflake_stream" "test_stream" {
	database    = snowflake_database.test_database.name
	schema      = snowflake_schema.test_schema.name
	name        = "%s"
	comment     = "Terraform acceptance test"
	on_stage    = "${snowflake_database.test_database.name}.${snowflake_schema.test_schema.name}.${snowflake_stage.test_stage.name}"
}
`
	return fmt.Sprintf(s, name, name, name, directory, name)
}
