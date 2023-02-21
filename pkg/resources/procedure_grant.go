package resources

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
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
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)

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
	shares := expandStringList(d.Get("shares").(*schema.Set).List())

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
		builder = snowflake.FutureProcedureGrant(databaseName, schemaName)
	} else {
		builder = snowflake.ProcedureGrant(databaseName, schemaName, procedureName, argumentDataTypes)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := NewProcedureGrantID(databaseName, schemaName, procedureName, argumentDataTypes, privilege, roles, shares, withGrantOption)
	d.SetId(grantID.String())
	return ReadProcedureGrant(d, meta)
}

// ReadProcedureGrant implements schema.ReadFunc.
func ReadProcedureGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseProcedureGrantID(d.Id())
	if err != nil {
		return err
	}
	if !grantID.IsOldID {
		if err := d.Set("shares", grantID.Shares); err != nil {
			return err
		}
	}

	if err := d.Set("roles", grantID.Roles); err != nil {
		return err
	}

	if err := d.Set("database_name", grantID.DatabaseName); err != nil {
		return err
	}

	if err := d.Set("schema_name", grantID.SchemaName); err != nil {
		return err
	}
	onFuture := false
	if grantID.ObjectName == "" {
		onFuture = true
	}

	if err := d.Set("procedure_name", grantID.ObjectName); err != nil {
		return err
	}

	if err := d.Set("argument_data_types", grantID.ArgumentDataTypes); err != nil {
		return err
	}

	if err := d.Set("on_future", onFuture); err != nil {
		return err
	}

	if err := d.Set("privilege", grantID.Privilege); err != nil {
		return err
	}

	if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureProcedureGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.ProcedureGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName, grantID.ArgumentDataTypes)
	}

	return readGenericGrant(d, meta, procedureGrantSchema, builder, onFuture, validProcedurePrivileges)
}

// DeleteProcedureGrant implements schema.DeleteFunc.
func DeleteProcedureGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseProcedureGrantID(d.Id())
	if err != nil {
		return err
	}

	procedureObjectName := grantID.ObjectName
	onFuture := (procedureObjectName == "")

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureProcedureGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.ProcedureGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName, grantID.ArgumentDataTypes)
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
	grantID, err := ParseProcedureGrantID(d.Id())
	if err != nil {
		return err
	}

	onFuture := (grantID.ObjectName == "")

	// create the builder
	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureProcedureGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.ProcedureGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName, grantID.ArgumentDataTypes)
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
	return ReadProcedureGrant(d, meta)
}

type ProcedureGrantID struct {
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

func NewProcedureGrantID(databaseName string, schemaName, objectName string, argumentDataTypes []string, privilege string, roles []string, shares []string, withGrantOption bool) *ProcedureGrantID {
	return &ProcedureGrantID{
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

func (v *ProcedureGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	shares := strings.Join(v.Shares, ",")
	argumentDataTypes := strings.Join(v.ArgumentDataTypes, ",")
	return fmt.Sprintf("%v❄️%v❄️%v❄️%v❄️%v❄️%v❄️%v❄️%v", v.DatabaseName, v.SchemaName, v.ObjectName, argumentDataTypes, v.Privilege, v.WithGrantOption, roles, shares)
}

func ParseProcedureGrantID(s string) (*ProcedureGrantID, error) {
	// is this an old ID format?
	if !strings.Contains(s, "❄️") {
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
		return &ProcedureGrantID{
			DatabaseName:      idParts[0],
			SchemaName:        idParts[1],
			ObjectName:        objectNameParts[0],
			ArgumentDataTypes: argumentDataTypes,
			Privilege:         idParts[3],
			Roles:             helpers.SplitStringToSlice(idParts[4], ","),
			Shares:            []string{},
			WithGrantOption:   idParts[5] == "true",
			IsOldID:           true,
		}, nil
	}
	idParts := strings.Split(s, "❄️")
	if len(idParts) != 8 {
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 8", len(idParts))
	}
	return &ProcedureGrantID{
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
