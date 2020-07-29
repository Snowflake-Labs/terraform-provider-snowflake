package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var validPipePrivileges = newPrivilegeSet(
	privilegeAll,
	privilegeOwnership,
)

var pipeGrantSchema = map[string]*schema.Schema{
	"pipe_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the pipe on which to grant privileges.",
		ForceNew:    true,
	},
	"schema_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the schema containing the current pipe on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current pipe on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the pipe.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validPipePrivileges.toList(), true),
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

// PipeGrant returns a pointer to the resource representing a pipe grant
func PipeGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreatePipeGrant,
		Read:   ReadPipeGrant,
		Delete: DeletePipeGrant,

		Schema: pipeGrantSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreatePipeGrant implements schema.CreateFunc
func CreatePipeGrant(data *schema.ResourceData, meta interface{}) error {
	var pipeName string
	if _, ok := data.GetOk("pipe_name"); ok {
		pipeName = data.Get("pipe_name").(string)
	} else {
		pipeName = ""
	}
	schemaName := data.Get("schema_name").(string)
	dbName := data.Get("database_name").(string)
	priv := data.Get("privilege").(string)

	var builder snowflake.GrantBuilder
	builder = snowflake.PipeGrant(dbName, schemaName, pipeName)

	err := createGenericGrant(data, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   pipeName,
		Privilege:    priv,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadPipeGrant(data, meta)
}

// ReadPipeGrant implements schema.ReadFunc
func ReadPipeGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	pipeName := grantID.ObjectName
	priv := grantID.Privilege

	err = data.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = data.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	err = data.Set("pipe_name", pipeName)
	if err != nil {
		return err
	}
	err = data.Set("privilege", priv)
	if err != nil {
		return err
	}

	builder := snowflake.PipeGrant(dbName, schemaName, pipeName)

	return readGenericGrant(data, meta, builder, false, validPipePrivileges)
}

// DeletePipeGrant implements schema.DeleteFunc
func DeletePipeGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	pipeName := grantID.ObjectName

	builder := snowflake.PipeGrant(dbName, schemaName, pipeName)

	return deleteGenericGrant(data, meta, builder)
}
