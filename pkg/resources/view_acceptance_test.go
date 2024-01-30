package resources_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_View(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckViewDestroy,
		Steps: []resource.TestStep{
			{
				Config: viewConfig(accName, false, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "Terraform test resource"),
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "false"),
					checkBool("snowflake_view.test", "is_secure", true), // this is from user_acceptance_test.go
				),
			},
		},
	})
}

func TestAcc_View2(t *testing.T) {
	accName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckViewDestroy,
		Steps: []resource.TestStep{
			{
				Config: viewConfig(accName, false, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES where ROLE_OWNER like 'foo%%';", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "Terraform test resource"),
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "false"),
					checkBool("snowflake_view.test", "is_secure", true), // this is from user_acceptance_test.go
				),
			},
		},
	})
}

func TestAcc_ViewWithCopyGrants(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckViewDestroy,
		Steps: []resource.TestStep{
			{
				Config: viewConfig(accName, true, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_view.test", "comment", "Terraform test resource"),
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "true"),
					checkBool("snowflake_view.test", "is_secure", true), // this is from user_acceptance_test.go
				),
			},
		},
	})
}

// Checks that copy_grants changes don't trigger a drop
func TestAcc_ViewChangeCopyGrants(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

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
				Config: viewConfig(accName, false, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_view.test", "copy_grants", "false"),
					resource.TestCheckResourceAttrWith("snowflake_view.test", "created_on", func(value string) error {
						createdOn = value
						return nil
					}),
					checkBool("snowflake_view.test", "is_secure", true),
				),
			},
			{
				Config: viewConfig(accName, true, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith("snowflake_view.test", "created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("View was recreated")
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
				Config: viewConfig(accName, true, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES", acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: viewConfig(accName, false, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith("snowflake_view.test", "created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("View was recreated")
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

func viewConfig(n string, copyGrants bool, q string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_view" "test" {
	name        = "%v"
	comment     = "Terraform test resource"
	database    = "%s"
	schema      = "%s"
	is_secure   = true
	or_replace  = %t
	copy_grants = %t
	statement   = "%s"
}
`, n, databaseName, schemaName, copyGrants, copyGrants, q)
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

func testAccCheckViewDestroy(s *terraform.State) error {
	db := acc.TestAccProvider.Meta().(*sql.DB)
	client := sdk.NewClientFromDB(db)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_materialized_view" {
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
