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
		Description: "TODO",
	},
	"revert": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "TODO",
	},
}

func UnsafeExecute() *schema.Resource {
	return &schema.Resource{
		Create: ApplyUnsafeMigration,
		Read:   schema.Noop,
		Delete: RevertUnsafeMigration,
		Update: schema.Noop,

		Schema: unsafeExecuteSchema,

		DeprecationMessage: "Experimental resource. Will be deleted in the upcoming versions. Use on your own risk.",
	}
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

	return nil
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
