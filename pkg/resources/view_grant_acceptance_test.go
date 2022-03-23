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

func TestAcc_ViewGrantBasic(t *testing.T) {
	viewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: viewGrantConfigFuture(t, databaseName, viewName, roleName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "view_name", viewName),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "privilege", "SELECT"),
				),
			},
		},
	})
}

func TestAcc_ViewGrantShares(t *testing.T) {

	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	viewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	shareName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: viewGrantConfigShares(t, databaseName, viewName, roleName, shareName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "view_name", viewName),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "privilege", "SELECT"),
				),
			},
		},
	})
}

func TestAcc_FutureViewGrantChange(t *testing.T) {
	viewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: viewGrantConfigFuture(t, databaseName, viewName, roleName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "view_name", viewName),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "on_future", "false"),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "privilege", "SELECT"),
				),
			},
			// CHANGE FROM CURRENT TO FUTURE VIEWS
			{
				Config: viewGrantConfigFuture(t, databaseName, viewName, roleName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "view_name", ""),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "privilege", "SELECT"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_view_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func viewGrantConfigShares(t *testing.T, database_name, view_name, role, share_name string) string {
	r := require.New(t)

	tmpl := template.Must(template.New("shares").Parse(`
resource "snowflake_database" "test" {
  name = "{{.database_name}}"
}

resource "snowflake_schema" "test" {
	name = "{{ .schema_name }}"
	database = snowflake_database.test.name
}

resource "snowflake_view" "test" {
  name      = "{{.view_name}}"
  database  = "{{.database_name}}"
  schema    = "{{ .schema_name }}"
  statement = "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
  is_secure = true

  depends_on = [snowflake_database.test, snowflake_schema.test]
}

resource "snowflake_role" "test" {
  name = "{{.role_name}}"
}

resource "snowflake_share" "test" {
  name     = "{{.share_name}}"
}

resource "snowflake_database_grant" "test" {
  database_name = "{{ .database_name }}"
  shares        = ["{{ .share_name }}"]

  depends_on = [snowflake_database.test, snowflake_share.test]
}

resource "snowflake_view_grant" "test" {
  view_name     = "{{ .view_name }}"
  database_name = "{{ .database_name }}"
  roles         = ["{{ .role_name }}"]
	shares        = ["{{ .share_name }}"]
	schema_name = "{{ .schema_name }}"

  // HACK(el): There is a problem with the provider where
  // in older versions of terraform referencing role.name will
  // trick the provider into thinking there are no roles inputted
  // so I hard-code the references.
  depends_on = [snowflake_database_grant.test, snowflake_role.test, snowflake_share.test, snowflake_view.test, snowflake_schema.test]
}`))

	out := bytes.NewBuffer(nil)
	err := tmpl.Execute(out, map[string]string{
		"share_name":    share_name,
		"database_name": database_name,
		"schema_name":   database_name,
		"role_name":     role,
		"view_name":     view_name,
	})
	r.NoError(err)

	return out.String()
}

func viewGrantConfigFuture(t *testing.T, database_name, view_name string, role string, future bool) string {
	r := require.New(t)

	view_name_config := "view_name = snowflake_view.test.name"
	if future {
		view_name_config = "on_future = true"
	}

	config := `
resource "snowflake_database" "test" {
  name = "{{ .database_name }}"
}

resource "snowflake_schema" "test" {
	name = "{{ .schema_name }}"
	database = snowflake_database.test.name
}

resource "snowflake_view" "test" {
  name      = "{{.view_name}}"
	database  = snowflake_database.test.name
	schema    = snowflake_schema.test.name
  statement = "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
  is_secure = true
}

resource "snowflake_role" "test" {
  name = "{{.role_name}}"
}

resource "snowflake_view_grant" "test" {
  {{.view_name_config}}
  database_name = snowflake_view.test.database
	roles         = ["{{.role_name}}"]
	schema_name   = snowflake_schema.test.name
	depends_on = [snowflake_role.test]
	privilege = "SELECT"
}
`

	out := bytes.NewBuffer(nil)
	tmpl := template.Must(template.New("view)").Parse(config))
	err := tmpl.Execute(out, map[string]string{
		"database_name":    database_name,
		"schema_name":      database_name,
		"view_name":        view_name,
		"role_name":        role,
		"view_name_config": view_name_config,
	})
	r.NoError(err)

	return out.String()
}
