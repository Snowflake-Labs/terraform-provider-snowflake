package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/pkg/errors"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
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
		Optional:    true,
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
		ValidateFunc: validation.StringInSlice(validFileFormatPrivileges.toList(), true),
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
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future file formats in the given schema. When this is true and no schema_name is provided apply this grant on all future file formats in the given database. The file_format_name field must be unset in order to use on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"file_format_name"},
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
	var (
		fileFormatName string
		schemaName     string
	)
	if _, ok := data.GetOk("file_format_name"); ok {
		fileFormatName = data.Get("file_format_name").(string)
	}
	if _, ok := data.GetOk("schema_name"); ok {
		schemaName = data.Get("schema_name").(string)
	}
	dbName := data.Get("database_name").(string)
	priv := data.Get("privilege").(string)
	futureFileFormats := data.Get("on_future").(bool)
	grantOption := data.Get("with_grant_option").(bool)

	if (schemaName == "") && !futureFileFormats {
		return errors.New("schema_name must be set unless on_future is true.")
	}

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
	futureFileFormatsEnabled := false
	if fileFormatName == "" {
		futureFileFormatsEnabled = true
	}
	err = data.Set("file_format_name", fileFormatName)
	if err != nil {
		return err
	}
	err = data.Set("on_future", futureFileFormatsEnabled)
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

	var builder snowflake.GrantBuilder
	if futureFileFormatsEnabled {
		builder = snowflake.FutureFileFormatGrant(dbName, schemaName)
	} else {
		builder = snowflake.FileFormatGrant(dbName, schemaName, fileFormatName)
	}

	return readGenericGrant(data, meta, fileFormatGrantSchema, builder, futureFileFormatsEnabled, validFileFormatPrivileges)
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

	futureFileFormats := (fileFormatName == "")

	var builder snowflake.GrantBuilder
	if futureFileFormats {
		builder = snowflake.FutureFileFormatGrant(dbName, schemaName)
	} else {
		builder = snowflake.FileFormatGrant(dbName, schemaName, fileFormatName)
	}
	return deleteGenericGrant(data, meta, builder)
}
