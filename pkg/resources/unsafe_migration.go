package resources

import (
	"context"
	"database/sql"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

var unsafeMigrationSchema = map[string]*schema.Schema{
	"up": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "TODO",
	},
	"down": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "TODO",
	},
}

func UnsafeMigration() *schema.Resource {
	return &schema.Resource{
		Create: ApplyUnsafeMigration,
		Read:   schema.Noop,
		Delete: RevertUnsafeMigration,
		Update: schema.Noop,

		Schema: unsafeMigrationSchema,
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

	upStatement := d.Get("up").(string)
	_, err = client.ExecUnsafe(ctx, upStatement)
	if err != nil {
		return err
	}

	d.SetId(id)
	log.Printf(`[DEBUG] SQL "%s" applied successfully\n`, upStatement)

	return nil
}

func RevertUnsafeMigration(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	ctx := context.Background()
	client := sdk.NewClientFromDB(db)

	downStatement := d.Get("down").(string)
	_, err := client.ExecUnsafe(ctx, downStatement)
	if err != nil {
		return err
	}

	d.SetId("")
	log.Printf(`[DEBUG] SQL "%s" applied successfully\n`, downStatement)

	return nil
}
