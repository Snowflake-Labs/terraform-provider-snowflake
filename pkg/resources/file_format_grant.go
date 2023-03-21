package resources

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validFileFormatPrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
)

var fileFormatGrantSchema = map[string]*schema.Schema{
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the current or future file formats on which to grant privileges.",
		ForceNew:    true,
	},
	"enable_multiple_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.",
		Default:     false,
		ForceNew:    true,
	},
	"file_format_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the file format on which to grant privileges immediately (only valid if on_future is false).",
		ForceNew:    true,
	},
	"on_future": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true and a schema_name is provided, apply this grant on all future file formats in the given schema. When this is true and no schema_name is provided apply this grant on all future file formats in the given database. The file_format_name field must be unset in order to use on_future.",
		Default:     false,
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future file format.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validFileFormatPrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Required:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Grants privilege to these roles.",
	},
	"schema_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the schema containing the current or future file formats on which to grant privileges.",
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

// FileFormatGrant returns a pointer to the resource representing a file format grant.
func FileFormatGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateFileFormatGrant,
			Read:   ReadFileFormatGrant,
			Delete: DeleteFileFormatGrant,
			Update: UpdateFileFormatGrant,

			Schema: fileFormatGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validFileFormatPrivileges,
	}
}

// CreateFileFormatGrant implements schema.CreateFunc.
func CreateFileFormatGrant(d *schema.ResourceData, meta interface{}) error {
	var fileFormatName string
	if name, ok := d.GetOk("file_format_name"); ok {
		fileFormatName = name.(string)
	}
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if (fileFormatName == "") && !onFuture {
		return errors.New("file_format_name must be set unless on_future is true")
	}
	if (fileFormatName != "") && onFuture {
		return errors.New("file_format_name must be empty if on_future is true")
	}
	if (schemaName == "") && !onFuture {
		return errors.New("schema_name must be set unless on_future is true")
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureFileFormatGrant(databaseName, schemaName)
	} else {
		builder = snowflake.FileFormatGrant(databaseName, schemaName, fileFormatName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := NewFileFormatGrantID(databaseName, schemaName, fileFormatName, privilege, roles, withGrantOption)
	d.SetId(grantID.String())

	return ReadFileFormatGrant(d, meta)
}

// ReadFileFormatGrant implements schema.ReadFunc.
func ReadFileFormatGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseFileFormatGrantID(d.Id())
	if err != nil {
		return err
	}
	if err := d.Set("roles", grantID.Roles); err != nil {
		return err
	}
	if err := d.Set("database_name", grantID.DatabaseName); err != nil {
		return err
	}
	if err := d.Set("schema_name", grantID.SchemaName); err != nil {
		return err
	}
	if err := d.Set("file_format_name", grantID.ObjectName); err != nil {
		return err
	}
	onFuture := false
	if grantID.ObjectName == "" {
		onFuture = true
	}
	if err := d.Set("on_future", onFuture); err != nil {
		return err
	}
	if err := d.Set("privilege", grantID.Privilege); err != nil {
		return err
	}
	if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
		return err
	}

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureFileFormatGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.FileFormatGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}

	return readGenericGrant(d, meta, fileFormatGrantSchema, builder, onFuture, validFileFormatPrivileges)
}

// DeleteFileFormatGrant implements schema.DeleteFunc.
func DeleteFileFormatGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseFileFormatGrantID(d.Id())
	if err != nil {
		return err
	}

	onFuture := (grantID.ObjectName == "")

	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureFileFormatGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.FileFormatGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}
	return deleteGenericGrant(d, meta, builder)
}

// UpdateFileFormatGrant implements schema.UpdateFunc.
func UpdateFileFormatGrant(d *schema.ResourceData, meta interface{}) error {
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

	grantID, err := ParseFileFormatGrantID(d.Id())
	if err != nil {
		return err
	}

	onFuture := (grantID.ObjectName == "")

	// create the builder
	var builder snowflake.GrantBuilder
	if onFuture {
		builder = snowflake.FutureFileFormatGrant(grantID.DatabaseName, grantID.SchemaName)
	} else {
		builder = snowflake.FileFormatGrant(grantID.DatabaseName, grantID.SchemaName, grantID.ObjectName)
	}

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, rolesToRevoke, []string{},
	); err != nil {
		return err
	}
	// then add
	if err := createGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, grantID.WithGrantOption, rolesToAdd, []string{},
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadFileFormatGrant(d, meta)
}

type FileFormatGrantID struct {
	DatabaseName    string
	SchemaName      string
	ObjectName      string
	Privilege       string
	Roles           []string
	WithGrantOption bool
}

func NewFileFormatGrantID(databaseName string, schemaName, objectName, privilege string, roles []string, withGrantOption bool) *FileFormatGrantID {
	return &FileFormatGrantID{
		DatabaseName:    databaseName,
		SchemaName:      schemaName,
		ObjectName:      objectName,
		Privilege:       privilege,
		Roles:           roles,
		WithGrantOption: withGrantOption,
	}
}

func (v *FileFormatGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v", v.DatabaseName, v.SchemaName, v.ObjectName, v.Privilege, v.WithGrantOption, roles)
}

func ParseFileFormatGrantID(s string) (*FileFormatGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		withGrantOption := false
		roles := []string{}
		if len(idParts) == 6 {
			withGrantOption = idParts[5] == "true"
			roles = helpers.SplitStringToSlice(idParts[4], ",")
		} else {
			withGrantOption = idParts[4] == "true"
		}
		return &FileFormatGrantID{
			DatabaseName:    idParts[0],
			SchemaName:      idParts[1],
			ObjectName:      idParts[2],
			Privilege:       idParts[3],
			Roles:           roles,
			WithGrantOption: withGrantOption,
		}, nil
	}
	idParts := strings.Split(s, "|")
	if len(idParts) < 6 {
		idParts = strings.Split(s, "❄️") // for that time in 0.56/0.57 when we used ❄️ as a separator
	}
	if len(idParts) != 6 {
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 6", len(idParts))
	}
	return &FileFormatGrantID{
		DatabaseName:    idParts[0],
		SchemaName:      idParts[1],
		ObjectName:      idParts[2],
		Privilege:       idParts[3],
		WithGrantOption: idParts[4] == "true",
		Roles:           helpers.SplitStringToSlice(idParts[5], ","),
	}, nil
}
