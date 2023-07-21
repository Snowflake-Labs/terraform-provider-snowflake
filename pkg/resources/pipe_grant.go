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

var validPipePrivileges = NewPrivilegeSet(
	privilegeMonitor,
	privilegeOperate,
	privilegeOwnership,
	privilegeAllPrivileges,
)

var pipeGrantSchema = map[string]*schema.Schema{
	"pipe_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the pipe on which to grant privileges immediately (only valid if on_future is false).",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema containing the current or future pipes on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future pipes on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future pipe. To grant all privileges, use the value `ALL PRIVILEGES`",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validPipePrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
	},
	"on_future": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future pipes in the given schema. When this is true and no schema_name is provided apply this grant on all future pipes in the given database. The pipe_name field must be unset in order to use on_future.",
		Default:     false,
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
			return snowflake.ValidateIdentifier(val, additionalCharsToIgnoreValidation)
		},
	},
}

// PipeGrant returns a pointer to the resource representing a pipe grant.
func PipeGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create:             CreatePipeGrant,
			Read:               ReadPipeGrant,
			Delete:             DeletePipeGrant,
			Update:             UpdatePipeGrant,
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_grant_privileges_to_role instead.",
			Schema:             pipeGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 7 {
						return nil, fmt.Errorf("pipe grant ID must be in the form database_name|schema_name|pipe_name|privilege|with_grant_option|on_future|roles")
					}
					if err := d.Set("database_name", parts[0]); err != nil {
						return nil, err
					}
					if err := d.Set("schema_name", parts[1]); err != nil {
						return nil, err
					}
					if parts[2] != "" {
						if err := d.Set("pipe_name", parts[2]); err != nil {
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
					if err := d.Set("roles", helpers.StringListToList(parts[6])); err != nil {
						return nil, err
					}
					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validPipePrivileges,
	}
}

// CreatePipeGrant implements schema.CreateFunc.
func CreatePipeGrant(d *schema.ResourceData, meta interface{}) error {
	pipeName := d.Get("pipe_name").(string)
	var schemaName string
	if name, ok := d.GetOk("schema_name"); ok {
		schemaName = name.(string)
	}

	databaseName := d.Get("database_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if (schemaName == "") && !onFuture {
		return errors.New("schema_name must be set unless on_future is true")
	}
	if (pipeName == "") && !onFuture {
		return errors.New("pipe_name must be set unless on_future is true")
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FuturePipeGrant(databaseName, schemaName)
	} else {
		builder = snowflake.PipeGrant(databaseName, schemaName, pipeName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, pipeName, privilege, withGrantOption, onFuture, roles)
	d.SetId(grantID)

	return ReadPipeGrant(d, meta)
}

// ReadPipeGrant implements schema.ReadFunc.
func ReadPipeGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	pipeName := d.Get("pipe_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	onFuture := d.Get("on_future").(bool)
	if pipeName == "" && !onFuture {
		return errors.New("pipe_name must be set unless on_future is true")
	}
	if pipeName != "" && onFuture {
		return errors.New("pipe_name must not be set when on_future is true")
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FuturePipeGrant(databaseName, schemaName)
	} else {
		builder = snowflake.PipeGrant(databaseName, schemaName, pipeName)
	}
	// TODO
	onAll := false

	err := readGenericGrant(d, meta, pipeGrantSchema, builder, onFuture, onAll, validPipePrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, pipeName, privilege, withGrantOption, onFuture, roles)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// DeletePipeGrant implements schema.DeleteFunc.
func DeletePipeGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	pipeName := d.Get("pipe_name").(string)
	onFuture := d.Get("on_future").(bool)

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FuturePipeGrant(databaseName, schemaName)
	} else {
		builder = snowflake.PipeGrant(databaseName, schemaName, pipeName)
	}
	return deleteGenericGrant(d, meta, builder)
}

// UpdatePipeGrant implements schema.UpdateFunc.
func UpdatePipeGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update are roles or shares
	// if nothing changed, nothing to update and we're done
	if !d.HasChanges("roles") {
		return nil
	}

	rolesToAdd := []string{}
	rolesToRevoke := []string{}

	if d.HasChange("roles") {
		rolesToAdd, rolesToRevoke = changeDiff(d, "roles")
	}
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	pipeName := d.Get("pipe_name").(string)
	onFuture := d.Get("on_future").(bool)
	privilege := d.Get("privilege").(string)
	reversionRole := d.Get("revert_ownership_to_role_name").(string)
	withGrantOption := d.Get("with_grant_option").(bool)

	// create the builder
	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FuturePipeGrant(databaseName, schemaName)
	} else {
		builder = snowflake.PipeGrant(databaseName, schemaName, pipeName)
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
	return ReadPipeGrant(d, meta)
}
