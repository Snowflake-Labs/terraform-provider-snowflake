package resources

import (
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

var validFunctionPrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
)

var functionGrantSchema = map[string]*schema.Schema{
	"function_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the function on which to grant privileges immediately (only valid if on_future is false).",
		ForceNew:    true,
	},
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
	},
	"return_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The return type of the function (must be present if function_name is present)",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the schema containing the current or future functions on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future functions on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future function.",
		Default:      "USAGE",
		ValidateFunc: validation.ValidatePrivilege(validFunctionPrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
		ForceNew:    true,
	},
	"shares": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these shares (only valid if on_future is false).",
		ForceNew:    true,
	},
	"on_future": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future functions in the given schema. When this is true and no schema_name is provided apply this grant on all future functions in the given database. The function_name, arguments, return_type, and shares fields must be unset in order to use on_future.",
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
}

// FunctionGrant returns a pointer to the resource representing a function grant
func FunctionGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateFunctionGrant,
			Read:   ReadFunctionGrant,
			Delete: DeleteFunctionGrant,

			Schema: functionGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validFunctionPrivileges,
	}
}

// CreateFunctionGrant implements schema.CreateFunc
func CreateFunctionGrant(d *schema.ResourceData, meta interface{}) error {
	var (
		functionName      string
		arguments         []interface{}
		returnType        string
		functionSignature string
		argumentTypes     []string
	)
	if name, ok := d.GetOk("function_name"); ok {
		functionName = name.(string)
		if ret, ok := d.GetOk("return_type"); ok {
			returnType = strings.ToUpper(ret.(string))
		} else {
			return errors.New("return_type must be set when specifying function_name.")
		}
	}
	dbName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	priv := d.Get("privilege").(string)
	futureFunctions := d.Get("on_future").(bool)
	grantOption := d.Get("with_grant_option").(bool)
	arguments = d.Get("arguments").([]interface{})
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if (functionName == "") && !futureFunctions {
		return errors.New("function_name must be set unless on_future is true.")
	}
	if (functionName != "") && futureFunctions {
		return errors.New("function_name must be empty if on_future is true.")
	}

	if functionName != "" {
		functionSignature, _, argumentTypes = formatCallableObjectName(functionName, returnType, arguments)
	} else {
		argumentTypes = make([]string, 0)
	}

	var builder snowflake.GrantBuilder
	if futureFunctions {
		builder = snowflake.FutureFunctionGrant(dbName, schemaName)
	} else {
		builder = snowflake.FunctionGrant(dbName, schemaName, functionName, argumentTypes)
	}

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   functionSignature,
		Privilege:    priv,
		GrantOption:  grantOption,
		Roles:        roles,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadFunctionGrant(d, meta)
}

// ReadFunctionGrant implements schema.ReadFunc
func ReadFunctionGrant(d *schema.ResourceData, meta interface{}) error {
	var (
		functionName  string
		returnType    string
		arguments     []interface{}
		argumentTypes []string
	)
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	functionSignature := grantID.ObjectName
	priv := grantID.Privilege

	err = d.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = d.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	futureFunctionsEnabled := false
	if functionSignature == "" {
		futureFunctionsEnabled = true
	} else {
		functionSignatureMap, err := parseCallableObjectName(functionSignature)
		if err != nil {
			return err
		}
		functionName = functionSignatureMap["callableName"].(string)
		returnType = functionSignatureMap["returnType"].(string)
		arguments = functionSignatureMap["arguments"].([]interface{})
		argumentTypes = functionSignatureMap["argumentTypes"].([]string)
	}
	err = d.Set("function_name", functionName)
	if err != nil {
		return err
	}
	err = d.Set("arguments", arguments)
	if err != nil {
		return err
	}
	err = d.Set("return_type", returnType)
	if err != nil {
		return err
	}
	err = d.Set("on_future", futureFunctionsEnabled)
	if err != nil {
		return err
	}
	err = d.Set("privilege", priv)
	if err != nil {
		return err
	}
	err = d.Set("with_grant_option", grantID.GrantOption)
	if err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if futureFunctionsEnabled {
		builder = snowflake.FutureFunctionGrant(dbName, schemaName)
	} else {
		builder = snowflake.FunctionGrant(dbName, schemaName, functionName, argumentTypes)
	}

	return readGenericGrant(d, meta, functionGrantSchema, builder, futureFunctionsEnabled, validFunctionPrivileges)
}

// DeleteFunctionGrant implements schema.DeleteFunc
func DeleteFunctionGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName

	futureFunctions := (grantID.ObjectName == "")

	var builder snowflake.GrantBuilder
	if futureFunctions {
		builder = snowflake.FutureFunctionGrant(dbName, schemaName)
	} else {
		functionSignatureMap, err := parseCallableObjectName(grantID.ObjectName)
		if err != nil {
			return err
		}
		functionName := functionSignatureMap["callableName"].(string)
		argumentTypes := functionSignatureMap["argumentTypes"].([]string)
		builder = snowflake.FunctionGrant(dbName, schemaName, functionName, argumentTypes)
	}
	return deleteGenericGrant(d, meta, builder)
}
