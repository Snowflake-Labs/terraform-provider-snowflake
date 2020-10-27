package resources

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var validFunctionPrivileges = newPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
)

var functionGrantSchema = map[string]*schema.Schema{
	"function_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the function on which to grant privileges immediately (only valid if on_future is unset).",
		ForceNew:    true,
	},
	"arguments": {
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "The argument name",
				},
				"type": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "The argument type",
				},
			},
		},
		Optional:    true,
		Description: "List of the arguments for the function (must be present if function_name is present)",
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
		Optional:    true,
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
		ValidateFunc: validation.StringInSlice(validFunctionPrivileges.toList(), true),
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
		Description: "Grants privilege to these shares (only valid if on_future is unset).",
		ForceNew:    true,
	},
	"on_future": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future functions in the given schema. When this is true and no schema_name is provided apply this grant on all future functions in the given database. The function_name, arguments, return_type, and shares fields must be unset in order to use on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"function_name", "arguments", "return_type", "shares"},
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
func FunctionGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateFunctionGrant,
		Read:   ReadFunctionGrant,
		Delete: DeleteFunctionGrant,

		Schema: functionGrantSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateFunctionGrant implements schema.CreateFunc
func CreateFunctionGrant(data *schema.ResourceData, meta interface{}) error {
	var (
		functionName      string
		schemaName        string
		arguments         []interface{}
		returnType        string
		functionSignature string
		argumentTypes     []string
	)
	if _, ok := data.GetOk("function_name"); ok {
		functionName = data.Get("function_name").(string)
		if _, ok = data.GetOk("arguments"); ok {
			arguments = data.Get("arguments").([]interface{})
		} else {
			return errors.New("arguments must be set when specifying function_name.")
		}
		if _, ok = data.GetOk("return_type"); ok {
			returnType = strings.ToUpper(data.Get("return_type").(string))
		} else {
			return errors.New("return_type must be set when specifying function_name.")
		}
	} else {
		functionName = ""
	}
	if _, ok := data.GetOk("schema_name"); ok {
		schemaName = data.Get("schema_name").(string)
	} else {
		schemaName = ""
	}
	dbName := data.Get("database_name").(string)
	priv := data.Get("privilege").(string)
	futureFunctions := data.Get("on_future").(bool)
	grantOption := data.Get("with_grant_option").(bool)

	if (schemaName == "") && !futureFunctions {
		return errors.New("schema_name must be set unless on_future is true.")
	}

	if (functionName == "") && !futureFunctions {
		return errors.New("function_name must be set unless on_future is true.")
	}
	if (functionName != "") && futureFunctions {
		return errors.New("function_name must be empty if on_future is true.")
	}

	if functionName != "" {
		functionSignature, _, argumentTypes = formatCallableObjectName(functionName, returnType, arguments)
	} else {
		functionSignature = ""
		argumentTypes = make([]string, 0)
	}

	var builder snowflake.GrantBuilder
	if futureFunctions {
		builder = snowflake.FutureFunctionGrant(dbName, schemaName)
	} else {
		builder = snowflake.FunctionGrant(dbName, schemaName, functionName, argumentTypes)
	}

	err := createGenericGrant(data, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   functionSignature,
		Privilege:    priv,
		GrantOption:  grantOption,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadFunctionGrant(data, meta)
}

// ReadFunctionGrant implements schema.ReadFunc
func ReadFunctionGrant(data *schema.ResourceData, meta interface{}) error {
	var (
		functionName  string
		returnType    string
		arguments     []interface{}
		argumentTypes []string
	)
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	functionSignature := grantID.ObjectName
	priv := grantID.Privilege

	err = data.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = data.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	futureFunctionsEnabled := false
	if functionSignature == "" {
		futureFunctionsEnabled = true
		functionName = ""
		returnType = ""
		arguments = nil
		argumentTypes = nil
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
	err = data.Set("function_name", functionName)
	if err != nil {
		return err
	}
	err = data.Set("arguments", arguments)
	if err != nil {
		return err
	}
	err = data.Set("return_type", returnType)
	if err != nil {
		return err
	}
	err = data.Set("on_future", futureFunctionsEnabled)
	if err != nil {
		return err
	}
	err = data.Set("privilege", priv)
	if err != nil {
		return err
	}
	err = data.Set("with_grant_option", grantID.GrantOption)
	if err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if futureFunctionsEnabled {
		builder = snowflake.FutureFunctionGrant(dbName, schemaName)
	} else {
		builder = snowflake.FunctionGrant(dbName, schemaName, functionName, argumentTypes)
	}

	return readGenericGrant(data, meta, builder, futureFunctionsEnabled, validFunctionPrivileges)
}

// DeleteFunctionGrant implements schema.DeleteFunc
func DeleteFunctionGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
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
	return deleteGenericGrant(data, meta, builder)
}
