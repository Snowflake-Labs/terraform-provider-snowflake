package resources_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TableWithSeparateDataRetentionObjectParameterWithoutLifecycle(t *testing.T) {
	accName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfig(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARIANT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.comment", ""),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableAndDataRetentionParameterConfigWithoutLifecycle(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.comment", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.comment", ""),
					resource.TestCheckResourceAttr("snowflake_object_parameter.data_retention_in_time", "value", "30"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAcc_TableWithSeparateDataRetentionObjectParameterWithLifecycle(t *testing.T) {
	accName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfig(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARIANT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.comment", ""),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableAndDataRetentionParameterConfigWithLifecycle(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.comment", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.comment", ""),
					resource.TestCheckResourceAttr("snowflake_object_parameter.data_retention_in_time", "value", "30"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: updatedTableAndDataRetentionParameterConfigWithLifecycle(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Table with a separate data retention parameter"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.comment", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.comment", ""),
					resource.TestCheckResourceAttr("snowflake_object_parameter.data_retention_in_time", "value", "30"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
		},
	})
}

func TestAcc_Table(t *testing.T) {
	table1Id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	table2Id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	table3Id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfig(table1Id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", table1Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", table1Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARIANT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.comment", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.schema_evolution_record", ""),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableConfig2(table1Id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", table1Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", table1Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.comment", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "FLOAT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.comment", ""),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableConfig3(table2Id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "fully_qualified_name", table2Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.name", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.name", "col2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.#", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.0", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.type", "FLOAT"),
				),
			},
			{
				Config: tableConfig4(table2Id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "fully_qualified_name", table2Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.name", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.name", "col2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.1", "\"col2\""),
				),
			},
			{
				Config: tableConfig5(table2Id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "fully_qualified_name", table2Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.name", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.name", "col2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.0", "\"col2\""),
				),
			},
			{
				Config: tableConfig6(table1Id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", table1Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", table1Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.nullable", "true"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "FLOAT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.nullable", "false"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableConfig7(table1Id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", table1Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", table1Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.nullable", "true"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "FLOAT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.nullable", "false"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "cluster_by.0"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "primary_key.0.keys.0", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "primary_key.0.name", ""),
				),
			},
			{
				Config: tableConfig8(table1Id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", table1Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", table1Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.nullable", "true"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "FLOAT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.nullable", "false"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "cluster_by.0"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "primary_key.0.keys.0", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "primary_key.0.keys.1", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "primary_key.0.name", "new_name"),
				),
			},
			{
				Config: tableConfig9CreateTableWithColumnComment(table2Id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "fully_qualified_name", table2Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.name", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.nullable", "true"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.comment", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.name", "COL2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.type", "FLOAT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.nullable", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.comment", "some comment"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table2", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table2", "primary_key.0"),
				),
			},
			{
				Config: tableConfig10AlterTableColumnComment(table2Id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "fully_qualified_name", table2Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.name", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.nullable", "true"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.comment", "other comment"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.name", "COL2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.type", "FLOAT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.nullable", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.comment", ""),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table2", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table2", "primary_key.0"),
				),
			},
			{
				Config: tableConfig11AlterTableAddColumnWithComment(table2Id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "fully_qualified_name", table2Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.name", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.nullable", "true"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.comment", "other comment"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.name", "COL2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.type", "FLOAT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.nullable", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.comment", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.2.name", "COL3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.2.type", "NUMBER(38,0)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.2.nullable", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.2.comment", "extra"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table2", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table2", "primary_key.0"),
				),
			},
			{
				Config: tableConfig12CreateTableWithDataRetention(table3Id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "name", table3Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "fully_qualified_name", table3Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "data_retention_time_in_days", "10"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.0.type", "VARIANT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.1.name", "column2"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table3", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table3", "primary_key.0"),
				),
			},
			{
				Config: tableConfig13AlterTableDataRetention(table3Id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "name", table3Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "fully_qualified_name", table3Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "data_retention_time_in_days", "0"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.0.type", "VARIANT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.1.name", "column2"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table3", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table3", "primary_key.0"),
				),
			},
			{
				Config: tableConfig14AlterTableEnableChangeTracking(table3Id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "name", table3Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "fully_qualified_name", table3Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "data_retention_time_in_days", "0"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "change_tracking", "true"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.0.type", "VARIANT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "column.1.name", "column2"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table3", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table3", "primary_key.0"),
				),
			},
			{
				Config: tableConfig15CreateTableWithChangeTracking(table1Id.Name(), acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", table1Id.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", table1Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "data_retention_time_in_days", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "true"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARIANT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "cluster_by.0"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
		},
	})
}

