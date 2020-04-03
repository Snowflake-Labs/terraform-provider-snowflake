package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var ValidViewPrivileges = newPrivilegeSet(
	privilegeSelect,
)

var viewGrantSchema = map[string]*schema.Schema{
	"view_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the view on which to grant privileges immediately (only valid if on_future is unset).",
		ForceNew:    true,
	},
	"schema_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema containing the current or future views on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future views on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future view.",
		Default:      "SELECT",
		ValidateFunc: validation.StringInSlice(ValidViewPrivileges.toList(), true),
		ForceNew:     true,
	},
	"roles": &schema.Schema{
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
		ForceNew:    true,
	},
	"shares": &schema.Schema{
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these shares (only valid if on_future is unset).",
		ForceNew:    true,
	},
	"on_future": &schema.Schema{
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future views in the given schema. When this is true and no schema_name is provided apply this grant on all future views in the given database. The view_name and shares fields must be unset in order to use on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"view_name", "shares"},
	},
}

// ViewGrant returns a pointer to the resource representing a view grant
func ViewGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateViewGrant,
		Read:   ReadViewGrant,
		Delete: DeleteViewGrant,

		Schema: viewGrantSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateViewGrant implements schema.CreateFunc
func CreateViewGrant(data *schema.ResourceData, meta interface{}) error {
	var (
		viewName   string
		schemaName string
	)
	if _, ok := data.GetOk("view_name"); ok {
		viewName = data.Get("view_name").(string)
	} else {
		viewName = ""
	}
	if _, ok := data.GetOk("schema_name"); ok {
		schemaName = data.Get("schema_name").(string)
	} else {
		schemaName = ""
	}
	dbName := data.Get("database_name").(string)
	priv := data.Get("privilege").(string)
	futureViews := data.Get("on_future").(bool)

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

	err := createGenericGrant(data, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   viewName,
		Privilege:    priv,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadViewGrant(data, meta)
}

// ReadViewGrant implements schema.ReadFunc
func ReadViewGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	viewName := grantID.ObjectName
	priv := grantID.Privilege

	err = data.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = data.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	futureViewsEnabled := false
	if viewName == "" {
		futureViewsEnabled = true
	}
	err = data.Set("view_name", viewName)
	if err != nil {
		return err
	}
	err = data.Set("on_future", futureViewsEnabled)
	if err != nil {
		return err
	}
	err = data.Set("privilege", priv)
	if err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if futureViewsEnabled {
		builder = snowflake.FutureViewGrant(dbName, schemaName)
	} else {
		builder = snowflake.ViewGrant(dbName, schemaName, viewName)
	}

	return readGenericGrant(data, meta, builder, futureViewsEnabled, ValidViewPrivileges)
}

// DeleteViewGrant implements schema.DeleteFunc
func DeleteViewGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
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
	return deleteGenericGrant(data, meta, builder)
}
