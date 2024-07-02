package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CortexSearchService_basic(t *testing.T) {
	t.Skipf("Skipped for now because of the <err: 090105 (22000): Cannot perform operation. This session does not have a current database. Call 'USE DATABASE', or use a qualified name.> problem.")
	resourceName := "snowflake_cortex_search_service.css"
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newWarehouseName := acc.TestClient().Ids.Alpha()
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(id.Name()),
			"on":         config.StringVariable("id"),
			"database":   config.StringVariable(acc.TestDatabaseName),
			"schema":     config.StringVariable(acc.TestSchemaName),
			"warehouse":  config.StringVariable(acc.TestWarehouseName),
			"query":      config.StringVariable(fmt.Sprintf(`select "id" from %s"`, tableId.FullyQualifiedName())),
			"comment":    config.StringVariable("Terraform acceptance test"),
			"table_name": config.StringVariable(tableId.Name()),
		}
	}
	variableSet2 := m()
	variableSet2["attributes"] = config.ListVariable(config.StringVariable("type"))
	variableSet2["warehouse"] = config.StringVariable(newWarehouseName)
	variableSet2["comment"] = config.StringVariable("Terraform acceptance test - updated")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.CortexSearchService),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "on", "id"),
					resource.TestCheckNoResourceAttr(resourceName, "attributes"),
					resource.TestCheckResourceAttr(resourceName, "warehouse", acc.TestWarehouseName),
					resource.TestCheckResourceAttr(resourceName, "target_lag", "2 minutes"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "query", fmt.Sprintf("select \"id\" from %s", tableId.FullyQualifiedName())),
					resource.TestCheckResourceAttrSet(resourceName, "created_on"),
				),
			},

			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: variableSet2,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "on", "id"),
					resource.TestCheckResourceAttr(resourceName, "attributes", "type"),
					resource.TestCheckResourceAttr(resourceName, "warehouse", newWarehouseName),
					resource.TestCheckResourceAttr(resourceName, "target_lag", "2 minutes"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - updated"),
					resource.TestCheckResourceAttr(resourceName, "query", fmt.Sprintf("select \"id\" from %s", tableId.FullyQualifiedName())),
					resource.TestCheckResourceAttrSet(resourceName, "created_on"),
				),
			},
			// test import
			{
				ConfigDirectory:   config.TestStepDirectory(),
				ConfigVariables:   variableSet2,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
