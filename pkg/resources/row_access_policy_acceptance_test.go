package resources_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_RowAccessPolicy(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":     config.StringVariable(name),
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
		CheckDestroy: testAccCheckRowAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "name", name),
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
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "name", name),
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
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "name", name),
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

func testAccCheckRowAccessPolicyDestroy(s *terraform.State) error {
	client := acc.TestAccProvider.Meta().(*provider.Context).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_row_access_policy" {
			continue
		}
		ctx := context.Background()
		id := sdk.NewSchemaObjectIdentifier(rs.Primary.Attributes["database"], rs.Primary.Attributes["schema"], rs.Primary.Attributes["name"])
		existingRowAccessPolicy, err := client.RowAccessPolicies.ShowByID(ctx, id)
		if err == nil {
			return fmt.Errorf("row access policy %v still exists", existingRowAccessPolicy.ID().FullyQualifiedName())
		}
	}
	return nil
}
