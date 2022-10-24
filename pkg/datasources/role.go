package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var roleSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The role for which to return metadata.",
	},
	"comment": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The comment on the role",
	},
}

// Role Snowflake Role resource.
func Role() *schema.Resource {
	return &schema.Resource{
		Read:   ReadRole,
		Schema: roleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// ReadRole Reads the database metadata information.
func ReadRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	roleName := d.Get("name").(string)

	row := snowflake.QueryRow(db, fmt.Sprintf("SHOW ROLES LIKE '%s'", roleName))
	role, err := snowflake.ScanRole(row)

	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] role (%s) not found", roleName)
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	d.SetId(role.Name.String)
	err = d.Set("name", role.Name.String)
	if err != nil {
		return err
	}
	err = d.Set("comment", role.Comment.String)
	if err != nil {
		return err
	}

	return nil
}