func tableAndDataRetentionParameterConfigWithoutLifecycle(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table" {
	name     = "%s"
	database = "%s"
	schema   = "%s"
	comment  = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
}

resource "snowflake_object_parameter" "data_retention_in_time" {
	key = "DATA_RETENTION_TIME_IN_DAYS"
	value = "30"
    object_type = "TABLE"
    object_identifier {
		name = "%s"
		database = "%s"
		schema = "%s"
	}
	depends_on = [snowflake_table.test_table]
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName, name, databaseName, schemaName)
}

func tableAndDataRetentionParameterConfigWithLifecycle(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table" {
	name     = "%s"
	database = "%s"
	schema   = "%s"
	comment  = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
	lifecycle {
		ignore_changes = [
			"data_retention_time_in_days"
		]
	}
}

resource "snowflake_object_parameter" "data_retention_in_time" {
	key = "DATA_RETENTION_TIME_IN_DAYS"
	value = "30"
    object_type = "TABLE"
    object_identifier {
		name = "%s"
		database = "%s"
		schema = "%s"
	}
	depends_on = [snowflake_table.test_table]
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName, name, databaseName, schemaName)
}

func updatedTableAndDataRetentionParameterConfigWithLifecycle(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table" {
	name     = "%s"
	database = "%s"
	schema   = "%s"
	comment  = "Table with a separate data retention parameter"
	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
	lifecycle {
		ignore_changes = [
			"data_retention_time_in_days"
		]
	}
}
resource "snowflake_object_parameter" "data_retention_in_time" {
	key = "DATA_RETENTION_TIME_IN_DAYS"
	value = "30"
    object_type = "TABLE"
    object_identifier {
		name = "%s"
		database = "%s"
		schema = "%s"
	}
	depends_on = [snowflake_table.test_table]
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName, name, databaseName, schemaName)
}

func tableConfig(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table" {
	name     = "%s"
	database = "%s"
	schema   = "%s"
	data_retention_time_in_days = 1
	comment  = "Terraform acceptance test"
	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName)
}

func tableConfig2(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table" {
	name     = "%s"
	database = "%s"
	schema   = "%s"
	comment  = "Terraform acceptance test"
	data_retention_time_in_days = 1
	column {
		name = "column2"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "column3"
		type = "FLOAT"
	}
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName)
}

func tableConfig3(table2Name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table2" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
	data_retention_time_in_days = 1
	comment  = "Terraform acceptance test"
	cluster_by = ["COL1"]
	column {
		name = "COL1"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "col2"
		type = "FLOAT"
	}
}
`
	return fmt.Sprintf(s, table2Name, databaseName, schemaName)
}

func tableConfig4(table2Name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table2" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
	comment  = "Terraform acceptance test"
	data_retention_time_in_days = 1
	cluster_by = ["COL1","\"col2\""]
	column {
		name = "COL1"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "col2"
		type = "FLOAT"
	}
}
`
	return fmt.Sprintf(s, table2Name, databaseName, schemaName)
}

func tableConfig5(table2Name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table2" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
	comment  = "Terraform acceptance test"
	data_retention_time_in_days = 1
	cluster_by = ["\"col2\"","COL1"]
	column {
		name = "COL1"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "col2"
		type = "FLOAT"
	}
}
`
	return fmt.Sprintf(s, table2Name, databaseName, schemaName)
}

func tableConfig6(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
	comment  = "Terraform acceptance test"
	data_retention_time_in_days = 1
	column {
		name = "column2"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "column3"
		type = "FLOAT"
		nullable = false
	}
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName)
}

func tableConfig7(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
	comment  = "Terraform acceptance test"
	data_retention_time_in_days = 1
	column {
		name = "column2"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "column3"
		type = "FLOAT"
		nullable = false
	}
	primary_key {
		keys = ["column2"]
	}
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName)
}

func tableConfig8(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
	comment  = "Terraform acceptance test"
	data_retention_time_in_days = 1
	column {
		name = "column2"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "column3"
		type = "FLOAT"
		nullable = false
	}
	primary_key {
		name = "new_name"
		keys = ["column2","column3"]
	}
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName)
}

func tableConfig9CreateTableWithColumnComment(table2Name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table2" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
	comment  = "Terraform acceptance test"
	data_retention_time_in_days = 1
	column {
		name = "COL1"
		type = "VARCHAR(16777216)"
	}
	column {
		name     = "COL2"
		type     = "FLOAT"
		nullable = false
		comment  = "some comment"
	}
}
`
	return fmt.Sprintf(s, table2Name, databaseName, schemaName)
}

