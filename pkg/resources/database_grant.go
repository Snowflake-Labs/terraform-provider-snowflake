package resources

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
					v, err := helpers.DecodeSnowflakeImportID(d.Id(), DatabaseGrantID{})
					if err != nil {
						return nil, err
					}
					id := v.(DatabaseGrantID)
					err = d.Set("database_name", id.DatabaseName)
					if err != nil {
						return nil, err
					}
					err = d.Set("privilege", id.Privilege)
					if err != nil {
						return nil, err
					}
					err = d.Set("roles", id.Roles)
					if err != nil {
						return nil, err
					}
					err = d.Set("shares", id.Shares)
					if err != nil {
						return nil, err
					}
					err = d.Set("with_grant_option", id.WithGrantOption)
					if err != nil {
						return nil, err
					}
					d.SetId(helpers.RandomSnowflakeID())
					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validDatabasePrivileges,
	}
}

type DatabaseGrantID struct {
	DatabaseName    string   `tf:"database_name"`
	Privilege       string   `tf:"privilege"`
	Roles           []string `tf:"roles"`
	Shares          []string `tf:"shares"`
	WithGrantOption bool     `tf:"with_grant_option"`
}

// CreateDatabaseGrant implements schema.CreateFunc.
func CreateDatabaseGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	builder := snowflake.DatabaseGrant(databaseName)
	if err := createGenericGrant(d, meta, builder); err != nil {
		return fmt.Errorf("error creating database grant err = %w", err)
	}

	d.SetId(helpers.RandomSnowflakeID())

	return ReadDatabaseGrant(d, meta)
}

// ReadDatabaseGrant implements schema.ReadFunc.
func ReadDatabaseGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	builder := snowflake.DatabaseGrant(databaseName)
	return readGenericGrant(d, meta, databaseGrantSchema, builder, false, validDatabasePrivileges)
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
	withGrantOption := d.Get("with_grant_option").(bool)
	// create the builder
	builder := snowflake.DatabaseGrant(databaseName)

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta,
		builder,
		privilege,
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
