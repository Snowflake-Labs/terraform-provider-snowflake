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
	role, err := snowflake.NewRoleBuilder(db, roleName).Show()

	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("[DEBUG] role (%s) not found", roleName)
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	d.SetId(role.Name.String)
	if err := d.Set("name", role.Name.String); err != nil {
		return err
	}
	if err := d.Set("comment", role.Comment.String); err != nil {
		return err
	}
	return nil
}
