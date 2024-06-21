package resources

import (
	"context"
	"errors"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validRowAccessPoilcyPrivileges = NewPrivilegeSet(
	privilegeOwnership,
	privilegeApply,
	privilegeAllPrivileges,
)

var rowAccessPolicyGrantSchema = map[string]*schema.Schema{
	"database_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the database containing the row access policy on which to grant privileges.",
		ForceNew:    true,
	},
	"row_access_policy_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the row access policy on which to grant privileges immediately.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the row access policy. To grant all privileges, use the value `ALL PRIVILEGES`",
		Default:      "APPLY",
		ValidateFunc: validation.StringInSlice(validRowAccessPoilcyPrivileges.ToList(), true),
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
		Description: "The name of the schema containing the row access policy on which to grant privileges.",
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
	},
}

// RowAccessPolicyGrant returns a pointer to the resource representing a row access policy grant.
func RowAccessPolicyGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create:             CreateRowAccessPolicyGrant,
			Read:               ReadRowAccessPolicyGrant,
			Delete:             DeleteRowAccessPolicyGrant,
			Update:             UpdateRowAccessPolicyGrant,
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_grant_privileges_to_account_role instead.",
			Schema:             rowAccessPolicyGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 6 {
						return nil, errors.New("invalid row access policy grant ID format. ID must be in the format database_name|schema_name|row_access_policy_name|privilege|with_grant_option|roles")
					}
					if err := d.Set("database_name", parts[0]); err != nil {
						return nil, err
					}
					if err := d.Set("schema_name", parts[1]); err != nil {
						return nil, err
					}
					if err := d.Set("row_access_policy_name", parts[2]); err != nil {
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
		ValidPrivs: validRowAccessPoilcyPrivileges,
	}
}

// CreateRowAccessPolicyGrant implements schema.CreateFunc.
func CreateRowAccessPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	var rowAccessPolicyName string
	if name, ok := d.GetOk("row_access_policy_name"); ok {
		rowAccessPolicyName = name.(string)
	}
	if err := d.Set("row_access_policy_name", rowAccessPolicyName); err != nil {
		return err
	}
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.RowAccessPolicyGrant(databaseName, schemaName, rowAccessPolicyName)

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, rowAccessPolicyName, privilege, withGrantOption, roles)
	d.SetId(grantID)

	return ReadRowAccessPolicyGrant(d, meta)
}

// ReadRowAccessPolicyGrant implements schema.ReadFunc.
func ReadRowAccessPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	rowAccessPolicyName := d.Get("row_access_policy_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.RowAccessPolicyGrant(databaseName, schemaName, rowAccessPolicyName)

	err := readGenericGrant(d, meta, rowAccessPolicyGrantSchema, builder, false, false, validRowAccessPoilcyPrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(databaseName, schemaName, rowAccessPolicyName, privilege, withGrantOption, roles)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// DeleteRowAccessPolicyGrant implements schema.DeleteFunc.
func DeleteRowAccessPolicyGrant(d *schema.ResourceData, meta interface{}) error {
	databaseName := d.Get("database_name").(string)
	schemaName := d.Get("schema_name").(string)
	rowAccessPolicyName := d.Get("row_access_policy_name").(string)

	builder := snowflake.RowAccessPolicyGrant(databaseName, schemaName, rowAccessPolicyName)

	return deleteGenericGrant(d, meta, builder)
}

// UpdateRowAccessPolicyGrant implements schema.UpdateFunc.
func UpdateRowAccessPolicyGrant(d *schema.ResourceData, meta interface{}) error {
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
	rowAccessPolicyName := d.Get("row_access_policy_name").(string)
	privilege := d.Get("privilege").(string)
	reversionRole := d.Get("revert_ownership_to_role_name").(string)
	withGrantOption := d.Get("with_grant_option").(bool)

	// create the builder
	builder := snowflake.RowAccessPolicyGrant(databaseName, schemaName, rowAccessPolicyName)

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
	return ReadRowAccessPolicyGrant(d, meta)
}
