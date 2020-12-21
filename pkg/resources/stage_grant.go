package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var validStagePrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
	// These privileges are only valid for internal stages
	privilegeRead,
	privilegeWrite,
)

var stageGrantSchema = map[string]*schema.Schema{
	"stage_name": {
		Type:          schema.TypeString,
		Optional:      true,
		Description:   "The name of the stage on which to grant privilege (only valid if on_future is false).",
		ForceNew:      true,
		ConflictsWith: []string{"on_future"},
	},
	"schema_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the schema containing the current stage on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current stage on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the stage.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validStagePrivileges.ToList(), true),
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
		Description: "Grants privilege to these shares (only valid if on_future is false).",
		ForceNew:    true,
	},
	"on_future": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future stages in the given schema. When this is true and no schema_name is provided apply this grant on all future stages in the given database. The stage_name and shares fields must be unset in order to use on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"stage_name", "shares"},
	},
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
		Default:     false,
		ForceNew:    true,
	},
}

// StageGrant returns a pointer to the resource representing a stage grant
func StageGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateStageGrant,
			Read:   ReadStageGrant,
			Delete: DeleteStageGrant,

			Schema: stageGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validStagePrivileges,
	}
}

// CreateStageGrant implements schema.CreateFunc
func CreateStageGrant(d *schema.ResourceData, meta interface{}) error {
	var stageName string
	if name, ok := d.GetOk("stage_name"); ok {
		stageName = name.(string)
	}
	dbName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	priv := d.Get("privilege").(string)
	futureStages := d.Get("on_future").(bool)
	grantOption := d.Get("with_grant_option").(bool)

	var builder snowflake.GrantBuilder
	if futureStages {
		builder = snowflake.FutureStageGrant(dbName, schemaName)
	} else {
		builder = snowflake.StageGrant(dbName, schemaName, stageName)
	}

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   stageName,
		Privilege:    priv,
		GrantOption:  grantOption,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadStageGrant(d, meta)
}

// ReadStageGrant implements schema.ReadFunc
func ReadStageGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	stageName := grantID.ObjectName
	priv := grantID.Privilege

	err = d.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = d.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	futureStagesEnabled := false
	if stageName == "" {
		futureStagesEnabled = true
	}
	err = d.Set("stage_name", stageName)
	if err != nil {
		return err
	}
	err = d.Set("on_future", futureStagesEnabled)
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
	if futureStagesEnabled {
		builder = snowflake.FutureStageGrant(dbName, schemaName)
	} else {
		builder = snowflake.StageGrant(dbName, schemaName, stageName)
	}

	return readGenericGrant(d, meta, stageGrantSchema, builder, futureStagesEnabled, validStagePrivileges)
}

// DeleteStageGrant implements schema.DeleteFunc
func DeleteStageGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	stageName := grantID.ObjectName

	futureStages := (stageName == "")

	var builder snowflake.GrantBuilder
	if futureStages {
		builder = snowflake.FutureStageGrant(dbName, schemaName)
	} else {
		builder = snowflake.StageGrant(dbName, schemaName, stageName)
	}

	return deleteGenericGrant(d, meta, builder)
}
