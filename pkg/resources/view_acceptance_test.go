package resources_test

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_View(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	query := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	otherQuery := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES where ROLE_OWNER like 'foo%%'"

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":        config.StringVariable(accName),
			"database":    config.StringVariable(acc.TestDatabaseName),
			"schema":      config.StringVariable(acc.TestSchemaName),
			"comment":     config.StringVariable("Terraform test resource"),
			"is_secure":   config.BoolVariable(true),
			"or_replace":  config.BoolVariable(false),
			"copy_grants": config.BoolVariable(false),
			"statement":   config.StringVariable(query),
		}
	}
	m2 := m()
	m2["comment"] = config.StringVariable("different comment")
	m2["is_secure"] = config.BoolVariable(false)
	m3 := m()
	m3["comment"] = config.StringVariable("different comment")
	m3["is_secure"] = config.BoolVariable(false)
	m3["statement"] = config.StringVariable(otherQuery)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckViewDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View_basic"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "statement", query),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_view.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "Terraform test resource"),
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "false"),
					checkBool("snowflake_view.test", "is_secure", true),
				),
			},
			// update parameters
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View_basic"),
				ConfigVariables: m2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "statement", query),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_view.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "different comment"),
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "false"),
					checkBool("snowflake_view.test", "is_secure", false),
				),
			},
			// change statement
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View_basic"),
				ConfigVariables: m3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "statement", otherQuery),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_view.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "different comment"),
					// copy grants is currently set to true for recreation
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "true"),
					checkBool("snowflake_view.test", "is_secure", false),
				),
			},
			// change statement externally
			{
				PreConfig: func() {
					alterViewQueryExternally(t, sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, accName), query)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View_basic"),
				ConfigVariables: m3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "statement", otherQuery),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_view.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "different comment"),
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "true"),
					checkBool("snowflake_view.test", "is_secure", false),
				),
			},
			// IMPORT
			{
				ConfigVariables:         m3,
				ResourceName:            "snowflake_view.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"or_replace"},
			},
		},
	})
}

func TestAcc_View_Tags(t *testing.T) {
	viewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tag1Name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tag2Name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	query := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":      config.StringVariable(viewName),
			"database":  config.StringVariable(acc.TestDatabaseName),
			"schema":    config.StringVariable(acc.TestSchemaName),
			"statement": config.StringVariable(query),
			"tag1Name":  config.StringVariable(tag1Name),
			"tag2Name":  config.StringVariable(tag2Name),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckViewDestroy,
		Steps: []resource.TestStep{
			// create tags
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_view.test", "tag.#", "1"),
					resource.TestCheckResourceAttr("snowflake_view.test", "tag.0.name", tag1Name),
					resource.TestCheckResourceAttr("snowflake_view.test", "tag.0.value", "some_value"),
				),
			},
			// update tags
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_view.test", "tag.#", "1"),
					resource.TestCheckResourceAttr("snowflake_view.test", "tag.0.name", tag2Name),
					resource.TestCheckResourceAttr("snowflake_view.test", "tag.0.value", "some_value"),
				),
			},
			// IMPORT
			{
				ConfigVariables:         m(),
				ResourceName:            "snowflake_view.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"or_replace", "tag"},
			},
		},
	})
}

func TestAcc_View_Rename(t *testing.T) {
	viewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	newViewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	query := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":        config.StringVariable(viewName),
			"database":    config.StringVariable(acc.TestDatabaseName),
			"schema":      config.StringVariable(acc.TestSchemaName),
			"comment":     config.StringVariable("Terraform test resource"),
			"is_secure":   config.BoolVariable(true),
			"or_replace":  config.BoolVariable(false),
			"copy_grants": config.BoolVariable(false),
			"statement":   config.StringVariable(query),
		}
	}
	m2 := m()
	m2["name"] = config.StringVariable(newViewName)
	m2["comment"] = config.StringVariable("new comment")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckViewDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View_basic"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", viewName),
				),
			},
			// rename with one param changed
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View_basic"),
				ConfigVariables: m2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", newViewName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "new comment"),
				),
			},
		},
	})
}

func TestAcc_ViewChangeCopyGrants(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":        config.StringVariable(accName),
			"database":    config.StringVariable(acc.TestDatabaseName),
			"schema":      config.StringVariable(acc.TestSchemaName),
			"comment":     config.StringVariable("Terraform test resource"),
			"is_secure":   config.BoolVariable(true),
			"or_replace":  config.BoolVariable(false),
			"copy_grants": config.BoolVariable(false),
			"statement":   config.StringVariable("SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"),
		}
	}
	m2 := m()
	m2["copy_grants"] = config.BoolVariable(true)
	m2["or_replace"] = config.BoolVariable(true)

	var createdOn string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckViewDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View_basic"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "Terraform test resource"),
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "false"),
					checkBool("snowflake_view.test", "is_secure", true),
					resource.TestCheckResourceAttrWith("snowflake_view.test", "created_on", func(value string) error {
						createdOn = value
						return nil
					}),
				),
			},
			// Checks that copy_grants changes don't trigger a drop
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View_basic"),
				ConfigVariables: m2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith("snowflake_view.test", "created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("view was recreated")
						}
						return nil
					}),
					checkBool("snowflake_view.test", "is_secure", true),
				),
			},
		},
	})
}

