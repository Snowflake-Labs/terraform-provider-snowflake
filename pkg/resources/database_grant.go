package resources

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

var validDatabasePrivileges = NewPrivilegeSet(
	privilegeCreateSchema,
	privilegeImportedPrivileges,
	privilegeModify,
	privilegeMonitor,
	privilegeOwnership,
	privilegeReferenceUsage,
	privilegeUsage,
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
		Description:  "The privilege to grant on the database.",
		Default:      "USAGE",
		ForceNew:     true,
		ValidateFunc: validation.ValidatePrivilege(validDatabasePrivileges.ToList(), true),
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		ForceNew:    true,
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
}

// DatabaseGrant returns a pointer to the resource representing a database grant
func DatabaseGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateDatabaseGrant,
			Read:   ReadDatabaseGrant,
			Delete: DeleteDatabaseGrant,
			Update: UpdateDatabaseGrant,

			Schema: databaseGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validDatabasePrivileges,
	}
}

// CreateDatabaseGrant implements schema.CreateFunc
func CreateDatabaseGrant(d *schema.ResourceData, meta interface{}) error {
	dbName := d.Get("database_name").(string)
	builder := snowflake.DatabaseGrant(dbName)
	priv := d.Get("privilege").(string)
	grantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return errors.Wrap(err, "error creating database grant")
	}

	grant := &grantID{
		ResourceName: dbName,
		Privilege:    priv,
		GrantOption:  grantOption,
		Roles:        roles,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return errors.Wrap(err, "error creating database grant")
	}
	d.SetId(dataIDInput)

	return ReadDatabaseGrant(d, meta)
}

// ReadDatabaseGrant implements schema.ReadFunc
func ReadDatabaseGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	err = d.Set("database_name", grantID.ResourceName)
	if err != nil {
		return err
	}
	err = d.Set("privilege", grantID.Privilege)
	if err != nil {
		return err
	}
	err = d.Set("with_grant_option", grantID.GrantOption)
	if err != nil {
		return err
	}

	// IMPORTED PRIVILEGES is not a real resource, so we can't actually verify
	// that it is still there. Just exit for now
	if grantID.Privilege == "IMPORTED PRIVILEGES" {
		return nil
	}

	builder := snowflake.DatabaseGrant(grantID.ResourceName)
	return readGenericGrant(d, meta, databaseGrantSchema, builder, false, validDatabasePrivileges)
}

// DeleteDatabaseGrant implements schema.DeleteFunc
func DeleteDatabaseGrant(d *schema.ResourceData, meta interface{}) error {
	dbName := d.Get("database_name").(string)
	builder := snowflake.DatabaseGrant(dbName)

	return deleteGenericGrant(d, meta, builder)
}

// UpdateDatabaseGrant implements schema.UpdateFunc
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

	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}

	// create the builder
	builder := snowflake.DatabaseGrant(grantID.ResourceName)

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta,
		builder,
		grantID.Privilege,
		rolesToRevoke,
		sharesToRevoke,
	); err != nil {
		return err
	}

	// then add
	if err := createGenericGrantRolesAndShares(
		meta,
		builder,
		grantID.Privilege,
		grantID.GrantOption,
		rolesToAdd,
		sharesToAdd,
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadDatabaseGrant(d, meta)
}
