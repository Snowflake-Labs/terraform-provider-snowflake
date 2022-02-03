package resources

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
)

var validFileFormatPrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
)

var fileFormatGrantSchema = map[string]*schema.Schema{
	"file_format_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the file format on which to grant privileges immediately (only valid if on_future is false).",
		ForceNew:    true,
	},
	"schema_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the schema containing the current or future file formats on which to grant privileges.",
		ForceNew:    true,
	},
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future file formats on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future file format.",
		Default:      "USAGE",
		ValidateFunc: validation.ValidatePrivilege(validFileFormatPrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
		ForceNew:    true,
	},
	"on_future": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future file formats in the given schema. When this is true and no schema_name is provided apply this grant on all future file formats in the given database. The file_format_name field must be unset in order to use on_future.",
		Default:     false,
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
func FileFormatGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateFileFormatGrant,
			Read:   ReadFileFormatGrant,
			Delete: DeleteFileFormatGrant,

			Schema: fileFormatGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validFileFormatPrivileges,
	}
}

// CreateFileFormatGrant implements schema.CreateFunc
func CreateFileFormatGrant(d *schema.ResourceData, meta interface{}) error {
	var fileFormatName string
	if name, ok := d.GetOk("file_format_name"); ok {
		fileFormatName = name.(string)
	}
	dbName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	priv := d.Get("privilege").(string)
	futureFileFormats := d.Get("on_future").(bool)
	grantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if (fileFormatName == "") && !futureFileFormats {
		return errors.New("file_format_name must be set unless on_future is true.")
	}
	if (fileFormatName != "") && futureFileFormats {
		return errors.New("file_format_name must be empty if on_future is true.")
	}

	var builder snowflake.GrantBuilder
	if futureFileFormats {
		builder = snowflake.FutureFileFormatGrant(dbName, schemaName)
	} else {
		builder = snowflake.FileFormatGrant(dbName, schemaName, fileFormatName)
	}

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: dbName,
		SchemaName:   schemaName,
		ObjectName:   fileFormatName,
		Privilege:    priv,
		GrantOption:  grantOption,
		Roles:        roles,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadFileFormatGrant(d, meta)
}

// ReadFileFormatGrant implements schema.ReadFunc
func ReadFileFormatGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	fileFormatName := grantID.ObjectName
	priv := grantID.Privilege

	err = d.Set("database_name", dbName)
	if err != nil {
		return err
	}
	err = d.Set("schema_name", schemaName)
	if err != nil {
		return err
	}
	futureFileFormatsEnabled := false
	if fileFormatName == "" {
		futureFileFormatsEnabled = true
	}
	err = d.Set("file_format_name", fileFormatName)
	if err != nil {
		return err
	}
	err = d.Set("on_future", futureFileFormatsEnabled)
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
	if futureFileFormatsEnabled {
		builder = snowflake.FutureFileFormatGrant(dbName, schemaName)
	} else {
		builder = snowflake.FileFormatGrant(dbName, schemaName, fileFormatName)
	}

	return readGenericGrant(d, meta, fileFormatGrantSchema, builder, futureFileFormatsEnabled, validFileFormatPrivileges)
}

// DeleteFileFormatGrant implements schema.DeleteFunc
func DeleteFileFormatGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	dbName := grantID.ResourceName
	schemaName := grantID.SchemaName
	fileFormatName := grantID.ObjectName

	futureFileFormats := (fileFormatName == "")

	var builder snowflake.GrantBuilder
	if futureFileFormats {
		builder = snowflake.FutureFileFormatGrant(dbName, schemaName)
	} else {
		builder = snowflake.FileFormatGrant(dbName, schemaName, fileFormatName)
	}
	return deleteGenericGrant(d, meta, builder)
}
