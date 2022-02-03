package resources

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var validWarehousePrivileges = NewPrivilegeSet(
	privilegeModify,
	privilegeMonitor,
	privilegeOperate,
	privilegeOwnership,
	privilegeUsage,
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
		Description:  "The privilege to grant on the warehouse.",
		Default:      privilegeUsage.String(),
		ForceNew:     true,
		ValidateFunc: validation.ValidatePrivilege(validWarehousePrivileges.ToList(), true),
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
}

// WarehouseGrant returns a pointer to the resource representing a warehouse grant
func WarehouseGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateWarehouseGrant,
			Read:   ReadWarehouseGrant,
			Delete: DeleteWarehouseGrant,
			Update: UpdateWarehouseGrant,

			Schema: warehouseGrantSchema,
			// FIXME - tests for this don't currently work
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validWarehousePrivileges,
	}
}

// CreateWarehouseGrant implements schema.CreateFunc
func CreateWarehouseGrant(d *schema.ResourceData, meta interface{}) error {
	w := d.Get("warehouse_name").(string)
	priv := d.Get("privilege").(string)
	grantOption := d.Get("with_grant_option").(bool)
	builder := snowflake.WarehouseGrant(w)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: w,
		Privilege:    priv,
		GrantOption:  grantOption,
		Roles:        roles,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadWarehouseGrant(d, meta)
}

// ReadWarehouseGrant implements schema.ReadFunc
func ReadWarehouseGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	w := grantID.ResourceName
	priv := grantID.Privilege

	err = d.Set("warehouse_name", w)
	if err != nil {
		return err
	}
	err = d.Set("privilege", priv)
	if err != nil {
		return err
	}
	err = d.Set("with_grant_option", grantID.GrantOption)
	if err != nil {
		return err
	}

	builder := snowflake.WarehouseGrant(w)

	return readGenericGrant(d, meta, warehouseGrantSchema, builder, false, validWarehousePrivileges)
}

// DeleteWarehouseGrant implements schema.DeleteFunc
func DeleteWarehouseGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	w := grantID.ResourceName

	builder := snowflake.WarehouseGrant(w)

	return deleteGenericGrant(d, meta, builder)
}

// UpdateWarehouseGrant implements schema.UpdateFunc
func UpdateWarehouseGrant(d *schema.ResourceData, meta interface{}) error {
	// for now the only thing we can update is roles. if nothing changed,
	// nothing to update and we're done.
	if !d.HasChanges("roles") {
		return nil
	}

	rolesToAdd, rolesToRevoke := changeDiff(d, "roles")

	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}

	// create the builder
	builder := snowflake.WarehouseGrant(grantID.ResourceName)

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta,
		builder,
		grantID.Privilege,
		rolesToRevoke,
		nil,
	); err != nil {
		return err
	}

	// then add
	if err := createGenericGrantRolesAndShares(
		meta,
		builder,
		grantID.Privilege,
		grantID.GrantOption,
		rolesToAdd,
		nil,
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadWarehouseGrant(d, meta)
}
