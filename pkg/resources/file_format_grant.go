package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var ValidFileFormatPrivileges = newPrivilegeSet(
	privilegeAll,
	privilegeOwnership,
	privilegeUsage,
)

var fileFormatGrantSchema = map[string]*schema.Schema{
	"file_format_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the file format on which to grant privileges.",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the schema containing the current file format on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current file format on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the file format.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(ValidFileFormatPrivileges.toList(), true),
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
	"with_grant_option": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, allows the recipient role to grant the privileges to other roles.",
		Default:     false,
		ForceNew:    true,
	},
}

// FileFormatGrant returns a pointer to the resource representing a file format grant
func FileFormatGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateFileFormatGrant,
		Read:   ReadFileFormatGrant,
		Delete: DeleteFileFormatGrant,

		Schema: fileFormatGrantSchema,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateFileFormatGrant implements schema.CreateFunc
func CreateFileFormatGrant(data *schema.ResourceData, meta interface{}) error {
	var fileFormatName string
	if _, ok := data.GetOk("file_format_name"); ok {
		fileFormatName = data.Get("file_format_name").(string)
	} else {
		fileFormatName = ""
	}
	schemaName := data.Get("schema_name").(string)
	dbName := data.Get("database_name").(string)
	priv := data.Get("privilege").(string)
	grantOption := data.Get("with_grant_option").(bool)

	var builder snowflake.GrantBuilder
	builder = snowflake.FileFormatGrant(dbName, schemaName, fileFormatName)

	err := createGenericGrant(data, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   fileFormatName,
		Privilege:    priv,
		GrantOption:  grantOption,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadFileFormatGrant(data, meta)
}

// ReadFileFormatGrant implements schema.ReadFunc
func ReadFileFormatGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	fileFormatName := grantID.ObjectName
	priv := grantID.Privilege

	err = data.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = data.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	err = data.Set("file_format_name", fileFormatName)
	if err != nil {
		return err
	}
	err = data.Set("privilege", priv)
	if err != nil {
		return err
	}
	err = data.Set("with_grant_option", grantID.GrantOption)
	if err != nil {
		return err
	}

	builder := snowflake.FileFormatGrant(dbName, schemaName, fileFormatName)

	return readGenericGrant(data, meta, builder, false, ValidFileFormatPrivileges)
}

// DeleteFileFormatGrant implements schema.DeleteFunc
func DeleteFileFormatGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	fileFormatName := grantID.ObjectName

	builder := snowflake.FileFormatGrant(dbName, schemaName, fileFormatName)

	return deleteGenericGrant(data, meta, builder)
}

