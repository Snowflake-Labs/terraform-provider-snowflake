package resources

import (
	"context"
	"database/sql"
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
	"query": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Optional SQL statement to do a read.",
	},
	"query_results": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of key-value maps retrieved after executing read query.",
		Elem:        &schema.Schema{Type: schema.TypeMap},
	},
}

func UnsafeExecute() *schema.Resource {
	return &schema.Resource{
		Create: CreateUnsafeExecute,
		Read:   ReadUnsafeExecute,
		Delete: DeleteUnsafeExecute,
		Update: UpdateUnsafeExecute,

		Schema: unsafeExecuteSchema,

		DeprecationMessage: "Experimental resource. Will be deleted in the upcoming versions. Use at your own risk.",
		Description:        "Experimental resource used for testing purposes only. Allows to execute ANY SQL statement.",
	}
}

func ReadUnsafeExecute(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)

	readStatement := d.Get("query").(string)

	if readStatement == "" {
		err := d.Set("query_results", nil)
		if err != nil {
			return err
		}
	} else {
		rows, err := client.QueryUnsafe(ctx, readStatement)
		log.Printf(`[DEBUG] SQL query "%s" executed successfully, returned rows count: %d`, readStatement, len(rows))
		err = d.Set("query_results", rows)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateUnsafeExecute(d *schema.ResourceData, meta interface{}) error {
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

	return ReadUnsafeExecute(d, meta)
}

func DeleteUnsafeExecute(d *schema.ResourceData, meta interface{}) error {
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

func UpdateUnsafeExecute(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("query") {
		return ReadUnsafeExecute(d, meta)
	}
	return nil
}