func tableConfig10AlterTableColumnComment(table2Name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table2" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
	comment             = "Terraform acceptance test"
	data_retention_time_in_days = 1
	column {
		name    = "COL1"
		type    = "VARCHAR(16777216)"
		comment = "other comment"
	}
	column {
		name     = "COL2"
		type     = "FLOAT"
		nullable = false
	}
}
`
	return fmt.Sprintf(s, table2Name, databaseName, schemaName)
}

func tableConfig11AlterTableAddColumnWithComment(table2Name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table2" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
	comment             = "Terraform acceptance test"
	data_retention_time_in_days = 1
	column {
		name    = "COL1"
		type    = "VARCHAR(16777216)"
		comment = "other comment"
	}
	column {
		name     = "COL2"
		type     = "FLOAT"
		nullable = false
	}
	column {
		name     = "COL3"
		type     = "NUMBER(38,0)"
		nullable = false
		comment  = "extra"
	}
}
`
	return fmt.Sprintf(s, table2Name, databaseName, schemaName)
}

func tableConfig12CreateTableWithDataRetention(table3Name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table3" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
	comment             = "Terraform acceptance test"
	data_retention_time_in_days = 10
	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
}
`
	return fmt.Sprintf(s, table3Name, databaseName, schemaName)
}

func tableConfig13AlterTableDataRetention(table3Name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table3" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
	comment             = "Terraform acceptance test"
	data_retention_time_in_days = 0
	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
}
`
	return fmt.Sprintf(s, table3Name, databaseName, schemaName)
}

func tableConfig14AlterTableEnableChangeTracking(table3Name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table3" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
	comment             = "Terraform acceptance test"
	data_retention_time_in_days = 0
	change_tracking     = true
	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
}
`
	return fmt.Sprintf(s, table3Name, databaseName, schemaName)
}

func tableConfig15CreateTableWithChangeTracking(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
	comment             = "Terraform acceptance test"
	data_retention_time_in_days = 1
	change_tracking     = true
	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName)
}

func TestAcc_TableDefaults(t *testing.T) {
	accName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableColumnWithDefaults(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARCHAR(16)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.default.0.constant", "hello"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.type.default.0.expression"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.type.default.0.sequence"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "TIMESTAMP_NTZ(9)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.1.type.default.0.constant"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.default.0.expression", "CURRENT_TIMESTAMP()"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.1.type.default.0.sequence"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.type", "NUMBER(38,0)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.2.type.default.0.constant"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.2.type.default.0.expression"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.default.0.sequence", fmt.Sprintf(`"%v"."%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName, accName)),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableColumnWithoutDefaults(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARCHAR(16)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.default.0"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "TIMESTAMP_NTZ(9)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.1.type.default"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.type", "NUMBER(38,0)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.2.type.default.0.constant"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.2.type.default.0.expression"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.default.0.sequence", fmt.Sprintf(`"%v"."%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName, accName)),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
		},
	})
}

func tableColumnWithDefaults(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_sequence" "test_seq" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
}

resource "snowflake_table" "test_table" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
	comment             = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARCHAR(16)"
		default {
			constant = "hello"
		}
	}
	column {
		name = "column2"
		type = "TIMESTAMP_NTZ(9)"
		default {
			expression = "CURRENT_TIMESTAMP()"
		}
	}
	column {
		name = "column3"
		type = "NUMBER(38,0)"
		default {
			sequence = snowflake_sequence.test_seq.fully_qualified_name
		}
	}
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName, name, databaseName, schemaName)
}

func tableColumnWithoutDefaults(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_sequence" "test_seq" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
}

resource "snowflake_table" "test_table" {
	name                = "%s"
	database            = "%s"
	schema              = "%s"
	comment             = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARCHAR(16)"
	}
	column {
		name = "column2"
		type = "TIMESTAMP_NTZ(9)"
	}
	column {
		name = "column3"
		type = "NUMBER(38,0)"
		default {
			sequence = snowflake_sequence.test_seq.fully_qualified_name
		}
	}
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName, name, databaseName, schemaName)
}

func TestAcc_TableTags(t *testing.T) {
	accName := acc.TestClient().Ids.Alpha()
	tagName := acc.TestClient().Ids.Alpha()
	tag2Name := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableWithTags(accName, tagName, tag2Name, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "tag.0.name", tagName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "tag.0.value", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "tag.0.database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "tag.0.schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "tag.1.name", tag2Name),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "tag.1.value", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "tag.1.database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "tag.1.schema", acc.TestSchemaName),
				),
			},
		},
	})
}

