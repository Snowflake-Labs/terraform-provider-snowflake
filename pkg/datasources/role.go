package datasources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var roleSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The role for which to return metadata.",
	},
}

// Role Snowflake Role resource
func Role() *schema.Resource {
	return &schema.Resource{
		Read:   ReadRole,
		Schema: roleSchema,
	}
}

// ReadRole Reads the database metadata information
func ReadRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	row := snowflake.QueryRow(db, fmt.Sprintf("SHOW ROLES LIKE '%s'", id))
	role, err := snowflake.ScanRole(row)

	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] role (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	err = d.Set("name", role.Name.String)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return err
}
