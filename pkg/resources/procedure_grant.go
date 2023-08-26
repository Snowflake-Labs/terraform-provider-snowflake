package resources

import (
	"context"
	"errors"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validProcedurePrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
	privilegeAllPrivileges,
)

var procedureGrantSchema = map[string]*schema.Schema{
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
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future procedures in the given schema. When this is true and no schema_name is provided apply this grant on all future procedures in the given database. The procedure_name and shares fields must be unset in order to use on_future. Cannot be used together with on_all.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"procedure_name"},
	},
	"on_all": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all procedures in the given schema. When this is true and no schema_name is provided apply this grant on all procedures in the given database. The procedure_name and shares fields must be unset in order to use on_all. Cannot be used together with on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"procedure_name"},
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future procedure. To grant all privileges, use the value `ALL PRIVILEGES`",
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
	"revert_ownership_to_role_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the role to revert ownership to on destroy. Has no effect unless `privilege` is set to `OWNERSHIP`",
		Default:     "",
		ValidateFunc: func(val interface{}, key string) ([]string, []error) {
			additionalCharsToIgnoreValidation := []string{".", " ", ":", "(", ")"}
			return sdk.ValidateIdentifier(val, additionalCharsToIgnoreValidation)
		},
	},
}

// ProcedureGrant returns a pointer to the resource representing a procedure grant.
func ProcedureGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create:             CreateProcedureGrant,
			Read:               ReadProcedureGrant,
			Delete:             DeleteProcedureGrant,
			Update:             UpdateProcedureGrant,
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_grant_privileges_to_role instead.",
			Schema:             procedureGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 10 {
						return nil, errors.New("incorrect ID format (expecting database_name|schema_name|procedure_name|argument_data_types|privilege|with_grant_option|on_future|roles|shares)")
					}
					if err := d.Set("database_name", parts[0]); err != nil {
						return nil, err
					}
					if err := d.Set("schema_name", parts[1]); err != nil {
						return nil, err
					}
					if parts[2] != "" {
						if err := d.Set("procedure_name", parts[2]); err != nil {
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
		ValidPrivs: validProcedurePrivileges,
	}
}

// CreateProcedureGrant implements schema.CreateFunc.
func CreateProcedureGrant(d *schema.ResourceData, meta interface{}) error {
	procedureName := d.Get("procedure_name").(string)
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
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

	if (procedureName == "") && !onFuture && !onAll {
		return errors.New("procedure_name must be set unless on_future or on_all is true")
	}
	if (schemaName == "") && !onFuture && !onAll {
		return errors.New("schema_name must be set unless on_future or on_all is true")
	}

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureProcedureGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllProcedureGrant(databaseName, schemaName)
	default:
		builder = snowflake.ProcedureGrant(databaseName, schemaName, procedureName, argumentDataTypes)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, procedureName, argumentDataTypes, privilege, withGrantOption, onFuture, onAll, roles, shares)
	d.SetId(grantID)
	return ReadProcedureGrant(d, meta)
}

// ReadProcedureGrant implements schema.ReadFunc.
func ReadProcedureGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	procedureName := d.Get("procedure_name").(string)
	argumentDataTypes := expandStringList(d.Get("argument_data_types").([]interface{}))
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureProcedureGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllProcedureGrant(databaseName, schemaName)
	default:
		builder = snowflake.ProcedureGrant(databaseName, schemaName, procedureName, argumentDataTypes)
	}

	err := readGenericGrant(d, meta, procedureGrantSchema, builder, onFuture, onAll, validProcedurePrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, procedureName, argumentDataTypes, privilege, withGrantOption, onFuture, onAll, roles, shares)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// DeleteProcedureGrant implements schema.DeleteFunc.
func DeleteProcedureGrant(d *schema.ResourceData, meta interface{}) error {
	procedureName := d.Get("procedure_name").(string)
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	argumentDataTypes := expandStringList(d.Get("argument_data_types").([]interface{}))
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureProcedureGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllProcedureGrant(databaseName, schemaName)
	default:
		builder = snowflake.ProcedureGrant(databaseName, schemaName, procedureName, argumentDataTypes)
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

	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	procedureName := d.Get("procedure_name").(string)
	argumentDataTypes := expandStringList(d.Get("argument_data_types").([]interface{}))
	privilege := d.Get("privilege").(string)
	reversionRole := d.Get("revert_ownership_to_role_name").(string)
	withGrantOption := d.Get("with_grant_option").(bool)

	// create the builder
	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureProcedureGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllProcedureGrant(databaseName, schemaName)
	default:
		builder = snowflake.ProcedureGrant(databaseName, schemaName, procedureName, argumentDataTypes)
	}

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta, builder, privilege, reversionRole, rolesToRevoke, sharesToRevoke,
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
	return ReadProcedureGrant(d, meta)
}
