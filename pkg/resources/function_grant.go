package resources

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validFunctionPrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
)

var functionGrantSchema = map[string]*schema.Schema{
	"arguments": {
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The argument name",
				},
				"type": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The argument type",
				},
			},
		},
		Optional:    true,
		Description: "List of the arguments for the function (must be present if function has arguments and function_name is present)",
		ForceNew:    true,
		Deprecated:  "Use argument_data_types instead",
	},
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
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future functions in the given schema. When this is true and no schema_name is provided apply this grant on all future functions in the given database. The function_name, arguments, return_type, and shares fields must be unset in order to use on_future.",
		Default:     false,
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future function. Must be one of `USAGE` or `OWNERSHIP`.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validFunctionPrivileges.ToList(), true),
		ForceNew:     true,
	},
	"return_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The return type of the function (must be present if function_name is present)",
		ForceNew:    true,
		Deprecated:  "Not used anymore",
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
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validFunctionPrivileges,
	}
}

// CreateFunctionGrant implements schema.CreateFunc.
func CreateFunctionGrant(d *schema.ResourceData, meta interface{}) error {
	var functionName string
	if name, ok := d.GetOk("function_name"); ok {
		functionName = name.(string)
	}

	dbName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	priv := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	grantOption := d.Get("with_grant_option").(bool)
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

	if (functionName == "") && !onFuture {
		return errors.New("function_name must be set unless on_future is true")
	}
	if (functionName != "") && onFuture {
		return errors.New("function_name must be empty if on_future is true")
	}
	if (schemaName == "") && !onFuture {
		return errors.New("schema_name must be set unless on_future is true")
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureFunctionGrant(dbName, schemaName)
	} else {
		builder = snowflake.FunctionGrant(dbName, schemaName, functionName, argumentDataTypes)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	// If this is a on_futures grant then the function name and arguments do not get set. This is only used for refresh purposes.
	var functionObjectName string
	if !onFuture {
		functionObjectName = fmt.Sprintf("%s(%s)", functionName, strings.Join(argumentDataTypes, ","))
	}
	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   functionObjectName,
		Privilege:    priv,
		GrantOption:  grantOption,
		Roles:        roles,
	}
	grantID, err := grant.String()
	if err != nil {
		return err
	}
	d.SetId(grantID)
	return ReadFunctionGrant(d, meta)
}

// ReadFunctionGrant implements schema.ReadFunc.
func ReadFunctionGrant(d *schema.ResourceData, meta interface{}) error {
	var functionName string
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	if err := d.Set("database_name", dbName); err != nil {
		return err
	}
	schemaName := grantID.SchemaName
	if err := d.Set("schema_name", schemaName); err != nil {
		return err
	}
	if err := d.Set("privilege", grantID.Privilege); err != nil {
		return err
	}
	if err := d.Set("with_grant_option", grantID.GrantOption); err != nil {
		return err
	}
	functionObjectName := grantID.ObjectName
	var argumentDataTypes []string
	onFuture := false
	if functionObjectName == "" {
		onFuture = true
	} else {
		functionName, argumentDataTypes = parseFunctionObjectName(functionObjectName)
	}
	if err := d.Set("function_name", functionName); err != nil {
		return err
	}
	if err := d.Set("argument_data_types", argumentDataTypes); err != nil {
		return err
	}
	if err := d.Set("on_future", onFuture); err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureFunctionGrant(dbName, schemaName)
	} else {
		builder = snowflake.FunctionGrant(dbName, schemaName, functionName, argumentDataTypes)
	}

	return readGenericGrant(d, meta, functionGrantSchema, builder, onFuture, validFunctionPrivileges)
}

// DeleteFunctionGrant implements schema.DeleteFunc.
func DeleteFunctionGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName

	onFuture := (grantID.ObjectName == "")

	functionObjectName := grantID.ObjectName
	var functionName string
	var argumentDataTypes []string
	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureFunctionGrant(dbName, schemaName)
	} else {
		functionName, argumentDataTypes = parseFunctionObjectName(functionObjectName)
		builder = snowflake.FunctionGrant(dbName, schemaName, functionName, argumentDataTypes)
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
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	functionObjectName := grantID.ObjectName
	onFuture := (functionObjectName == "")

	// create the builder
	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureFunctionGrant(dbName, schemaName)
	} else {
		functionName, argumentDataTypes := parseFunctionObjectName(functionObjectName)
		builder = snowflake.FunctionGrant(dbName, schemaName, functionName, argumentDataTypes)
	}

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, rolesToRevoke, sharesToRevoke,
	); err != nil {
		return err
	}
	// then add
	if err := createGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, grantID.GrantOption, rolesToAdd, sharesToAdd,
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadFunctionGrant(d, meta)
}
