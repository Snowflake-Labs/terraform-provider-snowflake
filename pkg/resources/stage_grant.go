package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var ValidStagePrivileges = newPrivilegeSet(
	privilegeAll,
	privilegeOwnership,
	privilegeUsage,
	// These privileges are only valid for internal stages
	privilegeRead,
	privilegeWrite,
)

var stageGrantSchema = map[string]*schema.Schema{
	"stage_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the stage on which to grant privileges.",
		ForceNew:    true,
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
		ValidateFunc: validation.StringInSlice(ValidStagePrivileges.toList(), true),
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
		Description: "Grants privilege to these shares.",
		ForceNew:    true,
	},
}

// StageGrant returns a pointer to the resource representing a stage grant
func StageGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateStageGrant,
		Read:   ReadStageGrant,
		Delete: DeleteStageGrant,

		Schema: stageGrantSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateStageGrant implements schema.CreateFunc
func CreateStageGrant(data *schema.ResourceData, meta interface{}) error {
	var stageName string
	if _, ok := data.GetOk("stage_name"); ok {
		stageName = data.Get("stage_name").(string)
	} else {
		stageName = ""
	}
	schemaName := data.Get("schema_name").(string)
	dbName := data.Get("database_name").(string)
	priv := data.Get("privilege").(string)

	var builder snowflake.GrantBuilder
	builder = snowflake.StageGrant(dbName, schemaName, stageName)

	err := createGenericGrant(data, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   stageName,
		Privilege:    priv,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadStageGrant(data, meta)
}

// ReadStageGrant implements schema.ReadFunc
func ReadStageGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	stageName := grantID.ObjectName
	priv := grantID.Privilege

	err = data.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = data.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	err = data.Set("stage_name", stageName)
	if err != nil {
		return err
	}
	err = data.Set("privilege", priv)
	if err != nil {
		return err
	}

	builder := snowflake.StageGrant(dbName, schemaName, stageName)

	return readGenericGrant(data, meta, builder, false, ValidStagePrivileges)
}

// DeleteStageGrant implements schema.DeleteFunc
func DeleteStageGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	stageName := grantID.ObjectName

	builder := snowflake.StageGrant(dbName, schemaName, stageName)

	return deleteGenericGrant(data, meta, builder)
}
