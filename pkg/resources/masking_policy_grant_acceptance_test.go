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

func TestAcc_MaskingPolicyGrant(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	maskingPolicyName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: maskingPolicyGrantConfig(t, databaseName, schemaName, maskingPolicyName, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_making_policy_grant.test", "database_name", databaseName),
					resource.TestCheckResourceAttr("snowflake_making_policy_grant.test", "schema_name", schemaName),
					resource.TestCheckResourceAttr("snowflake_making_policy_grant.test", "making_policy_name", maskingPolicyName),
					resource.TestCheckResourceAttr("snowflake_making_policy_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_making_policy_grant.test", "privilege", "APPLY"),
				),
			},
		},
	})
}

func maskingPolicyGrantConfig(t *testing.T, database_name, schema_name, masking_policy_name, role string) string {
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

resource "snowflake_masking_policy" "test" {
	name = "{{.masking_policy_name}}"
	database = snowflake_database.test.name
	schema = snowflake_schema.test.name
	value_data_type = "VARCHAR"
	masking_expression = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
	return_data_type = "VARCHAR(16777216)"
	comment = "Terraform acceptance test"
}

resource "snowflake_masking_policy_grant" "test" {
	masking_policy_name = snowflake_masking_policy.test.name
    database_name = snowflake_database.test.name
	roles         = [snowflake_role.test.name]
	schema_name   = snowflake_schema.test.name
	privilege = "APPLY"
}
`

	out := bytes.NewBuffer(nil)
	tmpl := template.Must(template.New("view)").Parse(config))
	err := tmpl.Execute(out, map[string]string{
		"database_name":       database_name,
		"schema_name":         schema_name,
		"masking_policy_name": masking_policy_name,
		"role_name":           role,
	})
	r.NoError(err)

	return out.String()
}
