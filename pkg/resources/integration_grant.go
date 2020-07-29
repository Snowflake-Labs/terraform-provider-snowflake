package resources

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var validIntegrationPrivileges = newPrivilegeSet(
	privilegeAll,
	privilegeUsage,
	privilegeOwnership,
)
var integrationGrantSchema = map[string]*schema.Schema{
	"integration_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier for the integration; must be unique for your account.",
		ForceNew:    true,
	},
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the integration.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validIntegrationPrivileges.toList(), true),
		ForceNew:     true,
	},
	"roles": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Grants privilege to these roles.",
		ForceNew:    true,
	},
}

// IntegrationGrant returns a pointer to the resource representing a integration grant
func IntegrationGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateIntegrationGrant,
		Read:   ReadIntegrationGrant,
		Delete: DeleteIntegrationGrant,

		Schema: integrationGrantSchema,
	}
}

// CreateIntegrationGrant implements schema.CreateFunc
func CreateIntegrationGrant(data *schema.ResourceData, meta interface{}) error {
	w := data.Get("integration_name").(string)
	priv := data.Get("privilege").(string)
	builder := snowflake.IntegrationGrant(w)

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

	return ReadIntegrationGrant(data, meta)
}

// ReadIntegrationGrant implements schema.ReadFunc
func ReadIntegrationGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	w := grantID.ResourceName
	priv := grantID.Privilege

	err = data.Set("integration_name", w)
	if err != nil {
		return err
	}
	err = data.Set("privilege", priv)
	if err != nil {
		return err
	}

	builder := snowflake.IntegrationGrant(w)

	return readGenericGrant(data, meta, builder, false, validIntegrationPrivileges)
}

// DeleteIntegrationGrant implements schema.DeleteFunc
func DeleteIntegrationGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	w := grantID.ResourceName

	builder := snowflake.IntegrationGrant(w)

	return deleteGenericGrant(data, meta, builder)
}
