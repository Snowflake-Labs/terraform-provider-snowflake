package resources_test

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAcc_FileFormatGrantFutureGrant(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: fileFormatGrantConfigFuture(t, databaseName, schemaName, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "database_name", databaseName),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "schema_name", schemaName),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "file_format_name", ""),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_file_format_grant.test", "privilege", "USAGE"),
				),
			},
		},
	})
}

func fileFormatGrantConfigFuture(t *testing.T, database_name, schema_name, role string) string {
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

resource "snowflake_file_format_grant" "test" {
    database_name = snowflake_database.test.name	
	roles         = [snowflake_role.test.name]
	schema_name   = snowflake_schema.test.name
	on_future = true
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
