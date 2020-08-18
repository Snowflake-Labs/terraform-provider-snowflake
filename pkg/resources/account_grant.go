package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/snowflake"
)

var validAccountPrivileges = newPrivilegeSet(
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
		Description:  "The privilege to grant on the schema.",
		Default:      "USAGE",
		ValidateFunc: validation.StringInSlice(validAccountPrivileges.toList(), true),
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

// ViewGrant returns a pointer to the resource representing a view grant
func AccountGrant() *schema.Resource {
	return &schema.Resource{
		Create: CreateAccountGrant,
		Read:   ReadAccountGrant,
		Delete: DeleteAccountGrant,

		Schema: accountGrantSchema,
	}
}

// CreateAccountGrant implements schema.CreateFunc
func CreateAccountGrant(data *schema.ResourceData, meta interface{}) error {
	priv := data.Get("privilege").(string)
	grantOption := data.Get("with_grant_option").(bool)

	builder := snowflake.AccountGrant()

	err := createGenericGrant(data, meta, builder)
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
	data.SetId(dataIDInput)

	return ReadAccountGrant(data, meta)
}

// ReadAccountGrant implements schema.ReadFunc
func ReadAccountGrant(data *schema.ResourceData, meta interface{}) error {
	grantID, err := grantIDFromString(data.Id())
	if err != nil {
		return err
	}
	err = data.Set("privilege", grantID.Privilege)
	if err != nil {
		return err
	}
	err = data.Set("with_grant_option", grantID.GrantOption)
	if err != nil {
		return err
	}

	builder := snowflake.AccountGrant()

	return readGenericGrant(data, meta, builder, false, validAccountPrivileges)
}

// DeleteAccountGrant implements schema.DeleteFunc
func DeleteAccountGrant(data *schema.ResourceData, meta interface{}) error {
	builder := snowflake.AccountGrant()

	return deleteGenericGrant(data, meta, builder)
}