func tableWithTags(name string, tagName string, tag2Name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_tag" "test_tag" {
	name     = "%[2]s"
	database = "%[4]s"
	schema   = "%[5]s"
	allowed_values = ["%[1]s"]
	comment  = "Terraform acceptance test"
}

resource "snowflake_tag" "test2_tag" {
	name     = "%[3]s"
	database = "%[4]s"
	schema   = "%[5]s"
	allowed_values = ["%[1]s"]
	comment  = "Terraform acceptance test"
}

resource "snowflake_table" "test_table" {
	database = "%[4]s"
	schema   = "%[5]s"
	name                = "%[1]s"
	comment             = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARCHAR(16)"
	}

	tag {
		name = snowflake_tag.test_tag.name
		schema = snowflake_tag.test_tag.schema
		database = snowflake_tag.test_tag.database
		value = "%[1]s"
	}

	tag {
		name = snowflake_tag.test2_tag.name
		schema = snowflake_tag.test2_tag.schema
		database = snowflake_tag.test2_tag.database
		value = "%[1]s"
	}
}
`
	return fmt.Sprintf(s, name, tagName, tag2Name, databaseName, schemaName)
}

func TestAcc_TableIdentity(t *testing.T) {
	accName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableColumnWithIdentityDefault(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "NUMBER(38,0)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.type.default.0.expression"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.type.default.0.sequence"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "TIMESTAMP_NTZ(9)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.1.type.default.0.constant"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.1.type.default.0.sequence"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.type", "NUMBER(38,0)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.2.type.default.0.constant"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.2.type.default.0.expression"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.identity.0.start_num", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.identity.0.step_num", "1"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableColumnWithIdentity(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "NUMBER(38,0)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.type.default.0.expression"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.type.default.0.sequence"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "TIMESTAMP_NTZ(9)"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.1.type.default.0.constant"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.1.type.default.0.sequence"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.2.type.default.0.constant"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.2.type.default.0.expression"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.identity.0.start_num"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "column.0.identity.0.step_num"),
					// we've dropped the previous identity column and making sure that adding a new column as an identity works
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.identity.0.start_num", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.identity.0.step_num", "4"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
		},
	})
}

func tableColumnWithIdentityDefault(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_sequence" "test_seq" {
	name     = "%s"
	database = "%s"
	schema   = "%s"
}

resource "snowflake_table" "test_table" {
	name     = "%s"
	database = "%s"
	schema   = "%s"
	comment             = "Terraform acceptance test"

	column {
		name = "column1"
		type = "NUMBER(38,0)"
	}
	column {
		name = "column2"
		type = "TIMESTAMP_NTZ(9)"
	}
	column {
		name = "column3"
		type = "NUMBER(38,0)"
		identity {
		}
	}
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName, name, databaseName, schemaName)
}

func tableColumnWithIdentity(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_sequence" "test_seq" {
	name     = "%s"
	database = "%s"
	schema   = "%s"
}

resource "snowflake_table" "test_table" {
	name     = "%s"
	database = "%s"
	schema   = "%s"
	comment             = "Terraform acceptance test"

	column {
		name = "column1"
		type = "NUMBER(38,0)"
	}
	column {
		name = "column2"
		type = "TIMESTAMP_NTZ(9)"
	}

	column {
		name = "column4"
		type = "NUMBER(38,0)"
		identity {
			start_num = 2
			step_num = 4
		}
	}
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName, name, databaseName, schemaName)
}

func TestAcc_TableCollate(t *testing.T) {
	accName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableColumnWithCollate(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.collate", "en"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.collate", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.collate", ""),
				),
			},
			{
				Config: addColumnWithCollate(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "4"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.name", "column4"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.collate", "utf8"),
				),
			},
			{
				Config:      alterTableColumnWithIncompatibleCollate(accName, acc.TestDatabaseName, acc.TestSchemaName),
				ExpectError: regexp.MustCompile("\"VARCHAR\\(100\\) COLLATE 'fr'\" because they have incompatible collations\\."),
			},
		},
	})
}

func tableColumnWithCollate(name string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[2]s"
	schema   = "%[3]s"
	name     = "%[1]s"

	column {
		name = "column1"
		type = "VARCHAR(100)"
		collate = "en"
	}
	column {
		name = "column2"
		type = "VARCHAR(100)"
		collate = ""
	}
	column {
		name = "column3"
		type = "VARCHAR(100)"
	}
}
`, name, databaseName, schemaName)
}

