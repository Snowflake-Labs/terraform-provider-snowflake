package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_TableWithSeparateDataRetentionObjectParameterWithoutLifecycle(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_TABLE_DATA_RETENTION_TESTS"); ok {
		t.Skip("Skipping TestAcc_TableWithSeparateDataRetentionObjectParameterWithoutLifecycle")
	}

	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableConfig(accName),
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
				Config: tableAndDataRetentionParameterConfigWithoutLifecycle(accName),
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
	if _, ok := os.LookupEnv("SKIP_TABLE_DATA_RETENTION_TESTS"); ok {
		t.Skip("Skipping TestAcc_TableWithSeparateDataRetentionObjectParameterWithLifecycle")
	}

	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableConfig(accName),
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
				Config: tableAndDataRetentionParameterConfigWithLifecycle(accName),
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
				Config: updatedTableAndDataRetentionParameterConfigWithLifecycle(accName),
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
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	table2Name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	table3Name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
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
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableConfig2(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
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
				Config: tableConfig3(table2Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Name),
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
				Config: tableConfig4(table2Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Name),
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
				Config: tableConfig5(table2Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Name),
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
				Config: tableConfig6(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
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
				Config: tableConfig7(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
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
				Config: tableConfig8(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
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
				Config: tableConfig9CreateTableWithColumnComment(table2Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Name),
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
				Config: tableConfig10AlterTableColumnComment(table2Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Name),
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
				Config: tableConfig11AlterTableAddColumnWithComment(table2Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Name),
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
				Config: tableConfig12CreateTableWithDataRetention(table3Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "name", table3Name),
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
				Config: tableConfig13AlterTableDataRetention(table3Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "name", table3Name),
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
				Config: tableConfig14AlterTableEnableChangeTracking(table3Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table3", "name", table3Name),
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
				Config: tableConfig15CreateTableWithChangeTracking(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
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

func tableAndDataRetentionParameterConfigWithoutLifecycle(name string) string {
	s := `
resource "snowflake_table" "test_table" {
	database = "terraform_test_database"
	schema   = "terraform_test_schema"
	name     = "%s"
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
		database = "terraform_test_database"
		schema = "terraform_test_schema"
		name = "%s"
	}
}
`
	return fmt.Sprintf(s, name, name)
}

func tableAndDataRetentionParameterConfigWithLifecycle(name string) string {
	s := `
resource "snowflake_table" "test_table" {
	database = "terraform_test_database"
	schema   = "terraform_test_schema"
	name     = "%s"
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
		database = "terraform_test_database"
		schema = "terraform_test_schema"
		name = "%s"
	}
}
`
	return fmt.Sprintf(s, name, name)
}

func updatedTableAndDataRetentionParameterConfigWithLifecycle(name string) string {
	s := `
resource "snowflake_table" "test_table" {
	database = "terraform_test_database"
	schema   = "terraform_test_schema"
	name     = "%s"
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
		database = "terraform_test_database"
		schema = "terraform_test_schema"
		name = "%s"
	}
}
`
	return fmt.Sprintf(s, name, name)
}

func tableConfig(name string) string {
	s := `
resource "snowflake_table" "test_table" {
	database = "terraform_test_database"
	schema   = "terraform_test_schema"
	data_retention_time_in_days = 1
	name     = "%s"
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
	return fmt.Sprintf(s, name)
}

func tableConfig2(name string) string {
	s := `
resource "snowflake_table" "test_table" {
	database = "terraform_test_database"
	schema   = "terraform_test_schema"
	name     = "%s"
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
	return fmt.Sprintf(s, name)
}

func tableConfig3(table2Name string) string {
	s := `
resource "snowflake_table" "test_table2" {
	database = "terraform_test_database"
	schema   = "terraform_test_schema"
	data_retention_time_in_days = 1
	name     = "%s"
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
	return fmt.Sprintf(s, table2Name)
}

func tableConfig4(table2Name string) string {
	s := `
resource "snowflake_table" "test_table2" {
	database = "terraform_test_database"
	schema   = "terraform_test_schema"
	name     = "%s"
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
	return fmt.Sprintf(s, table2Name)
}

func tableConfig5(table2Name string) string {
	s := `
resource "snowflake_table" "test_table2" {
	database = "terraform_test_database"
	schema   = "terraform_test_schema"
	name     = "%s"
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
	return fmt.Sprintf(s, table2Name)
}

func tableConfig6(name string) string {
	s := `
resource "snowflake_table" "test_table" {
	database = "terraform_test_database"
	schema   = "terraform_test_schema"
	name     = "%s"
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
	return fmt.Sprintf(s, name)
}

func tableConfig7(name string) string {
	s := `
resource "snowflake_table" "test_table" {
	database = "terraform_test_database"
	schema   = "terraform_test_schema"
	name     = "%s"
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
		name = ""
		keys = ["column2"]
	}
}
`
	return fmt.Sprintf(s, name)
}

func tableConfig8(name string) string {
	s := `
resource "snowflake_table" "test_table" {
	database = "terraform_test_database"
	schema   = "terraform_test_schema"
	name     = "%s"
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
	return fmt.Sprintf(s, name)
}

func tableConfig9CreateTableWithColumnComment(table2Name string) string {
	s := `
resource "snowflake_table" "test_table2" {
	database = "terraform_test_database"
	schema   = "terraform_test_schema"
	name     = "%s"
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
	return fmt.Sprintf(s, table2Name)
}

func tableConfig10AlterTableColumnComment(table2Name string) string {
	s := `
resource "snowflake_table" "test_table2" {
	database            = "terraform_test_database"
	schema              = "terraform_test_schema"
	name                = "%s"
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
	return fmt.Sprintf(s, table2Name)
}

func tableConfig11AlterTableAddColumnWithComment(table2Name string) string {
	s := `
resource "snowflake_table" "test_table2" {
	database            = "terraform_test_database"
	schema              = "terraform_test_schema"
	name                = "%s"
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
	return fmt.Sprintf(s, table2Name)
}

func tableConfig12CreateTableWithDataRetention(table3Name string) string {
	s := `
resource "snowflake_table" "test_table3" {
	database            = "terraform_test_database"
	schema              = "terraform_test_schema"
	name                = "%s"
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
	return fmt.Sprintf(s, table3Name)
}

func tableConfig13AlterTableDataRetention(table3Name string) string {
	s := `
resource "snowflake_table" "test_table3" {
	database            = "terraform_test_database"
	schema              = "terraform_test_schema"
	name                = "%s"
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
	return fmt.Sprintf(s, table3Name)
}

func tableConfig14AlterTableEnableChangeTracking(table3Name string) string {
	s := `
resource "snowflake_table" "test_table3" {
	database            = "terraform_test_database"
	schema              = "terraform_test_schema"
	name                = "%s"
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
	return fmt.Sprintf(s, table3Name)
}

func tableConfig15CreateTableWithChangeTracking(name string) string {
	s := `
resource "snowflake_table" "test_table" {
	database            = "terraform_test_database"
	schema              = "terraform_test_schema"
	name                = "%s"
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
	return fmt.Sprintf(s, name)
}

func TestAcc_TableDefaults(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableColumnWithDefaults(accName),
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
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.default.0.sequence", fmt.Sprintf(`"%v"."%v".%v`, acc.TestDatabaseName, acc.TestSchemaName, accName)),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableColumnWithoutDefaults(accName),
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
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.default.0.sequence", fmt.Sprintf(`"%v"."%v".%v`, acc.TestDatabaseName, acc.TestSchemaName, accName)),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
		},
	})
}

func tableColumnWithDefaults(name string) string {
	s := `
resource "snowflake_sequence" "test_seq" {
	database            = "terraform_test_database"
	schema              = "terraform_test_schema"
	name                = "%s"
}

resource "snowflake_table" "test_table" {
	database            = "terraform_test_database"
	schema              = "terraform_test_schema"
	name                = "%s"
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
	return fmt.Sprintf(s, name, name)
}

func tableColumnWithoutDefaults(name string) string {
	s := `
resource "snowflake_sequence" "test_seq" {
	database            = "terraform_test_database"
	schema              = "terraform_test_schema"
	name                = "%s"
}

resource "snowflake_table" "test_table" {
	database            = "terraform_test_database"
	schema              = "terraform_test_schema"
	name                = "%s"
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
	return fmt.Sprintf(s, name, name)
}

func TestAcc_TableTags(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tagName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tag2Name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableWithTags(accName, tagName, tag2Name),
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

func tableWithTags(name string, tagName string, tag2Name string) string {
	s := `
resource "snowflake_tag" "test_tag" {
	name     = "%[2]s"
	database = "terraform_test_database"
	schema   = "terraform_test_schema"
	allowed_values = ["%[1]s"]
	comment  = "Terraform acceptance test"
}

resource "snowflake_tag" "test2_tag" {
	name     = "%[3]s"
	database = "terraform_test_database"
	schema   = "terraform_test_schema"
	allowed_values = ["%[1]s"]
	comment  = "Terraform acceptance test"
}

resource "snowflake_table" "test_table" {
	database            = "terraform_test_database"
	schema              = "terraform_test_schema"
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
	return fmt.Sprintf(s, name, tagName, tag2Name)
}

func TestAcc_TableIdentity(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableColumnWithIdentityDefault(accName),
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
				Config: tableColumnWithIdentity(accName),
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

func tableColumnWithIdentityDefault(name string) string {
	s := `
resource "snowflake_sequence" "test_seq" {
	database            = "terraform_test_database"
	schema              = "terraform_test_schema"
	name                = "%s"
}

resource "snowflake_table" "test_table" {
	database            = "terraform_test_database"
	schema              = "terraform_test_schema"
	name                = "%s"
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
	return fmt.Sprintf(s, name, name)
}

func tableColumnWithIdentity(name string) string {
	s := `
resource "snowflake_sequence" "test_seq" {
	database            = "terraform_test_database"
	schema              = "terraform_test_schema"
	name                = "%s"
}

resource "snowflake_table" "test_table" {
	database            = "terraform_test_database"
	schema              = "terraform_test_schema"
	name                = "%s"
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
	return fmt.Sprintf(s, name, name)
}

func TestAcc_TableRename(t *testing.T) {
	oldTableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	newTableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableConfigWithName(oldTableName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", oldTableName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "change_tracking", "false"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARIANT"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "primary_key.0"),
				),
			},
			{
				Config: tableConfigWithName(newTableName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", newTableName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
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

func tableConfigWithName(tableName string) string {
	s := `
resource "snowflake_table" "test_table" {
	database = "terraform_test_database"
	schema   = "terraform_test_schema"
	name     = "%s"
	column {
		name = "column1"
		type = "VARIANT"
	}
}
`
	return fmt.Sprintf(s, tableName)
}
