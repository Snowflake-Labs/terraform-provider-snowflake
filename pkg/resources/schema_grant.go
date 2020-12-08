package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var validSchemaPrivileges = NewPrivilegeSet(
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
	privilegeCreateExternalTable,
	privilegeCreateMaterializedView,
	privilegeCreateTemporaryTable,
	privilegeCreateMaskingPolicy,
)

var schemaGrantSchema = map[string]*schema.Schema{
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the schema on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future schema. Note that if \"OWNERSHIP\" is specified, ensure that the role that terraform is using is granted access.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validSchemaPrivileges.ToList(), true),
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
		Description:   "When this is set to true, apply this grant on all future schemas in the given database. The schema_name and shares fields must be unset in order to use on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"schema_name", "shares"},
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
		Default:     false,
		ForceNew:    true,
	},
}

// SchemaGrant returns a pointer to the resource representing a view grant
func SchemaGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateSchemaGrant,
			Read:   ReadSchemaGrant,
			Delete: DeleteSchemaGrant,

			Schema: schemaGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validSchemaPrivileges,
	}
}

// CreateSchemaGrant implements schema.CreateFunc
func CreateSchemaGrant(d *schema.ResourceData, meta interface{}) error {
	var schema string
	if _, ok := d.GetOk("schema_name"); ok {
		schema = d.Get("schema_name").(string)
	} else {
		schema = ""
	}
	db := d.Get("database_name").(string)
	priv := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	grantOption := d.Get("with_grant_option").(bool)

	if (schema == "") && !onFuture {
		return errors.New("schema_name must be set unless on_future is true.")
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureSchemaGrant(db)
	} else {
		builder = snowflake.SchemaGrant(db, schema)
	}

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grantID := &grantID{
		ResourceName: db,
		SchemaName:   schema,
		Privilege:    priv,
		GrantOption:  grantOption,
	}
	dataIDInput, err := grantID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadSchemaGrant(d, meta)
}

// ReadSchemaGrant implements schema.ReadFunc
func ReadSchemaGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	err = d.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = d.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	onFuture := false
	if schemaName == "" {
		onFuture = true
	}
	err = d.Set("on_future", onFuture)
	if err != nil {
		return err
	}
	err = d.Set("privilege", grantID.Privilege)
	if err != nil {
		return err
	}
	err = d.Set("with_grant_option", grantID.GrantOption)
	if err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureSchemaGrant(dbName)
	} else {
		builder = snowflake.SchemaGrant(dbName, schemaName)
	}
	return readGenericGrant(d, meta, schemaGrantSchema, builder, onFuture, validSchemaPrivileges)
}

// DeleteSchemaGrant implements schema.DeleteFunc
func DeleteSchemaGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	onFuture := false
	if schemaName == "" {
		onFuture = true
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureSchemaGrant(dbName)
	} else {
		builder = snowflake.SchemaGrant(dbName, schemaName)
	}
	return deleteGenericGrant(d, meta, builder)
}