func addColumnWithCollate(name string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[2]s"
	schema   = "%[3]s"
	name     = "%[1]s"

	column {
		name = "column1"
		type = "VARCHAR(100)"
		collate = "en"
	}
	column {
		name = "column2"
		type = "VARCHAR(100)"
		collate = ""
	}
	column {
		name = "column3"
		type = "VARCHAR(100)"
	}
	column {
		name = "column4"
		type = "VARCHAR"
		collate = "utf8"
	}
}
`, name, databaseName, schemaName)
}

func alterTableColumnWithIncompatibleCollate(name string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	database = "%[2]s"
	schema   = "%[3]s"
	name     = "%[1]s"

	column {
		name = "column1"
		type = "VARCHAR(100)"
		collate = "fr"
	}
	column {
		name = "column2"
		type = "VARCHAR(100)"
		collate = ""
	}
	column {
		name = "column3"
		type = "VARCHAR(100)"
	}
	column {
		name = "column4"
		type = "VARCHAR"
		collate = "utf8"
	}
}
`, name, databaseName, schemaName)
}

func TestAcc_TableRename(t *testing.T) {
	oldId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	oldComment := acc.TestClient().Ids.Alpha()
	newComment := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfigWithName(oldId.Name(), acc.TestDatabaseName, acc.TestSchemaName, oldComment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", oldId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", oldId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", oldComment),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARIANT"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableConfigWithName(newId.Name(), acc.TestDatabaseName, acc.TestSchemaName, newComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_table.test_table", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", newId.Name()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "fully_qualified_name", newId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", newComment),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARIANT"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
		},
	})
}

func tableConfigWithName(tableName string, databaseName string, schemaName string, comment string) string {
	s := `
resource "snowflake_table" "test_table" {
	name     = "%s"
	database = "%s"
	schema   = "%s"
    comment  = "%s"
	column {
		name = "column1"
		type = "VARIANT"
	}
}
`
	return fmt.Sprintf(s, tableName, databaseName, schemaName, comment)
}

