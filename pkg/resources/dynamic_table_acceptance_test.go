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
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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
	variableSet2["warehouse"] = config.StringVariable(acc.TestWarehouseName2)
	variableSet2["comment"] = config.StringVariable("Terraform acceptance test - updated")

	variableSet3 := m()
	variableSet3["initialize"] = config.StringVariable(string(sdk.DynamicTableInitializeOnSchedule))

	variableSet4 := m()
	variableSet4["initialize"] = config.StringVariable(string(sdk.DynamicTableInitializeOnSchedule)) // keep the same setting from set 3
	variableSet4["refresh_mode"] = config.StringVariable(string(sdk.DynamicTableRefreshModeFull))

	// used to check whether a dynamic table was replaced
	var createdOn string

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
					resource.TestCheckResourceAttr(resourceName, "initialize", string(sdk.DynamicTableInitializeOnCreate)),
					resource.TestCheckResourceAttr(resourceName, "refresh_mode", string(sdk.DynamicTableRefreshModeAuto)),
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
					// - not used at this time
					// resource.TestCheckResourceAttrSet(resourceName, "automatic_clustering"),
					resource.TestCheckResourceAttrSet(resourceName, "scheduling_state"),
					resource.TestCheckResourceAttrSet(resourceName, "last_suspended_on"),
					resource.TestCheckResourceAttrSet(resourceName, "is_clone"),
					resource.TestCheckResourceAttrSet(resourceName, "is_replica"),
					resource.TestCheckResourceAttrSet(resourceName, "data_timestamp"),

					resource.TestCheckResourceAttrWith(resourceName, "created_on", func(value string) error {
						createdOn = value
						return nil
					}),
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
					resource.TestCheckResourceAttr(resourceName, "warehouse", acc.TestWarehouseName2),
					resource.TestCheckResourceAttr(resourceName, "target_lag.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "target_lag.0.downstream", "true"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test - updated"),

					resource.TestCheckResourceAttrWith(resourceName, "created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("created_on changed from %v to %v", createdOn, value)
						}
						return nil
					}),
				),
			},
			// test changing initialize setting
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: variableSet3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "initialize", string(sdk.DynamicTableInitializeOnSchedule)),

					resource.TestCheckResourceAttrWith(resourceName, "created_on", func(value string) error {
						if value == createdOn {
							return fmt.Errorf("expected created_on to change but was not changed")
						}
						createdOn = value
						return nil
					}),
				),
			},
			// test changing refresh_mode setting
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: variableSet4,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "initialize", string(sdk.DynamicTableInitializeOnSchedule)),
					resource.TestCheckResourceAttr(resourceName, "refresh_mode", string(sdk.DynamicTableRefreshModeFull)),

					resource.TestCheckResourceAttrWith(resourceName, "created_on", func(value string) error {
						if value == createdOn {
							return fmt.Errorf("expected created_on to change but was not changed")
						}
						return nil
					}),
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
	dynamicTableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := dynamicTableName + "_table"
	query := fmt.Sprintf(`select "id" from "%v"."%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName, tableName)
	otherSchema := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":         config.StringVariable(dynamicTableName),
			"database":     config.StringVariable(acc.TestDatabaseName),
			"schema":       config.StringVariable(acc.TestSchemaName),
			"warehouse":    config.StringVariable(acc.TestWarehouseName),
			"query":        config.StringVariable(query),
			"comment":      config.StringVariable("Terraform acceptance test for GH issue 2173"),
			"table_name":   config.StringVariable(tableName),
			"other_schema": config.StringVariable(otherSchema),
		}
	}

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
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_schema.other_schema", "name", otherSchema),
					resource.TestCheckResourceAttr("snowflake_table.t", "name", tableName),
				),
			},
			{
				PreConfig:       func() { createDynamicTableOutsideTerraform(t, otherSchema, dynamicTableName, query) },
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
				),
			},
			{
				// We use the same config here as in the previous step so the plan should be empty.
				ConfigDirectory: acc.ConfigurationSameAsStepN(2),
				ConfigVariables: m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					/*
					 * Before the fix this step resulted in
					 *     # snowflake_dynamic_table.dt will be updated in-place
					 *     ~ resource "snowflake_dynamic_table" "dt" {
					 *         + comment              = "Terraform acceptance test for GH issue 2173"
					 *           id                   = "terraform_test_database|terraform_test_schema|SFVNXKJFAA"
					 *           name                 = "SFVNXKJFAA"
					 *         ~ schema               = "MEYIYWUGGO" -> "terraform_test_schema"
					 *           # (14 unchanged attributes hidden)
					 *     }
					 * which matches the issue description exactly (issue mentioned also query being changed but here for simplicity the same underlying table and query were used).
					 */
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// TestAcc_DynamicTable_issue2134 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2134 issue.
func TestAcc_DynamicTable_issue2134(t *testing.T) {
	dynamicTableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := dynamicTableName + "_table"
	// whitespace (initial tab) is added on purpose here
	query := fmt.Sprintf(`	select "id" from "%v"."%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName, tableName)
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(dynamicTableName),
			"database":   config.StringVariable(acc.TestDatabaseName),
			"schema":     config.StringVariable(acc.TestSchemaName),
			"warehouse":  config.StringVariable(acc.TestWarehouseName),
			"query":      config.StringVariable(query),
			"comment":    config.StringVariable("Terraform acceptance test for GH issue 2134"),
			"table_name": config.StringVariable(tableName),
		}
	}
	m2 := m()
	m2["comment"] = config.StringVariable("Changed comment")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDynamicTableDestroy,
		Steps: []resource.TestStep{
			/*
			 * Before the fix the first step resulted in not empty plan (as expected)
			 *     # snowflake_dynamic_table.dt will be updated in-place
			 *     ~ resource "snowflake_dynamic_table" "dt" {
			 *         id                   = "terraform_test_database|terraform_test_schema|IKLBYWKSOV"
			 *         name                 = "IKLBYWKSOV"
			 *         ~ query                = "select \"id\" from \"terraform_test_database\".\"terraform_test_schema\".\"IKLBYWKSOV_table\"" -> "\tselect \"id\" from \"terraform_test_database\".\"terraform_test_schema\".\"IKLBYWKSOV_table\""
			 *         # (15 unchanged attributes hidden)
			 *     }
			 * which matches the issue description exactly.
			 */
			{
				ConfigDirectory: config.TestStepDirectory(),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
				),
			},
			/*
			 * Before the fix the second step resulted in SQL error (as expected)
			 *     Error: 001003 (42000): SQL compilation error:
			 *         syntax error line 1 at position 86 unexpected '<EOF>'.
			 * which matches the issue description exactly.
			 */
			{
				ConfigDirectory: acc.ConfigurationSameAsStepN(1),
				ConfigVariables: m2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
				),
			},
		},
	})
}

