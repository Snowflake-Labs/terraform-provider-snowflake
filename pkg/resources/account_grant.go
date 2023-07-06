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

var validAccountPrivileges = NewPrivilegeSet(
	privilegeApplyMaskingPolicy,
	privilegeApplyPasswordPolicy,
	privilegeApplyRowAccessPolicy,
	privilegeApplySessionPolicy,
	privilegeApplyTag,
	privilegeAttachPolicy,
	privilegeAudit,
	privilegeCreateAccount,
	privilegeCreateCredential,
	privilegeCreateDatabase,
	privilegeCreateDataExchangeListing,
	privilegeCreateFailoverGroup,
	privilegeCreateIntegration,
	privilegeCreateNetworkPolicy,
	privilegeCreateRole,
	privilegeCreateShare,
	privilegeCreateUser,
	privilegeCreateWarehouse,
	privilegeExecuteTask,
	privilegeImportShare,
	privilegeManageGrants,
	privilegeMonitor,
	privilegeMonitorUsage,
	privilegeMonitorExecution,
	privilegeMonitorSecurity,
	privilegeOverrideShareRestrictions,
	privilegeExecuteManagedTask,
	privilegeOrganizationSupportCases,
	privilegeProvisionApplication,
	privilegePurchaseDataExchangeListing,
	privilegeAccountSupportCases,
	privilegeUserSupportCases,
	privilegeAllPrivileges,
)

var accountGrantSchema = map[string]*schema.Schema{
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The account privilege to grant. Valid privileges are those in [globalPrivileges](https://docs.snowflake.com/en/sql-reference/sql/grant-privilege.html). To grant all privileges, use the value `ALL PRIVILEGES`.",
		Default:      privilegeMonitorUsage,
		ValidateFunc: validation.StringInSlice(validAccountPrivileges.ToList(), true),
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
	},
}

// AccountGrant returns a pointer to the resource representing an account grant.
func AccountGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateAccountGrant,
			Read:   ReadAccountGrant,
			Delete: DeleteAccountGrant,
			Update: UpdateAccountGrant,

			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_grant_privileges_to_role instead.", Schema: accountGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 3 {
						return nil, errors.New("id should be in the format 'privilege|with_grant_option|roles'")
					}
					if err := d.Set("privilege", parts[0]); err != nil {
						return nil, err
					}
					if err := d.Set("with_grant_option", helpers.StringToBool(parts[1])); err != nil {
						return nil, err
					}
					if err := d.Set("roles", helpers.StringListToList(parts[2])); err != nil {
						return nil, err
					}
					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validAccountPrivileges,
	}
}

// CreateAccountGrant implements schema.CreateFunc.
func CreateAccountGrant(d *schema.ResourceData, meta interface{}) error {
	builder := snowflake.AccountGrant()

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	privilege := d.Get("privilege").(string)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	withGrantOption := d.Get("with_grant_option").(bool)
	grantID := helpers.EncodeSnowflakeID(privilege, withGrantOption, roles)
	d.SetId(grantID)

	return ReadAccountGrant(d, meta)
}

// ReadAccountGrant implements schema.ReadFunc.
func ReadAccountGrant(d *schema.ResourceData, meta interface{}) error {
	privilege := d.Get("privilege").(string)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	withGrantOption := d.Get("with_grant_option").(bool)

	builder := snowflake.AccountGrant()
	err := readGenericGrant(d, meta, accountGrantSchema, builder, false, false, validAccountPrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(privilege, withGrantOption, roles)
	// if the ID is not in the new format, rewrite it
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// DeleteAccountGrant implements schema.DeleteFunc.
func DeleteAccountGrant(d *schema.ResourceData, meta interface{}) error {
	builder := snowflake.AccountGrant()
	return deleteGenericGrant(d, meta, builder)
}

// UpdateAccountGrant implements schema.UpdateFunc.
func UpdateAccountGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update is roles.
	// if nothing changed, nothing to update and we're done.
	if !d.HasChanges("roles") {
		return nil
	}

	rolesToAdd, rolesToRevoke := changeDiff(d, "roles")

	builder := snowflake.AccountGrant()
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)

	// first revoke
	if err := deleteGenericGrantRolesAndShares(meta, builder, privilege, "", rolesToRevoke, nil); err != nil {
		return err
	}

	// then add
	if err := createGenericGrantRolesAndShares(meta, builder, privilege, withGrantOption, rolesToAdd, nil); err != nil {
		return err
	}

	// done, refresh state
	return ReadAccountGrant(d, meta)
}
