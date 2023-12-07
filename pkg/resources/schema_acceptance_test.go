package resources_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Schema(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := "Terraform acceptance test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckSchemaDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(name),
					"database": config.StringVariable(acc.TestDatabaseName),
					"comment":  config.StringVariable(comment),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", comment),
					checkBool("snowflake_schema.test", "is_transient", false),
					checkBool("snowflake_schema.test", "is_managed", false),
				),
			},
		},
	})
}

func TestAcc_Schema_Rename(t *testing.T) {
	oldSchemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	newSchemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := "Terraform acceptance test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckSchemaDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(oldSchemaName),
					"database": config.StringVariable(acc.TestDatabaseName),
					"comment":  config.StringVariable(comment),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", oldSchemaName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", comment),
					checkBool("snowflake_schema.test", "is_transient", false),
					checkBool("snowflake_schema.test", "is_managed", false),
				),
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(newSchemaName),
					"database": config.StringVariable(acc.TestDatabaseName),
					"comment":  config.StringVariable(comment),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", newSchemaName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test", "comment", comment),
					checkBool("snowflake_schema.test", "is_transient", false),
					checkBool("snowflake_schema.test", "is_managed", false),
				),
			},
		},
	})
}

// TestAcc_Schema_TwoSchemasWithTheSameNameOnDifferentDatabases proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2209 issue.
func TestAcc_Schema_TwoSchemasWithTheSameNameOnDifferentDatabases(t *testing.T) {
	name := "test_schema"
	newDatabaseName := "SELDQBXEKC"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckSchemaDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(name),
					"database": config.StringVariable(acc.TestDatabaseName),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
				),
			},
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: map[string]config.Variable{
					"name":         config.StringVariable(name),
					"database":     config.StringVariable(acc.TestDatabaseName),
					"new_database": config.StringVariable(newDatabaseName),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.test", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_schema.test_2", "name", name),
					resource.TestCheckResourceAttr("snowflake_schema.test_2", "database", newDatabaseName),
				),
			},
		},
	})
}

func testAccCheckSchemaDestroy(s *terraform.State) error {
	db := acc.TestAccProvider.Meta().(*sql.DB)
	client := sdk.NewClientFromDB(db)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_schema" {
			continue
		}
		ctx := context.Background()
		id := sdk.NewDatabaseObjectIdentifier(rs.Primary.Attributes["database"], rs.Primary.Attributes["name"])
		schema, err := client.Schemas.ShowByID(ctx, id)
		if err == nil {
			return fmt.Errorf("schema %v still exists", schema.Name)
		}
	}
	return nil
}
