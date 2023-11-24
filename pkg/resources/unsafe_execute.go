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
}

func UnsafeExecute() *schema.Resource {
	return &schema.Resource{
		Create: ExecuteUnsafeSQLStatement,
		Read:   schema.Noop,
		Delete: RevertUnsafeSQLStatement,
		Update: schema.Noop,

		Schema: unsafeExecuteSchema,

		DeprecationMessage: "Experimental resource. Will be deleted in the upcoming versions. Use at your own risk.",
		Description:        "Experimental resource used for testing purposes only. Allows to execute ANY SQL statement.",
	}
}

func ExecuteUnsafeSQLStatement(d *schema.ResourceData, meta interface{}) error {
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

	return nil
}

func RevertUnsafeSQLStatement(d *schema.ResourceData, meta interface{}) error {
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
