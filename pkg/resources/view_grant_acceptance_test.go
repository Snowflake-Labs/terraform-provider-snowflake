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

func TestAcc_ViewGrantBasic(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: viewGrantConfig(name, normal),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "view_name", name),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "privilege", "SELECT"),
				),
			},
			{
				ResourceName:      "snowflake_view_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
					"on_all",                 // not defined in Snowflake, can't be imported
				},
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
		Providers:    providers(),
		CheckDestroy: nil,
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
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: viewGrantConfig(name, normal),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "view_name", name),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "on_future", "false"),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "privilege", "SELECT"),
				),
			},
			// CHANGE FROM CURRENT TO FUTURE VIEWS
			{
				Config: viewGrantConfig(name, onFuture),
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
					"on_all",                 // not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func viewGrantConfigShares(t *testing.T, databaseName, viewName, role, shareName string) string {
	t.Helper()
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
		"share_name":    shareName,
		"database_name": databaseName,
		"schema_name":   databaseName,
		"role_name":     role,
		"view_name":     viewName,
	})
	r.NoError(err)

	return out.String()
}

func viewGrantConfig(name string, grantType grantType) string {
	var viewNameConfig string
	switch grantType {
	case normal:
		viewNameConfig = "view_name = snowflake_view.test.name"
	case onFuture:
		viewNameConfig = "on_future = true"
	case onAll:
		viewNameConfig = "on_all = true"
	}

	return fmt.Sprintf(`
resource "snowflake_database" "test" {
  name = "%s"
}

resource "snowflake_schema" "test" {
	name = "%s"
	database = snowflake_database.test.name
}

resource "snowflake_view" "test" {
  name      = "%s"
  database  = snowflake_database.test.name
  schema    = snowflake_schema.test.name
  statement = "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
  is_secure = true
}

resource "snowflake_role" "test" {
  name = "%s"
}

resource "snowflake_view_grant" "test" {
  %s
  database_name = snowflake_view.test.database
  roles         = [snowflake_role.test.name]
  schema_name   = snowflake_schema.test.name
  privilege = "SELECT"
}
`, name, name, name, name, viewNameConfig)
}

func TestAcc_ViewGrantOnAll(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: viewGrantConfig(name, onAll),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "schema_name", name),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "on_all", "true"),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "privilege", "SELECT"),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "roles.#", "1"),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "roles.0", name),
					testRolesAndShares(t, "snowflake_view_grant.test", []string{name}),
				),
			},
		},
	})
}
