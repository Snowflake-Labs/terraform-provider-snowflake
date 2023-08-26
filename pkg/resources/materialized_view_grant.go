package resources

import (
	"context"
	"errors"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

/*
NewPrivilegeSet creates a set of privileges that are allowed
They are used for validation in the schema object below.
*/

var validMaterializedViewPrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeReferences,
	privilegeSelect,
	privilegeAllPrivileges,
)

// The schema holds the resource variables that can be provided in the Terraform.
var materializedViewGrantSchema = map[string]*schema.Schema{
	"materialized_view_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the materialized view on which to grant privileges immediately (only valid if on_future and on_all are false).",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema containing the current or future materialized views on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future materialized views on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future materialized view. To grant all privileges, use the value `ALL PRIVILEGES`",
		Default:      "SELECT",
		ValidateFunc: validation.StringInSlice(validMaterializedViewPrivileges.ToList(), true),
		ForceNew:     true,
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
		Description: "Grants privilege to these shares (only valid if on_future and on_all are false).",
	},
	"on_future": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future materialized views in the given schema. When this is true and no schema_name is provided apply this grant on all future materialized views in the given database. The materialized_view_name and shares fields must be unset in order to use on_future. Cannot be used together with on_all.",
		Default:     false,
		ForceNew:    true,
	},
	"on_all": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all materialized views in the given schema. When this is true and no schema_name is provided apply this grant on all materialized views in the given database. The materialized_view_name and shares fields must be unset in order to use on_all. Cannot be used together with on_future.",
		Default:     false,
		ForceNew:    true,
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
		ForceNew:    true,
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

// MaterializedViewGrant returns a pointer to the resource representing a view grant.
func MaterializedViewGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create:             CreateMaterializedViewGrant,
			Read:               ReadMaterializedViewGrant,
			Delete:             DeleteMaterializedViewGrant,
			Update:             UpdateMaterializedViewGrant,
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_grant_privileges_to_role instead.",
			Schema:             materializedViewGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 9 {
						return nil, errors.New("invalid ID specified for materialized view grant. Expecting {database}|{schema}|{materialized view}|{privilege}|{with_grant_option}|{on_future}|{on_all}|{roles}|{shares}")
					}
					if err := d.Set("database_name", parts[0]); err != nil {
						return nil, err
					}
					if err := d.Set("schema_name", parts[1]); err != nil {
						return nil, err
					}
					if parts[2] != "" {
						if err := d.Set("materialized_view_name", parts[2]); err != nil {
							return nil, err
						}
					}
					if err := d.Set("privilege", parts[3]); err != nil {
						return nil, err
					}
					if err := d.Set("with_grant_option", helpers.StringToBool(parts[4])); err != nil {
						return nil, err
					}
					if err := d.Set("on_future", helpers.StringToBool(parts[5])); err != nil {
						return nil, err
					}
					if err := d.Set("on_all", helpers.StringToBool(parts[6])); err != nil {
						return nil, err
					}
					if err := d.Set("roles", helpers.StringListToList(parts[7])); err != nil {
						return nil, err
					}
					if err := d.Set("shares", helpers.StringListToList(parts[8])); err != nil {
						return nil, err
					}

					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validMaterializedViewPrivileges,
	}
}

// CreateMaterializedViewGrant implements schema.CreateFunc.
func CreateMaterializedViewGrant(d *schema.ResourceData, meta interface{}) error {
	materializedViewName := d.Get("materialized_view_name").(string)
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())

	if (schemaName == "") && !onFuture && !onAll {
		return errors.New("schema_name must be set unless on_future or on_all is true")
	}

	if (materializedViewName == "") && !onFuture && !onAll {
		return errors.New("materialized_view_name must be set unless on_future or on_all is true")
	}
	if (materializedViewName != "") && onFuture && onAll {
		return errors.New("materialized_view_name must be empty if on_future and on_all is true")
	}

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureMaterializedViewGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllMaterializedViewGrant(databaseName, schemaName)
	default:
		builder = snowflake.MaterializedViewGrant(databaseName, schemaName, materializedViewName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, materializedViewName, privilege, withGrantOption, onFuture, onAll, roles, shares)
	d.SetId(grantID)

	return ReadMaterializedViewGrant(d, meta)
}

// ReadMaterializedViewGrant implements schema.ReadFunc.
func ReadMaterializedViewGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	materializedViewName := d.Get("materialized_view_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	if materializedViewName == "" && !onFuture && !onAll {
		return errors.New("materialized_view_name must be set unless on_future or on_all is true")
	}
	if materializedViewName != "" && (onFuture || onAll) {
		return errors.New("materialized_view_name must be empty if on_future or on_all is true")
	}
	if onAll && onFuture {
		return errors.New("on_future and on_all cannot both be true")
	}
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureMaterializedViewGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllMaterializedViewGrant(databaseName, schemaName)
	default:
		builder = snowflake.MaterializedViewGrant(databaseName, schemaName, materializedViewName)
	}

	err := readGenericGrant(d, meta, materializedViewGrantSchema, builder, onFuture, onAll, validMaterializedViewPrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, materializedViewName, privilege, withGrantOption, onFuture, onAll, roles, shares)
	// if the ID is not in the new format, rewrite it
	if d.Id() != grantID {
		d.SetId(grantID)
	}
	return nil
}

// DeleteMaterializedViewGrant implements schema.DeleteFunc.
func DeleteMaterializedViewGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	materializedViewName := d.Get("materialized_view_name").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureMaterializedViewGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllMaterializedViewGrant(databaseName, schemaName)
	default:
		builder = snowflake.MaterializedViewGrant(databaseName, schemaName, materializedViewName)
	}
	return deleteGenericGrant(d, meta, builder)
}

// UpdateMaterializedViewGrant implements schema.UpdateFunc.
func UpdateMaterializedViewGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update are roles or shares
	// if nothing changed, nothing to update, and we're done
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
	schemaName := d.Get("schema_name").(string)
	materializedViewName := d.Get("materialized_view_name").(string)
	privilege := d.Get("privilege").(string)
	reversionRole := d.Get("revert_ownership_to_role_name").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)

	// create the builder
	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureMaterializedViewGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllMaterializedViewGrant(databaseName, schemaName)
	default:
		builder = snowflake.MaterializedViewGrant(databaseName, schemaName, materializedViewName)
	}

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta, builder, privilege, reversionRole, rolesToRevoke, sharesToRevoke,
	); err != nil {
		return err
	}
	// then add
	if err := createGenericGrantRolesAndShares(
		meta, builder, privilege, withGrantOption, rolesToAdd, sharesToAdd,
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadMaterializedViewGrant(d, meta)
}
