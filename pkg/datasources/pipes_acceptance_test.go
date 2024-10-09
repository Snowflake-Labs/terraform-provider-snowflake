package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Pipes(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	pipeName := acc.TestClient().Ids.Alpha()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: pipes(databaseName, schemaName, pipeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_pipes.t", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_pipes.t", "schema", schemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_pipes.t", "pipes.#"),
					resource.TestCheckResourceAttr("data.snowflake_pipes.t", "pipes.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_pipes.t", "pipes.0.name", pipeName),
				),
			},
		},
	})
}

func pipes(databaseName string, schemaName string, pipeName string) string {
	s := `
resource "snowflake_database" "test" {
  name 	  = "%v"
  comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test" {
  name 	   = "%v"
  database = snowflake_database.test.name
  comment  = "Terraform acceptance test"
}

resource "snowflake_table" "test" {
  database = snowflake_database.test.name
  schema   = snowflake_schema.test.name
  name     = snowflake_schema.test.name
  column {
	name = "id"
	type = "NUMBER(5,0)"
  }
  column {
    name = "data"
	type = "VARCHAR(16)"
  }
}

resource "snowflake_stage" "test" {
  name = snowflake_schema.test.name
  database = snowflake_database.test.name
  schema = snowflake_schema.test.name
  comment = "Terraform acceptance test"
}

data snowflake_pipes "t" {
	database = snowflake_pipe.test.database
	schema = snowflake_pipe.test.schema
	depends_on = [snowflake_pipe.test]
}

resource "snowflake_pipe" "test" {
  database       = snowflake_database.test.name
  schema         = snowflake_schema.test.name
  name           = "%v"
  comment        = "Terraform acceptance test"
  copy_statement = <<CMD
COPY INTO "${snowflake_table.test.database}"."${snowflake_table.test.schema}"."${snowflake_table.test.name}"
  FROM @"${snowflake_stage.test.database}"."${snowflake_stage.test.schema}"."${snowflake_stage.test.name}"
  FILE_FORMAT = (TYPE = CSV)
CMD
  auto_ingest    = false
}
`
	return fmt.Sprintf(s, databaseName, schemaName, pipeName)
}
