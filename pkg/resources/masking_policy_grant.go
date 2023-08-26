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

var validMaskingPoilcyPrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeApply,
	privilegeAllPrivileges,
)

var maskingPolicyGrantSchema = map[string]*schema.Schema{
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the masking policy on which to grant privileges.",
		ForceNew:    true,
	},
	"masking_policy_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the masking policy on which to grant privileges immediately.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the masking policy. To grant all privileges, use the value `ALL PRIVILEGES`",
		Default:      "APPLY",
		ValidateFunc: validation.StringInSlice(validMaskingPoilcyPrivileges.ToList(), true),
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
		Description: "The name of the schema containing the masking policy on which to grant privileges.",
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

// MaskingPolicyGrant returns a pointer to the resource representing a masking policy grant.
func MaskingPolicyGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create:             CreateMaskingPolicyGrant,
			Read:               ReadMaskingPolicyGrant,
			Delete:             DeleteMaskingPolicyGrant,
			Update:             UpdateMaskingPolicyGrant,
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_grant_privileges_to_role instead.",
			Schema:             maskingPolicyGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 6 {
						return nil, fmt.Errorf("unexpected format of ID (%q), expected database_name|schema_name|masking_policy_name|privilege|with_grant_option|roles", d.Id())
					}
					if err := d.Set("database_name", parts[0]); err != nil {
						return nil, err
					}
					if err := d.Set("schema_name", parts[1]); err != nil {
						return nil, err
					}
					if parts[2] != "" {
						if err := d.Set("masking_policy_name", parts[2]); err != nil {
							return nil, err
						}
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
		ValidPrivs: validMaskingPoilcyPrivileges,
	}
}

// CreateMaskingPolicyGrant implements schema.CreateFunc.
func CreateMaskingPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	maskingPolicyName := d.Get("masking_policy_name").(string)
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.MaskingPolicyGrant(databaseName, schemaName, maskingPolicyName)
	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, maskingPolicyName, privilege, withGrantOption, roles)
	d.SetId(grantID)

	return ReadMaskingPolicyGrant(d, meta)
}

// ReadMaskingPolicyGrant implements schema.ReadFunc.
func ReadMaskingPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	maskingPolicyName := d.Get("masking_policy_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.MaskingPolicyGrant(databaseName, schemaName, maskingPolicyName)

	err := readGenericGrant(d, meta, maskingPolicyGrantSchema, builder, false, false, validMaskingPoilcyPrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, maskingPolicyName, privilege, withGrantOption, roles)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// DeleteMaskingPolicyGrant implements schema.DeleteFunc.
func DeleteMaskingPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	maskingPolicyName := d.Get("masking_policy_name").(string)

	builder := snowflake.MaskingPolicyGrant(databaseName, schemaName, maskingPolicyName)

	return deleteGenericGrant(d, meta, builder)
}

// UpdateMaskingPolicyGrant implements schema.UpdateFunc.
func UpdateMaskingPolicyGrant(d *schema.ResourceData, meta interface{}) error {
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
	maskingPolicyName := d.Get("masking_policy_name").(string)
	privilege := d.Get("privilege").(string)
	reversionRole := d.Get("revert_ownership_to_role_name").(string)
	withGrantOption := d.Get("with_grant_option").(bool)

	// create the builder
	builder := snowflake.MaskingPolicyGrant(databaseName, schemaName, maskingPolicyName)

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
	return ReadMaskingPolicyGrant(d, meta)
}
