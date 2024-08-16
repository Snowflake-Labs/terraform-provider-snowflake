package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_RowAccessPolicy(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":     config.StringVariable(id.Name()),
			"database": config.StringVariable(acc.TestDatabaseName),
			"schema":   config.StringVariable(acc.TestSchemaName),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.RowAccessPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "row_access_expression", "case when current_role() in ('ANALYST') then true else false end"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "signature.N", "VARCHAR"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "signature.V", "VARCHAR"),
				),
			},
			// change comment and expression
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "comment", "Terraform acceptance test - changed comment"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "row_access_expression", "case when current_role() in ('ANALYST') then false else true end"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "signature.N", "VARCHAR"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "signature.V", "VARCHAR"),
				),
			},
			// change signature
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "comment", "Terraform acceptance test - changed comment"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "row_access_expression", "case when current_role() in ('ANALYST') then false else true end"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "signature.V", "BOOLEAN"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "signature.X", "TIMESTAMP_NTZ"),
				),
			},
			// IMPORT
			{
				ConfigVariables:   m(),
				ResourceName:      "snowflake_row_access_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
