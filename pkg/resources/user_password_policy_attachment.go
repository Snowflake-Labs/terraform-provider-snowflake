package resources

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO: what happens if there is a conflict with account policies and user policies?

// TODO: how to return a reference to a struct literal? If I prepend a & it throws an error
var userPasswordPolicyAttachmentSchema = map[string]*schema.Schema{
	"user_name": {
		Type:     schema.TypeString,
		Required: true,
		// TODO: do I need this?
		ForceNew:    true,
		Description: "User name of the user you want to attach the password policy to",
	},
	"password_policy": {
		Type:     schema.TypeString,
		Required: true,
		// TODO: do I need this?
		ForceNew:    true,
		Description: "Qualified name (`\"db\".\"schema\".\"policy_name\"`) of the password policy to apply to the current account.",
	},
}

func UserPasswordPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "Specifies the password policy to use for the current account. To set the password policy of a different account, use a provider alias.",

		Create: CreateUserPasswordPolicyAttachment,
		Read:   ReadUserPasswordPolicyAttachment,
		Delete: DeleteUserPasswordPolicyAttachment,

		Schema: userPasswordPolicyAttachmentSchema,

		// TODO: importer. not sure what it is.
		// 		Importer: &schema.ResourceImporter{
		// 			StateContext: schema.ImportStatePassthroughContext,
		// 		},
	}
}

func CreateUserPasswordPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	// passwordPolicy := d.Get("password_policy").(string)

	// TODO: the password policy is a string, so I comment this for now
	// TODO: I would like to raise an exception if the identifier is not account based
	passwordPolicy, ok := sdk.NewObjectIdentifierFromFullyQualifiedName(d.Get("password_policy").(string)).(sdk.SchemaObjectIdentifier)
	if !ok {
		return fmt.Errorf("password_policy %s is not a valid password policy qualified name, expected format: `\"db\".\"schema\".\"policy\"`", d.Get("password_policy"))
	}

	// TODO: why the following line is commented? I guess indeed we would expect only an accountobjectidentifier, right?
	// passwordPolicy := sdk.NewAccountObjectIdentifier(d.Get("password_policy").(string))

	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("user_name").(string))

	err := client.Users.Alter(ctx, userName, &sdk.AlterUserOptions{
		Set: &sdk.UserSet{
			PasswordPolicy: passwordPolicy,
		},
	})
	if err != nil {
		return err
	}

	d.SetId(helpers.EncodeSnowflakeID(passwordPolicy))

	return nil
}

// TODO: I think this is not correct: this only reads if there is a certain password policy, not if a user has the password policy attached
func ReadUserPasswordPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	passwordPolicy := helpers.DecodeSnowflakeID(d.Id())
	if err := d.Set("password_policy", passwordPolicy.FullyQualifiedName()); err != nil {
		return err
	}

	return nil
}

// DeleteAccountPasswordPolicyAttachment implements schema.DeleteFunc.
func DeleteUserPasswordPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("user_name").(string))

	err := client.Users.Alter(ctx, userName, &sdk.AlterUserOptions{
		Unset: &sdk.UserUnset{
			PasswordPolicy: sdk.Bool(true),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
