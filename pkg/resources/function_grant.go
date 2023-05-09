package resources

import (
	"context"
	"errors"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validFunctionPrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
)

var functionGrantSchema = map[string]*schema.Schema{
	"argument_data_types": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "List of the argument data types for the function (must be present if function has arguments and function_name is present)",
		ForceNew:    true,
	},
	"enable_multiple_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.",
		Default:     false,
		ForceNew:    true,
	},
	"function_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the function on which to grant privileges immediately (only valid if on_future is false).",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future functions on which to grant privileges.",
		ForceNew:    true,
	},
	"on_future": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future functions in the given schema. When this is true and no schema_name is provided apply this grant on all future functions in the given database. The function_name, arguments, return_type, and shares fields must be unset in order to use on_future. Cannot be used together with on_all.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"function_name"},
	},
	"on_all": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all functions in the given schema. When this is true and no schema_name is provided apply this grant on all functions in the given database. The function_name, arguments, return_type, and shares fields must be unset in order to use on_all. Cannot be used together with on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"function_name"},
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future function. Must be one of `USAGE` or `OWNERSHIP`.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validFunctionPrivileges.ToList(), true),
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
		Description: "The name of the schema containing the current or future functions on which to grant privileges.",
		ForceNew:    true,
	},
	"shares": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these shares (only valid if on_future is false).",
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
		Default:     false,
		ForceNew:    true,
	},
}

// FunctionGrant returns a pointer to the resource representing a function grant.
func FunctionGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateFunctionGrant,
			Read:   ReadFunctionGrant,
			Delete: DeleteFunctionGrant,
			Update: UpdateFunctionGrant,

			Schema: functionGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 10 {
						return nil, errors.New("function grant ID must be specified as database_name|schema_name|function_name|argument_data_types|privilege|with_grant_option|on_future|roles|shares")
					}
					if err := d.Set("database_name", parts[0]); err != nil {
						return nil, err
					}
					if err := d.Set("schema_name", parts[1]); err != nil {
						return nil, err
					}
					if parts[2] != "" {
						if err := d.Set("function_name", parts[2]); err != nil {
							return nil, err
						}
					}
					if err := d.Set("argument_data_types", helpers.StringListToList(parts[3])); err != nil {
						return nil, err
					}
					if err := d.Set("privilege", parts[4]); err != nil {
						return nil, err
					}
					if err := d.Set("with_grant_option", helpers.StringToBool(parts[5])); err != nil {
						return nil, err
					}
					if err := d.Set("on_future", helpers.StringToBool(parts[6])); err != nil {
						return nil, err
					}
					if err := d.Set("on_all", helpers.StringToBool(parts[7])); err != nil {
						return nil, err
					}
					if err := d.Set("roles", helpers.StringListToList(parts[8])); err != nil {
						return nil, err
					}
					if err := d.Set("shares", helpers.StringListToList(parts[9])); err != nil {
						return nil, err
					}

					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validFunctionPrivileges,
	}
}

// CreateFunctionGrant implements schema.CreateFunc.
func CreateFunctionGrant(d *schema.ResourceData, meta interface{}) error {
	functionName := d.Get("function_name").(string)
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	var argumentDataTypes []string
	// support deprecated arguments
	if v, ok := d.GetOk("arguments"); ok {
		arguments := v.([]interface{})
		for _, argument := range arguments {
			argumentMap := argument.(map[string]interface{})
			argumentDataTypes = append(argumentDataTypes, argumentMap["type"].(string))
		}
	}
	if v, ok := d.GetOk("argument_data_types"); ok {
		argumentDataTypes = expandStringList(v.([]interface{}))
	}

	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())

	if (functionName == "") && !onFuture && !onAll {
		return errors.New("function_name must be set unless on_future or on_all is true")
	}
	if (schemaName == "") && !onFuture && !onAll {
		return errors.New("schema_name must be set unless on_future or on_all is true")
	}

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureFunctionGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllFunctionGrant(databaseName, schemaName)
	default:
		builder = snowflake.FunctionGrant(databaseName, schemaName, functionName, argumentDataTypes)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}
	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, functionName, argumentDataTypes, privilege, withGrantOption, onFuture, onAll, roles, shares)
	d.SetId(grantID)
	return ReadFunctionGrant(d, meta)
}

// ReadFunctionGrant implements schema.ReadFunc.
func ReadFunctionGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	functionName := d.Get("function_name").(string)
	argumentDataTypes := expandStringList(d.Get("argument_data_types").([]interface{}))
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureFunctionGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllFunctionGrant(databaseName, schemaName)
	default:
		builder = snowflake.FunctionGrant(databaseName, schemaName, functionName, argumentDataTypes)
	}

	err := readGenericGrant(d, meta, functionGrantSchema, builder, onFuture, onAll, validFunctionPrivileges)
	if err != nil {
		return err
	}
	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, functionName, argumentDataTypes, privilege, withGrantOption, onFuture, onAll, roles, shares)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// DeleteFunctionGrant implements schema.DeleteFunc.
func DeleteFunctionGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	functionName := d.Get("function_name").(string)
	argumentDataTypes := expandStringList(d.Get("argument_data_types").([]interface{}))
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureFunctionGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllFunctionGrant(databaseName, schemaName)
	default:
		builder = snowflake.FunctionGrant(databaseName, schemaName, functionName, argumentDataTypes)
	}

	return deleteGenericGrant(d, meta, builder)
}

// UpdateFunctionGrant implements schema.UpdateFunc.
func UpdateFunctionGrant(d *schema.ResourceData, meta interface{}) error {
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
	schemaName := d.Get("schema_name").(string)
	functionName := d.Get("function_name").(string)
	argumentDataTypes := expandStringList(d.Get("argument_data_types").([]interface{}))
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	// create the builder
	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureFunctionGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllFunctionGrant(databaseName, schemaName)
	default:
		builder = snowflake.FunctionGrant(databaseName, schemaName, functionName, argumentDataTypes)
	}

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta, builder, privilege, rolesToRevoke, sharesToRevoke,
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
	return ReadFunctionGrant(d, meta)
}
