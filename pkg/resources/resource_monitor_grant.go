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

var validResourceMonitorPrivileges = NewPrivilegeSet(
	privilegeModify,
	privilegeMonitor,
	privilegeAllPrivileges,
)

var resourceMonitorGrantSchema = map[string]*schema.Schema{
	"monitor_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier for the resource monitor; must be unique for your account.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the resource monitor. To grant all privileges, use the value `ALL PRIVILEGES`",
		Default:      "MONITOR",
		ValidateFunc: validation.StringInSlice(validResourceMonitorPrivileges.ToList(), true),
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
}

// ResourceMonitorGrant returns a pointer to the resource representing a resource monitor grant.
func ResourceMonitorGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create:             CreateResourceMonitorGrant,
			Read:               ReadResourceMonitorGrant,
			Delete:             DeleteResourceMonitorGrant,
			Update:             UpdateResourceMonitorGrant,
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_grant_privileges_to_role instead.",
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 4 {
						return nil, fmt.Errorf("incorrect ID %v: expected monitor_name|privilege|with_grant_option|roles", d.Id())
					}
					if err := d.Set("monitor_name", parts[0]); err != nil {
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
			Schema: resourceMonitorGrantSchema,
		},
		ValidPrivs: validResourceMonitorPrivileges,
	}
}

// CreateResourceMonitorGrant implements schema.CreateFunc.
func CreateResourceMonitorGrant(d *schema.ResourceData, meta interface{}) error {
	monitorName := d.Get("monitor_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	builder := snowflake.ResourceMonitorGrant(monitorName)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	if err := createGenericGrant(d, meta, builder); err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(monitorName, privilege, withGrantOption, roles)
	d.SetId(grantID)

	return ReadResourceMonitorGrant(d, meta)
}

// ReadResourceMonitorGrant implements schema.ReadFunc.
func ReadResourceMonitorGrant(d *schema.ResourceData, meta interface{}) error {
	monitorName := d.Get("monitor_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.ResourceMonitorGrant(monitorName)
	err := readGenericGrant(d, meta, resourceMonitorGrantSchema, builder, false, false, validResourceMonitorPrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(monitorName, privilege, withGrantOption, roles)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// DeleteResourceMonitorGrant implements schema.DeleteFunc.
func DeleteResourceMonitorGrant(d *schema.ResourceData, meta interface{}) error {
	monitorName := d.Get("monitor_name").(string)

	builder := snowflake.ResourceMonitorGrant(monitorName)

	return deleteGenericGrant(d, meta, builder)
}

// UpdateResourceMonitorGrant implements schema.UpdateFunc.
func UpdateResourceMonitorGrant(d *schema.ResourceData, meta interface{}) error {
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

	monitorName := d.Get("monitor_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	// create the builder
	builder := snowflake.ResourceMonitorGrant(monitorName)

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta, builder, privilege, "", rolesToRevoke, []string{},
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
	return ReadResourceMonitorGrant(d, meta)
}
