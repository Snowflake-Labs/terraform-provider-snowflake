package resources_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DynamicTable_basic(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_dynamic_table.dt"
	tableName := name + "_table"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(name),
			"database":   config.StringVariable(acc.TestDatabaseName),
			"schema":     config.StringVariable(acc.TestSchemaName),
			"warehouse":  config.StringVariable(acc.TestWarehouseName),
			"query":      config.StringVariable(fmt.Sprintf(`select "id" from "%v"."%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName, tableName)),
			"comment":    config.StringVariable("Terraform acceptance test"),
			"table_name": config.StringVariable(tableName),
		}
	}
	variableSet2 := m()
	variableSet2["comment"] = config.StringVariable("Terraform acceptance test - updated")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDynamicTableDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "warehouse", acc.TestWarehouseName),
					resource.TestCheckResourceAttr(resourceName, "target_lag.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "target_lag.0.maximum_duration", "2 minutes"),
					resource.TestCheckResourceAttr(resourceName, "query", fmt.Sprintf("select \"id\" from \"%v\".\"%v\".\"%v\"", acc.TestDatabaseName, acc.TestSchemaName, tableName)),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),

					// computed attributes

					// - not used at this time
					//  resource.TestCheckResourceAttrSet(resourceName, "cluster_by"),
					resource.TestCheckResourceAttrSet(resourceName, "rows"),
					resource.TestCheckResourceAttrSet(resourceName, "bytes"),
					resource.TestCheckResourceAttrSet(resourceName, "owner"),
					resource.TestCheckResourceAttrSet(resourceName, "refresh_mode"),
					// - not used at this time
					// resource.TestCheckResourceAttrSet(resourceName, "automatic_clustering"),
					resource.TestCheckResourceAttrSet(resourceName, "scheduling_state"),
					resource.TestCheckResourceAttrSet(resourceName, "last_suspended_on"),
					resource.TestCheckResourceAttrSet(resourceName, "is_clone"),
					resource.TestCheckResourceAttrSet(resourceName, "is_replica"),
					resource.TestCheckResourceAttrSet(resourceName, "data_timestamp"),
				),
			},
			// test target lag to downstream and change comment

			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: variableSet2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "target_lag.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "target_lag.0.downstream", "true"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - updated"),
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

// TestAcc_DynamicTable_issue2173 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2173 issue.
func TestAcc_DynamicTable_issue2173(t *testing.T) {
	resourceName := "snowflake_dynamic_table.dt"
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := name + "_table"
	query := fmt.Sprintf(`select "id" from "%v"."%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName, tableName)
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(name),
			"database":   config.StringVariable(acc.TestDatabaseName),
			"schema":     config.StringVariable(acc.TestSchemaName),
			"warehouse":  config.StringVariable(acc.TestWarehouseName),
			"query":      config.StringVariable(query),
			"comment":    config.StringVariable("Terraform acceptance test for GH issue 2173"),
			"table_name": config.StringVariable(tableName),
		}
	}
	otherSchema := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			createDynamicTableOutsideTerraform(t, otherSchema, name)
		},
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: dropAdditionalSchemaAndCheckDynamicTableDestroy(t, otherSchema),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationSameAsStepN(1),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					/*
					 * Before the fix this step resulted in
					 *     # snowflake_dynamic_table.dt will be updated in-place
					 *     ~ resource "snowflake_dynamic_table" "dt" {
					 *     + comment              = "Terraform acceptance test for GH issue 2173"
					 *     id                   = "terraform_test_database|terraform_test_schema|JJIVTACVOU"
					 *     name                 = "JJIVTACVOU"
					 *     ~ query                = "select id from \"terraform_test_database\".\"WFAUDTBOJW\".\"JJIVTACVOU_table\"" -> "select \"id\" from \"terraform_test_database\".\"terraform_test_schema\".\"JJIVTACVOU_table\""
					 *     ~ schema               = "WFAUDTBOJW" -> "terraform_test_schema"
					 *     # (13 unchanged attributes hidden) *
					 *     # (1 unchanged block hidden)
					 *     }
					 * which matches the issue description exactly.
					 */
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

func testAccCheckDynamicTableDestroy(s *terraform.State) error {
	db := acc.TestAccProvider.Meta().(*sql.DB)
	client := sdk.NewClientFromDB(db)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_dynamic_table" {
			continue
		}
		ctx := context.Background()
		id := sdk.NewSchemaObjectIdentifier(rs.Primary.Attributes["database"], rs.Primary.Attributes["schema"], rs.Primary.Attributes["name"])
		dynamicTable, err := client.DynamicTables.ShowByID(ctx, id)
		if err == nil {
			return fmt.Errorf("dynamic table %v still exists", dynamicTable.Name)
		}
	}
	return nil
}

func createDynamicTableOutsideTerraform(t *testing.T, schemaName string, dynamicTableName string) {
	t.Helper()
	client, err := sdk.NewDefaultClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	schemaId := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, schemaName)
	if err := client.Schemas.Create(ctx, schemaId, &sdk.CreateSchemaOptions{
		IfNotExists: sdk.Bool(true),
	}); err != nil {
		t.Fatal(fmt.Errorf("error creating schema: %w", err))
	}

	tableName := dynamicTableName + "_table"
	tableId := sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, schemaName, tableName)

	_, err = client.ExecForTests(ctx, fmt.Sprintf("CREATE TABLE %s (id NUMBER)", tableId.FullyQualifiedName()))
	if err != nil {
		t.Fatal(fmt.Errorf("error creating table: %w", err))
	}

	query := fmt.Sprintf(`select id from "%v"."%v"."%v"`, acc.TestDatabaseName, schemaName, tableName)
	dynamicTableId := sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, schemaName, dynamicTableName)
	if err := client.DynamicTables.Create(ctx, sdk.NewCreateDynamicTableRequest(dynamicTableId, sdk.NewAccountObjectIdentifier(acc.TestWarehouseName), sdk.TargetLag{MaximumDuration: sdk.String("2 minutes")}, query)); err != nil {
		t.Fatal(fmt.Errorf("error creating dynamic table: %w", err))
	}
}

func dropAdditionalSchemaAndCheckDynamicTableDestroy(t *testing.T, schemaName string) resource.TestCheckFunc {
	t.Helper()
	return func(s *terraform.State) error {
		client, err := sdk.NewDefaultClient()
		if err != nil {
			t.Fatal(err)
		}
		ctx := context.Background()

		schemaId := sdk.NewDatabaseObjectIdentifier(acc.TestDatabaseName, schemaName)
		_ = client.Schemas.Drop(ctx, schemaId, nil)

		return testAccCheckDynamicTableDestroy(s)
	}
}
