package resources

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		ValidateFunc: validation.ValidatePrivilege(validResourceMonitorPrivileges.ToList(), true),
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

// ResourceMonitorGrant returns a pointer to the resource representing a resource monitor grant
func ResourceMonitorGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateResourceMonitorGrant,
			Read:   ReadResourceMonitorGrant,
			Delete: DeleteResourceMonitorGrant,

			Schema: resourceMonitorGrantSchema,
		},
		ValidPrivs: validResourceMonitorPrivileges,
	}
}

// CreateResourceMonitorGrant implements schema.CreateFunc
func CreateResourceMonitorGrant(d *schema.ResourceData, meta interface{}) error {
	w := d.Get("monitor_name").(string)
	priv := d.Get("privilege").(string)
	grantOption := d.Get("with_grant_option").(bool)
	builder := snowflake.ResourceMonitorGrant(w)
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

	return ReadResourceMonitorGrant(d, meta)
}

// ReadResourceMonitorGrant implements schema.ReadFunc
func ReadResourceMonitorGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	w := grantID.ResourceName
	priv := grantID.Privilege

	err = d.Set("monitor_name", w)
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

	builder := snowflake.ResourceMonitorGrant(w)
	return readGenericGrant(d, meta, resourceMonitorGrantSchema, builder, false, validResourceMonitorPrivileges)
}

// DeleteResourceMonitorGrant implements schema.DeleteFunc
func DeleteResourceMonitorGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	w := grantID.ResourceName

	builder := snowflake.ResourceMonitorGrant(w)

	return deleteGenericGrant(d, meta, builder)
}
