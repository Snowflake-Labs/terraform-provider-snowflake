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

func TestAccStreamGrant_basic(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	streamName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: streamGrantConfigExisting(t, databaseName, schemaName, roleName, streamName, tableName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "database_name", databaseName),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "schema_name", schemaName),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "stream_name", streamName),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "on_future", "false"),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "privilege", "SELECT"),
				),
			},
		},
	})
}

func TestAccStreamGrante_future(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: streamGrantConfigFuture(t, databaseName, schemaName, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "database_name", databaseName),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "schema_name", schemaName),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "stream_name", ""),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "privilege", "SELECT"),
				),
			},
		},
	})
}

func streamGrantConfigExisting(t *testing.T, database_name, schema_name, role, stream_name, table_name string) string {
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
resource "snowflake_table" "test" {
	database        = snowflake_database.test.name
	schema          = snowflake_schema.test.name
	name            = "{{ .table_name }}"
	change_tracking = true
	comment         = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16777216)"
	}
}

resource "snowflake_stream" "test" {
	database = snowflake_database.test.name
	schema   = snowflake_schema.test.name
	name     = "{{ .stream_name }}"
	comment  = "Terraform acceptance test"
	on_table = "${snowflake_database.test.name}.${snowflake_schema.test.name}.${snowflake_table.test.name}"
}

resource "snowflake_stream_grant" "test" {
  database_name = snowflake_database.test.name
	roles         = [snowflake_role.test.name]
	schema_name   = snowflake_schema.test.name
	stream_name = snowflake_stream.test.name
	privilege = "SELECT"
}
`

	out := bytes.NewBuffer(nil)
	tmpl := template.Must(template.New("view)").Parse(config))
	err := tmpl.Execute(out, map[string]string{
		"database_name": database_name,
		"schema_name":   schema_name,
		"role_name":     role,
		"stream_name":   stream_name,
		"table_name":    table_name,
	})
	r.NoError(err)

	return out.String()
}

func streamGrantConfigFuture(t *testing.T, database_name, schema_name, role string) string {
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

resource "snowflake_stream_grant" "test" {
  database_name = snowflake_database.test.name
	roles         = [snowflake_role.test.name]
	schema_name   = snowflake_schema.test.name
	on_future = true
	depends_on = [snowflake_role.test]
	privilege = "SELECT"
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
