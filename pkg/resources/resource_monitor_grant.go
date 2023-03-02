package resources

import (
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
		Description:  "The privilege to grant on the resource monitor.",
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
			Create: CreateResourceMonitorGrant,
			Read:   ReadResourceMonitorGrant,
			Delete: DeleteResourceMonitorGrant,
			Update: UpdateResourceMonitorGrant,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
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

	grantID := NewResourceMonitorGrantID(monitorName, privilege, roles, withGrantOption)
	d.SetId(grantID.String())

	return ReadResourceMonitorGrant(d, meta)
}

// ReadResourceMonitorGrant implements schema.ReadFunc.
func ReadResourceMonitorGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := parseResourceMonitorGrantID(d.Id())
	if err != nil {
		return err
	}
	if err := d.Set("roles", grantID.Roles); err != nil {
		return err
	}
	if err := d.Set("monitor_name", grantID.ObjectName); err != nil {
		return err
	}
	if err := d.Set("privilege", grantID.Privilege); err != nil {
		return err
	}
	if err := d.Set("with_grant_option", grantID.WithGrantOption); err != nil {
		return err
	}

	builder := snowflake.ResourceMonitorGrant(grantID.ObjectName)
	return readGenericGrant(d, meta, resourceMonitorGrantSchema, builder, false, validResourceMonitorPrivileges)
}

// DeleteResourceMonitorGrant implements schema.DeleteFunc.
func DeleteResourceMonitorGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := parseResourceMonitorGrantID(d.Id())
	if err != nil {
		return err
	}

	builder := snowflake.ResourceMonitorGrant(grantID.ObjectName)

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

	grantID, err := parseResourceMonitorGrantID(d.Id())
	if err != nil {
		return err
	}

	// create the builder
	builder := snowflake.ResourceMonitorGrant(grantID.ObjectName)

	// first revoke
	if err := deleteGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, rolesToRevoke, []string{},
	); err != nil {
		return err
	}
	// then add
	if err := createGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, grantID.WithGrantOption, rolesToAdd, []string{},
	); err != nil {
		return err
	}

	// Done, refresh state
	return ReadResourceMonitorGrant(d, meta)
}

type ResourceMonitorGrantID struct {
	ObjectName      string
	Privilege       string
	Roles           []string
	WithGrantOption bool
	IsOldID         bool
}

func NewResourceMonitorGrantID(objectName string, privilege string, roles []string, withGrantOption bool) *ResourceMonitorGrantID {
	return &ResourceMonitorGrantID{
		ObjectName:      objectName,
		Privilege:       privilege,
		Roles:           roles,
		WithGrantOption: withGrantOption,
		IsOldID:         false,
	}
}

func (v *ResourceMonitorGrantID) String() string {
	roles := strings.Join(v.Roles, ",")
	return fmt.Sprintf("%v|%v|%v|%v", v.ObjectName, v.Privilege, v.WithGrantOption, roles)
}

func parseResourceMonitorGrantID(s string) (*ResourceMonitorGrantID, error) {
	if IsOldGrantID(s) {
		idParts := strings.Split(s, "|")
		return &ResourceMonitorGrantID{
			ObjectName:      idParts[0],
			Privilege:       idParts[3],
			Roles:           helpers.SplitStringToSlice(idParts[4], ","),
			WithGrantOption: idParts[5] == "true",
			IsOldID:         true,
		}, nil
	}
	idParts := helpers.SplitStringToSlice(s, "|")
	if len(idParts) < 4 {
		idParts = helpers.SplitStringToSlice(s, "❄️") // for that time in 0.56/0.57 when we used ❄️ as a separator
	}
	if len(idParts) != 4 {
		return nil, fmt.Errorf("unexpected number of ID parts (%d), expected 4", len(idParts))
	}
	return &ResourceMonitorGrantID{
		ObjectName:      idParts[0],
		Privilege:       idParts[1],
		WithGrantOption: idParts[2] == "true",
		Roles:           helpers.SplitStringToSlice(idParts[3], ","),
		IsOldID:         false,
	}, nil
}
