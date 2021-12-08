package resources

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var roleProperties = []string{"comment"}
var roleSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Optional: true,
		// TODO validation
	},
	"tag": tagReferenceSchema,
}

func Role() *schema.Resource {
	return &schema.Resource{
		Create: CreateRole,
		Read:   ReadRole,
		Delete: DeleteRole,
		Update: UpdateRole,

		Schema: roleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateRole(d *schema.ResourceData, meta interface{}) error {
	return CreateResource("role", roleProperties, roleSchema, snowflake.Role, ReadRole)(d, meta)
}

func ReadRole(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := d.Id()

	row := snowflake.QueryRow(db, fmt.Sprintf("SHOW ROLES LIKE '%s'", id))
	role, err := snowflake.ScanRole(row)
	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
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
	err = d.Set("comment", role.Comment.String)
	if err != nil {
		return err
	}

	return err
}

func UpdateRole(d *schema.ResourceData, meta interface{}) error {
	return UpdateResource("role", roleProperties, roleSchema, snowflake.Role, ReadRole)(d, meta)
}

func DeleteRole(d *schema.ResourceData, meta interface{}) error {
	return DeleteResource("role", snowflake.Role)(d, meta)
}
