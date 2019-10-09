package resources

import (
	"database/sql"
	"fmt"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var roleProperties = []string{"comment"}
var roleSchema = map[string]*schema.Schema{
	"name": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	},
	"comment": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		// TODO validation
	},
}

func Role() *schema.Resource {
	return &schema.Resource{
		Create: CreateRole,
		Read:   ReadRole,
		Delete: DeleteRole,
		Update: UpdateRole,

		Schema: roleSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func CreateRole(data *schema.ResourceData, meta interface{}) error {
	return CreateResource("role", roleProperties, roleSchema, snowflake.Role, ReadRole)(data, meta)
}

func ReadRole(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	id := data.Id()

	row := db.QueryRow(fmt.Sprintf("SHOW ROLES LIKE '%s'", id))
	var createdOn, name, isDefault, isCurrent, isInherited, assignedToUsers, grantedToRoles, grantedRoles, owner, comment sql.NullString
	err := row.Scan(&createdOn, &name, &isDefault, &isCurrent, &isInherited, &assignedToUsers, &grantedToRoles, &grantedRoles, &owner, &comment)
	if err != nil {
		return err
	}

	err = data.Set("name", name.String)
	if err != nil {
		return err
	}
	err = data.Set("comment", comment.String)
	if err != nil {
		return err
	}

	return err
}

func UpdateRole(data *schema.ResourceData, meta interface{}) error {
	return UpdateResource("role", roleProperties, roleSchema, snowflake.Role, ReadRole)(data, meta)
}

func DeleteRole(data *schema.ResourceData, meta interface{}) error {
	return DeleteResource("role", snowflake.Role)(data, meta)
}
