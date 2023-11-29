package resources

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var unsafeExecuteSchema = map[string]*schema.Schema{
	"execute": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "SQL statement to execute.",
	},
	"revert": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "SQL statement to revert the execute statement. Invoked when resource is deleted.",
	},
	"read": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Optional SQL statement to do a read.",
	},
	"read_results": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of key-value maps retrieved after executing read query.",
		Elem:        &schema.Schema{Type: schema.TypeMap},
	},
}

func UnsafeExecute() *schema.Resource {
	return &schema.Resource{
		Create: ApplyUnsafeMigration,
		Read:   ReadUnsafeMigration,
		Delete: RevertUnsafeMigration,
		Update: schema.Noop,

		Schema: unsafeExecuteSchema,

		DeprecationMessage: "Experimental resource. Will be deleted in the upcoming versions. Use at your own risk.",
		Description:        "Experimental resource used for testing purposes only. Allows to execute ANY SQL statement.",
	}
}

func ReadUnsafeMigration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)

	if readStatement := d.Get("read").(string); readStatement != "" {
		rows, err := client.QueryUnsafe(ctx, readStatement)
		if err != nil {
			return err
		}
		allRows, err := unsafeExecuteProcessRows(rows)
		if err != nil {
			return err
		}
		log.Printf(`[DEBUG] SQL "%s" queried successfully, returned rows count: %d`, readStatement, len(allRows))
		err = d.Set("read_results", allRows)
		if err != nil {
			return err
		}
	}

	return nil
}

func unsafeExecuteProcessRows(rows *sql.Rows) ([]map[string]string, error) {
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	allRows := make([]map[string]string, 0)

	unsafeExecuteProcessResultSet := func(rows *sql.Rows, columnNames []string) error {
		for rows.Next() {
			row, err := unsafeExecuteProcessRow(rows, columnNames)
			if err != nil {
				return err
			}
			allRows = append(allRows, row)
		}
		return nil
	}

	err = unsafeExecuteProcessResultSet(rows, columnNames)
	if err != nil {
		return nil, err
	}
	for rows.NextResultSet() {
		err := unsafeExecuteProcessResultSet(rows, columnNames)
		if err != nil {
			return nil, err
		}
	}

	return allRows, nil
}

func unsafeExecuteProcessRow(rows *sql.Rows, columnNames []string) (map[string]string, error) {
	values := make([]any, len(columnNames))
	for i, _ := range values {
		values[i] = new(any)
	}

	err := rows.Scan(values...)
	if err != nil {
		return nil, err
	}

	row := make(map[string]string)
	for i, col := range columnNames {
		row[col] = fmt.Sprintf("%v", *values[i].(*interface{}))
	}
	return row, nil
}

func ApplyUnsafeMigration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)

	id, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}

	executeStatement := d.Get("execute").(string)
	_, err = client.ExecUnsafe(ctx, executeStatement)
	if err != nil {
		return err
	}

	d.SetId(id)
	log.Printf(`[DEBUG] SQL "%s" applied successfully\n`, executeStatement)

	return ReadUnsafeMigration(d, meta)
}

func RevertUnsafeMigration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)

	revertStatement := d.Get("revert").(string)
	_, err := client.ExecUnsafe(ctx, revertStatement)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf(`[DEBUG] SQL "%s" applied successfully\n`, revertStatement)

	return nil
}
