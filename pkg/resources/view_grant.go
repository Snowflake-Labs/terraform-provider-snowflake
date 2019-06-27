package resources

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var validViewPrivileges = []string{"SELECT"}

var viewGrantSchema = map[string]*schema.Schema{
	"view_name": &schema.Schema{
		Type: schema.TypeString,
		Required: true,
		Description: "The name of the view on which to grant privileges.",
		ForceNew: true,
	},
	"schema_name": &schema.Schema{
		Type: schema.TypeString,
		Optional: true,
		Default: "PUBLIC",
		Description: "The name of the schema containing the view on which to grant privileges.",
		ForceNew: true,
	},
	"database_name": &schema.Schema{
		Type: schema.TypeString,
		Required: true,
		Description: "The name of the database containing the view on which to grant privileges.",
		ForceNew: true,
	},
	"privilege": &schema.Schema{
		Type: schema.TypeString,
		Optional: true,
		Description: "The privilege to grant on the view.",
		Default: "SELECT",
		ValidateFunc: validation.StringInSlice(validViewPrivileges, true),
		ForceNew: true,
	},
	"roles": &schema.Schema{
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
		ForceNew: true,
	},
	"shares": &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Schema{Type: schema.TypeString},
		Optional: true,
		Description: "Grants privilege to these shares.",
		ForceNew: true,
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
	view := data.Get("view_name").(string)
	schema := data.Get("schema_name").(string)
	db := data.Get("database_name").(string)
	priv := data.Get("privilege").(string)
	builder := snowflake.ViewGrant(db, schema, view)

	err := createGenericGrant(data, meta, builder)
	if err != nil {
		return err
	}

	// ID format is <db_name>|<schema_name>|<view_name>|<privilege>
	data.SetId(fmt.Sprintf("%v|%v|%v|%v", db, schema, view, priv))

	return ReadViewGrant(data, meta)
}

// ReadViewGrant implements schema.ReadFunc
func ReadViewGrant(data *schema.ResourceData, meta interface{}) error {
	db, schema, view, priv, err := splitGrantID(data.Id())
	if err != nil {
		return err
	}
	err = data.Set("database_name", db)
	if err != nil {
		return err
	}
	err = data.Set("schema_name", schema)
	if err != nil {
		return err
	}
	err = data.Set("view_name", view)
	if err != nil {
		return err
	}
	err = data.Set("privilege", priv)
	if err != nil {
		return err
	}

	builder := snowflake.ViewGrant(db, schema, view)

	return readGenericGrant(data, meta, builder)
}

// DeleteViewGrant implements schema.DeleteFunc
func DeleteViewGrant(data *schema.ResourceData, meta interface{}) error {
	db, schema, view, _, err := splitGrantID(data.Id())
	if err != nil {
		return err
	}

	builder := snowflake.ViewGrant(db, schema, view)

	return deleteGenericGrant(data, meta, builder)
}
