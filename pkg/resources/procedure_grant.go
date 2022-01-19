package resources

import (
	"strings"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

var validProcedurePrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
)

var procedureGrantSchema = map[string]*schema.Schema{
	"procedure_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the procedure on which to grant privileges immediately (only valid if on_future is false).",
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
		Description: "List of the arguments for the procedure (must be present if procedure has arguments and procedure_name is present)",
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
		Required:    true,
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
		ValidateFunc: validation.ValidatePrivilege(validProcedurePrivileges.ToList(), true),
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
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future procedures in the given schema. When this is true and no schema_name is provided apply this grant on all future procedures in the given database. The procedure_name and shares fields must be unset in order to use on_future.",
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

// ProcedureGrant returns a pointer to the resource representing a procedure grant
func ProcedureGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateProcedureGrant,
			Read:   ReadProcedureGrant,
			Delete: DeleteProcedureGrant,

			Schema: procedureGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validProcedurePrivileges,
	}
}

// CreateProcedureGrant implements schema.CreateFunc
func CreateProcedureGrant(d *schema.ResourceData, meta interface{}) error {
	var (
		procedureName      string
		arguments          []interface{}
		returnType         string
		procedureSignature string
		argumentTypes      []string
	)
	if name, ok := d.GetOk("procedure_name"); ok {
		procedureName = name.(string)
		if ret, ok := d.GetOk("return_type"); ok {
			returnType = strings.ToUpper(ret.(string))
		} else {
			return errors.New("return_type must be set when specifying procedure_name.")
		}
	}
	dbName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	priv := d.Get("privilege").(string)
	futureProcedures := d.Get("on_future").(bool)
	grantOption := d.Get("with_grant_option").(bool)
	arguments = d.Get("arguments").([]interface{})
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

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

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   procedureSignature,
		Privilege:    priv,
		GrantOption:  grantOption,
		Roles:        roles,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadProcedureGrant(d, meta)
}

// ReadProcedureGrant implements schema.ReadFunc
func ReadProcedureGrant(d *schema.ResourceData, meta interface{}) error {
	var (
		procedureName string
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
	procedureSignature := grantID.ObjectName
	priv := grantID.Privilege

	err = d.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = d.Set("schema_name", schemaName)
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
	err = d.Set("procedure_name", procedureName)
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
	err = d.Set("on_future", futureProceduresEnabled)
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
	if futureProceduresEnabled {
		builder = snowflake.FutureProcedureGrant(dbName, schemaName)
	} else {
		builder = snowflake.ProcedureGrant(dbName, schemaName, procedureName, argumentTypes)
	}

	return readGenericGrant(d, meta, procedureGrantSchema, builder, futureProceduresEnabled, validProcedurePrivileges)
}

// DeleteProcedureGrant implements schema.DeleteFunc
func DeleteProcedureGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
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
	return deleteGenericGrant(d, meta, builder)
}
