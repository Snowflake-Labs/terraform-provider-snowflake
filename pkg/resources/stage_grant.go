package resources

import (
	"context"
	"errors"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validStagePrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
	privilegeAllPrivileges,
	// These privileges are only valid for internal stages
	privilegeRead,
	privilegeWrite,
)

var stageGrantSchema = map[string]*schema.Schema{
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current stage on which to grant privileges.",
		ForceNew:    true,
	},
	"enable_multiple_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.",
		Default:     false,
		ForceNew:    true,
	},
	"on_future": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future stages in the given schema. When this is true and no schema_name is provided apply this grant on all future stages in the given database. The stage_name field must be unset in order to use on_future. Cannot be used together with on_all.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"stage_name"},
	},
	"on_all": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all stages in the given schema. When this is true and no schema_name is provided apply this grant on all stages in the given database. The stage_name field must be unset in order to use on_all. Cannot be used together with on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"stage_name"},
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the stage. To grant all privileges, use the value `ALL PRIVILEGES`.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validStagePrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Required:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Grants privilege to these roles.",
	},
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema containing the current stage on which to grant privileges.",
		ForceNew:    true,
	},
	"stage_name": {
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "The name of the stage on which to grant privilege (only valid if on_future and on_all are false).",
		ForceNew:      true,
		ConflictsWith: []string{"on_future"},
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
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

// StageGrant returns a pointer to the resource representing a stage grant.
func StageGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create:             CreateStageGrant,
			Read:               ReadStageGrant,
			Delete:             DeleteStageGrant,
			Update:             UpdateStageGrant,
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_grant_privileges_to_role instead.",
			Schema:             stageGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 8 {
						return nil, errors.New("incorrect ID format (expecting database_name|schema_name|stage_name|privilege|with_grant_option|on_future|on_all|roles")
					}
					if err := d.Set("database_name", parts[0]); err != nil {
						return nil, err
					}
					if err := d.Set("schema_name", parts[1]); err != nil {
						return nil, err
					}
					if parts[2] != "" {
						if err := d.Set("stage_name", parts[2]); err != nil {
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
					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validStagePrivileges,
	}
}

// CreateStageGrant implements schema.CreateFunc.
func CreateStageGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	if onFuture && onAll {
		return errors.New("on_future and on_all cannot both be true")
	}
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)

	var schemaName string
	if name, ok := d.GetOk("schema_name"); ok {
		schemaName = name.(string)
	}
	stageName := d.Get("stage_name").(string)
	if (schemaName == "") && !onFuture && !onAll {
		return errors.New("schema_name must be set unless on_future or on_all is true")
	}
	if (stageName == "") && !onFuture && !onAll {
		return errors.New("stage_name must be set unless on_future or on_all is true")
	}

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureStageGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllStageGrant(databaseName, schemaName)
	default:
		builder = snowflake.StageGrant(databaseName, schemaName, stageName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, stageName, privilege, withGrantOption, onFuture, onAll, roles)
	d.SetId(grantID)

	return ReadStageGrant(d, meta)
}

// ReadStageGrant implements schema.ReadFunc.
func ReadStageGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	stageName := d.Get("stage_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureStageGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllStageGrant(databaseName, schemaName)
	default:
		builder = snowflake.StageGrant(databaseName, schemaName, stageName)
	}

	err := readGenericGrant(d, meta, stageGrantSchema, builder, onFuture, onAll, validStagePrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, stageName, privilege, withGrantOption, onFuture, onAll, roles)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// UpdateStageGrant implements schema.UpdateFunc.
func UpdateStageGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update are roles or shares
	// if nothing changed, nothing to update, and we're done
	if !d.HasChanges("roles") {
		return nil
	}

	rolesToAdd := []string{}
	rolesToRevoke := []string{}

	if d.HasChange("roles") {
		rolesToAdd, rolesToRevoke = changeDiff(d, "roles")
	}

	var builder snowflake.GrantBuilder
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	stageName := d.Get("stage_name").(string)
	privilege := d.Get("privilege").(string)
	reversionRole := d.Get("revert_ownership_to_role_name").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	switch {
	case onFuture:
		builder = snowflake.FutureStageGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllStageGrant(databaseName, schemaName)
	default:
		builder = snowflake.StageGrant(databaseName, schemaName, stageName)
	}

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta, builder, privilege, reversionRole, rolesToRevoke, []string{},
	); err != nil {
		return err
	}
	// then add
	if err := createGenericGrantRolesAndShares(
		meta, builder, privilege, withGrantOption, rolesToAdd, []string{},
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadStageGrant(d, meta)
}

// DeleteStageGrant implements schema.DeleteFunc.
func DeleteStageGrant(d *schema.ResourceData, meta interface{}) error {
	var builder snowflake.GrantBuilder
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	stageName := d.Get("stage_name").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	switch {
	case onFuture:
		builder = snowflake.FutureStageGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllStageGrant(databaseName, schemaName)
	default:
		builder = snowflake.StageGrant(databaseName, schemaName, stageName)
	}

	return deleteGenericGrant(d, meta, builder)
}