// TestAcc_DynamicTable_issue2276 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2276 issue.
func TestAcc_DynamicTable_issue2276(t *testing.T) {
	dynamicTableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := dynamicTableName + "_table"
	query := fmt.Sprintf(`select "id" from "%v"."%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName, tableName)
	newQuery := fmt.Sprintf(`select "data" from "%v"."%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName, tableName)
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(dynamicTableName),
			"database":   config.StringVariable(acc.TestDatabaseName),
			"schema":     config.StringVariable(acc.TestSchemaName),
			"warehouse":  config.StringVariable(acc.TestWarehouseName),
			"query":      config.StringVariable(query),
			"comment":    config.StringVariable("Terraform acceptance test for GH issue 2276"),
			"table_name": config.StringVariable(tableName),
		}
	}
	m2 := m()
	m2["query"] = config.StringVariable(newQuery)

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
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "query", query),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationSameAsStepN(1),
				ConfigVariables: m2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "query", newQuery),
				),
			},
		},
	})
}

// TestAcc_DynamicTable_issue2329 proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2329 issue.
func TestAcc_DynamicTable_issue2329(t *testing.T) {
	dynamicTableName := strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha) + "AS" + acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	tableName := dynamicTableName + "_table"
	query := fmt.Sprintf(`select "id" from "%v"."%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName, tableName)
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":      config.StringVariable(dynamicTableName),
			"database":  config.StringVariable(acc.TestDatabaseName),
			"schema":    config.StringVariable(acc.TestSchemaName),
			"warehouse": config.StringVariable(acc.TestWarehouseName),
			// spaces added on purpose
			"query":      config.StringVariable("  " + query),
			"comment":    config.StringVariable("Comment with AS on purpose"),
			"table_name": config.StringVariable(tableName),
		}
	}

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
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "query", query),
				),
			},
			// No changes are expected
			{
				ConfigDirectory: acc.ConfigurationSameAsStepN(1),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "query", query),
				),
			},
		},
	})
}

// TestAcc_DynamicTable_issue2329_with_matching_comment proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2329 issue.
func TestAcc_DynamicTable_issue2329_with_matching_comment(t *testing.T) {
	dynamicTableName := strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha) + "AS" + acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	tableName := dynamicTableName + "_table"
	query := fmt.Sprintf(`with temp as (select "id" from "%v"."%v"."%v") select * from temp`, acc.TestDatabaseName, acc.TestSchemaName, tableName)
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(dynamicTableName),
			"database":   config.StringVariable(acc.TestDatabaseName),
			"schema":     config.StringVariable(acc.TestSchemaName),
			"warehouse":  config.StringVariable(acc.TestWarehouseName),
			"query":      config.StringVariable(query),
			"comment":    config.StringVariable("Comment with AS SELECT on purpose"),
			"table_name": config.StringVariable(tableName),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckDynamicTableDestroy,
		Steps: []resource.TestStep{
			// If we match more than one time (in this case in comment) we raise an explanation error.
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_DynamicTable_issue2329/1"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "name", dynamicTableName),
					resource.TestCheckResourceAttr("snowflake_dynamic_table.dt", "query", query),
				),
			},
		},
	})
}

func testAccCheckDynamicTableDestroy(s *terraform.State) error {
	client := acc.TestAccProvider.Meta().(*provider.Context).Client
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

// TODO [SNOW-926148]: currently this dynamic table is not cleaned in the test; it is removed when the whole database is removed - this currently happens in a sweeper
func createDynamicTableOutsideTerraform(t *testing.T, schemaName string, dynamicTableName string, query string) {
	t.Helper()
	client := acc.Client(t)
	ctx := context.Background()

	dynamicTableId := sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, schemaName, dynamicTableName)
	if err := client.DynamicTables.Create(ctx, sdk.NewCreateDynamicTableRequest(dynamicTableId, sdk.NewAccountObjectIdentifier(acc.TestWarehouseName), sdk.TargetLag{MaximumDuration: sdk.String("2 minutes")}, query)); err != nil {
		t.Fatal(fmt.Errorf("error creating dynamic table: %w", err))
	}
}
