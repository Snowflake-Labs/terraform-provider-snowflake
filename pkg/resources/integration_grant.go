package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validIntegrationPrivileges = NewPrivilegeSet(
	privilegeUsage,
	privilegeOwnership,
	privilegeAllPrivileges,
)

var integrationGrantSchema = map[string]*schema.Schema{
	"integration_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier for the integration; must be unique for your account.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the integration. To grant all privileges, use the value `ALL PRIVILEGES`",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validIntegrationPrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
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

// IntegrationGrant returns a pointer to the resource representing a integration grant.
func IntegrationGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create:             CreateIntegrationGrant,
			Read:               ReadIntegrationGrant,
			Delete:             DeleteIntegrationGrant,
			Update:             UpdateIntegrationGrant,
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_grant_privileges_to_role instead.",
			Schema:             integrationGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 4 {
						return nil, fmt.Errorf("invalid ID %v: expected integration_name|privilege|with_grant_option|roles", d.Id())
					}
					if err := d.Set("integration_name", parts[0]); err != nil {
						return nil, err
					}
					if err := d.Set("privilege", parts[1]); err != nil {
						return nil, err
					}
					if err := d.Set("with_grant_option", helpers.StringToBool(parts[2])); err != nil {
						return nil, err
					}
					if err := d.Set("roles", helpers.StringListToList(parts[3])); err != nil {
						return nil, err
					}
					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validIntegrationPrivileges,
	}
}

// CreateIntegrationGrant implements schema.CreateFunc.
func CreateIntegrationGrant(d *schema.ResourceData, meta interface{}) error {
	integrationName := d.Get("integration_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.IntegrationGrant(integrationName)
	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}
	grantID := helpers.EncodeSnowflakeID(integrationName, privilege, withGrantOption, roles)
	d.SetId(grantID)

	return ReadIntegrationGrant(d, meta)
}

// ReadIntegrationGrant implements schema.ReadFunc.
func ReadIntegrationGrant(d *schema.ResourceData, meta interface{}) error {
	integrationName := d.Get("integration_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.IntegrationGrant(integrationName)

	err := readGenericGrant(d, meta, integrationGrantSchema, builder, false, false, validIntegrationPrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(integrationName, privilege, withGrantOption, roles)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// DeleteIntegrationGrant implements schema.DeleteFunc.
func DeleteIntegrationGrant(d *schema.ResourceData, meta interface{}) error {
	integrationName := d.Get("integration_name").(string)
	builder := snowflake.IntegrationGrant(integrationName)

	return deleteGenericGrant(d, meta, builder)
}

// UpdateIntegrationGrant implements schema.UpdateFunc.
func UpdateIntegrationGrant(d *schema.ResourceData, meta interface{}) error {
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

	integrationName := d.Get("integration_name").(string)
	privilege := d.Get("privilege").(string)
	reversionRole := d.Get("revert_ownership_to_role_name").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	// create the builder
	builder := snowflake.IntegrationGrant(integrationName)

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
	return ReadIntegrationGrant(d, meta)
}
