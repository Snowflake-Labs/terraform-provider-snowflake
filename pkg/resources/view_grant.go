package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var validViewPrivileges = NewPrivilegeSet(
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
		ValidateFunc: validation.StringInSlice(validViewPrivileges.ToList(), true),
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
}

// ViewGrant returns a pointer to the resource representing a view grant
func ViewGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateViewGrant,
			Read:   ReadViewGrant,
			Delete: DeleteViewGrant,

			Schema: viewGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validViewPrivileges,
	}
}

// CreateViewGrant implements schema.CreateFunc
func CreateViewGrant(d *schema.ResourceData, meta interface{}) error {
	var (
		viewName   string
		schemaName string
	)
	if _, ok := d.GetOk("view_name"); ok {
		viewName = d.Get("view_name").(string)
	} else {
		viewName = ""
	}
	if _, ok := d.GetOk("schema_name"); ok {
		schemaName = d.Get("schema_name").(string)
	} else {
		schemaName = ""
	}
	dbName := d.Get("database_name").(string)
	priv := d.Get("privilege").(string)
	futureViews := d.Get("on_future").(bool)
	grantOption := d.Get("with_grant_option").(bool)

	if (schemaName == "") && !futureViews {
		return errors.New("schema_name must be set unless on_future is true.")
	}

	if (viewName == "") && !futureViews {
		return errors.New("view_name must be set unless on_future is true.")
	}
	if (viewName != "") && futureViews {
		return errors.New("view_name must be empty if on_future is true.")
	}

	var builder snowflake.GrantBuilder
	if futureViews {
		builder = snowflake.FutureViewGrant(dbName, schemaName)
	} else {
		builder = snowflake.ViewGrant(dbName, schemaName, viewName)
	}

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   viewName,
		Privilege:    priv,
		GrantOption:  grantOption,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadViewGrant(d, meta)
}

// ReadViewGrant implements schema.ReadFunc
func ReadViewGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	viewName := grantID.ObjectName
	priv := grantID.Privilege

	err = d.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = d.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	futureViewsEnabled := false
	if viewName == "" {
		futureViewsEnabled = true
	}
	err = d.Set("view_name", viewName)
	if err != nil {
		return err
	}
	err = d.Set("on_future", futureViewsEnabled)
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
	if futureViewsEnabled {
		builder = snowflake.FutureViewGrant(dbName, schemaName)
	} else {
		builder = snowflake.ViewGrant(dbName, schemaName, viewName)
	}

	return readGenericGrant(d, meta, viewGrantSchema, builder, futureViewsEnabled, validViewPrivileges)
}

// DeleteViewGrant implements schema.DeleteFunc
func DeleteViewGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	viewName := grantID.ObjectName

	futureViews := (viewName == "")

	var builder snowflake.GrantBuilder
	if futureViews {
		builder = snowflake.FutureViewGrant(dbName, schemaName)
	} else {
		builder = snowflake.ViewGrant(dbName, schemaName, viewName)
	}
	return deleteGenericGrant(d, meta, builder)
}
