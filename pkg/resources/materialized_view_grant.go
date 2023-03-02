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

/*
NewPriviligeSet creates a set of privileges that are allowed
They are used for validation in the schema object below.
*/

var validMaterializedViewPrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeReferences,
	privilegeSelect,
)

// The schema holds the resource variables that can be provided in the Terraform.
var materializedViewGrantSchema = map[string]*schema.Schema{
	"materialized_view_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the materialized view on which to grant privileges immediately (only valid if on_future is false).",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema containing the current or future materialized views on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future materialized views on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future materialized view view.",
		Default:      "SELECT",
		ValidateFunc: validation.StringInSlice(validMaterializedViewPrivileges.ToList(), true),
		ForceNew:     true,
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
		Description: "Grants privilege to these shares (only valid if on_future is false).",
	},
	"on_future": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future materialized views in the given schema. When this is true and no schema_name is provided apply this grant on all future materialized views in the given database. The materialized_view_name and shares fields must be unset in order to use on_future.",
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
	"enable_multiple_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.",
		Default:     false,
		ForceNew:    true,
	},
}

// ViewGrant returns a pointer to the resource representing a view grant.
func MaterializedViewGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateMaterializedViewGrant,
			Read:   ReadMaterializedViewGrant,
			Delete: DeleteMaterializedViewGrant,
			Update: UpdateMaterializedViewGrant,

			Schema: materializedViewGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validMaterializedViewPrivileges,
	}
}

// CreateViewGrant implements schema.CreateFunc.
func CreateMaterializedViewGrant(d *schema.ResourceData, meta interface{}) error {
	var materializedViewName string
	if name, ok := d.GetOk("materialized_view_name"); ok {
		materializedViewName = name.(string)
	}
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	futureMaterializedViews := d.Get("on_future").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())

	if (schemaName == "") && !futureMaterializedViews {
		return errors.New("schema_name must be set unless on_future is true")
	}

	if (materializedViewName == "") && !futureMaterializedViews {
		return errors.New("materialized_view_name must be set unless on_future is true")
	}
	if (materializedViewName != "") && futureMaterializedViews {
		return errors.New("materialized_view_name must be empty if on_future is true")
	}

	var builder snowflake.GrantBuilder
	if futureMaterializedViews {
		builder = snowflake.FutureMaterializedViewGrant(databaseName, schemaName)
	} else {
		builder = snowflake.MaterializedViewGrant(databaseName, schemaName, materializedViewName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := NewMaterializedViewGrantID(databaseName, schemaName, materializedViewName, privilege, roles, shares, withGrantOption)
	d.SetId(grantID.String())

	return ReadMaterializedViewGrant(d, meta)
}

// ReadViewGrant implements schema.ReadFunc.
func ReadMaterializedViewGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := parseMaterializedViewGrantID(d.Id())
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
	futureMaterializedViewsEnabled := false
	if grantID.ObjectName == "" {
		futureMaterializedViewsEnabled = true
	}
	if err := d.Set("materialized_view_name", grantID.ObjectName); err != nil {
		return err
	}
	if err := d.Set("on_future", futureMaterializedViewsEnabled); err != nil {
		return err
	}
	if err := d.Set("privilege", grantID.Privilege); err != nil {
		return err
	}
	if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if futureMaterializedViewsEnabled {
		builder = snowflake.FutureMaterializedViewGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.MaterializedViewGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}

	return readGenericGrant(d, meta, materializedViewGrantSchema, builder, futureMaterializedViewsEnabled, validMaterializedViewPrivileges)
}

// DeleteViewGrant implements schema.DeleteFunc.
func DeleteMaterializedViewGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := parseMaterializedViewGrantID(d.Id())
	if err != nil {
		return err
	}

	futureMaterializedViews := (grantID.ObjectName == "")

	var builder snowflake.GrantBuilder
	if futureMaterializedViews {
		builder = snowflake.FutureMaterializedViewGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.MaterializedViewGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}
	return deleteGenericGrant(d, meta, builder)
}

// UpdateMaterializedViewGrant implements schema.UpdateFunc.
func UpdateMaterializedViewGrant(d *schema.ResourceData, meta interface{}) error {
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
	grantID, err := parseMaterializedViewGrantID(d.Id())
	if err != nil {
		return err
	}

	futureMaterializedViews := (grantID.ObjectName == "")

	// create the builder
	var builder snowflake.GrantBuilder
	if futureMaterializedViews {
		builder = snowflake.FutureMaterializedViewGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.MaterializedViewGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
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
	return ReadMaterializedViewGrant(d, meta)
}

type MaterializedViewGrantID struct {
	DatabaseName    string
	SchemaName      string
	ObjectName      string
	Privilege       string
	Roles           []string
	Shares          []string
	WithGrantOption bool
	IsOldID         bool
}

func NewMaterializedViewGrantID(databaseName string, schemaName, objectName, privilege string, roles []string, shares []string, withGrantOption bool) *MaterializedViewGrantID {
	return &MaterializedViewGrantID{
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

func (v *MaterializedViewGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	shares := strings.Join(v.Shares, ",")
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", v.DatabaseName, v.SchemaName, v.ObjectName, v.Privilege, v.WithGrantOption, roles, shares)
}

func parseMaterializedViewGrantID(s string) (*MaterializedViewGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		return &MaterializedViewGrantID{
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
	if len(idParts) != 7 {
		idParts := strings.Split(s, "❄️") // for that time in 0.56/0.57 when we used ❄️ as a separator
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 7", len(idParts))
	}
	return &MaterializedViewGrantID{
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
