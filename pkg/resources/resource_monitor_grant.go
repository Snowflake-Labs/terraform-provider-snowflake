package resources

import (
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

			Schema: resourceMonitorGrantSchema,
		},
		ValidPrivs: validResourceMonitorPrivileges,
	}
}

// CreateResourceMonitorGrant implements schema.CreateFunc.
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

// ReadResourceMonitorGrant implements schema.ReadFunc.
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

// DeleteResourceMonitorGrant implements schema.DeleteFunc.
func DeleteResourceMonitorGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	w := grantID.ResourceName

	builder := snowflake.ResourceMonitorGrant(w)

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

	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}

	w := grantID.ResourceName

	// create the builder
	builder := snowflake.ResourceMonitorGrant(w)

	// first revoke
	err = deleteGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, rolesToRevoke, []string{})
	if err != nil {
		return err
	}
	// then add
	err = createGenericGrantRolesAndShares(
		meta, builder, grantID.Privilege, grantID.GrantOption, rolesToAdd, []string{})
	if err != nil {
		return err
	}

	// Done, refresh state
	return ReadResourceMonitorGrant(d, meta)
}
