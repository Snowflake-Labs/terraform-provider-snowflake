package resources

import (
	"context"

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
)

var accountGrantSchema = map[string]*schema.Schema{
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The account privilege to grant. Valid privileges are those in [globalPrivileges](https://docs.snowflake.com/en/sql-reference/sql/grant-privilege.html)",
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

			Schema: accountGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
					v, err := helpers.DecodeSnowflakeImportID(d.Id(), AccountGrantID{})
					if err != nil {
						return nil, err
					}
					id := v.(AccountGrantID)
					d.Set("privilege", id.Privilege)
					d.Set("roles", id.Roles)
					d.Set("with_grant_option", id.WithGrantOption)
					d.SetId(helpers.RandomSnowflakeID())
					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validAccountPrivileges,
	}
}

type AccountGrantID struct {
	Privilege       string   `tf:"privilege"`
	Roles           []string `tf:"roles"`
	WithGrantOption bool     `tf:"with_grant_option"`
}

// CreateAccountGrant implements schema.CreateFunc.
func CreateAccountGrant(d *schema.ResourceData, meta interface{}) error {
	builder := snowflake.AccountGrant()

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	d.SetId(helpers.RandomSnowflakeID())

	return ReadAccountGrant(d, meta)
}

// ReadAccountGrant implements schema.ReadFunc.
func ReadAccountGrant(d *schema.ResourceData, meta interface{}) error {
	builder := snowflake.AccountGrant()

	return readGenericGrant(d, meta, accountGrantSchema, builder, false, validAccountPrivileges)
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
	if err := deleteGenericGrantRolesAndShares(meta, builder, privilege, rolesToRevoke, nil); err != nil {
		return err
	}

	// then add
	if err := createGenericGrantRolesAndShares(meta, builder, privilege, withGrantOption, rolesToAdd, nil); err != nil {
		return err
	}

	// done, refresh state
	return ReadAccountGrant(d, meta)
}
