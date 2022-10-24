package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"
)

var validPipePrivileges = NewPrivilegeSet(
	privilegeMonitor,
	privilegeOperate,
	privilegeOwnership,
)

var pipeGrantSchema = map[string]*schema.Schema{
	"pipe_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the pipe on which to grant privileges immediately (only valid if on_future is false).",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the schema containing the current or future pipes on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future pipes on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future pipe.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validPipePrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
	},
	"on_future": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future pipes in the given schema. When this is true and no schema_name is provided apply this grant on all future pipes in the given database. The pipe_name field must be unset in order to use on_future.",
		Default:     false,
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
		ForceNew:    true,
	},
}

// PipeGrant returns a pointer to the resource representing a pipe grant.
func PipeGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreatePipeGrant,
			Read:   ReadPipeGrant,
			Delete: DeletePipeGrant,
			Update: UpdatePipeGrant,

			Schema: pipeGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validPipePrivileges,
	}
}

// CreatePipeGrant implements schema.CreateFunc.
func CreatePipeGrant(d *schema.ResourceData, meta interface{}) error {
	var pipeName string
	if name, ok := d.GetOk("pipe_name"); ok {
		pipeName = name.(string)
	}
	dbName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	priv := d.Get("privilege").(string)
	futurePipes := d.Get("on_future").(bool)
	grantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if (pipeName == "") && !futurePipes {
		return errors.New("pipe_name must be set unless on_future is true.")
	}
	if (pipeName != "") && futurePipes {
		return errors.New("pipe_name must be empty if on_future is true.")
	}

	var builder snowflake.GrantBuilder
	if futurePipes {
		builder = snowflake.FuturePipeGrant(dbName, schemaName)
	} else {
		builder = snowflake.PipeGrant(dbName, schemaName, pipeName)
	}

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   pipeName,
		Privilege:    priv,
		GrantOption:  grantOption,
		Roles:        roles,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadPipeGrant(d, meta)
}

// ReadPipeGrant implements schema.ReadFunc.
func ReadPipeGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	pipeName := grantID.ObjectName
	priv := grantID.Privilege

	err = d.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = d.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	futurePipesEnabled := false
	if pipeName == "" {
		futurePipesEnabled = true
	}
	err = d.Set("pipe_name", pipeName)
	if err != nil {
		return err
	}
	err = d.Set("on_future", futurePipesEnabled)
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
	if futurePipesEnabled {
		builder = snowflake.FuturePipeGrant(dbName, schemaName)
	} else {
		builder = snowflake.PipeGrant(dbName, schemaName, pipeName)
	}

	return readGenericGrant(d, meta, pipeGrantSchema, builder, futurePipesEnabled, validPipePrivileges)
}

// DeletePipeGrant implements schema.DeleteFunc.
func DeletePipeGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	pipeName := grantID.ObjectName

	futurePipes := (pipeName == "")

	var builder snowflake.GrantBuilder
	if futurePipes {
		builder = snowflake.FuturePipeGrant(dbName, schemaName)
	} else {
		builder = snowflake.PipeGrant(dbName, schemaName, pipeName)
	}
	return deleteGenericGrant(d, meta, builder)
}

// UpdatePipeGrant implements schema.UpdateFunc.
func UpdatePipeGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update are roles or shares
	// if nothing changed, nothing to update and we're done
	if !d.HasChanges("roles") {
		return nil
	}

	rolesToAdd := []string{}
	rolesToRevoke := []string{}

	if d.HasChange("roles") {
		rolesToAdd, rolesToRevoke = changeDiff(d, "roles")
	}

	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}

	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	pipeName := grantID.ObjectName
	futurePipes := (pipeName == "")

	// create the builder
	var builder snowflake.GrantBuilder
	if futurePipes {
		builder = snowflake.FuturePipeGrant(dbName, schemaName)
	} else {
		builder = snowflake.PipeGrant(dbName, schemaName, pipeName)
	}

	// first revoke
	err = deleteGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, rolesToRevoke, []string{})
	if err != nil {
		return err
	}
	// then add
	err = createGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, grantID.GrantOption, rolesToAdd, []string{})
	if err != nil {
		return err
	}

	// Done, refresh state
	return ReadPipeGrant(d, meta)
}