func TestAcc_ViewChangeCopyGrantsReversed(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":        config.StringVariable(accName),
			"database":    config.StringVariable(acc.TestDatabaseName),
			"schema":      config.StringVariable(acc.TestSchemaName),
			"comment":     config.StringVariable("Terraform test resource"),
			"is_secure":   config.BoolVariable(true),
			"or_replace":  config.BoolVariable(true),
			"copy_grants": config.BoolVariable(true),
			"statement":   config.StringVariable("SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"),
		}
	}
	m2 := m()
	m2["copy_grants"] = config.BoolVariable(false)

	var createdOn string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckViewDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View_basic"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "true"),
					resource.TestCheckResourceAttrWith("snowflake_view.test", "created_on", func(value string) error {
						createdOn = value
						return nil
					}),
					checkBool("snowflake_view.test", "is_secure", true),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_View_basic"),
				ConfigVariables: m2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith("snowflake_view.test", "created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("view was recreated")
						}
						return nil
					}),
					checkBool("snowflake_view.test", "is_secure", true),
				),
			},
		},
	})
}

func TestAcc_ViewStatementUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckViewDestroy,
		Steps: []resource.TestStep{
			{
				Config: viewConfigWithGrants(acc.TestDatabaseName, acc.TestSchemaName, `\"name\"`),
				Check: resource.ComposeTestCheckFunc(
					// there should be more than one privilege, because we applied grant all privileges and initially there's always one which is ownership
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "2"),
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.1.privilege", "SELECT"),
				),
			},
			{
				Config: viewConfigWithGrants(acc.TestDatabaseName, acc.TestSchemaName, "*"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "2"),
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.1.privilege", "SELECT"),
				),
			},
		},
	})
}

func TestAcc_View_copyGrants(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	query := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckViewDestroy,
		Steps: []resource.TestStep{
			{
				Config:      viewConfigWithCopyGrants(acc.TestDatabaseName, acc.TestSchemaName, accName, query, true),
				ExpectError: regexp.MustCompile("all of `copy_grants,or_replace` must be specified"),
			},
			{
				Config: viewConfigWithCopyGrantsAndOrReplace(acc.TestDatabaseName, acc.TestSchemaName, accName, query, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
				),
			},
			{
				Config: viewConfigWithOrReplace(acc.TestDatabaseName, acc.TestSchemaName, accName, query, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
				),
			},
		},
	})
}

func viewConfigWithGrants(databaseName string, schemaName string, selectStatement string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "table" {
  database = "%s"
  schema = "%s"
  name     = "view_test_table"

  column {
    name = "name"
    type = "text"
  }
}

resource "snowflake_view" "test" {
  depends_on = [snowflake_table.table]
  name = "test"
  comment = "created by terraform"
  database = "%s"
  schema = "%s"
  statement = "select %s from \"%s\".\"%s\".\"${snowflake_table.table.name}\""
  or_replace = true
  copy_grants = true
  is_secure = true
}

resource "snowflake_role" "test" {
  name = "test"
}

resource "snowflake_view_grant" "grant" {
  database_name = "%s"
  schema_name = "%s"
  view_name = snowflake_view.test.name
  privilege = "SELECT"
  roles = [snowflake_role.test.name]
}

data "snowflake_grants" "grants" {
  depends_on = [snowflake_view_grant.grant, snowflake_view.test]
  grants_on {
    object_name = "\"%s\".\"%s\".\"${snowflake_view.test.name}\""
    object_type = "VIEW"
  }
}
	`, databaseName, schemaName,
		databaseName, schemaName,
		selectStatement,
		databaseName, schemaName,
		databaseName, schemaName,
		databaseName, schemaName)
}

func viewConfigWithCopyGrants(databaseName string, schemaName string, name string, selectStatement string, copyGrants bool) string {
	return fmt.Sprintf(`
resource "snowflake_view" "test" {
  name = "%[3]s"
  database = "%[1]s"
  schema = "%[2]s"
  statement = "%[4]s"
  copy_grants = %[5]t
}
	`, databaseName, schemaName, name, selectStatement, copyGrants)
}

func viewConfigWithCopyGrantsAndOrReplace(databaseName string, schemaName string, name string, selectStatement string, copyGrants bool, orReplace bool) string {
	return fmt.Sprintf(`
resource "snowflake_view" "test" {
  name = "%[3]s"
  database = "%[1]s"
  schema = "%[2]s"
  statement = "%[4]s"
  copy_grants = %[5]t
  or_replace = %[6]t
}
	`, databaseName, schemaName, name, selectStatement, copyGrants, orReplace)
}

func viewConfigWithOrReplace(databaseName string, schemaName string, name string, selectStatement string, orReplace bool) string {
	return fmt.Sprintf(`
resource "snowflake_view" "test" {
  name = "%[3]s"
  database = "%[1]s"
  schema = "%[2]s"
  statement = "%[4]s"
  or_replace = %[5]t
}
	`, databaseName, schemaName, name, selectStatement, orReplace)
}

func testAccCheckViewDestroy(s *terraform.State) error {
	db := acc.TestAccProvider.Meta().(*sql.DB)
	client := sdk.NewClientFromDB(db)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_view" {
			continue
		}
		ctx := context.Background()
		id := sdk.NewSchemaObjectIdentifier(rs.Primary.Attributes["database"], rs.Primary.Attributes["schema"], rs.Primary.Attributes["name"])
		existingView, err := client.Views.ShowByID(ctx, id)
		if err == nil {
			return fmt.Errorf("view %v still exists", existingView.ID().FullyQualifiedName())
		}
	}
	return nil
}

func alterViewQueryExternally(t *testing.T, id sdk.SchemaObjectIdentifier, query string) {
	t.Helper()

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)
	ctx := context.Background()

	err = client.Views.Create(ctx, sdk.NewCreateViewRequest(id, query).WithOrReplace(sdk.Bool(true)))
	require.NoError(t, err)
}
