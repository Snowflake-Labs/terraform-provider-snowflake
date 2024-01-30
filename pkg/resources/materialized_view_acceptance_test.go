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
	"github.com/stretchr/testify/require"
)

func TestAcc_MaterializedView(t *testing.T) {
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	viewName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	warehouseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

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
				Config: materializedViewConfig(warehouseName, tableName, viewName, queryEscaped, acc.TestDatabaseName, acc.TestSchemaName, "Terraform test resource", true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "statement", query),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", warehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", "Terraform test resource"),
					checkBool("snowflake_materialized_view.test", "is_secure", true),
				),
			},
			// update parameters
			{
				Config: materializedViewConfig(warehouseName, tableName, viewName, queryEscaped, acc.TestDatabaseName, acc.TestSchemaName, "other comment", false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "statement", query),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", warehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", "other comment"),
					checkBool("snowflake_materialized_view.test", "is_secure", false),
				),
			},
			// change statement
			{
				Config: materializedViewConfig(warehouseName, tableName, viewName, otherQueryEscaped, acc.TestDatabaseName, acc.TestSchemaName, "other comment", false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "statement", otherQuery),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", warehouseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "comment", "other comment"),
					checkBool("snowflake_materialized_view.test", "is_secure", false),
				),
			},
			// change statement externally
			{
				PreConfig: func() {
					alterMaterializedViewQueryExternally(t, sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, viewName), query)
				},
				Config: materializedViewConfig(warehouseName, tableName, viewName, otherQueryEscaped, acc.TestDatabaseName, acc.TestSchemaName, "other comment", false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "name", viewName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "statement", otherQuery),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_materialized_view.test", "warehouse", warehouseName),
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

func materializedViewConfig(warehouseName string, tableName string, viewName string, q string, databaseName string, schemaName string, comment string, isSecure bool, orReplace bool) string {
	// convert the cluster from string slice to string
	return fmt.Sprintf(`
resource "snowflake_warehouse" "wh" {
	name = "%s"
}
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
	warehouse = snowflake_warehouse.wh.name
	is_secure = %t
	or_replace = %t
	statement = "%s"

	depends_on = [
  		snowflake_table.test
  	]
}
`, warehouseName, tableName, databaseName, schemaName, viewName, comment, databaseName, schemaName, isSecure, orReplace, q)
}

func testAccCheckMaterializedViewDestroy(s *terraform.State) error {
	db := acc.TestAccProvider.Meta().(*sql.DB)
	client := sdk.NewClientFromDB(db)
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

func alterMaterializedViewQueryExternally(t *testing.T, id sdk.SchemaObjectIdentifier, query string) {
	t.Helper()

	client, err := sdk.NewDefaultClient()
	require.NoError(t, err)
	ctx := context.Background()

	err = client.MaterializedViews.Create(ctx, sdk.NewCreateMaterializedViewRequest(id, query).WithOrReplace(sdk.Bool(true)))
	require.NoError(t, err)
}
