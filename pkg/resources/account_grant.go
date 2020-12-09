package resources

import (
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validAccountPrivileges = NewPrivilegeSet(
	privilegeCreateRole,
	privilegeCreateUser,
	privilegeCreateWarehouse,
	privilegeCreateDatabase,
	privilegeCreateIntegration,
	privilegeManageGrants,
	privilegeMonitorUsage,
	privilegeMonitorExecution,
	privilegeExecuteTask,
)

var accountGrantSchema = map[string]*schema.Schema{
	"privilege": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  "The privilege to grant on the account.",
		Default:      privilegeMonitorUsage,
		ValidateFunc: validation.StringInSlice(validAccountPrivileges.ToList(), true),
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

// AccountGrant returns a pointer to the resource representing an account grant
func AccountGrant() *TerraformGrantResource {
	return &TerraformGrantResource{
		Resource: &schema.Resource{
			Create: CreateAccountGrant,
			Read:   ReadAccountGrant,
			Delete: DeleteAccountGrant,

			Schema: accountGrantSchema,
		},
		ValidPrivs: validAccountPrivileges,
	}
}

// CreateAccountGrant implements schema.CreateFunc
func CreateAccountGrant(d *schema.ResourceData, meta interface{}) error {
	priv := d.Get("privilege").(string)
	grantOption := d.Get("with_grant_option").(bool)

	builder := snowflake.AccountGrant()

	err := createGenericGrant(d, meta, builder)
	if err != nil {
		return err
	}

	grantID := &grantID{
		ResourceName: "ACCOUNT",
		Privilege:    priv,
		GrantOption:  grantOption,
	}
	dataIDInput, err := grantID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadAccountGrant(d, meta)
}

// ReadAccountGrant implements schema.ReadFunc
func ReadAccountGrant(d *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(d.Id())
	if err != nil {
		return err
	}
	err = d.Set("privilege", grantID.Privilege)
	if err != nil {
		return err
	}
	err = d.Set("with_grant_option", grantID.GrantOption)
	if err != nil {
		return err
	}

	builder := snowflake.AccountGrant()

	return readGenericGrant(d, meta, accountGrantSchema, builder, false, validAccountPrivileges)
}

// DeleteAccountGrant implements schema.DeleteFunc
func DeleteAccountGrant(d *schema.ResourceData, meta interface{}) error {
	builder := snowflake.AccountGrant()

	return deleteGenericGrant(d, meta, builder)
}
