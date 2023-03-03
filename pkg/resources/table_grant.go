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
		Description: "The name of the table on which to grant privileges immediately (only valid if on_future is unset).",
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
		Description: "Grants privilege to these shares (only valid if on_future is unset).",
	},
	"on_future": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future tables in the given schema. When this is true and no schema_name is provided apply this grant on all future tables in the given database. The table_name and shares fields must be unset in order to use on_future.",
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
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validTablePrivileges,
	}
}

// CreateTableGrant implements schema.CreateFunc.
func CreateTableGrant(d *schema.ResourceData, meta interface{}) error {
	var tableName string
	if _, ok := d.GetOk("table_name"); ok {
		tableName = d.Get("table_name").(string)
	}
	var schemaName string
	if _, ok := d.GetOk("schema_name"); ok {
		schemaName = d.Get("schema_name").(string)
	}
	databaseName := d.Get("database_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())
	if (schemaName == "") && !onFuture {
		return errors.New("schema_name must be set unless on_future is true")
	}

	if (tableName == "") && !onFuture {
		return errors.New("table_name must be set unless on_future is true")
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureTableGrant(databaseName, schemaName)
	} else {
		builder = snowflake.TableGrant(databaseName, schemaName, tableName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := NewTableGrantID(databaseName, schemaName, tableName, privilege, roles, shares, withGrantOption)
	d.SetId(grantID.String())
	return ReadTableGrant(d, meta)
}

// ReadTableGrant implements schema.ReadFunc.
func ReadTableGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseTableGrantID(d.Id())
	if err != nil {
		return err
	}

	if !grantID.IsOldID {
		if err := d.Set("roles", grantID.Roles); err != nil {
			return err
		}
		if err := d.Set("shares", grantID.Shares); err != nil {
			return err
		}
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
	if err := d.Set("table_name", grantID.ObjectName); err != nil {
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
		builder = snowflake.FutureTableGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.TableGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}

	return readGenericGrant(d, meta, tableGrantSchema, builder, onFuture, validTablePrivileges)
}

// DeleteTableGrant implements schema.DeleteFunc.
func DeleteTableGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseTableGrantID(d.Id())
	if err != nil {
		return err
	}

	onFuture := false
	if grantID.ObjectName == "" {
		onFuture = true
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureTableGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.TableGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}
	return deleteGenericGrant(d, meta, builder)
}

// UpdateTableGrant implements schema.UpdateFunc.
func UpdateTableGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update are roles or shares
	// if nothing changed, nothing to update and we're done
	if !d.HasChanges("roles", "shares") {
		return nil
	}

	// difference calculates roles/shares to add/revoke
	difference := func(key string) (toAdd []string, toRevoke []string) {
		old, new := d.GetChange(key)
		oldSet := old.(*schema.Set)
		newSet := new.(*schema.Set)
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

	grantID, err := ParseTableGrantID(d.Id())
	if err != nil {
		return err
	}

	onFuture := (grantID.ObjectName == "")

	// create the builder
	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureTableGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.TableGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta,
		builder,
		grantID.Privilege,
		rolesToRevoke,
		sharesToRevoke,
	); err != nil {
		return err
	}

	// then add
	if err := createGenericGrantRolesAndShares(
		meta,
		builder,
		grantID.Privilege,
		grantID.WithGrantOption,
		rolesToAdd,
		sharesToAdd,
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadTableGrant(d, meta)
}

type TableGrantID struct {
	DatabaseName    string
	SchemaName      string
	ObjectName      string
	Privilege       string
	Roles           []string
	Shares          []string
	WithGrantOption bool
	IsOldID         bool
}

func NewTableGrantID(databaseName string, schemaName, objectName, privilege string, roles []string, shares []string, withGrantOption bool) *TableGrantID {
	return &TableGrantID{
		DatabaseName:    databaseName,
		SchemaName:      schemaName,
		ObjectName:      objectName,
		Privilege:       privilege,
		Roles:           roles,
		Shares:          shares,
		WithGrantOption: withGrantOption,
		IsOldID:         false,
	}
}

func (v *TableGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	shares := strings.Join(v.Shares, ",")
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", v.DatabaseName, v.SchemaName, v.ObjectName, v.Privilege, v.WithGrantOption, roles, shares)
}

func ParseTableGrantID(s string) (*TableGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		return &TableGrantID{
			DatabaseName:    idParts[0],
			SchemaName:      idParts[1],
			ObjectName:      idParts[2],
			Privilege:       idParts[3],
			Roles:           helpers.SplitStringToSlice(idParts[4], ","),
			Shares:          []string{},
			WithGrantOption: idParts[5] == "true",
			IsOldID:         true,
		}, nil
	}
	idParts := strings.Split(s, "|")
	if len(idParts) < 7 {
		idParts = strings.Split(s, "❄️") // for that time in 0.56/0.57 when we used ❄️ as a separator
	}
	if len(idParts) != 7 {
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 7", len(idParts))
	}
	return &TableGrantID{
		DatabaseName:    idParts[0],
		SchemaName:      idParts[1],
		ObjectName:      idParts[2],
		Privilege:       idParts[3],
		WithGrantOption: idParts[4] == "true",
		Roles:           helpers.SplitStringToSlice(idParts[5], ","),
		Shares:          helpers.SplitStringToSlice(idParts[6], ","),
		IsOldID:         false,
	}, nil
}
