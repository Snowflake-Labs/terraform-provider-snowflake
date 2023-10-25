// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var userOwnershipGrantSchema = map[string]*schema.Schema{
	"on_user_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the user ownership is granted on.",
	},
	"to_role_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the role to grant ownership. Please ensure that the role that terraform is using is granted access.",
		ValidateFunc: func(val interface{}, key string) ([]string, []error) {
			additionalCharsToIgnoreValidation := []string{".", " ", ":", "(", ")"}
			return sdk.ValidateIdentifier(val, additionalCharsToIgnoreValidation)
		},
	},
	"current_grants": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies whether to remove or transfer all existing outbound privileges on the object when ownership is transferred to a new role.",
		Default:     "COPY",
		ValidateFunc: validation.StringInSlice([]string{
			"COPY",
			"REVOKE",
		}, true),
	},
	"revert_ownership_to_role_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the role to revert ownership to on destroy.",
		Default:     "ACCOUNTADMIN",
		ValidateFunc: func(val interface{}, key string) ([]string, []error) {
			additionalCharsToIgnoreValidation := []string{".", " ", ":", "(", ")"}
			return sdk.ValidateIdentifier(val, additionalCharsToIgnoreValidation)
		},
	},
}

func UserOwnershipGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateUserOwnershipGrant,
		Read:   ReadUserOwnershipGrant,
		Delete: DeleteUserOwnershipGrant,
		Update: UpdateUserOwnershipGrant,
		Schema: userOwnershipGrantSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func CreateUserOwnershipGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	user := d.Get("on_user_name").(string)
	role := d.Get("to_role_name").(string)
	currentGrants := d.Get("current_grants").(string)

	g := snowflake.NewUserOwnershipGrantBuilder(user, currentGrants)
	err := snowflake.Exec(db, g.Role(role).Grant())
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf(`%s|%s|%s`, user, role, currentGrants))

	return ReadUserOwnershipGrant(d, meta)
}

func ReadUserOwnershipGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	log.Println(d.Id())
	user := strings.Split(d.Id(), "|")[0]
	currentGrants := strings.Split(d.Id(), "|")[2]

	stmt := fmt.Sprintf("SHOW USERS LIKE '%s'", user)
	row := snowflake.QueryRow(db, stmt)

	grant, err := snowflake.ScanUserOwnershipGrant(row)
	if errors.Is(err, sql.ErrNoRows) {
		// If not found, mark resource to be removed from state file during apply or refresh
		log.Printf("[DEBUG] user (%s) not found", d.Id())
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	if user != grant.Name.String {
		return fmt.Errorf("no user found like '%s'", user)
	}

	grant.Name.String = strings.TrimPrefix(grant.Name.String, `"`)
	grant.Name.String = strings.TrimSuffix(grant.Name.String, `"`)
	err = d.Set("on_user_name", grant.Name.String)
	if err != nil {
		return err
	}

	grant.Owner.String = strings.TrimPrefix(grant.Owner.String, `"`)
	grant.Owner.String = strings.TrimSuffix(grant.Owner.String, `"`)
	err = d.Set("to_role_name", grant.Owner.String)
	if err != nil {
		return err
	}

	err = d.Set("current_grants", currentGrants)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserOwnershipGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	user := d.Get("on_user_name").(string)
	role := d.Get("to_role_name").(string)
	currentGrants := d.Get("current_grants").(string)

	d.SetId(fmt.Sprintf(`%s|%s|%s`, user, role, currentGrants))

	g := snowflake.NewUserOwnershipGrantBuilder(user, currentGrants)
	err := snowflake.Exec(db, g.Role(role).Grant())
	if err != nil {
		return err
	}

	return ReadUserOwnershipGrant(d, meta)
}

func DeleteUserOwnershipGrant(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	user := d.Get("on_user_name").(string)
	currentGrants := d.Get("current_grants").(string)
	reversionRole := d.Get("revert_ownership_to_role_name").(string)

	g := snowflake.NewUserOwnershipGrantBuilder(user, currentGrants)
	err := snowflake.Exec(db, g.Role(reversionRole).Revoke())
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
