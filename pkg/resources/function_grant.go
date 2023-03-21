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
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					grantID, err := ParseFunctionGrantID(d.Id())
					if err != nil {
						return nil, err
					}
					if err := d.Set("argument_data_types", grantID.ArgumentDataTypes); err != nil {
						return nil, err
					}
					if err := d.Set("function_name", grantID.ObjectName); err != nil {
						return nil, err
					}
					if err := d.Set("schema_name", grantID.SchemaName); err != nil {
						return nil, err
					}
					if err := d.Set("database_name", grantID.DatabaseName); err != nil {
						return nil, err
					}
					if err := d.Set("privilege", grantID.Privilege); err != nil {
						return nil, err
					}
					if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
						return nil, err
					}
					if err := d.Set("roles", grantID.Roles); err != nil {
						return nil, err
					}
					if err := d.Set("shares", grantID.Shares); err != nil {
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
	var functionName string
	if name, ok := d.GetOk("function_name"); ok {
		functionName = name.(string)
	}

	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
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
		builder = snowflake.FutureFunctionGrant(databaseName, schemaName)
	} else {
		builder = snowflake.FunctionGrant(databaseName, schemaName, functionName, argumentDataTypes)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := NewFunctionGrantID(databaseName, schemaName, functionName, argumentDataTypes, privilege, roles, shares, withGrantOption)
	d.SetId(grantID.String())
	return ReadFunctionGrant(d, meta)
}

// ReadFunctionGrant implements schema.ReadFunc.
func ReadFunctionGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseFunctionGrantID(d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("database_name", grantID.DatabaseName); err != nil {
		return err
	}
	if err := d.Set("schema_name", grantID.SchemaName); err != nil {
		return err
	}
	if err := d.Set("privilege", grantID.Privilege); err != nil {
		return err
	}
	if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
		return err
	}
	onFuture := false
	if grantID.ObjectName == "" {
		onFuture = true
	}
	if err := d.Set("function_name", grantID.ObjectName); err != nil {
		return err
	}
	if err := d.Set("argument_data_types", grantID.ArgumentDataTypes); err != nil {
		return err
	}
	if err := d.Set("on_future", onFuture); err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureFunctionGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.FunctionGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName, grantID.ArgumentDataTypes)
	}

	return readGenericGrant(d, meta, functionGrantSchema, builder, onFuture, validFunctionPrivileges)
}

// DeleteFunctionGrant implements schema.DeleteFunc.
func DeleteFunctionGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseFunctionGrantID(d.Id())
	if err != nil {
		return err
	}

	onFuture := (grantID.ObjectName == "")

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureFunctionGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.FunctionGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName, grantID.ArgumentDataTypes)
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
	grantID, err := ParseFunctionGrantID(d.Id())
	if err != nil {
		return err
	}
	onFuture := (grantID.ObjectName == "")

	// create the builder
	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureFunctionGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.FunctionGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName, grantID.ArgumentDataTypes)
	}

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, rolesToRevoke, sharesToRevoke,
	); err != nil {
		return err
	}
	// then add
	if err := createGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, grantID.WithGrantOption, rolesToAdd, sharesToAdd,
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadFunctionGrant(d, meta)
}

type FunctionGrantID struct {
	DatabaseName      string
	SchemaName        string
	ObjectName        string
	ArgumentDataTypes []string
	Privilege         string
	Roles             []string
	Shares            []string
	WithGrantOption   bool
	IsOldID           bool
}

func NewFunctionGrantID(databaseName string, schemaName, objectName string, argumentDataTypes []string, privilege string, roles []string, shares []string, withGrantOption bool) *FunctionGrantID {
	return &FunctionGrantID{
		DatabaseName:      databaseName,
		SchemaName:        schemaName,
		ObjectName:        objectName,
		ArgumentDataTypes: argumentDataTypes,
		Privilege:         privilege,
		Roles:             roles,
		Shares:            shares,
		WithGrantOption:   withGrantOption,
		IsOldID:           false,
	}
}

func (v *FunctionGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	shares := strings.Join(v.Shares, ",")
	argumentDataTypes := strings.Join(v.ArgumentDataTypes, ",")
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v|%v", v.DatabaseName, v.SchemaName, v.ObjectName, argumentDataTypes, v.Privilege, v.WithGrantOption, roles, shares)
}

func ParseFunctionGrantID(s string) (*FunctionGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		objectIdentifier := idParts[2]
		if idx := strings.Index(objectIdentifier, ")"); idx != -1 {
			objectIdentifier = objectIdentifier[0:idx]
		}
		objectNameParts := strings.Split(objectIdentifier, "(")
		argumentDataTypes := []string{}
		if len(objectNameParts) > 1 {
			argumentDataTypes = helpers.SplitStringToSlice(objectNameParts[1], ",")
		}
		// remove the param names from the argument (if present)
		for i, argumentDataType := range argumentDataTypes {
			parts := strings.Split(argumentDataType, " ")
			if len(parts) > 1 {
				argumentDataTypes[i] = parts[1]
			}
		}
		withGrantOption := false
		roles := []string{}
		if len(idParts) == 6 {
			withGrantOption = idParts[5] == "true"
			roles = helpers.SplitStringToSlice(idParts[4], ",")
		} else {
			withGrantOption = idParts[4] == "true"
		}
		return &FunctionGrantID{
			DatabaseName:      idParts[0],
			SchemaName:        idParts[1],
			ObjectName:        objectNameParts[0],
			ArgumentDataTypes: argumentDataTypes,
			Privilege:         idParts[3],
			Roles:             roles,
			Shares:            []string{},
			WithGrantOption:   withGrantOption,
			IsOldID:           true,
		}, nil
	}
	idParts := strings.Split(s, "|")
	if len(idParts) < 8 {
		idParts = strings.Split(s, "❄️") // for that time in 0.56/0.57 when we used ❄️ as a separator
	}
	if len(idParts) != 8 {
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 8", len(idParts))
	}
	return &FunctionGrantID{
		DatabaseName:      idParts[0],
		SchemaName:        idParts[1],
		ObjectName:        idParts[2],
		ArgumentDataTypes: helpers.SplitStringToSlice(idParts[3], ","),
		Privilege:         idParts[4],
		WithGrantOption:   idParts[5] == "true",
		Roles:             helpers.SplitStringToSlice(idParts[6], ","),
		Shares:            helpers.SplitStringToSlice(idParts[7], ","),
		IsOldID:           false,
	}, nil
}
