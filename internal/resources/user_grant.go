// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validUserPrivileges = NewPrivilegeSet(
	privilegeMonitor,
	privilegeAllPrivileges,
)

var userGrantSchema = map[string]*schema.Schema{
	"user_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the user on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  "The privilege to grant on the user. To grant all privileges, use the value `ALL PRIVILEGES`.",
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice(validUserPrivileges.ToList(), true),
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
		Default:     false,
		ForceNew:    true,
	},
	"enable_multiple_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.",
		Default:     false,
	},
}

// UserGrant returns a pointer to the resource representing a user grant.
func UserGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create:             CreateUserGrant,
			Read:               ReadUserGrant,
			Delete:             DeleteUserGrant,
			Update:             UpdateUserGrant,
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_grant_privileges_to_role instead.",
			Schema:             userGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 4 {
						return nil, fmt.Errorf("unexpected format of ID (%q), expected user-name|privilege|with_grant_option|roles", d.Id())
					}
					if err := d.Set("user_name", parts[0]); err != nil {
						return nil, err
					}
					if err := d.Set("privilege", parts[1]); err != nil {
						return nil, err
					}
					if err := d.Set("with_grant_option", helpers.StringToBool(parts[2])); err != nil {
						return nil, err
					}
					if err := d.Set("roles", helpers.StringListToList(parts[3])); err != nil {
						return nil, err
					}
					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validUserPrivileges,
	}
}

// CreateUserGrant implements schema.CreateFunc.
func CreateUserGrant(d *schema.ResourceData, meta interface{}) error {
	userName := d.Get("user_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	builder := snowflake.UserGrant(userName)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}
	grantID := helpers.EncodeSnowflakeID(userName, privilege, withGrantOption, roles)
	d.SetId(grantID)

	return ReadUserGrant(d, meta)
}

// ReadUserGrant implements schema.ReadFunc.
func ReadUserGrant(d *schema.ResourceData, meta interface{}) error {
	userName := d.Get("user_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.UserGrant(userName)

	err := readGenericGrant(d, meta, userGrantSchema, builder, false, false, validUserPrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(userName, privilege, withGrantOption, roles)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// DeleteUserGrant implements schema.DeleteFunc.
func DeleteUserGrant(d *schema.ResourceData, meta interface{}) error {
	userName := d.Get("user_name").(string)

	builder := snowflake.UserGrant(userName)

	return deleteGenericGrant(d, meta, builder)
}

// UpdateUserGrant implements schema.UpdateFunc.
func UpdateUserGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update is roles. if nothing changed,
	// nothing to update and we're done.
	if !d.HasChanges("roles") {
		return nil
	}

	rolesToAdd, rolesToRevoke := changeDiff(d, "roles")

	userName := d.Get("user_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)

	// create the builder
	builder := snowflake.UserGrant(userName)

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta,
		builder,
		privilege,
		"",
		rolesToRevoke,
		nil,
	); err != nil {
		return err
	}

	// then add
	if err := createGenericGrantRolesAndShares(
		meta,
		builder,
		privilege,
		withGrantOption,
		rolesToAdd,
		nil,
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadUserGrant(d, meta)
}
