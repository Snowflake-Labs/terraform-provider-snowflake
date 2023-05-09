package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validSchemaPrivileges = NewPrivilegeSet(
	privilegeAddSearchOptimization,
	privilegeCreateExternalTable,
	privilegeCreateFileFormat,
	privilegeCreateFunction,
	privilegeCreateMaskingPolicy,
	privilegeCreateMaterializedView,
	privilegeCreatePipe,
	privilegeCreateProcedure,
	privilegeCreateRowAccessPolicy,
	privilegeCreateSequence,
	privilegeCreateSessionPolicy,
	privilegeCreateStage,
	privilegeCreateStream,
	privilegeCreateTable,
	privilegeCreateTag,
	privilegeCreateTask,
	privilegeCreateTemporaryTable,
	privilegeCreateView,
	privilegeModify,
	privilegeMonitor,
	privilegeOwnership,
	privilegeUsage,
)

var schemaGrantSchema = map[string]*schema.Schema{
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the schema on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future schema. Note that if \"OWNERSHIP\" is specified, ensure that the role that terraform is using is granted access.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validSchemaPrivileges.ToList(), true),
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
		Description: "Grants privilege to these shares (only valid if on_future and on_all are unset).",
	},
	"on_future": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true, apply this grant on all future schemas in the given database. The schema_name and shares fields must be unset in order to use on_future. Cannot be used together with on_all.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"schema_name", "shares"},
	},
	"on_all": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true, apply this grant on all schemas in the given database. The schema_name and shares fields must be unset in order to use on_all. Cannot be used together with on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"schema_name", "shares"},
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

// SchemaGrant returns a pointer to the resource representing a view grant.
func SchemaGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateSchemaGrant,
			Read:   ReadSchemaGrant,
			Delete: DeleteSchemaGrant,
			Update: UpdateSchemaGrant,

			Schema: schemaGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 8 {
						return nil, fmt.Errorf("unexpected format of ID (%q), expected database_name|schema_name|privilege|with_grant_option|on_future|on_all|roles|shares", d.Id())
					}

					if err := d.Set("database_name", parts[0]); err != nil {
						return nil, err
					}
					if parts[1] != "" {
						if err := d.Set("schema_name", parts[1]); err != nil {
							return nil, err
						}
					}
					if err := d.Set("privilege", parts[2]); err != nil {
						return nil, err
					}
					if err := d.Set("with_grant_option", helpers.StringToBool(parts[3])); err != nil {
						return nil, err
					}
					if err := d.Set("on_future", helpers.StringToBool(parts[4])); err != nil {
						return nil, err
					}
					if err := d.Set("on_all", helpers.StringToBool(parts[5])); err != nil {
						return nil, err
					}
					if err := d.Set("roles", helpers.StringListToList(parts[6])); err != nil {
						return nil, err
					}
					if err := d.Set("shares", helpers.StringListToList(parts[7])); err != nil {
						return nil, err
					}
					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validSchemaPrivileges,
	}
}

// CreateSchemaGrant implements schema.CreateFunc.
func CreateSchemaGrant(d *schema.ResourceData, meta interface{}) error {
	schemaName := d.Get("schema_name").(string)
	databaseName := d.Get("database_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	if onFuture && onAll {
		return errors.New("on_future and on_all cannot both be true")
	}
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())

	if (schemaName == "") && !onFuture && !onAll {
		return errors.New("schema_name must be set unless on_future or on_all is true")
	}

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureSchemaGrant(databaseName)
	case onAll:
		builder = snowflake.AllSchemaGrant(databaseName)
	default:
		builder = snowflake.SchemaGrant(databaseName, schemaName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, privilege, withGrantOption, onFuture, onAll, roles, shares)
	d.SetId(grantID)

	return ReadSchemaGrant(d, meta)
}

// UpdateSchemaGrant implements schema.UpdateFunc.
func UpdateSchemaGrant(d *schema.ResourceData, meta interface{}) error {
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
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)

	// create the builder
	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureSchemaGrant(databaseName)
	case onAll:
		builder = snowflake.AllSchemaGrant(databaseName)
	default:
		builder = snowflake.SchemaGrant(databaseName, schemaName)
	}

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
	return ReadSchemaGrant(d, meta)
}

// ReadSchemaGrant implements schema.ReadFunc.
func ReadSchemaGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureSchemaGrant(databaseName)
	case onAll:
		builder = snowflake.AllSchemaGrant(databaseName)
	default:
		builder = snowflake.SchemaGrant(databaseName, schemaName)
	}

	err := readGenericGrant(d, meta, schemaGrantSchema, builder, onFuture, onAll, validSchemaPrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, privilege, withGrantOption, onFuture, onAll, roles, shares)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// DeleteSchemaGrant implements schema.DeleteFunc.
func DeleteSchemaGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureSchemaGrant(databaseName)
	case onAll:
		builder = snowflake.AllSchemaGrant(databaseName)
	default:
		builder = snowflake.SchemaGrant(databaseName, schemaName)
	}
	return deleteGenericGrant(d, meta, builder)
}
