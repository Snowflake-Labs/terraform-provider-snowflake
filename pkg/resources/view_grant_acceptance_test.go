package resources_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAcc_ViewGrantBasic(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: viewGrantConfig(name, normal, "SELECT", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "view_name", name),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "privilege", "SELECT"),
				),
			},
			// UPDATE ALL PRIVILEGES
			{
				Config: viewGrantConfig(name, normal, "ALL PRIVILEGES", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "view_name", name),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "privilege", "ALL PRIVILEGES"),
				),
			},
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

func TestAcc_ViewGrantShares(t *testing.T) {
	viewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	roleName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	shareName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: viewGrantConfigShares(t, viewName, roleName, shareName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "view_name", viewName),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "privilege", "SELECT"),
				),
			},
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

func TestAcc_ViewGrantChange(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: viewGrantConfig(name, normal, "SELECT", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "view_name", name),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "on_future", "false"),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "privilege", "SELECT"),
				),
			},
			// CHANGE FROM CURRENT TO FUTURE VIEWS
			{
				Config: viewGrantConfig(name, onFuture, "SELECT", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("snowflake_view_grant.test", "view_name"),
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

func viewGrantConfigShares(t *testing.T, viewName, role, shareName string) string {
	t.Helper()
	r := require.New(t)

	tmpl := template.Must(template.New("shares").Parse(`
resource "snowflake_view" "test" {
	name      = "{{.view_name}}"
	database  = "{{.database_name}}"
	schema    = "{{ .schema_name }}"
	statement = "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	is_secure = true
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

	depends_on = [snowflake_share.test]
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
	depends_on = [snowflake_database_grant.test, snowflake_role.test, snowflake_share.test, snowflake_view.test]
}`))

	out := bytes.NewBuffer(nil)
	err := tmpl.Execute(out, map[string]string{
		"share_name":    shareName,
		"database_name": acc.TestDatabaseName,
		"schema_name":   acc.TestSchemaName,
		"role_name":     role,
		"view_name":     viewName,
	})
	r.NoError(err)

	return out.String()
}

func viewGrantConfig(name string, grantType grantType, privilege string, databaseName string, schemaName string) string {
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
resource "snowflake_view" "test" {
  	name      = "%s"
  	database  = "%s"
  	schema    = "%s"
  	statement = "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
  	is_secure = true
}

resource "snowflake_role" "test" {
  	name = "%s"
}

resource "snowflake_view_grant" "test" {
  	%s
  	database_name = "%s"
  	roles         = [snowflake_role.test.name]
  	schema_name   = "%s"
  	privilege = "%s"
}
`, name, databaseName, schemaName, name, viewNameConfig, databaseName, schemaName, privilege)
}

func TestAcc_ViewGrantOnAll(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: viewGrantConfig(name, onAll, "SELECT", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_view_grant.test", "schema_name", acc.TestSchemaName),
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
