package resources

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
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
		Delete: DeleteResource("role", snowflake.Role),
		Update: UpdateResource("role", roleProperties, roleSchema, snowflake.Role, ReadRole),

		Schema: roleSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func CreateRole(data *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	name := data.Get("name").(string)

	var sb strings.Builder

	_, err := sb.WriteString(fmt.Sprintf(`CREATE ROLE "%s"`, name))
	if err != nil {
		return err
	}

	for _, field := range roleProperties {
		log.Printf("prop %s", field)
		val, ok := data.GetOk(field)
		log.Printf("val, ok %#v, %#v", ok, val)
		if ok {
			valStr := val.(string)
			_, e := sb.WriteString(fmt.Sprintf(" %s='%s'", strings.ToUpper(field), snowflake.EscapeString(valStr)))
			if e != nil {
				return e
			}
		}
	}
	err = DBExec(db, sb.String())

	if err != nil {
		return errors.Wrap(err, "error creating role")
	}

	data.SetId(name)

	return ReadRole(data, meta)
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
	db := meta.(*sql.DB)
	if data.HasChange("name") {
		data.Partial(true)
		// I wish this could be done on one line.
		oldNameI, newNameI := data.GetChange("name")
		oldName := oldNameI.(string)
		newName := newNameI.(string)

		err := DBExec(db, `ALTER ROLE "%s" RENAME TO "%s"`, oldName, newName)

		if err != nil {
			return errors.Wrapf(err, "error renaming role %s to %s", oldName, newName)
		}
		data.SetId(newName)
		data.SetPartial("name")
		data.Partial(false)
	}

	changes := []string{}

	for _, prop := range roleProperties {
		if data.HasChange(prop) {
			changes = append(changes, prop)
		}
	}
	if len(changes) > 0 {
		name := data.Get("name").(string)
		var sb strings.Builder
		_, err := sb.WriteString(fmt.Sprintf(`ALTER ROLE "%s" SET`, name))
		if err != nil {
			return err
		}

		for _, change := range changes {
			val := data.Get(change).(string)
			_, e := sb.WriteString(fmt.Sprintf(" %s='%s'",
				strings.ToUpper(change), snowflake.EscapeString(val)))
			if e != nil {
				return e
			}
		}

		err = DBExec(db, sb.String())
		if err != nil {
			return errors.Wrap(err, "error altering role")
		}
	}
	return ReadRole(data, meta)
}
