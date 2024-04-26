package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var accountPasswordPolicyAttachmentSchema = map[string]*schema.Schema{
	"password_policy": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Qualified name (`\"db\".\"schema\".\"policy_name\"`) of the password policy to apply to the current account.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
}

// AccountPasswordPolicyAttachment returns a pointer to the resource representing an api integration.
func AccountPasswordPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "Specifies the password policy to use for the current account. To set the password policy of a different account, use a provider alias.",

		Create: CreateAccountPasswordPolicyAttachment,
		Read:   ReadAccountPasswordPolicyAttachment,
		Delete: DeleteAccountPasswordPolicyAttachment,

		Schema: accountPasswordPolicyAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateAccountPasswordPolicyAttachment implements schema.CreateFunc.
func CreateAccountPasswordPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()

	passwordPolicy, err := helpers.SafelyDecodeSnowflakeID[sdk.SchemaObjectIdentifier](d.Get("password_policy").(string))
	if err != nil {
		return err
	}

	err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
		Set: &sdk.AccountSet{
			PasswordPolicy: passwordPolicy,
		},
	})
	if err != nil {
		return err
	}

	d.SetId(helpers.EncodeSnowflakeID(passwordPolicy))

	return ReadAccountPasswordPolicyAttachment(d, meta)
}

func ReadAccountPasswordPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	passwordPolicy := helpers.DecodeSnowflakeID(d.Id())
	if err := d.Set("password_policy", passwordPolicy.FullyQualifiedName()); err != nil {
		return err
	}

	return nil
}

// DeleteAccountPasswordPolicyAttachment implements schema.DeleteFunc.
func DeleteAccountPasswordPolicyAttachment(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client
	ctx := context.Background()

	err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
		Unset: &sdk.AccountUnset{
			PasswordPolicy: sdk.Bool(true),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
