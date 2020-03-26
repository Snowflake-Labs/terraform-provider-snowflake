package resources

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var validResourceMonitorPrivileges = newPrivilegeSet(
	privilegeAll,
	privilegeModify,
	privilegeMonitor,
)

var resourceMonitorGrantSchema = map[string]*schema.Schema{
	"monitor_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier for the resource monitor; must be unique for your account.",
		ForceNew:    true,
	},
	"privilege": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the resource monitor.",
		Default:      "MONITOR",
		ValidateFunc: validation.StringInSlice(validResourceMonitorPrivileges.toList(), true),
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

// ResourceMonitorGrant returns a pointer to the resource representing a resource monitor grant
func ResourceMonitorGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateResourceMonitorGrant,
		Read:   ReadResourceMonitorGrant,
		Delete: DeleteResourceMonitorGrant,

		Schema: resourceMonitorGrantSchema,
	}
}

// CreateResourceMonitorGrant implements schema.CreateFunc
func CreateResourceMonitorGrant(data *schema.ResourceData, meta interface{}) error {
	w := data.Get("monitor_name").(string)
	priv := data.Get("privilege").(string)
	builder := snowflake.ResourceMonitorGrant(w)

	err := createGenericGrant(data, meta, builder)
	if err != nil {
		return err
	}

	grant := &grantID{
		ResourceName: w,
		Privilege:    priv,
	}
	dataIDInput, err := grant.String()
	if err != nil {
		return err
	}
	data.SetId(dataIDInput)

	return ReadResourceMonitorGrant(data, meta)
}

// ReadResourceMonitorGrant implements schema.ReadFunc
func ReadResourceMonitorGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	w := grantID.ResourceName
	priv := grantID.Privilege

	err = data.Set("monitor_name", w)
	if err != nil {
		return err
	}
	err = data.Set("privilege", priv)
	if err != nil {
		return err
	}

	builder := snowflake.ResourceMonitorGrant(w)

	return readGenericGrant(data, meta, builder, false, validResourceMonitorPrivileges)
}

// DeleteResourceMonitorGrant implements schema.DeleteFunc
func DeleteResourceMonitorGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	w := grantID.ResourceName

	builder := snowflake.ResourceMonitorGrant(w)

	return deleteGenericGrant(data, meta, builder)
}
