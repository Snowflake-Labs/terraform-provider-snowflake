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

var validFailoverGroupPrivileges = NewPrivilegeSet(
	privilegeFailover,
	privilegeMonitor,
	privilegeOwnership,
	privilegeReplicate,
	privilegeAllPrivileges,
)

var failoverGroupGrantSchema = map[string]*schema.Schema{
	"enable_multiple_grants": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "When this is set to true, multiple grants of the same type can be created. This will cause Terraform to not revoke grants applied to roles and objects outside Terraform.",
		Default:     false,
		ForceNew:    true,
	},
	"failover_group_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The name of the failover group on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the failover group. To grant all privileges, use the value `ALL PRIVILEGES`",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validFailoverGroupPrivileges.ToList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Required:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Grants privilege to these roles.",
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
			return snowflake.ValidateIdentifier(val, additionalCharsToIgnoreValidation)
		},
	},
}

// FailoverGroup returns a pointer to the resource representing a file format grant.
func FailoverGroupGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateFailoverGroupGrant,
			Read:   ReadFailoverGroupGrant,
			Delete: DeleteFailoverGroupGrant,
			Update: UpdateFailoverGroupGrant,

			Schema: failoverGroupGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 4 {
						return nil, fmt.Errorf("unexpected format of ID (%q), expected failover_group_name|privilege|with_grant_option|roles", d.Id())
					}

					if err := d.Set("failover_group_name", parts[0]); err != nil {
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
		ValidPrivs: validFailoverGroupPrivileges,
	}
}

// CreateFailoverGroupGrant implements schema.CreateFunc.
func CreateFailoverGroupGrant(d *schema.ResourceData, meta interface{}) error {
	failoverGroupName := d.Get("failover_group_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.FailoverGroupGrant(failoverGroupName)

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(failoverGroupName, privilege, withGrantOption, roles)
	d.SetId(grantID)

	return ReadFailoverGroupGrant(d, meta)
}

// ReadFailoverGroupGrant implements schema.ReadFunc.
func ReadFailoverGroupGrant(d *schema.ResourceData, meta interface{}) error {
	failoverGroupName := d.Get("failover_group_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())
	builder := snowflake.FailoverGroupGrant(failoverGroupName)

	err := readGenericGrant(d, meta, failoverGroupGrantSchema, builder, false, false, validFailoverGroupPrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(failoverGroupName, privilege, withGrantOption, roles)
	if d.Id() != grantID {
		d.SetId(grantID)
	}
	return nil
}

// DeleteFailoverGroupGrant implements schema.DeleteFunc.
func DeleteFailoverGroupGrant(d *schema.ResourceData, meta interface{}) error {
	failoverGroupName := d.Get("failover_group_name").(string)
	builder := snowflake.FailoverGroupGrant(failoverGroupName)
	return deleteGenericGrant(d, meta, builder)
}

// UpdateFailoverGroupGrant implements schema.UpdateFunc.
func UpdateFailoverGroupGrant(d *schema.ResourceData, meta interface{}) error {
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
	failoverGroupName := d.Get("failover_group_name").(string)
	privilege := d.Get("privilege").(string)
	reversionRole := d.Get("revert_ownership_to_role_name").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	builder := snowflake.FailoverGroupGrant(failoverGroupName)

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
