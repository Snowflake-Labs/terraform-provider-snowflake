package resources_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAccStageGrant_defaults(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: stageGrantConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database.d", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.s", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.s", "database", name),
					resource.TestCheckResourceAttr("snowflake_role.r", "name", name),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "schema_name", name),
					testRolesAndShares(t, "snowflake_stage_grant.g", []string{name}, []string{}),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_stage_grant.g",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func stageGrantConfig(n string) string {
	return fmt.Sprintf(`
	resource snowflake_database d {
		name = "%v"
		comment = "Terraform acceptance test"
	}
	
	resource snowflake_schema s {
		name = "%v"
		database = snowflake_database.d.name
		comment = "Terraform acceptance test"
	}
	
	resource snowflake_stage s {
		name = "%v"
		database = snowflake_database.d.name
		schema = snowflake_schema.s.name
		comment = "Terraform acceptance test"
	}

	resource snowflake_role r {
		name = "%s"
	}

	resource snowflake_stage_grant g {
		database_name = snowflake_database.d.name
		schema_name = snowflake_schema.s.name
		stage_name = snowflake_stage.s.name

		privilege = "READ"

		roles = [
			snowflake_role.r.name
		]
	}
`, n, n, n, n)
}

func TestAcc_StageFutureGrant(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: stageGrantConfigFuture(t, databaseName, schemaName, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stage_grant.test", "database_name", databaseName),
					resource.TestCheckResourceAttr("snowflake_stage_grant.test", "schema_name", schemaName),
					resource.TestCheckResourceAttr("snowflake_stage_grant.test", "stage_name", ""),
					resource.TestCheckResourceAttr("snowflake_stage_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_stage_grant.test", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_stage_grant.test", "privilege", "USAGE"),
				),
			},
		},
	})
}

func stageGrantConfigFuture(t *testing.T, database_name, schema_name, role string) string {
	r := require.New(t)

	config := `
resource "snowflake_database" "test" {
  name = "{{ .database_name }}"
}

resource "snowflake_schema" "test" {
	name = "{{ .schema_name }}"
	database = snowflake_database.test.name
}

resource "snowflake_role" "test" {
  name = "{{.role_name}}"
}

resource "snowflake_stage_grant" "test" {
    database_name = snowflake_database.test.name	
	roles         = ["{{.role_name}}"]
	schema_name   = snowflake_schema.test.name
	on_future = true
	depends_on = [snowflake_role.test]
	privilege = "USAGE"
}
`

	out := bytes.NewBuffer(nil)
	tmpl := template.Must(template.New("view)").Parse(config))
	err := tmpl.Execute(out, map[string]string{
		"database_name": database_name,
		"schema_name":   schema_name,
		"role_name":     role,
	})
	r.NoError(err)

	return out.String()
}
