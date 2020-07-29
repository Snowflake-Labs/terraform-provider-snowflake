package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/pkg/errors"

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
	privilegeCreateExternalTable,
	privilegeCreateMaterializedView,
	privilegeCreateTemporaryTable,
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
		ValidateFunc: validation.StringInSlice(validSchemaPrivileges.toList(), true),
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
}

// SchemaGrant returns a pointer to the resource representing a view grant
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
	var schema string
	if _, ok := data.GetOk("schema_name"); ok {
		schema = data.Get("schema_name").(string)
	} else {
		schema = ""
	}
	db := data.Get("database_name").(string)
	priv := data.Get("privilege").(string)
	onFuture := data.Get("on_future").(bool)

	if (schema == "") && !onFuture {
		return errors.New("schema_name must be set unless on_future is true.")
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureSchemaGrant(db)
	} else {
		builder = snowflake.SchemaGrant(db, schema)
	}

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

	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	err = data.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = data.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	onFuture := false
	if schemaName == "" {
		onFuture = true
	}
	err = data.Set("on_future", onFuture)
	if err != nil {
		return err
	}
	err = data.Set("privilege", grantID.Privilege)
	if err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureSchemaGrant(dbName)
	} else {
		builder = snowflake.SchemaGrant(dbName, schemaName)
	}
	return readGenericGrant(data, meta, builder, onFuture, validSchemaPrivileges)
}

// DeleteSchemaGrant implements schema.DeleteFunc
func DeleteSchemaGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
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
	return deleteGenericGrant(data, meta, builder)
}
