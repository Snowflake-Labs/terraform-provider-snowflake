package resources

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var validWarehousePrivileges = newPrivilegeSet(
	privilegeAll,
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
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validWarehousePrivileges.toList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
		ForceNew:    true,
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
func WarehouseGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateWarehouseGrant,
		Read:   ReadWarehouseGrant,
		Delete: DeleteWarehouseGrant,

		Schema: warehouseGrantSchema,
		// FIXME - tests for this don't currently work
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// CreateWarehouseGrant implements schema.CreateFunc
func CreateWarehouseGrant(data *schema.ResourceData, meta interface{}) error {
	w := data.Get("warehouse_name").(string)
	priv := data.Get("privilege").(string)
	grantOption := data.Get("with_grant_option").(bool)
	builder := snowflake.WarehouseGrant(w)

	err := createGenericGrant(data, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: w,
		Privilege:    priv,
		GrantOption:  grantOption,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadWarehouseGrant(data, meta)
}

// ReadWarehouseGrant implements schema.ReadFunc
func ReadWarehouseGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	w := grantID.ResourceName
	priv := grantID.Privilege

	err = data.Set("warehouse_name", w)
	if err != nil {
		return err
	}
	err = data.Set("privilege", priv)
	if err != nil {
		return err
	}
	err = data.Set("with_grant_option", grantID.GrantOption)
	if err != nil {
		return err
	}

	builder := snowflake.WarehouseGrant(w)

	return readGenericGrant(data, meta, builder, false, validWarehousePrivileges)
}

// DeleteWarehouseGrant implements schema.DeleteFunc
func DeleteWarehouseGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	w := grantID.ResourceName

	builder := snowflake.WarehouseGrant(w)

	return deleteGenericGrant(data, meta, builder)
}
