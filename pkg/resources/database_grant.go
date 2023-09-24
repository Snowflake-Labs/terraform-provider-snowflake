package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validDatabasePrivileges = NewPrivilegeSet(
	privilegeCreateDatabaseRole,
	privilegeCreateSchema,
	privilegeImportedPrivileges,
	privilegeModify,
	privilegeMonitor,
	privilegeOwnership,
	privilegeReferenceUsage,
	privilegeUsage,
	privilegeAllPrivileges,
)

var databaseGrantSchema = map[string]*schema.Schema{
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the database. To grant all privileges, use the value `ALL PRIVILEGES`.",
		Default:      "USAGE",
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice(validDatabasePrivileges.ToList(), true),
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
	},
	"shares": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these shares.",
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
	"revert_ownership_to_role_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the role to revert ownership to on destroy. Has no effect unless `privilege` is set to `OWNERSHIP`",
		Default:     "",
		ValidateFunc: func(val interface{}, key string) ([]string, []error) {
			additionalCharsToIgnoreValidation := []string{".", " ", ":", "(", ")"}
			return sdk.ValidateIdentifier(val, additionalCharsToIgnoreValidation)
		},
	},
}

// DatabaseGrant returns a pointer to the resource representing a database grant.
func DatabaseGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateDatabaseGrant,
			Read:   ReadDatabaseGrant,
			Delete: DeleteDatabaseGrant,
			Update: UpdateDatabaseGrant,

			Schema: databaseGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 5 {
						return nil, fmt.Errorf("invalid ID specified: %v, expected database_name|privilege|with_grant_option|roles|shares", d.Id())
					}
					if err := d.Set("database_name", parts[0]); err != nil {
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
					if err := d.Set("shares", helpers.StringListToList(parts[4])); err != nil {
						return nil, err
					}
					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validDatabasePrivileges,
	}
}

// CreateDatabaseGrant implements schema.CreateFunc.
func CreateDatabaseGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	builder := snowflake.DatabaseGrant(databaseName)
	if err := createGenericGrant(d, meta, builder); err != nil {
		return fmt.Errorf("error creating database grant err = %w", err)
	}

	privilege := d.Get("privilege").(string)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())
	withGrantOption := d.Get("with_grant_option").(bool)
	grantID := helpers.EncodeSnowflakeID(databaseName, privilege, withGrantOption, roles, shares)
	d.SetId(grantID)

	return ReadDatabaseGrant(d, meta)
}

// ReadDatabaseGrant implements schema.ReadFunc.
func ReadDatabaseGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	privilege := d.Get("privilege").(string)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())
	withGrantOption := d.Get("with_grant_option").(bool)

	// IMPORTED PRIVILEGES is not a real resource but
	// it is needed to grant usage for the snowflake database to custom roles.
	// We have to set to usage to check Snowflake for the grant
	if privilege == "IMPORTED PRIVILEGES" {
		err := d.Set("privilege", "USAGE")
		if err != nil {
			return fmt.Errorf("error setting privilege to USAGE: %w", err)
		}
	}

	builder := snowflake.DatabaseGrant(databaseName)
	err := readGenericGrant(d, meta, databaseGrantSchema, builder, false, false, validDatabasePrivileges)
	if err != nil {
		return fmt.Errorf("error reading database grant: %w", err)
	}

	// Then set it back to imported privledges for Terraform to execute the grant.
	if privilege == "IMPORTED PRIVILEGES" {
		err := d.Set("privilege", "IMPORTED PRIVILEGES")
		if err != nil {
			return fmt.Errorf("error setting privilege to IMPORTED PRIVILEGES: %w", err)
		}
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, privilege, withGrantOption, roles, shares)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// DeleteDatabaseGrant implements schema.DeleteFunc.
func DeleteDatabaseGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	builder := snowflake.DatabaseGrant(databaseName)

	return deleteGenericGrant(d, meta, builder)
}

// UpdateDatabaseGrant implements schema.UpdateFunc.
func UpdateDatabaseGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update are roles or shares
	// if nothing changed, nothing to update and we're done
	if !d.HasChanges("roles", "shares") {
		return nil
	}

	rolesToAdd := []string{}
	rolesToRevoke := []string{}
	sharesToAdd := []string{}
	sharesToRevoke := []string{}
	if d.HasChange("roles") {
		rolesToAdd, rolesToRevoke = changeDiff(d, "roles")
	}
	if d.HasChange("shares") {
		sharesToAdd, sharesToRevoke = changeDiff(d, "shares")
	}

	databaseName := d.Get("database_name").(string)
	privilege := d.Get("privilege").(string)
	reversionRole := d.Get("revert_ownership_to_role_name").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	// create the builder
	builder := snowflake.DatabaseGrant(databaseName)

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta,
		builder,
		privilege,
		reversionRole,
		rolesToRevoke,
		sharesToRevoke,
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
		sharesToAdd,
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadDatabaseGrant(d, meta)
}
