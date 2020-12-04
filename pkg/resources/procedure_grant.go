package resources

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var validProcedurePrivileges = newPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
)

var procedureGrantSchema = map[string]*schema.Schema{
	"procedure_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the procedure on which to grant privileges immediately (only valid if on_future is unset).",
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
		Description: "List of the arguments for the procedure (must be present if procedure_name is present)",
		ForceNew:    true,
	},
	"return_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The return type of the procedure (must be present if procedure_name is present)",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema containing the current or future procedures on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future procedures on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future procedure.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validProcedurePrivileges.toList(), true),
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
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future procedures in the given schema. When this is true and no schema_name is provided apply this grant on all future procedures in the given database. The procedure_name and shares fields must be unset in order to use on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"procedure_name", "arguments", "return_type", "shares"},
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
		Default:     false,
		ForceNew:    true,
	},
}

// ProcedureGrant returns a pointer to the resource representing a procedure grant
func ProcedureGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateProcedureGrant,
		Read:   ReadProcedureGrant,
		Delete: DeleteProcedureGrant,

		Schema: procedureGrantSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateProcedureGrant implements schema.CreateFunc
func CreateProcedureGrant(data *schema.ResourceData, meta interface{}) error {
	var (
		procedureName      string
		schemaName         string
		arguments          []interface{}
		returnType         string
		procedureSignature string
		argumentTypes      []string
	)
	if _, ok := data.GetOk("procedure_name"); ok {
		procedureName = data.Get("procedure_name").(string)
		if _, ok = data.GetOk("arguments"); ok {
			arguments = data.Get("arguments").([]interface{})
		} else {
			return errors.New("arguments must be set when specifying procedure_name.")
		}
		if _, ok = data.GetOk("return_type"); ok {
			returnType = strings.ToUpper(data.Get("return_type").(string))
		} else {
			return errors.New("return_type must be set when specifying procedure_name.")
		}
	}
	if _, ok := data.GetOk("schema_name"); ok {
		schemaName = data.Get("schema_name").(string)
	}
	dbName := data.Get("database_name").(string)
	priv := data.Get("privilege").(string)
	futureProcedures := data.Get("on_future").(bool)
	grantOption := data.Get("with_grant_option").(bool)

	if (schemaName == "") && !futureProcedures {
		return errors.New("schema_name must be set unless on_future is true.")
	}

	if (procedureName == "") && !futureProcedures {
		return errors.New("procedure_name must be set unless on_future is true.")
	}
	if (procedureName != "") && futureProcedures {
		return errors.New("procedure_name must be empty if on_future is true.")
	}

	if procedureName != "" {
		procedureSignature, _, argumentTypes = formatCallableObjectName(procedureName, returnType, arguments)
	} else {
		argumentTypes = make([]string, 0)
	}

	var builder snowflake.GrantBuilder
	if futureProcedures {
		builder = snowflake.FutureProcedureGrant(dbName, schemaName)
	} else {
		builder = snowflake.ProcedureGrant(dbName, schemaName, procedureName, argumentTypes)
	}

	err := createGenericGrant(data, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   procedureSignature,
		Privilege:    priv,
		GrantOption:  grantOption,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadProcedureGrant(data, meta)
}

// ReadProcedureGrant implements schema.ReadFunc
func ReadProcedureGrant(data *schema.ResourceData, meta interface{}) error {
	var (
		procedureName string
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
	procedureSignature := grantID.ObjectName
	priv := grantID.Privilege

	err = data.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = data.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	futureProceduresEnabled := false
	if procedureSignature == "" {
		futureProceduresEnabled = true
	} else {
		procedureSignatureMap, err := parseCallableObjectName(procedureSignature)
		if err != nil {
			return err
		}
		procedureName = procedureSignatureMap["callableName"].(string)
		returnType = procedureSignatureMap["returnType"].(string)
		arguments = procedureSignatureMap["arguments"].([]interface{})
		argumentTypes = procedureSignatureMap["argumentTypes"].([]string)
	}
	err = data.Set("procedure_name", procedureName)
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
	err = data.Set("on_future", futureProceduresEnabled)
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
	if futureProceduresEnabled {
		builder = snowflake.FutureProcedureGrant(dbName, schemaName)
	} else {
		builder = snowflake.ProcedureGrant(dbName, schemaName, procedureName, argumentTypes)
	}

	return readGenericGrant(data, meta, procedureGrantSchema, builder, futureProceduresEnabled, validProcedurePrivileges)
}

// DeleteProcedureGrant implements schema.DeleteFunc
func DeleteProcedureGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName

	futureProcedures := (grantID.ObjectName == "")

	var builder snowflake.GrantBuilder
	if futureProcedures {
		builder = snowflake.FutureProcedureGrant(dbName, schemaName)
	} else {
		procedureSignatureMap, err := parseCallableObjectName(grantID.ObjectName)
		if err != nil {
			return err
		}
		procedureName := procedureSignatureMap["callableName"].(string)
		argumentTypes := procedureSignatureMap["argumentTypes"].([]string)
		builder = snowflake.ProcedureGrant(dbName, schemaName, procedureName, argumentTypes)
	}
	return deleteGenericGrant(data, meta, builder)
}
