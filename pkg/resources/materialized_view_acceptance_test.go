package resources_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_MaterializedView(t *testing.T) {
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	viewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	queryEscaped := fmt.Sprintf("SELECT ID, DATA FROM \\\"%s\\\"", tableName)
	query := fmt.Sprintf(`SELECT ID, DATA FROM "%s"`, tableName)
	otherQueryEscaped := fmt.Sprintf("SELECT ID, DATA FROM \\\"%s\\\" WHERE ID LIKE 'foo%%'", tableName)
	otherQuery := fmt.Sprintf(`SELECT ID, DATA FROM "%s" WHERE ID LIKE 'foo%%'`, tableName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckMaterializedViewDestroy,
		Steps: []resource.TestStep{
			{
				Config: materializedViewConfig(acc.TestWarehouseName, tableName, viewName, queryEscaped, acc.TestDatabaseName, acc.TestSchemaName, "Terraform test resource", true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "statement", query),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", acc.TestWarehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", "Terraform test resource"),
					checkBool("snowflake_materialized_view.test", "is_secure", true),
				),
			},
			// update parameters
			{
				Config: materializedViewConfig(acc.TestWarehouseName, tableName, viewName, queryEscaped, acc.TestDatabaseName, acc.TestSchemaName, "other comment", false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "statement", query),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", acc.TestWarehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", "other comment"),
					checkBool("snowflake_materialized_view.test", "is_secure", false),
				),
			},
			// change statement
			{
				Config: materializedViewConfig(acc.TestWarehouseName, tableName, viewName, otherQueryEscaped, acc.TestDatabaseName, acc.TestSchemaName, "other comment", false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "statement", otherQuery),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", acc.TestWarehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", "other comment"),
					checkBool("snowflake_materialized_view.test", "is_secure", false),
				),
			},
			// change statement externally
			{
				PreConfig: func() {
					alterMaterializedViewQueryExternally(t, sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, viewName), query, acc.TestWarehouseName)
				},
				Config: materializedViewConfig(acc.TestWarehouseName, tableName, viewName, otherQueryEscaped, acc.TestDatabaseName, acc.TestSchemaName, "other comment", false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "statement", otherQuery),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", acc.TestWarehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", "other comment"),
					checkBool("snowflake_materialized_view.test", "is_secure", false),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_materialized_view.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"or_replace", "warehouse"},
			},
		},
	})
}

func TestAcc_MaterializedView_Tags(t *testing.T) {
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	viewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tag1Name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tag2Name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	queryEscaped := fmt.Sprintf("SELECT ID FROM \\\"%s\\\"", tableName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckMaterializedViewDestroy,
		Steps: []resource.TestStep{
			// create tags
			{
				Config: materializedViewConfigWithTags(acc.TestWarehouseName, tableName, viewName, queryEscaped, acc.TestDatabaseName, acc.TestSchemaName, "test_tag", tag1Name, tag2Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "tag.#", "1"),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "tag.0.name", tag1Name),
				),
			},
			// update tags
			{
				Config: materializedViewConfigWithTags(acc.TestWarehouseName, tableName, viewName, queryEscaped, acc.TestDatabaseName, acc.TestSchemaName, "test_tag_2", tag1Name, tag2Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "tag.#", "1"),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "tag.0.name", tag2Name),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_materialized_view.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"or_replace", "warehouse", "tag"},
			},
		},
	})
}

func TestAcc_MaterializedView_Rename(t *testing.T) {
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	viewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	newViewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	queryEscaped := fmt.Sprintf("SELECT ID FROM \\\"%s\\\"", tableName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckMaterializedViewDestroy,
		Steps: []resource.TestStep{
			{
				Config: materializedViewConfig(acc.TestWarehouseName, tableName, viewName, queryEscaped, acc.TestDatabaseName, acc.TestSchemaName, "Terraform test resource", true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
				),
			},
			// rename with one param change
			{
				Config: materializedViewConfig(acc.TestWarehouseName, tableName, newViewName, queryEscaped, acc.TestDatabaseName, acc.TestSchemaName, "Terraform test resource", false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", newViewName),
				),
			},
		},
	})
}

func materializedViewConfig(warehouseName string, tableName string, viewName string, q string, databaseName string, schemaName string, comment string, isSecure bool, orReplace bool) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test" {
	name     = "%s"
	database = "%s"
	schema   = "%s"

	column {
		name = "ID"
		type = "NUMBER(38,0)"
	}

	column {
		name = "DATA"
		type = "VARCHAR(16777216)"
	}
}

resource "snowflake_materialized_view" "test" {
	name      = "%s"
	comment   = "%s"
	database  = "%s"
	schema    = "%s"
	warehouse = "%s"
	is_secure = %t
	or_replace = %t
	statement = "%s"

	depends_on = [
  		snowflake_table.test
  	]
}
`, tableName, databaseName, schemaName, viewName, comment, databaseName, schemaName, warehouseName, isSecure, orReplace, q)
}

func materializedViewConfigWithTags(warehouseName string, tableName string, viewName string, q string, databaseName string, schemaName string, tag string, tag1Name string, tag2Name string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test" {
	name     = "%[1]s"
	database = "%[2]s"
	schema   = "%[3]s"

	column {
		name = "ID"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_tag" "test_tag" {
	name     = "%[8]s"
	database = "%[2]s"
	schema   = "%[3]s"
}

resource "snowflake_tag" "test_tag_2" {
	name     = "%[9]s"
	database = "%[2]s"
	schema   = "%[3]s"
}

resource "snowflake_materialized_view" "test" {
	name      = "%[4]s"
	database  = "%[2]s"
	schema    = "%[3]s"
	warehouse = "%[5]s"
	statement = "%[6]s"

	tag {
		name = snowflake_tag.%[7]s.name
		schema = snowflake_tag.%[7]s.schema
		database = snowflake_tag.%[7]s.database
		value = "some_value"
	}

	depends_on = [
		snowflake_table.test
	]
}
`, tableName, databaseName, schemaName, viewName, warehouseName, q, tag, tag1Name, tag2Name)
}

func testAccCheckMaterializedViewDestroy(s *terraform.State) error {
	client := acc.TestAccProvider.Meta().(*provider.Context).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_materialized_view" {
			continue
		}
		ctx := context.Background()
		id := sdk.NewSchemaObjectIdentifier(rs.Primary.Attributes["database"], rs.Primary.Attributes["schema"], rs.Primary.Attributes["name"])
		existingMaterializedView, err := client.MaterializedViews.ShowByID(ctx, id)
		if err == nil {
			return fmt.Errorf("materialized view %v still exists", existingMaterializedView.ID().FullyQualifiedName())
		}
	}
	return nil
}

func alterMaterializedViewQueryExternally(t *testing.T, id sdk.SchemaObjectIdentifier, query string, warehouse string) {
	t.Helper()

	client := acc.Client(t)
	ctx := context.Background()

	err := client.Sessions.UseWarehouse(ctx, sdk.NewAccountObjectIdentifier(warehouse))
	require.NoError(t, err)

	err = client.MaterializedViews.Create(ctx, sdk.NewCreateMaterializedViewRequest(id, query).WithOrReplace(sdk.Bool(true)))
	require.NoError(t, err)
}
