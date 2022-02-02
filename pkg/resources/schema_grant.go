package resources

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

var validSchemaPrivileges = NewPrivilegeSet(
	privilegeModify,
	privilegeMonitor,
	privilegeOwnership,
	privilegeUsage,
	privilegeCreateTable,
	privilegeCreateTag,
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
	privilegeCreateRowAccessPolicy,
	privilegeCreateTemporaryTable,
	privilegeCreateMaskingPolicy,
	privilegeAddSearchOptimization,
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
		ValidateFunc: validation.ValidatePrivilege(validSchemaPrivileges.ToList(), true),
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
		Description: "Grants privilege to these shares (only valid if on_future is unset).",
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
			Update: UpdateSchemaGrant,

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
	var schemaName string
	if _, ok := d.GetOk("schema_name"); ok {
		schemaName = d.Get("schema_name").(string)
	} else {
		schemaName = ""
	}
	db := d.Get("database_name").(string)
	priv := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	grantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if (schemaName == "") && !onFuture {
		return errors.New("schema_name must be set unless on_future is true.")
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureSchemaGrant(db)
	} else {
		builder = snowflake.SchemaGrant(db, schemaName)
	}

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grantID := &grantID{
		ResourceName: db,
		SchemaName:   schemaName,
		Privilege:    priv,
		GrantOption:  grantOption,
		Roles:        roles,
	}
	dataIDInput, err := grantID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadSchemaGrant(d, meta)
}

// UpdateSchemaGrant implements schema.UpdateFunc
func UpdateSchemaGrant(d *schema.ResourceData, meta interface{}) error {
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

	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	onFuture := d.Get("on_future").(bool)

	// create the builder
	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureSchemaGrant(dbName)
	} else {
		builder = snowflake.SchemaGrant(dbName, schemaName)
	}

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta,
		builder,
		grantID.Privilege,
		rolesToRevoke,
		sharesToRevoke,
	); err != nil {
		return err
	}

	// then add
	if err := createGenericGrantRolesAndShares(
		meta,
		builder,
		grantID.Privilege,
		grantID.GrantOption,
		rolesToAdd,
		sharesToAdd,
	); err != nil {
		return err
	}

	// Done, refresh state
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
