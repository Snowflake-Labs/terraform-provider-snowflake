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

var validTablePrivileges = NewPrivilegeSet(
	privilegeSelect,
	privilegeInsert,
	privilegeUpdate,
	privilegeDelete,
	privilegeTruncate,
	privilegeReferences,
	privilegeRebuild,
	privilegeOwnership,
)

var tableGrantSchema = map[string]*schema.Schema{
	"table_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the table on which to grant privileges immediately (only valid if on_future or on_all are unset).",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema containing the current or future tables on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future tables on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future table.",
		Default:      privilegeSelect.String(),
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice(validTablePrivileges.ToList(), true),
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
	},
	"shares": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these shares (only valid if on_future or on_all are unset).",
	},
	"on_future": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future tables in the given schema. When this is true and no schema_name is provided apply this grant on all future tables in the given database. The table_name and shares fields must be unset in order to use on_future. Cannot be used together with on_all.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"table_name", "shares"},
	},
	"on_all": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all tables in the given schema. When this is true and no schema_name is provided apply this grant on all tables in the given database. The table_name and shares fields must be unset in order to use on_all. Cannot be used together with on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"table_name", "shares"},
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
		Default:     false,
		ForceNew:    true,
	},
	"enable_multiple_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.",
		Default:     false,
	},
}

// TableGrant returns a pointer to the resource representing a Table grant.
func TableGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateTableGrant,
			Read:   ReadTableGrant,
			Delete: DeleteTableGrant,
			Update: UpdateTableGrant,

			Schema: tableGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 9 {
						return nil, errors.New("invalid ID specified for Table Grant, should be in format database_name|schema_name|table_name|privilege|with_grant_option|on_future|on_all|roles|shares")
					}
					if err := d.Set("database_name", parts[0]); err != nil {
						return nil, err
					}
					if err := d.Set("schema_name", parts[1]); err != nil {
						return nil, err
					}
					if parts[2] != "" {
						if err := d.Set("table_name", parts[2]); err != nil {
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
		ValidPrivs: validTablePrivileges,
	}
}

// CreateTableGrant implements schema.CreateFunc.
func CreateTableGrant(d *schema.ResourceData, meta interface{}) error {
	tableName := d.Get("table_name").(string)
	var schemaName string
	if _, ok := d.GetOk("schema_name"); ok {
		schemaName = d.Get("schema_name").(string)
	}
	databaseName := d.Get("database_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	if onFuture && onAll {
		return errors.New("on_future and on_all cannot both be true")
	}
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())
	if (schemaName == "") && !onFuture && !onAll {
		return errors.New("schema_name must be set unless on_future or on_all is true")
	}

	if (tableName == "") && !onFuture && !onAll {
		return errors.New("table_name must be set unless on_future or on_all is true")
	}

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureTableGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllTableGrant(databaseName, schemaName)
	default:
		builder = snowflake.TableGrant(databaseName, schemaName, tableName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, tableName, privilege, withGrantOption, onFuture, onAll, roles, shares)
	d.SetId(grantID)
	return ReadTableGrant(d, meta)
}

// ReadTableGrant implements schema.ReadFunc.
func ReadTableGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	tableName := d.Get("table_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureTableGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllTableGrant(databaseName, schemaName)
	default:
		builder = snowflake.TableGrant(databaseName, schemaName, tableName)
	}
	err := readGenericGrant(d, meta, tableGrantSchema, builder, onFuture, onAll, validTablePrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, tableName, privilege, withGrantOption, onFuture, onAll, roles, shares)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// DeleteTableGrant implements schema.DeleteFunc.
func DeleteTableGrant(d *schema.ResourceData, meta interface{}) error {
	var builder snowflake.GrantBuilder
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	tableName := d.Get("table_name").(string)
	switch {
	case onFuture:
		builder = snowflake.FutureTableGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllTableGrant(databaseName, schemaName)
	default:
		builder = snowflake.TableGrant(databaseName, schemaName, tableName)
	}
	return deleteGenericGrant(d, meta, builder)
}

// UpdateTableGrant implements schema.UpdateFunc.
func UpdateTableGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update are roles or shares
	// if nothing changed, nothing to update, and we're done
	if !d.HasChanges("roles", "shares") {
		return nil
	}

	// difference calculates roles/shares to add/revoke
	difference := func(key string) (toAdd []string, toRevoke []string) {
		o, n := d.GetChange(key)
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)
		toAdd = expandStringList(newSet.Difference(oldSet).List())
		toRevoke = expandStringList(oldSet.Difference(newSet).List())
		return
	}

	rolesToAdd := []string{}
	rolesToRevoke := []string{}
	sharesToAdd := []string{}
	sharesToRevoke := []string{}
	if d.HasChange("roles") {
		rolesToAdd, rolesToRevoke = difference("roles")
	}
	if d.HasChange("shares") {
		sharesToAdd, sharesToRevoke = difference("shares")
	}

	// create the builder
	var builder snowflake.GrantBuilder
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	tableName := d.Get("table_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)

	switch {
	case onFuture:
		builder = snowflake.FutureTableGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllTableGrant(databaseName, schemaName)
	default:
		builder = snowflake.TableGrant(databaseName, schemaName, tableName)
	}

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta,
		builder,
		privilege,
		rolesToRevoke,
		sharesToRevoke,
	); err != nil {
		return err
	}

	// then add
	if err := createGenericGrantRolesAndShares(
		meta,
		builder,
		privilege,
		withGrantOption,
		rolesToAdd,
		sharesToAdd,
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadTableGrant(d, meta)
}
