package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validFileFormatPrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeUsage,
	privilegeAllPrivileges,
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
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all future file formats in the given schema. When this is true and no schema_name is provided apply this grant on all future file formats in the given database. The file_format_name field must be unset in order to use on_future. Cannot be used together with on_all.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"file_format_name"},
	},
	"on_all": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When this is set to true and a schema_name is provided, apply this grant on all file formats in the given schema. When this is true and no schema_name is provided apply this grant on all file formats in the given database. The file_format_name field must be unset in order to use on_all. Cannot be used together with on_future.",
		Default:       false,
		ForceNew:      true,
		ConflictsWith: []string{"file_format_name"},
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the current or future file format. To grant all privileges, use the value `ALL PRIVILEGES`",
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
	"revert_ownership_to_role_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the role to revert ownership to on destroy. Has no effect unless `privilege` is set to `OWNERSHIP`",
		Default:     "",
		ValidateFunc: func(val interface{}, key string) ([]string, []error) {
			additionalCharsToIgnoreValidation := []string{".", " ", ":", "(", ")"}
			return sdk.ValidateIdentifier(val, additionalCharsToIgnoreValidation)
		},
	},
}

// FileFormatGrant returns a pointer to the resource representing a file format grant.
func FileFormatGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create:             CreateFileFormatGrant,
			Read:               ReadFileFormatGrant,
			Delete:             DeleteFileFormatGrant,
			Update:             UpdateFileFormatGrant,
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_grant_privileges_to_account_role instead.",
			Schema:             fileFormatGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 8 {
						return nil, fmt.Errorf("unexpected format of ID (%q), expected database_name|schema_name|file_format_name|privilege|with_grant_option|on_future|roles", d.Id())
					}
					if err := d.Set("database_name", parts[0]); err != nil {
						return nil, err
					}
					if err := d.Set("schema_name", parts[1]); err != nil {
						return nil, err
					}
					if parts[2] != "" {
						if err := d.Set("file_format_name", parts[2]); err != nil {
							return nil, err
						}
					}
					if err := d.Set("privilege", parts[3]); err != nil {
						return nil, err
					}
					if err := d.Set("with_grant_option", helpers.StringToBool(parts[4])); err != nil {
						return nil, err
					}
					if err := d.Set("on_future", helpers.StringToBool(parts[5])); err != nil {
						return nil, err
					}
					if err := d.Set("on_all", helpers.StringToBool(parts[6])); err != nil {
						return nil, err
					}
					if err := d.Set("roles", helpers.StringListToList(parts[7])); err != nil {
						return nil, err
					}
					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validFileFormatPrivileges,
	}
}

// CreateFileFormatGrant implements schema.CreateFunc.
func CreateFileFormatGrant(d *schema.ResourceData, meta interface{}) error {
	fileFormatName := d.Get("file_format_name").(string)
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if (fileFormatName == "") && !onFuture && !onAll {
		return errors.New("file_format_name must be set unless on_future or on_all is true")
	}
	if (schemaName == "") && !onFuture && !onAll {
		return errors.New("schema_name must be set unless on_future or on_all is true")
	}

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureFileFormatGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllFileFormatGrant(databaseName, schemaName)
	default:
		builder = snowflake.FileFormatGrant(databaseName, schemaName, fileFormatName)
	}

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, fileFormatName, privilege, withGrantOption, onFuture, onAll, roles)
	d.SetId(grantID)

	return ReadFileFormatGrant(d, meta)
}

// ReadFileFormatGrant implements schema.ReadFunc.
func ReadFileFormatGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	fileFormatName := d.Get("file_format_name").(string)
	privilege := d.Get("privilege").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureFileFormatGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllFileFormatGrant(databaseName, schemaName)
	default:
		builder = snowflake.FileFormatGrant(databaseName, schemaName, fileFormatName)
	}

	err := readGenericGrant(d, meta, fileFormatGrantSchema, builder, onFuture, onAll, validFileFormatPrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, fileFormatName, privilege, withGrantOption, onFuture, onAll, roles)
	if d.Id() != grantID {
		d.SetId(grantID)
	}
	return nil
}

// DeleteFileFormatGrant implements schema.DeleteFunc.
func DeleteFileFormatGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	fileFormatName := d.Get("file_format_name").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)

	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureFileFormatGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllFileFormatGrant(databaseName, schemaName)
	default:
		builder = snowflake.FileFormatGrant(databaseName, schemaName, fileFormatName)
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
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	fileFormatName := d.Get("file_format_name").(string)
	privilege := d.Get("privilege").(string)
	reversionRole := d.Get("revert_ownership_to_role_name").(string)
	onFuture := d.Get("on_future").(bool)
	onAll := d.Get("on_all").(bool)
	withGrantOption := d.Get("with_grant_option").(bool)

	// create the builder
	var builder snowflake.GrantBuilder
	switch {
	case onFuture:
		builder = snowflake.FutureFileFormatGrant(databaseName, schemaName)
	case onAll:
		builder = snowflake.AllFileFormatGrant(databaseName, schemaName)
	default:
		builder = snowflake.FileFormatGrant(databaseName, schemaName, fileFormatName)
	}

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta, builder, privilege, reversionRole, rolesToRevoke, []string{},
	); err != nil {
		return err
	}
	// then add
	if err := createGenericGrantRolesAndShares(
		meta, builder, privilege, withGrantOption, rolesToAdd, []string{},
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadFileFormatGrant(d, meta)
}
