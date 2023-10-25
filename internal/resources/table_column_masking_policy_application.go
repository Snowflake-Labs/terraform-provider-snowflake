// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	"database/sql"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var tableColumnMaskingPolicyApplicationSchema = map[string]*schema.Schema{
	"table": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The fully qualified name (`database.schema.table`) of the table to apply the masking policy to.",
	},
	"column": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The column to apply the masking policy to.",
	},
	"masking_policy": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Fully qualified name (`database.schema.policyname`) of the policy to apply.",
	},
}

func TableColumnMaskingPolicyApplication() *schema.Resource {
	return &schema.Resource{
		Description: "Applies a masking policy to a table column.",
		Create:      CreateTableColumnMaskingPolicyApplication,
		Read:        ReadTableColumnMaskingPolicyApplication,
		Delete:      DeleteTableColumnMaskingPolicyApplication,

		Schema: tableColumnMaskingPolicyApplicationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateTableColumnMaskingPolicyApplication implements schema.CreateFunc.
func CreateTableColumnMaskingPolicyApplication(d *schema.ResourceData, meta interface{}) error {
	manager := snowflake.NewTableColumnMaskingPolicyApplicationManager()

	input := &snowflake.TableColumnMaskingPolicyApplicationCreateInput{
		TableColumnMaskingPolicyApplication: snowflake.TableColumnMaskingPolicyApplication{
			Table:         snowflake.SchemaObjectIdentifierFromQualifiedName(d.Get("table").(string)),
			Column:        d.Get("column").(string),
			MaskingPolicy: snowflake.SchemaObjectIdentifierFromQualifiedName(d.Get("masking_policy").(string)),
		},
	}

	stmt := manager.Create(input)

	db := meta.(*sql.DB)
	_, err := db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("error applying masking policy: %w", err)
	}

	d.SetId(TableColumnMaskingPolicyApplicationID(&input.TableColumnMaskingPolicyApplication))

	return nil
}

// ReadTableColumnMaskingPolicyApplication implements schema.ReadFunc.
func ReadTableColumnMaskingPolicyApplication(d *schema.ResourceData, meta interface{}) error {
	manager := snowflake.NewTableColumnMaskingPolicyApplicationManager()

	table, column := TableColumnMaskingPolicyApplicationIdentifier(d.Id())

	if err := d.Set("table", table.QualifiedName()); err != nil {
		return fmt.Errorf("error setting table: %w", err)
	}
	if err := d.Set("column", column); err != nil {
		return fmt.Errorf("error setting column: %w", err)
	}

	input := &snowflake.TableColumnMaskingPolicyApplicationReadInput{
		Table:  table,
		Column: column,
	}

	stmt := manager.Read(input)

	db := meta.(*sql.DB)
	rows, err := db.Query(stmt)
	if err != nil {
		return fmt.Errorf("error querying password policy: %w", err)
	}

	defer rows.Close()
	maskingPolicy, err := manager.Parse(rows, column)
	if err != nil {
		return fmt.Errorf("failed to parse result of describe: %w", err)
	}

	if err = d.Set("masking_policy", maskingPolicy); err != nil {
		return fmt.Errorf("error setting masking_policy: %w", err)
	}

	return nil
}

// DeleteTableColumnMaskingPolicyApplication implements schema.DeleteFunc.
func DeleteTableColumnMaskingPolicyApplication(d *schema.ResourceData, meta interface{}) error {
	manager := snowflake.NewTableColumnMaskingPolicyApplicationManager()

	input := &snowflake.TableColumnMaskingPolicyApplicationDeleteInput{
		TableColumn: snowflake.TableColumn{
			Table:  snowflake.SchemaObjectIdentifierFromQualifiedName(d.Get("table").(string)),
			Column: d.Get("column").(string),
		},
	}

	stmt := manager.Delete(input)

	db := meta.(*sql.DB)
	_, err := db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("error executing drop statement: %w", err)
	}

	return nil
}

func TableColumnMaskingPolicyApplicationID(mpa *snowflake.TableColumnMaskingPolicyApplication) string {
	identifier := snowflake.ColumnIdentifier{
		Database:   mpa.Table.Database,
		Schema:     mpa.Table.Schema,
		ObjectName: mpa.Table.ObjectName,
		Column:     mpa.Column,
	}
	return identifier.QualifiedName()
}

func TableColumnMaskingPolicyApplicationIdentifier(id string) (table *snowflake.SchemaObjectIdentifier, column string) {
	columnIdentifier := snowflake.ColumnIdentifierFromQualifiedName(id)
	return &snowflake.SchemaObjectIdentifier{
		Database:   columnIdentifier.Database,
		Schema:     columnIdentifier.Schema,
		ObjectName: columnIdentifier.ObjectName,
	}, columnIdentifier.Column
}
