package resources

import (
	"context"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validTagPrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeApply,
	privilegeAllPrivileges,
)

var tagGrantSchema = map[string]*schema.Schema{
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the tag on which to grant privileges.",
		ForceNew:    true,
	},
	"tag_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the tag on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the tag. To grant all privileges, use the value `ALL PRIVILEGES`.",
		Default:      "APPLY",
		ValidateFunc: validation.StringInSlice(validTagPrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
	},
	"schema_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the schema containing the tag on which to grant privileges.",
		ForceNew:    true,
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

// TagGrant returns a pointer to the resource representing a tag grant.
func TagGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create:             CreateTagGrant,
			Read:               ReadTagGrant,
			Update:             UpdateTagGrant,
			Delete:             DeleteTagGrant,
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_grant_privileges_to_role instead.",
			Schema:             tagGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 6 {
						return nil, fmt.Errorf("unexpected format of ID (%q), expected database|schema|tag|privilege|with_grant_option|roles", d.Id())
					}
					if err := d.Set("database_name", parts[0]); err != nil {
						return nil, err
					}
					if err := d.Set("schema_name", parts[1]); err != nil {
						return nil, err
					}
					if err := d.Set("tag_name", parts[2]); err != nil {
						return nil, err
					}
					if err := d.Set("privilege", parts[3]); err != nil {
						return nil, err
					}
					if err := d.Set("with_grant_option", helpers.StringToBool(parts[4])); err != nil {
						return nil, err
					}
					if err := d.Set("roles", helpers.StringListToList(parts[5])); err != nil {
						return nil, err
					}
					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validTagPrivileges,
	}
}

// CreateTagGrant implements schema.CreateFunc.
func CreateTagGrant(d *schema.ResourceData, meta interface{}) error {
	tagName := d.Get("tag_name").(string)
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.TagGrant(databaseName, schemaName, tagName)

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}
	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, tagName, privilege, withGrantOption, roles)
	d.SetId(grantID)

	return ReadTagGrant(d, meta)
}

// ReadTagGrant implements schema.ReadFunc.
func ReadTagGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	tagName := d.Get("tag_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.TagGrant(databaseName, schemaName, tagName)

	err := readGenericGrant(d, meta, tagGrantSchema, builder, false, false, validTagPrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, tagName, privilege, withGrantOption, roles)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// UpdateTagGrant implements schema.UpdateFunc.
func UpdateTagGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update is roles. if nothing changed,
	// nothing to update and we're done.
	if !d.HasChanges("roles") {
		return nil
	}

	rolesToAdd, rolesToRevoke := changeDiff(d, "roles")

	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	tagName := d.Get("tag_name").(string)
	privilege := d.Get("privilege").(string)
	reversionRole := d.Get("revert_ownership_to_role_name").(string)
	withGrantOption := d.Get("with_grant_option").(bool)

	// create the builder
	builder := snowflake.TagGrant(databaseName, schemaName, tagName)

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta,
		builder,
		privilege,
		reversionRole,
		rolesToRevoke,
		nil,
	); err != nil {
		return err
	}

	// then add
	if err := createGenericGrantRolesAndShares(
		meta,
		builder,
		privilege,
		withGrantOption,
		rolesToAdd,
		nil,
	); err != nil {
		return err
	}

	return ReadTagGrant(d, meta)
}

// DeleteTagGrant implements schema.DeleteFunc.
func DeleteTagGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	tagName := d.Get("tag_name").(string)

	builder := snowflake.TagGrant(databaseName, schemaName, tagName)

	return deleteGenericGrant(d, meta, builder)
}
