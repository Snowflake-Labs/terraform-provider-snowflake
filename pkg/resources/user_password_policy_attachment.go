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
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "User name of the user you want to attach the password policy to",
	},
	"password_policy": {
		Type:        schema.TypeString,
		Required:    true,
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

		// TODO: importer, look into it because I am not really sure what is happening here
		// Importer: &schema.ResourceImporter{
		// 	StateContext: schema.ImportStatePassthroughContext,
		// },
	}
}

func CreateUserPasswordPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(d.Get("user_name").(string))
	passwordPolicy := sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(d.Get("password_policy").(string))

	err := client.Users.Alter(ctx, userName, &sdk.AlterUserOptions{
		Set: &sdk.UserSet{
			PasswordPolicy: passwordPolicy,
		},
	})
	if err != nil {
		return err
	}
	if err := d.Set("password_policy", passwordPolicy.FullyQualifiedName()); err != nil {
		return err
	}
	if err := d.Set("user_name", helpers.EncodeSnowflakeID(userName)); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf(`%s|%s`, helpers.EncodeSnowflakeID(passwordPolicy), helpers.EncodeSnowflakeID(userName)))

	return ReadUserPasswordPolicyAttachment(d, meta)
}

// TODO: the client does not incorporate an API to read the view POLICY REFERENCES yet. implement a PolicyReference in client, similar to the function getRowAccessPolicyFor in helpers_test.go
func ReadUserPasswordPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
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
