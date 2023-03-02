package resources

import (
	"context"
	"fmt"
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
					v, err := helpers.DecodeSnowflakeImportID(d.Id(), AccountGrantImporter{})
					if err != nil {
						return nil, err
					}
					importer := v.(AccountGrantImporter)
					err = d.Set("privilege", importer.Privilege)
					if err != nil {
						return nil, err
					}
					err = d.Set("roles", importer.Roles)
					if err != nil {
						return nil, err
					}
					err = d.Set("with_grant_option", importer.WithGrantOption)
					if err != nil {
						return nil, err
					}
					d.SetId(helpers.RandomSnowflakeID())
					return []*schema.ResourceData{d}, nil
				},
			},
		},
		ValidPrivs: validAccountPrivileges,
	}
}

type AccountGrantImporter struct {
	Privilege       string   `tf:"privilege"`
	Roles           []string `tf:"roles"`
	WithGrantOption bool     `tf:"with_grant_option"`
}

type AccountGrantID struct {
	Privilege       string
	Roles           []string
	WithGrantOption bool
	IsOldID         bool
}

func NewAccountGrantID(privilege string, roles []string, withGrantOption bool) *AccountGrantID {
	return &AccountGrantID{
		Privilege:       privilege,
		Roles:           roles,
		WithGrantOption: withGrantOption,
		IsOldID:         false,
	}
}

func (v *AccountGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	return fmt.Sprintf("%v❄️%v❄️%v", v.Privilege, v.WithGrantOption, roles)
}

func parseAccountGrantID(s string) (*AccountGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		return &AccountGrantID{
			Privilege:       idParts[3],
			Roles:           []string{},
			WithGrantOption: idParts[4] == "true",
			IsOldID:         true,
		}, nil
	}
	idParts := helpers.SplitStringToSlice(s, "|")
	if len(idParts) < 3 {
		idParts = helpers.SplitStringToSlice(s, "❄️") // for that time in 0.56/0.57 when we used ❄️ as a separator
	}
	if len(idParts) != 3 {
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 3", len(idParts))
	}
	return &AccountGrantID{
		Privilege:       idParts[0],
		WithGrantOption: idParts[1] == "true",
		Roles:           helpers.SplitStringToSlice(idParts[2], ","),
		IsOldID:         false,
	}, nil
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
	grantID := NewAccountGrantID(privilege, roles, withGrantOption)
	d.SetId(grantID.String())

	return ReadAccountGrant(d, meta)
}

// ReadAccountGrant implements schema.ReadFunc.
func ReadAccountGrant(d *schema.ResourceData, meta interface{}) error {
	builder := snowflake.AccountGrant()
	grantID, err := parseAccountGrantID(d.Id())
	if err != nil {
		return err
	}
	if err := d.Set("privilege", grantID.Privilege); err != nil {
		return err
	}
	if err := d.Set("roles", grantID.Roles); err != nil {
		return err
	}
	if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
		return err
	}

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
