package resources

import (
	"fmt"
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
}

// WarehouseGrant returns a pointer to the resource representing a warehouse grant.
func WarehouseGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateWarehouseGrant,
			Read:   ReadWarehouseGrant,
			Delete: DeleteWarehouseGrant,
			Update: UpdateWarehouseGrant,

			Schema: warehouseGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
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

	grantID := NewWarehouseGrantID(warehouseName, privilege, roles, withGrantOption)

	d.SetId(grantID.String())

	return ReadWarehouseGrant(d, meta)
}

// ReadWarehouseGrant implements schema.ReadFunc.
func ReadWarehouseGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseWarehouseGrantID(d.Id())
	if err != nil {
		return err
	}

	if !grantID.IsOldID {
		if err := d.Set("roles", grantID.Roles); err != nil {
			return err
		}
	}

	err = d.Set("warehouse_name", grantID.ObjectName)
	if err != nil {
		return err
	}
	err = d.Set("privilege", grantID.Privilege)
	if err != nil {
		return err
	}
	err = d.Set("with_grant_option", grantID.WithGrantOption)
	if err != nil {
		return err
	}

	builder := snowflake.WarehouseGrant(grantID.ObjectName)

	return readGenericGrant(d, meta, warehouseGrantSchema, builder, false, validWarehousePrivileges)
}

// DeleteWarehouseGrant implements schema.DeleteFunc.
func DeleteWarehouseGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := ParseWarehouseGrantID(d.Id())
	if err != nil {
		return err
	}

	builder := snowflake.WarehouseGrant(grantID.ObjectName)

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

	grantID, err := ParseWarehouseGrantID(d.Id())
	if err != nil {
		return err
	}

	// create the builder
	builder := snowflake.WarehouseGrant(grantID.ObjectName)

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
		grantID.WithGrantOption,
		rolesToAdd,
		nil,
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadWarehouseGrant(d, meta)
}

type WarehouseGrantID struct {
	ObjectName      string
	Privilege       string
	Roles           []string
	WithGrantOption bool
	IsOldID         bool
}

func NewWarehouseGrantID(objectName string, privilege string, roles []string, withGrantOption bool) *WarehouseGrantID {
	return &WarehouseGrantID{
		ObjectName:      objectName,
		Privilege:       privilege,
		Roles:           roles,
		WithGrantOption: withGrantOption,
		IsOldID:         false,
	}
}

func (v *WarehouseGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	return fmt.Sprintf("%v|%v|%v|%v", v.ObjectName, v.Privilege, v.WithGrantOption, roles)
}

func ParseWarehouseGrantID(s string) (*WarehouseGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		return &WarehouseGrantID{
			ObjectName:      idParts[0],
			Privilege:       idParts[3],
			Roles:           helpers.SplitStringToSlice(idParts[4], ","),
			WithGrantOption: idParts[5] == "true",
			IsOldID:         true,
		}, nil
	}
	idParts := strings.Split(s, "|")
	if len(idParts) < 4 {
		idParts = strings.Split(s, "❄️") // for that time in 0.56/0.57 when we used ❄️ as a separator
	}
	if len(idParts) != 4 {
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 4", len(idParts))
	}
	return &WarehouseGrantID{
		ObjectName:      idParts[0],
		Privilege:       idParts[1],
		WithGrantOption: idParts[2] == "true",
		Roles:           helpers.SplitStringToSlice(idParts[3], ","),
		IsOldID:         false,
	}, nil
}
