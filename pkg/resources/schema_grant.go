package resources

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var validSchemaPrivileges = []string{"USAGE"}

var schemaGrantSchema = map[string]*schema.Schema{
	"schema_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the schema on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the schema on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the schema.",
		Default:      "USAGE",
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
		Description: "Grants privilege to these shares.",
		ForceNew:    true,
	},
}

// ViewGrant returns a pointer to the resource representing a view grant
func SchemaGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateSchemaGrant,
		Read:   ReadSchemaGrant,
		Delete: DeleteSchemaGrant,

		Schema: schemaGrantSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateSchemaGrant implements schema.CreateFunc
func CreateSchemaGrant(data *schema.ResourceData, meta interface{}) error {
	schema := data.Get("schema_name").(string)
	db := data.Get("database_name").(string)
	priv := data.Get("privilege").(string)
	builder := snowflake.SchemaGrant(db, schema)

	err := createGenericGrant(data, meta, builder)
	if err != nil {
		return err
	}

	// ID format is <db_name>|<schema_name>||<privilege>
	data.SetId(fmt.Sprintf("%v|%v||%v", db, schema, priv))

	return ReadSchemaGrant(data, meta)
}

// ReadSchemaGrant implements schema.ReadFunc
func ReadSchemaGrant(data *schema.ResourceData, meta interface{}) error {
	db, schema, _, priv, err := splitGrantID(data.Id())
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
	err = data.Set("privilege", priv)
	if err != nil {
		return err
	}

	builder := snowflake.SchemaGrant(db, schema)

	return readGenericGrant(data, meta, builder)
}

// DeleteSchemaGrant implements schema.DeleteFunc
func DeleteSchemaGrant(data *schema.ResourceData, meta interface{}) error {
	db, schema, _, _, err := splitGrantID(data.Id())
	if err != nil {
		return err
	}

	builder := snowflake.SchemaGrant(db, schema)

	return deleteGenericGrant(data, meta, builder)
}
