package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

// Intentionally exclude the "ALL" alias because it is not a real privilege and
// might not interact well with this provider.
var validSchemaPrivileges = newPrivilegeSet(
	privilegeAll,
	privilegeModify,
	privilegeMonitor,
	privilegeOwnership,
	privilegeUsage,
	privilegeCreateTable,
	privilegeCreateView,
	privilegeCreateFileFormat,
	privilegeCreateStage,
	privilegeCreatePipe,
	privilegeCreateStream,
	privilegeCreateTask,
	privilegeCreateSequence,
	privilegeCreateFunction,
	privilegeCreateProcedure,
)

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
		Description:  "The privilege to grant on the schema.  Note that if \"OWNERSHIP\" is specified, ensure that the role that terraform is using is granted access.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validSchemaPrivileges.toList(), true),
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

	grantID := &grantID{
		ResourceName: db,
		SchemaName:   schema,
		Privilege:    priv,
	}
	dataIDInput, err := grantID.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadSchemaGrant(data, meta)
}

// ReadSchemaGrant implements schema.ReadFunc
func ReadSchemaGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	err = data.Set("database_name", grantID.ResourceName)
	if err != nil {
		return err
	}
	err = data.Set("schema_name", grantID.SchemaName)
	if err != nil {
		return err
	}
	err = data.Set("privilege", grantID.Privilege)
	if err != nil {
		return err
	}

	builder := snowflake.SchemaGrant(grantID.ResourceName, grantID.SchemaName)

	return readGenericGrant(data, meta, builder, false, validSchemaPrivileges)
}

// DeleteSchemaGrant implements schema.DeleteFunc
func DeleteSchemaGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}

	builder := snowflake.SchemaGrant(grantID.ResourceName, grantID.SchemaName)

	return deleteGenericGrant(data, meta, builder)
}
