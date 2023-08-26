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

var validWarehousePrivileges = NewPrivilegeSet(
	privilegeModify,
	privilegeMonitor,
	privilegeOperate,
	privilegeOwnership,
	privilegeUsage,
	privilegeAllPrivileges,
)

var warehouseGrantSchema = map[string]*schema.Schema{
	"warehouse_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the warehouse on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the warehouse. To grant all privileges, use the value `ALL PRIVILEGES`.",
		Default:      privilegeUsage.String(),
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice(validWarehousePrivileges.ToList(), true),
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

// WarehouseGrant returns a pointer to the resource representing a warehouse grant.
func WarehouseGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create:             CreateWarehouseGrant,
			Read:               ReadWarehouseGrant,
			Delete:             DeleteWarehouseGrant,
			Update:             UpdateWarehouseGrant,
			DeprecationMessage: "This resource is deprecated and will be removed in a future major version release. Please use snowflake_grant_privileges_to_role instead.",
			Schema:             warehouseGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
					parts := strings.Split(d.Id(), helpers.IDDelimiter)
					if len(parts) != 4 {
						return nil, fmt.Errorf("unexpected format of ID (%q), expected warehouse-name|privilege|with_grant_option|roles", d.Id())
					}
					if err := d.Set("warehouse_name", parts[0]); err != nil {
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
		ValidPrivs: validWarehousePrivileges,
	}
}

// CreateWarehouseGrant implements schema.CreateFunc.
func CreateWarehouseGrant(d *schema.ResourceData, meta interface{}) error {
	warehouseName := d.Get("warehouse_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	builder := snowflake.WarehouseGrant(warehouseName)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(warehouseName, privilege, withGrantOption, roles)
	d.SetId(grantID)

	return ReadWarehouseGrant(d, meta)
}

// ReadWarehouseGrant implements schema.ReadFunc.
func ReadWarehouseGrant(d *schema.ResourceData, meta interface{}) error {
	warehouseName := d.Get("warehouse_name").(string)
	privilege := d.Get("privilege").(string)
	withGrantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.WarehouseGrant(warehouseName)

	err := readGenericGrant(d, meta, warehouseGrantSchema, builder, false, false, validWarehousePrivileges)
	if err != nil {
		return err
	}

	grantID := helpers.EncodeSnowflakeID(warehouseName, privilege, withGrantOption, roles)
	if grantID != d.Id() {
		d.SetId(grantID)
	}
	return nil
}

// DeleteWarehouseGrant implements schema.DeleteFunc.
func DeleteWarehouseGrant(d *schema.ResourceData, meta interface{}) error {
	warehouseName := d.Get("warehouse_name").(string)
	builder := snowflake.WarehouseGrant(warehouseName)

	return deleteGenericGrant(d, meta, builder)
}

// UpdateWarehouseGrant implements schema.UpdateFunc.
func UpdateWarehouseGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update is roles. if nothing changed,
	// nothing to update and we're done.
	if !d.HasChanges("roles") {
		return nil
	}

	rolesToAdd, rolesToRevoke := changeDiff(d, "roles")

	warehouseName := d.Get("warehouse_name").(string)
	privilege := d.Get("privilege").(string)
	reversionRole := d.Get("revert_ownership_to_role_name").(string)
	withGrantOption := d.Get("with_grant_option").(bool)

	// create the builder
	builder := snowflake.WarehouseGrant(warehouseName)

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

	// Done, refresh state
	return ReadWarehouseGrant(d, meta)
}
