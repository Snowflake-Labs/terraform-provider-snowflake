package resources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var validViewPrivileges = []string{"SELECT"}

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
		Default:     "PUBLIC",
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
		ValidateFunc: validation.StringInSlice(validViewPrivileges, true),
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
		Description:   "When this is set to true, apply this grant on all future views in the given schema.  The view_name and shares fields must be unset in order to use on_future.",
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
	var viewName string
	if _, ok := data.GetOk("view_name"); ok {
		viewName = data.Get("view_name").(string)
	} else {
		viewName = ""
	}
	schemaName := data.Get("schema_name").(string)
	dbName := data.Get("database_name").(string)
	priv := data.Get("privilege").(string)
	futureViews := data.Get("on_future").(bool)

	if (viewName == "") && (futureViews == false) {
		return errors.New("view_name must be set unless on_future is true.")
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

	// ID format is <db_name>|<schema_name>|<view_name>|<privilege>
	// view_name is empty when on_future = true
	if futureViews {
		data.SetId(fmt.Sprintf("%v|%v||%v", dbName, schemaName, priv))
	} else {
		data.SetId(fmt.Sprintf("%v|%v|%v|%v", dbName, schemaName, viewName, priv))
	}

	return ReadViewGrant(data, meta)
}

// ReadViewGrant implements schema.ReadFunc
func ReadViewGrant(data *schema.ResourceData, meta interface{}) error {
	dbName, schemaName, viewName, priv, err := splitGrantID(data.Id())
	if err != nil {
		return err
	}
	err = data.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = data.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	futureViews := false
	if viewName == "" {
		futureViews = true
	}
	err = data.Set("view_name", viewName)
	if err != nil {
		return err
	}
	err = data.Set("on_future", futureViews)
	if err != nil {
		return err
	}
	err = data.Set("privilege", priv)
	if err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if futureViews {
		builder = snowflake.FutureViewGrant(dbName, schemaName)
		return readGenericGrant(data, meta, builder, true)
	} else {
		builder = snowflake.ViewGrant(dbName, schemaName, viewName)
		return readGenericGrant(data, meta, builder, false)
	}
}

// DeleteViewGrant implements schema.DeleteFunc
func DeleteViewGrant(data *schema.ResourceData, meta interface{}) error {
	dbName, schemaName, viewName, _, err := splitGrantID(data.Id())
	if err != nil {
		return err
	}

	futureViews := false
	if viewName == "" {
		futureViews = true
	}

	var builder snowflake.GrantBuilder
	if futureViews {
		builder = snowflake.FutureViewGrant(dbName, schemaName)
	} else {
		builder = snowflake.ViewGrant(dbName, schemaName, viewName)
	}
	return deleteGenericGrant(data, meta, builder)
}
