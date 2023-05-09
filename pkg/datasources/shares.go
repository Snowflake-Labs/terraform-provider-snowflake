package datasources

import (
	"database/sql"
	"errors"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
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
					Type:        schema.TypeString,
					Computed:    true,
					Description: "For the OUTBOUND share, list of consumers.",
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
	sharePattern := d.Get(pattern).(string)

	listShares, err := snowflake.ListShares(db, sharePattern)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("[DEBUG] no shares found in account (%s)", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Println("[DEBUG] failed to list shares")
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] list shares: %v", listShares)

	shares := []map[string]interface{}{}
	for _, share := range listShares {
		shareMap := map[string]interface{}{}
		if !share.Name.Valid {
			continue
		}
		shareMap["name"] = share.Name.String
		shareMap["comment"] = share.Comment.String
		shareMap["owner"] = share.Owner.String
		shareMap["kind"] = share.Kind.String
		shareMap["to"] = share.To.String
		shares = append(shares, shareMap)
	}

	err = d.Set("shares", shares)
	return err
}
