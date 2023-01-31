package resources

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validProcedurePrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
)

var procedureGrantSchema = map[string]*schema.Schema{
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
		Deprecated:  "use argument_data_types instead.",
	},
	"argument_data_types": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "List of the argument data types for the procedure (must be present if procedure has arguments and procedure_name is present)",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future procedures on which to grant privileges.",
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
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future procedures in the given schema. When this is true and no schema_name is provided apply this grant on all future procedures in the given database. The procedure_name and shares fields must be unset in order to use on_future.",
		Default:     false,
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future procedure.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validProcedurePrivileges.ToList(), true),
		ForceNew:     true,
	},
	"procedure_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the procedure on which to grant privileges immediately (only valid if on_future is false).",
		ForceNew:    true,
	},
	"return_type": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The return type of the procedure (must be present if procedure_name is present)",
		ForceNew:    true,
		Deprecated:  "return_type is no longer required. It will be removed in a future release.",
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
		Description: "The name of the schema containing the current or future procedures on which to grant privileges.",
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

// ProcedureGrant returns a pointer to the resource representing a procedure grant.
func ProcedureGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateProcedureGrant,
			Read:   ReadProcedureGrant,
			Delete: DeleteProcedureGrant,
			Update: UpdateProcedureGrant,

			Schema: procedureGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validProcedurePrivileges,
	}
}

// CreateProcedureGrant implements schema.CreateFunc.
func CreateProcedureGrant(d *schema.ResourceData, meta interface{}) error {
	var procedureName string
	if name, ok := d.GetOk("procedure_name"); ok {
		procedureName = name.(string)
	}
	dbName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	priv := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	grantOption := d.Get("with_grant_option").(bool)

	argumentDataTypes := make([]string, 0)

	if v, ok := d.GetOk("arguments"); ok {
		arguments := v.([]interface{})
		for _, argument := range arguments {
			argumentDataTypes = append(argumentDataTypes, argument.(map[string]interface{})["data_type"].(string))
		}
	}

	if v, ok := d.GetOk("argument_data_types"); ok {
		argumentDataTypes = expandStringList(v.([]interface{}))
	}

	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if (procedureName == "") && !onFuture {
		return errors.New("procedure_name must be set unless on_future is true")
	}
	if (procedureName != "") && onFuture {
		return errors.New("procedure_name must be empty if on_future is true")
	}
	if (schemaName == "") && !onFuture {
		return errors.New("schema_name must be set unless on_future is true")
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureProcedureGrant(dbName, schemaName)
	} else {
		builder = snowflake.ProcedureGrant(dbName, schemaName, procedureName, argumentDataTypes)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	// If this is a on_futures grant then the procedure name and arguments do not get set. This is only used for refresh purposes.
	var procedureObjectName string
	if !onFuture {
		procedureObjectName = fmt.Sprintf("%s(%s)", procedureName, strings.Join(argumentDataTypes, ", "))
	}
	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   procedureObjectName,
		Privilege:    priv,
		GrantOption:  grantOption,
		Roles:        roles,
	}
	grantID, err := grant.String()
	if err != nil {
		return err
	}
	d.SetId(grantID)
	return ReadProcedureGrant(d, meta)
}

// ReadProcedureGrant implements schema.ReadFunc.
func ReadProcedureGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	procedureObjectName := grantID.ObjectName
	priv := grantID.Privilege

	if err := d.Set("database_name", dbName); err != nil {
		return err
	}

	if err := d.Set("schema_name", schemaName); err != nil {
		return err
	}
	onFuture := false
	var procedureName string
	argumentDataTypes := make([]string, 0)
	if procedureObjectName == "" {
		onFuture = true
	} else {
		procedureName, argumentDataTypes = parseFunctionObjectName(procedureObjectName)
	}

	if err := d.Set("procedure_name", procedureName); err != nil {
		return err
	}

	if err := d.Set("argument_data_types", argumentDataTypes); err != nil {
		return err
	}

	if err := d.Set("on_future", onFuture); err != nil {
		return err
	}

	if err := d.Set("privilege", priv); err != nil {
		return err
	}

	if err := d.Set("with_grant_option", grantID.GrantOption); err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureProcedureGrant(dbName, schemaName)
	} else {
		builder = snowflake.ProcedureGrant(dbName, schemaName, procedureName, argumentDataTypes)
	}

	return readGenericGrant(d, meta, procedureGrantSchema, builder, onFuture, validProcedurePrivileges)
}

// DeleteProcedureGrant implements schema.DeleteFunc.
func DeleteProcedureGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName

	procedureObjectName := grantID.ObjectName
	onFuture := (procedureObjectName == "")

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureProcedureGrant(dbName, schemaName)
	} else {
		procedureName, argumentDataTypes := parseFunctionObjectName(procedureObjectName)
		builder = snowflake.ProcedureGrant(dbName, schemaName, procedureName, argumentDataTypes)
	}
	return deleteGenericGrant(d, meta, builder)
}

// UpdateProcedureGrant implements schema.UpdateFunc.
func UpdateProcedureGrant(d *schema.ResourceData, meta interface{}) error {
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
	procedureObjectName := grantID.ObjectName
	onFuture := (procedureObjectName == "")

	// create the builder
	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureProcedureGrant(dbName, schemaName)
	} else {
		procedureName, argumentDataTypes := parseFunctionObjectName(procedureObjectName)
		builder = snowflake.ProcedureGrant(dbName, schemaName, procedureName, argumentDataTypes)
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
	return ReadProcedureGrant(d, meta)
}
