// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package datasources

import (
	"database/sql"
	"errors"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var rolesSchema = map[string]*schema.Schema{
	"pattern": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the command output by object name.",
	},
	"roles": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of all the roles which you can view across your entire account, including the system-defined roles and any custom roles that exist.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Identifier for the role.",
				},
				"comment": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The comment on the role",
				},
				"owner": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The owner of the role",
				},
			},
		},
	},
}

// Roles Snowflake Roles resource.
func Roles() *schema.Resource {
	return &schema.Resource{
		Read:   ReadRoles,
		Schema: rolesSchema,
	}
}

// ReadRoles Reads the database metadata information.
func ReadRoles(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	d.SetId("roles_read")
	rolePattern := d.Get("pattern").(string)

	listRoles, err := snowflake.ListRoles(db, rolePattern)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("[DEBUG] no roles found in account (%s)", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		log.Println("[DEBUG] failed to list roles")
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] list roles: %v", listRoles)

	roles := []map[string]interface{}{}
	for _, role := range listRoles {
		roleMap := map[string]interface{}{}
		if !role.Name.Valid {
			continue
		}
		roleMap["name"] = role.Name.String
		roleMap["comment"] = role.Comment.String
		roleMap["owner"] = role.Owner.String
		roles = append(roles, roleMap)
	}

	if err := d.Set("roles", roles); err != nil {
		return err
	}
	return nil
}
