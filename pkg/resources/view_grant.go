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

var validViewPrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeReferences,
	privilegeSelect,
)

var viewGrantSchema = map[string]*schema.Schema{
	"view_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the view on which to grant privileges immediately (only valid if on_future is unset).",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema containing the current or future views on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future views on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future view.",
		Default:      privilegeSelect.String(),
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice(validViewPrivileges.ToList(), true),
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
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future views in the given schema. When this is true and no schema_name is provided apply this grant on all future views in the given database. The view_name and shares fields must be unset in order to use on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"view_name", "shares"},
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

// ViewGrant returns a pointer to the resource representing a view grant.
func ViewGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateViewGrant,
			Read:   ReadViewGrant,
			Delete: DeleteViewGrant,
			Update: UpdateViewGrant,

			Schema: viewGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validViewPrivileges,
	}
}

// CreateViewGrant implements schema.CreateFunc.
func CreateViewGrant(d *schema.ResourceData, meta interface{}) error {
	var viewName string
	if _, ok := d.GetOk("view_name"); ok {
		viewName = d.Get("view_name").(string)
	}
	var schemaName string
	if _, ok := d.GetOk("schema_name"); ok {
		schemaName = d.Get("schema_name").(string)
	}
	databaseName := d.Get("database_name").(string)
	privilege := d.Get("privilege").(string)
	futureViews := d.Get("on_future").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	shares := expandStringList(d.Get("shares").(*schema.Set).List())
	if (schemaName == "") && !futureViews {
		return errors.New("schema_name must be set unless on_future is true")
	}
	if (viewName == "") && !futureViews {
		return errors.New("view_name must be set unless on_future is true")
	}
	if (viewName != "") && futureViews {
		return errors.New("view_name must be empty if on_future is true")
	}

	var builder snowflake.GrantBuilder
	if futureViews {
		builder = snowflake.FutureViewGrant(databaseName, schemaName)
	} else {
		builder = snowflake.ViewGrant(databaseName, schemaName, viewName)
	}

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grantID := NewViewGrantID(databaseName, schemaName, viewName, privilege, roles, shares, withGrantOption)
	d.SetId(grantID.String())

	return ReadViewGrant(d, meta)
}

// ReadViewGrant implements schema.ReadFunc.
func ReadViewGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := parseViewGrantID(d.Id())
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

	futureViewsEnabled := false
	if grantID.ObjectName == "" {
		futureViewsEnabled = true
	}
	err = d.Set("view_name", grantID.ObjectName)
	if err != nil {
		return err
	}
	err = d.Set("on_future", futureViewsEnabled)
	if err != nil {
		return err
	}
	err = d.Set("privilege", grantID.Privilege)
	if err != nil {
		return err
	}
	err = d.Set("with_grant_option", grantID.WithGrantOption)
	if err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if futureViewsEnabled {
		builder = snowflake.FutureViewGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.ViewGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}

	return readGenericGrant(d, meta, viewGrantSchema, builder, futureViewsEnabled, validViewPrivileges)
}

// DeleteViewGrant implements schema.DeleteFunc.
func DeleteViewGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := parseViewGrantID(d.Id())
	if err != nil {
		return err
	}

	futureViews := (grantID.ObjectName == "")

	var builder snowflake.GrantBuilder
	if futureViews {
		builder = snowflake.FutureViewGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.ViewGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}
	return deleteGenericGrant(d, meta, builder)
}

// UpdateViewGrant implements schema.UpdateFunc.
func UpdateViewGrant(d *schema.ResourceData, meta interface{}) error {
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
	grantID, err := parseViewGrantID(d.Id())
	if err != nil {
		return err
	}

	futureViews := (grantID.ObjectName == "")

	// create the builder
	var builder snowflake.GrantBuilder
	if futureViews {
		builder = snowflake.FutureViewGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.ViewGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}

	// first revoke
	err = deleteGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, rolesToRevoke, sharesToRevoke)
	if err != nil {
		return err
	}
	// then add
	err = createGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, grantID.WithGrantOption, rolesToAdd, sharesToAdd)
	if err != nil {
		return err
	}

	// Done, refresh state
	return ReadViewGrant(d, meta)
}

type ViewGrantID struct {
	DatabaseName    string
	SchemaName      string
	ObjectName      string
	Privilege       string
	Roles           []string
	Shares          []string
	WithGrantOption bool
	IsOldID         bool
}

func NewViewGrantID(databaseName string, schemaName, objectName, privilege string, roles []string, shares []string, withGrantOption bool) *ViewGrantID {
	return &ViewGrantID{
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

func (v *ViewGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	shares := strings.Join(v.Shares, ",")
	return fmt.Sprintf("%v❄️%v❄️%v❄️%v❄️%v❄️%v❄️%v", v.DatabaseName, v.SchemaName, v.ObjectName, v.Privilege, roles, shares, v.WithGrantOption)
}

func parseViewGrantID(s string) (*ViewGrantID, error) {
	// is this an old ID format?
	if !strings.Contains(s, "❄️") {
		idParts := strings.Split(s, "|")
		return &ViewGrantID{
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
	idParts := strings.Split(s, "❄️")
	if len(idParts) != 7 {
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 7", len(idParts))
	}
	return &ViewGrantID{
		DatabaseName:    idParts[0],
		SchemaName:      idParts[1],
		ObjectName:      idParts[2],
		Privilege:       idParts[3],
		Roles:           helpers.SplitStringToSlice(idParts[4], ","),
		Shares:          helpers.SplitStringToSlice(idParts[5], ","),
		WithGrantOption: idParts[6] == "true",
		IsOldID:         false,
	}, nil
}
