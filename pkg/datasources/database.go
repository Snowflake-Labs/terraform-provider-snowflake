package datasources

import (
	"context"
	"database/sql"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var databaseSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return its metadata.",
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_default": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"is_current": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"origin": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"retention_time": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"options": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

// Database the Snowflake Database resource.
func Database() *schema.Resource {
	return &schema.Resource{
		Read:   ReadDatabase,
		Schema: databaseSchema,
	}
}

// ReadDatabase read the database meta-data information.
func ReadDatabase(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	database, err := client.Databases.ShowByID(ctx, id)
	if err != nil {
		return err
	}
	d.SetId(database.Name)
	if err := d.Set("name", database.Name); err != nil {
		return err
	}
	if err := d.Set("comment", database.Comment); err != nil {
		return err
	}
	if err := d.Set("owner", database.Owner); err != nil {
		return err
	}
	if err := d.Set("is_default", database.IsDefault); err != nil {
		return err
	}
	if err := d.Set("is_current", database.IsCurrent); err != nil {
		return err
	}
	if err := d.Set("origin", database.Origin); err != nil {
		return err
	}
	if err := d.Set("retention_time", database.RetentionTime); err != nil {
		return err
	}
	if err := d.Set("created_on", database.CreatedOn.String()); err != nil {
		return err
	}
	if err := d.Set("options", database.Options); err != nil {
		return err
	}
	return nil
}