func TestAcc_Table_MaskingPolicy(t *testing.T) {
	accName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableWithMaskingPolicy(accName, acc.TestDatabaseName, acc.TestSchemaName, "policy1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.masking_policy", acc.TestClient().Ids.NewSchemaObjectIdentifier(fmt.Sprintf("%s1", accName)).FullyQualifiedName()),
				),
			},
			// this step proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2186
			{
				Config: tableWithMaskingPolicy(accName, acc.TestDatabaseName, acc.TestSchemaName, "policy2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.masking_policy", acc.TestClient().Ids.NewSchemaObjectIdentifier(fmt.Sprintf("%s2", accName)).FullyQualifiedName()),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Table_DefaultDataRetentionTime(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	configWithDatabaseDataRetentionSet := func(databaseDataRetentionTime int) config.Variables {
		return config.Variables{
			"database":                     config.StringVariable(databaseId.Name()),
			"schema":                       config.StringVariable(schemaId.Name()),
			"table":                        config.StringVariable(tableId.Name()),
			"database_data_retention_time": config.IntegerVariable(databaseDataRetentionTime),
		}
	}

	configWithSchemaDataRetentionSet := func(databaseDataRetentionTime int, schemaDataRetentionTime int) config.Variables {
		vars := configWithDatabaseDataRetentionSet(databaseDataRetentionTime)
		vars["schema_data_retention_time"] = config.IntegerVariable(schemaDataRetentionTime)
		return vars
	}

	configWithTableDataRetentionSet := func(databaseDataRetentionTime int, schemaDataRetentionTime int, tableDataRetentionTime int) config.Variables {
		vars := configWithSchemaDataRetentionSet(databaseDataRetentionTime, schemaDataRetentionTime)
		vars["table_data_retention_time"] = config.IntegerVariable(tableDataRetentionTime)
		return vars
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithDatabaseDataRetentionSet"),
				ConfigVariables: configWithDatabaseDataRetentionSet(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 5, 5, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithSchemaDataRetentionSet"),
				ConfigVariables: configWithSchemaDataRetentionSet(5, 10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 5, 10, 10),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(10, 3, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "5"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(10, 3, 15),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "15"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 15),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithSchemaDataRetentionSet"),
				ConfigVariables: configWithSchemaDataRetentionSet(10, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 3),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithDatabaseDataRetentionSet"),
				ConfigVariables: configWithDatabaseDataRetentionSet(10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 10, 10),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(10, 5, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "0"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 5, 0),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(10, 5, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "3"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 5, 3),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Table_DefaultDataRetentionTime_SetOutsideOfTerraform(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	configWithDatabaseDataRetentionSet := func(databaseDataRetentionTime int) config.Variables {
		return config.Variables{
			"database":                     config.StringVariable(databaseId.Name()),
			"schema":                       config.StringVariable(schemaId.Name()),
			"table":                        config.StringVariable(tableId.Name()),
			"database_data_retention_time": config.IntegerVariable(databaseDataRetentionTime),
		}
	}

	configWithTableDataRetentionSet := func(databaseDataRetentionTime int, schemaDataRetentionTime int, tableDataRetentionTime int) config.Variables {
		vars := configWithDatabaseDataRetentionSet(databaseDataRetentionTime)
		vars["schema_data_retention_time"] = config.IntegerVariable(schemaDataRetentionTime)
		vars["table_data_retention_time"] = config.IntegerVariable(tableDataRetentionTime)
		return vars
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithDatabaseDataRetentionSet"),
				ConfigVariables: configWithDatabaseDataRetentionSet(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 5, 5, 5),
				),
			},
			{
				PreConfig: func() {
					acc.TestClient().Table.SetDataRetentionTime(t, tableId, 20)
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithDatabaseDataRetentionSet"),
				ConfigVariables: configWithDatabaseDataRetentionSet(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 5, 5, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(5, 10, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "3"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 5, 10, 3),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Table_DefaultDataRetentionTimeSettingUnsetting(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaId := acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	configWithDatabaseDataRetentionSet := func(databaseDataRetentionTime int) config.Variables {
		return config.Variables{
			"database":                     config.StringVariable(databaseId.Name()),
			"schema":                       config.StringVariable(schemaId.Name()),
			"table":                        config.StringVariable(tableId.Name()),
			"database_data_retention_time": config.IntegerVariable(databaseDataRetentionTime),
		}
	}

	configWithSchemaDataRetentionSet := func(databaseDataRetentionTime int, schemaDataRetentionTime int) config.Variables {
		vars := configWithDatabaseDataRetentionSet(databaseDataRetentionTime)
		vars["schema_data_retention_time"] = config.IntegerVariable(schemaDataRetentionTime)
		return vars
	}

	configWithTableDataRetentionSet := func(databaseDataRetentionTime int, schemaDataRetentionTime int, tableDataRetentionTime int) config.Variables {
		vars := configWithSchemaDataRetentionSet(databaseDataRetentionTime, schemaDataRetentionTime)
		vars["table_data_retention_time"] = config.IntegerVariable(tableDataRetentionTime)
		return vars
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(10, 3, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "5"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(10, 3, -1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 3),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithSchemaDataRetentionSet"),
				ConfigVariables: configWithSchemaDataRetentionSet(10, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 3),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(10, 3, -1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", r.IntDefaultString),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 3),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(10, 3, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "5"),
					checkDatabaseSchemaAndTableDataRetentionTime(tableId, 10, 3, 5),
				),
			},
		},
	})
}

func tableWithMaskingPolicy(name string, databaseName string, schemaName string, policy string) string {
	s := `
resource "snowflake_masking_policy" "policy1" {
	name 	 		   = "%[1]s1"
	database 	       = "%[2]s"
	schema   		   = "%[3]s"
	signature {
		column {
			name = "val"
			type = "VARCHAR"
		}
	}
	masking_expression = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
	return_data_type   = "VARCHAR(16777216)"
}

resource "snowflake_masking_policy" "policy2" {
	name 	 		   = "%[1]s2"
	database 	       = "%[2]s"
	schema   		   = "%[3]s"
	signature {
		column {
			name = "val"
			type = "VARCHAR"
		}
	}
	masking_expression = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
	return_data_type   = "VARCHAR(16777216)"
}

resource "snowflake_table" "test_table" {
	name     = "%[1]s"
	database = "%[2]s"
	schema   = "%[3]s"
	comment  = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARCHAR(16)"
		masking_policy = snowflake_masking_policy.%[4]s.fully_qualified_name
	}
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName, policy)
}

// proves issues https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2110 and https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2495
func TestAcc_Table_ClusterBy(t *testing.T) {
	accName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfigWithComplexClusterBy(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "cluster_by.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "cluster_by.0", "date_trunc('month', LAST_LOAD_TIME)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "cluster_by.1", "COL1"),
				),
			},
		},
	})
}

func tableConfigWithComplexClusterBy(name string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	name                = "%[1]s"
	database            = "%[2]s"
	schema              = "%[3]s"
	cluster_by = ["date_trunc('month', LAST_LOAD_TIME)", "COL1"]
	column {
		name = "COL1"
		type = "VARCHAR(16777216)"
	}
    column {
        name     = "LAST_LOAD_TIME"
        type     = "TIMESTAMP_LTZ(6)"
        nullable = true
    }
}
`, name, databaseName, schemaName)
}

// TODO [SNOW-1348114]: do not trim the data type (e.g. NUMBER(38,0) -> NUMBER(36,0) diff is ignored); finish the test
// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2588 is fixed
func TestAcc_ColumnTypeChangeWithNonTextType(t *testing.T) {
	t.Skipf("Will be fixed with tables redesign in SNOW-1348114")
	accName := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfigWithNumberColumnType(accName, acc.TestDatabaseName, acc.TestSchemaName, "NUMBER(38,0)"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "id"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "NUMBER(38,0)"),
				),
			},
			{
				Config: tableConfigWithNumberColumnType(accName, acc.TestDatabaseName, acc.TestSchemaName, "NUMBER(36,0)"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "id"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "NUMBER(36,0)"),
				),
			},
		},
	})
}

