package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var validWarehousePrivileges = []string{
	"ALL", "MODIFY", "MONITOR", "OPERATE", "OWNERSHIP", "USAGE",
}

var warehouseGrantSchema = map[string]*schema.Schema{
	"warehouse_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the warehouse on which to grant privileges.",
		ForceNew:    true,
	},
	"privilege": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the warehouse.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validWarehousePrivileges, true),
		ForceNew:     true,
	},
	"roles": &schema.Schema{
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
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
	builder := snowflake.WarehouseGrant(w)

	err := createGenericGrant(data, meta, builder)
	if err != nil {
		return err
	}

	// ID format is <warehouse_name>|||<privilege>
	// data.SetId(fmt.Sprintf("%v|||%v", w, priv))
	dataIdentifiers := make([][]string, 1)
	dataIdentifiers[0] = make([]string, 2)
	dataIdentifiers[0][0] = w
	dataIdentifiers[0][1] = priv
	grantID, err := createGrantID(dataIdentifiers)

	if err != nil {
		return err
	}

	data.SetId(grantID)
	return ReadWarehouseGrant(data, meta)
}

// ReadWarehouseGrant implements schema.ReadFunc
func ReadWarehouseGrant(data *schema.ResourceData, meta interface{}) error {
	// w, _, _, priv, err := splitGrantID(data.Id())
	grantIDArray, err := splitGrantID(data.Id())
	if err != nil {
		return err
	}
	w, priv := grantIDArray[0], grantIDArray[1]

	err = data.Set("warehouse_name", w)
	if err != nil {
		return err
	}
	err = data.Set("privilege", priv)
	if err != nil {
		return err
	}

	builder := snowflake.WarehouseGrant(w)

	return readGenericGrant(data, meta, builder, false)
}

// DeleteWarehouseGrant implements schema.DeleteFunc
func DeleteWarehouseGrant(data *schema.ResourceData, meta interface{}) error {
	// w, _, _, _, err := splitGrantID(data.Id())
	grantIDArray, err := splitGrantID(data.Id())
	if err != nil {
		return err
	}
	w := grantIDArray[0]

	builder := snowflake.WarehouseGrant(w)

	return deleteGenericGrant(data, meta, builder)
}
