package datasources

import (
	"context"
	"database/sql"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var sharesSchema = map[string]*schema.Schema{
	"pattern": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the command output by object name.",
	},
	"shares": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of all the shares available in the system.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Identifier for the share.",
				},
				"comment": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The comment on the share.",
				},
				"owner": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The owner of the share.",
				},
				"kind": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The kind of the share.",
				},
				"to": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "For the OUTBOUND share, list of consumers.",
					Elem:        schema.TypeString,
				},
			},
		},
	},
}

// Shares Snowflake Shares resource.
func Shares() *schema.Resource {
	return &schema.Resource{
		Read:   ReadShares,
		Schema: sharesSchema,
	}
}

// ReadShares Reads the database metadata information.
func ReadShares(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	d.SetId("shares_read")
	pattern := d.Get("pattern").(string)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	var opts sdk.ShowShareOptions
	if pattern != "" {
		opts.Like = &sdk.Like{
			Pattern: sdk.String(pattern),
		}
	}
	shares, err := client.Shares.Show(ctx, &opts)
	if err != nil {
		return err
	}
	sharesFlatten := []map[string]interface{}{}
	for _, share := range shares {
		m := map[string]interface{}{}
		m["name"] = share.Name.Name()
		m["comment"] = share.Comment
		m["owner"] = share.Owner
		m["kind"] = share.Kind
		var to []string
		for _, consumer := range share.To {
			to = append(to, consumer.Name())
		}
		m["to"] = to
		sharesFlatten = append(sharesFlatten, m)
	}

	if err := d.Set("shares", sharesFlatten); err != nil {
		return err
	}
	return nil
}