func tableConfigWithNumberColumnType(name string, databaseName string, schemaName string, columnType string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
	name                = "%[1]s"
	database            = "%[2]s"
	schema              = "%[3]s"

	column {
		name = "id"
		type = "%[4]s"
	}
}
`, name, databaseName, schemaName, columnType)
}

func checkDatabaseSchemaAndTableDataRetentionTime(id sdk.SchemaObjectIdentifier, expectedDatabaseRetentionDays int, expectedSchemaRetentionDays int, expectedTableRetentionsDays int) func(state *terraform.State) error {
	return func(state *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()

		database, err := client.Databases.ShowByID(ctx, id.DatabaseId())
		if err != nil {
			return err
		}

		if database.RetentionTime != expectedDatabaseRetentionDays {
			return fmt.Errorf("invalid database retention time, expected: %d, got: %d", expectedDatabaseRetentionDays, database.RetentionTime)
		}

		s, err := client.Schemas.ShowByID(ctx, id.SchemaId())
		if err != nil {
			return err
		}

		// "retention_time" may sometimes be an empty string instead of an integer
		var schemaRetentionTime int64
		{
			rt := s.RetentionTime
			if rt == "" {
				rt = "0"
			}

			schemaRetentionTime, err = strconv.ParseInt(rt, 10, 64)
			if err != nil {
				return err
			}
		}

		if schemaRetentionTime != int64(expectedSchemaRetentionDays) {
			return fmt.Errorf("invalid schema retention time, expected: %d, got: %d", expectedSchemaRetentionDays, schemaRetentionTime)
		}

		table, err := client.Tables.ShowByID(ctx, id)
		if err != nil {
			return err
		}

		if table.RetentionTime != expectedTableRetentionsDays {
			return fmt.Errorf("invalid table retention time, expected: %d, got: %d", expectedTableRetentionsDays, table.RetentionTime)
		}

		return nil
	}
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2733 is fixed
func TestAcc_Table_gh2733(t *testing.T) {
	name := acc.TestClient().Ids.Alpha()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Table),
		Steps: []resource.TestStep{
			{
				Config: tableConfigGh2733(acc.TestDatabaseName, acc.TestSchemaName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", name),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "MY_INT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "NUMBER(38,0)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "MY_STRING"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "VARCHAR(16777216)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "MY_DATE"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.type", "TIMESTAMP_NTZ(9)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.name", "MY_DATE2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.type", "TIMESTAMP_NTZ(9)"),
				),
			},
		},
	})
}

func tableConfigGh2733(database string, schema string, name string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
  database            = "%[1]s"
  schema              = "%[2]s"
  name                = "%[3]s"

  column {
    name = "MY_INT"
    type = "int"
    # type  = "NUMBER(38,0)" # Should be equivalent
  }

  column {
    name = "MY_STRING"
    type = "VARCHAR(16777216)"
    # type = "STRING" # Should be equivalent
  }

  column {
    name = "MY_DATE"
    type = "TIMESTAMP_NTZ"
    # type = "TIMESTAMP_NTZ(9)" # Should be equivalent
  }

  column {
    name = "MY_DATE2"
    type = "DATETIME"
    # type = "TIMESTAMP_NTZ" # Equivalent to TIMESTAMP_NTZ
  }
}
`, database, schema, name)
}

