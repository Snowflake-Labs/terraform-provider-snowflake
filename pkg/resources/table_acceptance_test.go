package resources_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/require"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TableWithSeparateDataRetentionObjectParameterWithoutLifecycle(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckTableDestroy,
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
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckTableDestroy,
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
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	table2Name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	table3Name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: tableConfig(accName, acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: tableConfig2(accName, acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: tableConfig3(table2Name, acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: tableConfig4(table2Name, acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: tableConfig5(table2Name, acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: tableConfig6(accName, acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: tableConfig7(accName, acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: tableConfig8(accName, acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: tableConfig9CreateTableWithColumnComment(table2Name, acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: tableConfig10AlterTableColumnComment(table2Name, acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: tableConfig11AlterTableAddColumnWithComment(table2Name, acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: tableConfig12CreateTableWithDataRetention(table3Name, acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: tableConfig13AlterTableDataRetention(table3Name, acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: tableConfig14AlterTableEnableChangeTracking(table3Name, acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: tableConfig15CreateTableWithChangeTracking(accName, acc.TestDatabaseName, acc.TestSchemaName),
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
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckTableDestroy,
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
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tagName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tag2Name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckTableDestroy,
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
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckTableDestroy,
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
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: tableColumnWithCollate(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
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
				Config: alterTableColumnWithCollate(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "4"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.collate", "en"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.collate", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.collate", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.name", "column4"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.collate", "utf8"),
				),
			},
			{
				Config:      alterTableColumnWithIncompatibleCollate(accName, acc.TestDatabaseName, acc.TestSchemaName),
				ExpectError: regexp.MustCompile("\"VARCHAR\\(200\\) COLLATE 'fr'\" because they have incompatible collations\\."),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "4"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.collate", "en"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.collate", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.2.collate", ""),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.name", "column4"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.3.collate", "utf8"),
				),
			},
		},
	})
}

func tableColumnWithCollate(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table" {
	name     = "%s"
	database = "%s"
	schema   = "%s"
	comment  = "Terraform acceptance test"

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
`
	return fmt.Sprintf(s, name, databaseName, schemaName)
}

func alterTableColumnWithCollate(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table" {
	name     = "%s"
	database = "%s"
	schema   = "%s"
	comment  = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARCHAR(200)"
		collate = "en"
	}
	column {
		name = "column2"
		type = "VARCHAR(200)"
		collate = ""
	}
	column {
		name = "column3"
		type = "VARCHAR(200)"
	}
	column {
		name = "column4"
		type = "VARCHAR"
		collate = "utf8"
	}
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName)
}

func alterTableColumnWithIncompatibleCollate(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table" {
	name     = "%s"
	database = "%s"
	schema   = "%s"
	comment  = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARCHAR(200)"
		collate = "fr"
	}
	column {
		name = "column2"
		type = "VARCHAR(200)"
		collate = ""
	}
	column {
		name = "column3"
		type = "VARCHAR(200)"
	}
	column {
		name = "column4"
		type = "VARCHAR"
		collate = "utf8"
	}
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName)
}

func TestAcc_TableRename(t *testing.T) {
	oldTableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	newTableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: tableConfigWithName(oldTableName, acc.TestDatabaseName, acc.TestSchemaName),
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
				Config: tableConfigWithName(newTableName, acc.TestDatabaseName, acc.TestSchemaName),
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

func tableConfigWithName(tableName string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test_table" {
	name     = "%s"
	database = "%s"
	schema   = "%s"
	column {
		name = "column1"
		type = "VARIANT"
	}
}
`
	return fmt.Sprintf(s, tableName, databaseName, schemaName)
}

func TestAcc_Table_MaskingPolicy(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: tableWithMaskingPolicy(accName, acc.TestDatabaseName, acc.TestSchemaName, "policy1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.masking_policy", sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, fmt.Sprintf("%s1", accName)).FullyQualifiedName()),
				),
			},
			// this step proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2186
			{
				Config: tableWithMaskingPolicy(accName, acc.TestDatabaseName, acc.TestSchemaName, "policy2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.masking_policy", sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, fmt.Sprintf("%s2", accName)).FullyQualifiedName()),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Table_DefaultDataRetentionTime(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, tableName)

	configWithDatabaseDataRetentionSet := func(databaseDataRetentionTime int) config.Variables {
		return config.Variables{
			"database":                     config.StringVariable(databaseName),
			"schema":                       config.StringVariable(schemaName),
			"table":                        config.StringVariable(tableName),
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
		CheckDestroy: testAccCheckTableDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithDatabaseDataRetentionSet"),
				ConfigVariables: configWithDatabaseDataRetentionSet(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "-1"),
					checkDatabaseSchemaAndTableDataRetentionTime(id, 5, 5, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithSchemaDataRetentionSet"),
				ConfigVariables: configWithSchemaDataRetentionSet(5, 10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "-1"),
					checkDatabaseSchemaAndTableDataRetentionTime(id, 5, 10, 10),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(10, 3, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "5"),
					checkDatabaseSchemaAndTableDataRetentionTime(id, 10, 3, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(10, 3, 15),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "15"),
					checkDatabaseSchemaAndTableDataRetentionTime(id, 10, 3, 15),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithSchemaDataRetentionSet"),
				ConfigVariables: configWithSchemaDataRetentionSet(10, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "-1"),
					checkDatabaseSchemaAndTableDataRetentionTime(id, 10, 3, 3),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithDatabaseDataRetentionSet"),
				ConfigVariables: configWithDatabaseDataRetentionSet(10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "-1"),
					checkDatabaseSchemaAndTableDataRetentionTime(id, 10, 10, 10),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(10, 5, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "0"),
					checkDatabaseSchemaAndTableDataRetentionTime(id, 10, 5, 0),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(10, 5, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "3"),
					checkDatabaseSchemaAndTableDataRetentionTime(id, 10, 5, 3),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2356 issue is fixed.
func TestAcc_Table_DefaultDataRetentionTime_SetOutsideOfTerraform(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, tableName)

	configWithDatabaseDataRetentionSet := func(databaseDataRetentionTime int) config.Variables {
		return config.Variables{
			"database":                     config.StringVariable(databaseName),
			"schema":                       config.StringVariable(schemaName),
			"table":                        config.StringVariable(tableName),
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
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithDatabaseDataRetentionSet"),
				ConfigVariables: configWithDatabaseDataRetentionSet(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "-1"),
					checkDatabaseSchemaAndTableDataRetentionTime(id, 5, 5, 5),
				),
			},
			{
				PreConfig:       setTableDataRetentionTime(t, id, 20),
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithDatabaseDataRetentionSet"),
				ConfigVariables: configWithDatabaseDataRetentionSet(5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "-1"),
					checkDatabaseSchemaAndTableDataRetentionTime(id, 5, 5, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(5, 10, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "3"),
					checkDatabaseSchemaAndTableDataRetentionTime(id, 5, 10, 3),
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
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, tableName)

	configWithDatabaseDataRetentionSet := func(databaseDataRetentionTime int) config.Variables {
		return config.Variables{
			"database":                     config.StringVariable(databaseName),
			"schema":                       config.StringVariable(schemaName),
			"table":                        config.StringVariable(tableName),
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
		CheckDestroy: testAccCheckTableDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(10, 3, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "5"),
					checkDatabaseSchemaAndTableDataRetentionTime(id, 10, 3, 5),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(10, 3, -1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "-1"),
					checkDatabaseSchemaAndTableDataRetentionTime(id, 10, 3, 3),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithSchemaDataRetentionSet"),
				ConfigVariables: configWithSchemaDataRetentionSet(10, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "-1"),
					checkDatabaseSchemaAndTableDataRetentionTime(id, 10, 3, 3),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(10, 3, -1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "-1"),
					checkDatabaseSchemaAndTableDataRetentionTime(id, 10, 3, 3),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Table_DefaultDataRetentionTime/WithTableDataRetentionSet"),
				ConfigVariables: configWithTableDataRetentionSet(10, 3, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test", "data_retention_time_in_days", "5"),
					checkDatabaseSchemaAndTableDataRetentionTime(id, 10, 3, 5),
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
		masking_policy = snowflake_masking_policy.%[4]s.qualified_name
	}
}
`
	return fmt.Sprintf(s, name, databaseName, schemaName, policy)
}

func testAccCheckTableDestroy(s *terraform.State) error {
	client := acc.TestAccProvider.Meta().(*provider.Context).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_table" {
			continue
		}
		ctx := context.Background()
		id := sdk.NewSchemaObjectIdentifier(rs.Primary.Attributes["database"], rs.Primary.Attributes["schema"], rs.Primary.Attributes["name"])
		existingTable, err := client.Tables.ShowByID(ctx, id)
		if err == nil {
			return fmt.Errorf("table %v still exists", existingTable.ID().FullyQualifiedName())
		}
	}
	return nil
}

// proves issues https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2110 and https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2495
func TestAcc_Table_ClusterBy(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckTableDestroy,
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

func checkDatabaseSchemaAndTableDataRetentionTime(id sdk.SchemaObjectIdentifier, expectedDatabaseRetentionDays int, expectedSchemaRetentionDays int, expectedTableRetentionsDays int) func(state *terraform.State) error {
	return func(state *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()

		database, err := client.Databases.ShowByID(ctx, sdk.NewAccountObjectIdentifier(id.DatabaseName()))
		if err != nil {
			return err
		}

		if database.RetentionTime != expectedDatabaseRetentionDays {
			return fmt.Errorf("invalid database retention time, expected: %d, got: %d", expectedDatabaseRetentionDays, database.RetentionTime)
		}

		s, err := client.Schemas.ShowByID(ctx, sdk.NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName()))
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

func setTableDataRetentionTime(t *testing.T, id sdk.SchemaObjectIdentifier, days int) func() {
	t.Helper()

	return func() {
		client, err := sdk.NewDefaultClient()
		require.NoError(t, err)
		ctx := context.Background()

		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithSet(sdk.NewTableSetRequest().WithDataRetentionTimeInDays(sdk.Int(days))))
		require.NoError(t, err)
	}
}
