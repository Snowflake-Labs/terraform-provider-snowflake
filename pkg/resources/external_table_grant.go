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

var validExternalTablePrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeReferences,
	privilegeSelect,
	privilegeAllPrivileges,
)

var externalTableGrantSchema = map[string]*schema.Schema{
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future external tables on which to grant privileges.",
		ForceNew:    true,
	},
	"enable_multiple_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.",
		Default:     false,
		ForceNew:    true,
	},
	"external_table_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the external table on which to grant privileges immediately (only valid if on_future is false).",
		ForceNew:    true,
	},
	"on_future": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future external tables in the given schema. When this is true and no schema_name is provided apply this grant on all future external tables in the given database. The external_table_name and shares fields must be unset in order to use on_future. Cannot be used together with on_all.",
		Default:     false,
		ForceNew:    true,
	},
	"on_all": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all external tables in the given schema. When this is true and no schema_name is provided apply this grant on all external tables in the given database. The external_table_name and shares fields must be unset in order to use on_all. Cannot be used together with on_future.",
		Default:     false,
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future external table. To grant all privileges, use the value `ALL PRIVILEGES`",
		Default:      "SELECT",
		ValidateFunc: validation.StringInSlice(validExternalTablePrivileges.ToList(), true),
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
		Description: "The name of the schema containing the current or future external tables on which to grant privileges.",
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
			return snowflake.ValidateIdentifier(val, additionalCharsToIgnoreValidation)
		},
	},
}

// ExternalTableGrant returns a pointer to the resource representing a external table grant.
func ExternalTableGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateExternalTableGrant,
			Read:   ReadExternalTableGrant,
			Delete: DeleteExternalTableGrant,
			Update: UpdateExternalTableGrant,

			Schema: externalTableGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 9 {
						return nil, errors.New("external table grant ID should be in the format database|schema|external_table|privilege|with_grant_option|on_future|roles|shares")
					}
					if err := d.Set("database_name", parts[0]); err != nil {
						return nil, err
					}
					if err := d.Set("schema_name", parts[1]); err != nil {
						return nil, err
					}
					if parts[2] != "" {
						if err := d.Set("external_table_name", parts[2]); err != nil {
							return nil, err
						}
					}
					if err := d.Set("privilege", parts[3]); err != nil {
						return nil, err
					}
					if err := d.Set("with_grant_option", helpers.StringToBool(parts[4])); err != nil {
						return nil, err
					}
					if err := d.Set("on_future", helpers.StringToBool(parts[5])); err != nil {
						return nil, err
					}
					if err := d.Set("on_all", helpers.StringToBool(parts[6])); err != nil {
						return nil, err
					}
					if err := d.Set("roles", helpers.StringListToList(parts[7])); err != nil {
						return nil, err
					}
					if err := d.Set("shares", helpers.StringListToList(parts[8])); err != nil {
						return nil, err
					}

					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validExternalTablePrivileges,
	}
}

// CreateExternalTableGrant implements schema.CreateFunc.
func CreateExternalTableGrant(d *schema.ResourceData, meta interface{}) error {
	externalTableName := d.Get("external_table_name").(string)
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	if onFuture && onAll {
		return errors.New("on_future and on_all cannot both be true")
	}
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())

	if (externalTableName == "") && !onFuture && !onAll {
		return errors.New("external_table_name must be set unless on_future or on_all is true")
	}
	if (externalTableName != "") && (onFuture || onAll) {
		return errors.New("external_table_name must be empty if on_future or on_all is true")
	}
	if (schemaName == "") && !onFuture && !onAll {
		return errors.New("schema_name must be set unless on_future or on_all is true")
	}

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureExternalTableGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllExternalTableGrant(databaseName, schemaName)
	default:
		builder = snowflake.ExternalTableGrant(databaseName, schemaName, externalTableName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, externalTableName, privilege, withGrantOption, onFuture, onAll, roles, shares)
	d.SetId(grantID)

	return ReadExternalTableGrant(d, meta)
}

// ReadExternalTableGrant implements schema.ReadFunc.
func ReadExternalTableGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	externalTableName := d.Get("external_table_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureExternalTableGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllExternalTableGrant(databaseName, schemaName)
	default:
		builder = snowflake.ExternalTableGrant(databaseName, schemaName, externalTableName)
	}

	err := readGenericGrant(d, meta, externalTableGrantSchema, builder, onFuture, onAll, validExternalTablePrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, externalTableName, privilege, withGrantOption, onFuture, onAll, roles, shares)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// DeleteExternalTableGrant implements schema.DeleteFunc.
func DeleteExternalTableGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	externalTableName := d.Get("external_table_name").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureExternalTableGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllExternalTableGrant(databaseName, schemaName)
	default:
		builder = snowflake.ExternalTableGrant(databaseName, schemaName, externalTableName)
	}
	return deleteGenericGrant(d, meta, builder)
}

// UpdateExternalTableGrant implements schema.UpdateFunc.
func UpdateExternalTableGrant(d *schema.ResourceData, meta interface{}) error {
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
	externalTableName := d.Get("external_table_name").(string)
	privilege := d.Get("privilege").(string)
	reversionRole := d.Get("revert_ownership_to_role_name").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	// create the builder
	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureExternalTableGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllExternalTableGrant(databaseName, schemaName)
	default:
		builder = snowflake.ExternalTableGrant(databaseName, schemaName, externalTableName)
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
	return ReadExternalTableGrant(d, meta)
}