func TestAcc_Table_migrateFromVersion_0_94_1(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_table.test_table"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: tableConfig(id.Name(), id.DatabaseName(), id.SchemaName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "qualified_name", id.FullyQualifiedName()),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfig(id.Name(), id.DatabaseName(), id.SchemaName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckNoResourceAttr(resourceName, "qualified_name"),
				),
			},
		},
	})
}

func TestAcc_Table_SuppressQuotingOnDefaultSequence_issue2644(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	name := acc.TestClient().Ids.Alpha()
	resourceName := "snowflake_table.test_table"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				ExpectNonEmptyPlan: true,
				Config:             tableConfigWithSequence(name, databaseName, schemaName),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigWithSequence(name, databaseName, schemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.0.default.0.sequence", sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name).FullyQualifiedName()),
				),
			},
		},
	})
}

func tableConfigWithSequence(name string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test_database" {
	name = "%[2]s"
}

resource "snowflake_schema" "test_schema" {
	depends_on = [snowflake_database.test_database]
	name = "%[3]s"
	database = "%[2]s"
}

resource "snowflake_sequence" "test_sequence" {
	depends_on = [snowflake_schema.test_schema]
	name     = "%[1]s"
	database = "%[2]s"
	schema   = "%[3]s"
}

resource "snowflake_table" "test_table" {
	depends_on = [snowflake_sequence.test_sequence]
	name     = "%[1]s"
	database = "%[2]s"
	schema   = "%[3]s"
	data_retention_time_in_days = 1
	comment  = "Terraform acceptance test"
	column {
		name = "column1"
		type = "NUMBER"
		default {
			sequence = "%[2]s.%[3]s.%[1]s"
		}
	}
}
`, name, databaseName, schemaName)
}

func TestAcc_Table_issue3007_textColumn(t *testing.T) {
	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_table.test_table"

	defaultVarchar := fmt.Sprintf("VARCHAR(%d)", sdk.DefaultVarcharLength)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "VARCHAR(3)"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.0.type", "NUMBER(11,2)"),
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "VARCHAR(3)"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "VARCHAR(256)"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(resourceName, "column.1.type", tfjson.ActionUpdate, sdk.String("VARCHAR(3)"), sdk.String("VARCHAR(256)")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "VARCHAR(256)"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "VARCHAR"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(resourceName, "column.1.type", tfjson.ActionUpdate, sdk.String("VARCHAR(256)"), sdk.String("VARCHAR")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultVarchar),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, defaultVarchar),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultVarchar),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "text"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultVarchar),
				),
			},
		},
	})
}

// TODO [SNOW-1348114]: visit with table rework (e.g. changing scale is not supported: err 040052 (22000): SQL compilation error: cannot change column SOME_COLUMN from type NUMBER(38,0) to NUMBER(11,2) because changing the scale of a number is not supported.)
func TestAcc_Table_issue3007_numberColumn(t *testing.T) {
	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_table.test_table"

	defaultNumber := fmt.Sprintf("NUMBER(%d,%d)", sdk.DefaultNumberPrecision, sdk.DefaultNumberScale)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "NUMBER"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.0.type", "NUMBER(11,2)"),
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "NUMBER(38,0)"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "NUMBER(11)"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(resourceName, "column.1.type", tfjson.ActionUpdate, sdk.String("NUMBER(38,0)"), sdk.String("NUMBER(11)")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.0.type", "NUMBER(11,2)"),
					resource.TestCheckResourceAttr(resourceName, "column.1.type", "NUMBER(11,0)"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "NUMBER"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(resourceName, "column.1.type", tfjson.ActionUpdate, sdk.String("NUMBER(11,0)"), sdk.String("NUMBER")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultNumber),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, defaultNumber),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultNumber),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   tableConfigIssue3007(tableId, "decimal"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "column.1.type", defaultNumber),
				),
			},
		},
	})
}

func tableConfigIssue3007(tableId sdk.SchemaObjectIdentifier, dataType string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "test_table" {
    name     = "%[1]s"
    database = "%[2]s"
    schema   = "%[3]s"
    comment  = "Issue 3007 confirmation"

    column {
        name = "ID"
        type = "NUMBER(11,2)"
    }
    
    column {
        name = "SOME_COLUMN"
        type = "%[4]s"
    }
}
`, tableId.Name(), tableId.DatabaseName(), tableId.SchemaName(), dataType)
}
