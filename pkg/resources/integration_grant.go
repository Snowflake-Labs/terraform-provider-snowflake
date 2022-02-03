package resources

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var validIntegrationPrivileges = NewPrivilegeSet(
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
		ValidateFunc: validation.ValidatePrivilege(validIntegrationPrivileges.ToList(), true),
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

// IntegrationGrant returns a pointer to the resource representing a integration grant
func IntegrationGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateIntegrationGrant,
			Read:   ReadIntegrationGrant,
			Delete: DeleteIntegrationGrant,

			Schema: integrationGrantSchema,
			Importer: &schema.ResourceImporter{
				StateContext: schema.ImportStatePassthroughContext,
			},
		},
		ValidPrivs: validIntegrationPrivileges,
	}
}

// CreateIntegrationGrant implements schema.CreateFunc
func CreateIntegrationGrant(d *schema.ResourceData, meta interface{}) error {
	w := d.Get("integration_name").(string)
	priv := d.Get("privilege").(string)
	grantOption := d.Get("with_grant_option").(bool)
	roles := expandStringList(d.Get("roles").(*schema.Set).List())

	builder := snowflake.IntegrationGrant(w)

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

	return ReadIntegrationGrant(d, meta)
}

// ReadIntegrationGrant implements schema.ReadFunc
func ReadIntegrationGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	w := grantID.ResourceName
	priv := grantID.Privilege

	err = d.Set("integration_name", w)
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

	builder := snowflake.IntegrationGrant(w)

	return readGenericGrant(d, meta, integrationGrantSchema, builder, false, validIntegrationPrivileges)
}

// DeleteIntegrationGrant implements schema.DeleteFunc
func DeleteIntegrationGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	w := grantID.ResourceName

	builder := snowflake.IntegrationGrant(w)

	return deleteGenericGrant(d, meta, builder)
}
